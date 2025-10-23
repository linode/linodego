package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestFirewallRuleSets_CRUD(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestFirewallRuleSets_CRUD")
	defer fixtureTeardown()

	ctx := context.Background()

	label := "rs-51452000"
	createOpts := linodego.RuleSetCreateOptions{
		Label:       label,
		Description: "Allow inbound HTTP",
		Type:        linodego.FirewallRuleSetTypeInbound,
		Rules: []linodego.FirewallRule{
			{
				Label:    "allow-http",
				Action:   "ACCEPT",
				Protocol: linodego.NetworkProtocol("TCP"),
				Ports:    "80",
				Addresses: linodego.NetworkAddresses{
					IPv4: &[]string{"0.0.0.0/0"},
				},
			},
		},
	}

	ruleSet, err := client.CreateFirewallRuleSet(ctx, createOpts)
	if err != nil {
		t.Fatalf("failed to create firewall rule set: %v", err)
	}

	t.Cleanup(func() {
		if ruleSet != nil {
			if err := client.DeleteFirewallRuleSet(ctx, ruleSet.ID); err != nil {
				t.Fatalf("failed to cleanup firewall rule set %d: %v", ruleSet.ID, err)
			}
		}
	})

	if ruleSet.Label != createOpts.Label {
		t.Fatalf("expected label %q but got %q", createOpts.Label, ruleSet.Label)
	}
	if ruleSet.Type != createOpts.Type {
		t.Fatalf("expected type %q but got %q", createOpts.Type, ruleSet.Type)
	}

	fetched, err := client.GetFirewallRuleSet(ctx, ruleSet.ID)
	if err != nil {
		t.Fatalf("failed to fetch firewall rule set %d: %v", ruleSet.ID, err)
	}
	if fetched.Label != createOpts.Label {
		t.Fatalf("expected fetched label %q but got %q", createOpts.Label, fetched.Label)
	}
	if len(fetched.Rules) == 0 {
		t.Fatal("expected fetched rule set to contain rules")
	}

	updatedLabel := label + "-updated"
	updatedDescription := "Updated description"
	updatedRules := []linodego.FirewallRule{
		{
			Label:    "allow-https",
			Action:   "ACCEPT",
			Protocol: linodego.NetworkProtocol("TCP"),
			Ports:    "443",
			Addresses: linodego.NetworkAddresses{
				IPv6: &[]string{"::/0"},
			},
		},
	}

	updateOpts := linodego.RuleSetUpdateOptions{
		Label:       &updatedLabel,
		Description: &updatedDescription,
		Rules:       &updatedRules,
	}

	updated, err := client.UpdateFirewallRuleSet(ctx, ruleSet.ID, updateOpts)
	if err != nil {
		t.Fatalf("failed to update firewall rule set %d: %v", ruleSet.ID, err)
	}
	if updated.Label != updatedLabel {
		t.Fatalf("expected updated label %q but got %q", updatedLabel, updated.Label)
	}
	if updated.Description != updatedDescription {
		t.Fatalf("expected updated description %q but got %q", updatedDescription, updated.Description)
	}
	if len(updated.Rules) != len(updatedRules) {
		t.Fatalf("expected %d updated rules but got %d", len(updatedRules), len(updated.Rules))
	}

	list, err := client.ListFirewallRuleSets(ctx, nil)
	if err != nil {
		t.Fatalf("failed to list firewall rule sets: %v", err)
	}

	found := false
	for _, rs := range list {
		if rs.ID == updated.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("updated firewall rule set %d not found in list", updated.ID)
	}

	if err := client.DeleteFirewallRuleSet(ctx, ruleSet.ID); err != nil {
		t.Fatalf("failed to delete firewall rule set %d: %v", ruleSet.ID, err)
	}
	ruleSet = nil
}
