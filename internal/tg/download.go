package tg

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/gotd/td/telegram/downloader"
	gotdtg "github.com/gotd/td/tg"

	"github.com/sorokin-vladimir/tele/internal/store"
)

func (c *GotdClient) DownloadPhoto(ctx context.Context, ref store.PhotoRef) (image.Image, error) {
	c.mu.RLock()
	api := c.api
	c.mu.RUnlock()
	if api == nil {
		return nil, fmt.Errorf("not connected")
	}

	loc := &gotdtg.InputPhotoFileLocation{
		ID:            ref.ID,
		AccessHash:    ref.AccessHash,
		FileReference: ref.FileReference,
		ThumbSize:     ref.ThumbSize,
	}

	var buf bytes.Buffer
	d := downloader.NewDownloader()
	if _, err := d.Download(api, loc).Stream(ctx, &buf); err != nil {
		return nil, fmt.Errorf("download photo %d: %w", ref.ID, err)
	}

	img, _, err := image.Decode(&buf)
	if err != nil {
		return nil, fmt.Errorf("decode photo %d: %w", ref.ID, err)
	}
	return img, nil
}
