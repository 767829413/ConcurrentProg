package base

import (
	"fmt"
	"hash/fnv"
	"sync"
)

func ExampleStructSetKey() {
	type mapKey struct {
		key int
	}

	tk := mapKey{key: 66}
	m := make(map[mapKey]string)
	m[tk] = "run"
	fmt.Printf("m[tk] = %s\n", m[tk])
	tk.key = 77
	fmt.Printf("reload m[tk] = %s\n", m[tk])
}

//easy Concurrent map
type Concurrent struct {
	m map[int]int
	l sync.RWMutex
}

func NewConcurrent(cap int) *Concurrent {
	return &Concurrent{m: make(map[int]int, cap),}
}

func (c *Concurrent) Set(k, v int) {
	defer c.l.Unlock()
	c.l.Lock()
	c.m[k] = v
}

func (c *Concurrent) Get(k int) (int, bool) {
	defer c.l.RUnlock()
	c.l.RLock()
	v, ok := c.m[k]
	return v, ok
}

func (c *Concurrent) Delete(k int) {
	defer c.l.Unlock()
	c.l.Lock()
	delete(c.m, k)
}

func (c *Concurrent) Len() int {
	defer c.l.RUnlock()
	c.l.RLock()
	return len(c.m)
}

func (c *Concurrent) Each(f func(k, v int) error) error {
	defer c.l.RUnlock()
	c.l.RLock()
	for k, v := range c.m {
		err := f(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

//Fragmentation Map
var ShardCount = 32

type ConcurrentMap []*ConcurrentMapShared

type ConcurrentMapShared struct {
	items map[string]interface{}
	l     sync.RWMutex
}

func NewConcurrentMap() ConcurrentMap {
	m := make(ConcurrentMap, ShardCount)
	for i := 0; i < ShardCount; i++ {
		m[i] = &ConcurrentMapShared{items: make(map[string]interface{}),}
	}
	return m
}

func (m ConcurrentMap) getShard(k string) (*ConcurrentMapShared, error) {
	h := fnv.New32()
	_, err := h.Write([]byte(k))
	if err != nil {
		return nil, err
	}
	return m[h.Sum32()%uint32(ShardCount)], nil
}

func (m ConcurrentMap) Set(k string, v interface{}) error {
	shard, err := m.getShard(k)
	if err != nil {
		return err
	}
	shard.l.Lock()
	shard.items[k] = v
	shard.l.Unlock()
	return nil
}

func (m ConcurrentMap) Get(k string) (interface{}, bool) {
	shard, err := m.getShard(k)
	if err != nil {
		return nil, false
	}
	shard.l.RLock()
	v, ok := shard.items[k]
	shard.l.RUnlock()
	return v, ok
}

func (m ConcurrentMap) Delete(k string) error {
	shard, err := m.getShard(k)
	if err != nil {
		return err
	}
	shard.l.Lock()
	delete(shard.items, k)
	shard.l.Unlock()
	return nil
}
