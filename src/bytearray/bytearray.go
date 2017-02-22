package bytearray

import (
	"bytearray/bufferflag"
	"unsafe"
	"errors"
)

const (
	MAXIMUM_PACKAGE = 256 * 1024
	HEADERLENGTH    = 5
)

var ByteArrayOperExceptionDeal func(ba *ByteArray, err error)

type ByteArray struct {
	Buffer      []byte
	ReaderIdx   int
	WriterIdx   int
	CRC         uint32
	HeaderFlags bufferflag.BufferFlag
	Capacity int
}

func CreateByteArray() *ByteArray {
	ba := &ByteArray{make([]byte, MAXIMUM_PACKAGE, MAXIMUM_PACKAGE), 0, 0, 0, 0, MAXIMUM_PACKAGE}
	return ba
}

func (ba *ByteArray) Length() int {
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
	l := *(*int)(unsafe.Pointer(&ba.Buffer[0]))
	if l != ba.WriterIdx {
		panic("buffer length not equal WriterIdx")
	}
	ba.GetHeaderFlags()
}

func (ba *ByteArray) BypassHeader() {
	ba.ReaderIdx = HEADERLENGTH
	ba.GetHeaderFlags()
}

func (ba *ByteArray) AdvanceReader(steps int) {
	ba.ReaderIdx += steps
}

func (ba *ByteArray) AdvanceWriter(steps int) {
	ba.WriterIdx += steps
}

func (ba *ByteArray) GetHeaderFlags() {
	ba.HeaderFlags = bufferflag.BufferFlag(ba.Buffer[HEADERLENGTH-1])
}

func (ba *ByteArray) RewindToInitPos() {
	ba.ReaderIdx = 0
}

func (ba *ByteArray) RewindSteps(steps int) {
	if ba.ReaderIdx < steps {
		panic("too many step to rewind")
	}
	ba.ReaderIdx -= steps
}

func (ba *ByteArray) ReadOverflowCheck(steps int) bool {
	if ba.ReaderIdx+steps >= ba.WriterIdx {
		return true
	}
	return false
}

func (ba *ByteArray) WriteOverflowCheck(steps int) bool {
	if ba.WriterIdx+steps >= ba.Capacity {
		return true
	}
	return false
}

func (ba *ByteArray) Reset() {
	ba.ReaderIdx = 0
	ba.WriterIdx = 0
}

func (ba *ByteArray) ReadByte() (d byte){
	if ba.ReadOverflowCheck(1) {
		if ByteArrayOperExceptionDeal != nil{
			ByteArrayOperExceptionDeal(ba, errors.New("Read Byte Overflow !"))
		} else {
			panic("Read Byte Overflow !")
		}
		d = 0
	} else {
		d = ba.Buffer[ba.ReaderIdx]
		ba.ReaderIdx++
	}
	return
}

func (ba *ByteArray) WriteByte(d byte){

}

func (ba *ByteArray) ReadInt8() int8{

}

func (ba *ByteArray)WriteInt8(d int8){

}

func (ba *ByteArray) ReadInt16() int16{

} 

func (ba *ByteArray) WriteInt16(d int16){

}

func (ba *ByteArray) ReadUInt16() uint16{

}

func (ba *ByteArray) WriteUInt16(d uint16){

}

func (ba *ByteArray) ReadInt32() int32{

}

func (ba *ByteArray) WriteInt32(d int32){

}

func (ba *ByteArray) ReadUInt32() uint32{

}

func (ba *ByteArray) WriteUInt32(d uint32){

}

func (ba *ByteArray) ReadInt64() int64{

}

func (ba *ByteArray) WriteInt64(d int64){

}

func (ba *ByteArray) ReadUInt64() uint64{

}

func (ba *ByteArray) WriteUInt64(d uint64){

}

func (ba *ByteArray) ReadFloat32() float32{

}

func (ba *ByteArray) WriteFloat32(d float32){

}

func (ba *ByteArray) ReadFloat64() float64{

}

func (ba *ByteArray) WriteFloat64(d float64){

}