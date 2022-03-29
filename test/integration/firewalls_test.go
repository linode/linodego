package integration

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/linode/linodego"
)

var testFirewallCreateOpts = linodego.FirewallCreateOptions{
	Label: "label",
	Rules: testFirewallRuleSet, // borrowed from firewall_rules.test.go
	Tags:  []string{"testing"},
}

// ignoreNetworkAddresses negates comparing IP addresses. Because of fixture sanitization,
// these addresses will be changed to bogus values when running tests.
var ignoreNetworkAddresses = cmpopts.IgnoreFields(linodego.FirewallRule{}, "Addresses")

// ignoreFirewallTimestamps negates comparing created and updated timestamps. Because of
// fixture sanitization, these addresses will be changed to bogus values when running tests.
var ignoreFirewallTimestamps = cmpopts.IgnoreFields(linodego.Firewall{}, "Created", "Updated")

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
				Label:    "test-label",
				Action:   "DROP",
				Protocol: linodego.ICMP,
				Addresses: linodego.NetworkAddresses{
					IPv4: &[]string{"0.0.0.0/0"},
					IPv6: &[]string{"::/0"},
				},
			},
		},
		InboundPolicy:  "ACCEPT",
		OutboundPolicy: "ACCEPT",
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

	if result.Rules.Inbound[0].Label != rules.Inbound[0].Label {
		t.Errorf("Expected firewall rules to be %#v but got %#v", rules, result.Rules)
	}
}

func TestUpdateFirewall(t *testing.T) {
	label := randString(12, lowerBytes, upperBytes) + "-linodego-testing"
	rules := linodego.FirewallRuleSet{
		InboundPolicy: "ACCEPT",
		Inbound: []linodego.FirewallRule{
			{
				Label:    "test-label",
				Action:   "DROP",
				Protocol: linodego.ICMP,
				Addresses: linodego.NetworkAddresses{
					IPv4: &[]string{"0.0.0.0/0"},
				},
			},
		},
		OutboundPolicy: "ACCEPT",
	}

	client, firewall, teardown, err := setupFirewall(t, []firewallModifier{
		func(createOpts *linodego.FirewallCreateOptions) {
			createOpts.Label = label
			createOpts.Rules = rules
			createOpts.Tags = []string{"test"}
		},
	}, "fixtures/TestUpdateFirewall")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	updateOpts := firewall.GetUpdateOptions()
	updateOpts.Status = linodego.FirewallDisabled
	updateOpts.Label = "updatedFirewallLabel"
	updateOpts.Tags = &[]string{}

	updated, err := client.UpdateFirewall(context.Background(), firewall.ID, updateOpts)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(updated.Tags, *updateOpts.Tags) {
		t.Errorf("expected tags to be updated: %s", cmp.Diff(updated.Tags, *updateOpts.Tags))
	}
	if updated.Status != updateOpts.Status {
		t.Errorf("expected status %s but got %s", updateOpts.Status, updated.Status)
	}
	if updated.Label != updateOpts.Label {
		t.Errorf(`expected label to be "%s" but got "%s"`, updateOpts.Label, updated.Label)
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
