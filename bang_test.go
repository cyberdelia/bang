package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))
	defer ts.Close()
	request, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case timer := <-run(request, time.Millisecond, 1):
		if timer.Count() <= 0 {
			t.Fatal("no requests processed")
		}
	}
}
