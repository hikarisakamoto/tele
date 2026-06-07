package ui

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenInViewer_CreatesFileInTmpDir(t *testing.T) {
	tmpDir := t.TempDir()

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})

	openInViewer(img, tmpDir)

	entries, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	require.Len(t, entries, 1, "expected exactly one temp file in tmpDir")

	filePath := filepath.Join(tmpDir, entries[0].Name())
	info, err := os.Stat(filePath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

func TestWriteTempMediaFile_WritesBytesWithExtension(t *testing.T) {
	dir := t.TempDir()
	data := []byte("opus-bytes")

	path, err := writeTempMediaFile(data, dir, ".ogg")
	require.NoError(t, err)

	assert.Equal(t, ".ogg", filepath.Ext(path))
	got, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, data, got)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}
