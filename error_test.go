package rest

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"
)

func TestErrorWriteJSON(t *testing.T) {
	type TestCase struct {
		Input  Error
		C      int
		Output string
		Err    error
	}

	testCases := []TestCase{
		// Empty error
		{Input: Err, C: 200, Output: `{"message":""}`},

		// Standard errors
		{Input: ErrPermanentRedirect, C: 308, Output: `{"message":"Permanent Redirect"}`},
		{Input: ErrNotFound, C: 404, Output: `{"message":"Not Found"}`},
		{Input: ErrInternalServerError, C: 500, Output: `{"message":"Internal Server Error"}`},

		// Error with changed message
		{Input: ErrBadRequest.WithMessage("Invalid Recipe"), C: 400, Output: `{"message":"Invalid Recipe"}`},

		// Error with data
		{
			Input:  ErrGatewayTimeout.WithData(map[string]any{"service": "RecipeDatabase"}),
			C:      504,
			Output: `{"message":"Gateway Timeout","data":{"service":"RecipeDatabase"}}`,
		},

		// Error with value
		{
			Input:  ErrGatewayTimeout.WithValue("service", "RecipeDatabase"),
			C:      504,
			Output: `{"message":"Gateway Timeout","data":{"service":"RecipeDatabase"}}`,
		},

		// Error with error
		{
			Input:  ErrInternalServerError.WithError(errors.New("recipe is too delicious")),
			C:      500,
			Output: `{"message":"Internal Server Error","data":{"error":"recipe is too delicious"}}`,
		},
	}

	for i, tc := range testCases {
		t.Logf("(%d) Testing %v", i, tc.Input)

		rec := httptest.NewRecorder()

		err := tc.Input.WriteJSON(rec)
		if err != tc.Err {
			t.Errorf("Expected error %v, got %v", tc.Err, err)
		}
		if err != nil {
			continue
		}

		res := rec.Result()
		if res.StatusCode != tc.C {
			t.Errorf("Expected status code %d, got %d", tc.C, res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Unexpected error reading response body: %v", err)
			continue
		}

		if string(body) != tc.Output {
			t.Errorf("Expected body %q, got %q", tc.Output, string(body))
		}
	}
}
