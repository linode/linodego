package unit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestFirewall_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("firewall_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet(formatMockAPIPath("networking/firewalls"), fixtureData)

	firewalls, err := base.Client.ListFirewalls(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.Len(t, firewalls, 1)

	firewall := firewalls[0]
	assert.Equal(t, 123, firewall.ID)
	assert.Equal(t, "firewall123", firewall.Label)
	assert.Equal(t, linodego.FirewallStatus("enabled"), firewall.Status)

	assert.Equal(t, "DROP", firewall.Rules.InboundPolicy)
	assert.Len(t, firewall.Rules.Inbound, 1)

	inboundRule := firewall.Rules.Inbound[0]
	assert.Equal(t, "firewallrule123", inboundRule.Label)
	assert.Equal(t, "ACCEPT", inboundRule.Action)
	assert.Equal(t, linodego.NetworkProtocol("TCP"), inboundRule.Protocol)
	assert.Equal(t, "22-24, 80, 443", inboundRule.Ports)

	assert.Equal(t, "DROP", firewall.Rules.OutboundPolicy)
	assert.Len(t, firewall.Rules.Outbound, 1)

	outboundRule := firewall.Rules.Outbound[0]
	assert.Equal(t, "firewallrule123", outboundRule.Label)
	assert.Equal(t, "ACCEPT", outboundRule.Action)
	assert.Equal(t, linodego.NetworkProtocol("TCP"), outboundRule.Protocol)
	assert.Equal(t, "22-24, 80, 443", outboundRule.Ports)
	assert.ElementsMatch(t, []string{"192.0.2.0/24", "198.51.100.2/32"}, *outboundRule.Addresses.IPv4)
	assert.ElementsMatch(t, []string{"2001:DB8::/128"}, *outboundRule.Addresses.IPv6)

	assert.ElementsMatch(t, []string{"example tag", "another example"}, firewall.Tags)

	assert.Equal(t, "2018-01-01T00:01:01Z", firewall.Created.Format(time.RFC3339))
	assert.Equal(t, "2018-01-02T00:01:01Z", firewall.Updated.Format(time.RFC3339))
}

func TestFirewall_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("firewall_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.FirewallCreateOptions{
		Label: "firewall123",
		Rules: linodego.FirewallRuleSet{
			InboundPolicy:  "DROP",
			OutboundPolicy: "DROP",
			Inbound: []linodego.FirewallRule{
				{
					Action: "ACCEPT",
					Addresses: linodego.NetworkAddresses{
						IPv4: &[]string{"192.0.2.0/24", "198.51.100.2/32"},
						IPv6: &[]string{"2001:DB8::/128"},
					},
					Description: "An example firewall rule description.",
					Label:       "firewallrule123",
					Ports:       "22-24, 80, 443",
					Protocol:    "TCP",
				},
			},
			Outbound: []linodego.FirewallRule{
				{
					Action: "ACCEPT",
					Addresses: linodego.NetworkAddresses{
						IPv4: &[]string{"192.0.2.0/24", "198.51.100.2/32"},
						IPv6: &[]string{"2001:DB8::/128"},
					},
					Description: "An example firewall rule description.",
					Label:       "firewallrule123",
					Ports:       "22-24, 80, 443",
					Protocol:    "TCP",
				},
			},
		},
		Tags: []string{"example tag", "another example"},
		Devices: linodego.DevicesCreationOptions{
			Interfaces: []int{1, 2, 3},
		},
	}

	base.MockPost(formatMockAPIPath("networking/firewalls"), fixtureData)

	firewall, err := base.Client.CreateFirewall(context.Background(), requestData)
	assert.NoError(t, err)

	assert.NotNil(t, firewall)
	assert.Equal(t, 123, firewall.ID)
	assert.Equal(t, "firewall123", firewall.Label)
	assert.Equal(t, linodego.FirewallStatus("enabled"), firewall.Status)
	assert.ElementsMatch(t, []string{"example tag", "another example"}, firewall.Tags)

	assert.NotNil(t, firewall.Rules)
	assert.Equal(t, "DROP", firewall.Rules.InboundPolicy)
	assert.Equal(t, "DROP", firewall.Rules.OutboundPolicy)

	assert.Len(t, firewall.Rules.Inbound, 1)
	inboundRule := firewall.Rules.Inbound[0]
	assert.Equal(t, "ACCEPT", inboundRule.Action)
	assert.Equal(t, "firewallrule123", inboundRule.Label)
	assert.Equal(t, "An example firewall rule description.", inboundRule.Description)
	assert.Equal(t, "22-24, 80, 443", inboundRule.Ports)
	assert.Equal(t, linodego.NetworkProtocol("TCP"), inboundRule.Protocol)
	assert.ElementsMatch(t, []string{"192.0.2.0/24", "198.51.100.2/32"}, *inboundRule.Addresses.IPv4)
	assert.ElementsMatch(t, []string{"2001:DB8::/128"}, *inboundRule.Addresses.IPv6)

	assert.Len(t, firewall.Rules.Outbound, 1)
	outboundRule := firewall.Rules.Outbound[0]
	assert.Equal(t, "ACCEPT", outboundRule.Action)
	assert.Equal(t, "firewallrule123", outboundRule.Label)
	assert.Equal(t, "An example firewall rule description.", outboundRule.Description)
	assert.Equal(t, "22-24, 80, 443", outboundRule.Ports)
	assert.Equal(t, linodego.NetworkProtocol("TCP"), outboundRule.Protocol)
	assert.ElementsMatch(t, []string{"192.0.2.0/24", "198.51.100.2/32"}, *outboundRule.Addresses.IPv4)
	assert.ElementsMatch(t, []string{"2001:DB8::/128"}, *outboundRule.Addresses.IPv6)
}

