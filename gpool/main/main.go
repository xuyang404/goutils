package main

import (
	"github.com/xuyang404/goutils/gpool"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Pool *gpool.Pool

func test(w http.ResponseWriter, r *http.Request)  {
	if r.RequestURI == "/favicon.ico" {
		return
	}
	vals := r.URL.Query()
	ids := vals.Get("ids")
	idsArr := strings.Split(ids, ",")
	wg := &sync.WaitGroup{}
	for range idsArr {
		wg.Add(1)
		Pool.Add(func() {
			defer wg.Done()
			time.Sleep(1*time.Second)
		})
	}
	wg.Wait()
	num := strconv.Itoa(runtime.NumGoroutine())
	w.Write([]byte(num))
}

func main()  {
	Pool = gpool.NewPool(5)

	http.HandleFunc("/", test)
	http.ListenAndServe(":7890", nil)
}
