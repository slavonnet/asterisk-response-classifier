package classifier

import (
	"log"
	"os"

	"github.com/slavonnet/asterisk-response-classifier/internal/audio"
	"github.com/slavonnet/asterisk-response-classifier/internal/config"
)

// Classifier maps raw ulaw audio to positive / negative / uncertain.
type Classifier interface {
	ProcessAudio(ulaw []byte) Result
}

// SimilarityClassifier compares audio to reference clips (no STT, no text).
type SimilarityClassifier struct {
	loader *config.Loader
}

func NewSimilarityClassifier(loader *config.Loader) *SimilarityClassifier {
	return &SimilarityClassifier{loader: loader}
}

type refSet struct {
	positive []audio.FeatureVector
	negative []audio.FeatureVector
	minScore float64
	minMargin float64
}

func (s *SimilarityClassifier) ProcessAudio(ulaw []byte) Result {
	cfg, err := s.loader.Load()
	if err != nil {
		log.Printf("config: %v", err)
		return uncertain(0)
	}
	refs, err := s.loadRefs(cfg)
	if err != nil {
		log.Printf("refs: %v", err)
		return uncertain(0)
	}
	if len(ulaw) == 0 {
		return uncertain(0)
	}

	feat := audio.ExtractFeatures(audio.DecodeULaw(ulaw))
	if feat == nil {
		return uncertain(0)
	}

	posBest := bestMatch(feat, refs.positive)
	negBest := bestMatch(feat, refs.negative)

	if len(refs.positive) == 0 || len(refs.negative) == 0 {
		log.Printf("need both positive and negative reference clips")
		return uncertain(0)
	}

	var label Label
	var score float64
	margin := posBest - negBest

	switch {
	case posBest >= negBest && posBest >= refs.minScore && margin >= refs.minMargin:
		label = Positive
		score = posBest * 100
	case negBest > posBest && negBest >= refs.minScore && -margin >= refs.minMargin:
		label = Negative
		score = negBest * 100
	default:
		return uncertain(mathMax(posBest, negBest) * 100)
	}

	log.Printf("audio sim pos=%.2f neg=%.2f margin=%.2f → %s", posBest, negBest, margin, label)
	return Result{Text: string(label), Score: score, Label: label}
}

func (s *SimilarityClassifier) loadRefs(cfg *config.Config) (*refSet, error) {
	rs := &refSet{
		minScore:  cfg.References.MinScore,
		minMargin: cfg.References.MinMargin,
	}
	for _, p := range cfg.References.Positive {
		f, err := loadRefFeature(cfg.Resolve(p))
		if err != nil {
			return nil, err
		}
		rs.positive = append(rs.positive, f)
	}
	for _, n := range cfg.References.Negative {
		f, err := loadRefFeature(cfg.Resolve(n))
		if err != nil {
			return nil, err
		}
		rs.negative = append(rs.negative, f)
	}
	return rs, nil
}

func loadRefFeature(path string) (audio.FeatureVector, error) {
	ulaw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f := audio.ExtractFeatures(audio.DecodeULaw(ulaw))
	if f == nil {
		return nil, os.ErrInvalid
	}
	return f, nil
}

func bestMatch(f audio.FeatureVector, refs []audio.FeatureVector) float64 {
	best := 0.0
	for _, r := range refs {
		if sim := audio.CosineSimilarity(f, r); sim > best {
			best = sim
		}
	}
	return best
}

func uncertain(score float64) Result {
	return Result{Text: string(Uncertain), Score: score, Label: Uncertain}
}

func mathMax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
