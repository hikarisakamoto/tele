package store_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sorokin-vladimir/tele/internal/domain"
	"github.com/sorokin-vladimir/tele/internal/store"
)

func mergeMsg(id int, sec int64) domain.Message {
	return domain.Message{ID: id, ChatID: 7, Date: time.Unix(sec, 0)}
}

func storeWithChat(t *testing.T) store.Store {
	t.Helper()
	s := store.NewMemory()
	s.SetChat(domain.Chat{ID: 7, Peer: domain.Peer{ID: 7, Type: domain.PeerUser}})
	return s
}

func TestMergeMessages_ReportsWhatThePageAdded(t *testing.T) {
	s := storeWithChat(t)
	s.SetMessages(7, []domain.Message{mergeMsg(3, 3), mergeMsg(4, 4)})

	added := s.MergeMessages(7, []domain.Message{mergeMsg(1, 1), mergeMsg(2, 2), mergeMsg(3, 3)})

	assert.Equal(t, 2, added, "the page overlapped the store by one message")
	assert.Equal(t, []int{1, 2, 3, 4}, storedIDs(s, 7))
}

func TestMergeMessages_APageOfWhatIsHeldAddsNothing(t *testing.T) {
	s := storeWithChat(t)
	s.SetMessages(7, []domain.Message{mergeMsg(1, 1), mergeMsg(2, 2)})

	added := s.MergeMessages(7, []domain.Message{mergeMsg(1, 1), mergeMsg(2, 2)})

	assert.Zero(t, added)
	assert.Equal(t, []int{1, 2}, storedIDs(s, 7))
}

// The reason the merge lives in the store at all. A fetch reads the held
// history, spends a round trip on the network and comes back with a page: a
// message that arrived in the meantime is in the store but not in what the
// caller read, and a plain SetMessages would write it back out of existence.
func TestMergeMessages_KeepsAMessageThatArrivedDuringTheFetch(t *testing.T) {
	s := storeWithChat(t)
	s.SetMessages(7, []domain.Message{mergeMsg(5, 5)})

	// What the caller read before it went to the network.
	_ = s.Messages(7)
	// What arrived while it was there.
	s.AppendMessage(mergeMsg(6, 6))
	// What it came back with.
	added := s.MergeMessages(7, []domain.Message{mergeMsg(3, 3), mergeMsg(4, 4)})

	require.Equal(t, 2, added)
	assert.Equal(t, []int{3, 4, 5, 6}, storedIDs(s, 7), "the arrival must survive the page")
}

func TestMergeMessages_TakesTheFetchedCopyOfAMessageEditedInTheGap(t *testing.T) {
	s := storeWithChat(t)
	stale := mergeMsg(1, 1)
	stale.Text = "before"
	s.SetMessages(7, []domain.Message{stale})

	fresh := mergeMsg(1, 1)
	fresh.Text = "after"
	added := s.MergeMessages(7, []domain.Message{fresh})

	assert.Zero(t, added, "an edit changes no count")
	held := s.Messages(7)
	require.Len(t, held, 1)
	assert.Equal(t, "after", held[0].Text)
}

func storedIDs(s store.Store, chatID int64) []int {
	msgs := s.Messages(chatID)
	out := make([]int, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, m.ID)
	}
	return out
}
