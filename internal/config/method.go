package config

import (
	"fmt"
	"strings"
)

type HTTPMethod string

const (
	MethodGET     HTTPMethod = "GET"
	MethodPOST    HTTPMethod = "POST"
	MethodPUT     HTTPMethod = "PUT"
	MethodPATCH   HTTPMethod = "PATCH"
	MethodDELETE  HTTPMethod = "DELETE"
	MethodHEAD    HTTPMethod = "HEAD"
	MethodOPTIONS HTTPMethod = "OPTIONS"
	MethodCONNECT HTTPMethod = "CONNECT"
	MethodTRACE   HTTPMethod = "TRACE"
)

var validMethods = map[HTTPMethod]bool{
	MethodGET:     true,
	MethodPOST:    true,
	MethodPUT:     true,
	MethodPATCH:   true,
	MethodDELETE:  true,
	MethodHEAD:    true,
	MethodOPTIONS: true,
	MethodCONNECT: true,
	MethodTRACE:   true,
}

// ParseHTTPMethod validates and returns an HTTPMethod.
func ParseHTTPMethod(s string) (HTTPMethod, error) {
	method := HTTPMethod(strings.ToUpper(s))
	if !validMethods[method] {
		return "", fmt.Errorf("invalid HTTP method: %q", s)
	}
	return method, nil
}

// String implements the Stringer interface.
func (m HTTPMethod) String() string {
	return string(m)
}

// IsValid checks if the method is valid.
func (m HTTPMethod) IsValid() bool {
	return validMethods[m]
}
