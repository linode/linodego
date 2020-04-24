package linodego

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

func TestClientClone(t *testing.T) {
	bogusRetryConditional := func(_ *resty.Response, _ error) bool { return true }

	original := NewClient(nil)
	original.addRetryConditional(bogusRetryConditional)
	originalUserAgent := original.resty.Header.Get("User-Agent")
	originalBaseURL := original.resty.HostURL
	originalRetryMaxWaitTime := original.resty.RetryMaxWaitTime
	originalRetryConditionals := len(original.retryConditionals)

	clone := original.clone()
	if reflect.DeepEqual(original, clone) {
		t.Errorf("expected original %#v and cloned %#v clients to be equal", original, clone)
	}

	clone.SetDebug(true)
	if original.debug != false {
		t.Error("expected original client to have debug disabled")
	}

	clone.addRetryConditional(bogusRetryConditional)
	if len(clone.retryConditionals) != originalRetryConditionals+1 {
		t.Error("expected retry conditional to have been added to cloned client")
	}
	if len(original.retryConditionals) != originalRetryConditionals {
		t.Errorf("expected original client to have %d retryConditionals; got %d", originalRetryConditionals, len(original.retryConditionals))
	}

	newUserAgent := "test"
	clone.SetUserAgent(newUserAgent)
	if clone.resty.Header.Get("User-Agent") != newUserAgent {
		t.Errorf("expected cloned client to have user agent '%s'; got '%s'", newUserAgent, clone.resty.Header.Get("User-Agent"))
	}
	if original.resty.Header.Get("User-Agent") != originalUserAgent {
		t.Errorf("expected original client to have user agent '%s'; got '%s'", originalUserAgent, original.resty.Header.Get("User-Agent"))
	}

	newBaseURL := "http://0.0.0.0/api/v1beta"
	clone.SetBaseURL(newBaseURL)
	if clone.resty.HostURL != newBaseURL {
		t.Errorf("expected cloned client to have base url '%s'; got '%s'", newUserAgent, clone.resty.HostURL)
	}
	if original.resty.HostURL != originalBaseURL {
		t.Errorf("expected original client to have base url '%s'; got '%s'", originalBaseURL, original.resty.HostURL)
	}

	newRetryMaxWaitTime := time.Minute * 3
	clone.SetRetryMaxWaitTime(newRetryMaxWaitTime)
	if clone.resty.RetryMaxWaitTime != newRetryMaxWaitTime {
		t.Errorf("expected cloned client to have retry max wait time of %d; got %d", newRetryMaxWaitTime, clone.resty.RetryMaxWaitTime)
	}
	if original.resty.RetryMaxWaitTime != originalRetryMaxWaitTime {
		t.Errorf("expected original client to have retry max wait time of %d; got %d", originalRetryMaxWaitTime, original.resty.RetryMaxWaitTime)
	}
}
