// Package audio provides in-app playback of Telegram voice messages: it decodes
// Opus/Ogg to PCM (pure Go, via pion/opus) and plays it through the system
// audio device.
//
// Audio output is cgo-free on every platform, so binaries stay fully static:
//   - macOS / Windows: ebitengine/oto (purego), in sink_other.go.
//   - Linux: the system audio server over the pure-Go PulseAudio/PipeWire
//     protocol (jfreymuth/pulse), in sink_linux.go — avoiding oto's ALSA/cgo
//     dependency.
//
// The Opus/Ogg decode path and the Player/position logic are shared and always
// compiled and tested. When no audio server is reachable, NewPlayer returns an
// error and callers degrade gracefully (no playback).
package audio
