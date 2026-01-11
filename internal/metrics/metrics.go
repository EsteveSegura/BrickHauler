package metrics

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// Metrics collects load test statistics in a thread-safe manner.
type Metrics struct {
	successCount atomic.Int64
	failureCount atomic.Int64
	totalTime    atomic.Int64 // nanoseconds

	mu        sync.Mutex
	durations []time.Duration
}

// New creates a new Metrics collector with pre-allocated capacity.
func New(expectedRequests int) *Metrics {
	return &Metrics{
		durations: make([]time.Duration, 0, expectedRequests),
	}
}

// RecordSuccess records a successful request with its duration.
func (m *Metrics) RecordSuccess(d time.Duration) {
	m.successCount.Add(1)
	m.totalTime.Add(int64(d))

	m.mu.Lock()
	m.durations = append(m.durations, d)
	m.mu.Unlock()
}

// RecordFailure records a failed request.
func (m *Metrics) RecordFailure() {
	m.failureCount.Add(1)
}

// Snapshot represents a point-in-time copy of metrics.
type Snapshot struct {
	SuccessCount int64
	FailureCount int64
	TotalTime    time.Duration
	Durations    []time.Duration // sorted
}

// Snapshot returns a copy of current metrics for reporting.
func (m *Metrics) Snapshot() Snapshot {
	m.mu.Lock()
	durations := make([]time.Duration, len(m.durations))
	copy(durations, m.durations)
	m.mu.Unlock()

	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	return Snapshot{
		SuccessCount: m.successCount.Load(),
		FailureCount: m.failureCount.Load(),
		TotalTime:    time.Duration(m.totalTime.Load()),
		Durations:    durations,
	}
}

// TotalRequests returns the total number of requests (success + failure).
func (s Snapshot) TotalRequests() int64 {
	return s.SuccessCount + s.FailureCount
}

// Percentile calculates the Nth percentile from sorted durations.
func (s Snapshot) Percentile(p float64) time.Duration {
	if len(s.Durations) == 0 {
		return 0
	}
	idx := int(p / 100 * float64(len(s.Durations)))
	if idx >= len(s.Durations) {
		idx = len(s.Durations) - 1
	}
	return s.Durations[idx]
}

// AverageTime returns the average request duration.
func (s Snapshot) AverageTime() time.Duration {
	if s.SuccessCount == 0 {
		return 0
	}
	return s.TotalTime / time.Duration(s.SuccessCount)
}
