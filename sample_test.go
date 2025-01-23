package audio_test

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"

	"golang.org/x/exp/constraints"

	"github.com/BeatGlow/audio"
)

// Shared test vectors between int/uint/int32/uint32/int64/uint64
var (
	testSamplesInt32 = audio.Samples[int32]{
		math.MinInt32,
		math.MinInt32 + 1,
		math.MinInt32 + 2,
		0,
		math.MaxInt32 - 3,
		math.MaxInt32 - 2,
		math.MaxInt32 - 1,
		math.MaxInt32,
	}
	testEncodedInt32 = []byte{
		0x80, 0x00, 0x00, 0x00,
		0x80, 0x00, 0x00, 0x01,
		0x80, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00,
		0x7f, 0xff, 0xff, 0xfc,
		0x7f, 0xff, 0xff, 0xfd,
		0x7f, 0xff, 0xff, 0xfe,
		0x7f, 0xff, 0xff, 0xff,
	}
	testSamplesUint32 = audio.Samples[uint32]{
		0, 1, 2, 3,
		math.MaxUint32 - 3,
		math.MaxUint32 - 2,
		math.MaxUint32 - 1,
		math.MaxUint32,
	}
	testEncodedUint32 = []byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x03,
		0xff, 0xff, 0xff, 0xfc,
		0xff, 0xff, 0xff, 0xfd,
		0xff, 0xff, 0xff, 0xfe,
		0xff, 0xff, 0xff, 0xff,
	}
	testSamplesInt64 = audio.Samples[int64]{
		math.MinInt64,
		math.MinInt64 + 1,
		math.MinInt64 + 2,
		0,
		math.MaxInt64 - 3,
		math.MaxInt64 - 2,
		math.MaxInt64 - 1,
		math.MaxInt64,
	}
	testEncodedInt64 = []byte{
		0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfc,
		0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfd,
		0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe,
		0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}
	testSamplesUint64 = audio.Samples[uint64]{
		0, 1, 2, 3,
		math.MaxUint64 - 3,
		math.MaxUint64 - 2,
		math.MaxUint64 - 1,
		math.MaxUint64,
	}
	testEncodedUint64 = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfc,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfd,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}
)

const (
	testMeanInt32  = 268435455
	testRMSInt32   = 2008787011
	testMeanUint32 = 2147483647
	testRMSUint32  = 3037000498
	testMeanInt64  = 1152921504606846976
	testRMSInt64   = 8627674528165471232
	//testMeanUint64 = 9223372036854775808
	//testRMSUint64  = 13043817825332783104
)

