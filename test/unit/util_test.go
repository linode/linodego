package unit

import (
	"regexp"
	"testing"

	"github.com/linode/linodego"
	"github.com/linode/linodego/internal/testutil"

	"github.com/jarcoal/httpmock"
)

func mockRequestBodyValidate(t *testing.T, expected interface{}, response interface{}) httpmock.Responder {
	return testutil.MockRequestBodyValidate(t, expected, response)
}

func mockRequestURL(t *testing.T, path string) *regexp.Regexp {
	return testutil.MockRequestURL(path)
}

func createMockClient(t *testing.T) *linodego.Client {
	return testutil.CreateMockClient(t, linodego.NewClient)
}
