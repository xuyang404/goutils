package WaitGo_test

import (
	"fmt"
	"github.com/xuyang404/goutils/WaitGo"
	"sync/atomic"
	"testing"
	"time"
)

var runTimes = 1000000
var sum int64

func TestWaitGo(t *testing.T) {
	waitgo := WaitGo.NewWaitGo(10)

	for i := 0; i < 10; i++ {
		a := i
		waitgo.Add(func() {
			time.Sleep(2 * time.Second)
			fmt.Println("a", a)
		})
	}

	waitgo.Wait()

	fmt.Println("finish")
}

func TestPoolRecover (t *testing.T) {
	wg := WaitGo.NewWaitGo(20)
	//wg.PanicHandler = func(r interface{}) {
	//	fmt.Println(r)
	//}

	for i := 0; i < 10; i++ {
		wg.Add(func() {
			b := 0
			fmt.Println(1/b)
		})
	}

	wg.Wait()
}

func BenchmarkWaitGo(b *testing.B) {
	wg := WaitGo.NewWaitGo(20)

	for i := 0; i < runTimes; i++ {
		wg.Add(func() {
			for a := 0; a < 100; a++ {
				atomic.AddInt64(&sum, int64(a))
			}
		})
	}

	wg.Wait()
}