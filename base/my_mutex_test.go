package base

import (
	"sync"
	"testing"
)

func TestGeneralCountTest(t *testing.T) {
	GeneralCount()
}

func TestPrintConst(t *testing.T) {
	t.Log(mutexLocked)
	t.Log(mutexWoken)
	t.Log(1 << mutexWaiterShift)
}

func TestGetGoId(t *testing.T) {
	t.Log(GetGoId())
}

func TestUnUseLock(t *testing.T) {
	UnUseLock()
}

func TestCopyLock(t *testing.T) {
	CopyLock()
}

func TestReentrantLock(t *testing.T) {
	l := &sync.Mutex{}
	ReentrantLock(l)
}
func TestDeadLock(t *testing.T) {
	DeadLock()
}
