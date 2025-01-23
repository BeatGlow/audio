package filter

import (
	"github.com/BeatGlow/audio/dsp/filter"
	"github.com/BeatGlow/audio/dsp/window"
)

// LowPass is a basic LowPass filter cutting off
// CutOffFreq is where the filter would be at -3db.
// TODO: param to say how efficient we want the low pass to be.
// matlab: lpFilt = designfilt('lowpassfir','PassbandFrequency',0.25, ...
//
//	'StopbandFrequency',0.35,'PassbandRipple',0.5, ...
//	'StopbandAttenuation',65,'DesignMethod','kaiserwin');
func LowPass(dst, src []float64, cutOffFrequency float64, sampleRate int) []float64 {
	s := &filter.Sinc{
		// TODO: find the right taps number to do a proper
		// audio low pass based in the sample rate
		// there should be a magical function to get that number.
		Taps:            62,
		SampleRate:      sampleRate,
		CutOffFrequency: cutOffFrequency,
		Window:          window.Hamming,
	}
	fir := &filter.FIR{Sinc: s}
	return fir.LowPass(dst, src)
}

// HighPass is a basic LowPass filter cutting off
// the audio buffer frequencies below the cutOff frequency.
func HighPass(dst, src []float64, cutOffFrequency float64, sampleRate int) []float64 {
	s := &filter.Sinc{
		Taps:            62,
		SampleRate:      sampleRate,
		CutOffFrequency: cutOffFrequency,
		Window:          window.Blackman,
	}
	fir := &filter.FIR{Sinc: s}
	return fir.HighPass(dst, src)
}
