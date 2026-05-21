package components_test

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/sorokin-vladimir/tele/internal/ui/components"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComposer_Value(t *testing.T) {
	c := components.NewComposer(60)
	c.SetValue("hello")
	assert.Equal(t, "hello", c.Value())
}

func TestComposer_Reset(t *testing.T) {
	c := components.NewComposer(60)
	c.SetValue("text")
	c.Reset()
	assert.Empty(t, c.Value())
}

func TestComposer_View_NotEmpty(t *testing.T) {
	c := components.NewComposer(60)
	assert.NotEmpty(t, c.View())
}

func TestComposer_Update_AcceptsInput(t *testing.T) {
	c := components.NewComposer(60)
	c.Focus()
	newC, _ := c.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	assert.Equal(t, "a", newC.Value())
}

func TestComposer_VisualHeight_Empty(t *testing.T) {
	c := components.NewComposer(60)
	assert.Equal(t, 3, c.VisualHeight()) // 1 textarea line + 2 border
}

func TestComposer_VisualHeight_WithPreview(t *testing.T) {
	c := components.NewComposer(60)
	c.SetReplyPreview("↩ Reply: some message")
	assert.Equal(t, 5, c.VisualHeight()) // 1 textarea + 2 border + 1 preview + 1 blank sep
}

func TestComposer_ClearReplyPreview(t *testing.T) {
	c := components.NewComposer(60)
	c.SetReplyPreview("something")
	c.ClearReplyPreview()
	assert.Equal(t, 3, c.VisualHeight())
}

func TestComposer_Reset_ClearsPreview(t *testing.T) {
	c := components.NewComposer(60)
	c.SetValue("hello")
	c.SetReplyPreview("something")
	c.Reset()
	assert.Empty(t, c.Value())
	assert.Equal(t, 3, c.VisualHeight())
}

func TestComposer_VisualHeight_WithTwoLinePreview(t *testing.T) {
	c := components.NewComposer(60)
	c.SetReplyPreview("line1\nline2")
	assert.Equal(t, 6, c.VisualHeight()) // 1 textarea + 2 border + 2 preview lines + 1 blank sep
}

func TestComposer_View_ReplyPreviewBlankSeparator(t *testing.T) {
	c := components.NewComposer(60)
	c.SetReplyPreview("↩ Reply: hello")
	view := stripANSI(c.View())
	lines := strings.Split(view, "\n")

	previewIdx := -1
	for i, l := range lines {
		if strings.Contains(l, "↩ Reply: hello") {
			previewIdx = i
			break
		}
	}
	require.GreaterOrEqual(t, previewIdx, 0, "preview line not found")
	require.Less(t, previewIdx+1, len(lines), "no line after preview")

	// The line immediately after the preview must be blank inside the border.
	lineAfter := lines[previewIdx+1]
	trimmed := strings.TrimSpace(strings.Trim(lineAfter, "│ "))
	assert.Empty(t, trimmed, "expected blank line after reply preview in composer")
}
