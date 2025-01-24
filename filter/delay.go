package filter

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/BeatGlow/audio"
)

var (
	ErrChannels      = errors.New("filter: need more than 0 channels")
	ErrDelayNegative = errors.New("filter: delay can't be negative")
)

// Delay is a general purpose delay line.
type Delay[T audio.Sample] struct {
	audio.Reader[T]
	audio.Samples[T]
	time.Duration
}

// NewDelay introduces a fixed delay for reading samples from r.
func NewDelay[T audio.Sample](r audio.Reader[T], channels, sampleRate int, delay time.Duration) (audio.Reader[T], error) {
	if channels < 1 {
		return nil, ErrChannels
	}

	if delay < 0 {
		return nil, ErrDelayNegative
	} else if delay == 0 {
		return r, nil
	}

	samplesPerDelay := channels * int(math.Round(float64(sampleRate)*float64(delay)/float64(time.Second)))
	return &Delay[T]{
		Reader:   r,
		Samples:  make(audio.Samples[T], samplesPerDelay),
		Duration: delay,
	}, nil
}

func (d *Delay[T]) String() string {
	return fmt.Sprintf("delay %s", d.Duration)
}

func (d *Delay[T]) ReadSamples(buffer audio.Samples[T]) (int, error) {
	if len(d.Samples) == 0 {
		return d.Reader.ReadSamples(buffer)
	}

	// Consume delay line samples.
	n := copy(buffer, d.Samples)
	if r := len(buffer) - n; r > 0 {
		// We need more, consume from Reader.
		if _, err := d.Reader.ReadSamples(buffer[n:]); err != nil {
			return n, err
		}
	}

	// Advance our delay line.
	copy(d.Samples, d.Samples[n:])

	// Replenish our delay line by consuming samples from Reader.
	if n > 0 {
		if _, err := d.Reader.ReadSamples(d.Samples[len(d.Samples)-n:]); err != nil {
			return n, err
		}
	}

	return len(buffer), nil
}
