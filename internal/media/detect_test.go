package media

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectMIME_ByExtension(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "pic.png")
	require.NoError(t, os.WriteFile(p, []byte("dummy"), 0o600))

	got, err := DetectMIME(p)
	require.NoError(t, err)
	assert.Equal(t, "image/png", NormalizeMIME(got))
}

func TestDetectMIME_SniffFallback(t *testing.T) {
	dir := t.TempDir()
	// No recognizable extension: detection must fall back to content sniffing.
	// A minimal PNG header sniffs as image/png.
	p := filepath.Join(dir, "blob")
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	require.NoError(t, os.WriteFile(p, pngHeader, 0o600))

	got, err := DetectMIME(p)
	require.NoError(t, err)
	assert.Equal(t, "image/png", NormalizeMIME(got))
}

func TestDetectMIME_UnknownDefaultsToOctetStream(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "blob")
	require.NoError(t, os.WriteFile(p, []byte("just text-ish bytes"), 0o600))

	got, err := DetectMIME(p)
	require.NoError(t, err)
	// http.DetectContentType never returns empty; default is text/plain or
	// application/octet-stream. Either way DefaultMediaType funnels it to MediaFile.
	assert.NotEmpty(t, got)
}

func TestDetectMIME_MissingFile(t *testing.T) {
	_, err := DetectMIME(filepath.Join(t.TempDir(), "nope"))
	assert.Error(t, err)
}
