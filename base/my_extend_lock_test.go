package base

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMyMutexTryLock(t *testing.T) {
	var myMutex MyMutex
	go func() {
		defer myMutex.Unlock(11)
		myMutex.Lock(11)
		fmt.Println("我锁了")
		time.Sleep(5 * time.Second)
	}()
	for {
		time.Sleep(1 * time.Second)
		if myMutex.TryLock(12) {
			fmt.Println("我拿锁了")
			myMutex.Unlock(12)
			break
		}
	}
	fmt.Println("end")
}

func TestMyMutexWaitCount(t *testing.T) {
	var wg sync.WaitGroup
	var myMutex MyMutex
	for i := 11; i < 17; i++ {
		wg.Add(1)
		go func(i int) {
			defer myMutex.Unlock(int64(i))
			defer wg.Done()
			myMutex.Lock(int64(i))
			if i == 11 {
				time.Sleep(2 * time.Second)
			}
		}(i)
	}
	time.Sleep(1 * time.Second)
	fmt.Println(6, myMutex.WaitCount())
	wg.Wait()
}
