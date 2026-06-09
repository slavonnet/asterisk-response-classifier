package stt

import "encoding/json"

// GrammarJSON builds Vosk grammar from phrase list.
func GrammarJSON(phrases []string) (string, error) {
	b, err := json.Marshal(phrases)
	return string(b), err
}
