package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestAccountMaintenances_List(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestAccountMaintenances_List")
	defer fixtureTeardown()

	listOpts := linodego.NewListOptions(0, "")
	_, err := client.ListMaintenances(context.Background(), listOpts)
	if err != nil {
		t.Errorf("Error listing maintenances, expected array, got error %v", err)
	}
}

func TestAccountMaintenancePolicies_List(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestAccountMaintenancePolicies_List")
	defer fixtureTeardown()

	policies, err := client.ListMaintenancePolicies(context.Background())
	if err != nil {
		t.Fatalf("Error listing maintenance policies: %v", err)
	}

	if len(policies) == 0 {
		t.Fatal("Expected at least one maintenance policy, got none")
	}

	for _, policy := range policies {
		if policy.Sulg == "" {
			t.Error("Policy Sulg should not be empty")
		}
		if policy.Label == "" {
			t.Error("Policy Label should not be empty")
		}
		if policy.Type == "" {
			t.Error("Policy Type should not be empty")
		}
		if policy.NotificationPeriodSec <= 0 {
			t.Errorf("NotificationPeriodSec should be positive, got %d", policy.NotificationPeriodSec)
		}
	}
}
