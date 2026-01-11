package runner

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/EsteveSegura/BrickHauler/internal/config"
	"github.com/EsteveSegura/BrickHauler/internal/httpclient"
	"github.com/EsteveSegura/BrickHauler/internal/metrics"
	"github.com/EsteveSegura/BrickHauler/internal/output"
	"github.com/EsteveSegura/BrickHauler/internal/version"
)

// Runner executes load tests.
type Runner struct {
	cfg     *config.Config
	client  *http.Client
	metrics *metrics.Metrics
	output  *output.Writer
}

// New creates a new Runner.
func New(cfg *config.Config, w io.Writer) *Runner {
	return &Runner{
		cfg: cfg,
		client: httpclient.New(httpclient.Config{
			ProxyURL: cfg.ProxyURL,
			Timeout:  30 * time.Second,
		}),
		metrics: metrics.New(cfg.Requests),
		output:  output.New(w),
	}
}

// Run executes the load test with graceful shutdown support.
func (r *Runner) Run(ctx context.Context) error {
	startTime := time.Now()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	requestsPerWorker := r.cfg.RequestsPerWorker()

	// Launch worker goroutines
	for i := 0; i < r.cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.worker(ctx, requestsPerWorker)
		}()
	}

	// Progress reporting goroutine for live feed
	var progressWg sync.WaitGroup
	if r.cfg.LiveFeed {
		progressWg.Add(1)
		go func() {
			defer progressWg.Done()
			r.progressReporter(ctx, startTime)
		}()
	}

	// Wait for all workers to complete
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Handle completion or cancellation
	var err error
	select {
	case <-done:
		// Normal completion
	case <-ctx.Done():
		// Cancelled - wait for workers to finish current requests
		<-done
		err = ctx.Err()
	}

	// Stop progress reporter
	cancel()
	progressWg.Wait()

	if r.cfg.LiveFeed {
		r.output.PrintNewline()
	}

	duration := time.Since(startTime)
	r.output.PrintResults(r.cfg, r.metrics.Snapshot(), duration)

	return err
}

// worker sends requests for a single virtual user.
func (r *Runner) worker(ctx context.Context, numRequests int) {
	for i := 0; i < numRequests; i++ {
		select {
		case <-ctx.Done():
			return // Graceful shutdown
		default:
			r.sendRequest(ctx)
		}
	}
}

// sendRequest sends a single HTTP request and records metrics.
func (r *Runner) sendRequest(ctx context.Context) {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, r.cfg.Method.String(), r.cfg.URI.String(), nil)
	if err != nil {
		r.metrics.RecordFailure()
		return
	}

	req.Header.Set("User-Agent", version.UserAgent)

	for _, cookie := range r.cfg.Cookies {
		req.AddCookie(cookie)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		r.metrics.RecordFailure()
		return
	}
	defer resp.Body.Close()

	// Drain body to enable connection reuse
	io.Copy(io.Discard, resp.Body)

	duration := time.Since(start)

	if resp.StatusCode < 400 {
		r.metrics.RecordSuccess(duration)
	} else {
		r.metrics.RecordFailure()
	}
}

// progressReporter periodically prints progress.
func (r *Runner) progressReporter(ctx context.Context, startTime time.Time) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			snap := r.metrics.Snapshot()
			r.output.PrintProgress(snap.TotalRequests(), int64(r.cfg.Requests), time.Since(startTime))
		}
	}
}
