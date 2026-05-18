package integration

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/linode/linodego"
)

var (
	testFirewallRuleInbound = linodego.FirewallRuleInbound{
		Label:    "go-fwrule-test",
		Action:   "ACCEPT",
		Ports:    "22",
		Protocol: "TCP",
		Addresses: linodego.NetworkAddresses{
			IPv4: []string{"0.0.0.0/0"},
			IPv6: []string{"::0/0"},
		},
	}
	testFirewallRuleOutbound = linodego.FirewallRuleOutbound{
		Label:    "go-fwrule-test",
		Action:   "ACCEPT",
		Ports:    "22",
		Protocol: "TCP",
		Addresses: linodego.NetworkAddresses{
			IPv4: []string{"0.0.0.0/0"},
			IPv6: []string{"::0/0"},
		},
	}

	testFirewallRuleSet = linodego.FirewallRules{
		Inbound:        []linodego.FirewallRuleInbound{testFirewallRuleInbound},
		InboundPolicy:  "ACCEPT",
		Outbound:       []linodego.FirewallRuleOutbound{testFirewallRuleOutbound},
		OutboundPolicy: "ACCEPT",
	}
)

func TestFirewallRules_Get_smoke(t *testing.T) {
	client, firewall, teardown, err := setupFirewall(t, []firewallModifier{func(createOpts *linodego.FirewallCreateOptions) {
		createOpts.Rules = testFirewallRuleSet
	}}, "fixtures/TestFirewallRules_Get")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	rules, err := client.GetFirewallRules(context.Background(), firewall.ID)
	if err != nil {
		t.Error(err)
	}

	if rules.Version <= 0 {
		t.Errorf("expected non-zero rules version, got %d", rules.Version)
	}

	if rules.Fingerprint == "" {
		t.Error("expected non-empty rules fingerprint")
	}

	if rules.InboundPolicy != testFirewallRuleSet.InboundPolicy {
		t.Errorf("expected inbound policy %q, got %q", testFirewallRuleSet.InboundPolicy, rules.InboundPolicy)
	}

	if rules.OutboundPolicy != testFirewallRuleSet.OutboundPolicy {
		t.Errorf("expected outbound policy %q, got %q", testFirewallRuleSet.OutboundPolicy, rules.OutboundPolicy)
	}

	if !cmp.Equal(rules.Inbound, testFirewallRuleSet.Inbound, ignoreNetworkAddresses) {
		t.Errorf("expected inbound rules to match, but got diff: %s", cmp.Diff(rules.Inbound, testFirewallRuleSet.Inbound, ignoreNetworkAddresses))
	}

	if !cmp.Equal(rules.Outbound, testFirewallRuleSet.Outbound, ignoreNetworkAddresses) {
		t.Errorf("expected outbound rules to match, but got diff: %s", cmp.Diff(rules.Outbound, testFirewallRuleSet.Outbound, ignoreNetworkAddresses))
	}
}

func TestFirewallRules_Update(t *testing.T) {
	client, firewall, teardown, err := setupFirewall(t, []firewallModifier{}, "fixtures/TestFirewallRules_Update")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	newRules := linodego.FirewallRuleSetUpdateOptions{
		Inbound: []linodego.FirewallRuleInbound{
			{
				Label:    testFirewallRuleInbound.Label + "_r",
				Action:   "DROP",
				Ports:    "22",
				Protocol: "TCP",
				Addresses: linodego.NetworkAddresses{
					IPv4: []string{"0.0.0.0/0"},
					IPv6: []string{"::0/0"},
				},
			},
		},
		InboundPolicy:  "ACCEPT",
		OutboundPolicy: "ACCEPT",
	}

	if _, err := client.UpdateFirewallRules(context.Background(), firewall.ID, newRules); err != nil {
		t.Error(err)
	}

	rules, err := client.GetFirewallRules(context.Background(), firewall.ID)
	if err != nil {
		t.Error(err)
	}

	if rules.Version <= 0 {
		t.Errorf("expected non-zero rules version, got %d", rules.Version)
	}

	if rules.Fingerprint == "" {
		t.Error("expected non-empty rules fingerprint")
	}

	if rules.InboundPolicy != newRules.InboundPolicy {
		t.Errorf("expected inbound policy %q, got %q", newRules.InboundPolicy, rules.InboundPolicy)
	}

	if rules.OutboundPolicy != newRules.OutboundPolicy {
		t.Errorf("expected outbound policy %q, got %q", newRules.OutboundPolicy, rules.OutboundPolicy)
	}

	if !cmp.Equal(rules.Inbound, newRules.Inbound, ignoreNetworkAddresses) {
		t.Errorf("expected inbound rules to match, but got diff: %s", cmp.Diff(rules.Inbound, newRules.Inbound, ignoreNetworkAddresses))
	}

	if !cmp.Equal(rules.Outbound, newRules.Outbound, ignoreNetworkAddresses) {
		t.Errorf("expected outbound rules to match, but got diff: %s", cmp.Diff(rules.Outbound, newRules.Outbound, ignoreNetworkAddresses))
	}
}
