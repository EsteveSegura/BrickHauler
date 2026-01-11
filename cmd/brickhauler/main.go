package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/EsteveSegura/BrickHauler/internal/config"
	"github.com/EsteveSegura/BrickHauler/internal/runner"
	"github.com/EsteveSegura/BrickHauler/internal/version"
)

// stringSlice implements flag.Value for repeatable string flags.
type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		method      string
		uri         string
		concurrency int
		requests    int
		cookies     stringSlice
		proxy       string
		liveFeed    bool
		showVersion bool
	)

	flag.StringVar(&method, "verb", "GET", "HTTP method (GET, POST, PUT, PATCH, DELETE, etc.)")
	flag.StringVar(&uri, "uri", "", "Target URL for load testing")
	flag.IntVar(&concurrency, "concurrent", 0, "Number of concurrent virtual users")
	flag.IntVar(&requests, "request", 0, "Total number of requests to send")
	flag.Var(&cookies, "cookie", "Cookie in name=value format (repeatable)")
	flag.StringVar(&proxy, "proxy", "", "HTTP proxy URL")
	flag.BoolVar(&liveFeed, "feed", false, "Show real-time progress")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.BoolVar(&showVersion, "version", false, "Show version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "BrickHauler %s - HTTP Load Testing Tool\n\n", version.Version)
		fmt.Fprintf(os.Stderr, "Usage: brickhauler [options]\n\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if showVersion {
		fmt.Printf("BrickHauler %s\n", version.Version)
		return nil
	}

	// Validate required flags
	if uri == "" {
		return fmt.Errorf("--uri is required")
	}
	if concurrency == 0 {
		return fmt.Errorf("--concurrent is required")
	}
	if requests == 0 {
		return fmt.Errorf("--request is required")
	}

	// Build and validate config
	cfg, err := buildConfig(method, uri, concurrency, requests, cookies, proxy, liveFeed)
	if err != nil {
		return err
	}

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Fprintln(os.Stderr, "\nShutting down gracefully...")
		cancel()
	}()

	// Run the load test
	r := runner.New(cfg, os.Stdout)
	return r.Run(ctx)
}

func buildConfig(method, uri string, concurrency, requests int, cookies stringSlice, proxy string, liveFeed bool) (*config.Config, error) {
	parsedMethod, err := config.ParseHTTPMethod(method)
	if err != nil {
		return nil, err
	}

	parsedURI, err := config.NewURI(uri)
	if err != nil {
		return nil, err
	}

	parsedCookies, err := config.ParseCookies(cookies)
	if err != nil {
		return nil, err
	}

	var proxyURL *url.URL
	if proxy != "" {
		proxyURL, err = url.Parse(proxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}
	}

	cfg := &config.Config{
		URI:         parsedURI,
		Method:      parsedMethod,
		Concurrency: concurrency,
		Requests:    requests,
		Cookies:     parsedCookies,
		ProxyURL:    proxyURL,
		LiveFeed:    liveFeed,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
