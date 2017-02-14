package util

import (
	"testing"
	"time"
)

func Test_Time_IsSameYear(t *testing.T) {

	t1 := time.Now()
	t2 := time.Date(2016, 10, 23, 22, 23, 45, 23432, t1.Location())
	if IsSameYear(&t1, &t2) == false {
		t.Errorf("IsSameYear Test Fail : T1=%v, T2=%v", t1, t2)
	}
	t3 := time.Date(2015, 10, 23, 22, 23, 45, 23432, t1.Location())
	if IsSameYear(&t1, &t3) == true {
		t.Errorf("IsSameYear Test Fail : T1=%v, T3=%v", t1, t3)
	}
}

func Benchmark_Time_IsSameYear(b *testing.B) {
	t1 := time.Now()
	t2 := time.Date(2016, 10, 23, 22, 23, 45, 23432, t1.Location())
	for i := 0; i < b.N; i++ {
		IsSameYear(&t1, &t2)
	}
}
