package scanner

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
)

type Result struct {
	Word       string
	URL        string
	StatusCode int
	Err        error
}

type Job struct {
	Word string
	URL  string
}

type Scanner struct {
	Client  *http.Client
	Workers int
}

func New(client *http.Client, workers int) *Scanner {
	return &Scanner{
		Client:  client,
		Workers: workers,
	}
}

func (s *Scanner) Run(words []string, targetURL string) (<-chan Result, <-chan int) {
	jobs := make(chan Job)
	results := make(chan Result)
	progress := make(chan int)

	var wg sync.WaitGroup
	var completedCount int64

	for i := 0; i < s.Workers; i++ {
		wg.Add(1)
		go s.worker(jobs, results, &wg, &completedCount)
	}

	go func() {
		for _, word := range words {
			url := fmt.Sprintf("%s/%s", strings.TrimSuffix(targetURL, "/"), word)
			jobs <- Job{
				Word: word,
				URL:  url,
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		defer close(progress)

		var last int64 = -1
		for {
			current := atomic.LoadInt64(&completedCount)
			if current != last {
				progress <- int(current)
				last = current
			}
			if current >= int64(len(words)) {
				return
			}
		}
	}()

	return results, progress
}

func (s *Scanner) worker(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup, counter *int64) {
	defer wg.Done()
	for job := range jobs {
		resp, err := s.Client.Get(job.URL)
		atomic.AddInt64(counter, 1)

		if err != nil {
			results <- Result{
				Word: job.Word,
				URL:  job.URL,
				Err:  err,
			}
			continue
		}

		status := resp.StatusCode
		results <- Result{
			Word:       job.Word,
			URL:        job.URL,
			StatusCode: status,
		}

		_ = resp.Body.Close()
	}
}
