package linodego

import (
	"bytes"
	"context"
	"errors"
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
