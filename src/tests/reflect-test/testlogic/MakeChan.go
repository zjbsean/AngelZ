package testlogic

import (
	"fmt"
	"reflect"
)

func MakeChanTest() {
	type SS struct {
		Name string
	}

	SSC := make(chan *SS, 1)
	//var i string = "abc"
	v := reflect.MakeChan(reflect.TypeOf(SSC), 1)
	cv := v.Interface().(chan *SS)
	fmt.Println(reflect.TypeOf(cv))

	ss := &SS{
		Name: "abc",
	}
	var i interface{} = ss
	fmt.Println(reflect.TypeOf(i))
}

type two [2]uintptr

func dummy(b byte, c int, d byte, e two, f byte, g float32, h byte) (i byte, j int, k byte, l two, m byte, n float32, o byte) {
	return b, c, d, e, f, g, h
}

func MakeFuncTest() {
	f := dummy
	//f(1, 2, 3, two{4, 5}, 6, 7, 8)

	fv := reflect.MakeFunc(reflect.TypeOf(f), func(in []reflect.Value) []reflect.Value { return in })

	reflect.ValueOf(&f).Elem().Set(fv)

	g := dummy
	g(1, 2, 3, two{4, 5}, 6, 7, 8)

	i, j, k, l, m, n, o := f(10, 20, 30, two{40, 50}, 60, 70, 80)
	if i != 10 || j != 20 || k != 30 || l != (two{40, 50}) || m != 60 || n != 70 || o != 80 {
		fmt.Printf("Call returned %d, %d, %d, %v, %d, %g, %d; want 10, 20, 30, [40, 50], 60, 70, 80", i, j, k, l, m, n, o)
	}

	fmt.Println("finish")
}
