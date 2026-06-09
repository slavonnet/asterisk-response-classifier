package classifier

import (
	"strings"

	"github.com/slavonnet/asterisk-response-classifier/internal/config"
)

// KeywordClassifier maps known phrases to labels without ML (MVP / fallback).
type KeywordClassifier struct {
	positive  []string
	negative  []string
	threshold float64
}

func NewKeywordClassifier(p config.PhraseGroups) *KeywordClassifier {
	return &KeywordClassifier{
		positive:  normalizeList(p.Positive),
		negative:  normalizeList(p.Negative),
		threshold: p.MinConfidence,
	}
}

func normalizeList(items []string) []string {
	out := make([]string, len(items))
	for i, s := range items {
		out[i] = strings.ToLower(strings.TrimSpace(s))
	}
	return out
}

func (k *KeywordClassifier) Classify(text string) Result {
	t := strings.ToLower(strings.TrimSpace(text))
	if t == "" {
		return Result{Text: string(Uncertain), Score: 0, Label: Uncertain}
	}

	for _, p := range k.positive {
		if containsPhrase(t, p) {
			return Result{Text: string(Positive), Score: 90, Label: Positive}
		}
	}
	for _, n := range k.negative {
		if containsPhrase(t, n) {
			return Result{Text: string(Negative), Score: 90, Label: Negative}
		}
	}
	return Result{Text: string(Uncertain), Score: 30, Label: Uncertain}
}

func containsPhrase(text, phrase string) bool {
	return strings.Contains(text, phrase)
}
