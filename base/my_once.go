package base

import (
	"sync"
	"sync/atomic"
)

type MyOnce struct {
	m    sync.Mutex
	done uint32
}

func (o *MyOnce) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 {
		return nil
	}
	return o.doSlow(f)
}

func (o *MyOnce) doSlow(f func() error) error {
	defer o.m.Unlock()
	o.m.Lock()
	var err error
	if o.done == 0 {
		err = f()
		if err != nil {
			atomic.StoreUint32(&o.done, 1)
		}
	}
	return err
}
