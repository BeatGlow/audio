package audio

import (
	"encoding/binary"
	"math"

	"golang.org/x/exp/constraints"
)

// Sample is an audio sample, it can be any of the native Go numeric types.
type Sample interface {
	constraints.Integer | constraints.Float
}

// Samples represents a mono audio buffer, or a single channel.
type Samples[T Sample] []T

// Push adds samples to the slice.
func Push[T Sample](buffer *Samples[T], samples ...T) {
	*buffer = append(*buffer, samples...)
}

// Pop removes one sample from the end of the slice.
func Pop[T Sample](buffer *Samples[T]) T {
	var removed T //nolint:gosimple
	removed = (*buffer)[len(*buffer)-1]
	*buffer = (*buffer)[:len(*buffer)-1]
	return removed
}

// Shift removes n samples from the start of the slice.
func Shift[T Sample](buffer *Samples[T], n int) Samples[T] {
	var shifted Samples[T] //nolint:gosimple
	shifted = (*buffer)[:n]
	*buffer = (*buffer)[n:]
	return shifted
}

// Unshift adds samples to the start of the slice.
func Unshift[T Sample](buffer *Samples[T], samples ...T) {
	*buffer = append(samples, *buffer...)
}

// Decodes values encoded in b using the specified byte order.
func (s Samples[T]) Decode(b []byte, order binary.ByteOrder) {
	var zero T
	switch any(zero).(type) {
	case int:
		switch s.BitsPerSample() {
		case 32:
			for i := range s {
				s[i] = T(order.Uint32(b[i<<2:]))
			}
		case 64:
			for i := range s {
				s[i] = T(order.Uint64(b[i<<3:]))
			}
		}
	case int8:
		for i, x := range b {
			s[i] = T(x)
		}
	case int16:
		for i := range s {
			s[i] = T(order.Uint16(b[i<<1:]))
		}
	case int32:
		for i := range s {
			s[i] = T(order.Uint32(b[i<<2:]))
		}
	case int64:
		for i := range s {
			s[i] = T(order.Uint64(b[i<<3:]))
		}
	case uint:
		switch s.BitsPerSample() {
		case 32:
			for i := range s {
				s[i] = T(order.Uint32(b[i<<2:]))
			}
		case 64:
			for i := range s {
				s[i] = T(order.Uint64(b[i<<3:]))
			}
		}
	case uint8:
		for i, x := range b {
			s[i] = T(x)
		}
	case uint16:
		for i := range s {
			s[i] = T(order.Uint16(b[i<<1:]))
		}
	case uint32:
		for i := range s {
			s[i] = T(order.Uint32(b[i<<2:]))
		}
	case uint64:
		for i := range s {
			s[i] = T(order.Uint64(b[i<<3:]))
		}
	case float32:
		for i := range s {
			s[i] = T(math.Float32frombits(order.Uint32(b[i<<2:])))
		}
	case float64:
		for i := range s {
			s[i] = T(math.Float64frombits(order.Uint64(b[i<<3:])))
		}
	}
}

// Encode the sample slice as bytes in b.
func (s Samples[T]) Encode(b []byte, order binary.ByteOrder) {
	var zero T
	switch any(zero).(type) {
	case int:
		switch s.BitsPerSample() {
		case 32:
			for i, v := range s {
				order.PutUint32(b[i<<2:], uint32(v))
			}
		case 64:
			for i, v := range s {
				order.PutUint64(b[i<<3:], uint64(v))
			}
		}
	case int8:
		for i, v := range s {
			b[i] = byte(v)
		}
	case int16:
		for i, v := range s {
			order.PutUint16(b[i<<1:], uint16(v))
		}
	case int32:
		for i, v := range s {
			order.PutUint32(b[i<<2:], uint32(v))
		}
	case int64:
		for i, v := range s {
			order.PutUint64(b[i<<3:], uint64(v))
		}
	case uint:
		switch s.BitsPerSample() {
		case 32:
			for i, v := range s {
				order.PutUint32(b[i<<2:], uint32(v))
			}
		case 64:
			for i, v := range s {
				order.PutUint64(b[i<<3:], uint64(v))
			}
		}
	case uint8:
		for i, v := range s {
			b[i] = byte(v)
		}
	case uint16:
		for i, v := range s {
			order.PutUint16(b[i<<1:], uint16(v))
		}
	case uint32:
		for i, v := range s {
			order.PutUint32(b[i<<2:], uint32(v))
		}
	case uint64:
		for i, v := range s {
			order.PutUint64(b[i<<3:], uint64(v))
		}
	case float32:
		for i, v := range s {
			order.PutUint32(b[i<<2:], math.Float32bits(float32(v)))
		}
	case float64:
		for i, v := range s {
			order.PutUint64(b[i<<3:], math.Float64bits(float64(v)))
		}
	}
}

// intSize is the number of bits used by the int/uint type, dependent on CPU architecture
const intSize = 32 << (^uint(0) >> 63)

// BitsPerSample is the number of bits required to store one sample.
func (s Samples[T]) BitsPerSample() int {
	var zero T
	switch any(zero).(type) {
	case int, uint:
		return intSize // math.intSize
	case int8, uint8:
		return 8
	case int16, uint16:
		return 16
	case int32, uint32, float32:
		return 32
	case int64, uint64, float64:
		return 64
	case complex64:
		return 128
	case complex128:
		return 256
	default:
		return 0
	}
}

// Buffer represents a multi channel audio buffer.
type Buffer[T Sample] [][]T

// Channels is the number of channels in the buffer.
func (b Buffer[T]) Channels() int {
	return len(b)
}

