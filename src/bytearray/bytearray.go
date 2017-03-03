package bytearray

import (
	"bytearray/bufferflag"
	"bytes"
	"encoding/binary"
	"errors"
	"math"
)

const (
	PACKAGE_DEFAULT_SIZE = 1024
	HEADER_LENGTH        = 5
)

var readFromByteArrayOverflowMsg = "Read Byte Overflow !"

var ByteArrayOperExceptionDeal func(ba *ByteArray, err error)

type ByteArray struct {
	Buffer      *bytes.Buffer
	ReaderIdx   int
	WriterIdx   int
	CRC         uint32
	HeaderFlags bufferflag.BufferFlag
	IsBigEndian bool
}

func isBigEndian(initFlag bufferflag.BufferFlag) bool {
	if (initFlag & bufferflag.BigEndian) != 0 {
		return true
	}
	return false
}

func CreateByteArrayWithFlag(initFlag bufferflag.BufferFlag) *ByteArray {
	buff := bytes.NewBuffer(make([]byte, HEADER_LENGTH, PACKAGE_DEFAULT_SIZE))
	headflag := isBigEndian(initFlag)
	ba := &ByteArray{buff, HEADER_LENGTH, HEADER_LENGTH, 0, initFlag, headflag}
	return ba
}

func CreateByteArrayWithBytes(d []byte) *ByteArray {
	//ba := &ByteArray{d}
}

func (ba *ByteArray) Length() int {
	return ba.WriterIdx
}

func (ba *ByteArray) RefreshHeader() {
	buf := make([]byte, 0, 4)
	if isBigEndian(ba.HeaderFlags) {
		binary.BigEndian.PutUint32(buf, uint32(ba.WriterIdx))
	} else {
		binary.LittleEndian.PutUint32(buf, uint32(ba.WriterIdx))
	}
	bb := ba.Buffer.Bytes()
	bb[0] = buf[0]
	bb[1] = buf[1]
	bb[2] = buf[2]
	bb[3] = buf[3]
	bb[4] = byte(ba.HeaderFlags)
}

func (ba *ByteArray) InitHeaderWithBuffer() {
	tb := ba.Buffer.Bytes()
	ba.HeaderFlags = bufferflag.BufferFlag(tb[HEADER_LENGTH-1])
	ba.IsBigEndian = isBigEndian(ba.HeaderFlags)
	var lenth int
	if ba.IsBigEndian {
		lenth = int(binary.BigEndian.Uint32(tb))
	} else {
		lenth = int(binary.LittleEndian.Uint32(tb))
	}

	if lenth != ba.Buffer.Len() {
		panic("buffer length not equal WriterIdx")
	}
	ba.ReaderIdx = HEADER_LENGTH
	ba.WriterIdx = lenth
}

func getHeaderInfoFromBytes(header []byte) (len int, headerFlag byte) {
	headerFlag = header[HEADER_LENGTH-1]
	if isBigEndian(bufferflag.BufferFlag(headerFlag)) {
		len = int(binary.BigEndian.Uint32(header))
	} else {
		len = int(binary.LittleEndian.Uint32(header))
	}
	return
}

func (ba *ByteArray) BypassHeader() {
	ba.ReaderIdx = HEADER_LENGTH
}

func (ba *ByteArray) AdvanceReader(steps int) {
	ba.ReaderIdx += steps
}

func (ba *ByteArray) AdvanceWriter(steps int) {
	ba.WriterIdx += steps
}

func (ba *ByteArray) RewindSteps(steps int) {
	if (ba.ReaderIdx - HEADER_LENGTH) < steps {
		panic("too many step to rewind")
	}
	ba.ReaderIdx -= steps
}

func (ba *ByteArray) ReadOverflowCheck(steps int) bool {
	if (ba.ReaderIdx + steps) >= ba.WriterIdx {
		return true
	}
	return false
}

func (ba *ByteArray) Reset() {
	ba.ReaderIdx = HEADER_LENGTH
	ba.WriterIdx = HEADER_LENGTH
}

