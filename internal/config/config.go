package config

import (
	"fmt"
	"net/http"
	"net/url"
)

// Config holds all configuration for a load test run.
type Config struct {
	URI         URI
	Method      HTTPMethod
	Concurrency int
	Requests    int
	Cookies     []*http.Cookie
	ProxyURL    *url.URL
	LiveFeed    bool
}

// Validate checks all configuration values.
func (c *Config) Validate() error {
	if c.Concurrency <= 0 {
		return fmt.Errorf("concurrency must be greater than 0, got %d", c.Concurrency)
	}

	if c.Requests <= 0 {
		return fmt.Errorf("requests must be greater than 0, got %d", c.Requests)
	}

	if c.Requests%c.Concurrency != 0 {
		return fmt.Errorf(
			"requests (%d) must be evenly divisible by concurrency (%d)",
			c.Requests, c.Concurrency,
		)
	}

	if !c.Method.IsValid() {
		return fmt.Errorf("invalid HTTP method: %s", c.Method)
	}

	return nil
}

// RequestsPerWorker returns how many requests each worker should make.
func (c *Config) RequestsPerWorker() int {
	return c.Requests / c.Concurrency
}
