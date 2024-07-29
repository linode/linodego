package linodego

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-cmp/cmp"
)

type testStringer string

func (t testStringer) String() string {
	return string(t)
}

type testError string

func (e testError) Error() string {
	return string(e)
}

func restyError(reason, field string) *resty.Response {
	var reasons []APIErrorReason

	// allow for an empty reasons
	if reason != "" && field != "" {
		reasons = append(reasons, APIErrorReason{
			Reason: reason,
			Field:  field,
		})
	}

	return &resty.Response{
		RawResponse: &http.Response{
			StatusCode: 500,
		},
		Request: &resty.Request{
			Error: &APIError{
				Errors: reasons,
			},
		},
	}
}

func TestNewError(t *testing.T) {
	if NewError(nil) != nil {
		t.Errorf("nil error should return nil")
	}
	if NewError(struct{}{}).Code != ErrorUnsupported {
		t.Error("empty struct should return unsupported error type")
	}

	err := errors.New("test")
	newErr := NewError(err)

	if newErr.Message != err.Error() && newErr.Code != ErrorFromError {
		t.Error("error should return ErrorFromError")
	}

	if newErr.Error() != "[002] test" {
		t.Error("Error should support Error() formatter with code")
	}

	if NewError(newErr) != newErr {
		t.Error("Error should be itself")
	}

	if err := NewError(&resty.Response{Request: &resty.Request{}}); err.Message != "Unexpected Resty Error Response, no error" {
		t.Error("Unexpected Resty Error Response, no error")
	}

	if err := NewError(restyError("testreason", "testfield")); err.Message != "[testfield] testreason" {
		t.Error("rest response error should should be set")
	}

	if err := NewError("stringerror"); err.Message != "stringerror" || err.Code != ErrorFromString {
		t.Errorf("string error should be set")
	}

	if err := NewError(testStringer("teststringer")); err.Message != "teststringer" || err.Code != ErrorFromStringer {
		t.Errorf("error should be set for a stringer interface")
	}

	if err := NewError(testError("testerror")); err.Message != "testerror" || err.Code != ErrorFromError {
		t.Errorf("error should be set for an error interface")
	}
}

func createTestServer(method, route, contentType, body string, statusCode int) (*httptest.Server, *Client) {
	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == method && r.URL.Path == route {
			rw.Header().Add("Content-Type", contentType)
			rw.WriteHeader(statusCode)
			rw.Write([]byte(body))
			return
		}
		rw.WriteHeader(http.StatusNotImplemented)
	})
	ts := httptest.NewServer(h)

	client := NewClient(nil)
	client.SetBaseURL(ts.URL)
	return ts, &client
}

func TestCoupleAPIErrors(t *testing.T) {
	t.Run("not nil error generates error", func(t *testing.T) {
		err := errors.New("test")
		if _, err := coupleAPIErrors(nil, err); !cmp.Equal(err, NewError(err)) {
			t.Errorf("expect a not nil error to be returned as an Error")
		}
	})

	t.Run("resty 500 response error with reasons", func(t *testing.T) {
		if _, err := coupleAPIErrors(restyError("testreason", "testfield"), nil); err.Error() != "[500] [testfield] testreason" {
			t.Error("resty error should return with proper format [code] [field] reason")
		}
	})

	t.Run("resty 500 response error without reasons", func(t *testing.T) {
		if _, err := coupleAPIErrors(restyError("", ""), nil); err != nil {
			t.Error("resty error with no reasons should return no error")
		}
	})

	t.Run("resty response with nil error", func(t *testing.T) {
		emptyErr := &resty.Response{
			RawResponse: &http.Response{
				StatusCode: 500,
			},
			Request: &resty.Request{
				Error: nil,
			},
		}
		if _, err := coupleAPIErrors(emptyErr, nil); err != nil {
			t.Error("resty error with no reasons should return no error")
		}
	})

	t.Run("generic html error", func(t *testing.T) {
		rawResponse := `<html>
<head><title>500 Internal Server Error</title></head>
<body bgcolor="white">
<center><h1>500 Internal Server Error</h1></center>
<hr><center>nginx</center>
</body>
</html>`
		route := "/v4/linode/instances/123"
		ts, client := createTestServer(http.MethodGet, route, "text/html", rawResponse, http.StatusInternalServerError)
		// client.SetDebug(true)
		defer ts.Close()

		expectedError := Error{
			Code:    http.StatusInternalServerError,
			Message: "Unexpected Content-Type: Expected: application/json, Received: text/html\nResponse body: " + rawResponse,
		}

		_, err := coupleAPIErrors(client.R(context.Background()).SetResult(&Instance{}).Get(ts.URL + route))
		if diff := cmp.Diff(expectedError, err); diff != "" {
			t.Errorf("expected error to match but got diff:\n%s", diff)
		}
	})

	t.Run("bad gateway error", func(t *testing.T) {
		rawResponse := []byte(`<html>
<head><title>502 Bad Gateway</title></head>
<body bgcolor="white">
<center><h1>502 Bad Gateway</h1></center>
<hr><center>nginx</center>
</body>
</html>`)
		buf := io.NopCloser(bytes.NewBuffer(rawResponse))

		resp := &resty.Response{
			Request: &resty.Request{
				Error: errors.New("Bad Gateway"),
			},
			RawResponse: &http.Response{
				Header: http.Header{
					"Content-Type": []string{"text/html"},
				},
				StatusCode: http.StatusBadGateway,
				Body:       buf,
			},
		}

		expectedError := Error{
			Code:    http.StatusBadGateway,
			Message: http.StatusText(http.StatusBadGateway),
		}

		if _, err := coupleAPIErrors(resp, nil); !cmp.Equal(err, expectedError) {
			t.Errorf("expected error %#v to match error %#v", err, expectedError)
		}
	})
}

