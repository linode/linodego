package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestFirewallRule_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("firewall_rule_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	firewallID := 123
	base.MockGet(formatMockAPIPath("networking/firewalls/%d/rules", firewallID), fixtureData)

	firewallRule, err := base.Client.GetFirewallRules(context.Background(), firewallID)
	assert.NoError(t, err)
	assert.NotNil(t, firewallRule)

	assert.Equal(t, "DROP", firewallRule.InboundPolicy)
	assert.Equal(t, 1, len(firewallRule.Inbound))
	assert.Equal(t, "ACCEPT", firewallRule.Inbound[0].Action)
	assert.Equal(t, "firewallrule123", firewallRule.Inbound[0].Label)
	assert.Equal(t, "An example firewall rule description.", firewallRule.Inbound[0].Description)
	assert.Equal(t, "22-24, 80, 443", firewallRule.Inbound[0].Ports)
	assert.Equal(t, linodego.NetworkProtocol("TCP"), firewallRule.Inbound[0].Protocol)
	assert.ElementsMatch(t, []string{"192.0.2.0/24", "198.51.100.2/32"}, *firewallRule.Inbound[0].Addresses.IPv4)
	assert.ElementsMatch(t, []string{"2001:DB8::/128"}, *firewallRule.Inbound[0].Addresses.IPv6)

	assert.Equal(t, "DROP", firewallRule.OutboundPolicy)
	assert.Equal(t, 1, len(firewallRule.Outbound))
	assert.Equal(t, "ACCEPT", firewallRule.Outbound[0].Action)
	assert.Equal(t, "firewallrule123", firewallRule.Outbound[0].Label)
	assert.Equal(t, "An example firewall rule description.", firewallRule.Outbound[0].Description)
	assert.Equal(t, "22-24, 80, 443", firewallRule.Outbound[0].Ports)
	assert.Equal(t, linodego.NetworkProtocol("TCP"), firewallRule.Outbound[0].Protocol)
	assert.ElementsMatch(t, []string{"192.0.2.0/24", "198.51.100.2/32"}, *firewallRule.Outbound[0].Addresses.IPv4)
	assert.ElementsMatch(t, []string{"2001:DB8::/128"}, *firewallRule.Outbound[0].Addresses.IPv6)
}

func TestFirewallRule_GetExpansion(t *testing.T) {
	inboundIPv4 := []string{"pl::vpcs:1234"}
	inboundIPv6 := []string{"pl::vpcs:<current>"}
	outboundIPv4 := []string{"pl::vpcs:1234"}
	outboundIPv6 := []string{"pl::vpcs:<current>"}

	mockResponse := linodego.FirewallRuleSet{
		Inbound: []linodego.FirewallRule{
			{
				Action:      "ACCEPT",
				Label:       "accept-inbound-ssh",
				Description: "Accept inbound SSH within the VPC",
				Ports:       "22",
				Protocol:    linodego.NetworkProtocol("TCP"),
				Addresses: linodego.NetworkAddresses{
					IPv4: &inboundIPv4,
					IPv6: &inboundIPv6,
				},
			},
		},
		InboundPolicy: "DROP",
		Outbound: []linodego.FirewallRule{
			{
				Action:      "ACCEPT",
				Label:       "accept-outbound-ssh",
				Description: "Accept outbound SSH within the VPC",
				Ports:       "22",
				Protocol:    linodego.NetworkProtocol("TCP"),
				Addresses: linodego.NetworkAddresses{
					IPv4: &outboundIPv4,
					IPv6: &outboundIPv6,
				},
			},
		},
		OutboundPolicy: "DROP",
	}

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	firewallID := 456
	base.MockGet(formatMockAPIPath("networking/firewalls/%d/rules/expansion", firewallID), mockResponse)

	firewallRuleSet, err := base.Client.GetFirewallRulesExpansion(context.Background(), firewallID)
	assert.NoError(t, err)
	assert.NotNil(t, firewallRuleSet)

	if assert.Len(t, firewallRuleSet.Inbound, 1) {
		assert.Equal(t, "ACCEPT", firewallRuleSet.Inbound[0].Action)
		assert.Equal(t, "accept-inbound-ssh", firewallRuleSet.Inbound[0].Label)
		assert.Equal(t, "Accept inbound SSH within the VPC", firewallRuleSet.Inbound[0].Description)
		assert.Equal(t, "22", firewallRuleSet.Inbound[0].Ports)
		assert.Equal(t, linodego.NetworkProtocol("TCP"), firewallRuleSet.Inbound[0].Protocol)
		if assert.NotNil(t, firewallRuleSet.Inbound[0].Addresses.IPv4) {
			assert.ElementsMatch(t, []string{"pl::vpcs:1234"}, *firewallRuleSet.Inbound[0].Addresses.IPv4)
		}
		if assert.NotNil(t, firewallRuleSet.Inbound[0].Addresses.IPv6) {
			assert.ElementsMatch(t, []string{"pl::vpcs:<current>"}, *firewallRuleSet.Inbound[0].Addresses.IPv6)
		}
	}

	assert.Equal(t, "DROP", firewallRuleSet.InboundPolicy)

	if assert.Len(t, firewallRuleSet.Outbound, 1) {
		assert.Equal(t, "ACCEPT", firewallRuleSet.Outbound[0].Action)
		assert.Equal(t, "accept-outbound-ssh", firewallRuleSet.Outbound[0].Label)
		assert.Equal(t, "Accept outbound SSH within the VPC", firewallRuleSet.Outbound[0].Description)
		assert.Equal(t, "22", firewallRuleSet.Outbound[0].Ports)
		assert.Equal(t, linodego.NetworkProtocol("TCP"), firewallRuleSet.Outbound[0].Protocol)
		if assert.NotNil(t, firewallRuleSet.Outbound[0].Addresses.IPv4) {
			assert.ElementsMatch(t, []string{"pl::vpcs:1234"}, *firewallRuleSet.Outbound[0].Addresses.IPv4)
		}
		if assert.NotNil(t, firewallRuleSet.Outbound[0].Addresses.IPv6) {
			assert.ElementsMatch(t, []string{"pl::vpcs:<current>"}, *firewallRuleSet.Outbound[0].Addresses.IPv6)
		}
	}

	assert.Equal(t, "DROP", firewallRuleSet.OutboundPolicy)
}

