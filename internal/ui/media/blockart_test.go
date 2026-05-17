package media_test

import (
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/sorokin-vladimir/tele/internal/ui/media"
	"github.com/stretchr/testify/require"
)

func TestRenderBlockArt_OutputLines(t *testing.T) {
	// 4x4 image → after scaling to 4 cols, ~8 rows → 4 terminal lines
	// But we just need at least 1 line returned
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x * 50), G: uint8(y * 50), B: 100, A: 255})
		}
	}
	lines := media.RenderBlockArt(img, 4)
	require.NotEmpty(t, lines, "should produce at least one line")
	for _, l := range lines {
		stripped := stripANSI(l)
		require.Equal(t, 4, len([]rune(stripped)), "each line should have 4 half-block chars")
	}
}

func TestRenderBlockArt_EmptyImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))
	lines := media.RenderBlockArt(img, 4)
	require.Nil(t, lines)
}

func TestRenderBlockArt_ZeroCols(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	lines := media.RenderBlockArt(img, 0)
	require.Nil(t, lines)
}

func TestPhotoTermLines_Basic(t *testing.T) {
	// 100x100 image at 10 cols: targetH = 10*100/100 = 10, termRows = 5
	result := media.PhotoTermLines(100, 100, 10)
	require.Equal(t, 5, result)
}

func TestPhotoTermLines_ZeroWidth(t *testing.T) {
	result := media.PhotoTermLines(0, 100, 10)
	require.Equal(t, 1, result)
}

// stripANSI removes ANSI escape sequences for visual width checking.
func stripANSI(s string) string {
	var out strings.Builder
	skip := false
	for _, r := range s {
		if r == '\x1b' {
			skip = true
			continue
		}
		if skip {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				skip = false
			}
			continue
		}
		out.WriteRune(r)
	}
	return out.String()
}
