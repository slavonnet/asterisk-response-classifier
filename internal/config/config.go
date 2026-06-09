package config

import (
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SpeechToPhrase string   `yaml:"speech_to_phrase"` // tcp://127.0.0.1:10300
	Language       string   `yaml:"language"`
	Positive       []string `yaml:"positive"`
	Negative       []string `yaml:"negative"`
}

type Loader struct {
	path string
	mu   sync.Mutex
	pos  map[string]bool
	neg  map[string]bool
}

func NewLoader(path string) *Loader {
	return &Loader{path: path}
}

func (l *Loader) Load() (*Config, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	data, err := os.ReadFile(l.path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.SpeechToPhrase == "" {
		cfg.SpeechToPhrase = "tcp://127.0.0.1:10300"
	}
	if cfg.Language == "" {
		cfg.Language = "ru"
	}
	l.pos = toSet(cfg.Positive)
	l.neg = toSet(cfg.Negative)
	return &cfg, nil
}

func (l *Loader) MapPhrase(heard string) string {
	h := norm(heard)
	if l.pos[h] {
		return "positive"
	}
	if l.neg[h] {
		return "negative"
	}
	for p := range l.pos {
		if strings.Contains(h, p) {
			return "positive"
		}
	}
	for n := range l.neg {
		if strings.Contains(h, n) {
			return "negative"
		}
	}
	return "uncertain"
}

func toSet(items []string) map[string]bool {
	m := make(map[string]bool, len(items))
	for _, s := range items {
		m[norm(s)] = true
	}
	return m
}

func norm(s string) string { return strings.ToLower(strings.TrimSpace(s)) }
