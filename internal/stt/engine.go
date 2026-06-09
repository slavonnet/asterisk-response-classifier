package stt

// Engine recognizes speech limited to the given phrase list (speech-to-phrase style).
// When phrases change, the next call uses the new list — no retraining on the server.
type Engine interface {
	Transcribe(pcm []byte, sampleRate int, phrases []string) (text string, confidence float64, err error)
}
