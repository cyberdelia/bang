package main

import (
	"flag"
	"fmt"
	"github.com/rcrowley/go-metrics"
	"net/http"
	"os"
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
	flag.StringVar(&url, "url", "", "URL to hit")
	flag.IntVar(&concurrency, "concurrency", 10, "Concurrency")
	flag.StringVar(&duration, "duration", "10s", "Duration")
	flag.StringVar(&method, "method", "GET", "HTTP method")
	flag.StringVar(&contentType, "content-type", "text/plain", "Content-Type")
	flag.StringVar(&body, "body", "", "Request body")
}

func worker(request *http.Request, timer *metrics.StandardTimer) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Fprintf(os.Stderr, "%s\n", e)
			os.Exit(1)
		}
	}()
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
	percentiles := timer.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
	fmt.Printf("Successful calls  \t\t %9d\n", timer.Count())
	fmt.Printf("Total time        \t\t %12.2fs\n", duration.Seconds())
	fmt.Printf("Fastest           \t\t %12.2fs\n", time.Duration(timer.Min()).Seconds())
	fmt.Printf("Slowest           \t\t %12.2fs\n", time.Duration(timer.Max()).Seconds())
	fmt.Printf("Mean              \t\t %12.2fs\n", time.Duration(timer.Mean()).Seconds())
	fmt.Printf("Standard deviation\t\t %12.2fs\n", time.Duration(timer.StdDev()).Seconds())
	fmt.Printf("Median            \t\t %12.2fs\n", time.Duration(percentiles[0]).Seconds())
	fmt.Printf("75th percentile   \t\t %12.2fs\n", time.Duration(percentiles[1]).Seconds())
	fmt.Printf("95th percentile   \t\t %12.2fs\n", time.Duration(percentiles[2]).Seconds())
	fmt.Printf("99th percentile   \t\t %12.2fs\n", time.Duration(percentiles[3]).Seconds())
	fmt.Printf("99.9th percentile \t\t %12.2fs\n", time.Duration(percentiles[4]).Seconds())
	fmt.Printf("Mean rate         \t\t %12.2f\n", timer.RateMean())
	fmt.Printf("1-min rate        \t\t %12.2f\n", timer.Rate1())
	fmt.Printf("5-min rate        \t\t %12.2f\n", timer.Rate5())
	fmt.Printf("15-min rate       \t\t %12.2f\n", timer.Rate15())
}

func main() {
	flag.Parse()

	if url == "" {
		fmt.Fprintf(os.Stderr, "please specify an url to run bang\n")
		os.Exit(1)
	}

	timer := metrics.NewTimer()

	d, err := time.ParseDuration(duration)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't parse duration: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Running %d workers for at least %s\n", concurrency, duration)
	fmt.Println("Starting to load the server")

	go func() {
		request, err := http.NewRequest(method, url, strings.NewReader(body))
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't build request: %s\n", err)
			os.Exit(1)
		}
		request.ContentLength = -1
		request.Header.Add("Content-Type", contentType)
		for i := 0; i < concurrency; i++ {
			go worker(request, timer)
		}
	}()

	select {
	case <-time.After(d):
		summary(d, timer)
	}
}
