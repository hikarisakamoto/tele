package ui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	xansi "github.com/charmbracelet/x/ansi"
	kitty "github.com/charmbracelet/x/ansi/kitty"
	"github.com/sorokin-vladimir/tele/internal/ui/components"
	"github.com/stretchr/testify/assert"
)

func TestDimBackground_StripsColorAndAppliesGray(t *testing.T) {
	// A line with a green foreground SGR around visible text.
	colored := "\x1b[32mhello\x1b[0m world"
	out := dimBackground(colored, true)

	// Original green (32) color is gone.
	assert.NotContains(t, out, "[32m")
	// Visible text is preserved.
	assert.Equal(t, "hello world", xansi.Strip(out))
	// A faded gray foreground is applied (256-color index 240 for dark bg).
	assert.Contains(t, out, "240")
}

func TestDimBackground_LightBackgroundUsesLighterGray(t *testing.T) {
	out := dimBackground("plain text", false)
	assert.Contains(t, out, "250")
	assert.Equal(t, "plain text", xansi.Strip(out))
}

func TestDimBackground_LeavesKittyPlaceholderLinesUntouched(t *testing.T) {
	imageLine := "\x1b[38;5;7m" + string(kitty.Placeholder) + string(kitty.Diacritic(0)) + "\x1b[0m"
	out := dimBackground(imageLine, true)
	assert.Equal(t, imageLine, out)
}

func TestDimBackground_PreservesLineCountAndWidth(t *testing.T) {
	in := "\x1b[1;31mline one\x1b[0m\nplain two\n" + string(kitty.Placeholder)
	out := dimBackground(in, true)

	inLines := strings.Split(in, "\n")
	outLines := strings.Split(out, "\n")
	assert.Equal(t, len(inLines), len(outLines))

	// Non-placeholder line widths are preserved.
	assert.Equal(t, lipgloss.Width(inLines[0]), lipgloss.Width(outLines[0]))
	assert.Equal(t, lipgloss.Width(inLines[1]), lipgloss.Width(outLines[1]))
	// Placeholder line is identical.
	assert.Equal(t, inLines[2], outLines[2])
}

func TestAnchorMenu_Outgoing_LeftOfBubble(t *testing.T) {
	area := components.Rect{Top: 1, Left: 10, Height: 20, Width: 40}
	bubble := components.Rect{Top: 5, Left: 30, Height: 3, Width: 18}
	top, left := anchorMenu(bubble, area, 8, 6, true) // onLeft
	assert.Equal(t, 5, top)
	assert.Equal(t, 22, left) // bubble.Left - menuW = 30 - 8
}

func TestAnchorMenu_Incoming_RightOfBubble(t *testing.T) {
	area := components.Rect{Top: 1, Left: 10, Height: 20, Width: 40}
	bubble := components.Rect{Top: 5, Left: 10, Height: 3, Width: 18}
	top, left := anchorMenu(bubble, area, 8, 6, false) // onRight
	assert.Equal(t, 5, top)
	assert.Equal(t, 28, left) // bubble.Left + bubble.Width = 10 + 18
}

func TestAnchorMenu_ClampsLeftEdge(t *testing.T) {
	area := components.Rect{Top: 1, Left: 10, Height: 20, Width: 40}
	bubble := components.Rect{Top: 5, Left: 12, Height: 3, Width: 18} // near left edge
	_, left := anchorMenu(bubble, area, 8, 6, true)                   // 12-8 = 4 < area.Left
	assert.Equal(t, 10, left)                                         // clamped to area.Left
}

func TestAnchorMenu_ClampsBottomEdge(t *testing.T) {
	area := components.Rect{Top: 1, Left: 10, Height: 10, Width: 40} // bottom = 11
	bubble := components.Rect{Top: 9, Left: 12, Height: 3, Width: 18}
	top, _ := anchorMenu(bubble, area, 8, 6, false) // 9+6 = 15 > 11 -> top = 11-6 = 5
	assert.Equal(t, 5, top)
}

func TestAnchorMenu_ClampsRightEdge(t *testing.T) {
	area := components.Rect{Top: 1, Left: 10, Height: 20, Width: 40} // right = 50
	bubble := components.Rect{Top: 5, Left: 30, Height: 3, Width: 18}
	_, left := anchorMenu(bubble, area, 8, 6, false) // 30+18 = 48; 48+8 = 56 > 50 -> 50-8 = 42
	assert.Equal(t, 42, left)
}

func TestOverlayAt_StampsAtPosition(t *testing.T) {
	base := strings.Repeat(strings.Repeat(".", 20)+"\n", 4) + strings.Repeat(".", 20)
	result := overlayAt(base, "AB\nCD", 20, 5, 1, 3)
	lines := strings.Split(result, "\n")
	assert.Equal(t, "...AB...............", lines[1])
	assert.Equal(t, "...CD...............", lines[2])
}

func TestMeasureBox_WidthAndHeight(t *testing.T) {
	w, h := measureBox("ab\ncdef\ng")
	assert.Equal(t, 4, w) // widest line "cdef"
	assert.Equal(t, 3, h)
}

func TestOverlayBottomRight_PlacesOverlayAtBottomRight(t *testing.T) {
	// 20×6 base; overlay is 3 lines × 5 cols
	base := strings.Repeat(strings.Repeat(".", 20)+"\n", 5) + strings.Repeat(".", 20)
	overlay := "AAAAA\nBBBBB\nCCCCC"
	result := overlayBottomRight(base, overlay, 20, 6, 0)
	lines := strings.Split(result, "\n")

	// overlayH=3, overlayW=5
	// top = 6 - 3 - 1 = 2
	// left = 20 - 5 - 2 = 13
	assert.Equal(t, 6, len(lines), "line count preserved")
	assert.Contains(t, lines[2], "AAAAA", "first overlay row at top=2")
	assert.Contains(t, lines[3], "BBBBB")
	assert.Contains(t, lines[4], "CCCCC")
}

func TestOverlayBottomRight_SingleLine(t *testing.T) {
	base := strings.Repeat(strings.Repeat(" ", 10)+"\n", 4) + strings.Repeat(" ", 10)
	overlay := "XYZ"
	result := overlayBottomRight(base, overlay, 10, 5, 0)
	lines := strings.Split(result, "\n")
	// overlayH=1, top = 5-1-1 = 3; overlayW=3, left = 10-3-2 = 5
	assert.Equal(t, 5, len(lines))
	assert.Contains(t, lines[3], "XYZ")
}
