package testlogic

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

func Bytesbufftest() {
	buf := new(bytes.Buffer)
	var pi float64 = math.Pi
	buf.Grow(100)
	err := binary.Write(buf, binary.LittleEndian, pi)
	if err != nil {
		fmt.Println("binary.Write failed : ", err)
	}
	var px int64 = 100
	binary.Write(buf, binary.LittleEndian, px)
	fmt.Printf("Len=%d, Cap=%d, % x\n", buf.Len(), cap(buf.Bytes()), buf.Bytes())
	for index := 0; index < 20; index++ {
		binary.Write(buf, binary.LittleEndian, px)
	}
	fmt.Printf("Len=%d, Cap=%d, % x\n", buf.Len(), cap(buf.Bytes()), buf.Bytes())
	cb := make([]byte, 4)
	buf.Read(cb)
	fmt.Printf("Len=%d, Cap=%d, % x\n", buf.Len(), cap(buf.Bytes()), buf.Bytes())
}

func Bytesbufftest_1() {
	buff := bytes.NewBuffer(make([]byte, 5, 10))
	buf := bytes.NewBuffer(make([]byte, 0, 4))
	binary.Write(buf, binary.LittleEndian, int32(10))
	bb := buff.Bytes()
	tb := buf.Bytes()
	bb[0] = tb[0]
	bb[1] = tb[1]
	bb[2] = tb[2]
	bb[3] = tb[3]
	fmt.Printf("% x", buff.Bytes())
}
