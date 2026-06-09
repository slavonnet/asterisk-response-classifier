//go:build vosk

package stt

func NewEngine(modelPath string) (Engine, error) {
	if modelPath == "" {
		return Noop{}, nil
	}
	return NewVoskLocal(modelPath)
}
