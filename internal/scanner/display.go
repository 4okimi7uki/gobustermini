package scanner

import (
	"fmt"
	"net/http"

	"github.com/fatih/color"
)

var (
	white  = color.New(color.FgWhite).SprintfFunc()
	yellow = color.New(color.FgYellow).SprintfFunc()
	green  = color.New(color.FgGreen).SprintfFunc()
	blue   = color.New(color.FgBlue).SprintfFunc()
	red    = color.New(color.FgRed).SprintfFunc()
	cyan   = color.New(color.FgCyan).SprintfFunc()
)

func FormatResult(result Result) string {
	const wordWidth = 40

	if result.Err != nil {
		return fmt.Sprintf("%-*s %s", wordWidth, white(result.Word), red(result.Err.Error()))
	}

	statusCodeColor := white
	switch {
	case result.StatusCode == http.StatusOK:
		statusCodeColor = green
	case result.StatusCode >= 300 && result.StatusCode < 400:
		statusCodeColor = cyan
	case result.StatusCode >= 400 && result.StatusCode < 500:
		statusCodeColor = yellow
	case result.StatusCode >= 500 && result.StatusCode < 600:
		statusCodeColor = red
	}

	status := statusCodeColor(fmt.Sprintf(" ( Status: %d )", result.StatusCode))

	return fmt.Sprintf("%-*s %s", wordWidth, white(result.Word), status)
}