/*
func TestInt(t *testing.T) {
	var (
		samples  audio.Samples[int]
		encoded  []byte
		shift    uint
		intSize  = 32 << (^uint(0) >> 63)
		min      = -1 << (intSize - 1)
		max      = 1<<(intSize-1) - 1
		wantRMS  int
		wantMean int
	)
	switch intSize {
	case 32:
		shift = 2
		encoded = testEncodedInt32
		samples = make(audio.Samples[int], len(testSamplesInt32))
		for i, v := range testSamplesInt32 {
			samples[i] = int(v)
		}
		wantRMS, wantMean = testRMSInt32, testMeanInt32
	case 64:
		shift = 3
		encoded = testEncodedInt64
		samples = make(audio.Samples[int], len(testSamplesInt64))
		for i, v := range testSamplesInt64 {
			samples[i] = int(v)
		}
		wantRMS, wantMean = testRMSInt64, testMeanInt64
	}
	t.Run("min", func(it *testing.T) {
		if v := audio.Min(samples); v != min {
			it.Fatalf("expected Min(%v) to return %d, got %d", samples, min, v)
		}
	})
	t.Run("max", func(it *testing.T) {
		if v := audio.Max(samples); v != max {
			it.Fatalf("expected Max(%v) to return %d, got %d", samples, max, v)
		}
	})
	t.Run("mean", func(it *testing.T) {
		if v := audio.Mean(samples); v != wantMean {
			it.Fatalf("expected Mean(%v) to return %d, got %d", samples, wantMean, v)
		}
	})
	t.Run("rms", func(it *testing.T) {
		if v := audio.RMS(samples); v != wantRMS {
			it.Fatalf("expected RMS(%v) to return %d, got %d", samples, wantRMS, v)
		}
	})
	t.Run("decode", func(it *testing.T) {
		test := make(audio.Samples[int], len(samples))
		test.Decode(encoded, binary.BigEndian)
		for i, v := range test {
			if v != samples[i] {
				it.Fatalf("expected values to decode to %d, got %d", samples, test)
			}
		}
	})
	t.Run("encode", func(it *testing.T) {
		test := make([]byte, len(samples)<<shift)
		samples.Encode(test, binary.BigEndian)
		for i, v := range test {
			if v != encoded[i] {
				it.Fatalf("expected values to encode to %#02v, got %#02v", encoded, test)
			}
		}
	})
}

func TestUint(t *testing.T) {
	var (
		samples  audio.Samples[uint]
		encoded  []byte
		shift    uint
		intSize       = 32 << (^uint(0) >> 63)
		max      uint = 1<<intSize - 1
		wantMean uint
		wantRMS  uint
	)
	switch intSize {
	case 32:
		shift = 2
		encoded = testEncodedUint32
		samples = make(audio.Samples[uint], len(testSamplesUint32))
		for i, v := range testSamplesUint32 {
			samples[i] = uint(v)
		}
		wantMean, wantRMS = testMeanUint32, testRMSUint32
	case 64:
		shift = 3
		encoded = testEncodedUint64
		samples = make(audio.Samples[uint], len(testSamplesUint64))
		for i, v := range testSamplesUint64 {
			samples[i] = uint(v)
		}
		wantMean, wantRMS = 9223372036854775808, 13043817825332783104
	}
	t.Run("min", func(it *testing.T) {
		if v := audio.Min(samples); v != 0 {
			it.Fatalf("expected Min(%v) to return 0, got %d", samples, v)
		}
	})
	t.Run("max", func(it *testing.T) {
		if v := audio.Max(samples); v != max {
			it.Fatalf("expected Max(%v) to return %d, got %d", samples, max, v)
		}
	})
	t.Run("mean", func(it *testing.T) {
		if v := audio.Mean(samples); v != wantMean {
			it.Fatalf("expected Mean(%v) to return %d, got %d", samples, wantMean, v)
		}
	})
	t.Run("rms", func(it *testing.T) {
		if v := audio.RMS(samples); v != wantRMS {
			it.Fatalf("expected RMS(%v) to return %d, got %d", samples, wantRMS, v)
		}
	})
	t.Run("decode", func(it *testing.T) {
		test := make(audio.Samples[uint], len(samples))
		test.Decode(encoded, binary.BigEndian)
		for i, v := range test {
			if v != samples[i] {
				it.Fatalf("expected values to decode to %d, got %d", samples, test)
			}
		}
	})
	t.Run("encode", func(it *testing.T) {
		test := make([]byte, len(samples)<<shift)
		samples.Encode(test, binary.BigEndian)
		for i, v := range test {
			if v != encoded[i] {
				it.Fatalf("expected values to encode to %#02v, got %#02v", encoded, test)
			}
		}
	})
}
*/

func TestInt8(t *testing.T) {
	samples := audio.Samples[int8]{0, 1, 2, 3, -4, -5, -6, -7}
	encoded := []byte{0, 1, 2, 3, 252, 251, 250, 249}
	t.Run("min", func(it *testing.T) {
		if v := audio.Min(samples); v != -7 {
			it.Fatalf("expected Min(%v) to return -7, got %d", samples, v)
		}
	})
	t.Run("max", func(it *testing.T) {
		if v := audio.Max(samples); v != 3 {
			it.Fatalf("expected Max(%v) to return 3, got %d", samples, v)
		}
	})
	t.Run("mean", func(it *testing.T) {
		if v := audio.Mean(samples); v != -2 {
			it.Fatalf("expected Mean(%v) to return -2, got %d", samples, v)
		}
	})
	t.Run("rms", func(it *testing.T) {
		if v := audio.RMS(samples); v != 4 {
			it.Fatalf("expected RMS(%v) to return 4, got %d", samples, v)
		}
	})
	t.Run("normalize", func(it *testing.T) {
		// Make a copy so we don't modify the test vector.
		test := make([]int8, len(samples))
		want := []int8{0, 0, 0, 0, 0, 0, 0, -1}
		copy(test, samples)
		audio.Normalize(test)
		for i, v := range test {
			if v != want[i] {
				it.Errorf("expected value %d to be %d, got %d", i, want[i], v)
			}
		}
	})
	t.Run("decode", func(it *testing.T) {
		test := make(audio.Samples[int8], len(samples))
		test.Decode(encoded, nil) // doesn't have a byte order
		for i, v := range test {
			if v != samples[i] {
				it.Errorf("expected value %d to decode to %d, got %d", i, samples[i], v)
			}
		}
	})
	t.Run("encode", func(it *testing.T) {
		test := make([]byte, len(samples))
		samples.Encode(test, nil) // doesn't have a byte order
		for i, v := range test {
			if v != encoded[i] {
				it.Fatalf("expected values to encode to %#02v, got %#02v", encoded, test)
			}
		}
	})
}

