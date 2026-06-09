//go:build phrase

package phrase

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	vosk "github.com/alphacep/vosk-api/go"
	"github.com/slavonnet/asterisk-response-classifier/internal/audio"
	"github.com/slavonnet/asterisk-response-classifier/internal/config"
)

// Engine — phrase-limited recognition (как speech-to-phrase: какая из известных фраз).
type Engine struct {
	modelPath string
	model     *vosk.VoskModel
	mu        sync.Mutex
}

func New(modelPath string) (*Engine, error) {
	vosk.SetLogLevel(-1)
	model, err := vosk.NewModel(modelPath)
	if err != nil {
		return nil, fmt.Errorf("model: %w", err)
	}
	return &Engine{modelPath: modelPath, model: model}, nil
}

func (e *Engine) Recognize(ulaw []byte, cfg *config.Config) (label string, heard string, score float64) {
	phrases, labelOf := cfg.AllPhrases()
	if len(phrases) == 0 || len(ulaw) == 0 {
		return "uncertain", "", 0
	}

	grammar, _ := json.Marshal(phrases)
	pcm := audio.DecodeULaw(ulaw)

	e.mu.Lock()
	defer e.mu.Unlock()

	rec, err := vosk.NewRecognizerGrm(e.model, 8000, string(grammar))
	if err != nil {
		return "uncertain", "", 0
	}
	defer rec.Free()

	_ = rec.AcceptWaveform(pcm)

	var parsed struct {
		Text string `json:"text"`
	}
	_ = json.Unmarshal([]byte(rec.FinalResult()), &parsed)
	heard = strings.TrimSpace(parsed.Text)
	if heard == "" {
		return "uncertain", "", 0
	}

	label = labelOf[norm(heard)]
	if label == "" {
		// fuzzy: substring match
		for p, l := range labelOf {
			if strings.Contains(norm(heard), p) {
				label = l
				break
			}
		}
	}
	if label != "positive" && label != "negative" {
		return "uncertain", heard, 50
	}
	return label, heard, 90
}

func (e *Engine) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.model != nil {
		e.model.Free()
		e.model = nil
	}
}

func norm(s string) string { return strings.ToLower(strings.TrimSpace(s)) }
