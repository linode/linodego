package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestLKEType_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLKEType_List")
	defer teardown()

	lkeTypes, err := client.ListLKETypes(context.Background(), nil)
	require.NoError(t, err)
	require.Greater(t, len(lkeTypes), 0)

	for _, lkeType := range lkeTypes {
		validateLKEType(t, lkeType)
	}
}

func validateLKEType(
	t *testing.T,
	lkeType linodego.LKEType,
) {
	require.NotEmpty(t, lkeType.ID)
	require.NotEmpty(t, lkeType.Label)

	// NOTE: We use >= 0 here because this is treated as an additional
	// cost on top of the base LKE cluster price, meaning SA has its
	// prices set to 0.
	require.GreaterOrEqual(t, lkeType.Price.Hourly, 0.0)
	require.GreaterOrEqual(t, lkeType.Price.Monthly, 0.0)
	require.GreaterOrEqual(t, lkeType.Transfer, 0)

	for _, regionPrice := range lkeType.RegionPrices {
		require.NotEmpty(t, regionPrice.ID)
		require.GreaterOrEqual(t, regionPrice.Hourly, 0.0)
		require.GreaterOrEqual(t, regionPrice.Monthly, 0.0)
	}
}