func (ba *ByteArray) ReadByte() (d byte) {
	if ba.ReadOverflowCheck(1) {
		if ByteArrayOperExceptionDeal != nil {
			ByteArrayOperExceptionDeal(ba, errors.New(readFromByteArrayOverflowMsg))
		} else {
			panic(readFromByteArrayOverflowMsg)
		}
		d = 0
	} else {
		td := ba.Buffer.Bytes()
		d = td[ba.ReaderIdx]
		ba.ReaderIdx++
	}
	return
}

func (ba *ByteArray) WriteByte(d byte) {
	binary.Write(ba.Buffer, binary.LittleEndian, d)
	ba.WriterIdx++
}

func (ba *ByteArray) ReadInt8() (d int8) {
	if ba.ReadOverflowCheck(1) {
		if ByteArrayOperExceptionDeal != nil {
			ByteArrayOperExceptionDeal(ba, errors.New(readFromByteArrayOverflowMsg))
		} else {
			panic(readFromByteArrayOverflowMsg)
		}
		d = 0
	} else {
		td := ba.Buffer.Bytes()
		d = int8(td[ba.ReaderIdx])
		ba.ReaderIdx++
	}
	return
}

func (ba *ByteArray) WriteInt8(d int8) {
	binary.Write(ba.Buffer, binary.LittleEndian, d)
	ba.WriterIdx++
}

func (ba *ByteArray) ReadInt16() (d int16) {
	if ba.ReadOverflowCheck(2) {
		if ByteArrayOperExceptionDeal != nil {
			ByteArrayOperExceptionDeal(ba, errors.New(readFromByteArrayOverflowMsg))
		} else {
			panic(readFromByteArrayOverflowMsg)
		}
		d = 0
	} else {
		td := ba.Buffer.Bytes()
		var ud uint16
		if ba.IsBigEndian == true {
			binary.BigEndian.PutUint16(td[ba.ReaderIdx:ba.ReaderIdx+2], ud)
		} else {
			binary.LittleEndian.PutUint16(td[ba.ReaderIdx:ba.ReaderIdx+2], ud)
		}
		d = int16(ud)
		ba.ReaderIdx += 2
	}
	return
}

func (ba *ByteArray) WriteInt16(d int16) {
	if ba.IsBigEndian {
		binary.Write(ba.Buffer, binary.BigEndian, d)
	} else {
		binary.Write(ba.Buffer, binary.LittleEndian, d)
	}
	ba.WriterIdx += 2
}

func (ba *ByteArray) ReadUInt16() (d uint16) {
	if ba.ReadOverflowCheck(2) {
		if ByteArrayOperExceptionDeal != nil {
			ByteArrayOperExceptionDeal(ba, errors.New(readFromByteArrayOverflowMsg))
		} else {
			panic(readFromByteArrayOverflowMsg)
		}
		d = 0
	} else {
		td := ba.Buffer.Bytes()
		if ba.IsBigEndian == true {
			binary.BigEndian.PutUint16(td[ba.ReaderIdx:ba.ReaderIdx+2], d)
		} else {
			binary.LittleEndian.PutUint16(td[ba.ReaderIdx:ba.ReaderIdx+2], d)
		}
		ba.ReaderIdx += 2
	}
	return
}

func (ba *ByteArray) WriteUInt16(d uint16) {
	if ba.IsBigEndian {
		binary.Write(ba.Buffer, binary.BigEndian, d)
	} else {
		binary.Write(ba.Buffer, binary.LittleEndian, d)
	}
	ba.WriterIdx += 2
}

func (ba *ByteArray) ReadInt32() (d int32) {
	if ba.ReadOverflowCheck(4) {
		if ByteArrayOperExceptionDeal != nil {
			ByteArrayOperExceptionDeal(ba, errors.New(readFromByteArrayOverflowMsg))
		} else {
			panic(readFromByteArrayOverflowMsg)
		}
		d = 0
	} else {
		td := ba.Buffer.Bytes()
		var ud uint32
		if ba.IsBigEndian == true {
			binary.BigEndian.PutUint32(td[ba.ReaderIdx:ba.ReaderIdx+4], ud)
		} else {
			binary.LittleEndian.PutUint32(td[ba.ReaderIdx:ba.ReaderIdx+4], ud)
		}
		d = int32(ud)
		ba.ReaderIdx += 4
	}
	return
}

