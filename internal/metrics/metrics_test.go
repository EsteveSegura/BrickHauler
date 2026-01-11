package metrics

import (
	"sync"
	"testing"
	"time"
)

func TestMetrics_RecordSuccess(t *testing.T) {
	m := New(10)

	m.RecordSuccess(100 * time.Millisecond)
	m.RecordSuccess(200 * time.Millisecond)

	snap := m.Snapshot()

	if snap.SuccessCount != 2 {
		t.Errorf("SuccessCount = %d, want 2", snap.SuccessCount)
	}
	if snap.TotalTime != 300*time.Millisecond {
		t.Errorf("TotalTime = %v, want 300ms", snap.TotalTime)
	}
	if len(snap.Durations) != 2 {
		t.Errorf("len(Durations) = %d, want 2", len(snap.Durations))
	}
}

func TestMetrics_RecordFailure(t *testing.T) {
	m := New(10)

	m.RecordFailure()
	m.RecordFailure()
	m.RecordFailure()

	snap := m.Snapshot()

	if snap.FailureCount != 3 {
		t.Errorf("FailureCount = %d, want 3", snap.FailureCount)
	}
	if snap.SuccessCount != 0 {
		t.Errorf("SuccessCount = %d, want 0", snap.SuccessCount)
	}
}

func TestMetrics_ConcurrentAccess(t *testing.T) {
	m := New(1000)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				m.RecordSuccess(100 * time.Millisecond)
			}
		}()
	}
	wg.Wait()

	snap := m.Snapshot()
	if snap.SuccessCount != 1000 {
		t.Errorf("SuccessCount = %d, want 1000", snap.SuccessCount)
	}
	if len(snap.Durations) != 1000 {
		t.Errorf("len(Durations) = %d, want 1000", len(snap.Durations))
	}
}

func TestMetrics_ConcurrentMixed(t *testing.T) {
	m := New(2000)

	var wg sync.WaitGroup

	// 100 goroutines recording successes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				m.RecordSuccess(50 * time.Millisecond)
			}
		}()
	}

	// 100 goroutines recording failures
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				m.RecordFailure()
			}
		}()
	}

	wg.Wait()

	snap := m.Snapshot()
	if snap.SuccessCount != 1000 {
		t.Errorf("SuccessCount = %d, want 1000", snap.SuccessCount)
	}
	if snap.FailureCount != 1000 {
		t.Errorf("FailureCount = %d, want 1000", snap.FailureCount)
	}
}

func TestSnapshot_TotalRequests(t *testing.T) {
	snap := Snapshot{
		SuccessCount: 80,
		FailureCount: 20,
	}
	if snap.TotalRequests() != 100 {
		t.Errorf("TotalRequests() = %d, want 100", snap.TotalRequests())
	}
}

func TestSnapshot_Percentile(t *testing.T) {
	m := New(100)
	for i := 1; i <= 100; i++ {
		m.RecordSuccess(time.Duration(i) * time.Millisecond)
	}

	snap := m.Snapshot()

	// P50 should be around 50ms
	p50 := snap.Percentile(50)
	if p50 < 49*time.Millisecond || p50 > 51*time.Millisecond {
		t.Errorf("P50 = %v, expected ~50ms", p50)
	}

	// P90 should be around 90ms
	p90 := snap.Percentile(90)
	if p90 < 89*time.Millisecond || p90 > 91*time.Millisecond {
		t.Errorf("P90 = %v, expected ~90ms", p90)
	}

	// P100 should be 100ms
	p100 := snap.Percentile(100)
	if p100 != 100*time.Millisecond {
		t.Errorf("P100 = %v, expected 100ms", p100)
	}
}

func TestSnapshot_Percentile_Empty(t *testing.T) {
	snap := Snapshot{Durations: []time.Duration{}}
	if snap.Percentile(50) != 0 {
		t.Errorf("Percentile(50) on empty = %v, want 0", snap.Percentile(50))
	}
}

func TestSnapshot_AverageTime(t *testing.T) {
	snap := Snapshot{
		SuccessCount: 4,
		TotalTime:    400 * time.Millisecond,
	}
	if snap.AverageTime() != 100*time.Millisecond {
		t.Errorf("AverageTime() = %v, want 100ms", snap.AverageTime())
	}
}

func TestSnapshot_AverageTime_Zero(t *testing.T) {
	snap := Snapshot{
		SuccessCount: 0,
		TotalTime:    0,
	}
	if snap.AverageTime() != 0 {
		t.Errorf("AverageTime() on zero success = %v, want 0", snap.AverageTime())
	}
}

func TestSnapshot_DurationsAreSorted(t *testing.T) {
	m := New(5)

	// Add durations in random order
	m.RecordSuccess(300 * time.Millisecond)
	m.RecordSuccess(100 * time.Millisecond)
	m.RecordSuccess(500 * time.Millisecond)
	m.RecordSuccess(200 * time.Millisecond)
	m.RecordSuccess(400 * time.Millisecond)

	snap := m.Snapshot()

	for i := 1; i < len(snap.Durations); i++ {
		if snap.Durations[i] < snap.Durations[i-1] {
			t.Errorf("Durations not sorted: %v comes after %v", snap.Durations[i], snap.Durations[i-1])
		}
	}
}
