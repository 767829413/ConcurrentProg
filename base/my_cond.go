package base

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type Queue struct {
	capacity int
	Data     []int
	c        *sync.Cond
}

func NewQueue(num int) *Queue {
	return &Queue{
		capacity: num,
		Data:     make([]int, 0, num),
		c:        sync.NewCond(&sync.Mutex{}),
	}
}

func (q *Queue) push(val int) {
	defer q.c.L.Unlock()
	q.c.L.Lock()
	for len(q.Data) == q.capacity {
		q.c.Wait()
	}
	q.Data = append(q.Data, val)
	q.c.Broadcast()
}

func (q *Queue) pop() int {
	defer q.c.L.Unlock()
	q.c.L.Lock()
	for len(q.Data) == 0 {
		q.c.Wait()
	}
	val := q.Data[0]
	q.Data = q.Data[1:]
	q.c.Broadcast()
	return val
}

// Athlete race
func ExampleOne() {
	c := sync.NewCond(&sync.Mutex{})
	defer c.L.Unlock()
	var ready int

	for i := 0; i < 10; i++ {
		go func(i int) {
			defer c.L.Unlock()
			time.Sleep(time.Duration(rand.Int63n(5)) * time.Second)
			c.L.Lock()
			ready++
			log.Printf("Athlete: %d Ready", i)
			//Wake up waiter
			c.Signal()
		}(i)
	}
	c.L.Lock()
	for ready != 10 {
		c.Wait()
		log.Printf("The referee is awakened")
	}
	log.Printf("Game start")
}
