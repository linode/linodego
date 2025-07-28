package integration

import (
	"context"
	"testing"
)

func TestFirewallTemplates_List(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestFirewallTemplates_List")
	defer fixtureTeardown()

	result, err := client.ListFirewallTemplates(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing firewall templates, expected struct, got error %v", err)
	}

	if len(result) == 0 {
		t.Errorf("Expected a list of Firewalls, but got none: %v", err)
	}
}

func TestFirewallTemplate_Get(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestFirewallTemplate_Get")
	defer fixtureTeardown()

	result, err := client.GetFirewallTemplate(context.Background(), "public")
	if err != nil {
		t.Errorf("Error listing firewall templates, expected struct, got error %v", err)
	}

	if result.Rules.InboundPolicy != "DROP" {
		t.Errorf(
			"Expected inbound_policy for the public firewall template to be 'DROP', but got: %q",
			result.Rules.InboundPolicy,
		)
	}

	if result.Rules.OutboundPolicy != "ACCEPT" {
		t.Errorf(
			"Expected outbound_policy for the public firewall template to be 'ACCEPT', but got: %q",
			result.Rules.OutboundPolicy,
		)
	}
}
