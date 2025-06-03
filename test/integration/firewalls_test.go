package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/linode/linodego"
)

// Helper to convert int to *int (for restoring firewall IDs)
func intToPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

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

	if settings == nil {
		t.Fatal("Expected firewall settings, got nil")
	}

	// Validate each individual default firewall ID is present (≥0)
	if settings.DefaultFirewallIDs.Linode < 0 {
		t.Errorf("Invalid Linode default firewall ID: %d", settings.DefaultFirewallIDs.Linode)
	}
	if settings.DefaultFirewallIDs.NodeBalancer < 0 {
		t.Errorf("Invalid NodeBalancer default firewall ID: %d", settings.DefaultFirewallIDs.NodeBalancer)
	}
	if settings.DefaultFirewallIDs.PublicInterface < 0 {
		t.Errorf("Invalid PublicInterface default firewall ID: %d", settings.DefaultFirewallIDs.PublicInterface)
	}
	if settings.DefaultFirewallIDs.VPCInterface < 0 {
		t.Errorf("Invalid VPCInterface default firewall ID: %d", settings.DefaultFirewallIDs.VPCInterface)
	}

	// Optional: ensure at least one ID is non-zero
	if settings.DefaultFirewallIDs.Linode == 0 &&
		settings.DefaultFirewallIDs.NodeBalancer == 0 &&
		settings.DefaultFirewallIDs.PublicInterface == 0 &&
		settings.DefaultFirewallIDs.VPCInterface == 0 {
		t.Log("No default firewall IDs are set — may be acceptable in a fresh or isolated environment.")
	}
}

func TestFirewallSettings_UpdateAllFields(t *testing.T) {
	label := fmt.Sprintf("fw-allfields-%s", getUniqueText())

	client, firewall, teardown, err := setupFirewall(t, []firewallModifier{
		func(opts *linodego.FirewallCreateOptions) {
			opts.Label = label
		},
	}, "fixtures/TestFirewallSettings_UpdateAllFields")
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	// Backup original firewall settings before changing them
	originalSettings, err := client.GetFirewallSettings(context.Background())
	if err != nil {
		t.Fatalf("Failed to get original firewall settings: %v", err)
	}

	defer func() {
		restoreOpts := linodego.FirewallSettingsUpdateOptions{
			DefaultFirewallIDs: linodego.DefaultFirewallIDsOptions{
				Linode:          intToPtr(originalSettings.DefaultFirewallIDs.Linode),
				NodeBalancer:    intToPtr(originalSettings.DefaultFirewallIDs.NodeBalancer),
				PublicInterface: intToPtr(originalSettings.DefaultFirewallIDs.PublicInterface),
				VPCInterface:    intToPtr(originalSettings.DefaultFirewallIDs.VPCInterface),
			},
		}
		_, err := client.UpdateFirewallSettings(context.Background(), restoreOpts)
		if err != nil {
			t.Fatalf("Failed to restore original default firewall IDs: %v", err)
		}
	}()

	// Update all default firewall settings to the test firewall
	updateOpts := linodego.FirewallSettingsUpdateOptions{
		DefaultFirewallIDs: linodego.DefaultFirewallIDsOptions{
			Linode:          &firewall.ID,
			NodeBalancer:    &firewall.ID,
			PublicInterface: &firewall.ID,
			VPCInterface:    &firewall.ID,
		},
	}
	updated, err := client.UpdateFirewallSettings(context.Background(), updateOpts)
	if err != nil {
		t.Fatalf("Error updating firewall settings: %v", err)
	}

	// Validate the updates
	if updated.DefaultFirewallIDs.Linode != firewall.ID {
		t.Errorf("Expected Linode default firewall ID %d, got %d", firewall.ID, updated.DefaultFirewallIDs.Linode)
	}
	if updated.DefaultFirewallIDs.NodeBalancer != firewall.ID {
		t.Errorf("Expected NodeBalancer default firewall ID %d, got %d", firewall.ID, updated.DefaultFirewallIDs.NodeBalancer)
	}
	if updated.DefaultFirewallIDs.PublicInterface != firewall.ID {
		t.Errorf("Expected PublicInterface default firewall ID %d, got %d", firewall.ID, updated.DefaultFirewallIDs.PublicInterface)
	}
	if updated.DefaultFirewallIDs.VPCInterface != firewall.ID {
		t.Errorf("Expected VPCInterface default firewall ID %d, got %d", firewall.ID, updated.DefaultFirewallIDs.VPCInterface)
	}
}

func TestFirewallSettings_UpdatePartial(t *testing.T) {
	label := fmt.Sprintf("fw-partial-%s", getUniqueText())

	client, firewall, teardown, err := setupFirewall(t, []firewallModifier{
		func(opts *linodego.FirewallCreateOptions) {
			opts.Label = label
		},
	}, "fixtures/TestFirewallSettings_UpdatePartial")
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	if firewall == nil {
		t.Fatal("setupFirewall returned nil firewall")
	}

	// Backup original default firewall settings
	originalSettings, err := client.GetFirewallSettings(context.Background())
	if err != nil {
		t.Fatalf("Failed to get original firewall settings: %v", err)
	}

	// Restore original settings after test
	defer func() {
		restoreOpts := linodego.FirewallSettingsUpdateOptions{
			DefaultFirewallIDs: linodego.DefaultFirewallIDsOptions{
				Linode:          intToPtr(originalSettings.DefaultFirewallIDs.Linode),
				NodeBalancer:    intToPtr(originalSettings.DefaultFirewallIDs.NodeBalancer),
				PublicInterface: intToPtr(originalSettings.DefaultFirewallIDs.PublicInterface),
				VPCInterface:    intToPtr(originalSettings.DefaultFirewallIDs.VPCInterface),
			},
		}
		_, err := client.UpdateFirewallSettings(context.Background(), restoreOpts)
		if err != nil {
			t.Fatalf("Failed to restore original default firewall IDs: %v", err)
		}
	}()

	// Update only the Linode default firewall ID
	opts := linodego.FirewallSettingsUpdateOptions{
		DefaultFirewallIDs: linodego.DefaultFirewallIDsOptions{
			Linode: &firewall.ID,
		},
	}
	updated, err := client.UpdateFirewallSettings(context.Background(), opts)
	if err != nil {
		t.Fatalf("Error updating firewall settings: %v", err)
	}

	if updated.DefaultFirewallIDs.Linode != firewall.ID {
		t.Errorf("Expected Linode default firewall ID %d, got %d", firewall.ID, updated.DefaultFirewallIDs.Linode)
	}
}
