package base

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestExampleStructSetKey(t *testing.T) {
	ExampleStructSetKey()
}

func TestMapReturn(t *testing.T) {
	var m = make(map[string]int)
	m["a"] = 0
	t.Logf("a=%d; b=%d\n", m["a"], m["b"])

	av, aexisted := m["a"]
	bv, bexisted := m["b"]
	t.Logf("a=%d, existed: %t; b=%d, existed: %t\n", av, aexisted, bv, bexisted)
}

func TestEasyConcurrentMap(t *testing.T) {
	n := 6
	m := NewConcurrent(n)
	go func() {
		for {
			k := rand.Intn(6)
			v := rand.Intn(6)
			m.Set(k, v)
			t.Log("SET KEY VALUE ", k, v)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			k := rand.Intn(6)
			v, ok := m.Get(k)
			t.Log("GET KEY VALUE IS_EXIST ", k, v, ok)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			k := rand.Intn(6)
			m.Delete(rand.Intn(6))
			t.Log("DELETE KEY ", k)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			t.Log("LEN ", m.Len())
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			f := func(k, v int) error {
				_, err := fmt.Println("PRINT K V ", k, v)
				return err
			}
			err := m.Each(f)
			if err != nil {
				t.Fatal(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	select {}
}

func TestConcurrentShardMap(t *testing.T) {
	m := NewConcurrentMap()
	go func() {
		for {
			k := strconv.Itoa(rand.Intn(6))
			v := rand.Intn(6)
			err := m.Set(k, v)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("SET KEY VALUE ", k, v)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			k := strconv.Itoa(rand.Intn(6))
			v, ok := m.Get(k)
			t.Log("GET KEY VALUE IS_EXIST ", k, v, ok)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			k := strconv.Itoa(rand.Intn(6))
			err := m.Delete(k)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("DELETE KEY ", k)
			time.Sleep(1 * time.Second)
		}
	}()
	select {}
}
