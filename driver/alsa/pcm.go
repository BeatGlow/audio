package alsa

type xferI struct {
	Result sFrames
	Buf    uintptr
	Frames uFrames
}

//const xferISize = 24
