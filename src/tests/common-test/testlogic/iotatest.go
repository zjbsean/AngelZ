package testlogic

import "fmt"

const (
	B = 1 << (10 * iota)
	KB
	MB
	GB
	TB
	PB
)

func IotaTest() {
	fmt.Println(B, KB, MB, GB, TB, PB)
}
