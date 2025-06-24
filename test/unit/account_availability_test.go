package unit

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestAccountAvailabilities_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_availability_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/availability", fixtureData)

	availabilities, err := base.Client.ListAccountAvailabilities(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	// Check specific region "us-central"
	var usCentralAvailability *linodego.AccountAvailability
	for _, availability := range availabilities {
		if availability.Region == "us-central" {
			usCentralAvailability = &availability
			break
		}
	}
	if usCentralAvailability == nil {
		t.Errorf("Expected region 'us-central' to be in the response, but it was not found")
	} else {
		expectedAvailable := []string{"Linodes", "NodeBalancers", "Block Storage", "Kubernetes"}
		if !equalSlices(usCentralAvailability.Available, expectedAvailable) {
			t.Errorf("Expected available resources for 'us-central' to be %v, but got %v", expectedAvailable, usCentralAvailability.Available)
		}

		if len(usCentralAvailability.Unavailable) != 0 {
			t.Errorf("Expected no unavailable resources for 'us-central', but got %v", usCentralAvailability.Unavailable)
		}
	}

	expectedRegionsCount := 40
	if len(availabilities) != expectedRegionsCount {
		t.Errorf("Expected %d regions, but got %d", expectedRegionsCount, len(availabilities))
	}
}

func TestAccountAvailability_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_availability_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	regionID := "us-east"

	base.MockGet(fmt.Sprintf("account/availability/%s", regionID), fixtureData)

	availability, err := base.Client.GetAccountAvailability(context.Background(), regionID)
	assert.NoError(t, err)

	assert.Equal(t, "us-east", availability.Region, "Expected region to be 'us-east'")

	expectedAvailable := []string{"Linodes", "NodeBalancers"}
	assert.ElementsMatch(t, expectedAvailable, availability.Available, "Available resources do not match the expected list")

	expectedUnavailable := []string{"Kubernetes", "Block Storage"}
	assert.ElementsMatch(t, expectedUnavailable, availability.Unavailable, "Unavailable resources do not match the expected list")
}

// Helper function to compare slices in assertion
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aMap := make(map[string]bool)
	for _, v := range a {
		aMap[v] = true
	}
	for _, v := range b {
		if !aMap[v] {
			return false
		}
	}
	return true
}
