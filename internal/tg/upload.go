package tg

import (
	"context"
	"fmt"
	"sync"

	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
)

// UploadParams configures a single file upload.
type UploadParams struct {
	Path       string
	OnProgress func(sent, total int64) // optional; nil-safe
}

// progressAdapter implements uploader.Progress. gotd may call Chunk concurrently
// when uploading with multiple threads, so it serializes calls under a mutex and
// forwards a clean (sent, total) stream to the user callback.
type progressAdapter struct {
	mu sync.Mutex
	cb func(sent, total int64)
}

func newProgressAdapter(cb func(sent, total int64)) *progressAdapter {
	return &progressAdapter{cb: cb}
}

func (p *progressAdapter) Chunk(_ context.Context, state uploader.ProgressState) error {
	if p.cb == nil {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cb(state.Uploaded, state.Total)
	return nil
}

// UploadFile uploads a local file in chunks and returns the resulting InputFile
// (small) or InputFileBig (>~10MB). gotd's uploader selects upload.saveFilePart vs
// upload.saveBigFilePart by size internally. Cancel via ctx.
func (c *GotdClient) UploadFile(ctx context.Context, p UploadParams) (tg.InputFileClass, error) {
	api, err := c.acquireAPI()
	if err != nil {
		return nil, err
	}
	u := uploader.NewUploader(api).WithProgress(newProgressAdapter(p.OnProgress))
	f, err := u.FromPath(ctx, p.Path)
	if err != nil {
		return nil, fmt.Errorf("upload %s: %w", p.Path, err)
	}
	return f, nil
}
