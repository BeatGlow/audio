package audio_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/BeatGlow/audio"
)

func TestDecodeFrom(t *testing.T) {
	testCases := []struct {
		Name      string
		TestBytes []byte
		Order     binary.ByteOrder
		Want      audio.Samples[int16]
		WantError error
		Samples   int
	}{
		{
			"empty",
			nil,
			binary.BigEndian,
			nil,
			io.ErrShortBuffer,
			0,
		},
		{
			"single",
			[]byte{0x7f, 0xff},
			binary.BigEndian,
			audio.Samples[int16]{0x7fff},
			nil,
			1,
		},
		{
			"random",
			[]byte{
				0xe8, 0x48, 0xd8, 0xc2, 0xd6, 0x07, 0x86, 0x29,
				0x2e, 0xe2, 0xcd, 0xb9, 0x83, 0x2c, 0x71, 0x4a,
				0x05, 0x67, 0x13, 0x7b, 0x9e, 0xc5, 0x1d, 0xf2,
				0x67, 0xf3, 0xa3, 0xee, 0xfe, 0x09, 0x5a, 0x5b,
				0x24, 0x08, 0x56, 0x0d, 0x00, 0xed, 0xfe, 0x98,
				0x53, 0xfb, 0xc7, 0x86, 0x74, 0x14, 0xc4, 0x04,
				0xd7, 0x4a, 0x28, 0xb9, 0x6f, 0xf7, 0x45, 0x21,
				0xc7, 0x1f, 0x3a, 0x84, 0xdd, 0x0e, 0x32, 0xd1,
			},
			binary.BigEndian,
			audio.Samples[int16]{
				-6072, -10046, -10745, -31191,
				12002, -12871, -31956, 29002,
				1383, 4987, -24891, 7666,
				26611, -23570, -503, 23131,
				9224, 22029, 237, -360,
				21499, -14458, 29716, -15356,
				-10422, 10425, 28663, 17697,
				-14561, 14980, -8946, 13009,
			},
			nil,
			32,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(it *testing.T) {
			buffer := make(audio.Samples[int16], test.Samples)

			n, err := buffer.DecodeFrom(bytes.NewBuffer(test.TestBytes), test.Order)
			if test.WantError != nil {
				if err != test.WantError {
					it.Fatalf("expected error %q, got %q", test.WantError, err)
				}
				return
			} else if err != nil {
				it.Fatal(err)
			}

			if n != test.Samples {
				it.Errorf("expected %d samples, got %d", test.Samples, n)
			}
			for i, v := range buffer {
				if v != test.Want[i] {
					t.Errorf("expected sample %d to be %d, got %d", i, test.Want[i], v)
				}
			}
		})
	}
}
