package classifier

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/slavonnet/asterisk-response-classifier/internal/config"
)

type mockSTT struct {
	text string
}

func (m mockSTT) Transcribe(_ []byte, _ int, _ []string) (string, float64, error) {
	return m.text, 0.9, nil
}

func TestService_ProcessAudio(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "phrases.yaml")
	yaml := "phrases:\n  positive: [\"да\"]\n  negative: [\"нет\"]\n"
	if err := os.WriteFile(path, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := NewService(config.NewLoader(path), mockSTT{text: "да конечно"})
	r := svc.ProcessAudio([]byte{0xff, 0x7f})
	if r.Label != Positive {
		t.Fatalf("label = %q, want positive", r.Label)
	}
	if r.Transcript != "да конечно" {
		t.Fatalf("transcript = %q", r.Transcript)
	}
}

func TestService_Negative(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "phrases.yaml")
	yaml := "phrases:\n  positive: [\"да\"]\n  negative: [\"нет\"]\n"
	_ = os.WriteFile(path, []byte(yaml), 0o644)

	svc := NewService(config.NewLoader(path), mockSTT{text: "нет не надо"})
	r := svc.ProcessAudio([]byte{0x00})
	if r.Label != Negative {
		t.Fatalf("label = %q, want negative", r.Label)
	}
}

func TestService_UncertainUnknown(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "phrases.yaml")
	yaml := "phrases:\n  positive: [\"да\"]\n  negative: [\"нет\"]\n"
	_ = os.WriteFile(path, []byte(yaml), 0o644)

	svc := NewService(config.NewLoader(path), mockSTT{text: "может быть завтра"})
	r := svc.ProcessAudio([]byte{0x00})
	if r.Label != Uncertain {
		t.Fatalf("label = %q, want uncertain", r.Label)
	}
}
