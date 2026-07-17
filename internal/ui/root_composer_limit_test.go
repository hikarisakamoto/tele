package ui

import (
	"testing"

	"github.com/sorokin-vladimir/tele/internal/ui/components"
	"github.com/stretchr/testify/assert"
)

func TestComposerLimitText(t *testing.T) {
	for _, tc := range []struct {
		name string
		msg  components.ComposerLimitMsg
		want string
	}{
		{
			name: "typing at the message limit",
			msg:  components.ComposerLimitMsg{Kind: components.ComposerLimitTyping, Limit: 4096},
			want: "Message limit reached — 4096 characters max",
		},
		{
			name: "typing at the caption limit",
			msg:  components.ComposerLimitMsg{Kind: components.ComposerLimitTyping, Limit: 1024, Caption: true},
			want: "Caption limit reached — 1024 characters max",
		},
		{
			name: "paste over the message limit",
			msg:  components.ComposerLimitMsg{Kind: components.ComposerLimitPaste, Limit: 4096},
			want: "Pasted text exceeds the 4096-character message limit — trim it to send",
		},
		{
			name: "draft already over the caption limit",
			msg:  components.ComposerLimitMsg{Kind: components.ComposerLimitOver, Limit: 1024, Caption: true},
			want: "Caption is over the 1024-character limit — trim it to send",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, composerLimitText(tc.msg))
		})
	}
}