func TestUint8(t *testing.T) {
	samples := audio.Samples[uint8]{0, 1, 2, 3, 252, 253, 254, 255}
	encoded := []byte{0, 1, 2, 3, 252, 253, 254, 255}
	t.Run("min", func(it *testing.T) {
		if v := audio.Min(samples); v != 0 {
			it.Fatalf("expected Min(%v) to return 0, got %d", samples, v)
		}
	})
	t.Run("max", func(it *testing.T) {
		if v := audio.Max(samples); v != 255 {
			it.Fatalf("expected Max(%v) to return 255, got %d", samples, v)
		}
	})
	t.Run("mean", func(it *testing.T) {
		if v := audio.Mean(samples); v != 127 {
			it.Fatalf("expected Mean(%v) to return 127, got %d", samples, v)
		}
	})
	t.Run("rms", func(it *testing.T) {
		if v := audio.RMS(samples); v != 179 {
			it.Fatalf("expected RMS(%v) to return 179, got %d", samples, v)
		}
	})
	t.Run("normalize", func(it *testing.T) {
		// Make a copy so we don't modify the test vector.
		test := make([]uint8, len(samples))
		want := []uint8{0, 0, 0, 0, 0, 0, 0, 1}
		copy(test, samples)
		audio.Normalize(test)
		for i, v := range test {
			if v != want[i] {
				it.Errorf("expected value %d to be %d, got %d", i, want[i], v)
			}
		}
	})
	t.Run("decode", func(it *testing.T) {
		test := make(audio.Samples[uint8], len(samples))
		test.Decode(encoded, nil) // doesn't have a byte order
		for i, v := range test {
			if v != samples[i] {
				it.Fatalf("expected values to decode to %#02x, got %#02x", samples, test)
			}
		}
	})
	t.Run("encode", func(it *testing.T) {
		test := make([]byte, len(samples))
		samples.Encode(test, nil) // doesn't have a byte order
		for i, v := range test {
			if v != encoded[i] {
				it.Fatalf("expected values to encode to %#02v, got %#02v", encoded, test)
			}
		}
	})
}

func TestInt16(t *testing.T) {
	samples := audio.Samples[int16]{0, 1, 2, 3, -4, -5, -6, -7}
	encoded := []byte{
		0x00, 0x00,
		0x00, 0x01,
		0x00, 0x02,
		0x00, 0x03,
		0xff, 0xfc,
		0xff, 0xfb,
		0xff, 0xfa,
		0xff, 0xf9,
	}
	t.Run("min", func(it *testing.T) {
		if v := audio.Min(samples); v != -7 {
			it.Fatalf("expected Min(%v) to return -7, got %d", samples, v)
		}
	})
	t.Run("max", func(it *testing.T) {
		if v := audio.Max(samples); v != 3 {
			it.Fatalf("expected Max(%v) to return 3, got %d", samples, v)
		}
	})
	t.Run("mean", func(it *testing.T) {
		if v := audio.Mean(samples); v != -2 {
			it.Fatalf("expected Mean(%v) to return -2, got %d", samples, v)
		}
	})
	t.Run("rms", func(it *testing.T) {
		if v := audio.RMS(samples); v != 4 {
			it.Fatalf("expected RMS(%v) to return 4, got %d", samples, v)
		}
	})
	t.Run("decode", func(it *testing.T) {
		test := make(audio.Samples[int16], len(samples))
		test.Decode(encoded, binary.BigEndian)
		for i, v := range test {
			if v != samples[i] {
				it.Fatalf("expected values to decode to %d, got %d", samples, test)
			}
		}
	})
	t.Run("encode", func(it *testing.T) {
		test := make([]byte, len(samples)<<1)
		samples.Encode(test, binary.BigEndian)
		for i, v := range test {
			if v != encoded[i] {
				it.Fatalf("expected values to encode to %#02v, got %#02v", encoded, test)
			}
		}
	})
}

