package base

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken
	mutexStarving
	mutexWaiterShift      = iota
	starvationThresholdNs = 1e6
)

type MyMutex struct {
	sync.Mutex
	token     int64 //自定义 goroutine 的标识
	recursion int32 //可重入次数
}

func (m *MyMutex) Lock(token int64) {
	//如果传入的token和持有锁的token一致，说明是递归调用
	if atomic.LoadInt64(&m.token) == token {
		m.recursion++
		return
	}
	//传入的token不一致，说明不是递归调用
	m.Mutex.Lock()
	//抢到锁记录这个token
	atomic.StoreInt64(&m.token, token)
	m.recursion = 1
}

func (m *MyMutex) Unlock(token int64) {
	if atomic.LoadInt64(&m.token) != token {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.token, token))
	}
	m.recursion--
	if m.recursion != 0 {
		return
	}
	atomic.StoreInt64(&m.token, 0)
	m.Mutex.Unlock()
}

func (m *MyMutex) TryLock(token int64) bool {
	if atomic.LoadInt64(&m.token) == token {
		m.recursion++
		return true
	}
	//成功抢到锁
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked) {
		//抢到锁记录这个token
		atomic.StoreInt64(&m.token, token)
		m.recursion = 1
		return true
	}
	// 处于唤醒,加锁,饥饿则不参加此次竞争
	oldV := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if oldV&(mutexLocked|mutexStarving|mutexWoken) != 0 {
		return false
	}
	//尝试竞争拿锁
	newV := oldV | mutexLocked
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), oldV, newV) {
		//抢到锁记录这个token
		atomic.StoreInt64(&m.token, token)
		m.recursion = 1
		return true
	}
	return false
}

func (m *MyMutex) WaitCount() int {
	//通过state字段来获取等待锁的 goroutine
	num := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	num = num >> mutexWaiterShift
	num = num + (num & mutexLocked)
	return int(num)
}
