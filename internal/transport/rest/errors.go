package rest

import "errors"

// Sentinel errors for the REST transport layer.
var (
	ErrMethodNotAllowed   = errors.New("method not allowed")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrInternalServer     = errors.New("internal server error")
)
