package classifier

// Label is the sentiment/decision label returned to Asterisk dialplan.
type Label string

const (
	Positive  Label = "positive"
	Negative  Label = "negative"
	Uncertain Label = "uncertain"
)

// Result holds recognition output for AEAP.
type Result struct {
	Text       string  // label for dialplan: positive|negative|uncertain
	Score      float64 // confidence 0..100
	Label      Label
	Transcript string // raw STT text (logs/debug)
}

// Classifier maps recognized speech to a decision label.
type Classifier interface {
	// Classify returns label for transcribed text.
	Classify(text string) Result
	// ProcessAudio runs STT + classification on raw audio (ulaw 8kHz).
	// MVP uses placeholder; ONNX STT plugs in here.
	ProcessAudio(data []byte) Result
}
