package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsAuthenticated(t *testing.T) {
	type TestCase struct {
		Req      *http.Request
		Token    string
		Expected bool
	}

	testCases := []TestCase{}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("authorization", "bearer abcd")
	testCases = append(testCases, TestCase{Req: req, Token: "abcd", Expected: true})

	req = httptest.NewRequest("POST", "/", nil)
	req.Header.Add("authorization", "Bearer defg hijk")
	testCases = append(testCases, TestCase{Req: req, Token: "defg hijk", Expected: true})

	req = httptest.NewRequest("DELETE", "/", nil)
	testCases = append(testCases, TestCase{Req: req, Token: "", Expected: true})

	req = httptest.NewRequest("GET", "/", nil)
	testCases = append(testCases, TestCase{Req: req, Token: "lmno"})

	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("authorization", "Bearer pqrs")
	testCases = append(testCases, TestCase{Req: req, Expected: true})

	for i, tc := range testCases {
		t.Logf("(%d) Testing request authorization header against %q", i, tc.Token)

		actual := IsAuthenticated(tc.Req, tc.Token)
		if actual != tc.Expected {
			t.Errorf("Expected %v, got %v", tc.Expected, actual)
		}
	}
}

func TestReadBearerToken(t *testing.T) {
	type TestCase struct {
		Req      *http.Request
		Expected string
	}

	testCases := []TestCase{}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("authorization", "bearer abcd")
	testCases = append(testCases, TestCase{Req: req, Expected: "abcd"})

	req = httptest.NewRequest("POST", "/", nil)
	req.Header.Add("authorization", "Bearer defg hijk")
	testCases = append(testCases, TestCase{Req: req, Expected: "defg hijk"})

	req = httptest.NewRequest("DELETE", "/", nil)
	testCases = append(testCases, TestCase{Req: req, Expected: ""})

	for i, tc := range testCases {
		t.Logf("(%d) Testing request authorization header against %q", i, tc.Expected)

		actual := ReadBearerToken(tc.Req)
		if actual != tc.Expected {
			t.Errorf("Expected %q, got %q", tc.Expected, actual)
		}
	}
}
