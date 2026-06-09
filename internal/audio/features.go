package audio

import (
	"math"
)

// FeatureVector — compact acoustic fingerprint (no text, no ASR).
type FeatureVector []float64

// ExtractFeatures builds a vector from 16-bit LE mono PCM @ 8 kHz.
func ExtractFeatures(pcm []byte) FeatureVector {
	samples := bytesToSamples(pcm)
	if len(samples) == 0 {
		return nil
	}
	samples = trimSilence(samples, 200)
	if len(samples) == 0 {
		return nil
	}

	rms := sampleRMS(samples)
	zcr := zeroCrossRate(samples)
	dur := float64(len(samples)) / 8000.0

	frames := 8
	frameLen := len(samples) / frames
	if frameLen < 1 {
		frameLen = 1
	}

	out := make(FeatureVector, 0, 2+frames)
	out = append(out, dur, rms, zcr)
	for i := 0; i < frames; i++ {
		start := i * frameLen
		end := start + frameLen
		if i == frames-1 {
			end = len(samples)
		}
		out = append(out, sampleRMS(samples[start:end]))
	}
	return normalize(out)
}

func bytesToSamples(pcm []byte) []int16 {
	n := len(pcm) / 2
	out := make([]int16, n)
	for i := 0; i < n; i++ {
		out[i] = int16(pcm[i*2]) | int16(pcm[i*2+1])<<8
	}
	return out
}

func trimSilence(samples []int16, threshold int16) []int16 {
	start, end := 0, len(samples)
	for start < end && abs16(samples[start]) < threshold {
		start++
	}
	for end > start && abs16(samples[end-1]) < threshold {
		end--
	}
	return samples[start:end]
}

func abs16(v int16) int16 {
	if v < 0 {
		return -v
	}
	return v
}

func sampleRMS(samples []int16) float64 {
	if len(samples) == 0 {
		return 0
	}
	var sum float64
	for _, s := range samples {
		f := float64(s)
		sum += f * f
	}
	return math.Sqrt(sum / float64(len(samples)))
}

func zeroCrossRate(samples []int16) float64 {
	if len(samples) < 2 {
		return 0
	}
	var crosses int
	for i := 1; i < len(samples); i++ {
		if (samples[i-1] >= 0 && samples[i] < 0) || (samples[i-1] < 0 && samples[i] >= 0) {
			crosses++
		}
	}
	return float64(crosses) / float64(len(samples)-1)
}

func normalize(v FeatureVector) FeatureVector {
	var sum float64
	for _, x := range v {
		sum += x * x
	}
	if sum == 0 {
		return v
	}
	n := math.Sqrt(sum)
	out := make(FeatureVector, len(v))
	for i, x := range v {
		out[i] = x / n
	}
	return out
}

// CosineSimilarity of two normalized vectors (0..1).
func CosineSimilarity(a, b FeatureVector) float64 {
	if len(a) == 0 || len(a) != len(b) {
		return 0
	}
	var dot float64
	for i := range a {
		dot += a[i] * b[i]
	}
	if dot < 0 {
		return 0
	}
	if dot > 1 {
		return 1
	}
	return dot
}
