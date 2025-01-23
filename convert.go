package audio

import (
	"fmt"
	"math"
)

const (
	intToFloat   = math.MaxInt + 1
	int8ToFloat  = math.MaxInt8 + 1
	int16ToFloat = math.MaxInt16 + 1
	int32ToFloat = math.MaxInt32 + 1
	int64ToFloat = math.MaxInt64 + 1
)

// ToInt16 converts samples in s to []int16 in dst.
//
// The provided buffer dst should be nil, or adequately sized to fit all samples.
func (s Samples[T]) ToInt16(dst []int16) []int16 {
	if dst == nil {
		dst = make([]int16, len(s))
	}

	var value T
	switch any(value).(type) {
	case []int:
		var shift = 48 // for 64-bit
		if intSize == 32 {
			shift = 16 // for 32-bit
		}
		for i, v := range s {
			dst[i] = int16(int(v) >> shift)
		}
	case []int8:
		for i, v := range s {
			dst[i] = int16(v)<<8 | int16(v)
		}
	case []int16:
		for i, v := range s {
			dst[i] = int16(v)
		}
	case []int32:
		for i, v := range s {
			dst[i] = int16(int32(v) >> 16)
		}
	case []int64:
		for i, v := range s {
			dst[i] = int16(int64(v) >> 48)
		}
	case []uint:
		var shift = 48 // for 64-bit
		if intSize == 32 {
			shift = 16 // for 32-bit
		}
		for i, v := range s {
			dst[i] = int16(uint(v) >> shift)
		}
	case []uint8:
		for i, v := range s {
			dst[i] = int16(uint16(v)<<8 | uint16(v))
		}
	case []uint16:
		for i, v := range s {
			dst[i] = int16(v)
		}
	case []uint32:
		for i, v := range s {
			dst[i] = int16(uint32(v) >> 16)
		}
	case []uint64:
		for i, v := range s {
			dst[i] = int16(uint64(v) >> 48)
		}
	case []float32:
		for i, v := range s {
			dst[i] = int16(float32(v)*int16ToFloat + .5)
		}
	case []float64:
		for i, v := range s {
			dst[i] = int16(float64(v)*int16ToFloat + .5)
		}
	}

	return dst
}

// ToFloat converts samples in s to []float64 in dst.
//
// The provided buffer dst should be nil, or adequately sized to fit all samples.
func (s Samples[T]) ToFloat(dst []float64) []float64 {
	if dst == nil {
		dst = make([]float64, len(s))
	}

	var value T
	switch any(value).(type) {
	case int:
		for i, v := range s {
			dst[i] = float64(v) / intToFloat
		}
	case int8:
		for i, v := range s {
			dst[i] = float64(v) / int8ToFloat
		}
	case int16:
		for i, v := range s {
			dst[i] = float64(v) / int16ToFloat
		}
	case int32:
		for i, v := range s {
			dst[i] = float64(v) / int32ToFloat
		}
	case int64:
		for i, v := range s {
			dst[i] = float64(v) / int64ToFloat
		}
	case uint:
		for i, v := range s {
			dst[i] = (float64(v) - intToFloat) / intToFloat
		}
	case uint8:
		for i, v := range s {
			dst[i] = (float64(v) - int8ToFloat) / int8ToFloat
		}
	case uint16:
		for i, v := range s {
			dst[i] = (float64(v) - int16ToFloat) / int16ToFloat
		}
	case uint32:
		for i, v := range s {
			dst[i] = (float64(v) - int32ToFloat) / int32ToFloat
		}
	case uint64:
		for i, v := range s {
			dst[i] = (float64(v) - int64ToFloat) / int64ToFloat
		}
	case float32, float64:
		for i, v := range s {
			dst[i] = float64(v)
		}
	default:
		panic(fmt.Sprintf("can't convert %T to []float64", value))
	}

	return dst
}

// ToComplex converts the samples in s to []complex128 in dst.
func (s Samples[T]) ToComplex(dst []complex128) []complex128 {
	if dst == nil {
		dst = make([]complex128, len(s))
	}

	var value T
	switch any(value).(type) {
	case int:
		for i, v := range s {
			dst[i] = complex(float64(v)/intToFloat, 0)
		}
	case int8:
		for i, v := range s {
			dst[i] = complex(float64(v)/int8ToFloat, 0)
		}
	case int16:
		for i, v := range s {
			dst[i] = complex(float64(v)/int16ToFloat, 0)
		}
	case int32:
		for i, v := range s {
			dst[i] = complex(float64(v)/int32ToFloat, 0)
		}
	case int64:
		for i, v := range s {
			dst[i] = complex(float64(v)/int64ToFloat, 0)
		}
	case uint:
		for i, v := range s {
			dst[i] = complex((float64(v)-intToFloat)/intToFloat, 0)
		}
	case uint8:
		for i, v := range s {
			dst[i] = complex((float64(v)-int8ToFloat)/int8ToFloat, 0)
		}
	case uint16:
		for i, v := range s {
			dst[i] = complex((float64(v)-int16ToFloat)/int16ToFloat, 0)
		}
	case uint32:
		for i, v := range s {
			dst[i] = complex((float64(v)-int32ToFloat)/int32ToFloat, 0)
		}
	case uint64:
		for i, v := range s {
			dst[i] = complex((float64(v)-int64ToFloat)/int64ToFloat, 0)
		}
	case float32, float64:
		for i, v := range s {
			dst[i] = complex(float64(v), 0)
		}
	default:
		panic(fmt.Sprintf("can't convert %T to []complex128", value))
	}

	return dst
}

// ToFloat converts the buffers in b to [][]float64 in dst.
//
// The provided buffer dst should be nil, or adequately sized to fit all samples in all channels.
func (b Buffer[T]) ToFloat(dst [][]float64) [][]float64 {
	if dst == nil {
		dst = make(Buffer[float64], len(b))
		for i := range b {
			dst[i] = make(Samples[float64], len(b[i]))
		}
	}

	for i := range b {
		Samples[T](b[0]).ToFloat(dst[i])
	}

	return dst
}

// ToMono converts a multi channel buffer to a mono channel by averaging the sounds per channel.
func (b Buffer[T]) ToMono(dst Samples[T]) int {
	var channels = b.Channels()
	if channels == 1 {
		// Fast path, only a single channel.
		return copy(dst, b[0])
	}

	var samples = b.Samples()
	for sample := 0; sample < samples; sample++ {
		var total float64
		for channel := 0; channel < channels; channel++ {
			total += float64(b[channel][sample])
		}
		dst[sample] = T(total / float64(channels))
	}

	return samples
}
