//go:build vosk

package stt

import (
	"encoding/json"
	"fmt"
	"sync"

	vosk "github.com/alphacep/vosk-api/go"
)

// VoskLocal: base acoustic model + phrase grammar rebuilt from phrases.yaml each utterance.
// Same idea as speech-to-phrase (limited vocabulary), not open-ended STT.
type VoskLocal struct {
	model *vosk.VoskModel
	mu    sync.Mutex
}

func NewVoskLocal(modelPath string) (*VoskLocal, error) {
	vosk.SetLogLevel(-1)
	model, err := vosk.NewModel(modelPath)
	if err != nil {
		return nil, fmt.Errorf("vosk model %q: %w", modelPath, err)
	}
	return &VoskLocal{model: model}, nil
}

func (v *VoskLocal) Transcribe(pcm []byte, sampleRate int, phrases []string) (string, float64, error) {
	if len(pcm) == 0 {
		return "", 0, nil
	}
	if len(phrases) == 0 {
		return "", 0, fmt.Errorf("empty phrase list")
	}

	grammar, err := json.Marshal(phrases)
	if err != nil {
		return "", 0, err
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	rec, err := vosk.NewRecognizerGrm(v.model, float64(sampleRate), string(grammar))
	if err != nil {
		return "", 0, err
	}
	defer rec.Free()

	_ = rec.AcceptWaveform(pcm)

	var parsed struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(rec.FinalResult()), &parsed); err != nil {
		return "", 0, err
	}
	text := parsed.Text
	if text == "" {
		return "", 0, nil
	}
	return text, 0.9, nil
}

func (v *VoskLocal) Close() {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.model != nil {
		v.model.Free()
		v.model = nil
	}
}
