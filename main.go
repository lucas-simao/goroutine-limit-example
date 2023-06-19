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
}

const channelLimit int = 50

var (
	channelReadPeople = make(chan struct{}, channelLimit)
)

func process(p People) {
	fmt.Printf("People is: %+v: Total goroutine: %v\n", p, runtime.NumGoroutine())
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

	start := time.Now()

	var wg = sync.WaitGroup{}

	for i := 0; i < value; i++ {
		wg.Add(1)

		channelReadPeople <- struct{}{}

		go func(i int) {
			defer wg.Done()

			process(People{Name: name, Age: i})
			<-channelReadPeople
		}(i)
	}

	wg.Wait()

	w.WriteHeader(201)

	fmt.Fprintf(w, "Time used to process %v is %vs", value, time.Since(start).Seconds())
}
