package main

import (
	"flag"
	"fmt"
	"github.com/rcrowley/go-metrics"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	url         string
	concurrency int
	count       int
	duration    string
	method      string
	contentType string
	body        string
	auth        string
)

func init() {
	flag.StringVar(&url, "url", "", "URL to hit")
	flag.IntVar(&concurrency, "concurrency", 10, "Concurrency")
	flag.StringVar(&duration, "duration", "10s", "Duration")
	flag.IntVar(&count, "requests", 0, "Number of requests")
	flag.StringVar(&method, "method", "GET", "HTTP method")
	flag.StringVar(&contentType, "content-type", "text/plain", "Content-Type")
	flag.StringVar(&body, "body", "", "Request body")
	flag.StringVar(&auth, "auth", "", "Credentials as user:password")
}

type worker struct {
	timer       *metrics.StandardTimer
	done        chan *metrics.StandardTimer
	request     *http.Request
	concurrency int
}

func newWorker(request *http.Request, concurrency int) *worker {
	return &worker{
		request:     request,
		concurrency: concurrency,
		timer:       metrics.NewTimer(),
		done:        make(chan *metrics.StandardTimer),
	}
}

func (r *worker) call() {
	r.timer.Time(func() {
		response, err := http.DefaultClient.Do(r.request)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		defer response.Body.Close()
	})
}

func (r *worker) run(count int, duration time.Duration) {
	if count > 0 {
		go r.counter(count)
	} else {
		go r.duration(duration)
	}
}

func (r *worker) counter(count int) {
	fmt.Printf("Running %d times per %d workers\n", r.concurrency, count)
	wg := sync.WaitGroup{}
	for i := 0; i < r.concurrency; i++ {
		for j := 0; j < count; j++ {
			wg.Add(1)
			go func() {
				r.call()
				wg.Done()
			}()
		}
	}
	wg.Wait()
	r.done <- r.timer
}

func (r *worker) duration(duration time.Duration) {
	fmt.Printf("Running %d workers for at least %s\n", r.concurrency, duration)
	for i := 0; i < r.concurrency; i++ {
		go func() {
			for {
				r.call()
			}
		}()
	}
	select {
	case <-time.After(duration):
		r.done <- r.timer
	}
}

func runner(request *http.Request, concurrency int, duration time.Duration, count int) <-chan *metrics.StandardTimer {
	runner := newWorker(request, concurrency)
	runner.run(count, duration)
	return runner.done
}

func summary(timer *metrics.StandardTimer) {
	percentiles := timer.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
	fmt.Printf("Successful calls  \t\t %9d\n", timer.Count())
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

	d, err := time.ParseDuration(duration)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't parse duration: %s\n", err)
		os.Exit(1)
	}

	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't build request: %s\n", err)
		os.Exit(1)
	}
	request.ContentLength = -1
	request.Header.Add("Content-Type", contentType)

	if auth != "" {
		credentials := strings.Split(auth, ":")
		request.SetBasicAuth(credentials[0], credentials[1])
	}

	fmt.Println("Starting to load the server")
	select {
	case timer := <-runner(request, concurrency, d, count):
		summary(timer)
	}
}
