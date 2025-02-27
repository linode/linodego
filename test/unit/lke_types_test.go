package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLKETypes_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_types_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("lke/types", fixtureData)

	types, err := base.Client.ListLKETypes(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, types, 2)

	// Validate first LKE type
	lkeType1 := types[0]
	assert.Equal(t, "g6-standard-1", lkeType1.ID)
	assert.Equal(t, "Standard 1GB", lkeType1.Label)

	assert.Equal(t, 0.0075, lkeType1.Price.Hourly)
	assert.Equal(t, 5.00, lkeType1.Price.Monthly)

	assert.Len(t, lkeType1.RegionPrices, 2)
	assert.Equal(t, 0.0074, lkeType1.RegionPrices[0].Hourly)
	assert.Equal(t, 4.99, lkeType1.RegionPrices[0].Monthly)
	assert.Equal(t, 0.0076, lkeType1.RegionPrices[1].Hourly)
	assert.Equal(t, 5.01, lkeType1.RegionPrices[1].Monthly)

	// Validate second LKE type
	lkeType2 := types[1]
	assert.Equal(t, "g6-standard-2", lkeType2.ID)
	assert.Equal(t, "Standard 2GB", lkeType2.Label)

	assert.Equal(t, 0.015, lkeType2.Price.Hourly)
	assert.Equal(t, 10.00, lkeType2.Price.Monthly)

	assert.Len(t, lkeType2.RegionPrices, 2)
	assert.Equal(t, 0.0148, lkeType2.RegionPrices[0].Hourly)
	assert.Equal(t, 9.90, lkeType2.RegionPrices[0].Monthly)
	assert.Equal(t, 0.0152, lkeType2.RegionPrices[1].Hourly)
	assert.Equal(t, 10.10, lkeType2.RegionPrices[1].Monthly)
}
