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
tree:
  start: greeting
  nodes:
    - id: greeting
      prompt: test
      on:
        positive: next
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
	if cfg.Phrases.MinConfidence != 0.6 {
		t.Fatalf("default min_confidence = %v, want 0.6", cfg.Phrases.MinConfidence)
	}
	if cfg.Tree.Start != "greeting" {
		t.Fatalf("tree.start = %q", cfg.Tree.Start)
	}
}