// Samples is the number of samples per channel.
func (b Buffer[T]) Samples() int {
	if len(b) == 0 {
		return 0
	}

	min := len(b[0])
	for _, samples := range b[1:] {
		if l := len(samples); l < min {
			min = l
		}
	}
	return min
}

// BitsPerSample is the number of bits required to store one sample.
func (b Buffer[T]) BitsPerSample() int {
	var zero Samples[T]
	return zero.BitsPerSample()
}

// minOfType returns the minimum value for the given type.
func minOfType[T Sample](value T) T {
	switch any(value).(type) {
	case int:
		var x int = math.MinInt
		return T(x)
	case int8:
		var x int8 = math.MinInt8
		return T(x)
	case int16:
		var x int16 = math.MinInt16
		return T(x)
	case int32:
		var x int32 = math.MinInt32
		return T(x)
	case int64:
		var x int64 = math.MinInt64
		return T(x)
	case float32:
		var x float32 = -math.MaxFloat32
		return T(x)
	case float64:
		var x float64 = -math.MaxFloat64
		return T(x)
	default:
		var zero T
		return zero
	}
}

// maxOfType returns the maximum value for the given type.
func maxOfType[T Sample](value T) T {
	switch any(value).(type) {
	case int:
		var x int = math.MaxInt
		return T(x)
	case int8:
		var x int8 = math.MaxInt8
		return T(x)
	case uint8:
		var x uint8 = math.MaxUint8
		return T(x)
	case int16:
		var x int16 = math.MaxInt16
		return T(x)
	case uint16:
		var x uint16 = math.MaxUint16
		return T(x)
	case int32:
		var x int32 = math.MaxInt32
		return T(x)
	case uint32:
		var x uint32 = math.MaxUint32
		return T(x)
	case int64:
		var x int64 = math.MaxInt64
		return T(x)
	case uint64:
		var x uint64 = math.MaxUint64
		return T(x)
	case float32:
		var x float32 = math.MaxFloat32
		return T(x)
	case float64:
		var x float64 = math.MaxFloat64
		return T(x)
	default: /* unreachable */
		var zero T
		return zero
	}
}

// Abs returns the absolute value of x.
func Abs[T Sample](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Min returns the smallest value.
//
// Example:
//
//	Min([]int8{-127, -63, 0, 63, 128}) = -127
func Min[T Sample](values Samples[T]) T {
	if len(values) == 0 {
		return T(0)
	}

	var min T
	min = maxOfType(min)
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

// Max returns the largest value.
//
// Example:
//
//	Max([]int8{-127, -63, 0, 63, 128}) = 128
func Max[T Sample](values Samples[T]) T {
	if len(values) == 0 {
		return T(0)
	}

	var max T
	max = minOfType(max)
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

// Mean of values.
func Mean[T Sample](samples Samples[T]) T {
	if len(samples) == 0 {
		return 0
	}

	// an algorithm that attempts to retain accuracy
	// with widely different values.
	var parts []float64
	for _, v := range samples {
		var (
			x = float64(v)
			i int
		)
		for _, p := range parts {
			sum := p + x
			var err float64
			switch ax, ap := math.Abs(x), math.Abs(p); {
			case ax < ap:
				err = x - (sum - p)
			case ap < ax:
				err = p - (sum - x)
			}
			if err != 0 {
				parts[i] = err
				i++
			}
			x = sum
		}
		parts = append(parts[:i], x)
	}

	var sum float64
	for _, x := range parts {
		sum += x
	}
	return T(sum / float64(len(samples)))
}

// Clip samples between min and max.
func Clip[T Sample](min, max T, samples Samples[T]) {
	for i, v := range samples {
		if v < min {
			samples[i] = min
		}
		if v > max {
			samples[i] = max
		}
	}
}

// RMS returns the root mean square of values.
func RMS[T Sample](samples Samples[T]) T {
	if len(samples) == 0 {
		return T(0)
	}

	var squares float64
	for _, v := range samples {
		squares += float64(v) * float64(v)
	}
	return T(math.Sqrt(squares / float64(len(samples))))
}

// Normalize samples to fit in [-1..1].
//
// It doesn't make a whole lot of sense to call this function for non-float values, but you could.
//
// NB: this function may change to bias around the centerpoint of math.MaxUint* for unsigned types.
func Normalize[T Sample](samples Samples[T]) {
	if len(samples) == 0 {
		return
	}

	max := Abs(Max(samples))
	min := Abs(Min(samples))
	if max < min {
		max = min
	}
	for i := range samples {
		samples[i] /= max
	}
}

// Deinterleave splits interleaved samples into a multi channel buffer.
//
// The provided buffer dst should be nil, or adequately sized to fit channels of samples.
func Deinterleave[T Sample](dst Buffer[T], src Samples[T], channels int) Buffer[T] {
	// Fast path
	if channels == 0 {
		return nil
	}

	samples := len(src) / channels

	if dst == nil {
		dst = make(Buffer[T], channels)
		for i := 0; i < channels; i++ {
			dst[i] = make(Samples[T], samples)
		}
	}

	if channels == 1 {
		copy(dst[0], src)
	} else {
		for i, sample := range src {
			dst[i%channels][i/channels] = sample
		}
	}

	return dst
}

// Interleave a buffer into a single samples slice.
//
// The provided samples dst should be nil, or adequately sized to fit channels of samples.
func Interleave[T Sample](dst Samples[T], src Buffer[T]) Samples[T] {
	if len(src) == 0 {
		return dst
	}

	channels := len(src)
	samples := len(src[0])

	if dst == nil {
		dst = make(Samples[T], channels*samples)
	}

	for i := 0; i < samples; i++ {
		o := i * channels
		for c := 0; c < channels; c++ {
			dst[o+c] = src[c][i]
		}
	}

	return dst
}
