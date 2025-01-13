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

	// Check if a specific region exists using slices.ContainsFunc
	exists := slices.ContainsFunc(regions, func(region linodego.Region) bool {
		return region.ID == "us-east"
	})
	assert.True(t, exists, "Expected region list to contain 'us-east'")
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
}
