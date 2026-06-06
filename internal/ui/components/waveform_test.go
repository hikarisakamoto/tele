package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeWaveform_Unpacks5BitSamples(t *testing.T) {
	// Samples [1, 2, 3] packed LSB-first: 1 | 2<<5 | 3<<10 = 0x0C41 -> LE {0x41, 0x0C}.
	got := decodeWaveform([]byte{0x41, 0x0C})
	assert.Equal(t, []byte{1, 2, 3}, got)
}

func TestDecodeWaveform_Empty(t *testing.T) {
	assert.Empty(t, decodeWaveform(nil))
}

func TestRenderWaveform_MapsAmplitudeToBlocks(t *testing.T) {
	// Min sample -> lowest block, max (31) -> full block.
	assert.Equal(t, "▁█", renderWaveform([]byte{0, 31}, 2))
}

func TestRenderWaveform_DownsamplesToWidth(t *testing.T) {
	// Four samples bucketed into two bars: avg(0,0)=low, avg(31,31)=full.
	assert.Equal(t, "▁█", renderWaveform([]byte{0, 0, 31, 31}, 2))
}

func TestRenderWaveform_EmptyIsEmpty(t *testing.T) {
	assert.Equal(t, "", renderWaveform(nil, 8))
}

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		secs int
		want string
	}{
		{0, "0:00"},
		{15, "0:15"},
		{75, "1:15"},
		{200, "3:20"},
		{3661, "1:01:01"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, formatDuration(tc.secs))
	}
}
