package bytearray

import (
	"bytearray/bufferflag"
	"unsafe"
)

const (
	MAXIMUM_PACKAGE = 256 * 1024
	HEADERLENGTH    = 5
)

type ByteArray struct {
	Buffer      []byte
	ReaderIdx   uint32
	WriterIdx   uint32
	CRC         uint32
	HeaderFlags bufferflag.BufferFlag
	Capacity    uint32
}

func CreateByteArray() *ByteArray {
	ba := &ByteArray{make([]byte, MAXIMUM_PACKAGE), 0, 0, 0, 0, MAXIMUM_PACKAGE}
	return ba
}

func (ba *ByteArray) Length() uint32 {
	return ba.WriterIdx
}

func (ba *ByteArray) RefreshHeader() {
	lp := unsafe.Pointer(&ba.WriterIdx)
	ba.Buffer[0] = *(*byte)(lp)
	var lpp uintptr = uintptr(lp)
	ba.Buffer[1] = *(*byte)(unsafe.Pointer(lpp + 1))
	ba.Buffer[2] = *(*byte)(unsafe.Pointer(lpp + 2))
	ba.Buffer[3] = *(*byte)(unsafe.Pointer(lpp + 3))
	ba.Buffer[4] = byte(ba.HeaderFlags)
}

func (ba *ByteArray) ReadHeader() {
	l := *(*uint32)(unsafe.Pointer(&ba.Buffer[0]))
	if l != ba.WriterIdx {
		panic("buffer length not equal WriterIdx")
	}
	ba.GetHeaderFlags()
}

func (ba *ByteArray) BypassHeader() {
	ba.ReaderIdx = HEADERLENGTH
	ba.GetHeaderFlags()
}

func (ba *ByteArray) AdvanceReader(steps uint32) {
	ba.ReaderIdx += steps
}

func (ba *ByteArray) AdvanceWriter(steps uint32) {
	ba.WriterIdx += steps
}

func (ba *ByteArray) GetHeaderFlags() {
	ba.HeaderFlags = bufferflag.BufferFlag(ba.Buffer[HEADERLENGTH-1])
}

func (ba *ByteArray) RewindToInitPos() {
	ba.ReaderIdx = 0
}

func (ba *ByteArray) RewindSteps(steps uint32) {
	if ba.ReaderIdx < steps {
		panic("too many step to rewind")
	}
	ba.ReaderIdx -= steps
}

func (ba *ByteArray) OverflowCheck(steps uint32) bool {
	if ba.ReaderIdx+steps >= ba.WriterIdx {
		return true
	}
	return false
}

func (ba *ByteArray) Reset() {
	ba.ReaderIdx = 0
	ba.WriterIdx = 0
}
