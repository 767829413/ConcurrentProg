package base

import (
	"math/rand"
	"testing"
	"time"
)

func TestExampleOne(t *testing.T) {
	ExampleOne()
}

func TestMyCondQueue(t *testing.T) {
	q := NewQueue(5)
	go func() {
		for {
			q.push(rand.Intn(20))
		}

	}()
	for {
		t.Log(q.pop())
		time.Sleep(1 * time.Second)
	}
}
