package media

import (
	"testing"

	"github.com/sorokin-vladimir/tele/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeMIME(t *testing.T) {
	assert.Equal(t, "image/jpeg", NormalizeMIME("Image/JPEG"))
	assert.Equal(t, "text/plain", NormalizeMIME("text/plain; charset=utf-8"))
	assert.Equal(t, "image/png", NormalizeMIME("  image/png  "))
	assert.Equal(t, "", NormalizeMIME(""))
}

func TestDefaultMediaType(t *testing.T) {
	cases := []struct {
		mime string
		want store.MediaKind
	}{
		{"image/jpeg", store.MediaPhoto},
		{"image/png", store.MediaPhoto},
		{"video/mp4", store.MediaVideo},
		{"video/quicktime", store.MediaVideo},
		{"audio/ogg", store.MediaVoice},
		{"audio/mpeg", store.MediaAudio},
		{"application/pdf", store.MediaFile},
		{"", store.MediaFile},
		{"IMAGE/JPEG; foo=bar", store.MediaPhoto},
	}
	for _, c := range cases {
		assert.Equalf(t, c.want, DefaultMediaType(c.mime), "mime=%q", c.mime)
	}
}
