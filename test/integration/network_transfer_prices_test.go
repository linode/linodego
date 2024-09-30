package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestNetworkTransferPrice_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestNetworkTransferPrice_List")
	defer teardown()

	prices, err := client.ListNetworkTransferPrices(context.Background(), nil)
	require.NoError(t, err)
	require.Greater(t, len(prices), 0)

	for _, price := range prices {
		validateNetworkTransferPrice(t, price)
	}
}

func validateNetworkTransferPrice(
	t *testing.T,
	price linodego.NetworkTransferPrice,
) {
	require.NotEmpty(t, price.ID)
	require.NotEmpty(t, price.Label)

	// NOTE: We do not check for monthly prices here because it is
	// explicitly set to null.
	require.Greater(t, price.Price.Hourly, 0.0)
	require.GreaterOrEqual(t, price.Transfer, 0)

	for _, regionPrice := range price.RegionPrices {
		require.NotEmpty(t, regionPrice.ID)
		require.Greater(t, regionPrice.Hourly, 0.0)
	}
}
