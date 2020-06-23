package retry

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/linode/linodego/pkg/errors"
)

func TestLinodeBusyRetryCondition(t *testing.T) {
	var retry bool

	request := resty.Request{}
	rawResponse := http.Response{StatusCode: http.StatusBadRequest}
	response := resty.Response{
		Request:     &request,
		RawResponse: &rawResponse,
	}

	retry = LinodeBusyRetryCondition(&response, nil)

	if retry {
		t.Errorf("Should not have retried")
	}

	apiError := errors.APIError{
		Errors: []errors.APIErrorReason{
			{Reason: "Linode busy."},
		},
	}
	request.SetError(&apiError)

	retry = LinodeBusyRetryCondition(&response, nil)

	if !retry {
		t.Errorf("Should have retried")
	}
}

func TestLinodeServiceUnavailableRetryCondition(t *testing.T) {
	request := resty.Request{}
	rawResponse := http.Response{StatusCode: http.StatusServiceUnavailable, Header: http.Header{
		retryAfterHeaderName: []string{"20"},
	}}
	response := resty.Response{
		Request:     &request,
		RawResponse: &rawResponse,
	}

	if retry := ServiceUnavailableRetryCondition(&response, nil); !retry {
		t.Error("expected request to be retried")
	}

	if retryAfter, err := RespectRetryAfter(resty.New(), &response); err != nil {
		t.Errorf("expected error to be nil but got %s", err)
	} else if retryAfter != time.Second*20 {
		t.Errorf("expected retryAfter to be 20 but got %d", retryAfter)
	}
}
