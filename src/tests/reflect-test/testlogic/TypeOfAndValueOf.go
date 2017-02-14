package testlogic

import (
	"fmt"
	"reflect"
)

type Speaker interface {
	Speak() string
	Set(string)
}

type Teacher struct {
	Name string
	Addr string
}

func (this *Teacher) Speak() string {
	return this.Name
}

func (this *Teacher) Set(name string) {
	this.Name = name
}

func (this *Teacher) SetAddr(addr string) {
	this.Addr = addr
}

func TestRef(s Speaker) (t reflect.Type) {
	reflectVal := reflect.ValueOf(s)
	t = reflectVal.Elem().Type()
	fmt.Printf("reflect.ValueOf(%v).Elem().Type()=%v, Value=%v\n", s, t, reflectVal)
	return
}

func TestRef2(s Speaker) (t reflect.Type) {
	t = reflect.TypeOf(s)
	fmt.Printf("reflect.TypeOf(%v)=%v\n", s, t)
	return
}

func TestRef3(s *Teacher) (t reflect.Type) {
	t = reflect.TypeOf(s).Elem()
	fmt.Printf("reflect.TypeOf(%v).Elem()=%v\n", s, t)
	return
}

func TestRef4() {
	te := &Teacher{
		Name: "newyear",
		Addr: "SH",
	}
	fmt.Printf("source=%v\n", te)
	fmt.Printf("\n")

	var tea Speaker = te
	t1 := TestRef(tea)

	if m, ok := reflect.New(t1).Interface().(Speaker); ok {
		fmt.Printf("reflect.New(%v).Interface().(Speaker)=%v\n", t1, m)
		fmt.Printf("se.Speak()=%v\n", m.Speak())
		m.Set("2014")
		fmt.Printf("reflect.New(%v).Interface().(Speaker)=%v\n", t1, m)
		fmt.Printf("se.Speak()=%v\n", m.Speak())
	}
}

func TestRef5() {
	te := &Teacher{
		Name: "newyear",
		Addr: "SH",
	}
	fmt.Printf("source=%v\n", te)
	fmt.Printf("\n")

	var tea Speaker = te
	t2 := TestRef2(tea)

	if m, ok := reflect.New(t2.Elem()).Interface().(Speaker); ok {
		fmt.Printf("reflect.New(%v).Interface().(Speaker)=%v\n", t2, m)
		fmt.Printf("se.Speak()=%v\n", m.Speak())
		m.Set("2014")
		fmt.Printf("reflect.New(%v).Interface().(Speaker)=%v\n", t2, m)
		fmt.Printf("se.Speak()=%v\n", m.Speak())
	}
}

func TestRef6() {
	/*te := &Teacher{
		Name: "newyear",
		Addr: "SH",
	}
	var tea Speaker = te
	*/
	te := "dafsafd"
	tp := &te
	t := reflect.TypeOf(tp)
	fmt.Printf("Type = %v, Type Elem = %v\n", t, t.Elem())
	//fmt.Printf("TypeOf te Elem()=%v\n", t.Elem())
	//fmt.Printf("TypeOf te Elem()=%v\n", t.Elem().Kind())
	v := reflect.ValueOf(tp)
	fmt.Printf("ValueOf = %v\n", v.Elem())
	fmt.Printf("abc \n")
}
