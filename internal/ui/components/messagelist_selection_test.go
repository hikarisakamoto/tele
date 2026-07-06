package components_test

import (
	"testing"

	"github.com/sorokin-vladimir/tele/internal/store"
	"github.com/sorokin-vladimir/tele/internal/ui/components"
	"github.com/stretchr/testify/assert"
)

func TestMessageList_SelectedMessageText_ReturnsFocusedText(t *testing.T) {
	ml := components.NewMessageList(20, 40)
	ml.SetMessages(makeMessages(5))
	text, ok := ml.SelectedMessageText()
	assert.True(t, ok)
	assert.Equal(t, "msg 5", text)
}

func TestMessageList_SelectedMessageText_FollowsCursor(t *testing.T) {
	ml := components.NewMessageList(20, 40)
	ml.SetMessages(makeMessages(5))
	ml.CursorUp()
	ml.CursorUp() // cursor on msg 3
	text, ok := ml.SelectedMessageText()
	assert.True(t, ok)
	assert.Equal(t, "msg 3", text)
}

func TestMessageList_SelectedMessageText_EmptyOnMediaOnly(t *testing.T) {
	ml := components.NewMessageList(20, 40)
	ml.SetMessages([]store.Message{{ID: 1, ChatID: 1, Photo: &store.PhotoRef{ID: 42}}})
	text, ok := ml.SelectedMessageText()
	assert.False(t, ok)
	assert.Equal(t, "", text)
}
