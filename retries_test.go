package linodego

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestLinodeBusyRetryCondition(t *testing.T) {
	var retry bool

	request := resty.Request{}
	rawResponse := http.Response{StatusCode: http.StatusBadRequest}
	response := resty.Response{
		Request:     &request,
		RawResponse: &rawResponse,
	}

	retry = linodeBusyRetryCondition(&response, nil)

	if retry {
		t.Errorf("Should not have retried")
	}

	apiError := APIError{
		Errors: []APIErrorReason{
			APIErrorReason{Reason: "Linode busy."},
		},
	}
	request.SetError(&apiError)

	retry = linodeBusyRetryCondition(&response, nil)

	if !retry {
		t.Errorf("Should have retried")
	}
}
