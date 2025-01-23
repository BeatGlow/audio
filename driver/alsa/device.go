package alsa

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/BeatGlow/audio"
)

type DeviceType int

const (
	UnknownDeviceType DeviceType = iota
	PCM
)

func (t DeviceType) String() string {
	switch t {
	case PCM:
		return "PCM"
	default:
		return fmt.Sprintf("unknown type %d", t)
	}
}

type Device struct {
	Type      DeviceType
	Index     int
	Path      string
	Name      string
	CanPlay   bool
	CanRecord bool

	fh       *os.File
	fd       uintptr
	info     pcmInfo
	version  Version
	hwParams hwParams
	swParams swParams
	ready    bool
}

func (dev Device) String() string {
	return dev.Name
}

func (dev *Device) Close() error {
	return dev.fh.Close()
}

func (dev *Device) Open() error {
	var err error
	if dev.fh, err = os.OpenFile(dev.Path, os.O_RDWR, 0644); err != nil {
		return err
	}
	dev.fd = dev.fh.Fd()

	/*
		if err = ioctl(dev.fd, ioctlEncodePointer(cmdRead, &dev.version, cmdPCMVersion), uintptr(unsafe.Pointer(&dev.version))); err != nil {
			_ = dev.fh.Close()
			return err
		}

		timestamp := uint32(pcmTimestampTypeGettimeofday)
		if err = ioctl(dev.fd, ioctlEncodePointer(cmdRead, &timestamp, cmdPCMTimestampType), uintptr(unsafe.Pointer(&timestamp))); err != nil {
			_ = dev.fh.Close()
			return err
		}
	*/

	dev.hwParams = makeHwParams()
	if err = dev.refineHwParams(); err != nil {
		_ = dev.fh.Close()
		return err
	}

	dev.hwParams.Cmask = 0
	dev.hwParams.Rmask = 0xffffffff
	dev.hwParams.SetAccess(RWInterleaved)
	if err := dev.refineHwParams(); err != nil {
		return err
	}

	return nil
}

func (dev *Device) refineHwParams() (err error) {
	//return ioctl(dev.fd, ioctlEncodePointer(cmdRead|cmdWrite, &dev.hwParams, cmdPCMHwRefine), uintptr(unsafe.Pointer(&dev.hwParams)))
	return nil
}

func (dev *Device) updateSwParams() (err error) {
	//return ioctl(dev.fd, ioctlEncodePointer(cmdRead|cmdWrite, &dev.swParams, cmdPCMSwParams), uintptr(unsafe.Pointer(&dev.swParams)))
	return nil
}

func (dev *Device) bytesPerFrame() int {
	var (
		ss = int(dev.hwParams.Intervals[paramSampleBits-paramFirstInterval].Max) / 8
		ch = int(dev.hwParams.Intervals[paramChannels-paramFirstInterval].Max)
	)
	return ss * ch
}

func (dev *Device) prepare() error {
	bufSize := int(dev.hwParams.Intervals[paramBufferSize-paramFirstInterval].Max)

	dev.swParams = swParams{
		PeriodStep:     1,
		AvailMin:       uFrames(bufSize),
		XferAlign:      1,
		StartThreshold: uFrames(bufSize),
		StopThreshold:  uFrames(bufSize * 2),
		Proto:          dev.version,
		TstampType:     1,
	}
	if err := dev.updateSwParams(); err != nil {
		return err
	}
	/*
		if err := ioctl(dev.fd, ioctlEncode(0, 0, cmdPCMPrepare), 0); err != nil {
			return err
		}
	*/

	dev.ready = true
	return nil
}

type format struct {
	sampleFormat
	channels int
	rate     float64
}

func (f format) Channels() int       { return f.channels }
func (f format) SampleRate() float64 { return f.rate }

func (dev *Device) BufferFormat() (audio.Format, error) {
	if !dev.ready {
		if err := dev.prepare(); err != nil {
			return format{}, err
		}
	}

	var s sampleFormat
	for s = formatTypeFirst; s <= formatTypeLast; s++ {
		if s >= formatInt24 && s <= formatUint24BE {
			// Go doesn't have a native 24-bit type, skip (for now)
			continue
		}
		if dev.hwParams.GetFormatSupport(s) {
			break
		}
	}

	ch, _ := dev.hwParams.IntervalRange(paramChannels)
	rt, _ := dev.hwParams.IntervalRange(paramRate)

	return format{
		sampleFormat: s,
		channels:     int(ch),
		rate:         float64(rt),
	}, nil
}

func (dev *Device) Read(p []byte) (n int, err error) {
	n = len(p) / dev.bytesPerFrame()
	x := xferI{
		Buf:    uintptr(unsafe.Pointer(&p[0])),
		Frames: uFrames(n),
	}
	_ = x
	return
	//return n, ioctl(dev.fd, ioctlEncodePointer(cmdRead, &x, cmdPCMReadIFrames), uintptr(unsafe.Pointer(&x)))
}

func (dev *Device) Write(p []byte) (n int, err error) {
	n = len(p) / dev.bytesPerFrame()
	x := xferI{
		Buf:    uintptr(unsafe.Pointer(&p[0])),
		Frames: uFrames(n),
	}
	_ = x
	return
	//return n, ioctl(dev.fd, ioctlEncodePointer(cmdWrite, &x, cmdPCMWriteIFrames), uintptr(unsafe.Pointer(&x)))
}
