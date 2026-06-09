package stt

// Engine transcribes 16-bit LE PCM mono audio.
type Engine interface {
	Transcribe(pcm []byte, sampleRate int) (text string, confidence float64, err error)
}
