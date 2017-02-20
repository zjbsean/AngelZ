package testlogic

import (
	"fmt"
	"net"
	"os"
)

func ResloveIPAddrTest(url string) {
	addr, err := net.ResolveIPAddr("ip", url)
	if err != nil {
		fmt.Printf("ResolveIPAddr Fail : %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("IP : %v\n", addr.IP)
}

func ParseIPTest(ipStr string) {
	ip := net.ParseIP(ipStr)
	fmt.Printf("IP : %v\n", ip)
}

func AddrTest() {

}
