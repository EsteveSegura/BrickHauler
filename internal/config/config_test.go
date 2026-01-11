package config

import (
	"testing"
)

func TestParseHTTPMethod(t *testing.T) {
	tests := []struct {
		input   string
		want    HTTPMethod
		wantErr bool
	}{
		{"GET", MethodGET, false},
		{"get", MethodGET, false},
		{"Post", MethodPOST, false},
		{"PUT", MethodPUT, false},
		{"PATCH", MethodPATCH, false},
		{"DELETE", MethodDELETE, false},
		{"HEAD", MethodHEAD, false},
		{"OPTIONS", MethodOPTIONS, false},
		{"CONNECT", MethodCONNECT, false},
		{"TRACE", MethodTRACE, false},
		{"INVALID", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseHTTPMethod(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHTTPMethod(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseHTTPMethod(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestHTTPMethod_IsValid(t *testing.T) {
	if !MethodGET.IsValid() {
		t.Error("GET should be valid")
	}
	if !MethodPOST.IsValid() {
		t.Error("POST should be valid")
	}
	if HTTPMethod("INVALID").IsValid() {
		t.Error("INVALID should not be valid")
	}
	if HTTPMethod("").IsValid() {
		t.Error("empty should not be valid")
	}
}

func TestHTTPMethod_String(t *testing.T) {
	if MethodGET.String() != "GET" {
		t.Errorf("MethodGET.String() = %q, want %q", MethodGET.String(), "GET")
	}
}

func TestNewURI(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid http", "http://example.com", false},
		{"valid https", "https://example.com", false},
		{"valid with path", "https://example.com/path", false},
		{"valid with query", "https://example.com?foo=bar", false},
		{"empty", "", true},
		{"invalid scheme", "ftp://example.com", true},
		{"no scheme", "example.com", true},
		{"invalid url", "://invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uri, err := NewURI(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewURI(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && uri.String() != tt.input {
				t.Errorf("URI.String() = %q, want %q", uri.String(), tt.input)
			}
		})
	}
}

func TestURI_URL(t *testing.T) {
	uri, err := NewURI("https://example.com/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u := uri.URL()
	if u == nil {
		t.Fatal("URL() returned nil")
	}
	if u.Host != "example.com" {
		t.Errorf("URL().Host = %q, want %q", u.Host, "example.com")
	}
	if u.Path != "/path" {
		t.Errorf("URL().Path = %q, want %q", u.Path, "/path")
	}
}

func TestParseCookie(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantName  string
		wantValue string
		wantErr   bool
		wantNil   bool
	}{
		{"valid", "foo=bar", "foo", "bar", false, false},
		{"empty value", "foo=", "foo", "", false, false},
		{"value with equals", "foo=bar=baz", "foo", "bar=baz", false, false},
		{"empty string", "", "", "", false, true},
		{"no equals", "foobar", "", "", true, false},
		{"empty name", "=value", "", "", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookie, err := ParseCookie(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCookie(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if tt.wantNil {
				if cookie != nil {
					t.Errorf("ParseCookie(%q) = %v, want nil", tt.input, cookie)
				}
				return
			}
			if tt.wantErr {
				return
			}
			if cookie.Name != tt.wantName {
				t.Errorf("cookie.Name = %q, want %q", cookie.Name, tt.wantName)
			}
			if cookie.Value != tt.wantValue {
				t.Errorf("cookie.Value = %q, want %q", cookie.Value, tt.wantValue)
			}
		})
	}
}

func TestParseCookies(t *testing.T) {
	cookies, err := ParseCookies([]string{"foo=bar", "baz=qux"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cookies) != 2 {
		t.Errorf("got %d cookies, want 2", len(cookies))
	}

	// Test with invalid cookie
	_, err = ParseCookies([]string{"foo=bar", "invalid"})
	if err == nil {
		t.Error("expected error for invalid cookie")
	}

	// Test empty slice
	cookies, err = ParseCookies([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cookies) != 0 {
		t.Errorf("got %d cookies, want 0", len(cookies))
	}
}

func TestConfig_Validate(t *testing.T) {
	validURI, _ := NewURI("https://example.com")

	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid",
			cfg: Config{
				URI:         validURI,
				Method:      MethodGET,
				Concurrency: 2,
				Requests:    10,
			},
			wantErr: false,
		},
		{
			name: "zero concurrency",
			cfg: Config{
				URI:         validURI,
				Method:      MethodGET,
				Concurrency: 0,
				Requests:    10,
			},
			wantErr: true,
		},
		{
			name: "negative concurrency",
			cfg: Config{
				URI:         validURI,
				Method:      MethodGET,
				Concurrency: -1,
				Requests:    10,
			},
			wantErr: true,
		},
		{
			name: "zero requests",
			cfg: Config{
				URI:         validURI,
				Method:      MethodGET,
				Concurrency: 2,
				Requests:    0,
			},
			wantErr: true,
		},
		{
			name: "requests not divisible",
			cfg: Config{
				URI:         validURI,
				Method:      MethodGET,
				Concurrency: 3,
				Requests:    10,
			},
			wantErr: true,
		},
		{
			name: "invalid method",
			cfg: Config{
				URI:         validURI,
				Method:      HTTPMethod("INVALID"),
				Concurrency: 2,
				Requests:    10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_RequestsPerWorker(t *testing.T) {
	cfg := Config{
		Concurrency: 5,
		Requests:    100,
	}
	if got := cfg.RequestsPerWorker(); got != 20 {
		t.Errorf("RequestsPerWorker() = %d, want 20", got)
	}
}
