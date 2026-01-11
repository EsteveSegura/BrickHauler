package output

import (
	"fmt"
	"io"
	"time"

	"github.com/EsteveSegura/BrickHauler/internal/config"
	"github.com/EsteveSegura/BrickHauler/internal/metrics"
	"github.com/EsteveSegura/BrickHauler/internal/version"
)

// Writer handles all output formatting.
type Writer struct {
	w io.Writer
}

// New creates a new output Writer.
func New(w io.Writer) *Writer {
	return &Writer{w: w}
}

// PrintResults outputs the final test results.
func (w *Writer) PrintResults(cfg *config.Config, snap metrics.Snapshot, duration time.Duration) {
	fmt.Fprintf(w.w, "\nBrickHauler %s\n", version.Version)
	fmt.Fprintf(w.w, "================================================\n\n")

	fmt.Fprintf(w.w, "Target URL:              %s\n", cfg.URI)
	fmt.Fprintf(w.w, "HTTP Method:             %s\n", cfg.Method)
	fmt.Fprintf(w.w, "Concurrency:             %d\n", cfg.Concurrency)
	fmt.Fprintf(w.w, "Total Requests:          %d\n\n", cfg.Requests)

	fmt.Fprintf(w.w, "Results:\n")
	fmt.Fprintf(w.w, "--------\n")
	fmt.Fprintf(w.w, "Successful:              %d\n", snap.SuccessCount)
	fmt.Fprintf(w.w, "Failed:                  %d\n", snap.FailureCount)

	totalRequests := snap.TotalRequests()
	if totalRequests > 0 {
		rps := float64(totalRequests) / duration.Seconds()
		fmt.Fprintf(w.w, "Requests/sec:            %.2f\n", rps)
	}

	if snap.SuccessCount > 0 {
		fmt.Fprintf(w.w, "Avg Response Time:       %v\n", snap.AverageTime())
	}

	fmt.Fprintf(w.w, "Total Duration:          %v\n\n", duration.Round(time.Millisecond))

	w.printPercentiles(snap)
}

func (w *Writer) printPercentiles(snap metrics.Snapshot) {
	if len(snap.Durations) == 0 {
		return
	}

	fmt.Fprintf(w.w, "Response Time Percentiles:\n")
	fmt.Fprintf(w.w, "--------------------------\n")

	percentiles := []float64{50, 66, 75, 80, 90, 95, 98, 99, 100}
	for _, p := range percentiles {
		fmt.Fprintf(w.w, "  %3.0f%%  <= %v\n", p, snap.Percentile(p))
	}
	fmt.Fprintln(w.w)
}

// PrintProgress outputs real-time progress during the test.
func (w *Writer) PrintProgress(completed, total int64, duration time.Duration) {
	rps := float64(completed) / duration.Seconds()
	fmt.Fprintf(w.w, "\rProgress: %d/%d requests (%.1f req/s)", completed, total, rps)
}

// PrintShutdown outputs a shutdown message.
func (w *Writer) PrintShutdown() {
	fmt.Fprintln(w.w, "\nShutting down gracefully...")
}

// PrintNewline outputs a newline.
func (w *Writer) PrintNewline() {
	fmt.Fprintln(w.w)
}
