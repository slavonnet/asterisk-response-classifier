package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAllPhrases(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sentences.yaml")
	content := `language: ru
lists:
  yes_word:
    values:
      - in: "да"
  no_word:
    values:
      - in: "нет"
intents:
  Positive:
    data:
      - sentences: ["{yes_word}"]
  Negative:
    data:
      - sentences: ["{no_word}"]
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	phrases, labels := cfg.AllPhrases()
	if len(phrases) != 2 {
		t.Fatalf("phrases = %v", phrases)
	}
	if labels["да"] != "positive" || labels["нет"] != "negative" {
		t.Fatalf("labels = %v", labels)
	}
}
