package errors

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-cmp/cmp"
)

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

	if _, err := CoupleAPIErrors(resp, nil); !cmp.Equal(err, expectedError) {
		t.Errorf("expected error %#v to match error %#v", err, expectedError)
	}
}
