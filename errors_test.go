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
	"net/url"
	"strconv"
	"strings"
	"testing"

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

func httpError(reason, field string) *http.Response {
	var reasons []APIErrorReason

	if reason != "" {
		reasons = append(reasons, APIErrorReason{
			Reason: reason,
			Field:  field,
		})
	}

	apiError := &APIError{
		Errors: reasons,
	}

	body, err := json.Marshal(apiError)
	if err != nil {
		panic("Failed to marshal APIError")
	}

	return &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(bytes.NewReader(body)),
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

	if newErr.Message != err.Error() || newErr.Code != ErrorFromError {
		t.Error("error should return ErrorFromError")
	}

	if newErr.Error() != "[002] test" {
		t.Error("Error should support Error() formatter with code")
	}

	if NewError(newErr) != newErr {
		t.Error("Error should be itself")
	}

	if err := NewError(&http.Response{Request: &http.Request{}}); err.Message != "Unexpected HTTP Error Response, no error" {
		t.Error("Unexpected HTTP Error Response, no error")
	}

	if err := NewError(httpError("testreason", "testfield")); err.Message != "[testfield] testreason" || err.Code != http.StatusInternalServerError {
		t.Errorf("unexpected HTTP response error: %v", err)
	}

	if err := NewError(httpError("testreason", "")); err.Message != "testreason" || err.Code != http.StatusInternalServerError {
		t.Errorf("unexpected fieldless HTTP response error: %v", err)
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

func createTestServer(t *testing.T, method, route, contentType, body string, statusCode int) (*httptest.Server, *Client) {
	t.Helper()

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

	client := newTestClient(t, nil)
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

		_, err := coupleAPIErrors(resp, nil)
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

		_, err := coupleAPIErrors(resp, nil)
		if err != nil {
			t.Error("http error with no reasons should return no error")
		}
	})

	t.Run("nil response", func(t *testing.T) {
		resp, err := coupleAPIErrors(nil, nil)
		if resp != nil || err == nil || err.Error() != "[002] response is nil" {
			t.Fatalf("got response %v, error %v; want nil response and response-is-nil error", resp, err)
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

		client := &Client{
			httpClient: ts.Client(),
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

		expectedError := &Error{
			Response: resp,
			Code:     http.StatusInternalServerError,
			Message:  "Unexpected Content-Type: Expected: application/json, Received: text/html\nResponse body: " + rawResponse,
		}

		defer resp.Body.Close()

		_, err = coupleAPIErrors(resp, nil)
		if err == nil {
			t.Fatal("expected an error for an HTML response")
		}

		if !strings.Contains(err.Error(), strconv.Itoa(expectedError.Code)) {
			t.Errorf("expected error code %d, got %d", expectedError.Code, resp.StatusCode)
		}

		if !strings.Contains(err.Error(), expectedError.Message) {
			t.Errorf("expected error message %s, got %s", expectedError.Message, err.Error())
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

		expectedError := &Error{
			Response: resp,
			Code:     http.StatusBadGateway,
			Message:  http.StatusText(http.StatusBadGateway),
		}

		_, err := coupleAPIErrors(resp, nil)
		if err == nil {
			t.Fatal("expected a bad gateway error")
		}
		if !strings.Contains(err.Error(), strconv.Itoa(expectedError.Code)) {
			t.Errorf("expected error code %d, got %d", expectedError.Code, resp.StatusCode)
		}

		if !strings.Contains(err.Error(), expectedError.Message) {
			t.Errorf("expected error message %s, got %s", expectedError.Message, err.Error())
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
			testName:       "default and Error from empty http error",
			err1:           NewError(httpError("", "")),
			err2:           defaultError,
			expectedResult: true,
		},
		{
			testName:       "default and Error from http error with field",
			err1:           NewError(httpError("", "test field")),
			err2:           defaultError,
			expectedResult: true,
		},
		{
			testName:       "default and Error from http error with field and reason",
			err1:           NewError(httpError("test reason", "test field")),
			err2:           defaultError,
			expectedResult: true,
		},
		{
			testName:       "default and Error from http error with reason",
			err1:           NewError(httpError("test reason", "")),
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

			if got := errors.Is(tc.err1, tc.err2); got != tc.expectedResult {
				t.Errorf("errors.Is(%v, %v) = %t, want %t", tc.err1, tc.err2, got, tc.expectedResult)
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
			if got := IsNotFound(err); got != tt.match {
				t.Errorf("IsNotFound(%v) = %t, want %t", err, got, tt.match)
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
		{
			name:  "WrappedError",
			err:   fmt.Errorf("request failed: %w", &Error{Code: http.StatusForbidden}),
			codes: []int{http.StatusForbidden},
			match: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ErrHasStatus(tt.err, tt.codes...)
			if got != tt.match {
				t.Errorf("ErrHasStatus(%v, %v) = %t, want %t", tt.err, tt.codes, got, tt.match)
			}
		})
	}
}

func TestClassifyInstanceCreateError(t *testing.T) {
	const (
		capacityMsg     = "The Linode plan you chose is not currently available in the selected region. Please select another region or plan type, or contact Support for assistance."
		accountLimitMsg = "You've reached a limit for the number of active services on your account. Please contact Support to request an increase and provide the total number of services you may need."
	)

	capacityErr := &Error{Code: http.StatusForbidden, Message: capacityMsg}
	accountLimitErr := &Error{Code: http.StatusBadRequest, Message: accountLimitMsg}

	type testCase struct {
		name string
		err  error
		want InstanceCreateErrorReason
	}

	tests := []testCase{
		{
			name: "InsufficientCapacity",
			err:  capacityErr,
			want: InstanceCreateErrorInsufficientCapacity,
		},
		{
			name: "InsufficientCapacityDifferentCase",
			err:  &Error{Code: http.StatusForbidden, Message: strings.ToUpper(capacityMsg)},
			want: InstanceCreateErrorInsufficientCapacity,
		},
		{
			name: "InsufficientCapacityWrapped",
			err:  fmt.Errorf("Error creating Linode instance: %w", capacityErr),
			want: InstanceCreateErrorInsufficientCapacity,
		},
		{
			name: "AccountLimitReached",
			err:  accountLimitErr,
			want: InstanceCreateErrorAccountLimitReached,
		},
		{
			name: "AccountLimitReachedWrapped",
			err:  fmt.Errorf("create instance: %w", accountLimitErr),

			want: InstanceCreateErrorAccountLimitReached,
		},
		{
			name: "AccountLimitReachedDifferentCase",
			err:  &Error{Code: http.StatusBadRequest, Message: strings.ToUpper(accountLimitMsg)},
			want: InstanceCreateErrorAccountLimitReached,
		},
		{
			name: "AccountLimitTextWithWrongStatus",
			err:  &Error{Code: http.StatusForbidden, Message: accountLimitMsg},
		},
		{
			name: "BothTextsWith403",
			err:  &Error{Code: http.StatusForbidden, Message: accountLimitMsg + "; " + capacityMsg},
			want: InstanceCreateErrorInsufficientCapacity,
		},
		{
			name: "BothTextsWith400",
			err:  &Error{Code: http.StatusBadRequest, Message: capacityMsg + "; " + accountLimitMsg},
			want: InstanceCreateErrorAccountLimitReached,
		},
		{
			name: "CapacityTextWithWrongStatus",
			err:  &Error{Code: http.StatusBadRequest, Message: capacityMsg},
		},
		{
			name: "LegacyCapacityWording",
			err:  &Error{Code: http.StatusForbidden, Message: "The selected region does not have any availability for the requested service plan"},
		},
		{
			name: "UnknownError",
			err:  &Error{Code: http.StatusBadRequest, Message: "The region this linode is in does not support disk encryption at this time"},
		},
		{
			name: "UnknownForbiddenError",
			err:  &Error{Code: http.StatusForbidden, Message: "You do not have permission to create Linodes."},
		},
		{
			name: "ClientSideError",
			err:  NewError(errors.New("context deadline exceeded")),
		},
		{
			name: "ContextDeadlineExceeded",
			err:  context.DeadlineExceeded,
		},
		{
			name: "WrappedURLError",
			err: fmt.Errorf("Error creating Linode instance: %w", &url.Error{
				Op:  http.MethodPost,
				URL: "https://api.linode.com/v4/linode/instances",
				Err: context.DeadlineExceeded,
			}),
		},
		{
			name: "NotALinodeError",
			err:  errors.New("[403] " + capacityMsg),
		},
		{
			name: "NilError",
		},
		{
			name: "TypedNilError",
			err:  (*Error)(nil),
		},
	}

	for _, code := range []int{
		http.StatusOK,
		http.StatusNotFound,
		http.StatusRequestTimeout,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusHTTPVersionNotSupported,
		http.StatusInsufficientStorage,
	} {
		tests = append(tests, testCase{
			name: "UnrecognizedStatus" + strconv.Itoa(code),
			err:  &Error{Code: code, Message: capacityMsg + "; " + accountLimitMsg},
		})
	}

	for _, code := range []int{ErrorUnsupported, ErrorFromString, ErrorFromError, ErrorFromStringer} {
		tests = append(tests, testCase{
			name: "PseudoCode" + strconv.Itoa(code),
			err:  &Error{Code: code, Message: capacityMsg + "; " + accountLimitMsg},
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClassifyInstanceCreateError(tt.err); got != tt.want {
				t.Errorf("ClassifyInstanceCreateError() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClassifyInstanceCreateErrorAPIResponse(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		want   InstanceCreateErrorReason
	}{
		{
			name:   "CapacityWithUnrelatedPermissionError",
			status: http.StatusForbidden,
			body:   `{"errors":[{"reason":"You do not have permission to perform this action."},{"reason":"The Linode plan you chose is not currently available in the selected region."}]}`,
			want:   InstanceCreateErrorInsufficientCapacity,
		},
		{
			name:   "AccountLimitWithUnrelatedLabelError",
			status: http.StatusBadRequest,
			body:   `{"errors":[{"reason":"You've reached a limit for the number of active services on your account."},{"field":"label","reason":"Label must be unique"}]}`,
			want:   InstanceCreateErrorAccountLimitReached,
		},
		{
			name:   "LabelConflictIsNotClassified",
			status: http.StatusBadRequest,
			body:   `{"errors":[{"field":"label","reason":"Label must be unique"}]}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := &http.Response{
				StatusCode: tc.status,
				Body:       io.NopCloser(strings.NewReader(tc.body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Request:    &http.Request{Header: http.Header{"Accept": []string{"application/json"}}},
			}
			defer response.Body.Close()

			_, err := coupleAPIErrors(response, nil)
			if err == nil {
				t.Fatal("expected an API error")
			}
			if got := ClassifyInstanceCreateError(err); got != tc.want {
				t.Errorf("ClassifyInstanceCreateError(%v) = %q, want %q", err, got, tc.want)
			}
		})
	}
}
