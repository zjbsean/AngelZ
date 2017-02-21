package testlogic

import (
	"bytearray/bufferflag"
	"fmt"
	"unsafe"
)

func EnumTest() {
	var e bufferflag.BufferFlag
	e = bufferflag.HadEncrypt
	fmt.Println(bufferflag.GetInstance().TestFlag(e, bufferflag.HadEncrypt))

	var i int = 0x00331201
	pi := unsafe.Pointer(&i)
	pi = unsafe.Pointer(uintptr(pi) + 0)
	fmt.Println(*(*uint16)(pi))
}