func TestUint16(t *testing.T) {
	samples := audio.Samples[uint16]{0, 1, 2, 3, 65532, 65533, 65534, 65535}
	encoded := []byte{
		0x00, 0x00,
		0x00, 0x01,
		0x00, 0x02,
		0x00, 0x03,
		0xff, 0xfc,
		0xff, 0xfd,
		0xff, 0xfe,
		0xff, 0xff,
	}
	t.Run("min", func(it *testing.T) {
		if v := audio.Min(samples); v != 0 {
			it.Fatalf("expected Min(%v) to return 0, got %d", samples, v)
		}
	})
	t.Run("max", func(it *testing.T) {
		if v := audio.Max(samples); v != 65535 {
			it.Fatalf("expected Max(%v) to return 65535, got %d", samples, v)
		}
	})
	t.Run("mean", func(it *testing.T) {
		if v := audio.Mean(samples); v != 32767 {
			it.Fatalf("expected Mean(%v) to return 32767, got %d", samples, v)
		}
	})
	t.Run("normalize", func(it *testing.T) {
		// Make a copy so we don't modify the test vector.
		test := make([]uint16, len(samples))
		want := []uint16{0, 0, 0, 0, 0, 0, 0, 1}
		copy(test, samples)
		audio.Normalize(test)
		for i, v := range test {
			if v != want[i] {
				it.Errorf("expected value %d to be %d, got %d", i, want[i], v)
			}
		}
	})
	t.Run("decode", func(it *testing.T) {
		test := make(audio.Samples[uint16], len(samples))
		test.Decode(encoded, binary.BigEndian)
		for i, v := range test {
			if v != samples[i] {
				it.Fatalf("expected values to decode to %d, got %d", samples, test)
			}
		}
	})
	t.Run("encode", func(it *testing.T) {
		test := make([]byte, len(samples)<<1)
		samples.Encode(test, binary.BigEndian)
		for i, v := range test {
			if v != encoded[i] {
				it.Fatalf("expected values to encode to %#02v, got %#02v", encoded, test)
			}
		}
	})
}

func TestInt32(t *testing.T) {
	samples := testSamplesInt32
	encoded := testEncodedInt32
	t.Run("min", func(it *testing.T) {
		if v := audio.Min(samples); v != math.MinInt32 {
			it.Fatalf("expected Min(%v) to return %d, got %d", samples, math.MinInt32, v)
		}
	})
	t.Run("max", func(it *testing.T) {
		if v := audio.Max(samples); v != math.MaxInt32 {
			it.Fatalf("expected Max(%v) to return %d, got %d", samples, math.MaxInt32, v)
		}
	})
	t.Run("mean", func(it *testing.T) {
		if v := audio.Mean(samples); v != testMeanInt32 {
			it.Fatalf("expected Mean(%v) to return %d got %d", samples, testMeanInt32, v)
		}
	})
	t.Run("rms", func(it *testing.T) {
		if v := audio.RMS(samples); v != testRMSInt32 {
			it.Fatalf("expected RMS(%v) to return %d, got %d", samples, testRMSInt32, v)
		}
	})
	t.Run("decode", func(it *testing.T) {
		test := make(audio.Samples[int32], len(samples))
		test.Decode(encoded, binary.BigEndian)
		for i, v := range test {
			if v != samples[i] {
				it.Fatalf("expected values to decode to %d, got %d", samples, test)
			}
		}
	})
	t.Run("encode", func(it *testing.T) {
		test := make([]byte, len(samples)<<2)
		samples.Encode(test, binary.BigEndian)
		for i, v := range test {
			if v != encoded[i] {
				it.Fatalf("expected values to encode to %#02v, got %#02v", encoded, test)
			}
		}
	})
}

