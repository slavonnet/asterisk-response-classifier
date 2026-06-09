//go:build !phrase

package phrase

import (
	"fmt"

	"github.com/slavonnet/asterisk-response-classifier/internal/config"
)

type Engine struct{}

func New(modelPath string) (*Engine, error) {
	if modelPath != "" {
		return nil, fmt.Errorf("build with -tags phrase required (see README)")
	}
	return &Engine{}, nil
}

func (e *Engine) Recognize(_ []byte, _ *config.Config) (label, heard string, score float64) {
	return "uncertain", "", 0
}

func (e *Engine) Close() {}
