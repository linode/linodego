package linodego

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-cmp/cmp"
)

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

func TestCoupleAPIErrors_genericHtmlError(t *testing.T) {
	rawResponse := `<html>
<head><title>500 Internal Server Error</title></head>
<body bgcolor="white">
<center><h1>500 Internal Server Error</h1></center>
<hr><center>nginx</center>
</body>
</html>`
	route := "/v4/linode/instances/123"
	ts, client := createTestServer(http.MethodGet, route, "text/html", rawResponse, http.StatusInternalServerError)
	client.SetDebug(true)
	defer ts.Close()

	expectedError := Error{
		Code:    http.StatusInternalServerError,
		Message: "Unexpected Content-Type: Expected: application/json, Received: text/html\nResponse body: " + rawResponse,
	}

	_, err := coupleAPIErrors(client.R(context.Background()).SetResult(&Instance{}).Get(ts.URL + route))
	if diff := cmp.Diff(expectedError, err); diff != "" {
		t.Errorf("expected error to match but got diff:\n%s", diff)
	}
}

func TestCoupleAPIErrors_badGatewayError(t *testing.T) {
	rawResponse := []byte(`<html>
<head><title>502 Bad Gateway</title></head>
<body bgcolor="white">
<center><h1>502 Bad Gateway</h1></center>
<hr><center>nginx</center>
</body>
</html>`)
	buf := ioutil.NopCloser(bytes.NewBuffer(rawResponse))

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
}
