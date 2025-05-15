package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestListVolumeTypes(t *testing.T) {
	// Load the mock fixture for volume types
	fixtureData, err := fixtures.GetFixture("volume_types_list")
	assert.NoError(t, err, "Expected no error when getting fixture")

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request for the volume types endpoint
	base.MockGet("volumes/types", fixtureData)

	// Call the ListVolumeTypes method
	volumeTypes, err := base.Client.ListVolumeTypes(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err, "Expected no error when listing volume types")
	assert.NotEmpty(t, volumeTypes, "Expected non-empty volume types list")

	// Validate the first volume type's details
	assert.Equal(t, "standard", volumeTypes[0].ID, "Expected volume type ID to match")
	assert.Equal(t, "Standard Volume", volumeTypes[0].Label, "Expected volume type label to match")
	assert.Equal(t, 0.10, volumeTypes[0].Price.Hourly, "Expected hourly price to match")
	assert.Equal(t, 10.00, volumeTypes[0].Price.Monthly, "Expected monthly price to match")

	// Validate regional pricing for the first volume type
	assert.NotEmpty(t, volumeTypes[0].RegionPrices, "Expected region prices to be non-empty")
	assert.Equal(t, 0.08, volumeTypes[0].RegionPrices[0].Hourly, "Expected regional hourly price to match")
	assert.Equal(t, 8.00, volumeTypes[0].RegionPrices[0].Monthly, "Expected regional monthly price to match")
}
