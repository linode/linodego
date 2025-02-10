package unit

import (
	"context"
	"fmt"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestListRegions(t *testing.T) {
	// Load the fixture data for regions
	fixtureData, err := fixtures.GetFixture("regions_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("regions", fixtureData)

	regions, err := base.Client.ListRegions(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, regions, "Expected non-empty region list")

	// Validate a specific region using slices.ContainsFunc
	exists := slices.ContainsFunc(regions, func(region linodego.Region) bool {
		return region.ID == "us-east"
	})
	assert.True(t, exists, "Expected region list to contain 'us-east'")

	// Additional assertions
	for _, region := range regions {
		assert.NotEmpty(t, region.Country, "Expected region country to be set")
		assert.NotEmpty(t, region.Capabilities, "Expected region capabilities to be set")
		assert.NotEmpty(t, region.Status, "Expected region status to be set")
		assert.NotEmpty(t, region.Label, "Expected region label to be set")
		assert.NotEmpty(t, region.SiteType, "Expected region site type to be set")
		assert.NotNil(t, region.Resolvers, "Expected region resolvers to be set")
		assert.NotEmpty(t, region.Resolvers.IPv4, "Expected IPv4 resolver to be set")
		assert.NotEmpty(t, region.Resolvers.IPv6, "Expected IPv6 resolver to be set")
		assert.NotNil(t, region.PlacementGroupLimits, "Expected placement group limits to be set")
		if region.PlacementGroupLimits != nil {
			assert.Greater(t, region.PlacementGroupLimits.MaximumPGsPerCustomer, 0, "Expected MaximumPGsPerCustomer to be greater than 0")
			assert.Greater(t, region.PlacementGroupLimits.MaximumLinodesPerPG, 0, "Expected MaximumLinodesPerPG to be greater than 0")
		}
		assert.Contains(t, region.Capabilities, linodego.CapabilityLinodes, "Expected region to support Linodes")
	}
}

func TestGetRegion(t *testing.T) {
	// Load the fixture data for a specific region
	fixtureData, err := fixtures.GetFixture("region_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	regionID := "us-east"
	base.MockGet(fmt.Sprintf("regions/%s", regionID), fixtureData)

	region, err := base.Client.GetRegion(context.Background(), regionID)
	assert.NoError(t, err)
	assert.NotNil(t, region, "Expected region object to be returned")
	assert.Equal(t, "us-east", region.ID, "Expected region ID to be 'us-east'")
	assert.NotEmpty(t, region.Country, "Expected Country field to be populated")
	assert.NotEmpty(t, region.Capabilities, "Expected Capabilities field to be populated")
	assert.NotEmpty(t, region.Status, "Expected Status field to be populated")
	assert.NotEmpty(t, region.Label, "Expected Label field to be populated")
	assert.NotEmpty(t, region.SiteType, "Expected SiteType field to be populated")
	assert.NotNil(t, region.Resolvers, "Expected Resolvers field to be populated")
	assert.NotEmpty(t, region.Resolvers.IPv4, "Expected IPv4 resolver to be set")
	assert.NotEmpty(t, region.Resolvers.IPv6, "Expected IPv6 resolver to be set")
	assert.NotNil(t, region.PlacementGroupLimits, "Expected PlacementGroupLimits field to be set")
	if region.PlacementGroupLimits != nil {
		assert.Greater(t, region.PlacementGroupLimits.MaximumPGsPerCustomer, 0, "Expected MaximumPGsPerCustomer to be greater than 0")
		assert.Greater(t, region.PlacementGroupLimits.MaximumLinodesPerPG, 0, "Expected MaximumLinodesPerPG to be greater than 0")
	}
	assert.Contains(t, region.Capabilities, linodego.CapabilityLinodes, "Expected region to support Linodes")
}

func TestListRegionsAvailability(t *testing.T) {
	// Load the fixture data for region availability
	fixtureData, err := fixtures.GetFixture("regions_availability_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("regions/availability", fixtureData)

	availability, err := base.Client.ListRegionsAvailability(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, availability, "Expected non-empty region availability list")

	// Check if a specific region availability exists using slices.ContainsFunc
	exists := slices.ContainsFunc(availability, func(a linodego.RegionAvailability) bool {
		return a.Region == "us-east" && a.Available
	})
	assert.True(t, exists, "Expected region availability list to contain 'us-east' with available status")

	// Additional assertions
	for _, avail := range availability {
		assert.NotEmpty(t, avail.Plan, "Expected plan to be set")
	}
}

func TestGetRegionAvailability(t *testing.T) {
	// Load the fixture data for a specific region availability
	fixtureData, err := fixtures.GetFixture("region_availability_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	regionID := "us-east"
	base.MockGet(fmt.Sprintf("regions/%s/availability", regionID), fixtureData)

	availability, err := base.Client.GetRegionAvailability(context.Background(), regionID)
	assert.NoError(t, err)
	assert.NotNil(t, availability, "Expected region availability object to be returned")
	assert.Equal(t, "us-east", availability.Region, "Expected region ID to be 'us-east'")
	assert.True(t, availability.Available, "Expected region to be available")
	assert.NotEmpty(t, availability.Plan, "Expected plan to be set")
}
