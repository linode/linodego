package integration

import (
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
)

var (
	testFirewallRuleSet = linodego.FirewallRuleSet{
		Inbound:  []linodego.FirewallRule{testFirewallRule},
		Outbound: []linodego.FirewallRule{testFirewallRule},
	}
)
