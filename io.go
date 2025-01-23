package audio

import (
	"encoding/binary"
	"io"
	"log"
)

// Reader can read samples.
type Reader[T Sample] interface {
	ReadSamples(Samples[T]) (int, error)
}

type reader[T Sample] struct {
	r     io.Reader
	order binary.ByteOrder
}

// NewReader returns a Reader that can read from any io.Reader.
func NewReader[T Sample](r io.Reader, order binary.ByteOrder) Reader[T] {
	return reader[T]{
		r:     r,
		order: order,
	}
}

func (r reader[T]) ReadSamples(samples Samples[T]) (int, error) {
	i, err := samples.DecodeFrom(r.r, r.order)
	return int(i), err
}

// Writer can write samples.
type Writer[T Sample] interface {
	WriteSamples(Samples[T]) (int, error)
}

type writer[T Sample] struct {
	w     io.Writer
	order binary.ByteOrder
}

// NewWriter returns a Writer that can write to any io.Writer.
func NewWriter[T Sample](w io.Writer, order binary.ByteOrder) Writer[T] {
	return writer[T]{
		w:     w,
		order: order,
	}
}

func (w writer[T]) WriteSamples(samples Samples[T]) (int, error) {
	i, err := samples.EncodeTo(w.w, w.order)
	return int(i), err
}

// DecodeFrom reads samples from r to s.
func (s Samples[T]) DecodeFrom(r io.Reader, order binary.ByteOrder) (n int, err error) {
	return s.DecodeFromChunked(r, order, len(s))
}

// DecodeFromChunked reads samples from r to s.
//
// If s doesn't contain a multiple of chunkSize samples, an additional smaller chunk will be read
// to complete to read.
func (s Samples[T]) DecodeFromChunked(r io.Reader, order binary.ByteOrder, chunkSize int) (n int, err error) {
	if len(s) == 0 || chunkSize < 1 {
		return 0, io.ErrShortBuffer
	}

	var (
		samples        = len(s)
		bytesPerSample = s.BitsPerSample() / 8
		bytesPerChunk  = bytesPerSample * chunkSize
		buf            = make([]byte, bytesPerChunk)
	)

	// Read chunks of samples.
	for ; n < samples; n += chunkSize {
		if _, err = r.Read(buf); err != nil {
			return
		}
		log.Printf("decode %T %d:%d from %d bytes", s, n, n+chunkSize, len(buf))
		s[n:n+chunkSize].Decode(buf, order)
	}

	// Read remaining bytes if chunkSize doesn't align with the number of samples.
	if remain := samples % chunkSize; remain > 0 {
		buf = buf[:bytesPerSample*int(remain)]
		if _, err = r.Read(buf); err != nil {
			return
		}
		log.Printf("decode %T %d:%d from %d bytes remaining", s, n, n+remain, len(buf))
		s[n:n+remain].Decode(buf, order)
		n += remain
	}

	return
}

// Write samples contained in src to w.
func (s Samples[T]) EncodeTo(w io.Writer, order binary.ByteOrder) (n int, err error) {
	return s.EncodeToChunked(w, order, len(s))
}

// WriteChunked writes samples contained in src to w.
func (s Samples[T]) EncodeToChunked(w io.Writer, order binary.ByteOrder, chunkSize int) (n int, err error) {
	if len(s) == 0 || chunkSize < 1 {
		return
	}

	var (
		samples        = len(s)
		bytesPerSample = s.BitsPerSample() / 8
		bytesPerChunk  = bytesPerSample * chunkSize
		buf            = make([]byte, bytesPerChunk)
	)

	// Write chunks of samples.
	for ; n < samples; n += chunkSize {
		s[n:n+chunkSize].Encode(buf, order)
		if _, err = w.Write(buf); err != nil {
			return
		}
	}

	// Write remaining bytes if chunkSize doesn't align with the number of samples.
	if remain := samples % chunkSize; remain > 0 {
		buf = buf[:bytesPerSample*int(remain)]
		s[n:].Encode(buf, order)
		if _, err = w.Write(buf); err != nil {
			return
		}
		n += remain
	}

	return
}

/*
// BufferReader can read deinterleaved audio samples.
type BufferWriter[T Sample] interface {
	WriteBuffer(Buffer[T]) (int, error)
}

type bufferWriter[T Sample] struct {
	w     io.Writer
	order binary.ByteOrder
}

func NewBufferWriter[T Sample](w io.Writer, order binary.ByteOrder) BufferWriter[T] {
	return bufferWriter[T]{
		w:     w,
		order: order,
	}
}

func (w bufferWriter[T]) WriteBuffer(buffer Buffer[T]) (int, error) {
	switch len(buffer) {
	case 0:
		// Fast path, nothing to do.
		return 0, nil

	case 1:
		// Fast path, single buffer.
		return Samples[T](buffer[0]).EncodeTo(w.w, w.order)

	default:
	}
}
*/
