package classifier

import (
	"testing"

	"github.com/slavonnet/asterisk-response-classifier/internal/config"
)

func TestKeywordClassifier_Classify(t *testing.T) {
	cls := NewKeywordClassifier(config.PhraseGroups{
		Positive: []string{"да", "конечно"},
		Negative: []string{"нет", "не надо"},
	})

	tests := []struct {
		text  string
		label Label
	}{
		{"да", Positive},
		{"  Да, конечно  ", Positive},
		{"нет", Negative},
		{"не надо так", Negative},
		{"может быть", Uncertain},
		{"", Uncertain},
	}

	for _, tt := range tests {
		got := cls.Classify(tt.text)
		if got.Label != tt.label {
			t.Errorf("Classify(%q) = %q, want %q", tt.text, got.Label, tt.label)
		}
	}
}

func TestKeywordClassifier_ProcessAudioEmpty(t *testing.T) {
	cls := NewKeywordClassifier(config.PhraseGroups{})
	got := cls.ProcessAudio([]byte{0xff, 0x7f})
	if got.Label != Uncertain {
		t.Fatalf("ProcessAudio label = %q, want uncertain", got.Label)
	}
}
