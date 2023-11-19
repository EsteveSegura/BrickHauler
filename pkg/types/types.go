package types

import (
	"net/url"
)

// TYPE HttpMethod
// Enum for verbs
const (
	HEAD    = "HEAD"
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	PATCH   = "PATCH"
	DELETE  = "DELETE"
	CONNECT = "CONNECT"
	OPTIONS = "OPTIONS"
	TRACE   = "TRACE"
)

// HttpMethod type
type HttpMethod string

// validator for HttpMethod type
func IsHttpMethod(method HttpMethod) bool {
	switch method {
	case HEAD, GET, POST, PUT, PATCH, DELETE, CONNECT, OPTIONS, TRACE:
		return true
	}
	return false
}

// TYPE URI
type Uri string

func IsUri(uri Uri) bool {
	_, err := url.ParseRequestURI(string(uri))

	if err != nil {
		return false
	}
	return true
}
