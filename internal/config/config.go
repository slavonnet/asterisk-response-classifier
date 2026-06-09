package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type PhraseGroups struct {
	Positive      []string `yaml:"positive"`
	Negative      []string `yaml:"negative"`
	MinConfidence float64  `yaml:"min_confidence"`
}

type TreeNode struct {
	ID       string            `yaml:"id"`
	Prompt   string            `yaml:"prompt"`   // sound file or TTS key
	On       map[string]string `yaml:"on"`       // positive|negative|uncertain -> next node id or action
	Timeout  int               `yaml:"timeout"`  // seconds
	MaxRetry int               `yaml:"max_retry"`
}

type DecisionTree struct {
	Start string     `yaml:"start"`
	Nodes []TreeNode `yaml:"nodes"`
}

type Config struct {
	Phrases PhraseGroups `yaml:"phrases"`
	Tree    DecisionTree `yaml:"tree"`
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
		cfg.Phrases.MinConfidence = 0.6
	}
	return &cfg, nil
}
