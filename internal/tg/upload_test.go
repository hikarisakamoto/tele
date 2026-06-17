package tg

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/gotd/td/telegram/uploader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProgressAdapter_ForwardsSentAndTotal(t *testing.T) {
	var lastSent, lastTotal int64
	a := newProgressAdapter(func(sent, total int64) {
		lastSent, lastTotal = sent, total
	})
	require.NoError(t, a.Chunk(context.Background(), uploader.ProgressState{Uploaded: 50, Total: 200}))
	assert.Equal(t, int64(50), lastSent)
	assert.Equal(t, int64(200), lastTotal)
}

func TestProgressAdapter_NilCallbackIsSafe(t *testing.T) {
	a := newProgressAdapter(nil)
	require.NoError(t, a.Chunk(context.Background(), uploader.ProgressState{Uploaded: 1, Total: 2}))
}

func TestProgressAdapter_ConcurrentChunksSerialized(t *testing.T) {
	var calls int64
	a := newProgressAdapter(func(sent, total int64) {
		atomic.AddInt64(&calls, 1)
	})
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = a.Chunk(context.Background(), uploader.ProgressState{Uploaded: int64(n), Total: 100})
		}(i)
	}
	wg.Wait()
	assert.Equal(t, int64(100), atomic.LoadInt64(&calls))
}
