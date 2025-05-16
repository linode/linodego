package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestInstanceFirewalls_List(t *testing.T) {
	client, instance, _, err := setupInstanceFirewall(t, []firewallModifier{
		func(createOpts *linodego.FirewallCreateOptions) {
			createOpts.Label = "linodego-fw-ins-test"
		},
	}, "fixtures/TestInstanceFirewalls_List")
	if err != nil {
		t.Error(err)
	}

	result, err := client.ListInstanceFirewalls(context.Background(), instance.ID, nil)
	if err != nil {
		t.Errorf("Error listing Firewalls, expected struct, got error %v", err)
	}

	if len(result) == 0 {
		t.Errorf("Expected a list of Firewalls, but got none: %v", err)
	}
}

func setupInstanceFirewall(t *testing.T, firewallModifiers []firewallModifier, fixturesYaml string) (*linodego.Client, *linodego.Instance, *linodego.Firewall, error) {
	t.Helper()
	client, firewall, firewallTeardown, err := setupFirewall(t, firewallModifiers, fixturesYaml)

	instance, err := createInstance(t, client, false,
		func(client *linodego.Client, opts *linodego.InstanceCreateOptions) {
			opts.Label = "linodego-fw-inst-test"
			opts.FirewallID = firewall.ID
		},
	)

	t.Cleanup(
		func() {
			if err := client.DeleteInstance(context.Background(), instance.ID); err != nil {
				if t != nil {
					t.Errorf("failed deleting test Instance: %s", err)
				}
			}
			firewallTeardown()
		},
	)
	return client, instance, firewall, err
}