func TestUint32(t *testing.T) {
	samples := testSamplesUint32
	encoded := testEncodedUint32
	t.Run("min", func(it *testing.T) {
		if v := audio.Min(samples); v != 0 {
			it.Fatalf("expected Min(%v) to return 0, got %d", samples, v)
		}
	})
	t.Run("max", func(it *testing.T) {
		if v := audio.Max(samples); v != math.MaxUint32 {
			it.Fatalf("expected Max(%v) to return 4294967295, got %d", samples, v)
		}
	})
	t.Run("mean", func(it *testing.T) {
		if v := audio.Mean(samples); v != testMeanUint32 {
			it.Fatalf("expected Mean(%v) to return 2147483647, got %d", samples, v)
		}
	})
	t.Run("rms", func(it *testing.T) {
		if v := audio.RMS(samples); v != testRMSUint32 {
			it.Fatalf("expected RMS(%v) to return 3037000498, got %d", samples, v)
		}
	})
	t.Run("normalize", func(it *testing.T) {
		// Make a copy so we don't modify the test vector.
		test := make([]uint32, len(samples))
		want := []uint32{0, 0, 0, 0, 0, 0, 0, 1}
		copy(test, samples)
		audio.Normalize(test)
		for i, v := range test {
			if v != want[i] {
				it.Errorf("expected value %d to be %d, got %d", i, want[i], v)
			}
		}
	})
	t.Run("decode", func(it *testing.T) {
		test := make(audio.Samples[uint32], len(samples))
		test.Decode(encoded, binary.BigEndian)
		for i, v := range test {
			if v != samples[i] {
				it.Fatalf("expected values to decode to %d, got %d", samples, test)
			}
		}
	})
	t.Run("encode", func(it *testing.T) {
		test := make([]byte, len(samples)<<2)
		samples.Encode(test, binary.BigEndian)
		for i, v := range test {
			if v != encoded[i] {
				it.Fatalf("expected values to encode to %#02v, got %#02v", encoded, test)
			}
		}
	})
}

func TestInt64(t *testing.T) {
	samples := testSamplesInt64
	encoded := testEncodedInt64
	t.Run("min", func(it *testing.T) {
		if v := audio.Min(samples); v != math.MinInt64 {
			it.Fatalf("expected Min(%v) to return -9223372036854775808, got %d", samples, v)
		}
	})
	t.Run("max", func(it *testing.T) {
		if v := audio.Max(samples); v != math.MaxInt64 {
			it.Fatalf("expected Max(%v) to return 9223372036854775807, got %d", samples, v)
		}
	})
	t.Run("mean", func(it *testing.T) {
		if v := audio.Mean(samples); v != 1152921504606846976 {
			it.Fatalf("expected Mean(%v) to return 1152921504606846976, got %d", samples, v)
		}
	})
	t.Run("rms", func(it *testing.T) {
		if v := audio.RMS(samples); v != 8627674528165471232 {
			it.Fatalf("expected RMS(%v) to return 8627674528165471232, got %d", samples, v)
		}
	})
	t.Run("decode", func(it *testing.T) {
		test := make(audio.Samples[int64], len(samples))
		test.Decode(encoded, binary.BigEndian)
		for i, v := range test {
			if v != samples[i] {
				it.Fatalf("expected values to decode to %d, got %d", samples, test)
			}
		}
	})
	t.Run("encode", func(it *testing.T) {
		test := make([]byte, len(samples)<<3)
		samples.Encode(test, binary.BigEndian)
		for i, v := range test {
			if v != encoded[i] {
				it.Fatalf("expected values to encode to %#02v, got %#02v", encoded, test)
			}
		}
	})
}

func TestUint64(t *testing.T) {
	samples := testSamplesUint64
	encoded := testEncodedUint64
	t.Run("min", func(it *testing.T) {
		if v := audio.Min(samples); v != 0 {
			it.Fatalf("expected Min(%v) to return 0, got %d", samples, v)
		}
	})
	t.Run("max", func(it *testing.T) {
		if v := audio.Max(samples); v != math.MaxUint64 {
			it.Fatalf("expected Max(%v) to return 18446744073709551615, got %d", samples, v)
		}
	})
	t.Run("mean", func(it *testing.T) {
		if v := audio.Mean(samples); v != 9223372036854775808 {
			it.Fatalf("expected Mean(%v) to return 9223372036854775808, got %d", samples, v)
		}
	})
	t.Run("rms", func(it *testing.T) {
		if v := audio.RMS(samples); v != 13043817825332783104 {
			it.Fatalf("expected RMS(%v) to return 13043817825332783104, got %d", samples, v)
		}
	})
	t.Run("normalize", func(it *testing.T) {
		// Make a copy so we don't modify the test vector.
		test := make([]uint64, len(samples))
		want := []uint64{0, 0, 0, 0, 0, 0, 0, 1}
		copy(test, samples)
		audio.Normalize(test)
		for i, v := range test {
			if v != want[i] {
				it.Errorf("expected value %d to be %d, got %d", i, want[i], v)
			}
		}
	})
	t.Run("decode", func(it *testing.T) {
		test := make(audio.Samples[uint64], len(samples))
		test.Decode(encoded, binary.BigEndian)
		for i, v := range test {
			if v != samples[i] {
				it.Fatalf("expected values to decode to %d, got %d", samples, test)
			}
		}
	})
	t.Run("encode", func(it *testing.T) {
		test := make([]byte, len(samples)<<3)
		samples.Encode(test, binary.BigEndian)
		for i, v := range test {
			if v != encoded[i] {
				it.Fatalf("expected values to encode to %#02v, got %#02v", encoded, test)
			}
		}
	})
}

