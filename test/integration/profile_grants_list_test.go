package integration

import (
	"context"
	"testing"
	//"github.com/google/go-cmp/cmp"
	//"github.com/linode/linodego"
)

func TestGrantsList(t *testing.T) {
	//username := usernamePrefix + "grantslist"

	client, teardown := createTestClient(t, "fixtures/TestGrantsList")
	client.GrantsList(context.Background())
	defer teardown()
}
