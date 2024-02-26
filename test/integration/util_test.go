package integration

import (
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/linode/linodego/internal/testutil"

	"github.com/jarcoal/httpmock"
)

func assertDateSet(t *testing.T, compared *time.Time) {
	emptyTime := time.Time{}
	if compared == nil || *compared == emptyTime {
		t.Errorf("Expected date to be set, got %v", compared)
	}
}

func mockRequestBodyValidate(t *testing.T, expected interface{}, response interface{}) httpmock.Responder {
	return testutil.MockRequestBodyValidate(t, expected, response)
}

func mockRequestURL(t *testing.T, path string) *regexp.Regexp {
	return testutil.MockRequestURL(path)
}

func assertSliceContains[T comparable](t *testing.T, slice []T, target T) {
	for _, v := range slice {
		if v == target {
			return
		}
	}

	t.Fatalf("value %v not found in slice", target)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// return the current nanosecond in string type as a unique text.
func getUniqueText() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
