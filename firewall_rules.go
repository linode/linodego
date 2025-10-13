package linodego

import (
	"context"
	"encoding/json"
)

// NetworkProtocol enum type
type NetworkProtocol string

// NetworkProtocol enum values
const (
	TCP     NetworkProtocol = "TCP"
	UDP     NetworkProtocol = "UDP"
	ICMP    NetworkProtocol = "ICMP"
	IPENCAP NetworkProtocol = "IPENCAP"
)

// NetworkAddresses are arrays of ipv4 and v6 addresses
type NetworkAddresses struct {
	IPv4 *[]string `json:"ipv4,omitempty"`
	IPv6 *[]string `json:"ipv6,omitempty"`
}

// A FirewallRule is a whitelist of ports, protocols, and addresses for which traffic should be allowed.
// The ipv4/ipv6 address lists may contain Prefix List tokens (for example, "pl::..." or "pl:system:...")
// in addition to literal IP addresses.
type FirewallRule struct {
	Action      string           `json:"action"`
	Label       string           `json:"label"`
	Description string           `json:"description,omitempty"`
	Ports       string           `json:"ports,omitempty"`
	Protocol    NetworkProtocol  `json:"protocol"`
	Addresses   NetworkAddresses `json:"addresses"`

	// FirewallRule references one `Rule Set` by ID. When provided, this entry
	// represents a reference and should be mutually exclusive with ordinary
	// rule fields according to the API contract.
	Ruleset int `json:"ruleset,omitempty"`
}

// MarshalJSON ensures that when a rule references a Rule Set (Ruleset != 0),
// only the reference shape { "ruleset": <id> } is emitted. Otherwise, the
// ordinary rule fields are emitted without the ruleset key.
func (r FirewallRule) MarshalJSON() ([]byte, error) {
	if r.Ruleset != 0 {
		type rulesetOnly struct {
			Ruleset int `json:"ruleset"`
		}

		return json.Marshal(rulesetOnly{Ruleset: r.Ruleset})
	}

	type normal struct {
		Action      string           `json:"action"`
		Label       string           `json:"label"`
		Description string           `json:"description,omitempty"`
		Ports       string           `json:"ports,omitempty"`
		Protocol    NetworkProtocol  `json:"protocol"`
		Addresses   NetworkAddresses `json:"addresses"`
	}

	return json.Marshal(normal{
		Action:      r.Action,
		Label:       r.Label,
		Description: r.Description,
		Ports:       r.Ports,
		Protocol:    r.Protocol,
		Addresses:   r.Addresses,
	})
}

// FirewallRuleSet is a pair of inbound and outbound rules that specify what network traffic should be allowed.
type FirewallRuleSet struct {
	Inbound        []FirewallRule `json:"inbound"`
	InboundPolicy  string         `json:"inbound_policy"`
	Outbound       []FirewallRule `json:"outbound"`
	OutboundPolicy string         `json:"outbound_policy"`
}

// GetFirewallRules gets the FirewallRuleSet for the given Firewall.
func (c *Client) GetFirewallRules(ctx context.Context, firewallID int) (*FirewallRuleSet, error) {
	e := formatAPIPath("networking/firewalls/%d/rules", firewallID)
	return doGETRequest[FirewallRuleSet](ctx, c, e)
}

// UpdateFirewallRules updates the FirewallRuleSet for the given Firewall
func (c *Client) UpdateFirewallRules(ctx context.Context, firewallID int, rules FirewallRuleSet) (*FirewallRuleSet, error) {
	e := formatAPIPath("networking/firewalls/%d/rules", firewallID)
	return doPUTRequest[FirewallRuleSet](ctx, c, e, rules)
}
