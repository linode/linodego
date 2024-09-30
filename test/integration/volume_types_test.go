package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestVolumeType_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestVolumeType_List")
	defer teardown()

	volumeTypes, err := client.ListVolumeTypes(context.Background(), nil)
	require.NoError(t, err)
	require.Greater(t, len(volumeTypes), 0)

	for _, volumeType := range volumeTypes {
		validateVolumeType(t, volumeType)
	}
}

func validateVolumeType(
	t *testing.T,
	volumeType linodego.VolumeType,
) {
	require.NotEmpty(t, volumeType.ID)
	require.NotEmpty(t, volumeType.Label)

	require.Greater(t, volumeType.Price.Hourly, 0.0)
	require.Greater(t, volumeType.Price.Monthly, 0.0)
	require.GreaterOrEqual(t, volumeType.Transfer, 0)

	for _, regionPrice := range volumeType.RegionPrices {
		require.NotEmpty(t, regionPrice.ID)
		require.Greater(t, regionPrice.Hourly, 0.0)
		require.Greater(t, regionPrice.Monthly, 0.0)
	}
}
