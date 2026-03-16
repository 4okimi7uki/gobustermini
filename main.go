package main

import "github.com/4okimi7uki/gobustermini/cmd"

func main() {
	cmd.Excute()
}

// import (
// 	"bufio"
// 	"crypto/tls"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"sync"
// 	"sync/atomic"
// 	"time"

// 	"github.com/fatih/color"
// )

// var (
// 	white  = color.New(color.FgWhite).SprintfFunc()
// 	yellow = color.New(color.FgYellow).SprintfFunc()
// 	green  = color.New(color.FgGreen).SprintfFunc()
// 	blue   = color.New(color.FgBlue).SprintfFunc()
// 	red    = color.New(color.FgRed).SprintfFunc()
// 	cyan   = color.New(color.FgCyan).SprintfFunc()
// )

// func worker(jobs <-chan string, results chan<- Result, wg *sync.WaitGroup, client *http.Client, counter *int64) {
// 	defer wg.Done()
// 	for url := range jobs {
// 		resp, err := client.Get(url)

// 		// counter
// 		atomic.AddInt64(counter, 1)
// 		if err != nil {
// 			continue
// 		}

// 		statusCodeColor := white

// 		if resp != nil {
// 			switch {
// 			case resp.StatusCode == http.StatusOK:
// 				statusCodeColor = green
// 			case resp.StatusCode >= 300 && resp.StatusCode < 400:
// 				statusCodeColor = cyan
// 			case resp.StatusCode >= 400 && resp.StatusCode < 500:
// 				statusCodeColor = yellow
// 			case resp.StatusCode >= 500 && resp.StatusCode < 600:
// 				statusCodeColor = red
// 			}
// 			statusCode := statusCodeColor(fmt.Sprintf("Status: %d", resp.StatusCode))
// 			results <- Result{URL: url, Status: statusCode}

// 			resp.Body.Close()
// 		}
// 	}
// }

// func main() {
// 	targetURL := "https://133.18.178.100"
// 	wordListPath := "./wordlist.txt"
// 	workerCount := 20

// 	client := &http.Client{
// 		Transport: &http.Transport{
// 			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 		},
// 		Timeout: 5 * time.Second,
// 	}

// 	total, err := countLine(wordListPath)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
// 		return
// 	}

// 	jobs := make(chan string)
// 	results := make(chan Result)
// 	var wg sync.WaitGroup
// 	var completedCount int64

// 	for range workerCount {
// 		wg.Add(1)
// 		go worker(jobs, results, &wg, client, &completedCount)
// 	}

// 	go func() {
// 		for res := range results {
// 			fmt.Fprintf(os.Stdout, "\r\033[K %s %s\n", res.URL, res.Status)
// 		}
// 	}()

// 	go func() {
// 		for {
// 			current := atomic.LoadInt64(&completedCount)
// 			percent := float64(current) / float64(total) * 100

// 			fmt.Fprintf(os.Stderr, "\r\033[KProgress: %.2f%% (%d/%d)", percent, current, total)

// 			if current >= int64(total) {
// 				fmt.Fprint(os.Stderr, "\r\033[KScan Completed!\n")
// 				break
// 			}
// 			time.Sleep(100 * time.Millisecond)
// 		}
// 	}()

// 	file, err := os.Open(wordListPath)
// 	if err != nil {
// 		fmt.Printf("Error opening wordlist %v\n ", err)
// 		return
// 	}
// 	defer func() {
// 		_ = file.Close()
// 	}()

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		word := scanner.Text()
// 		url := fmt.Sprintf("%s/%s", strings.TrimSuffix(targetURL, "/"), word)
// 		jobs <- url
// 	}
// 	close(jobs)

// 	if err := scanner.Err(); err != nil {
// 		fmt.Printf("error: %v\n", err)
// 	}

// 	wg.Wait()
// 	close(results)

// 	time.Sleep(200 * time.Millisecond)
// }