func (ba *ByteArray) WriteInt32(d int32) {
	if ba.IsBigEndian {
		binary.Write(ba.Buffer, binary.BigEndian, d)
	} else {
		binary.Write(ba.Buffer, binary.LittleEndian, d)
	}
	ba.WriterIdx += 4
}

func (ba *ByteArray) ReadUInt32() (d uint32) {
	if ba.ReadOverflowCheck(4) {
		if ByteArrayOperExceptionDeal != nil {
			ByteArrayOperExceptionDeal(ba, errors.New(readFromByteArrayOverflowMsg))
		} else {
			panic(readFromByteArrayOverflowMsg)
		}
		d = 0
	} else {
		td := ba.Buffer.Bytes()
		if ba.IsBigEndian == true {
			binary.BigEndian.PutUint32(td[ba.ReaderIdx:ba.ReaderIdx+4], d)
		} else {
			binary.LittleEndian.PutUint32(td[ba.ReaderIdx:ba.ReaderIdx+4], d)
		}
		ba.ReaderIdx += 4
	}
	return
}

func (ba *ByteArray) WriteUInt32(d uint32) {
	if ba.IsBigEndian {
		binary.Write(ba.Buffer, binary.BigEndian, d)
	} else {
		binary.Write(ba.Buffer, binary.LittleEndian, d)
	}
	ba.WriterIdx += 4
}

func (ba *ByteArray) ReadInt64() (d int64) {

	if ba.ReadOverflowCheck(8) {
		if ByteArrayOperExceptionDeal != nil {
			ByteArrayOperExceptionDeal(ba, errors.New(readFromByteArrayOverflowMsg))
		} else {
			panic(readFromByteArrayOverflowMsg)
		}
		d = 0
	} else {
		td := ba.Buffer.Bytes()
		var ud uint64
		if ba.IsBigEndian == true {
			binary.BigEndian.PutUint64(td[ba.ReaderIdx:ba.ReaderIdx+8], ud)
		} else {
			binary.LittleEndian.PutUint64(td[ba.ReaderIdx:ba.ReaderIdx+8], ud)
		}
		d = int64(ud)
		ba.ReaderIdx += 8
	}
	return
}

func (ba *ByteArray) WriteInt64(d int64) {
	if ba.IsBigEndian {
		binary.Write(ba.Buffer, binary.BigEndian, d)
	} else {
		binary.Write(ba.Buffer, binary.LittleEndian, d)
	}
	ba.WriterIdx += 8
}

func (ba *ByteArray) ReadUInt64() (d uint64) {
	if ba.ReadOverflowCheck(8) {
		if ByteArrayOperExceptionDeal != nil {
			ByteArrayOperExceptionDeal(ba, errors.New(readFromByteArrayOverflowMsg))
		} else {
			panic(readFromByteArrayOverflowMsg)
		}
		d = 0
	} else {
		td := ba.Buffer.Bytes()
		if ba.IsBigEndian == true {
			binary.BigEndian.PutUint64(td[ba.ReaderIdx:ba.ReaderIdx+8], d)
		} else {
			binary.LittleEndian.PutUint64(td[ba.ReaderIdx:ba.ReaderIdx+8], d)
		}
		ba.ReaderIdx += 8
	}
	return
}

func (ba *ByteArray) WriteUInt64(d uint64) {
	if ba.IsBigEndian {
		binary.Write(ba.Buffer, binary.BigEndian, d)
	} else {
		binary.Write(ba.Buffer, binary.LittleEndian, d)
	}
	ba.WriterIdx += 8
}

func (ba *ByteArray) ReadFloat32() (d float32) {
	if ba.ReadOverflowCheck(4) {
		if ByteArrayOperExceptionDeal != nil {
			ByteArrayOperExceptionDeal(ba, errors.New(readFromByteArrayOverflowMsg))
		} else {
			panic(readFromByteArrayOverflowMsg)
		}
		d = 0.0
	} else {
		td := ba.Buffer.Bytes()
		if ba.IsBigEndian == true {
			bits := binary.BigEndian.Uint32(td[ba.ReaderIdx : ba.ReaderIdx+4])
			d = math.Float32frombits(bits)
		} else {
			bits := binary.LittleEndian.Uint32(td[ba.ReaderIdx : ba.ReaderIdx+4])
			d = math.Float32frombits(bits)
		}
		ba.ReaderIdx += 4
	}
	return
}