func TestFloat32(t *testing.T) {
	samples := audio.Samples[float32]{0, 1, 2, 3, -4, -5, -6, -7}
	encoded := []byte{
		0x00, 0x00, 0x00, 0x00,
		0x3f, 0x80, 0x00, 0x00,
		0x40, 0x00, 0x00, 0x00,
		0x40, 0x40, 0x00, 0x00,
		0xc0, 0x80, 0x00, 0x00,
		0xc0, 0xa0, 0x00, 0x00,
		0xc0, 0xc0, 0x00, 0x00,
		0xc0, 0xe0, 0x00, 0x00,
	}
	t.Run("min", func(it *testing.T) {
		if v := audio.Min(samples); v != -7 {
			it.Fatalf("expected Min(%v) to return -7, got %f", samples, v)
		}
	})
	t.Run("max", func(it *testing.T) {
		if v := audio.Max(samples); v != 3 {
			it.Fatalf("expected Max(%v) to return 3, got %f", samples, v)
		}
	})
	t.Run("mean", func(it *testing.T) {
		if v := audio.Mean(samples); v != -2 {
			it.Fatalf("expected Mean(%v) to return -2, got %f", samples, v)
		}
	})
	t.Run("rms", func(it *testing.T) {
		if v := audio.RMS(samples); !testAlmostEqual(v, 4.183300) {
			it.Fatalf("expected RMS(%v) to return 4.183300, got %f", samples, v)
		}
	})
	t.Run("normalize", func(it *testing.T) {
		// Make a copy so we don't modify the test vector.
		test := make([]float32, len(samples))
		want := []float32{0, 0.142857, 0.285714, 0.428571, -0.571429, -0.714286, -0.857143, -1}
		copy(test, samples)
		audio.Normalize(test)
		for i, v := range test {
			if !testAlmostEqual(v, want[i]) {
				it.Errorf("expected value %d to be %f, got %f", i, want[i], v)
			}
		}
	})
	t.Run("decode", func(it *testing.T) {
		test := make(audio.Samples[float32], len(samples))
		test.Decode(encoded, binary.BigEndian)
		for i, v := range test {
			if v != samples[i] {
				it.Fatalf("expected values to decode to %f, got %f", samples, test)
			}
		}
	})
	t.Run("encode", func(it *testing.T) {
		test := make([]byte, len(samples)<<2)
		samples.Encode(test, binary.BigEndian)
		for i, v := range test {
			if v != encoded[i] {
				it.Fatalf("expected values to encode to %#02v, got %#02v", encoded, test)
			}
		}
	})
}

func TestFloat64(t *testing.T) {
	samples := audio.Samples[float64]{0, 1, 2, 3, -4, -5, -6, -7}
	encoded := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x40, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xc0, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xc0, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xc0, 0x18, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xc0, 0x1c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	t.Run("min", func(it *testing.T) {
		if v := audio.Min(samples); v != -7 {
			it.Fatalf("expected Min(%v) to return -7, got %f", samples, v)
		}
	})
	t.Run("max", func(it *testing.T) {
		if v := audio.Max(samples); v != 3 {
			it.Fatalf("expected Max(%v) to return 3, got %f", samples, v)
		}
	})
	t.Run("mean", func(it *testing.T) {
		if v := audio.Mean(samples); v != -2 {
			it.Fatalf("expected Mean(%v) to return -2, got %f", samples, v)
		}
	})
	t.Run("rms", func(it *testing.T) {
		if v := audio.RMS(samples); !testAlmostEqual(v, 4.183300) {
			it.Fatalf("expected RMS(%v) to return 4.183300, got %f", samples, v)
		}
	})
	t.Run("normalize", func(it *testing.T) {
		// Make a copy so we don't modify the test vector.
		test := make([]float64, len(samples))
		want := []float64{0, 0.142857, 0.285714, 0.428571, -0.571429, -0.714286, -0.857143, -1}
		copy(test, samples)
		audio.Normalize(test)
		for i, v := range test {
			if !testAlmostEqual(v, want[i]) {
				it.Errorf("expected value %d to be %f, got %f", i, want[i], v)
			}
		}
	})
	t.Run("decode", func(it *testing.T) {
		test := make(audio.Samples[float64], len(samples))
		test.Decode(encoded, binary.BigEndian)
		for i, v := range test {
			if v != samples[i] {
				it.Fatalf("expected values to decode to %f, got %f", samples, test)
			}
		}
	})
	t.Run("encode", func(it *testing.T) {
		test := make([]byte, len(samples)<<3)
		samples.Encode(test, binary.BigEndian)
		for i, v := range test {
			if v != encoded[i] {
				it.Fatalf("expected values to encode to %#02v, got %#02v", encoded, test)
			}
		}
	})
}

