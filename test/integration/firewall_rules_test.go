package integration

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/linode/linodego"
)

var (
	testFirewallRule = linodego.FirewallRule{
		Ports:    "22",
		Protocol: "TCP",
		Addresses: linodego.NetworkAddresses{
			IPv4: []string{"0.0.0.0/0"},
			IPv6: []string{"::0/0"},
		},
	}

	testFirewallRuleSet = linodego.FirewallRuleSet{
		Inbound:  []linodego.FirewallRule{testFirewallRule},
		Outbound: []linodego.FirewallRule{testFirewallRule},
	}
)

func TestGetFirewallRules(t *testing.T) {
	client, firewall, teardown, err := setupFirewall(t, []firewallModifier{func(createOpts *linodego.FirewallCreateOptions) {
		createOpts.Rules = testFirewallRuleSet
	}}, "fixtures/TestGetFirewallRules")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	rules, err := client.GetFirewallRules(context.Background(), firewall.ID)
	if !cmp.Equal(rules, &testFirewallRuleSet, ignoreNetworkAddresses) {
		t.Errorf("expected rules to match test rules, but got diff: %s", cmp.Diff(rules, testFirewallRuleSet, ignoreNetworkAddresses))
	}
}

func TestUpdateFirewallRules(t *testing.T) {
	client, firewall, teardown, err := setupFirewall(t, []firewallModifier{}, "fixtures/TestUpdateFirewallRules")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	newRules := linodego.FirewallRuleSet{
		Inbound: []linodego.FirewallRule{
			{
				Ports:    "22",
				Protocol: "TCP",
				Addresses: linodego.NetworkAddresses{
					IPv4: []string{"0.0.0.0/0"},
					IPv6: []string{"::0/0"},
				},
			},
		},
	}

	if _, err := client.UpdateFirewallRules(context.Background(), firewall.ID, newRules); err != nil {
		t.Error(err)
	}

	rules, err := client.GetFirewallRules(context.Background(), firewall.ID)
	if !cmp.Equal(rules, &newRules, ignoreNetworkAddresses) {
		t.Errorf("expected rules to have been updated but got diff: %s", cmp.Diff(rules, &newRules, ignoreNetworkAddresses))
	}
}
