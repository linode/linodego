package integration

import (
	"context"
	"slices"
	"testing"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/require"
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
	require.GreaterOrEqual(t, regionIdx, 0, "no region with Placement Group capability found")

	region := regions[regionIdx]

	require.NotNil(t, region.PlacementGroupLimits)
	require.NotZero(t, region.PlacementGroupLimits.MaximumLinodesPerPG)
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
	require.GreaterOrEqual(t, regionIdx, 0, "no region with Block Storage Encryption capability found")
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
	require.GreaterOrEqual(t, regionIdx, 0, "no region with Kubernetes Enterprise capability found")
}

func TestRegions_customVPCIPv4Ranges(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestRegions_customVPCIPv4Ranges")
	defer teardown()

	regions, err := client.ListRegions(context.Background(), nil)
	require.NoError(t, err)

	// Filtering is not currently supported on capabilities
	regionIdx := slices.IndexFunc(regions, func(region linodego.Region) bool {
		return slices.Contains(region.Capabilities, "Custom VPC IPv4 Ranges")
	})
	require.GreaterOrEqual(t, regionIdx, 0, "no region with Custom VPC IPv4 Ranges capability found")
}

func TestRegionsMonitorsSection(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestRegionsMonitorsSection")
	defer teardown()

	regions, err := client.ListRegions(context.Background(), nil)
	require.NoError(t, err)
	found := false
	for _, region := range regions {
		if region.Monitors.Alerts != nil || region.Monitors.Metrics != nil {
			found = true
			// Validate Alerts
			for _, alert := range region.Monitors.Alerts {
				require.NotEmpty(t, alert, "Alert should not be empty")
			}
			// Validate Metrics
			for _, metric := range region.Monitors.Metrics {
				require.NotEmpty(t, metric, "Metric should not be empty")
			}
		}
	}
	require.True(t, found, "At least one region should have monitors section populated")
}
