package runner

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/EsteveSegura/BrickHauler/internal/config"
)

func TestRunner_SuccessfulRun(t *testing.T) {
	var requestCount int64

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&requestCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	uri, err := config.NewURI(server.URL)
	if err != nil {
		t.Fatalf("failed to create URI: %v", err)
	}

	cfg := &config.Config{
		URI:         uri,
		Method:      config.MethodGET,
		Concurrency: 2,
		Requests:    10,
	}

	r := New(cfg, io.Discard)
	err = r.Run(context.Background())

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if atomic.LoadInt64(&requestCount) != 10 {
		t.Errorf("expected 10 requests, got %d", atomic.LoadInt64(&requestCount))
	}
}

func TestRunner_GracefulShutdown(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	uri, err := config.NewURI(server.URL)
	if err != nil {
		t.Fatalf("failed to create URI: %v", err)
	}

	cfg := &config.Config{
		URI:         uri,
		Method:      config.MethodGET,
		Concurrency: 5,
		Requests:    100,
	}

	r := New(cfg, io.Discard)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	err = r.Run(ctx)

	if err != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
}

func TestRunner_UserAgentHeader(t *testing.T) {
	var userAgent string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	uri, err := config.NewURI(server.URL)
	if err != nil {
		t.Fatalf("failed to create URI: %v", err)
	}

	cfg := &config.Config{
		URI:         uri,
		Method:      config.MethodGET,
		Concurrency: 1,
		Requests:    1,
	}

	r := New(cfg, io.Discard)
	_ = r.Run(context.Background())

	if userAgent != "BrickHauler/0.2.0" {
		t.Errorf("User-Agent = %q, want %q", userAgent, "BrickHauler/0.2.0")
	}
}

func TestRunner_CookiesAreSent(t *testing.T) {
	var receivedCookies []*http.Cookie

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedCookies = r.Cookies()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	uri, err := config.NewURI(server.URL)
	if err != nil {
		t.Fatalf("failed to create URI: %v", err)
	}

	cfg := &config.Config{
		URI:         uri,
		Method:      config.MethodGET,
		Concurrency: 1,
		Requests:    1,
		Cookies: []*http.Cookie{
			{Name: "foo", Value: "bar"},
			{Name: "baz", Value: "qux"},
		},
	}

	r := New(cfg, io.Discard)
	_ = r.Run(context.Background())

	if len(receivedCookies) != 2 {
		t.Errorf("expected 2 cookies, got %d", len(receivedCookies))
	}
}

func TestRunner_HTTPMethod(t *testing.T) {
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	uri, err := config.NewURI(server.URL)
	if err != nil {
		t.Fatalf("failed to create URI: %v", err)
	}

	cfg := &config.Config{
		URI:         uri,
		Method:      config.MethodPOST,
		Concurrency: 1,
		Requests:    1,
	}

	r := New(cfg, io.Discard)
	_ = r.Run(context.Background())

	if receivedMethod != "POST" {
		t.Errorf("Method = %q, want %q", receivedMethod, "POST")
	}
}

func TestRunner_FailedRequests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	uri, err := config.NewURI(server.URL)
	if err != nil {
		t.Fatalf("failed to create URI: %v", err)
	}

	cfg := &config.Config{
		URI:         uri,
		Method:      config.MethodGET,
		Concurrency: 1,
		Requests:    5,
	}

	r := New(cfg, io.Discard)
	_ = r.Run(context.Background())

	snap := r.metrics.Snapshot()
	if snap.FailureCount != 5 {
		t.Errorf("FailureCount = %d, want 5", snap.FailureCount)
	}
	if snap.SuccessCount != 0 {
		t.Errorf("SuccessCount = %d, want 0", snap.SuccessCount)
	}
}
