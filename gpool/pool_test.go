package gpool_test

import (
	"fmt"
	"github.com/xuyang404/goutils/gpool"
	"sync"
	"sync/atomic"
	"testing"
)

var wg sync.WaitGroup
var sum int64
var runTimes = 1000000


func TestPoolRecover (t *testing.T) {
	p := gpool.NewPool(20)
	p.Recover = false
	p.PanicHandler = func(r interface{}) {
		fmt.Println(r)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		p.Add(func() {
			defer wg.Done()
			b := 0
			fmt.Println(1/b)
		})
	}

	wg.Wait()
}

func TestPool(t *testing.T) {
	p := gpool.NewPool(20)

	p.PanicHandler = func(r interface{}) {
		fmt.Println(r)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		p.Add(func() {
			defer wg.Done()
			b := 0
			fmt.Println(1/b)
		})
	}

	wg.Wait()
}

func BenchmarkPool(b *testing.B) {
	p := gpool.NewPool(20)

	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		p.Add(func() {
			defer wg.Done()
			for a := 0; a < 100; a++ {
				atomic.AddInt64(&sum, int64(a))
			}
		})
	}

	wg.Wait()
}
