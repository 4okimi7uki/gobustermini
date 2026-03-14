package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Result struct {
	URL    string
	Status int
}

func worker(jobs <-chan string, results chan<- Result, wg *sync.WaitGroup, client *http.Client, counter *int64) {
	defer wg.Done()
	for url := range jobs {
		resp, err := client.Get(url)

		atomic.AddInt64(counter, 1)

		if err != nil {
			continue
		}

		if resp != nil {
			if resp.StatusCode == http.StatusOK {
				results <- Result{URL: url, Status: resp.StatusCode}
			}
			resp.Body.Close()
		}
	}
}

func countLine(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}

	return count, nil
}

func main() {
	targetURL := "https://133.18.178.100"
	wordListPath := "./wordlist.txt"
	workerCount := 20

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 5 * time.Second,
	}

	total, err := countLine(wordListPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	jobs := make(chan string)
	results := make(chan Result)
	var wg sync.WaitGroup
	var completedCount int64

	for range workerCount {
		wg.Add(1)
		go worker(jobs, results, &wg, client, &completedCount)
	}

	go func() {
		for res := range results {
			fmt.Fprintf(os.Stdout, "\r\033[K[+] Found: %s (Status: %d)\n", res.URL, res.Status)
		}
	}()

	go func() {
		for {
			current := atomic.LoadInt64(&completedCount)
			percent := float64(current) / float64(total) * 100

			fmt.Fprintf(os.Stderr, "\r\033[KProgress: %.2f%% (%d/%d)", percent, current, total)

			if current >= int64(total) {
				fmt.Fprint(os.Stderr, "\r\033[KScan Completed!\n")
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	file, err := os.Open(wordListPath)
	if err != nil {
		fmt.Printf("Error opening wordlist %v\n ", err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		url := fmt.Sprintf("%s/%s", strings.TrimSuffix(targetURL, "/"), word)
		jobs <- url
	}
	close(jobs)

	if err := scanner.Err(); err != nil {
		fmt.Printf("error: %v\n", err)
	}

	wg.Wait()
	close(results)

	time.Sleep(200 * time.Millisecond)
}