func TestCoupleAPIErrorsHTTP(t *testing.T) {
	t.Run("not nil error generates error", func(t *testing.T) {
		err := errors.New("test")
		if _, err := coupleAPIErrorsHTTP(nil, err); !cmp.Equal(err, NewError(err)) {
			t.Errorf("expect a not nil error to be returned as an Error")
		}
	})

	t.Run("http 500 response error with reasons", func(t *testing.T) {
		// Create the simulated HTTP response with a 500 status and a JSON body containing the error details
		apiError := APIError{
			Errors: []APIErrorReason{
				{Reason: "testreason", Field: "testfield"},
			},
		}
		apiErrorBody, _ := json.Marshal(apiError)
		bodyReader := io.NopCloser(bytes.NewBuffer(apiErrorBody))

		resp := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       bodyReader,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Request:    &http.Request{Header: http.Header{"Accept": []string{"application/json"}}},
		}

		_, err := coupleAPIErrorsHTTP(resp, nil)
		expectedMessage := "[500] [testfield] testreason"
		if err == nil || err.Error() != expectedMessage {
			t.Errorf("expected error message %q, got: %v", expectedMessage, err)
		}
	})

	t.Run("http 500 response error without reasons", func(t *testing.T) {
		// Create the simulated HTTP response with a 500 status and an empty errors array
		apiError := APIError{
			Errors: []APIErrorReason{},
		}
		apiErrorBody, _ := json.Marshal(apiError)
		bodyReader := io.NopCloser(bytes.NewBuffer(apiErrorBody))

		resp := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       bodyReader,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Request:    &http.Request{Header: http.Header{"Accept": []string{"application/json"}}},
		}

		_, err := coupleAPIErrorsHTTP(resp, nil)
		if err != nil {
			t.Error("http error with no reasons should return no error")
		}
	})

	t.Run("http response with nil error", func(t *testing.T) {
		// Create the simulated HTTP response with a 500 status and a nil error
		resp := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewBuffer([]byte(`{"errors":[]}`))), // empty errors array in body
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Request:    &http.Request{Header: http.Header{"Accept": []string{"application/json"}}},
		}

		_, err := coupleAPIErrorsHTTP(resp, nil)
		if err != nil {
			t.Error("http error with no reasons should return no error")
		}
	})

	t.Run("generic html error", func(t *testing.T) {
		rawResponse := `<html>
<head><title>500 Internal Server Error</title></head>
<body bgcolor="white">
<center><h1>500 Internal Server Error</h1></center>
<hr><center>nginx</center>
</body>
</html>`

		route := "/v4/linode/instances/123"
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(rawResponse))
		}))
		defer ts.Close()

		client := &httpClient{
			httpClient: ts.Client(),
		}

		expectedError := Error{
			Code:    http.StatusInternalServerError,
			Message: "Unexpected Content-Type: Expected: application/json, Received: text/html\nResponse body: " + rawResponse,
		}

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, ts.URL+route, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Accept", "application/json")

		resp, err := client.httpClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		_, err = coupleAPIErrorsHTTP(resp, nil)
		if diff := cmp.Diff(expectedError, err); diff != "" {
			t.Errorf("expected error to match but got diff:\n%s", diff)
		}
	})

	t.Run("bad gateway error", func(t *testing.T) {
		rawResponse := `<html>
<head><title>502 Bad Gateway</title></head>
<body bgcolor="white">
<center><h1>502 Bad Gateway</h1></center>
<hr><center>nginx</center>
</body>
</html>`
		buf := io.NopCloser(bytes.NewBuffer([]byte(rawResponse)))

		resp := &http.Response{
			StatusCode: http.StatusBadGateway,
			Body:       buf,
			Header: http.Header{
				"Content-Type": []string{"text/html"},
			},
			Request: &http.Request{
				Header: http.Header{"Accept": []string{"application/json"}},
			},
		}

		expectedError := Error{
			Code:    http.StatusBadGateway,
			Message: http.StatusText(http.StatusBadGateway),
		}

		_, err := coupleAPIErrorsHTTP(resp, nil)
		if !cmp.Equal(err, expectedError) {
			t.Errorf("expected error %#v to match error %#v", err, expectedError)
		}
	})
}