func TestInterleave(t *testing.T) {
	test := [][]byte{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{2, 2, 2, 2},
		{3, 3, 3, 3},
	}
	want := []byte{
		0, 1, 2, 3,
		0, 1, 2, 3,
		0, 1, 2, 3,
		0, 1, 2, 3,
	}

	v := make([]byte, 16)
	audio.Interleave(v, test)
	if len(v) != len(want) {
		t.Fatalf("expected %d samples returned, got %d", len(want), len(v))
		return
	}
	for i, s := range want {
		if v[i] != s {
			t.Errorf("expected sample %d to be %#02x, got %#02x", i, s, v[i])
		}
	}
}

func TestInterleaveNil(t *testing.T) {
	var test [][]byte
	if v := audio.Interleave(nil, test); v != nil {
		t.Fatalf("expected Interleave(nil) to return nil, got %#+v", v)
	}
}

func TestInterleaveDstNil(t *testing.T) {
	test := [][]byte{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{2, 2, 2, 2},
		{3, 3, 3, 3},
	}
	want := []byte{
		0, 1, 2, 3,
		0, 1, 2, 3,
		0, 1, 2, 3,
		0, 1, 2, 3,
	}

	v := audio.Interleave(nil, test)
	if len(v) != len(want) {
		t.Fatalf("expected %d samples returned, got %d", len(want), len(v))
		return
	}
	for i, s := range want {
		if v[i] != s {
			t.Errorf("expected sample %d to be %#02x, got %#02x", i, s, v[i])
		}
	}
}

func TestDeinterleave(t *testing.T) {
	test := []byte{
		0, 1, 2, 3,
		0, 1, 2, 3,
		0, 1, 2, 3,
		0, 1, 2, 3,
	}
	want := [][]byte{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{2, 2, 2, 2},
		{3, 3, 3, 3},
	}

	v := make([][]byte, 4)
	for i := 0; i < 4; i++ {
		v[i] = make([]byte, 4)
	}

	audio.Deinterleave(v, test, 4)
	if len(v) != len(want) {
		t.Fatalf("expected %d buffers returned, got %d", len(want), len(v))
		return
	}
	for c, vs := range want {
		if len(vs) != len(want[c]) {
			t.Errorf("expected channel %d to contain %d samples, got %d", c, len(want[c]), len(vs))
			continue
		}
		for i, s := range want[c] {
			if vs[i] != s {
				t.Errorf("expected sample %d to be %#02x, got %#02x", i, s, vs[i])
			}
		}
	}
}

func TestDeinterleaveMono(t *testing.T) {
	test := []byte{
		0, 1, 2, 3,
		0, 1, 2, 3,
		0, 1, 2, 3,
		0, 1, 2, 3,
	}
	want := [][]byte{test}

	v := make([][]byte, 1)
	v[0] = make([]byte, 16)
	audio.Deinterleave(v, test, 1)

	if len(v) != len(want) {
		t.Fatalf("expected %d buffers returned, got %d", len(want), len(v))
		return
	}

	if !bytes.Equal(v[0], want[0]) {
		t.Fatalf("expected %#02v, got %#02v", want, v)
	}
}

func TestDeinterleaveNil(t *testing.T) {
	var test []byte
	if v := audio.Deinterleave(nil, test, 0); v != nil {
		t.Fatalf("expected Deinterleave(nil) to return nil, got %#+v", v)
	}
}

func TestDeinterleaveDstNil(t *testing.T) {
	test := []byte{
		0, 1, 2, 3,
		0, 1, 2, 3,
		0, 1, 2, 3,
		0, 1, 2, 3,
	}
	want := [][]byte{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{2, 2, 2, 2},
		{3, 3, 3, 3},
	}

	v := audio.Deinterleave(nil, test, 4)
	if len(v) != len(want) {
		t.Fatalf("expected %d buffers returned, got %d", len(want), len(v))
		return
	}
	for c, vs := range want {
		if len(vs) != len(want[c]) {
			t.Errorf("expected channel %d to contain %d samples, got %d", c, len(want[c]), len(vs))
			continue
		}
		for i, s := range want[c] {
			if vs[i] != s {
				t.Errorf("expected sample %d to be %#02x, got %#02x", i, s, vs[i])
			}
		}
	}
}

