package integration

import (
	"context"
	"testing"
)

func TestRegionsAvailability_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestRegionsAvailability_List")
	defer teardown()

	testFunc := func(retryT *TRetry) {
		regions, err := client.ListRegionsAvailability(context.Background(), nil)
		if err != nil {
			t.Errorf("Error listing regions availability, expected struct - error %v", err)
		}
		if len(regions) == 0 {
			t.Errorf("Expected a list of regions availability - %v", regions)
		}
	}

	retryStatement(t, 3, testFunc)
}

func TestRegionsVPCAvailability_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestRegionsVPCAvailability_List")
	defer teardown()

	testFunc := func(retryT *TRetry) {
		regions, err := client.ListRegionsVPCAvailability(context.Background(), nil)
		if err != nil {
			t.Errorf("Error listing regions vpc availability, expected struct - error %v", err)
		}
		if len(regions) == 0 {
			t.Errorf("Expected a list of regions vpc availability - %v", regions)
		}
	}

	retryStatement(t, 3, testFunc)
}

func TestRegionVPCAvailability_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestRegionVPCAvailability_Get")
	defer teardown()

	testFunc := func(retryT *TRetry) {
		region, err := client.GetRegionVPCAvailability(context.Background(), "nl-ams")
		if err != nil {
			t.Errorf("Error getting region vpc availability, expected struct - error %v", err)
		}
		if region == nil {
			t.Errorf("Expected a region vpc availability object - %v", region)
		}
	}

	retryStatement(t, 3, testFunc)
}
