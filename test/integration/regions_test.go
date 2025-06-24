package integration

import (
	"context"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/linode/linodego"
)

func TestRegions_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestRegions_List")
	defer teardown()

	testFunc := func(retryT *TRetry) {
		regions, err := client.ListRegions(context.Background(), nil)
		if err != nil {
			t.Errorf("Error listing regions, expected struct - error %v", err)
		}
		if len(regions) == 0 {
			t.Errorf("Expected a list of regions - %v", regions)
		}
	}

	retryStatement(t, 3, testFunc)
}

func TestRegions_pgLimits(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestRegions_pgLimits")
	defer teardown()

	regions, err := client.ListRegions(context.Background(), nil)
	require.NoError(t, err)

	// Filtering is not currently supported on capabilities
	regionIdx := slices.IndexFunc(regions, func(region linodego.Region) bool {
		return slices.Contains(region.Capabilities, "Placement Group")
	})
	require.NotZero(t, regionIdx)

	region := regions[regionIdx]

	require.NotNil(t, region.PlacementGroupLimits)
	require.NotZero(t, region.PlacementGroupLimits.MaximumLinodesPerPG)
	require.NotZero(t, region.PlacementGroupLimits.MaximumPGsPerCustomer)
}

func TestRegions_blockStorageEncryption(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestRegions_blockStorageEncryption")
	defer teardown()

	regions, err := client.ListRegions(context.Background(), nil)
	require.NoError(t, err)

	// Filtering is not currently supported on capabilities
	regionIdx := slices.IndexFunc(regions, func(region linodego.Region) bool {
		return slices.Contains(region.Capabilities, "Block Storage Encryption")
	})
	require.NotZero(t, regionIdx)
}

func TestRegions_kubernetesEnterprise(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestRegions_kubernetesEnterprise")
	defer teardown()

	regions, err := client.ListRegions(context.Background(), nil)
	require.NoError(t, err)

	// Filtering is not currently supported on capabilities
	regionIdx := slices.IndexFunc(regions, func(region linodego.Region) bool {
		return slices.Contains(region.Capabilities, "Kubernetes Enterprise")
	})
	require.NotZero(t, regionIdx)
}
