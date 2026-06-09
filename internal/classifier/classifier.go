package classifier

import (
	"log"

	"github.com/slavonnet/asterisk-response-classifier/internal/config"
	"github.com/slavonnet/asterisk-response-classifier/internal/phrase"
)

type Classifier interface {
	ProcessAudio(ulaw []byte) Result
}

type Label string

const (
	Positive  Label = "positive"
	Negative  Label = "negative"
	Uncertain Label = "uncertain"
)

type Result struct {
	Text  string
	Score float64
	Label Label
	Heard string // matched phrase (debug)
}

type Service struct {
	loader *config.Loader
	phrase *phrase.Engine
}

func NewService(loader *config.Loader, pe *phrase.Engine) *Service {
	return &Service{loader: loader, phrase: pe}
}

func (s *Service) ProcessAudio(ulaw []byte) Result {
	cfg, err := s.loader.Load()
	if err != nil {
		log.Printf("config: %v", err)
		return uncertain(0, "")
	}
	label, heard, score := s.phrase.Recognize(ulaw, cfg)
	log.Printf("phrase=%q → %s (%.0f)", heard, label, score)
	return Result{
		Text:  label,
		Score: score,
		Label: Label(label),
		Heard: heard,
	}
}

func uncertain(score float64, heard string) Result {
	return Result{Text: "uncertain", Score: score, Label: Uncertain, Heard: heard}
}
