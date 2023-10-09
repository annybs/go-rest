package rest

import "net/http"

// Err is an empty error, representing any REST API Error.
var Err = Error{}

// ErrBadRequest represents HTTP 400 Bad Request.
var ErrBadRequest = Error{
	StatusCode: http.StatusBadRequest,
	Message:    "Bad request",
}

// ErrInternalServerError represents HTTP 500 Internal Server Error.
var ErrInternalServerError = Error{
	StatusCode: http.StatusInternalServerError,
	Message:    "Internal server error",
}

// ErrNotFound represents HTTP 404 Not Found.
var ErrNotFound = Error{
	StatusCode: http.StatusNotFound,
	Message:    "Not found",
}

// ErrUnavailable represents HTTP 503 Service Unavailable.
var ErrUnavailable = Error{
	StatusCode: http.StatusServiceUnavailable,
	Message:    "Service unavailable",
}
