package alsa

import (
	"fmt"
	"reflect"
	"syscall"
)

//nolint:unused
const (
	cmdWrite = 1
	cmdRead  = 2

	cmdPCMInfo              uintptr = 0x4101
	cmdPCMVersion           uintptr = 0x4100
	cmdPCMTimestamp         uintptr = 0x4102
	cmdPCMTimestampType     uintptr = 0x4103
	cmdPCMHwRefine          uintptr = 0x4110
	cmdPCMHwParams          uintptr = 0x4111
	cmdPCMSwParams          uintptr = 0x4113
	cmdPCMStatus            uintptr = 0x4120
	cmdPCMPrepare           uintptr = 0x4140
	cmdPCMReset             uintptr = 0x4141
	cmdPCMStart             uintptr = 0x4142
	cmdPCMDrop              uintptr = 0x4143
	cmdPCMDrain             uintptr = 0x4144
	cmdPCMPause             uintptr = 0x4145 // int
	cmdPCMRewind            uintptr = 0x4146 // snd_pcm_uframes_t
	cmdPCMResume            uintptr = 0x4147
	cmdPCMXrun              uintptr = 0x4148
	cmdPCMForward           uintptr = 0x4149
	cmdPCMWriteIFrames      uintptr = 0x4150 // snd_xferi
	cmdPCMReadIFrames       uintptr = 0x4151 // snd_xferi
	cmdPCMWriteNFrames      uintptr = 0x4152 // snd_xfern
	cmdPCMReadNFrames       uintptr = 0x4153 // snd_xfern
	cmdPCMLink              uintptr = 0x4160 // int
	cmdPCMUnlink            uintptr = 0x4161
	cmdControlVersion       uintptr = 0x5500
	cmdControlCardInfo      uintptr = 0x5501
	cmdControlPCMNextDevice uintptr = 0x5530
	cmdControlPCMInfo       uintptr = 0x5531
)

//nolint:unused
const (
	pcmTimestampTypeGettimeofday = iota
	pcmTimestampTypeMonotonic
	pcmTimestampTypeMonotonicRaw
	pcmTimestampTypeLast
)

type ioctlCommand uintptr

func (c ioctlCommand) String() string {
	var (
		mode = c >> 30 & 0x03
		size = c >> 16 & 0x3fff
		cmd  = c & 0xffff
		str  string
	)
	if mode&cmdWrite > 0 {
		str += " write"
	}
	if mode&cmdRead > 0 {
		str += " read "
	}
	return fmt.Sprintf("ioctl%s (%d bytes) 0x%04x", str, size, uintptr(cmd))
}

func ioctl(fd uintptr, command ioctlCommand, ptr interface{}) error {
	var p uintptr

	if ptr != nil {
		v := reflect.ValueOf(ptr)
		p = v.Pointer()
	}

	//fmt.Printf("%s :: %d bytes\n", c, reflect.TypeOf(ptr).Elem().Size())
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(command), p)
	if e != 0 {
		return fmt.Errorf("ioctl %s failed: %v", command, e)
	}
	return nil
}

func ioctlEncode(mode byte, size uint16, cmd uintptr) ioctlCommand {
	return ioctlCommand(mode)<<30 | ioctlCommand(size)<<16 | ioctlCommand(cmd)
}

func ioctlPointer(mode byte, ref interface{}, cmd uintptr) ioctlCommand {
	return ioctlEncode(mode, uint16(reflect.TypeOf(ref).Elem().Size()), cmd)
}
