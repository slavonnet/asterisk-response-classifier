package classifier

import (
	"log"
	"strings"

	"github.com/slavonnet/asterisk-response-classifier/internal/audio"
	"github.com/slavonnet/asterisk-response-classifier/internal/config"
	"github.com/slavonnet/asterisk-response-classifier/internal/stt"
)

// Service: STT → keyword classify. Config reloads on every utterance.
type Service struct {
	loader *config.Loader
	stt    stt.Engine
}

func NewService(loader *config.Loader, engine stt.Engine) *Service {
	return &Service{loader: loader, stt: engine}
}

func (s *Service) ProcessAudio(ulaw []byte) Result {
	cfg, err := s.loader.Load()
	if err != nil {
		log.Printf("config reload: %v", err)
		return Result{Text: string(Uncertain), Score: 0, Label: Uncertain}
	}
	kw := NewKeywordClassifier(cfg.Phrases)

	if len(ulaw) == 0 {
		return Result{Text: string(Uncertain), Score: 0, Label: Uncertain}
	}

	pcm := audio.DecodeULaw(ulaw)
	text, conf, err := s.stt.Transcribe(pcm, 8000)
	if err != nil {
		log.Printf("stt: %v", err)
		return Result{Text: string(Uncertain), Score: 0, Label: Uncertain}
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return Result{Text: string(Uncertain), Score: 0, Label: Uncertain, Transcript: ""}
	}

	r := kw.Classify(text)
	r.Transcript = text
	if conf > 0 && r.Score < conf*100 {
		r.Score = conf * 100
	}
	if r.Score < cfg.Phrases.MinConfidence*100 {
		r.Label = Uncertain
		r.Text = string(Uncertain)
	}
	log.Printf("stt=%q label=%s score=%.0f", text, r.Label, r.Score)
	return r
}

// Classify exposes text-only classification (tests).
func (s *Service) Classify(text string) Result {
	cfg, err := s.loader.Load()
	if err != nil {
		return Result{Text: string(Uncertain), Label: Uncertain}
	}
	return NewKeywordClassifier(cfg.Phrases).Classify(text)
}
