package dsp

import (
	"math"

	"github.com/BeatGlow/audio"
	"github.com/BeatGlow/audio/dsp/fourier"
	"github.com/BeatGlow/audio/dsp/window"
)

// SQNR is the signal-to-quantization-noise ratio.
func SQNR(bits int) float64 {
	return 20 * math.Log10(math.Pow(2, float64(bits)))
}

// FrequencyPower denotes a single frequency and its magnitude in a Fast
// Fourier Transform of a signal.
type FrequencyPower struct {
	Frequency int
	Magnitude float64
}

type FrequencyPowerCalculator[T audio.Sample] struct {
	// SampleRate in samples per second.
	SampleRate int

	// Window function
	Window window.Window

	// buffer gets dynamically allocated
	samples audio.Samples[float64]
	complex []complex128
}

func NewFrequencyPowerCalculator[T audio.Sample](sampleRate int, window window.Window) *FrequencyPowerCalculator[T] {
	return &FrequencyPowerCalculator[T]{
		SampleRate: sampleRate,
		Window:     window,
	}
}

// Apply applies a Fast Fourier Transform (FFT) on a slice of float64 `data`,
// with sample rate `sampleRate`. It returns a slice of FrequencyPower.
func (c *FrequencyPowerCalculator[T]) Apply(dst []FrequencyPower, samples audio.Samples[T]) []FrequencyPower {
	var length = len(samples)
	if dst == nil {
		dst = make([]FrequencyPower, (length/2)-1)
	}

	// Convert samples to floats.
	c.samples = samples.ToFloat(c.samples)

	// Apply a window function to the values.
	if c.Window != nil {
		c.Window.Apply(c.samples)
	}

	// Convert samples to complex.
	c.complex = c.samples.ToComplex(c.complex)

	// apply a fast Fourier transform on the data; exclude index 0, no 0Hz-freq results
	spectrum := fourier.FFT(c.complex)

	for i := 1; i < length/2; i++ {
		freqReal := real(spectrum[i])
		freqImag := imag(spectrum[i])
		// map the magnitude for each frequency bin to the corresponding value in the map
		// using math.Sqrt(re*re + im*im) is faster than using math.Hypot(re, im)
		// see fft_test.go for details
		dst[i-1] = FrequencyPower{
			Frequency: i * c.SampleRate / length,
			Magnitude: math.Sqrt(freqReal*freqReal + freqImag*freqImag),
		}
	}

	return dst
}
