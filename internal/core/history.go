package core

import (
	"context"

	"go.uber.org/zap"

	"github.com/sorokin-vladimir/tele/internal/core/project"
	"github.com/sorokin-vladimir/tele/internal/domain"
	"github.com/sorokin-vladimir/tele/internal/store"
)

// fill brings a chat window up to date from Telegram, in the order the two
// reasons to fetch deserve. A gap is a defect and is closed first; a window the
// store came up short on is ordinary and is filled after. They share one guard
// rather than racing for it: run separately, whichever started first would
// silently deny the other, and the loser would wait for a scroll that may never
// come.
//
// One fetch per subscription is in flight at a time (issue #120).
func (o *Owner) fill(ctx context.Context, id project.SubID, w project.ChatWindow) {
	if !o.beginFetch(id) {
		return
	}
	defer o.endFetch(id)

	o.forwardFill(ctx, w.ChatID)
	// Re-read the window: the repair may have filled it, and asking again costs
	// a projection build against messages already in memory.
	if needsBackfill(project.BuildChat(o.reader(), w), w) {
		o.backfill(ctx, w)
	}
}

// backfill fetches older history for a chat subscription whose window the store
// could not fill and applies it to state; the registry then emits the resulting
// delta through the same path as any other change. Caller holds the fetch guard.
func (o *Owner) backfill(ctx context.Context, w project.ChatWindow) {
	chat, ok := o.state.Store().GetChat(w.ChatID)
	if !ok {
		return
	}
	held := o.state.Store().Messages(w.ChatID)
	// Page backwards from the oldest message held; zero means "from the newest",
	// which is what an empty window needs.
	offsetID := 0
	if len(held) > 0 {
		offsetID = held[0].ID
	}
	// Read here rather than kept in a field: a changed history limit applies to
	// the next backfill, which is this one.
	fetched, err := o.client.GetHistory(ctx, chat.Peer, offsetID, o.Config().UI.HistoryLimit)
	if err != nil {
		o.log.Warn("history backfill failed", zap.Int64("chat", w.ChatID), zap.Error(err))
		// The client asked for a window it cannot fill itself, so it has to be
		// told: otherwise the pane waits on a load that will never arrive.
		o.publishFailure(Failure{ChatID: w.ChatID, Op: "load history", Err: err})
		return
	}
	// The store merges under its own lock, so a message that arrives mid-fetch
	// is not written back out of existence. Nothing added means the page was
	// entirely history the store already had: the chat has no more to give, and
	// there is no window to rebuild.
	_, added := o.state.MergeHistory(w.ChatID, fetched)
	o.log.Debug("history backfill",
		zap.Int64("chat", w.ChatID),
		zap.Int("offset_id", offsetID),
		zap.Int("held", len(held)),
		zap.Int("fetched", len(fetched)),
		zap.Bool("added", added),
		zap.Int("want", w.Before+w.After+1))
}

// forwardFill closes a chat's recorded gap, fetching from the position it opens
// after towards the tail one page at a time. Each page joins the one before it,
// so the history is contiguous at every step and never holds a hole nothing
// knows about; see docs/adr/0014.
//
// The repair records its progress as it goes, so an interruption costs the page
// in flight rather than the whole run, and the next trigger continues from
// where this one stopped. It gives up on a gap wider than the chat could hold
// anyway and reloads the tail instead: filling that much would only push out
// what it just fetched.
//
// Caller holds the fetch guard.
func (o *Owner) forwardFill(ctx context.Context, chatID int64) {
	st := o.state.Store()
	after, ok := st.Gap(chatID)
	if !ok {
		return
	}
	chat, ok := st.GetChat(chatID)
	if !ok {
		return
	}
	limit := o.Config().UI.HistoryLimit
	// The chat cannot hold more than the cap, so a gap wider than it is one the
	// repair cannot win: every page past that point evicts one it already
	// fetched.
	budget := store.MaxMessagesPerChat
	for {
		page, err := o.client.GetHistoryAfter(ctx, chat.Peer, after, limit)
		if err != nil {
			// Nobody asked for this fetch, so nobody is waiting to hear that it
			// failed: reporting it would interrupt a person over work they did
			// not start. The gap keeps its position and the next open,
			// reconnect or too-long carries on from there.
			o.log.Warn("gap repair failed",
				zap.Int64("chat", chatID), zap.Int("after", after), zap.Error(err))
			return
		}
		newest := 0
		if len(page) > 0 {
			newest = page[len(page)-1].ID
		}
		if newest <= after {
			// Telegram answers a request for messages newer than the newest by
			// sliding the window back over ones we have. Nothing moved forward,
			// so there is nothing left to close.
			o.log.Debug("gap closed", zap.Int64("chat", chatID), zap.Int("at", after))
			st.ClearGap(chatID)
			return
		}
		o.state.RepairHistory(chatID, page)
		after = newest
		st.AdvanceGap(chatID, after)
		o.log.Debug("gap repair page",
			zap.Int64("chat", chatID), zap.Int("fetched", len(page)), zap.Int("now_at", after))

		if chat.TopMessageID > 0 && after >= chat.TopMessageID {
			// Caught up with where the server was when the dialog list was
			// read. Anything newer than that arrived through the update stream.
			st.ClearGap(chatID)
			return
		}
		if budget -= len(page); budget <= 0 {
			o.tailReload(ctx, chat)
			return
		}
	}
}

// tailReload throws a chat's stored history away and takes a fresh page. It is
// the admission that the gap cannot be closed: what is held and what the server
// has cannot be joined, and joining them anyway would draw two ranges as one
// conversation.
//
// It is silent. The scrollback it drops comes back by scrolling, which is
// backfill's ordinary job, and a person who did nothing to cause this has
// nothing to do about it either.
func (o *Owner) tailReload(ctx context.Context, chat domain.Chat) {
	page, err := o.client.GetHistory(ctx, chat.Peer, 0, o.Config().UI.HistoryLimit)
	if err != nil {
		o.log.Warn("tail reload failed", zap.Int64("chat", chat.ID), zap.Error(err))
		return
	}
	if len(page) == 0 {
		return
	}
	o.state.ApplyHistory(chat.ID, page)
	o.state.Store().ClearGap(chat.ID)
	o.log.Info("tail reloaded: the gap was wider than the chat can hold",
		zap.Int64("chat", chat.ID), zap.Int("kept", len(page)))
}

func (o *Owner) beginFetch(id project.SubID) bool {
	o.fetchMu.Lock()
	defer o.fetchMu.Unlock()
	if o.fetching[id] {
		return false
	}
	o.fetching[id] = true
	return true
}

func (o *Owner) endFetch(id project.SubID) {
	o.fetchMu.Lock()
	defer o.fetchMu.Unlock()
	delete(o.fetching, id)
}
