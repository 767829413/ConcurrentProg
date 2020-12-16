package base

import (
	"fmt"
	"sync"
	"time"
)

type UserError struct {
	sync.Mutex
	Num int
}

func UnUseLock() {
	var mu sync.Mutex
	defer mu.Unlock()
	fmt.Println("hello world!")
}

func CopyLock() {
	var u UserError
	u.Lock()
	defer u.Unlock()
	u.Num++
	copyLock(u)
}

func copyLock(u UserError) {
	u.Lock()
	defer u.Unlock()
	fmt.Println("in example")
}

func ReentrantLock(l sync.Locker) {
	defer l.Unlock()
	l.Lock()
	fmt.Println("out lock")
	reentrantLock(l)
}

func reentrantLock(l sync.Locker) {
	defer l.Unlock()
	l.Lock()
	fmt.Println("inner lock")
}

func DeadLock() {
	var lock1 sync.Mutex
	var lock2 sync.Mutex
	var wg sync.WaitGroup
	var resource int
	wg.Add(2)
	go func(resource int) {
		defer lock1.Unlock()
		defer wg.Done()
		lock1.Lock()
		fmt.Println("我要锁资源1")
		resource++
		time.Sleep(3 * time.Second)
		lock2.Lock()
		lock2.Unlock()

	}(resource)

	go func(resource int) {
		defer lock2.Unlock()
		defer wg.Done()
		lock2.Lock()
		fmt.Println("我要锁资源2")
		resource++
		time.Sleep(3 * time.Second)
		lock1.Lock()
		lock1.Unlock()
	}(resource)

	wg.Wait()

}
