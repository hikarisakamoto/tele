package audio

import (
	"sync"
	"time"
)

// sink is the platform audio output for a single PCM stream. Implementations
// (oto on macOS/Windows, PulseAudio on Linux) are selected by build tag; both
// are cgo-free. newSink and checkBackend are provided per platform.
type sink interface {
	pause()
	resume()
	stop()
}

type current struct {
	snk       sink
	docID     int64
	duration  time.Duration
	start     time.Time     // reference time playback began
	accPaused time.Duration // total time spent paused
	pausedAt  time.Time     // when the current pause began
	paused    bool
}

// Player plays a single voice message at a time through the system audio device.
// Playback position is tracked by wall clock, which is accurate enough for the
// waveform playhead and works identically across backends.
type Player struct {
	mu  sync.Mutex
	cur *current
}

// NewPlayer returns a player, failing when no audio backend is available so
// callers can degrade gracefully (no voice playback).
func NewPlayer() (*Player, error) {
	if err := checkBackend(); err != nil {
		return nil, err
	}
	return &Player{}, nil
}

// Play decodes and starts playing the given voice document, replacing any
// current playback.
func (p *Player) Play(docID int64, ogg []byte) error {
	pcm, err := decodeVoicePCM(ogg)
	if err != nil {
		return err
	}
	p.Stop()

	snk, err := newSink(pcm)
	if err != nil {
		return err
	}
	dur := time.Duration(len(pcm)) * time.Second / time.Duration(bytesPerSecond)

	p.mu.Lock()
	p.cur = &current{snk: snk, docID: docID, duration: dur, start: time.Now()}
	p.mu.Unlock()
	return nil
}

// Toggle pauses or resumes the active playback. Returns true if docID is the one
// playing (so callers can distinguish toggle from starting a new message).
func (p *Player) Toggle(docID int64) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cur == nil || p.cur.docID != docID {
		return false
	}
	if p.cur.paused {
		p.cur.accPaused += time.Since(p.cur.pausedAt)
		p.cur.paused = false
		p.cur.snk.resume()
	} else {
		p.cur.pausedAt = time.Now()
		p.cur.paused = true
		p.cur.snk.pause()
	}
	return true
}

// Stop halts and clears any playback.
func (p *Player) Stop() {
	p.mu.Lock()
	cur := p.cur
	p.cur = nil
	p.mu.Unlock()
	if cur != nil {
		cur.snk.stop()
	}
}

// State reports the current playback. active is false when nothing is playing or
// the message has finished (which also clears it).
func (p *Player) State() (docID int64, progress float64, posSecs int, active bool) {
	p.mu.Lock()
	cur := p.cur
	p.mu.Unlock()
	if cur == nil {
		return 0, 0, 0, false
	}

	elapsed := time.Since(cur.start) - cur.accPaused
	if cur.paused {
		elapsed = cur.pausedAt.Sub(cur.start) - cur.accPaused
	}
	if elapsed < 0 {
		elapsed = 0
	}

	if elapsed >= cur.duration {
		p.Stop()
		return 0, 0, 0, false
	}

	if cur.duration > 0 {
		progress = float64(elapsed) / float64(cur.duration)
	}
	return cur.docID, progress, int(elapsed / time.Second), true
}
