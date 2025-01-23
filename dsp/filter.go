package dsp

import "github.com/BeatGlow/audio"

// DCRemoval removes the DC component from a signal.
func DCRemval[T audio.Sample](signal []T) {
	if len(signal) == 0 {
		return
	}

	var mean T
	for _, v := range signal {
		mean += v
	}

	mean /= T(len(signal))
	for i := range signal {
		signal[i] -= mean
	}
}
