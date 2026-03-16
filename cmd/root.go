package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/4okimi7uki/gobustermini/internal/client"
	"github.com/4okimi7uki/gobustermini/internal/scanner"

	"github.com/4okimi7uki/gobustermini/internal/wordlist"
	"github.com/spf13/cobra"
)

type Config struct {
	TargerURL   string
	Wordlist    string
	Workers     int
	Timeout     time.Duration
	TnsecureTLS bool
}

var rootCmd = &cobra.Command{
	Use:   "gobustermini",
	Short: "",
	RunE: func(cmd *cobra.Command, arg []string) error {
		targetURL := "https://133.18.178.100"
		wordListPath := "./wordlist.txt"
		workerCount := 20

		httpClient := client.New()
		words, err := wordlist.ReadWordList(wordListPath)
		if err != nil {
			return err
		}

		s := scanner.New(httpClient, workerCount)
		results, _ := s.Run(words, targetURL)

		for res := range results {
			if res.Err != nil {
				fmt.Printf("error: %v\n", res.Err)
				continue
			}
			resolvedResult := scanner.FormatResult(res)

			if res.StatusCode < 500 {
				fmt.Printf("\r\033[K%s\n", resolvedResult)
			}

		}

		return nil
	},
}

func Excute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.Flags()
	// rootCmd.PersistentFlags().StringVarP(&e)
}
