package unit

import (
	"testing"

	"github.com/linode/linodego"
	"github.com/linode/linodego/internal/testutil"
)

func createMockClient(t *testing.T) *linodego.Client {
	return testutil.CreateMockClient(t, linodego.NewClient)
}
