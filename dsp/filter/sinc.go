package filter

import (
	"math"

	"github.com/BeatGlow/audio/dsp/window"
)

// Sinc represents a sinc function
// The sinc function also called the "sampling function," is a function that
// arises frequently in signal processing and the theory of Fourier transforms.
// The full name of the function is "sine cardinal," but it is commonly referred to by
// its abbreviation, "sinc."
// http://mathworld.wolfram.com/SincFunction.html
type Sinc struct {
	CutOffFrequency float64

	SampleRate int

	// Taps are the numbers of samples we go back in time when processing the sync function.
	// The tap numbers will affect the shape of the filter. The more taps, the more
	// shape but the more delays being injected.
	Taps int

	// Window function to apply.
	Window window.GeneratorFunc

	lpCoefficients []float64
	hpCoefficients []float64
}

// LowPassCoefficients returns the coeficients to create a low pass filter
func (s *Sinc) LowPassCoefficients() []float64 {
	if s == nil {
		return nil
	}
	if len(s.lpCoefficients) > 0 {
		return s.lpCoefficients
	}
	size := s.Taps + 1
	// sample rate is 2 pi radians per second.
	// we get the cutt off frequency in radians per second
	b := (2 * math.Pi) * s.TransitionFrequency()
	s.lpCoefficients = make([]float64, size)
	// we use a window of size taps + 1
	winData := s.Window(size)

	// we only do half the taps because the coefs are symmetric
	// but we fill up all the coefs
	for i := 0; i < (s.Taps / 2); i++ {
		c := float64(i) - float64(s.Taps)/2
		y := math.Sin(c*b) / (math.Pi * c)
		s.lpCoefficients[i] = y * winData[i]
		s.lpCoefficients[size-1-i] = s.lpCoefficients[i]
	}

	// then we do the ones we missed in case we have an odd number of taps
	s.lpCoefficients[s.Taps/2] = 2 * s.TransitionFrequency() * winData[s.Taps/2]
	return s.lpCoefficients
}

// HighPassCoefficients returns the coeficients to create a high pass filter
func (s *Sinc) HighPassCoefficients() []float64 {
	if s == nil {
		return nil
	}
	if len(s.hpCoefficients) > 0 {
		return s.hpCoefficients
	}

	// we take the low pass coesf and invert them
	size := s.Taps + 1
	s.hpCoefficients = make([]float64, size)
	lowPassCoefs := s.LowPassCoefficients()
	winData := s.Window(size)

	for i := 0; i < (s.Taps / 2); i++ {
		s.hpCoefficients[i] = -lowPassCoefs[i]
		s.hpCoefficients[size-1-i] = s.hpCoefficients[i]
	}
	s.hpCoefficients[s.Taps/2] = (1 - 2*s.TransitionFrequency()) * winData[s.Taps/2]
	return s.hpCoefficients
}

// TransitionFrequency returns a ratio of the cutoff frequency and the sample rate.
func (s *Sinc) TransitionFrequency() float64 {
	if s == nil {
		return 0
	}
	return s.CutOffFrequency / float64(s.SampleRate)
}
