package unit

import (
	"github.com/linode/linodego/internal/testutil"
	"regexp"
	"testing"

	"github.com/jarcoal/httpmock"
)

func mockRequestBodyValidate(t *testing.T, expected interface{}, response interface{}) httpmock.Responder {
	return testutil.MockRequestBodyValidate(t, expected, response)
}

func mockRequestURL(t *testing.T, path string) *regexp.Regexp {
	return testutil.MockRequestURL(path)
}