func TestFirewall_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("firewall_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	firewallID := 123
	base.MockGet(formatMockAPIPath("networking/firewalls/%d", firewallID), fixtureData)

	firewall, err := base.Client.GetFirewall(context.Background(), firewallID)

	assert.NoError(t, err)
	assert.NotNil(t, firewall)

	assert.Equal(t, 123, firewall.ID)
	assert.Equal(t, "firewall123", firewall.Label)
	assert.Equal(t, linodego.FirewallStatus("enabled"), firewall.Status)
	assert.Equal(t, "2018-01-01T00:01:01Z", firewall.Created.Format(time.RFC3339))
	assert.Equal(t, "2018-01-02T00:01:01Z", firewall.Updated.Format(time.RFC3339))
	assert.ElementsMatch(t, []string{"example tag", "another example"}, firewall.Tags)

	assert.NotNil(t, firewall.Rules)
	assert.Equal(t, "DROP", firewall.Rules.InboundPolicy)
	assert.Equal(t, "DROP", firewall.Rules.OutboundPolicy)

	assert.Len(t, firewall.Rules.Inbound, 1)
	inboundRule := firewall.Rules.Inbound[0]
	assert.Equal(t, "ACCEPT", inboundRule.Action)
	assert.Equal(t, "firewallrule123", inboundRule.Label)
	assert.Equal(t, "An example firewall rule description.", inboundRule.Description)
	assert.Equal(t, "22-24, 80, 443", inboundRule.Ports)
	assert.Equal(t, linodego.NetworkProtocol("TCP"), inboundRule.Protocol)
	assert.ElementsMatch(t, []string{"192.0.2.0/24", "198.51.100.2/32"}, *inboundRule.Addresses.IPv4)
	assert.ElementsMatch(t, []string{"2001:DB8::/128"}, *inboundRule.Addresses.IPv6)

	assert.Len(t, firewall.Rules.Outbound, 1)
	outboundRule := firewall.Rules.Outbound[0]
	assert.Equal(t, "ACCEPT", outboundRule.Action)
	assert.Equal(t, "firewallrule123", outboundRule.Label)
	assert.Equal(t, "An example firewall rule description.", outboundRule.Description)
	assert.Equal(t, "22-24, 80, 443", outboundRule.Ports)
	assert.Equal(t, linodego.NetworkProtocol("TCP"), outboundRule.Protocol)
	assert.ElementsMatch(t, []string{"192.0.2.0/24", "198.51.100.2/32"}, *outboundRule.Addresses.IPv4)
	assert.ElementsMatch(t, []string{"2001:DB8::/128"}, *outboundRule.Addresses.IPv6)
}

