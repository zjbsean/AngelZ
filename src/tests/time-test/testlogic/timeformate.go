package testlogic

import (
	"fmt"
	"math"
	"time"
)

//时间格式所设定的时间必须是：2006-01-02 15:04:05
const _DateFormat = "2006-01-02"
const _DateTimeHFormat = "2006-01-02 15"
const _DateTimeHMFormat = "2006-01-02 15:04"
const _DateTimeHMSFormat = "2006-01-02 15:04:05"
const _DateTimeHMSMFormat = "2006-01-02 15:04:05.000"

func TimeFormation() {
	t := time.Now()
	y, m, d := t.Date()
	today := time.Now().Format(_DateFormat)
	datetime := time.Now().Format(_DateTimeHFormat) //后面的参数是固定的 否则将无法正常输出
	datetimeM := time.Now().Format(_DateTimeHMSMFormat)

	fmt.Println("time is : ", t)
	fmt.Println("y m d is : ", y, m, d)
	fmt.Println("now is :", today)
	fmt.Println("now is :", datetime)
	fmt.Println("now is :", datetimeM)
	nt, _ := time.Parse(_DateTimeHFormat, datetime)
	fmt.Println(nt)

	fmt.Println(math.MaxInt64 - 2)
	fmt.Println(int64(math.Max(float64(math.MaxInt64-2), float64(math.MaxInt64-1))))

	a := 1
	b := 2
	x := calc(&a, &b)
	defer calc(&a, &x)
	a = 0
	//x = calc(&a, &b)
	//defer calc(&a, &x)
	b = 1
	fmt.Println("The End")
	f := 1.1
	xx := int64(f)
	fmt.Println(xx)
}

func calc(a, b *int) int {
	ret := *a + *b
	fmt.Println(ret)
	return ret
}
