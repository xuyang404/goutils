package WaitGo

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var lock sync.Mutex

type WaitGo struct {
	token        chan struct{}
	cap          int
	PanicHandler func(interface{})
}

func NewWaitGo(cap int) *WaitGo {
	gp := &WaitGo{
		token: make(chan struct{}, cap),
		cap:   cap,
	}

	for i := 0; i < gp.cap; i++ {
		gp.token <- struct{}{}
	}

	return gp
}

func (gp *WaitGo) Add(task func()) {
	<-gp.token
	go func() {
		defer func() {
			gp.token <- struct{}{}
			if r := recover(); r != nil {
				pc, file, line, ok := runtime.Caller(3)
				if ok {
					if gp.PanicHandler != nil {
						gp.PanicHandler(r)
					} else {
						//log.Printf("task paniced: %s", r)
						lock.Lock()
						funcName := runtime.FuncForPC(pc).Name()
						fmt.Println("[goutils.waitGo]")
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
		}()

		task()
	}()
}

func (gp *WaitGo) Wait() {
	for i := 0; i < gp.cap; i++ {
		<-gp.token
	}

	close(gp.token)
}
