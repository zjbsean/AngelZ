package main

import (
	"fmt"
	"log"
	"logger"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	logger.SetLevel(logger.INFO)
	logger.SetConsole(false)
	logger.SetHourRollFile("./log", "testHour.log")
	for i := 0; i < 1000; i++ {
		go putlog()
	}

	var i int
	fmt.Scanf("%d", &i)
}

func putlog() {
	for i := 1; ; i++ {
		logger.Info("the index asdfsafsafsafsafadsfasfasdfdsafsadfasdfasdfsadfasffasdfsaf = %d", i)
		time.Sleep(1 * time.Millisecond)
	}
}
