package classifier

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/slavonnet/asterisk-response-classifier/internal/audio"
	"github.com/slavonnet/asterisk-response-classifier/internal/config"
)

func writeRef(t *testing.T, dir, name string, samples []int16) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, audio.EncodeULaw(samples), 0o644); err != nil {
		t.Fatal(err)
	}
}

func synthPositive() []int16 {
	s := make([]int16, 1200)
	for i := range s {
		s[i] = int16(3000 * (i % 40) / 40)
	}
	return s
}

func synthNegative() []int16 {
	s := make([]int16, 3500)
	for i := range s {
		s[i] = int16(2000 * ((i * 3) % 100) / 100)
	}
	return s
}

func setupTestClassifier(t *testing.T) *SimilarityClassifier {
	t.Helper()
	dir := t.TempDir()
	writeRef(t, dir, "refs/p.ulaw", synthPositive())
	writeRef(t, dir, "refs/n.ulaw", synthNegative())

	yaml := `references:
  min_score: 0.5
  min_margin: 0.05
  positive: ["refs/p.ulaw"]
  negative: ["refs/n.ulaw"]
`
	cfgPath := filepath.Join(dir, "references.yaml")
	if err := os.WriteFile(cfgPath, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	return NewSimilarityClassifier(config.NewLoader(cfgPath))
}

func TestSimilarity_Positive(t *testing.T) {
	cls := setupTestClassifier(t)
	r := cls.ProcessAudio(audio.EncodeULaw(synthPositive()))
	if r.Label != Positive {
		t.Fatalf("label = %q, score=%.0f", r.Label, r.Score)
	}
}

func TestSimilarity_Negative(t *testing.T) {
	cls := setupTestClassifier(t)
	r := cls.ProcessAudio(audio.EncodeULaw(synthNegative()))
	if r.Label != Negative {
		t.Fatalf("label = %q, score=%.0f", r.Label, r.Score)
	}
}

func TestSimilarity_UncertainEmpty(t *testing.T) {
	cls := setupTestClassifier(t)
	r := cls.ProcessAudio(nil)
	if r.Label != Uncertain {
		t.Fatalf("label = %q", r.Label)
	}
}
