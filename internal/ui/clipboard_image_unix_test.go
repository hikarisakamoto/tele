//go:build linux || freebsd || openbsd || netbsd || dragonfly

package ui

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClipboardTypesHavePNG(t *testing.T) {
	assert.True(t, clipboardTypesHavePNG("TARGETS\nimage/png\nimage/bmp"))
	assert.True(t, clipboardTypesHavePNG("  image/png  \n"))
	assert.False(t, clipboardTypesHavePNG("text/plain\nUTF8_STRING"))
	assert.False(t, clipboardTypesHavePNG("image/jpeg"))
	assert.False(t, clipboardTypesHavePNG(""))
}

func TestClipboardArgvBuilders(t *testing.T) {
	assert.Equal(t, []string{"wl-paste", "--list-types"}, wlPasteListCmd().Args)
	assert.Equal(t, []string{"wl-paste", "--no-newline", "--type", "image/png"}, wlPasteExtractCmd().Args)
	assert.Equal(t, []string{"xclip", "-selection", "clipboard", "-t", "TARGETS", "-o"}, xclipTargetsCmd().Args)
	assert.Equal(t, []string{"xclip", "-selection", "clipboard", "-t", "image/png", "-o"}, xclipExtractCmd().Args)
}

func TestReadClipboardImagePNG_Success(t *testing.T) {
	list := exec.Command("printf", "TARGETS\nimage/png\n")
	extract := exec.Command("printf", "PNGDATA")
	data, ext, err := readClipboardImagePNG(list, extract)
	require.NoError(t, err)
	assert.Equal(t, []byte("PNGDATA"), data)
	assert.Equal(t, ".png", ext)
}

func TestReadClipboardImagePNG_NoImage_Silent(t *testing.T) {
	list := exec.Command("printf", "text/plain\nUTF8_STRING\n")
	extract := exec.Command("printf", "unused")
	data, ext, err := readClipboardImagePNG(list, extract)
	assert.NoError(t, err)
	assert.Nil(t, data)
	assert.Empty(t, ext)
}

func TestReadClipboardImagePNG_ToolMissing_Silent(t *testing.T) {
	list := exec.Command("tele-no-such-clipboard-tool")
	extract := exec.Command("printf", "unused")
	data, _, err := readClipboardImagePNG(list, extract)
	assert.NoError(t, err, "a missing tool degrades to text, not an error")
	assert.Nil(t, data)
}

func TestReadClipboardImagePNG_ExtractFails_ReturnsError(t *testing.T) {
	list := exec.Command("printf", "image/png\n")
	extract := exec.Command("tele-no-such-clipboard-tool")
	_, _, err := readClipboardImagePNG(list, extract)
	assert.Error(t, err, "an advertised image that fails to extract must report an error (toast)")
}

func TestReadClipboardImagePNG_EmptyExtract_ReturnsError(t *testing.T) {
	list := exec.Command("printf", "image/png\n")
	extract := exec.Command("printf", "")
	_, _, err := readClipboardImagePNG(list, extract)
	assert.Error(t, err)
}
