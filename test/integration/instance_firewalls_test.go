package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestInstanceFirewalls_List(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceFirewall(t, []firewallModifier{
		func(createOpts *linodego.FirewallCreateOptions) {
			createOpts.Label = "linodego-fw-ins-test"
		},
	}, "fixtures/TestInstanceFirewalls_List")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	result, err := client.ListInstanceFirewalls(context.Background(), instance.ID, nil)
	if err != nil {
		t.Errorf("Error listing Firewalls, expected struct, got error %v", err)
	}

	if len(result) == 0 {
		t.Errorf("Expected a list of Firewalls, but got none: %v", err)
	}
}

func setupInstanceFirewall(t *testing.T, firewallModifiers []firewallModifier, fixturesYaml string) (*linodego.Client, *linodego.Instance, *linodego.Firewall, func(), error) {
	t.Helper()
	client, instance, instanceTeardown, err := setupInstance(t, fixturesYaml,
		func(client *linodego.Client, opts *linodego.InstanceCreateOptions) {
			opts.Label = "linodego-fw-inst-test"
		})
	device := linodego.DevicesCreationOptions{Linodes: []int{instance.ID}}
	firewallModifiers = append(firewallModifiers,
		func(createOpts *linodego.FirewallCreateOptions) {
			createOpts.Devices = device
		})
	firewall, firewallTeardown, err := createFirewall(t, client, firewallModifiers...)

	teardown := func() {
		firewallTeardown()
		instanceTeardown()
	}
	return client, instance, firewall, teardown, err
}
