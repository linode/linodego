package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccountTransfer_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestAccountTransfer_Get")
	defer teardown()

	transfer, err := client.GetAccountTransfer(context.Background())
	require.NoError(t, err, "Error getting Account Transfer, expected struct")

	require.NotEqual(t, 0, transfer.Billable, "Expected non-zero value for Billable")
	require.NotEqual(t, 0, transfer.Quota, "Expected non-zero value for Quota")
	require.NotEqual(t, 0, transfer.Used, "Expected non-zero value for Used")

	require.NotEmpty(t, transfer.RegionTransfers, "Expected to see region transfers")

	for _, regionTransfer := range transfer.RegionTransfers {
		require.NotEmpty(t, regionTransfer.ID, "Expected region ID to be non-empty")
		require.NotEqual(t, 0, regionTransfer.Billable, "Expected non-zero value for Billable in region %s", regionTransfer.ID)
		require.NotEqual(t, 0, regionTransfer.Quota, "Expected non-zero value for Quota in region %s", regionTransfer.ID)
		require.NotEqual(t, 0, regionTransfer.Used, "Expected non-zero value for Used in region %s", regionTransfer.ID)
	}
}
