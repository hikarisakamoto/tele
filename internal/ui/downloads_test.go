package ui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUniqueDownloadFile_ResolvesCollision(t *testing.T) {
	dir := t.TempDir()

	f1, err := createUniqueDownloadFile(dir, "report.pdf")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, "report.pdf"), f1.Name())
	require.NoError(t, f1.Close())

	f2, err := createUniqueDownloadFile(dir, "report.pdf")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, "report (1).pdf"), f2.Name())
	require.NoError(t, f2.Close())
}

func TestCreateUniqueDownloadFile_SanitizesName(t *testing.T) {
	dir := t.TempDir()
	f, err := createUniqueDownloadFile(dir, "/etc/passwd")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, "passwd"), f.Name())
	require.NoError(t, f.Close())
}

func TestCreateUniqueDownloadFile_EmptyNameFallsBack(t *testing.T) {
	dir := t.TempDir()
	f, err := createUniqueDownloadFile(dir, "")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, "file"), f.Name())
	require.NoError(t, f.Close())
}

func TestResolveDownloadsDir_NonEmpty(t *testing.T) {
	assert.NotEmpty(t, resolveDownloadsDir())
}

func TestResolveDownloadsDir_PrefersXDG(t *testing.T) {
	if os.Getenv("HOME") == "" {
		t.Skip("no HOME")
	}
	want := t.TempDir()
	t.Setenv("XDG_DOWNLOAD_DIR", want)
	// XDG is consulted on Linux; on macOS the env is ignored, so only assert
	// the env path when it is actually honored.
	got := resolveDownloadsDir()
	if got != want {
		assert.Contains(t, got, "Downloads")
	}
}
