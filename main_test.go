package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	tests := map[string]struct {
		totalPeople int
		expect      string
	}{
		"should test: 100": {
			totalPeople: 100,
			expect:      "People is: {Name:should test: 100 Age:99",
		},
		"should test: 500": {
			totalPeople: 500,
			expect:      "People is: {Name:should test: 500 Age:499",
		},
		"should test: 1000": {
			totalPeople: 1000,
			expect:      "People is: {Name:should test: 1000 Age:999",
		},
		"should test: 10000": {
			totalPeople: 10000,
			expect:      "People is: {Name:should test: 10000 Age:9999",
		},
		"should test: 100000": {
			totalPeople: 100000,
			expect:      "People is: {Name:should test: 100000 Age:99999",
		},
		"should test: 1000000": {
			totalPeople: 1000000,
			expect:      "People is: {Name:should test: 1000000 Age:999999",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tempFile, err := os.CreateTemp("", "")
			if err != nil {
				log.Fatal(err)
			}

			old := os.Stdout

			os.Stdout = tempFile

			req, err := http.NewRequest("POST", fmt.Sprintf("/?value=%v&name=%s", test.totalPeople, name), nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(handlerFunc)

			startTime := time.Now()
			handler.ServeHTTP(rr, req)
			endTime := time.Since(startTime).Seconds()

			os.Stdout = old

			_, err = tempFile.Seek(0, 0)
			if err != nil {
				t.Fatal(err)
			}

			b, err := io.ReadAll(tempFile)
			if err != nil {
				t.Fatal(err)
			}

			var words []string

			words = append(words, strings.Split(string(b), "\n")...)

			var exist bool

			for _, s := range words {
				if strings.Contains(s, test.expect) {
					exist = true
				}
			}

			if !exist {
				t.Errorf("expect contains: %s, but don't find", test.expect)
				return
			}

			fmt.Printf("%s finished: %v\n", name, endTime)
		})
	}
}
