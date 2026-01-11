package config

import (
	"fmt"
	"net/http"
	"strings"
)

// ParseCookie parses a "name=value" string into an http.Cookie.
func ParseCookie(s string) (*http.Cookie, error) {
	if s == "" {
		return nil, nil
	}

	name, value, found := strings.Cut(s, "=")
	if !found {
		return nil, fmt.Errorf("invalid cookie format %q: must be name=value", s)
	}

	if name == "" {
		return nil, fmt.Errorf("cookie name cannot be empty")
	}

	return &http.Cookie{Name: name, Value: value}, nil
}

// ParseCookies parses multiple cookie strings.
func ParseCookies(ss []string) ([]*http.Cookie, error) {
	var cookies []*http.Cookie
	for _, s := range ss {
		c, err := ParseCookie(s)
		if err != nil {
			return nil, err
		}
		if c != nil {
			cookies = append(cookies, c)
		}
	}
	return cookies, nil
}
