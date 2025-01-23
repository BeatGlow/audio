package fourier

import (
	"math"
)

const (
	tau = 2 * math.Pi

	// DefaultMagnitudeThreshold describes the default value where a certain
	// frequency is strong enough to be considered relevant to the spectrum filter.
	DefaultMagnitudeThreshold = 10
)

// FFT applies a Fast Fourier Transform to the input slice of complex128 values, to
// retrieve the frequency spectrum of a digital signal.
func FFT(src []complex128) []complex128 {
	var (
		valueLen = len(src)
		factors  = GetRadix2Factors(valueLen)
		temp     = make([]complex128, valueLen) // temp
	)

	src = ReorderData(src)

	// stage increases by a power of two
	for stage := 2; stage <= valueLen; stage <<= 1 {
		var (
			blocks      = valueLen / stage
			stage2Value = stage / 2
		)

		// iterate through each item in the batch, increasing by the stage value
		for batchIdx := 0; batchIdx < valueLen; batchIdx += stage {
			if stage == 2 { // "first stage" scenario
				var (
					reorderIdx  = src[batchIdx]
					reorderNext = src[batchIdx+1]
				)

				temp[batchIdx] = reorderIdx + reorderNext
				temp[batchIdx+1] = reorderIdx - reorderNext

				continue
			}

			for iter := 0; iter < stage2Value; iter++ {
				var (
					idx        = iter + batchIdx
					idx2       = idx + stage2Value
					reorderIdx = src[idx]
					factorized = src[idx2] * factors[blocks*iter]
				)

				temp[idx] = reorderIdx + factorized
				temp[idx2] = reorderIdx - factorized
			}
		}

		src, temp = temp, src
	}

	return src
}

// IFFT returns the Inverse Fast Fourier Transform of a given complex128 slice.
func IFFT(value []complex128) []complex128 {
	var (
		ln     = len(value)
		output = make([]complex128, ln)
		factor = complex(float64(ln), 0)
	)

	// Reverse inputs, which is calculated with modulo factor, hence value[0] as an outlier
	output[0] = value[0]

	for i, j := 1, ln-1; i < ln; i, j = i+1, j-1 {
		output[i] = value[j]
	}

	output = FFT(output)

	for i := range output {
		output[i] /= factor
	}

	return output
}

// Convolve returns the convolution of x ∗ y, applied to the complex128 slice x.
func Convolve(x, y []complex128) []complex128 {
	if len(x) != len(y) {
		return nil
	}

	x = FFT(x)
	y = FFT(y)

	for i := 0; i < len(x); i++ {
		x[i] *= y[i]
	}

	return IFFT(x)
}
