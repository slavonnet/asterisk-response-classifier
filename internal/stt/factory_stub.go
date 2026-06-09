//go:build !vosk

package stt

import "fmt"

// NewEngine creates the STT backend. With -tags vosk links libvosk in-process.
func NewEngine(modelPath string) (Engine, error) {
	if modelPath == "" {
		return Noop{}, nil
	}
	return nil, fmt.Errorf("STT requires build with -tags vosk and libvosk installed; see README")
}
