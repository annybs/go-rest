package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// ReadRequestJSON reads the body of an HTTP request into a target reference.
func ReadRequestJSON(req *http.Request, v any) error {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// WriteError writes an error to an HTTP response.
// If the error is a REST API error, it is written as standard per Error.Write.
// Otherwise, ErrInternalServerError is written with the given error attached as data.
func WriteError(w http.ResponseWriter, err error) error {
	if errors.Is(err, Err) {
		_, ret := err.(Error).Write(w)
		return ret
	}
	_, ret := ErrInternalServerError.WithError(err).Write(w)
	return ret
}

// WriteErrorJSON writes an error as JSON to an HTTP response.
// If the error is a REST API error, it is written as standard per Error.WriteJSON.
// Otherwise, ErrInternalServerError is written with the given error attached as data.
func WriteErrorJSON(w http.ResponseWriter, err error) error {
	if errors.Is(err, Err) {
		return err.(Error).WriteJSON(w)
	}
	return ErrInternalServerError.WithError(err).WriteJSON(w)
}

// WriteResponseJSON writes an HTTP response as JSON.
func WriteResponseJSON(w http.ResponseWriter, statusCode int, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(b)
	return nil
}
