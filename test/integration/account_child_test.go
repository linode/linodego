package integration

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestAccountChild_basic(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestAccountChild_basic")
	defer teardown()

	childAccounts, err := client.ListChildAccounts(context.Background(), nil)
	require.NoError(t, err)
	require.Greater(
		t,
		len(childAccounts),
		0,
		"number of child accounts should be > 0",
	)

	childAccount, err := client.GetChildAccount(context.Background(), childAccounts[0].EUUID)
	require.NoError(t, err)
	require.True(
		t,
		reflect.DeepEqual(*childAccount, childAccounts[0]),
		"child accounts should be equal",
		cmp.Diff(*childAccount, childAccounts[0]),
	)

	token, err := client.CreateChildAccountToken(context.Background(), childAccount.EUUID)
	require.NoError(t, err)
	require.Greater(t, len(token.Token), 0)
}
