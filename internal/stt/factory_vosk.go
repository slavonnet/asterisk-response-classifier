//go:build vosk

package stt

import "fmt"

func NewEngine(modelPath string) (Engine, error) {
	if modelPath == "" {
		return Noop{}, nil
	}
	return NewVoskLocal(modelPath)
}
