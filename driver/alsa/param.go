package alsa

import "fmt"

type param uint32

const (
	paramAccess    param = 0
	paramFormat    param = 1
	paramSubformat param = 2
	paramFirstMask param = paramAccess
	paramLastMask  param = paramSubformat

	paramSampleBits    param = 8
	paramFrameBits     param = 9
	paramChannels      param = 10
	paramRate          param = 11
	paramPeriodTime    param = 12
	paramPeriodSize    param = 13
	paramPeriodBytes   param = 14
	paramPeriods       param = 15
	paramBufferTime    param = 16
	paramBufferSize    param = 17
	paramBufferBytes   param = 18
	paramTickTime      param = 19
	paramFirstInterval param = paramSampleBits
	paramLastInterval  param = paramTickTime
)

const (
	maskMax = 256
)

type mask struct {
	Bits [(maskMax + 31) / 32]uint32
}

type interval struct {
	Min, Max uint32
	Flags    Flags
}

func (i interval) String() string {
	return fmt.Sprintf("Interval(%d/%d 0x%x)", i.Min, i.Max, i.Flags)
}

type hwParams struct {
	Flags     uint32
	Masks     [paramLastMask - paramFirstMask + 1]mask
	_         [5]mask
	Intervals [paramLastInterval - paramFirstInterval + 1]interval
	_         [9]interval
	Rmask     uint32
	Cmask     uint32
	Info      uint32
	Msbits    uint32
	RateNum   uint32
	RateDen   uint32
	FifoSize  uFrames
	_         [64]byte
}

func makeHwParams() hwParams {
	var p hwParams
	for i := range p.Masks {
		for ii := 0; ii < 2; ii++ {
			p.Masks[i].Bits[ii] = 0xffffffff
		}
	}
	for i := range p.Intervals {
		p.Intervals[i].Max = 0xffffffff
	}
	p.Rmask = 0xffffffff
	return p
}

func (p *hwParams) SetAccess(a Access) {
	p.SetMask(paramAccess, uint32(1<<uint(a)))
}
func (p *hwParams) SetFormat(f sampleFormat) {
	p.SetMask(paramFormat, uint32(1<<uint(f)))
}
func (p *hwParams) SetMask(param param, v uint32) {
	p.Masks[param-paramFirstMask].Bits[0] = v
}
func (p *hwParams) GetFormatSupport(f sampleFormat) bool {
	bits := p.Masks[paramFormat-paramFirstMask].Bits[0]
	b := bits & (1 << uint(f))
	return b != 0
}

func (p *hwParams) SetInterval(param param, min, max uint32, flags Flags) {
	p.Intervals[param-paramFirstInterval].Min = min
	p.Intervals[param-paramFirstInterval].Max = max
	p.Intervals[param-paramFirstInterval].Flags = flags
}
func (p *hwParams) SetIntervalToMin(param param) {
	p.Intervals[param-paramFirstInterval].Max = p.Intervals[param-paramFirstInterval].Min
}
func (p *hwParams) IntervalInRange(param param, v uint32) bool {
	min, max := p.IntervalRange(param)
	if min > v {
		return false
	}
	if max < v {
		return false
	}
	return true
}

func (p *hwParams) IntervalRange(param param) (uint32, uint32) {
	return p.Intervals[param-paramFirstInterval].Min, p.Intervals[param-paramFirstInterval].Max
}
