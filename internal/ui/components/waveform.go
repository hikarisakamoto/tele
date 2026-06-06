package components

import (
	"fmt"
	"strings"

	"github.com/sorokin-vladimir/tele/internal/store"
)

// waveformBlocks maps an amplitude level (0..7) to a Unicode block glyph.
var waveformBlocks = []rune("▁▂▃▄▅▆▇█")

// maxWaveformBars bounds the inline voice waveform width; longer waveforms are
// downsampled into this many bars.
const maxWaveformBars = 32

// voiceLabel renders a voice message as its waveform plus duration, falling
// back to a plain label when no waveform is present.
func voiceLabel(m *store.MediaRef) string {
	bars := renderWaveform(decodeWaveform(m.Waveform), maxWaveformBars)
	if bars == "" {
		if m.Duration > 0 {
			return "🎤 voice " + formatDuration(m.Duration)
		}
		return "🎤 voice"
	}
	return "🎤 " + bars + " " + formatDuration(m.Duration)
}

// audioLabel renders an audio (music) message as performer/title or filename,
// with a trailing duration when known.
func audioLabel(m *store.MediaRef) string {
	var name string
	switch {
	case m.Title != "" && m.Performer != "":
		name = m.Performer + " — " + m.Title
	case m.Title != "":
		name = m.Title
	case m.FileName != "":
		name = m.FileName
	default:
		name = "audio"
	}
	label := "🎵 " + name
	if m.Duration > 0 {
		label += " " + formatDuration(m.Duration)
	}
	return label
}

// decodeWaveform unpacks Telegram's bitpacked voice waveform: a series of
// little-endian 5-bit amplitude samples (each 0..31). Trailing bits that do
// not form a complete sample are discarded.
func decodeWaveform(packed []byte) []byte {
	if len(packed) == 0 {
		return nil
	}
	count := len(packed) * 8 / 5
	out := make([]byte, count)
	for i := 0; i < count; i++ {
		bit := i * 5
		idx := bit / 8
		shift := uint(bit % 8)
		v := int(packed[idx]) >> shift
		if idx+1 < len(packed) {
			v |= int(packed[idx+1]) << (8 - shift)
		}
		out[i] = byte(v & 0x1F)
	}
	return out
}

// renderWaveform draws amplitude samples as a Unicode block sparkline of the
// given width. When there are more samples than bars, samples are averaged into
// width buckets; otherwise one bar is drawn per sample.
func renderWaveform(samples []byte, width int) string {
	if len(samples) == 0 || width <= 0 {
		return ""
	}
	bars := width
	if len(samples) < bars {
		bars = len(samples)
	}
	var b strings.Builder
	for i := 0; i < bars; i++ {
		lo := i * len(samples) / bars
		hi := (i + 1) * len(samples) / bars
		if hi <= lo {
			hi = lo + 1
		}
		sum := 0
		for j := lo; j < hi; j++ {
			sum += int(samples[j])
		}
		avg := sum / (hi - lo)
		level := avg * (len(waveformBlocks) - 1) / 31
		b.WriteRune(waveformBlocks[level])
	}
	return b.String()
}

// formatDuration renders a media duration as m:ss, or h:mm:ss past an hour.
func formatDuration(seconds int) string {
	if seconds < 0 {
		seconds = 0
	}
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d", m, s)
}