func (ba *ByteArray) WriteFloat32(d float32) {
	bits := math.Float32bits(d)
	if ba.IsBigEndian {
		binary.Write(ba.Buffer, binary.BigEndian, bits)
	} else {
		binary.Write(ba.Buffer, binary.LittleEndian, bits)
	}
	ba.WriterIdx += 4
}

func (ba *ByteArray) ReadFloat64() (d float64) {
	if ba.ReadOverflowCheck(8) {
		if ByteArrayOperExceptionDeal != nil {
			ByteArrayOperExceptionDeal(ba, errors.New(readFromByteArrayOverflowMsg))
		} else {
			panic(readFromByteArrayOverflowMsg)
		}
		d = 0.0
	} else {
		td := ba.Buffer.Bytes()
		if ba.IsBigEndian == true {
			bits := binary.BigEndian.Uint64(td[ba.ReaderIdx : ba.ReaderIdx+8])
			d = math.Float64frombits(bits)
		} else {
			bits := binary.LittleEndian.Uint64(td[ba.ReaderIdx : ba.ReaderIdx+8])
			d = math.Float64frombits(bits)
		}
		ba.ReaderIdx += 8
	}
	return
}

func (ba *ByteArray) WriteFloat64(d float64) {
	bits := math.Float64bits(d)
	if ba.IsBigEndian {
		binary.Write(ba.Buffer, binary.BigEndian, bits)
	} else {
		binary.Write(ba.Buffer, binary.LittleEndian, bits)
	}
	ba.WriterIdx += 8
}

func (ba *ByteArray) ReadBool() (d bool) {
	bb := ba.ReadByte()
	if bb == 1 {
		d = true
	} else {
		d = false
	}
	return
}

func (ba *ByteArray) WriteBool(d bool) {
	var bb byte
	if d {
		bb = 1
	} else {
		bb = 0
	}
	ba.WriteByte(bb)
}

func (ba *ByteArray) ReadString() (d string) {
	strLen := ba.ReadUInt16()
	dd := ba.Buffer.Bytes()
	if ba.ReadOverflowCheck(int(strLen)) {
		if ByteArrayOperExceptionDeal != nil {
			ByteArrayOperExceptionDeal(ba, errors.New(readFromByteArrayOverflowMsg))
		} else {
			panic(readFromByteArrayOverflowMsg)
		}
		d = ""
	} else {
		d = string(dd[ba.ReaderIdx : ba.ReaderIdx+int(strLen)])
	}
	ba.ReaderIdx += int(strLen)
	return
}

func (ba *ByteArray) WriteString(d string) {
	strBytes := []byte(d)
	strLen := len(strBytes)
	if strLen > 0xFFFF {
		if ByteArrayOperExceptionDeal != nil {
			ByteArrayOperExceptionDeal(ba, errors.New(readFromByteArrayOverflowMsg))
		} else {
			panic(readFromByteArrayOverflowMsg)
		}
		strLen = 0xFFFF
	}
	ba.WriteUInt16(uint16(strLen))
	ba.Buffer.Write(strBytes[:strLen])
	ba.WriterIdx += strLen
}

func (ba *ByteArray) ReadByteArray() (d *ByteArray) {
	if ba.ReadOverflowCheck(HEADER_LENGTH) {
		if ByteArrayOperExceptionDeal != nil {
			ByteArrayOperExceptionDeal(ba, errors.New(readFromByteArrayOverflowMsg))
		} else {
			panic(readFromByteArrayOverflowMsg)
		}
		d = nil
	} else {
		dd := ba.Buffer.Bytes()
		length, headFlags := getHeaderInfoFromBytes(dd[ba.ReaderIdx : ba.ReaderIdx+HEADER_LENGTH])
		if ba.ReadOverflowCheck(length) {
			if ByteArrayOperExceptionDeal != nil {
				ByteArrayOperExceptionDeal(ba, errors.New(readFromByteArrayOverflowMsg))
			} else {
				panic(readFromByteArrayOverflowMsg)
			}
			d = nil
		} else {

		}
	}
}

func (ba *ByteArray) WriteByteArray(d *ByteArray) {

}
