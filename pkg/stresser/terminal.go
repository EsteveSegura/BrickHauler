package stresser

import (
	"fmt"
	"load-tester/pkg/types"
	"os"
	"os/exec"
	"runtime"
	"time"

	banner "moul.io/banner"
)

// printStatistics prints the final statistics of the stress test.
func PrintStatistics(url types.Uri, virtualUsers, totalRequests int, duration time.Duration) {
	fmt.Println(banner.Inline("brick hauler") + "    0.1.0\n")
	fmt.Printf("URL:                             %s\n", url)
	fmt.Printf("Concurrency:                     %d\n", virtualUsers)
	fmt.Printf("Requests:                        %d\n", totalRequests)
	fmt.Println("")
	fmt.Println("")
	fmt.Printf("Completed Requests:              %d\n", successfulRequests)
	fmt.Printf("Failed Requests:                 %d\n", failedRequests)
	fmt.Printf("Requests Per Second:             %.2f\n", float64(totalRequests)/duration.Seconds())
	fmt.Printf("Time Per Request:                %.2f ms\n", float64(totalRequestTime.Milliseconds())/float64(totalRequests))
	fmt.Printf("Time to Complete Test:           %v\n", duration)
}

// printPrematureStats prints the current statistics of the ongoing stress test.
func PrintPrematureStats(url types.Uri, virtualUsers, totalRequests int) {
	fmt.Println(banner.Inline("brick hauler") + "    0.1.0\n")
	fmt.Printf("URL:                             %s\n", url)
	fmt.Printf("Concurrency:                     %d\n", virtualUsers)
	fmt.Printf("Requests:                        %d\n", totalRequests)
	fmt.Println("")
	fmt.Println("")
	fmt.Printf("Completed Requests:              %d\n", successfulRequests)
	fmt.Printf("Failed Requests:                 %d\n", failedRequests)
	fmt.Println("Requests Per Second:            N/A")
	fmt.Println("Time Per Request:               N/A")
	fmt.Println("Time to Complete Test:          N/A")
}

// printPercentiles prints the percentile distribution of the request durations.
func PrintPercentiles() {
	percentiles := []float64{50, 66, 75, 80, 90, 95, 98, 99, 100}
	fmt.Println("")
	fmt.Println("Response time percentiles (ms):")
	for _, percentile := range percentiles {
		pIndex := int((percentile / 100) * float64(len(requestDurations)))
		if pIndex > 0 {
			pIndex--
		}
		if len(requestDurations) != 0 {
			if percentile == 100 {
				fmt.Printf("%v%%            <=           %v ms\n", percentile, requestDurations[pIndex].Milliseconds())
			} else {
				fmt.Printf("%v%%             <=           %v ms\n", percentile, requestDurations[pIndex].Milliseconds())
			}
		}
	}
	fmt.Println("")
}

// clearScreen clears the terminal screen based on the OS.
func ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}
