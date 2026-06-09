package config

import (
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// References lists ulaw reference clips (positive / negative variants).
type References struct {
	Positive  []string `yaml:"positive"`
	Negative  []string `yaml:"negative"`
	MinScore  float64  `yaml:"min_score"`  // min cosine similarity 0..1
	MinMargin float64  `yaml:"min_margin"` // pos_best - neg_best (or vice versa)
}

type Config struct {
	References References `yaml:"references"`
	baseDir    string
}

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
	if cfg.References.MinScore == 0 {
		cfg.References.MinScore = 0.55
	}
	if cfg.References.MinMargin == 0 {
		cfg.References.MinMargin = 0.08
	}
	cfg.baseDir = filepath.Dir(path)
	return &cfg, nil
}

func (c *Config) Resolve(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(c.baseDir, path)
}

func (c *Config) BaseDir() string { return c.baseDir }
