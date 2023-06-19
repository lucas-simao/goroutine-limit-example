package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type People struct {
	Name string
	Age  int
	wg   *sync.WaitGroup
}

var (
	channelLimit      int = 100
	channelReadPeople     = make(chan People)
	currentLimit      int32
)

func init() {
	go func() {
		for {
			if int(currentLimit) < channelLimit {
				var v = <-channelReadPeople
				go process(v)
			}
		}
	}()
}

func process(p People) {
	atomic.AddInt32(&currentLimit, 1)
	defer func() {
		atomic.AddInt32(&currentLimit, -1)
		p.wg.Done()
	}()
	fmt.Printf("People is: %+v: Total goroutine: %v - currentLimit: %v\n", p, runtime.NumGoroutine(), currentLimit)
}

func main() {
	http.HandleFunc("/", handlerFunc)

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Panic(err)
	}
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}

	start := time.Now()

	value, err := strconv.Atoi(r.FormValue("value"))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "should pass params value with value to process")
		return
	}

	name := r.FormValue("name")
	if len(name) == 0 {
		w.WriteHeader(500)
		fmt.Fprint(w, "should pass params name")
		return
	}

	var wg = sync.WaitGroup{}

	for i := 0; i < value; i++ {
		wg.Add(1)
		channelReadPeople <- People{Name: name, Age: i, wg: &wg}
	}

	wg.Wait()

	w.WriteHeader(201)

	fmt.Fprintf(w, "Time used to process %v is %vs", value, time.Since(start).Seconds())
}
