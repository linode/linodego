package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

var (
	testFirewallCreateOpts = linodego.FirewallCreateOptions{
		Label: "label",
		Rules: testFirewallRuleSet, // borrowed from firewall_rules.test.go
		Tags:  []string{"testing"},
	}
)

// TestListFirewalls should return a paginated list of Firewalls
func TestListFirewalls(t *testing.T) {
	client, _, teardown, err := setupFirewall(t, []firewallModifier{
		func(createOpts *linodego.FirewallCreateOptions) {
			createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
		},
	}, "fixtures/TestListFirewalls")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	result, err := client.ListFirewalls(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing Firewalls, expected struct, got error %v", err)
	}

	if len(result) == 0 {
		t.Errorf("Expected a list of Firewalls, but got none: %v", err)
	}
}

type firewallModifier func(*linodego.FirewallCreateOptions)

func setupFirewall(t *testing.T, firewallModifiers []firewallModifier, fixturesYaml string) (*linodego.Client, *linodego.Firewall, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	createOpts := testFirewallCreateOpts
	for _, modifier := range firewallModifiers {
		modifier(&createOpts)
	}

	firewall, err := client.CreateFirewall(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error creating Firewall, expected struct, got error %v", err)
	}

	teardown := func() {
		if err := client.DeleteFirewall(context.Background(), firewall.ID); err != nil {
			t.Errorf("Expected to delete a Firewall, but got %v", err)
		}
		fixtureTeardown()
	}
	return client, firewall, teardown, err
}
