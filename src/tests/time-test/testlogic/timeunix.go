package testlogic

import (
	"fmt"
	"time"
)

func TimeUnixTest() {
	ct := time.Now()
	fmt.Println("Cur Time : ", ct)

	local := ct.Location()

	//zn, off := ct.Zone()
	//fmt.Println("Cur Time Zone Name and Offset : ", zn, off)

	dr := ct.Round(24 * time.Hour)
	fmt.Println("Cur Time Round 24*time.Hour : ", dr)

	dt := ct.Truncate(24 * time.Hour)
	fmt.Println("Cur Time Truncate 24*time.Hour : ", dt)

	ctu := ct.Unix()
	fmt.Println("Cur Time Unix : ", ctu)

	ctut := time.Unix(ctu, 0)
	fmt.Println("Cur Time Unix Time : ", ctut)

	ctut = time.Unix(ctu+3600, 0)
	fmt.Println("Cur Time Unix + 3600 Time : ", ctut)

	fmt.Println("---------------------------------------------")

	fct, _ := time.ParseInLocation("2006-01-02 15", ct.Format("2006-01-02 15"), local)

	fmt.Println("Cur Format Time : ", fct)

	fctu := fct.Unix()
	fmt.Println("Cur Format Time Unix : ", fctu)

	fctut := time.Unix(fctu, 0)
	fmt.Println("Cur Time Unix Time : ", fctut)

	fctut = time.Unix(fctu+3600, 0)
	fmt.Println("Cur Time Unix + 3600 Time : ", fctut)
}
