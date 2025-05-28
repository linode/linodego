package integration

import (
	"context"
	"testing"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/linode/linodego"
)

var testFirewallCreateOpts = linodego.FirewallCreateOptions{
	Label: "linodego-fw-test",
	Rules: testFirewallRuleSet, // borrowed from firewall_rules.test.go
	Tags:  []string{"testing"},
}

// ignoreNetworkAddresses negates comparing IP addresses. Because of fixture sanitization,
// these addresses will be changed to bogus values when running tests.
var ignoreNetworkAddresses = cmpopts.IgnoreFields(linodego.FirewallRule{}, "Addresses")

// ignoreFirewallTimestamps negates comparing created and updated timestamps. Because of
// fixture sanitization, these addresses will be changed to bogus values when running tests.
var ignoreFirewallTimestamps = cmpopts.IgnoreFields(linodego.Firewall{}, "Created", "Updated")

func TestFirewalls_List_smoke(t *testing.T) {
	client, _, teardown, err := setupFirewall(t, []firewallModifier{
		func(createOpts *linodego.FirewallCreateOptions) {
			createOpts.Label = "linodego-fw-test"
		},
	}, "fixtures/TestFirewalls_List")
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

func TestFirewall_Get(t *testing.T) {
	rules := linodego.FirewallRuleSet{
		Inbound: []linodego.FirewallRule{
			{
				Label:    "linodego-fwrule-test",
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
			createOpts.Label = "linodego-fw-test"
			createOpts.Rules = rules
		},
	}, "fixtures/TestFirewall_Get")
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

func TestFirewall_Update(t *testing.T) {
	rules := linodego.FirewallRuleSet{
		InboundPolicy: "ACCEPT",
		Inbound: []linodego.FirewallRule{
			{
				Label:    "linodego-fwrule-test",
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
			createOpts.Label = "linodego-fw-test"
			createOpts.Rules = rules
			createOpts.Tags = []string{"test"}
		},
	}, "fixtures/TestFirewall_Update")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	updateOpts := firewall.GetUpdateOptions()
	updateOpts.Status = linodego.FirewallDisabled
	updateOpts.Label = firewall.Label + "-updated"
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

func TestFirewallSettings_Get(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestFirewallSettings_Get")
	defer fixtureTeardown()

	settings, err := client.GetFirewallSettings(context.Background())
	if err != nil {
		t.Fatalf("Error getting firewall settings: %v", err)
	}

	if settings.DefaultFirewallIDs.Linode == 0 &&
	settings.DefaultFirewallIDs.NodeBalancer == 0 &&
	settings.DefaultFirewallIDs.PublicInterface == 0 &&
	settings.DefaultFirewallIDs.VPCInterface == 0 {
	t.Log("No default firewall IDs set — this is acceptable in a fresh test environment.")
	}

}

func TestFirewallSettings_UpdateAllFields(t *testing.T) {
	label := fmt.Sprintf("fw-allfields-%s", getUniqueText())
	if len(label) < 3 || len(label) > 32 {
		t.Fatalf("generated label %q is %d chars; must be 3–32", label, len(label))
	}

	client, created, teardown, err := setupFirewall(t, []firewallModifier{
		func(opts *linodego.FirewallCreateOptions) {
			opts.Label = label
		},
	}, "fixtures/TestFirewallSettings_UpdateAllFields")
	if err != nil {
		t.Fatal(err)
	}
	if created == nil {
		teardown()
		t.Fatal("setupFirewall returned nil firewall")
	}

	defer func() {
		updateOpts := linodego.FirewallSettingsUpdateOptions{
			DefaultFirewallIDs: linodego.DefaultFirewallIDsOptions{
				Linode:          nil,
				NodeBalancer:    nil,
				PublicInterface: nil,
				VPCInterface:    nil,
			},
		}
		_, err := client.UpdateFirewallSettings(context.Background(), updateOpts)
		if err != nil {
			t.Fatalf("failed to unset default firewall IDs: %v", err)
		}
	}()

	opts := linodego.FirewallSettingsUpdateOptions{
		DefaultFirewallIDs: linodego.DefaultFirewallIDsOptions{
			Linode:          &created.ID,
			NodeBalancer:    &created.ID,
			PublicInterface: &created.ID,
			VPCInterface:    &created.ID,
		},
	}
	updated, err := client.UpdateFirewallSettings(context.Background(), opts)
	if err != nil {
		t.Fatalf("Error updating firewall settings: %v", err)
	}
	if updated.DefaultFirewallIDs.Linode != created.ID {
		t.Errorf("Expected Linode default firewall ID %d, got %d", created.ID, updated.DefaultFirewallIDs.Linode)
	}
}

func TestFirewallSettings_UpdatePartial(t *testing.T) {
	label := fmt.Sprintf("fw-partial-%s", getUniqueText())
	if len(label) < 3 || len(label) > 32 {
		t.Fatalf("generated label %q is %d chars; must be 3–32", label, len(label))
	}

	client, created, teardown, err := setupFirewall(t, []firewallModifier{
		func(opts *linodego.FirewallCreateOptions) {
			opts.Label = label
		},
	}, "fixtures/TestFirewallSettings_UpdatePartial")
	if err != nil {
		t.Fatal(err)
	}
	if created == nil {
		teardown()
		t.Fatal("setupFirewall returned nil firewall")
	}

	defer func() {
		updateOpts := linodego.FirewallSettingsUpdateOptions{
			DefaultFirewallIDs: linodego.DefaultFirewallIDsOptions{
				Linode:          nil,
				NodeBalancer:    nil,
				PublicInterface: nil,
				VPCInterface:    nil,
			},
		}
		_, err := client.UpdateFirewallSettings(context.Background(), updateOpts)
		if err != nil {
			t.Fatalf("failed to unset default firewall IDs: %v", err)
		}
	}()

	opts := linodego.FirewallSettingsUpdateOptions{
		DefaultFirewallIDs: linodego.DefaultFirewallIDsOptions{
			Linode: &created.ID,
		},
	}
	updated, err := client.UpdateFirewallSettings(context.Background(), opts)
	if err != nil {
		t.Fatalf("Error updating firewall settings: %v", err)
	}
	if updated.DefaultFirewallIDs.Linode != created.ID {
		t.Errorf("Expected Linode default firewall ID %d, got %d", created.ID, updated.DefaultFirewallIDs.Linode)
	}
}