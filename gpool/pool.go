package gpool

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var lock sync.Mutex

type Pool struct {
	queue        chan struct{}
	worker       chan func()
	PanicHandler func(interface{})
	Recover      bool
}

func NewPool(cap int) *Pool {
	return &Pool{
		queue:        make(chan struct{}, cap),
		worker:       make(chan func()),
		PanicHandler: nil,
		Recover:      true,
	}
}

func (p *Pool) Add(task func()) {
	select {
	case p.worker <- task:
	case p.queue <- struct{}{}:
		go p.work(task)
	}
}

func (p *Pool) work(task func()) {
	defer func() {
		if p.Recover {
			if r := recover(); r != nil {
				pc, file, line, ok := runtime.Caller(3)
				if ok {
					if p.PanicHandler != nil {
						p.PanicHandler(r)
					} else {
						funcName := runtime.FuncForPC(pc).Name()
						lock.Lock()
						fmt.Println("[goutils.gpool]")
						fmt.Println("-------------------------------------------------------------------")
						fmt.Println("time", time.Now().Format("2006-01-02 15:04:05"))
						fmt.Println("func", funcName)
						fmt.Println("file", file)
						fmt.Println("line", line)
						fmt.Println(r)
						fmt.Println("-------------------------------------------------------------------")
						lock.Unlock()
					}
				}
			}
		}
		<-p.queue
	}()

	for {
		task()
		task = <-p.worker
	}
}
