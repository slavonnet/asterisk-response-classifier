package audio

import "testing"

func TestCosineIdentical(t *testing.T) {
	a := FeatureVector{1, 0, 0}
	b := FeatureVector{1, 0, 0}
	if CosineSimilarity(a, b) < 0.99 {
		t.Fatal("expected ~1")
	}
}

func TestExtractDifferentPatterns(t *testing.T) {
	short := make([]int16, 800)
	for i := range short {
		short[i] = int16(i * 30)
	}
	long := make([]int16, 4000)
	for i := range long {
		long[i] = int16((i * 7) % 5000)
	}
	fa := ExtractFeatures(encodePCM(short))
	fb := ExtractFeatures(encodePCM(long))
	if CosineSimilarity(fa, fb) > 0.95 {
		t.Fatal("expected different features")
	}
}

func encodePCM(samples []int16) []byte {
	pcm := make([]byte, len(samples)*2)
	for i, s := range samples {
		pcm[i*2] = byte(s)
		pcm[i*2+1] = byte(s >> 8)
	}
	return pcm
}
