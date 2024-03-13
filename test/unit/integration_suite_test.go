package unit

import (
	"github.com/linode/linodego"
	"github.com/linode/linodego/internal/testutil"
	"testing"
)

func createMockClient(t *testing.T) *linodego.Client {
	return testutil.CreateMockClient(t, linodego.NewClient)
}
