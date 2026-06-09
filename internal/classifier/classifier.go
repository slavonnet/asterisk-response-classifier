package classifier

import (
	"log"

	"github.com/slavonnet/asterisk-response-classifier/internal/audio"
	"github.com/slavonnet/asterisk-response-classifier/internal/config"
	"github.com/slavonnet/asterisk-response-classifier/internal/wyoming"
)

type Classifier interface {
	ProcessAudio(ulaw []byte) Result
}

type Result struct {
	Text  string
	Score float64
	Label string
	Heard string
}

type Service struct {
	loader *config.Loader
}

func NewService(loader *config.Loader) *Service {
	return &Service{loader: loader}
}

func (s *Service) ProcessAudio(ulaw []byte) Result {
	cfg, err := s.loader.Load()
	if err != nil {
		log.Printf("config: %v", err)
		return Result{Text: "uncertain", Label: "uncertain"}
	}
	if len(ulaw) == 0 {
		return Result{Text: "uncertain", Label: "uncertain"}
	}

	pcm8 := audio.DecodeULaw(ulaw)
	pcm16 := wyoming.Upsample8kTo16k(pcm8)

	heard, err := wyoming.Transcribe(cfg.SpeechToPhrase, pcm16)
	if err != nil {
		log.Printf("speech-to-phrase: %v", err)
		return Result{Text: "uncertain", Label: "uncertain"}
	}
	if heard == "" {
		return Result{Text: "uncertain", Label: "uncertain"}
	}

	label := s.loader.MapPhrase(heard)
	log.Printf("stp=%q → %s", heard, label)
	return Result{Text: label, Label: label, Heard: heard, Score: 90}
}
