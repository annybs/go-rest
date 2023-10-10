package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadRequestJSON(t *testing.T) {
	type TestInput struct {
		Title  string `json:"title"`
		Serves int    `json:"serves"`
		Author string `json:"author,omitempty"`
	}

	type TestCase struct {
		Req   *http.Request
		Input TestInput
		Err   error
	}

	testCases := []TestCase{
		{
			Req: httptest.NewRequest("POST", "/recipes", nil),
			Err: errors.New("unexpected end of JSON input"),
		},
		{
			Req:   httptest.NewRequest("POST", "/recipes", strings.NewReader(`{"title":"Gnocchi","serves":2}`)),
			Input: TestInput{Title: "Gnocchi", Serves: 2},
		},
		{
			Req:   httptest.NewRequest("POST", "/recipes", strings.NewReader(`{"title":"Spaghetti","serves":4,"author":"Mom"}`)),
			Input: TestInput{Title: "Spaghetti", Serves: 4, Author: "Mom"},
		},
	}

	for i, tc := range testCases {
		t.Logf("(%d) Testing request body against %+v", i, tc.Input)

		input := TestInput{}
		err := ReadRequestJSON(tc.Req, &input)

		if err != nil {
			if tc.Err != nil {
				// Compare error strings, as json.SyntaxError isn't directly comparable
				if err.Error() != tc.Err.Error() {
					t.Errorf("Expected error %v, got %v", tc.Err, err)
				}
			} else {
				t.Errorf("Expected error %v, got %v", tc.Err, err)
			}
			continue
		}

		if input != tc.Input {
			t.Errorf("Expected %v, got %v", tc.Input, input)
		}
	}
}
