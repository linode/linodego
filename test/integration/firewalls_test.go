package integration

import (
	"context"
	"reflect"
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
			createOpts.Label = randString(12, lowerBytes, upperBytes) + "-linodego-testing"
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

func TestGetFirewall(t *testing.T) {
	label := randString(12, lowerBytes, upperBytes) + "-linodego-testing"
	rules := linodego.FirewallRuleSet{
		Inbound: []linodego.FirewallRule{
			{
				Protocol: linodego.ICMP,
				Addresses: linodego.NetworkAddresses{
					IPv4: []string{"10.20.30.40/0"},
					IPv6: []string{"1234::5678/0"},
				},
			},
		},
	}
	client, created, teardown, err := setupFirewall(t, []firewallModifier{
		func(createOpts *linodego.FirewallCreateOptions) {
			createOpts.Label = label
			createOpts.Rules = rules
		},
	}, "fixtures/TestGetFirewall")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	result, err := client.GetFirewall(context.Background(), created.ID)
	if err != nil {
		t.Errorf("failed to get newly created firewall %d: %s", created.ID, err)
	}

	if !reflect.DeepEqual(result.Rules, rules) {
		t.Errorf("Expected firewall rules to be %#v but got %#v", rules, result.Rules)
	}
}

type firewallModifier func(*linodego.FirewallCreateOptions)

func createFirewall(t *testing.T, client *linodego.Client, firewallModifiers ...firewallModifier) (*linodego.Firewall, func(), error) {
	t.Helper()

	createOpts := testFirewallCreateOpts
	for _, modifier := range firewallModifiers {
		modifier(&createOpts)
	}

	firewall, err := client.CreateFirewall(context.Background(), createOpts)
	if err != nil {
		t.Errorf("failed to create firewall: %s", err)
	}

	teardown := func() {
		if err := client.DeleteFirewall(context.Background(), firewall.ID); err != nil {
			t.Errorf("failed to delete firewall: %s", err)
		}
	}
	return firewall, teardown, nil
}

func setupFirewall(t *testing.T, firewallModifiers []firewallModifier, fixturesYaml string) (*linodego.Client, *linodego.Firewall, func(), error) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	firewall, firewallTeardown, err := createFirewall(t, client, firewallModifiers...)

	teardown := func() {
		firewallTeardown()
		fixtureTeardown()
	}
	return client, firewall, teardown, err
}
