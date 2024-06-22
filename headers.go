package rest

import (
	"net/http"
)

// IsAuthenticated returns true if the bearer token in a request's authorization is equal to a user-defined token.
// This function always returns true if the user-defined token is empty i.e. no authentication required.
func IsAuthenticated(req *http.Request, token string) bool {
	if token == "" {
		return true
	}

	read := ReadBearerToken(req)
	return read == token
}

// ReadBearerToken reads the token portion of a bearer token in a request's authorization header.
// This function returns an empty string if the header is not provided or is not a bearer token.
func ReadBearerToken(req *http.Request) string {
	header := req.Header.Get("authorization")
	if len(header) > 8 {
		bearer := header[0:7]
		if bearer == "bearer " || bearer == "Bearer " {
			return header[7:]
		}
	}
	return ""
}
