package rest

import (
	"encoding/json"
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
