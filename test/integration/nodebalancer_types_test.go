package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestNodeBalancerType_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestNodeBalancerType_List")
	defer teardown()

	nbTypes, err := client.ListNodeBalancerTypes(context.Background(), nil)
	require.NoError(t, err)
	require.Greater(t, len(nbTypes), 0)

	for _, nbType := range nbTypes {
		validateNodeBalancerType(t, nbType)
	}
}

func validateNodeBalancerType(
	t *testing.T,
	nbType linodego.NodeBalancerType,
) {
	require.NotEmpty(t, nbType.ID)
	require.NotEmpty(t, nbType.Label)

	require.Greater(t, nbType.Price.Hourly, 0.0)
	require.Greater(t, nbType.Price.Monthly, 0.0)
	require.GreaterOrEqual(t, nbType.Transfer, 0)

	for _, regionPrice := range nbType.RegionPrices {
		require.NotEmpty(t, regionPrice.ID)
		require.Greater(t, regionPrice.Hourly, 0.0)
		require.Greater(t, regionPrice.Monthly, 0.0)
	}
}
