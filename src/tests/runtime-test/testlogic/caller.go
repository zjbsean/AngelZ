package testlogic

import (
	"fmt"
	"runtime"
)

func Caller() {
	for skip := 0; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if ok == false {
			break
		}

		fmt.Printf("PC:%d, File:%s, Line:%d, OK:%T\n", pc, file, line, ok)
	}

	fmt.Println("------------------------------------------------")

	for skip := 0; ; skip++ {
		pc, _, _, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		p := runtime.FuncForPC(pc)
		file, line := p.FileLine(0)

		fmt.Printf("skip = %v, pc = %v\n", skip, pc)
		fmt.Printf("  file = %v, line = %d\n", file, line)
		fmt.Printf("  entry = %v\n", p.Entry())
		fmt.Printf("  name = %v\n", p.Name())
	}

}
