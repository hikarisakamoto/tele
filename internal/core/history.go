package core

import (
	"context"

	"go.uber.org/zap"

	"github.com/sorokin-vladimir/tele/internal/core/project"
)

// backfill fetches older history for a chat subscription whose window the store
// could not fill and applies it to state; the registry then emits the resulting
// delta through the same path as any other change. One fetch per subscription is
// in flight at a time (issue #120).
func (o *Owner) backfill(ctx context.Context, id project.SubID, w project.ChatWindow) {
	if !o.beginFetch(id) {
		return
	}
	defer o.endFetch(id)

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
