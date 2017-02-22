package testlogic

import (
	"fmt"
	"unsafe"
)

func SleceTest() {
	var s []int = make([]int, 5, 10)
	fmt.Printf("s[0]=%v, s[4]=%v, Length=%v, Cap=%v, Addr=%v\n", s[0], s[4], len(s), cap(s), unsafe.Pointer(&s))
	for i := 0; i < len(s); i++ {
		fmt.Println(unsafe.Pointer(&s[i]))
	}
	s = append(s, 6, 7, 8, 9, 10)
	fmt.Printf("s[5]=%v, Len=%v, Cap=%v, Addr=%v\n", s[5], len(s), cap(s), unsafe.Pointer(&s))
	for i := 0; i < len(s); i++ {
		fmt.Println(unsafe.Pointer(&s[i]))
	}
	s = append(s, 11)
	fmt.Printf("Len=%v, Cap=%v, Addr=%v\n", len(s), cap(s), unsafe.Pointer(&s))
	for i := 0; i < len(s); i++ {
		fmt.Println(unsafe.Pointer(&s[i]))
	}
	cs := s[1:4]
	fmt.Println("Print cs:")
	for i := 0; i < len(cs); i++ {
		fmt.Println(unsafe.Pointer(&cs[i]))
	}
	fmt.Println("--------------------")
	s = append(s, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21)
	cs = append(cs, 50)
	for i := 0; i < len(cs); i++ {
		fmt.Println(unsafe.Pointer(&cs[i]))
	}

	fmt.Println(s[4])
	fmt.Println("Print s:")
	for i := 0; i < len(s); i++ {
		fmt.Println(unsafe.Pointer(&s[i]))
	}
}