func TestFirewall_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("firewall_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	firewallID := 123
	base.MockPut(formatMockAPIPath("networking/firewalls/%d", firewallID), fixtureData)

	requestData := linodego.FirewallUpdateOptions{
		Label:  "firewall123",
		Status: "enabled",
		Tags:   &[]string{"updated tag", "another updated tag"},
	}

	firewall, err := base.Client.UpdateFirewall(context.Background(), firewallID, requestData)

	assert.NoError(t, err)
	assert.NotNil(t, firewall)

	assert.Equal(t, 123, firewall.ID)
	assert.Equal(t, "firewall123", firewall.Label)
	assert.Equal(t, linodego.FirewallStatus("enabled"), firewall.Status)
	assert.Equal(t, "2018-01-01T00:01:01Z", firewall.Created.Format(time.RFC3339))
	assert.Equal(t, "2018-01-02T00:01:01Z", firewall.Updated.Format(time.RFC3339))
	assert.ElementsMatch(t, []string{"updated tag", "another updated tag"}, firewall.Tags)

	assert.NotNil(t, firewall.Rules)
	assert.Equal(t, "DROP", firewall.Rules.InboundPolicy)
	assert.Equal(t, "DROP", firewall.Rules.OutboundPolicy)

	assert.Len(t, firewall.Rules.Inbound, 1)
	inboundRule := firewall.Rules.Inbound[0]
	assert.Equal(t, "ACCEPT", inboundRule.Action)
	assert.Equal(t, "firewallrule123", inboundRule.Label)
	assert.Equal(t, "An example firewall rule description.", inboundRule.Description)
	assert.Equal(t, "22-24, 80, 443", inboundRule.Ports)
	assert.Equal(t, linodego.NetworkProtocol("TCP"), inboundRule.Protocol)
	assert.ElementsMatch(t, []string{"192.0.2.0/24", "198.51.100.2/32"}, *inboundRule.Addresses.IPv4)
	assert.ElementsMatch(t, []string{"2001:DB8::/128"}, *inboundRule.Addresses.IPv6)

	assert.Len(t, firewall.Rules.Outbound, 1)
	outboundRule := firewall.Rules.Outbound[0]
	assert.Equal(t, "ACCEPT", outboundRule.Action)
	assert.Equal(t, "firewallrule123", outboundRule.Label)
	assert.Equal(t, "An example firewall rule description.", outboundRule.Description)
	assert.Equal(t, "22-24, 80, 443", outboundRule.Ports)
	assert.Equal(t, linodego.NetworkProtocol("TCP"), outboundRule.Protocol)
	assert.ElementsMatch(t, []string{"192.0.2.0/24", "198.51.100.2/32"}, *outboundRule.Addresses.IPv4)
	assert.ElementsMatch(t, []string{"2001:DB8::/128"}, *outboundRule.Addresses.IPv6)
}

func TestFirewall_Delete(t *testing.T) {
	client := createMockClient(t)

	firewallID := 123

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, fmt.Sprintf("networking/firewalls/%d", firewallID)),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteFirewall(context.Background(), firewallID); err != nil {
		t.Fatal(err)
	}
}

func TestDefaultFirewall_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("default_firewalls_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet(formatMockAPIPath("networking/firewalls/settings"), fixtureData)

	defaultFirewalls, err := base.Client.GetFirewallSettings(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, defaultFirewalls)

	assert.Equal(t, 101, *defaultFirewalls.DefaultFirewallIDs.NodeBalancer)
	assert.Equal(t, 100, *defaultFirewalls.DefaultFirewallIDs.Linode)
	assert.Equal(t, 200, *defaultFirewalls.DefaultFirewallIDs.PublicInterface)
	assert.Equal(t, 200, *defaultFirewalls.DefaultFirewallIDs.VPCInterface)
}

func TestDefaultFirewall_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("default_firewalls_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut(formatMockAPIPath("networking/firewalls/settings"), fixtureData)

	requestData := linodego.FirewallSettingsUpdateOptions{
		DefaultFirewallIDs: &linodego.DefaultFirewallIDsOptions{
			Linode:          linodego.DoublePointer(1),
			NodeBalancer:    linodego.DoublePointer(1),
			VPCInterface:    linodego.DoublePointer(1),
			PublicInterface: linodego.DoublePointer(1),
		},
	}

	defaultFirewalls, err := base.Client.UpdateFirewallSettings(context.Background(), requestData)

	assert.NoError(t, err)
	assert.NotNil(t, defaultFirewalls)

	assert.Equal(t, 1, *defaultFirewalls.DefaultFirewallIDs.NodeBalancer)
	assert.Equal(t, 1, *defaultFirewalls.DefaultFirewallIDs.Linode)
	assert.Equal(t, 1, *defaultFirewalls.DefaultFirewallIDs.PublicInterface)
	assert.Equal(t, 1, *defaultFirewalls.DefaultFirewallIDs.VPCInterface)
}
