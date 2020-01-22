package linodego

import (
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// type RetryConditional func(r *resty.Response) (shouldRetry bool)
type RetryConditional resty.RetryConditionFunc

var retryConditionals []RetryConditional

func addRetryConditional(retryConditional RetryConditional) {
	retryConditionals = append(retryConditionals, retryConditional)
}

func checkRetryConditionals(r *resty.Response, err error) bool {
	for _, retryConditional := range retryConditionals {
		retry := retryConditional(r, err)
		if retry {
			log.Printf("[INFO] Received error %s - Retrying", r.Error())
			return true
		}
	}
	return false
}

// SetLinodeBusyRetry configures resty to retry specifically on "Linode busy." errors
// The retry wait time is configured in SetPollDelay
func linodeBusyRetryCondition(r *resty.Response, _ error) bool {
	apiError, ok := r.Error().(*APIError)
	linodeBusy := ok && apiError.Error() == "Linode busy."
	retry := r.StatusCode() == http.StatusBadRequest && linodeBusy
	return retry
}

func tooManyRequestsRetryCondition(r *resty.Response, _ error) bool {
	return r.StatusCode() == http.StatusTooManyRequests
}
