package classifier

// Label is returned to Asterisk dialplan via SPEECH_TEXT(0).
type Label string

const (
	Positive  Label = "positive"
	Negative  Label = "negative"
	Uncertain Label = "uncertain"
)

type Result struct {
	Text  string  // positive | negative | uncertain
	Score float64 // similarity 0..100
	Label Label
}
