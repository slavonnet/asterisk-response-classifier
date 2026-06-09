package stt

import "fmt"

// Noop returns empty transcript when STT backend is not configured.
type Noop struct{}

func (Noop) Transcribe(_ []byte, _ int, _ []string) (string, float64, error) {
	return "", 0, fmt.Errorf("stt not configured")
}
