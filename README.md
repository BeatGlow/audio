# Go package for audio processing

This package makes use of Golang's [generics](https://go.dev/blog/intro-generics) for processing
audio. Don't be initimidated by their appearance :-)

```go
package main

import (
    "encoding/binary"
    "io"
    "os"

    "github.com/BeatGlow/audio"
)

func main() {
    f, err := os.Open("/tmp/mpd.fifo")
    if err != nil {
        panic(err)
    }
    defer func() { _ = f.Close() }()

    // Allocate a buffer of 1024 flaot64 samples.
    samples := make(audio.Samples[float64], 1024)

    r := audio.NewReader[float64](f, binary.BigEndian)
    for {
        n, err := r.ReadSamples(samples)
        if err == io.EOF {
            break
        } else if err != nil {
            panic(err)
        }
        println("Read", n, "samples")
    }
}
```