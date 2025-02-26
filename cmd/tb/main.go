package main

import (
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tb",
	Short: "tb benchmark",
	RunE: func(cmd *cobra.Command, args []string) error {
		return tb()
	},
}

const (
	ArgURL         = "url"
	ArgRequests    = "requests"
	ArgConcurrency = "concurrency"

	defaultRequests    = 30
	defaultConcurrency = 3
)

type inputArgs struct {
	url         *string
	requests    *int
	concurrency *int
}

var args = inputArgs{}

func tb() error {
	url := *args.url
	requests := *args.requests
	concurrency := *args.concurrency

	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	results := make(chan int, concurrency)

	println("Processing, wait...")
	println("")

	start := time.Now()
	for range concurrency {
		go func() {
			defer wg.Done()

			for range requests {
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					println("Failed to create request:", err.Error())
				}

				res, err := http.DefaultClient.Do(req)
				if err != nil {
					println("Failed to perform request:", err.Error())
				}

				results <- res.StatusCode
			}
		}()
	}

	statusCodes := make(map[int]int)
	go func() {
		for r := range results {
			statusCodes[r]++
		}
	}()
	wg.Wait()

	time.Sleep(time.Second)

	close(results)

	elapsed := time.Since(start)

	requestsPerformed := 0
	orderedCodeResults := make([]string, 0, len(statusCodes))
	okResult := "200 requests: 0"
	for k, v := range statusCodes {
		if k == 200 {
			okResult = fmt.Sprintf("%d requests: %d", k, v)
			requestsPerformed += v
			continue
		}

		orderedCodeResults = append(orderedCodeResults, fmt.Sprintf("%d requests: %d", k, v))
		requestsPerformed += v
	}
	slices.Sort(orderedCodeResults)

	println("Url:", url)
	println("Concurrency:", concurrency)
	println("Requests:", requests)
	println("Total Requests:", concurrency*requests)
	println("Total Requests complete:", requestsPerformed)
	println("Total time:", elapsed.String())
	println(okResult)
	for _, codeResult := range orderedCodeResults {
		println(codeResult)
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		println(err)
	}
}

func init() {
	args.url = rootCmd.Flags().String(ArgURL, "", "URL to benchmark")
	args.requests = rootCmd.Flags().Int(ArgRequests, defaultRequests, "Number of requests")
	args.concurrency = rootCmd.Flags().Int(ArgConcurrency, defaultConcurrency, "Number of concurrent requests")

	rootCmd.MarkFlagRequired(ArgURL)
}
