package modulation

import (
	"math"
)

// QPSKEncode encodes a binary stream into QPSK phase shifts
func QPSKEncode(bits []byte) []complex64 {
	var symbols []complex64
	for _, b := range bits {
		// Map bits to one of the four phase states
		switch b {
		case 0:
			symbols = append(symbols, complex(1, 0)) // 0° phase
		case 1:
			symbols = append(symbols, complex(0, 1)) // 90° phase
		case 2:
			symbols = append(symbols, complex(-1, 0)) // 180° phase
		case 3:
			symbols = append(symbols, complex(0, -1)) // 270° phase
		}
	}
	return symbols
}

// QPSKDecode decodes QPSK phase shifts into a binary stream
func QPSKDecode(symbols []complex64) []byte {
	var bits []byte
	for _, symbol := range symbols {
		angle := math.Atan2(float64(imag(symbol)), float64(real(symbol)))
		switch {
		case angle == 0:
			bits = append(bits, 0)
		case angle == math.Pi/2:
			bits = append(bits, 1)
		case angle == math.Pi:
			bits = append(bits, 2)
		case angle == -math.Pi/2:
			bits = append(bits, 3)
		}
	}
	return bits
}
