package testlogic

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func NoObjOut() {
	fmt.Println("Flags : ", log.Flags())
	log.SetPrefix("Debug ")
	fmt.Println("Prefix : ", log.Prefix())

	var byteBUff = bytes.Buffer{}
	log.SetOutput(&byteBUff)
	log.SetOutput(os.Stdout)
	log.Println("xxxxxx")
	fmt.Println("-----------------------")
	fmt.Println(&byteBUff)

	defer func() {
		fmt.Println("run defer !")
		if err := recover(); err != nil {
			fmt.Println("defer and recover")
			fmt.Println(err)
		}
	}()

	var i, j = 10, 1
	x := i / j
	fmt.Println(x)
	fmt.Println("Func End!")
	var d = 1
	var s = "dasfas"
	var f = 1.1
	fmt.Printf("%v - %v - %v", d, s, f)
}

func FileOutput() {
	fileName := "debug.log"
	logFile, err := os.Create(fileName)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
	}
	debugLog := log.New(logFile, "[Debug]", log.Lshortfile)
	debugLog.Println("A debug message here")
	debugLog.SetPrefix("[Info]")
	debugLog.Println("A Info Message here ")
	debugLog.SetFlags(debugLog.Flags() | log.Ldate | log.Lmicroseconds)
	logStr := fmt.Sprintf("A different prefix, context=%s", "adsfsdaf")
	debugLog.Println(logStr)
}
