package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type PhraseGroups struct {
	Positive      []string `yaml:"positive"`
	Negative      []string `yaml:"negative"`
	MinConfidence float64  `yaml:"min_confidence"`
}

type Config struct {
	Phrases PhraseGroups `yaml:"phrases"`
}

// Loader reloads phrases.yaml on every read (phase 3 — правка скрипта без рестарта).
type Loader struct {
	path string
	mu   sync.Mutex
}

func NewLoader(path string) *Loader {
	return &Loader{path: path}
}

func (l *Loader) Load() (*Config, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return Load(l.path)
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Phrases.MinConfidence == 0 {
		cfg.Phrases.MinConfidence = 0.55
	}
	return &cfg, nil
}
