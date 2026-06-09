package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReferences(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "references.yaml")
	content := `
references:
  positive: ["refs/p.ulaw"]
  negative: ["refs/n.ulaw"]
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.References.MinScore != 0.55 {
		t.Fatalf("min_score = %v", cfg.References.MinScore)
	}
	want := filepath.Join(dir, "refs/p.ulaw")
	if cfg.Resolve("refs/p.ulaw") != want {
		t.Fatalf("resolve = %q", cfg.Resolve("refs/p.ulaw"))
	}
}
