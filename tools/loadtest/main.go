package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	url := flag.String("url", "http://127.0.0.1:8080/", "target URL")
	requests := flag.Int("n", 10000, "number of requests")
	concurrency := flag.Int("c", 100, "number of concurrent workers")
	flag.Parse()

	var success int64
	latencies := make([]time.Duration, *requests)
	client := &http.Client{}

	sem := make(chan struct{}, *concurrency)
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < *requests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sem <- struct{}{}
			t0 := time.Now()
			resp, err := client.Get(*url)
			lat := time.Since(t0)
			latencies[i] = lat
			if err == nil && resp.StatusCode == http.StatusOK {
				atomic.AddInt64(&success, 1)
			}
			if resp != nil {
				resp.Body.Close()
			}
			<-sem
		}(i)
	}
	wg.Wait()
	duration := time.Since(start)

	var total time.Duration
	for _, l := range latencies {
		total += l
	}
	avg := total / time.Duration(len(latencies))

	rps := float64(*requests) / duration.Seconds()
	successRate := float64(success) / float64(*requests) * 100

	fmt.Printf("Total requests: %d\n", *requests)
	fmt.Printf("Total time: %v\n", duration)
	fmt.Printf("RPS: %.2f\n", rps)
	fmt.Printf("Avg latency: %v\n", avg)
	fmt.Printf("Success rate: %.2f%%\n", successRate)
}