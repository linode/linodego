package linodego

import (
	"context"
)

type FirewallTemplate struct {
	Slug  string          `json:"slug"`
	Rules FirewallRuleSet `json:"rules"`
}

// GetFirewallDevice gets a FirewallDevice given an ID
func (c *Client) GetFirewallTemplate(ctx context.Context, slug string) (*FirewallTemplate, error) {
	e := formatAPIPath("networking/firewalls/templates/%s", slug)
	return doGETRequest[FirewallTemplate](ctx, c, e)
}

// ListFirewallDevices get devices associated with a given Firewall
func (c *Client) ListFirewallTemplates(ctx context.Context, opts *ListOptions) ([]FirewallTemplate, error) {
	return getPaginatedResults[FirewallTemplate](ctx, c, "networking/firewalls/templates", opts)
}
