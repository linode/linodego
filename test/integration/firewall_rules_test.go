package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	testFirewallRuleSet = linodego.FirewallRulesCreateOptions{
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

func TestFirewallRules_ExtendedProtocols(t *testing.T) {
	label := fmt.Sprintf("fw-ext-%s", getUniqueText())

	rules := linodego.FirewallRulesCreateOptions{
		Inbound: []linodego.FirewallRuleInbound{
			{
				Label:    "linodego-fwrule-all",
				Action:   "ACCEPT",
				Protocol: linodego.AllNetworkProtocols,
				Addresses: linodego.NetworkAddresses{
					IPv4: []string{"0.0.0.0/0"},
				},
			},
			{
				Label:    "linodego-fwrule-numeric",
				Action:   "ACCEPT",
				Protocol: linodego.NetworkProtocol("50"),
				Addresses: linodego.NetworkAddresses{
					IPv4: []string{"0.0.0.0/0"},
				},
			},
			{
				Label:    "linodego-fwrule-tcp-numeric",
				Action:   "ACCEPT",
				Ports:    "443",
				Protocol: linodego.NetworkProtocol("6"),
				Addresses: linodego.NetworkAddresses{
					IPv4: []string{"0.0.0.0/0"},
				},
			},
		},
		InboundPolicy:  "DROP",
		OutboundPolicy: "ACCEPT",
	}

	client, firewall, teardown, err := setupFirewall(t, []firewallModifier{
		func(createOpts *linodego.FirewallCreateOptions) {
			createOpts.Label = label
			createOpts.Rules = rules
		},
	}, "fixtures/TestFirewall_CreateWithExtendedProtocols")
	t.Cleanup(teardown)
	require.NoError(t, err)

	require.Len(t, firewall.Rules.Inbound, 3)

	assert.Equal(t, linodego.AllNetworkProtocols, firewall.Rules.Inbound[0].Protocol)
	assert.Empty(t, firewall.Rules.Inbound[0].Ports)

	assert.Equal(t, linodego.NetworkProtocol("50"), firewall.Rules.Inbound[1].Protocol)
	assert.Empty(t, firewall.Rules.Inbound[1].Ports)

	assert.Equal(t, linodego.NetworkProtocol("6"), firewall.Rules.Inbound[2].Protocol)
	assert.Equal(t, "443", firewall.Rules.Inbound[2].Ports)

	result, err := client.GetFirewall(context.Background(), firewall.ID)
	require.NoErrorf(t, err, "failed to get firewall %d", firewall.ID)

	require.Len(t, result.Rules.Inbound, 3)

	assert.Equal(t, linodego.AllNetworkProtocols, result.Rules.Inbound[0].Protocol)
	assert.Empty(t, result.Rules.Inbound[0].Ports)

	assert.Equal(t, linodego.NetworkProtocol("50"), result.Rules.Inbound[1].Protocol)
	assert.Empty(t, result.Rules.Inbound[1].Ports)

	assert.Equal(t, linodego.NetworkProtocol("6"), result.Rules.Inbound[2].Protocol)
	assert.Equal(t, "443", result.Rules.Inbound[2].Ports)
}

func TestFirewallRules_Update(t *testing.T) {
	client, firewall, teardown, err := setupFirewall(t, []firewallModifier{}, "fixtures/TestFirewallRules_Update")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	newRules := linodego.FirewallRulesUpdateOptions{
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
