package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/sorokin-vladimir/tele/internal/ui/components"
)

// maxCaptionChars mirrors components' caption limit for the attach-time warning.
const maxCaptionChars = 1024

// composerLimitText words a composer limit event. The composer reports the event
// structurally and root owns the wording (see Composer.SetPlaceholder).
func composerLimitText(msg components.ComposerLimitMsg) string {
	noun := "Message"
	if msg.Caption {
		noun = "Caption"
	}
	switch msg.Kind {
	case components.ComposerLimitPaste:
		return fmt.Sprintf("Pasted text exceeds the %d-character %s limit — trim it to send",
			msg.Limit, strings.ToLower(noun))
	case components.ComposerLimitOver:
		return fmt.Sprintf("%s is over the %d-character limit — trim it to send", noun, msg.Limit)
	default: // ComposerLimitTyping
		return fmt.Sprintf("%s limit reached — %d characters max", noun, msg.Limit)
	}
}

// handleComposerLimit surfaces a composer limit event as a warning toast. The
// composer already gates these to one per over-limit episode.
func (m RootModel) handleComposerLimit(msg components.ComposerLimitMsg) (RootModel, tea.Cmd) {
	return m.handleStatusErr(StatusErrMsg{
		Sev:  components.SeverityWarning,
		Text: composerLimitText(msg),
	})
}
