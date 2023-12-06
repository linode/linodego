package integration

import (
	"context"
	"testing"
)

func TestNodeBalancerFirewalls_List(t *testing.T) {
	client, nodebalancer, teardown, err := setupNodeBalancer(t,
		"fixtures/TestNodeBalancerFirewalls_List")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	result, err := client.ListNodeBalancerFirewalls(context.Background(), nodebalancer.ID, nil)
	if err != nil {
		t.Errorf("Error listing Firewalls, expected struct, got error %v", err)
	}

	if len(result) == 0 {
		t.Errorf("Expected a list of Firewalls, but got none: %v", err)
	}
}
