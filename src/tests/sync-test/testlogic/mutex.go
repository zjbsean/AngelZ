package testlogic

import (
	"fmt"
	"sync"
)

type mtest struct {
	mu   *sync.Mutex
	data int
}

func MutexTest() {
	mt := new(mtest)
	mt.mu = new(sync.Mutex)
	mt.data = 1
	func() {
		mt.mu.Lock()
		defer func() {
			mt.mu.Unlock()
			fmt.Println("first unlock")
		}()
		fmt.Println("first lock")
	}()

	mt.mu.Lock()
	defer func() {
		mt.mu.Unlock()
		fmt.Println("second unlock")
	}()
	fmt.Println("second lock")

}
