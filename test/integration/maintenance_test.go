package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMaintenancePolicies_List(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestMaintenancePolicies_List")
	defer fixtureTeardown()

	policies, err := client.ListMaintenancePolicies(context.Background(), nil)
	require.NoError(t, err)

	if len(policies) == 0 {
		t.Fatal("Expected at least one maintenance policy, got none")
	}
	for _, policy := range policies {
		if policy.Slug == "" {
			t.Error("Policy Slug should not be empty")
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