func TestMinZero(t *testing.T) {
	if v := audio.Min([]int{}); v != 0 {
		t.Fatalf("expected Min(empty) to return 0, got %d", v)
	}
}

func TestMaxZero(t *testing.T) {
	if v := audio.Max([]int{}); v != 0 {
		t.Fatalf("expected Max(empty) to return 0, got %d", v)
	}
}

func TestMeanZero(t *testing.T) {
	if v := audio.Mean([]int{}); v != 0 {
		t.Fatalf("expected Mean(empty) to return 0, got %d", v)
	}
}

func TestNormalizeZero(t *testing.T) {
	audio.Normalize([]int{})
}

func TestRMSZero(t *testing.T) {
	if v := audio.RMS([]int{}); v != 0 {
		t.Fatalf("expected RMS(empty) to return 0, got %d", v)
	}
}

func TestSampleBitsPerSample(t *testing.T) {
	const intSize = 32 << (^uint(0) >> 63)
	t.Run("int", func(it *testing.T) {
		var test audio.Samples[int]
		if v := test.BitsPerSample(); v != intSize {
			t.Fatalf("expected %T to contain %d bits per sample, got %d", test, intSize, v)
		}
	})
	t.Run("int8", func(it *testing.T) {
		var test audio.Samples[int8]
		if v := test.BitsPerSample(); v != 8 {
			t.Fatalf("expected %T to contain 8 bits per sample, got %d", test, v)
		}
	})
	t.Run("int16", func(it *testing.T) {
		var test audio.Samples[int16]
		if v := test.BitsPerSample(); v != 16 {
			t.Fatalf("expected %T to contain 16 bits per sample, got %d", test, v)
		}
	})
	t.Run("int32", func(it *testing.T) {
		var test audio.Samples[int32]
		if v := test.BitsPerSample(); v != 32 {
			t.Fatalf("expected %T to contain 32 bits per sample, got %d", test, v)
		}
	})
	t.Run("int64", func(it *testing.T) {
		var test audio.Samples[int64]
		if v := test.BitsPerSample(); v != 64 {
			t.Fatalf("expected %T to contain 64 bits per sample, got %d", test, v)
		}
	})
	t.Run("uint", func(it *testing.T) {
		var test audio.Samples[uint]
		if v := test.BitsPerSample(); v != intSize {
			t.Fatalf("expected %T to contain %d bits per sample, got %d", test, intSize, v)
		}
	})
	t.Run("uint8", func(it *testing.T) {
		var test audio.Samples[uint8]
		if v := test.BitsPerSample(); v != 8 {
			t.Fatalf("expected %T to contain 8 bits per sample, got %d", test, v)
		}
	})
	t.Run("uint16", func(it *testing.T) {
		var test audio.Samples[uint16]
		if v := test.BitsPerSample(); v != 16 {
			t.Fatalf("expected %T to contain 16 bits per sample, got %d", test, v)
		}
	})
	t.Run("uint32", func(it *testing.T) {
		var test audio.Samples[uint32]
		if v := test.BitsPerSample(); v != 32 {
			t.Fatalf("expected %T to contain 32 bits per sample, got %d", test, v)
		}
	})
	t.Run("uint64", func(it *testing.T) {
		var test audio.Samples[uint64]
		if v := test.BitsPerSample(); v != 64 {
			t.Fatalf("expected %T to contain 64 bits per sample, got %d", test, v)
		}
	})
	t.Run("float32", func(it *testing.T) {
		var test audio.Samples[float32]
		if v := test.BitsPerSample(); v != 32 {
			t.Fatalf("expected %T to contain 32 bits per sample, got %d", test, v)
		}
	})
	t.Run("float64", func(it *testing.T) {
		var test audio.Samples[float64]
		if v := test.BitsPerSample(); v != 64 {
			t.Fatalf("expected %T to contain 64 bits per sample, got %d", test, v)
		}
	})
}

func testAlmostEqual[T constraints.Float](a, b T) bool {
	const testE = 1e-6
	return audio.Abs(a-b) < T(testE)
}
