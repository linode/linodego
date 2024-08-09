package linodego

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestHTTPLinodeBusyRetryCondition(t *testing.T) {
	var retry bool

	// Initialize response body
	rawResponse := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewBuffer(nil)),
	}

	retry = httpLinodeBusyRetryCondition(rawResponse, nil)

	if retry {
		t.Errorf("Should not have retried")
	}

	apiError := APIError{
		Errors: []APIErrorReason{
			{Reason: "Linode busy."},
		},
	}
	rawResponse.Body = createResponseBody(apiError)

	retry = httpLinodeBusyRetryCondition(rawResponse, nil)

	if !retry {
		t.Errorf("Should have retried")
	}
}

func TestHTTPServiceUnavailableRetryCondition(t *testing.T) {
	rawResponse := &http.Response{
		StatusCode: http.StatusServiceUnavailable,
		Header:     http.Header{httpRetryAfterHeaderName: []string{"20"}},
		Body:       io.NopCloser(bytes.NewBuffer(nil)), // Initialize response body
	}

	if retry := httpServiceUnavailableRetryCondition(rawResponse, nil); !retry {
		t.Error("expected request to be retried")
	}

	if retryAfter, err := httpRespectRetryAfter(rawResponse); err != nil {
		t.Errorf("expected error to be nil but got %s", err)
	} else if retryAfter != time.Second*20 {
		t.Errorf("expected retryAfter to be 20 but got %d", retryAfter)
	}
}

func TestHTTPServiceMaintenanceModeRetryCondition(t *testing.T) {
	rawResponse := &http.Response{
		StatusCode: http.StatusServiceUnavailable,
		Header: http.Header{
			httpRetryAfterHeaderName:      []string{"20"},
			httpMaintenanceModeHeaderName: []string{"Currently in maintenance mode."},
		},
		Body: io.NopCloser(bytes.NewBuffer(nil)), // Initialize response body
	}

	if retry := httpServiceUnavailableRetryCondition(rawResponse, nil); retry {
		t.Error("expected retry to be skipped due to maintenance mode header")
	}
}

// Helper function to create a response body from an object
func createResponseBody(obj interface{}) io.ReadCloser {
	body, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return io.NopCloser(bytes.NewBuffer(body))
}
