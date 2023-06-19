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
	channelLimit      int = 50
	channelReadPeople     = make(chan People, channelLimit)
	currentLimit      int32
)

func init() {
	go func() {
		for v := range channelReadPeople {
		GOLOOP:
			if int(currentLimit) < channelLimit {
				goto BEGIN
			} else {
				goto GOLOOP
			}
		BEGIN:
			atomic.AddInt32(&currentLimit, 1)
			go process(v)
		}
	}()
}

func process(p People) {
	fmt.Printf("People is: %+v: Total goroutine: %v Size limit: %v\n", People{Name: p.Name, Age: p.Age}, runtime.NumGoroutine(), currentLimit)
	atomic.AddInt32(&currentLimit, -1)
	p.wg.Done()
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
