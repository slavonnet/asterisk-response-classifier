package config

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// Sentences — как custom sentences в speech-to-phrase (ограниченный набор фраз).
type Sentences struct {
	Language string            `yaml:"language"`
	Lists    map[string]List   `yaml:"lists"`
	Intents  map[string]Intent `yaml:"intents"`
}

type List struct {
	Values []ListValue `yaml:"values"`
}

type ListValue struct {
	In  string `yaml:"in"`
	Out string `yaml:"out"`
}

type Intent struct {
	Data []IntentData `yaml:"data"`
}

type IntentData struct {
	Sentences []string `yaml:"sentences"`
}

type Config struct {
	Sentences Sentences
	path      string
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
	var s Sentences
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &Config{Sentences: s, path: path}, nil
}

func (c *Config) Dir() string { return filepath.Dir(c.path) }

// AllPhrases returns every speakable phrase and its label (positive/negative).
func (c *Config) AllPhrases() (phrases []string, labelOf map[string]string) {
	labelOf = make(map[string]string)
	seen := make(map[string]bool)

	add := func(phrase, label string) {
		phrase = strings.TrimSpace(phrase)
		if phrase == "" || seen[phrase] {
			return
		}
		seen[phrase] = true
		phrases = append(phrases, phrase)
		labelOf[norm(phrase)] = label
	}

	for intentName, intent := range c.Sentences.Intents {
		label := intentLabel(intentName)
		for _, d := range intent.Data {
			for _, sent := range d.Sentences {
				for _, p := range expand(sent, c.Sentences.Lists) {
					add(p, label)
				}
			}
		}
	}
	return phrases, labelOf
}

func intentLabel(name string) string {
	switch strings.ToLower(name) {
	case "positive":
		return "positive"
	case "negative":
		return "negative"
	default:
		return strings.ToLower(name)
	}
}

func expand(tmpl string, lists map[string]List) []string {
	tmpl = strings.TrimSpace(tmpl)
	if strings.HasPrefix(tmpl, "{") && strings.HasSuffix(tmpl, "}") {
		key := tmpl[1 : len(tmpl)-1]
		if list, ok := lists[key]; ok {
			out := make([]string, 0, len(list.Values))
			for _, v := range list.Values {
				if v.In != "" {
					out = append(out, v.In)
				}
			}
			return out
		}
	}
	return []string{tmpl}
}

func norm(s string) string { return strings.ToLower(strings.TrimSpace(s)) }
