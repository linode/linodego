package unit

import (
	"fmt"
	"net/url"
	"regexp"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/linode/linodego/internal/testutil"
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

func formatMockAPIPath(format string, args ...any) string {
	escapedArgs := make([]any, len(args))
	for i, arg := range args {
		if typeStr, ok := arg.(string); ok {
			arg = url.PathEscape(typeStr)
		}

		escapedArgs[i] = arg
	}

	return fmt.Sprintf(format, escapedArgs...)
}
