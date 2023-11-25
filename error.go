package rest

import (
	"net/http"
)

// REST API error.
var (
	Err = Error{}

	ErrMovedPermanently  = NewError(http.StatusMovedPermanently, "")  // 301
	ErrFound             = NewError(http.StatusFound, "")             // 302
	ErrTemporaryRedirect = NewError(http.StatusTemporaryRedirect, "") // 307
	ErrPermanentRedirect = NewError(http.StatusPermanentRedirect, "") // 308

	ErrBadRequest       = NewError(http.StatusBadRequest, "")       // 400
	ErrUnauthorized     = NewError(http.StatusUnauthorized, "")     // 401
	ErrPaymentRequired  = NewError(http.StatusPaymentRequired, "")  // 402
	ErrForbidden        = NewError(http.StatusForbidden, "")        // 403
	ErrNotFound         = NewError(http.StatusNotFound, "")         // 404
	ErrMethodNotAllowed = NewError(http.StatusMethodNotAllowed, "") // 405
	ErrNotAcceptable    = NewError(http.StatusNotAcceptable, "")    // 406

	ErrInternalServerError = NewError(http.StatusInternalServerError, "") // 500
	ErrNotImplemented      = NewError(http.StatusNotImplemented, "")      // 501
	ErrBadGateway          = NewError(http.StatusBadGateway, "")          // 502
	ErrServiceUnavailable  = NewError(http.StatusServiceUnavailable, "")  // 503
	ErrGatewayTimeout      = NewError(http.StatusGatewayTimeout, "")      // 504
)

// Error represents a REST API error.
// It can be marshaled to JSON with ease and provides a standard format for printing errors and additional data.
type Error struct {
	StatusCode int                    `json:"-"`              // HTTP status code (200, 404, 500 etc.)
	Message    string                 `json:"message"`        // Status message ("OK", "Not found", "Internal server error" etc.)
	Data       map[string]interface{} `json:"data,omitempty"` // Optional additional data.
}

// Error retrieves the message of a REST API error.
func (e Error) Error() string {
	return e.Message
}

// Is determines whether the Error is an instance of the target.
// https://pkg.go.dev/errors#Is
//
// If the target is a REST API error and specifies a status code, this function returns true if the status codes match.
// If the target is an empty REST API error, this function always returns true.
func (e Error) Is(target error) bool {
	if t, ok := target.(Error); ok {
		return t.StatusCode == e.StatusCode || t.StatusCode == 0
	}
	return false
}

// WithData returns a copy of the HTTP error with the given data merged in.
func (e Error) WithData(data map[string]interface{}) Error {
	if e.Data == nil {
		e.Data = map[string]any{}
	}
	if data != nil {
		for key, value := range data {
			e.Data[key] = value
		}
	}
	return e
}

// WithError returns a copy of the HTTP error with the given error added either as the Message, if it empty, or as additional data.
func (e Error) WithError(err error) Error {
	if e.Message == "" {
		return e.WithMessage(err.Error())
	}

	return e.WithData(map[string]interface{}{
		"error": err.Error(),
	})
}

// WithMessage returns a copy of the HTTP error with the given message.
func (e Error) WithMessage(message string) Error {
	e.Message = message
	return e
}

// WithValue returns a copy of the HTTP error with a single data value added.
func (e Error) WithValue(name string, value any) Error {
	return e.WithData(map[string]any{
		name: value,
	})
}

// Write writes the HTTP error to an HTTP response as plain text.
// Additional data is omitted.
func (e Error) Write(w http.ResponseWriter) (int, error) {
	if e.StatusCode == 0 {
		e.StatusCode = 200
	}
	w.WriteHeader(e.StatusCode)
	return w.Write([]byte(e.Message))
}

// WriteJSON writes the HTTP error to an HTTP response as JSON.
func (e Error) WriteJSON(w http.ResponseWriter) error {
	if e.StatusCode == 0 {
		e.StatusCode = 200
	}
	return WriteResponseJSON(w, e.StatusCode, e)
}

// NewError creates a new REST API error.
// If the message is empty, the standard text provided by http.StatusText is substituted.
func NewError(statusCode int, message string) Error {
	if len(message) == 0 {
		message = http.StatusText(statusCode)
	}
	return Error{
		StatusCode: statusCode,
		Message:    message,
	}
}
