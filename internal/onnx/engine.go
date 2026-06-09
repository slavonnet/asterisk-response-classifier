package onnx

import (
	"fmt"
	"os"
	"sync"
)

// Engine wraps ONNX Runtime for phrase-limited STT / classification.
// Requires libonnxruntime.so on Linux (ARM/x86) — see README for setup.
//
// Go binding: github.com/yalue/onnxruntime_go (CGO, links prebuilt .so, no C++ compile).
type Engine struct {
	modelPath string
	ready     bool
	mu        sync.Mutex
}

type Config struct {
	ModelPath          string
	SharedLibraryPath  string // path to libonnxruntime.so
	IntraOpThreads     int    // CPU threads, default 1 for weak machines
	InterOpThreads     int
}

func New(cfg Config) (*Engine, error) {
	if cfg.ModelPath == "" {
		return &Engine{}, nil // disabled
	}
	if _, err := os.Stat(cfg.ModelPath); err != nil {
		return nil, fmt.Errorf("model not found: %s", cfg.ModelPath)
	}
	// TODO: ort.SetSharedLibraryPath + InitializeEnvironment + load session
	return &Engine{modelPath: cfg.ModelPath, ready: false}, nil
}

func (e *Engine) Enabled() bool {
	return e != nil && e.modelPath != ""
}

// Transcribe runs ONNX inference on PCM audio (16-bit LE mono).
func (e *Engine) Transcribe(pcm []byte, sampleRate int) (text string, confidence float64, err error) {
	if !e.Enabled() {
		return "", 0, fmt.Errorf("onnx engine not configured")
	}
	_ = pcm
	_ = sampleRate
	return "", 0, fmt.Errorf("onnx inference not yet implemented")
}
