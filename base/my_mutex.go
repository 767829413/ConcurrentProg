package base

import (
	"fmt"
	myGoid "github.com/petermattis/goid"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

func GeneralCount() {
	var count Counter
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100000; j++ {
				count.Incr()
			}
		}()
	}
	wg.Wait()
	fmt.Println(count.Count())
}

type Counter struct {
	CounterType int
	Name        string

	mu    sync.Mutex
	count uint64
}

func (c *Counter) Incr() {
	defer c.mu.Unlock()
	c.mu.Lock()
	c.count++
}

func (c *Counter) Count() uint64 {
	defer c.mu.Unlock()
	c.mu.Lock()
	return c.count
}

func GetGoId() (int, error) {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	str := strings.Split(string(buf[:n]), " ")
	return strconv.Atoi(str[1])
}

type RecursiveMutexGoId struct {
	sync.Mutex
	owner     int64 //当前 goroutine 的 id
	recursion int32 //可重入次数
}

func (r *RecursiveMutexGoId) Lock() {
	goId := myGoid.Get()
	//如果当前持有锁的goroutine就是这次调用的goroutine,说明是重入
	if atomic.LoadInt64(&r.owner) == goId {
		r.recursion++
		return
	}
	r.Mutex.Lock()
	//如果当前持有锁的goroutine就是这次调用的goroutine,说明是重入
	atomic.StoreInt64(&r.owner, goId)
	r.recursion = 1
}

func (r *RecursiveMutexGoId) Unlock() {
	goId := myGoid.Get()
	//非持有锁的goroutine尝试释放锁，错误的使用
	if atomic.LoadInt64(&r.owner) != goId {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", r.owner, goId))
	}
	//调用次数减一
	r.recursion--
	//未完全释放,直接返回
	if r.recursion != 0 {
		return
	}
	atomic.StoreInt64(&r.owner, 0)
	r.Mutex.Unlock()
}

type RecursiveMutexToken struct {
	sync.Mutex
	token     int64 //自定义 goroutine 的标识
	recursion int32 //可重入次数
}

func (r *RecursiveMutexToken) Lock(token int64) {
	//如果传入的token和持有锁的token一致，说明是递归调用
	if atomic.LoadInt64(&r.token) == token {
		r.recursion++
		return
	}
	//传入的token不一致，说明不是递归调用
	r.Mutex.Lock()
	//抢到锁记录这个token
	atomic.StoreInt64(&r.token, token)
	r.recursion = 1
}

func (r *RecursiveMutexToken) Unlock(token int64) {
	if atomic.LoadInt64(&r.token) != token {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", r.token, token))
	}
	r.recursion--
	if r.recursion != 0 {
		return
	}
	atomic.StoreInt64(&r.token, 0)
	r.Mutex.Unlock()
}
