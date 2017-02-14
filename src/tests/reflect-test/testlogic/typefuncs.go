package testlogic

import (
	"fmt"
	"reflect"
)

func GetType(a interface{}) {
	t := reflect.TypeOf(a)
	fmt.Println(t.Kind())
}
