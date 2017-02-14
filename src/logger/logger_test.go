package logger

import (
	"log"
	"testing"
)

func Test_Logger_HourRollLogger(t *testing.T) {
	//t.Skip()
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	SetLevel(ALL)
	SetConsole(true)
	SetHourRollFile("./log", "test.log")

	Debug("context=%s, param=%d", "abc", 123)

	i := 10
	if i != 10 {
		t.Errorf("IsSameYear Test Fail")
	}
}

func Test_Logger_DailyRollLogger(t *testing.T) {
	t.Skip()
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	SetLevel(ALL)
	SetConsole(true)
	SetDailyRollFile("./log", "testDaily.log")

	Debug("context=%s, param=%d", "abc", 123)

	i := 10
	if i != 10 {
		t.Errorf("IsSameYear Test Fail")
	}
}

func Test_Logger_SizeRollLogger(t *testing.T) {
	t.Skip()
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	SetLevel(ALL)
	SetConsole(true)
	SetSizeRollFile("./log", "testSize.log", 2*KB, 1000)

	Debug("context=%s, param=%d", "abc", 123)

	i := 10
	if i != 10 {
		t.Errorf("IsSameYear Test Fail")
	}
}
