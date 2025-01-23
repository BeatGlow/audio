package filter

// FIR represents a Finite Impulse Response filter taking a sinc.
// https://en.wikipedia.org/wiki/Finite_impulse_response
type FIR struct {
	Sinc *Sinc
}

// LowPass applies a low pass filter using the FIR
func (f *FIR) LowPass(dst, src []float64) []float64 {
	return f.Convolve(dst, src, f.Sinc.LowPassCoefficients())
}

func (f *FIR) HighPass(dst, src []float64) []float64 {
	return f.Convolve(dst, src, f.Sinc.HighPassCoefficients())
}

// Convolve "mixes" two signals together
// kernels is the imput that is not part of our signal, it might be shorter
// than the origin signal.
func (f *FIR) Convolve(dst, src, kernels []float64) []float64 {
	if f == nil {
		return nil
	}
	if !(len(src) > len(kernels)) {
		// Provided data set is not greater than the filter weights.
		return nil
	}

	if len(dst) == 0 {
		dst = make([]float64, len(src))
	}
	for i := 0; i < len(kernels); i++ {
		var sum float64

		for j := 0; j < i; j++ {
			sum += (src[j] * kernels[len(kernels)-(1+i-j)])
		}
		dst[i] = sum
	}

	for i := len(kernels); i < len(src); i++ {
		var sum float64
		for j := 0; j < len(kernels); j++ {
			sum += (src[i-j] * kernels[j])
		}
		dst[i] = sum
	}

	return dst
}
