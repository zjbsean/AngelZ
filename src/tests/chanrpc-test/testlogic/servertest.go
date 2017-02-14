package testlogic

import (
	"chanrpc"
	"fmt"
	"time"
)

func regFunc0(agrs []interface{}) {
	fmt.Printf("regFunc0 %v \n", agrs)
}

func regFunc1(args []interface{}) interface{} {
	fmt.Printf("regFunc1 %v \n", args)
	return "regFunc1"
}

func regFuncN(args []interface{}) []interface{} {
	fmt.Printf("regFunc2 %v \n", args)

	ret := make([]interface{}, 3)
	ret[0] = "regFunc2"
	ret[1] = "p1"
	ret[2] = "p2"
	return ret
}

func cbFunc0(err error) {
	if err != nil {
		fmt.Printf("cbFunc0 - err=%v\n", err)
	} else {
		fmt.Printf("cbFunc0 - succ\n")
	}
}

func cbFunc1(ret interface{}, err error) {
	if err != nil {
		fmt.Printf("cbFunc1 - err=%v\n", err)
	} else {
		fmt.Printf("cbFunc1 - succ: ret=%v\n", ret)
	}
}

func cbFuncN(ret []interface{}, err error) {
	if err != nil {
		fmt.Printf("cbFuncN - err=%v\n", err)
	} else {
		fmt.Printf("cbFuncN - succ: ret=%v\n", ret)
	}
}

func serverExec(server *chanrpc.Server, runSign chan int) {
	var deal int
breakloop:
	for {
		fmt.Scanf("%d\n", &deal)
		fmt.Println(deal)
		switch deal {
		case -1:
			fmt.Println("lock runSign")
			runSign <- 1
			fmt.Println("unlock runSign")
			break breakloop
		case 0:
			server.Go("f0", 0, 1, 2)
		case 1:
			server.Go("f1", 1, 2, 3)
		case 2:
			server.Go("fn", 2, 3, 4)
		case 10:
			server.Call0("f0", 10, 11, 12)
		case 11:
			server.Call1("f1", 11, 12, 13)
		case 12:
			server.CallN("fn", 12, 13, 14)
		default:
			fmt.Println("unknow deal type : ", deal)
		}
	}
}

func clientExec(server *chanrpc.Server) {
	client := server.Open(1000)

	doCbCount := 0

	go func() {
		for {
			if client.Idle() == false {
				fmt.Println("Client Not Idle !!! ", client.PendingAsynCllCount())
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			client.AsynCall("fn", 110, 111, 112, cbFuncN)
		}
		for {
			select {
			case ri := <-client.ChanAsynRet:
				client.Cb(ri)
				doCbCount = doCbCount + 1
			}

		}
	}()

	/*
		var deal int

			breakloop:
				for {
					fmt.Scanf("%d\n", &deal)
					fmt.Println(deal)
					switch deal {
					case -1:
						fmt.Println("lock runSign")
						stopSign <- 1
						fmt.Println("unlock runSign")
						break breakloop
					case 0:
						err := client.Call0("f0", 0, 1, 2)
						if err != nil {
							fmt.Println("call0 error : ", err)
						} else {
							fmt.Println("call0 succ")
						}
					case 1:
						ret, err := client.Call1("f1", 1, 2, 3)
						if err != nil {
							fmt.Println("call1 error : ", err)
						} else {
							fmt.Println("call1 succ : ", ret)
						}
					case 2:
						ret, err := client.CallN("fn", 1, 2, 3)
						if err != nil {
							fmt.Println("callN error : ", err)
						} else {
							fmt.Println("callN succ : ", ret)
						}
					case 10:
						client.AsynCall("f0", 10, 11, 12, cbFunc0)
					case 11:
						client.AsynCall("f1", 11, 12, 13, cbFunc1)
					case 12:
						client.AsynCall("fn", 12, 13, 14, cbFuncN)
					case 100:
						for i := 0; i < 10000; i++ {
							client.AsynCall("fn", 110, 111, 112, cbFuncN)
						}
					default:
						fmt.Println("unknow deal type : ", deal)
					}
				}
	*/
	fmt.Println("do cb count : ", doCbCount)
	var d int
	fmt.Scan(&d)
}

func TestRun() {
	server := chanrpc.NewServer(10)
	server.Register("f0", regFunc0)
	server.Register("f1", regFunc1)
	server.Register("fn", regFuncN)

	runSign := make(chan int, 1)
	defer func() {
		fmt.Println("end 1")
		close(runSign)
		server.Close()
		fmt.Println("end 2")
	}()
	go func() {
	stopRunLoop:
		for {
			select {
			case ci := <-server.ChanCall:
				server.Exec(ci)
			case <-runSign:
				break stopRunLoop
			}
		}
	}()

	clientExec(server)
}
