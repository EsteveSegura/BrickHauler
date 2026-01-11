package config

import (
	"fmt"
	"net/url"
)

type URI struct {
	raw    string
	parsed *url.URL
}

// NewURI validates and constructs a URI.
func NewURI(s string) (URI, error) {
	if s == "" {
		return URI{}, fmt.Errorf("URI cannot be empty")
	}

	parsed, err := url.ParseRequestURI(s)
	if err != nil {
		return URI{}, fmt.Errorf("invalid URI %q: %w", s, err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return URI{}, fmt.Errorf("URI must use http or https scheme: %q", s)
	}

	return URI{raw: s, parsed: parsed}, nil
}

// String returns the raw URI string.
func (u URI) String() string {
	return u.raw
}

// URL returns the parsed URL.
func (u URI) URL() *url.URL {
	return u.parsed
}
