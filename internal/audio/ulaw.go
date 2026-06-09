package audio

// DecodeULaw converts 8-bit G.711 µ-law (Asterisk ulaw) to 16-bit signed PCM mono.
func DecodeULaw(ulaw []byte) []byte {
	pcm := make([]byte, len(ulaw)*2)
	for i, u := range ulaw {
		s := ulawToLinear(u)
		pcm[i*2] = byte(s)
		pcm[i*2+1] = byte(s >> 8)
	}
	return pcm
}

func ulawToLinear(u byte) int16 {
	u = ^u
	sign := u & 0x80
	exponent := (u >> 4) & 0x07
	mantissa := u & 0x0F
	sample := int16(mantissa<<4) + 0x108
	sample <<= exponent
	sample -= 0x108
	if sign != 0 {
		sample = -sample
	}
	return sample
}
