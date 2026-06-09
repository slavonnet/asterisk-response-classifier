package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "phrases.yaml")
	content := `
phrases:
  positive: ["да"]
  negative: ["нет"]
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Phrases.Positive) != 1 || cfg.Phrases.Positive[0] != "да" {
		t.Fatalf("positive phrases: %+v", cfg.Phrases.Positive)
	}
	if cfg.Phrases.MinConfidence != 0.55 {
		t.Fatalf("default min_confidence = %v, want 0.55", cfg.Phrases.MinConfidence)
	}
}

func TestLoaderReload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "phrases.yaml")
	if err := os.WriteFile(path, []byte("phrases:\n  positive: [\"да\"]\n  negative: [\"нет\"]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	loader := NewLoader(path)
	cfg1, _ := loader.Load()
	if err := os.WriteFile(path, []byte("phrases:\n  positive: [\"ага\"]\n  negative: [\"нет\"]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg2, _ := loader.Load()
	if cfg1.Phrases.Positive[0] == cfg2.Phrases.Positive[0] {
		t.Fatal("expected reloaded config to differ")
	}
}
