package linodego

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// type RetryConditional func(r *resty.Response) (shouldRetry bool)
type RetryConditional resty.RetryConditionFunc

var retryConditionals []RetryConditional

// Configures resty to
// lock until enough time has passed to retry the request as determined by the Retry-After response header.
// If the Retry-After header is not set, we fall back to value of SetPollDelay.
func configureRestyRetries(resty *resty.Client) {
	resty.
		SetRetryCount(1000).
		SetRetryMaxWaitTime(30 * time.Second).
		AddRetryCondition(checkRetryConditionals).
		SetRetryAfter(retryAfter)
}

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

func retryAfter(client *resty.Client, resp *resty.Response) (time.Duration, error) {
	retryAfterStr := resp.Header().Get("Retry-After")
	if retryAfterStr == "" {
		return 0, nil
	}

	retryAfter, err := strconv.Atoi(retryAfterStr)
	if err != nil {
		return 0, err
	}

	return time.Duration(retryAfter) * time.Second, nil
}
