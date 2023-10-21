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
