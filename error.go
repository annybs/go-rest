package rest

import (
	"net/http"
)

// Error represents a REST API error.
// It can be marshaled to JSON with ease.
type Error struct {
	StatusCode int                    `json:"statusCode"`     // HTTP status code (200, 404, 500 etc.)
	Message    string                 `json:"message"`        // Status message ("OK", "Not found", "Internal server error" etc.)
	Data       map[string]interface{} `json:"data,omitempty"` // Optional additional data.
}

// Error retrieves the message of a REST API Error.
// If it has a "error" string attached using WithData or WithError, that message is returned.
// Otherwise, the Error's own message is returned.
func (e Error) Error() string {
	if e.Data != nil && e.Data["error"] != nil {
		if value, ok := e.Data["error"].(string); ok {
			return value
		}
	}
	return e.Message
}

// Is determines whether the Error is an instance of the target.
// https://pkg.go.dev/errors#Is
//
// This function supports matching both any error and a specific error, based on the status code.
// Use rest.Err (an empty error) for the former:
//
//	errors.Is(err, rest.Err) // True if any REST API Error
//	errors.Is(err, rest.ErrBadRequest) // True if error is REST 400 Bad Request
func (e Error) Is(target error) bool {
	t, ok := target.(Error)
	if !ok {
		return false
	}
	if t.StatusCode == 0 {
		return true
	}
	return t.StatusCode == e.StatusCode
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

// WithError returns a copy of the HTTP error with the given error's message merged in to its additional data.
func (e Error) WithError(err error) Error {
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
func (e Error) Write(w http.ResponseWriter) {
	w.WriteHeader(e.StatusCode)
	w.Write([]byte(e.Message))
}

// WriteJSON writes the HTTP error to an HTTP response as JSON.
func (e Error) WriteJSON(w http.ResponseWriter) error {
	return WriteResponseJSON(w, e.StatusCode, e)
}
