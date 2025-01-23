package alsa

import (
	"fmt"
	"strings"
)

type Access int

const (
	MmapInterleaved Access = iota
	MmapNonInterleaved
	MmapComplex
	RWInterleaved
	RWNonInterleaved
	AccessTypeLast  = RWNonInterleaved
	AccessTypeFirst = MmapInterleaved
)

func (a Access) String() string {
	switch a {
	case MmapInterleaved:
		return "MmapInterleaved"
	case MmapNonInterleaved:
		return "MmapNonInterleaved"
	case MmapComplex:
		return "MmapComplex"
	case RWInterleaved:
		return "RWInterleaved"
	case RWNonInterleaved:
		return "RWNonInterleaved"
	default:
		return fmt.Sprintf("Invalid AccessType (%d)", a)
	}
}

type sampleFormat int

const (
	Unknown sampleFormat = -1
)
const (
	formatInt8 sampleFormat = iota
	formatUint8
	formatInt16
	formatInt16BE
	formatUint16
	formatUint16BE
	formatInt24
	formatInt24BE
	formatUint24
	formatUint24BE
	formatInt32
	formatInt32BE
	formatUint32
	formatUint32BE
	formatFloat32
	formatFloat32BE
	formatFloat64
	formatFloat64BE
	// There are so many more...
	formatTypeFirst = formatInt8
	formatTypeLast  = formatFloat64BE
)

func (f sampleFormat) BitsPerSample() int {
	switch f {
	case formatInt8,
		formatUint8:
		return 8
	case formatInt16,
		formatInt16BE,
		formatUint16,
		formatUint16BE:
		return 16
	case formatInt24,
		formatInt24BE,
		formatUint24,
		formatUint24BE:
		return 24
	case formatInt32,
		formatInt32BE,
		formatUint32,
		formatUint32BE,
		formatFloat32,
		formatFloat32BE:
		return 32
	case formatFloat64,
		formatFloat64BE:
		return 64
	default:
		return 0
	}
}

func (f sampleFormat) String() string {
	switch f {
	case formatInt8:
		return "int8"
	case formatUint8:
		return "uint8"
	case formatInt16:
		return "int16"
	case formatUint16:
		return "uint16"
	case formatInt16BE:
		return "int16be"
	case formatUint16BE:
		return "uint16be"
	case formatInt24:
		return "int24"
	case formatUint24:
		return "uint24"
	case formatInt24BE:
		return "int24be"
	case formatUint24BE:
		return "uint24be"
	case formatInt32:
		return "int32"
	case formatInt32BE:
		return "int32be"
	case formatUint32:
		return "uint32"
	case formatUint32BE:
		return "uint32be"
	case formatFloat32:
		return "float32"
	case formatFloat32BE:
		return "float32be"
	case formatFloat64:
		return "float64"
	case formatFloat64BE:
		return "float64be"
	default:
		return fmt.Sprintf("unsupported type (%d)", f)
	}
}

type SubformatType int

const (
	StandardSubformat  SubformatType = iota
	SubformatTypeFirst               = StandardSubformat
	SubformatTypeLast                = StandardSubformat
)

func (f SubformatType) String() string {
	switch f {
	case StandardSubformat:
		return "standard"
	default:
		return fmt.Sprintf("unknown subformat type (%d)", f)
	}
}

type Version uint32

func (v Version) Major() int { return int(v>>16) & 0xffff }
func (v Version) Minor() int { return int(v>>8) & 0xff }
func (v Version) Patch() int { return int(v) & 0xff }

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major(), v.Minor(), v.Patch())
}

type Flags uint32

const (
	OpenMin Flags = 1 << iota
	OpenMax
	Integer
	Empty
)

func (f Flags) String() string {
	r := ""
	if f&OpenMin != 0 {
		r += "OpenMin "
	}
	if f&OpenMax != 0 {
		r += "OpenMax "
	}
	if f&Integer != 0 {
		r += "Integer "
	}
	if f&Empty != 0 {
		r += "Empty "
	}
	return strings.TrimSpace(r)
}

type Timespec struct {
	Sec  int
	Nsec int
}

type swParams struct {
	TstampMode       int32
	PeriodStep       uint32
	SleepMin         uint32
	AvailMin         uFrames
	XferAlign        uFrames
	StartThreshold   uFrames
	StopThreshold    uFrames
	SilenceThreshold uFrames
	SilenceSize      uFrames
	Boundary         uFrames
	Proto            Version
	TstampType       uint32
	Reserved         [56]byte
}

func cstr(b []byte) string {
	for i, v := range b {
		if v == 0x00 {
			return string(b[:i])
		}
	}
	return string(b)
}
