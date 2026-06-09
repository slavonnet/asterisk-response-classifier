package audio

import "testing"

func TestDecodeULaw(t *testing.T) {
	pcm := DecodeULaw([]byte{0xff, 0x7f})
	if len(pcm) != 4 {
		t.Fatalf("len = %d, want 4", len(pcm))
	}
}
