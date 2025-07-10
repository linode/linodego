package integration

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestMaintenancePolicies_List(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestAccountMaintenancePolicies_List")
	defer fixtureTeardown()

	resp, err := client.R(context.Background()).Get("maintenance/policies")
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(resp.Body(), &result)
	require.NoError(t, err)

	dataRaw, ok := result["data"]
	require.True(t, ok, "Expected 'data' key in response")

	dataJSON, err := json.Marshal(dataRaw)
	require.NoError(t, err)

	var policies []linodego.MaintenancePolicy
	err = json.Unmarshal(dataJSON, &policies)
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
