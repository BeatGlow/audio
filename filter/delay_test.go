package filter

import (
	"bytes"
	"testing"
	"time"

	"github.com/BeatGlow/audio"
)

func TestDelayShorterThanBufferLength(t *testing.T) {
	testData := make([]byte, 1024)
	for i := range testData {
		testData[i] = byte(i & 0xff)
	}

	wantData := make([]byte, 1024)
	copy(wantData[80:], testData) // 20ms equals 80 samples, which is less than our buffer length of 128

	reader := audio.NewReader[byte](bytes.NewBuffer(testData), nil)
	d, err := NewDelay(reader, 1, 8000, 10*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("reader: %s (%d samples)", d, len(d.(*Delay[byte]).Samples))

	buffer := make(audio.Samples[byte], 128)
	for i := 0; i < len(testData); i += len(buffer) {
		if _, err = d.ReadSamples(buffer); err != nil {
			t.Fatal(err)
		}
		//t.Logf("samples (%d): %#02v", n, buffer)

		for i, v := range wantData[i : i+len(buffer)] {
			if v != buffer[i] {
				t.Errorf("expected sample %d to be %#02x, got %#02x", i, v, buffer[i])
			}
		}
	}
}

func TestDelayEqualToBufferLength(t *testing.T) {
	testData := make([]byte, 1024)
	for i := range testData {
		testData[i] = byte(i & 0xff)
	}

	wantData := make([]byte, 1024)
	copy(wantData[128:], testData) // 16ms equals 128 samples, which is equal to our buffer length of 128

	reader := audio.NewReader[byte](bytes.NewBuffer(testData), nil)
	d, err := NewDelay(reader, 1, 8000, 16*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("reader: %s (%d samples)", d, len(d.(*Delay[byte]).Samples))

	buffer := make(audio.Samples[byte], 128)
	for i := 0; i < len(testData); i += len(buffer) {
		if _, err = d.ReadSamples(buffer); err != nil {
			t.Fatal(err)
		}
		//t.Logf("samples (%d): %#02v", n, buffer)

		for i, v := range wantData[i : i+len(buffer)] {
			if v != buffer[i] {
				t.Errorf("expected sample %d to be %#02x, got %#02x", i, v, buffer[i])
			}
		}
	}
}

func TestDelayLongerThanBufferLength(t *testing.T) {
	testData := make([]byte, 1024)
	for i := range testData {
		testData[i] = byte(i & 0xff)
	}

	wantData := make([]byte, 1024)
	copy(wantData[160:], testData) // 20ms equals 160 samples, which is more than our buffer length of 128

	reader := audio.NewReader[byte](bytes.NewBuffer(testData), nil)
	d, err := NewDelay(reader, 1, 8000, 20*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("reader: %s (%d samples)", d, len(d.(*Delay[byte]).Samples))

	buffer := make(audio.Samples[byte], 128)
	for i := 0; i < len(testData); i += len(buffer) {
		if _, err = d.ReadSamples(buffer); err != nil {
			t.Fatal(err)
		}
		//t.Logf("samples (%d): %#02v", n, buffer)

		for i, v := range wantData[i : i+len(buffer)] {
			if v != buffer[i] {
				t.Errorf("expected sample %d to be %#02x, got %#02x", i, v, buffer[i])
			}
		}
	}
}
