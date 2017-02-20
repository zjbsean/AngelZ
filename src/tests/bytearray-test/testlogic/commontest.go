package testlogic

import (
	"bytearray/bufferflag"
	"fmt"
	"unsafe"
)

func EnumTest() {
	var e bytearrayflag.BufferFlag
	e = bytearrayflag.HadEncrypt
	fmt.Println(bytearrayflag.GetInstance().TestFlag(e, bytearrayflag.HadEncrypt))

	var i int = 0x00331201
	pi := unsafe.Pointer(&i)
	pi = unsafe.Pointer(uintptr(pi) + 0)
	fmt.Println(*(*uint16)(pi))
}