func TestFirewallRule_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("firewall_rule_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	firewallID := 123
	base.MockPut(formatMockAPIPath("networking/firewalls/%d/rules", firewallID), fixtureData)

	requestData := linodego.FirewallRuleSet{
		Inbound: []linodego.FirewallRule{
			{
				Action:      "ACCEPT",
				Label:       "firewallrule123",
				Description: "An example firewall rule description.",
				Ports:       "22-24, 80, 443",
				Protocol:    "TCP",
				Addresses: linodego.NetworkAddresses{
					IPv4: &[]string{"192.0.2.0/24", "198.51.100.2/32"},
					IPv6: &[]string{"2001:DB8::/128"},
				},
			},
		},
		InboundPolicy: "DROP",
		Outbound: []linodego.FirewallRule{
			{
				Action:      "ACCEPT",
				Label:       "firewallrule123",
				Description: "An example firewall rule description.",
				Ports:       "22-24, 80, 443",
				Protocol:    "TCP",
				Addresses: linodego.NetworkAddresses{
					IPv4: &[]string{"192.0.2.0/24", "198.51.100.2/32"},
					IPv6: &[]string{"2001:DB8::/128"},
				},
			},
		},
		OutboundPolicy: "DROP",
	}

	firewallRule, err := base.Client.UpdateFirewallRules(context.Background(), firewallID, requestData)

	assert.NoError(t, err)
	assert.NotNil(t, firewallRule)

	assert.Equal(t, "DROP", firewallRule.InboundPolicy)
	assert.Equal(t, 1, len(firewallRule.Inbound))
	assert.Equal(t, "ACCEPT", firewallRule.Inbound[0].Action)
	assert.Equal(t, "firewallrule123", firewallRule.Inbound[0].Label)
	assert.Equal(t, "An example firewall rule description.", firewallRule.Inbound[0].Description)
	assert.Equal(t, "22-24, 80, 443", firewallRule.Inbound[0].Ports)
	assert.Equal(t, linodego.NetworkProtocol("TCP"), firewallRule.Inbound[0].Protocol)
	assert.ElementsMatch(t, []string{"192.0.2.0/24", "198.51.100.2/32"}, *firewallRule.Inbound[0].Addresses.IPv4)
	assert.ElementsMatch(t, []string{"2001:DB8::/128"}, *firewallRule.Inbound[0].Addresses.IPv6)

	assert.Equal(t, "DROP", firewallRule.OutboundPolicy)
	assert.Equal(t, 1, len(firewallRule.Outbound))
	assert.Equal(t, "ACCEPT", firewallRule.Outbound[0].Action)
	assert.Equal(t, "firewallrule123", firewallRule.Outbound[0].Label)
	assert.Equal(t, "An example firewall rule description.", firewallRule.Outbound[0].Description)
	assert.Equal(t, "22-24, 80, 443", firewallRule.Outbound[0].Ports)
	assert.Equal(t, linodego.NetworkProtocol("TCP"), firewallRule.Outbound[0].Protocol)
	assert.ElementsMatch(t, []string{"192.0.2.0/24", "198.51.100.2/32"}, *firewallRule.Outbound[0].Addresses.IPv4)
	assert.ElementsMatch(t, []string{"2001:DB8::/128"}, *firewallRule.Outbound[0].Addresses.IPv6)
}
