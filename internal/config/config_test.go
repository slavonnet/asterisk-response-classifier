package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ivr.yaml")
	content := `
speech_to_phrase: "tcp://127.0.0.1:10300"
positive: ["да"]
negative: ["нет"]
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	loader := NewLoader(path)
	cfg, err := loader.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.SpeechToPhrase != "tcp://127.0.0.1:10300" {
		t.Fatalf("addr = %q", cfg.SpeechToPhrase)
	}
	if loader.MapPhrase("да") != "positive" {
		t.Fatal("expected positive")
	}
	if loader.MapPhrase("может") != "uncertain" {
		t.Fatal("expected uncertain")
	}
}
