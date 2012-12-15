package main

import (
	"flag"
	"fmt"
	"github.com/rcrowley/go-metrics"
	"net/http"
	"strings"
	"time"
)

var (
	url         string
	concurrency int
	duration    string
	method      string
	contentType string
	body        string
)

func init() {
	flag.StringVar(&url, "url", "http://localhost:80", "URL to hit")
	flag.IntVar(&concurrency, "concurrency", 10, "Concurrency")
	flag.StringVar(&duration, "duration", "10s", "Duration")
	flag.StringVar(&method, "method", "GET", "HTTP method")
	flag.StringVar(&contentType, "content-type", "text/plain", "Content-Type")
	flag.StringVar(&body, "body", "", "Request body")
}

func worker(request *http.Request, timer *metrics.StandardTimer) {
	for {
		timer.Time(func() {
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				panic(err)
			}
			defer response.Body.Close()
		})
	}
}

func summary(duration time.Duration, timer *metrics.StandardTimer) {
	fmt.Printf("Successful calls  \t\t %9d\n", timer.Count())
	fmt.Printf("Total time        \t\t %12.2fs\n", duration.Seconds())
	fmt.Printf("Fastest           \t\t %12.2fs\n", time.Duration(timer.Min()).Seconds())
	fmt.Printf("Slowest           \t\t %12.2fs\n", time.Duration(timer.Max()).Seconds())
	fmt.Printf("Average           \t\t %12.2fs\n", time.Duration(timer.Mean()).Seconds())
	fmt.Printf("Standard deviation\t\t %12.2fs\n", time.Duration(timer.StdDev()).Seconds())
	fmt.Printf("Mean rate         \t\t %12.2f\n", timer.RateMean())
	fmt.Printf("1-min rate        \t\t %12.2f\n", timer.Rate1())
	fmt.Printf("5-min rate        \t\t %12.2f\n", timer.Rate5())
	fmt.Printf("15-min rate       \t\t %12.2f\n", timer.Rate15())
}

func main() {
	flag.Parse()

	timer := metrics.NewTimer()

	d, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}

	timeout := time.NewTimer(d)
	defer timeout.Stop()

	fmt.Printf("Running %d workers for at least %s\n", concurrency, duration)
	fmt.Println("Starting to load the server")

	go func() {
		request, err := http.NewRequest(method, url, strings.NewReader(body))
		if err != nil {
			panic(err)
		}
		request.ContentLength = -1
		request.Header.Add("Content-Type", contentType)
		for i := 0; i < concurrency; i++ {
			go worker(request, timer)
		}
	}()

	select {
	case <-timeout.C:
		summary(d, timer)
	}
}
