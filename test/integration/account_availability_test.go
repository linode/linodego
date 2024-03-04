package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestAccountAvailability_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestAccountAvailability_List")
	defer teardown()

	availabilities, err := client.ListAccountAvailabilities(context.Background(), &linodego.ListOptions{})
	if err != nil {
		t.Errorf("Error getting Account Availabilities, expected struct, got error %v", err)
	}

	if len(availabilities) == 0 {
		t.Errorf("Expected to see account availabilities returned.")
	}
}

func TestAccountAvailability_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestAccountAvailability_Get")
	defer teardown()

	regionID := "us-east"
	availability, err := client.GetAccountAvailability(context.Background(), regionID)
	if err != nil {
		t.Errorf("Error getting Account Availability, expected struct, got error %v", err)
	}

	if availability.Region != regionID {
		t.Errorf("expected region ID to be %s; got %s", regionID, availability.Region)
	}
}
