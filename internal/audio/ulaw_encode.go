package audio

func EncodeULaw(samples []int16) []byte {
	out := make([]byte, len(samples))
	for i, s := range samples {
		out[i] = linearToUlaw(s)
	}
	return out
}

func linearToUlaw(sample int16) byte {
	// G.711 µ-law encode (simplified)
	const bias = 0x84
	sign := 0
	if sample < 0 {
		sign = 0x80
		sample = -sample
		if sample < 0 {
			sample = 0x7FFF
		}
	}
	if sample > 0x7FFF {
		sample = 0x7FFF
	}
	sample += bias
	exponent := 7
	for expMask := 0x4000; exponent > 0; exponent-- {
		if sample&expMask != 0 {
			break
		}
		expMask >>= 1
	}
	mantissa := (sample >> (exponent + 3)) & 0x0F
	return ^byte(sign | byte(exponent<<4) | byte(mantissa))
}
