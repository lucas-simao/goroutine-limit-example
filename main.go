package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"sync"
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
)

func init() {
	var limit int

	var mx = sync.Mutex{}

	go func() {
		for v := range channelReadPeople {
		GOLOOP:
			if limit < channelLimit {
				goto BEGIN
			} else {
				goto GOLOOP
			}
		BEGIN:
			mx.Lock()
			limit += 1
			mx.Unlock()
			go process(v, &limit, &mx)
		}
	}()
}

func process(p People, limit *int, mx *sync.Mutex) {
	fmt.Printf("People is: %+v: Total goroutine: %v Size limit: %v\n", p, runtime.NumGoroutine(), *limit)
	mx.Lock()
	defer mx.Unlock()
	*limit -= 1
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
