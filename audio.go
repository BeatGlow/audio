// Package audio ...
package audio

// Format describes the audio format.
type Format interface {
	// Channels is the number of channels.
	Channels() int

	// BitsPerSample is the number of bits required to store one sample.
	BitsPerSample() int
}
