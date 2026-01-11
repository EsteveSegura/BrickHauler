package httpclient

import (
	"net/http"
	"net/url"
	"time"
)

// Config for HTTP client creation.
type Config struct {
	ProxyURL *url.URL
	Timeout  time.Duration
}

// New creates a configured HTTP client with connection pooling.
func New(cfg Config) *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     100,
		IdleConnTimeout:     90 * time.Second,
	}

	if cfg.ProxyURL != nil {
		transport.Proxy = http.ProxyURL(cfg.ProxyURL)
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
}
