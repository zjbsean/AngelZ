package logger

import (
	"fmt"
	"log"
	"time"
)

//LOG级别
type LEVEL int32

const (
	ALL LEVEL = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

var (
	consoleAppender       = true //是否Console输出
	logLevel        LEVEL = 1    //LOG级别

)

//SetConsole : 设置是否Console输出
func SetConsole(isConsole bool) {
	consoleAppender = isConsole
	if consoleAppender == true {
		log.SetFlags(log.Lmicroseconds)
	}
}

//SetLevel : set LOG level
func SetLevel(level LEVEL) {
	logLevel = level
}

//Debug : Debug级别日志
func Debug(format string, params ...interface{}) {
	logto("debug", DEBUG, format, params...)
}

//Info : Info级别日志
func Info(format string, params ...interface{}) {
	logto("info", INFO, format, params...)
}

//Warn : Warn级别日志
func Warn(format string, params ...interface{}) {
	logto("warn", WARN, format, params...)
}

//Error : Error级别日志
func Error(format string, params ...interface{}) {
	logto("Error", ERROR, format, params...)
}

//Fatal : Fatal级别日志
func Fatal(format string, params ...interface{}) {
	logto("Fatal", FATAL, format, params...)
}

func logto(levelFlag string, level LEVEL, format string, params ...interface{}) {
	defer catchError()

	if logLevel > level {
		return
	}

	logStr := fmt.Sprintf("["+levelFlag+"] "+format, params...)

	if fileObj != nil {
		fileObj.mu.RLock()
		defer fileObj.mu.RUnlock()
		curTime := time.Now()
		fileObj.chackAndRename(&curTime)

		fileObj.lg.Println(logStr)
	}

	console(logStr)
}

//console: console输出
func console(context string) {
	if consoleAppender {
		log.Println(context)
	}
}

//catchError：捕获输出处理
func catchError() {
	if err := recover(); err != nil {
		log.Println("err", err)
	}
}