func TestErrorIs(t *testing.T) {
	t.Parallel()

	defaultError := &Error{
		Message: "default error",
		Code:    http.StatusInternalServerError,
	}

	for _, tc := range []struct {
		testName       string
		err1           error
		err2           error
		expectedResult bool
	}{
		{
			testName:       "base errors.Is comparision",
			err1:           defaultError,
			err2:           defaultError,
			expectedResult: true,
		},
		{
			testName:       "wrapped default",
			err1:           fmt.Errorf("test wrap: %w", defaultError),
			err2:           defaultError,
			expectedResult: true,
		},
		{
			testName:       "deeply wrapped error",
			err1:           fmt.Errorf("wrap 1: %w", fmt.Errorf("wrap 2: %w", defaultError)),
			err2:           defaultError,
			expectedResult: true,
		},
		{
			testName:       "default and Error from empty resty error",
			err1:           NewError(restyError("", "")),
			err2:           defaultError,
			expectedResult: true,
		},
		{
			testName:       "default and Error from resty error with field",
			err1:           NewError(restyError("", "test field")),
			err2:           defaultError,
			expectedResult: true,
		},
		{
			testName:       "default and Error from resty error with field and reason",
			err1:           NewError(restyError("test reason", "test field")),
			err2:           defaultError,
			expectedResult: true,
		},
		{
			testName:       "default and Error from resty error with reason",
			err1:           NewError(restyError("test reason", "")),
			err2:           defaultError,
			expectedResult: true,
		},
		{
			testName:       "error and nil",
			err1:           defaultError,
			err2:           nil,
			expectedResult: false,
		},
		{
			testName:       "wrapped nil",
			err1:           fmt.Errorf("test wrap: %w", nil),
			err2:           defaultError,
			expectedResult: false,
		},
		{
			testName:       "both errors are different nil", // NOTE: nils of different types are never equal
			err1:           nil,
			err2:           (*Error)(nil),
			expectedResult: false,
		},
		{
			testName:       "different error types",
			err1:           errors.New("different error type"),
			err2:           defaultError,
			expectedResult: false,
		},
	} {
		tc := tc

		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			if errors.Is(tc.err1, tc.err2) != tc.expectedResult {
				t.Errorf("expected %+#v to be equal %+#v", tc.err1, tc.err2)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		code  int
		match bool
	}{
		{code: http.StatusNotFound, match: true},
		{code: http.StatusInternalServerError},
		{code: http.StatusFound},
		{code: http.StatusOK},
	}

	for _, tt := range tests {
		name := http.StatusText(tt.code)
		t.Run(name, func(t *testing.T) {
			err := &Error{Code: tt.code}
			if matches := IsNotFound(err); !matches && tt.match {
				t.Errorf("should have matched %d", tt.code)
			} else if matches && !tt.match {
				t.Errorf("shoudl not have matched %d", tt.code)
			}
		})
	}
}

func TestErrHasStatusCode(t *testing.T) {
	tests := []struct {
		name  string
		err   error
		codes []int
		match bool
	}{
		{
			name:  "NotFound",
			err:   &Error{Code: http.StatusNotFound},
			codes: []int{http.StatusNotFound},
			match: true,
		},
		{
			name: "NoCodes",
			err:  &Error{Code: http.StatusInternalServerError},
		},
		{
			name:  "MultipleCodes",
			err:   &Error{Code: http.StatusTeapot},
			codes: []int{http.StatusBadRequest, http.StatusTeapot, http.StatusUnavailableForLegalReasons},
			match: true,
		},
		{
			name:  "NotALinodeError",
			err:   io.EOF,
			codes: []int{http.StatusTeapot},
		},
		{
			name:  "NoMatch",
			err:   &Error{Code: http.StatusTooEarly},
			codes: []int{http.StatusLocked, http.StatusTooManyRequests},
		},
		{
			name:  "NilError",
			codes: []int{http.StatusGone},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ErrHasStatus(tt.err, tt.codes...)
			if !got && tt.match {
				t.Errorf("should have matched")
			} else if got && !tt.match {
				t.Errorf("should not have matched")
			}
		})
	}
}
