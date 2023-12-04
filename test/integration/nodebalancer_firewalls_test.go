package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestNodeBalancerFirewalls_List(t *testing.T) {
	client, nodebalancer, _, teardown, err := setupNodeBalancerFirewall(t, []firewallModifier{
		func(createOpts *linodego.FirewallCreateOptions) {
			createOpts.Label = "linodego-fw-ins-test"
		},
	}, "fixtures/TestNodeBalancerFirewalls_List")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	result, err := client.ListNodeBalancerFirewalls(context.Background(), nil, nodebalancer.ID)
	if err != nil {
		t.Errorf("Error listing Firewalls, expected struct, got error %v", err)
	}

	if len(result) == 0 {
		t.Errorf("Expected a list of Firewalls, but got none: %v", err)
	}
}

func setupNodeBalancerFirewall(t *testing.T, firewallModifiers []firewallModifier, fixturesYaml string) (*linodego.Client, *linodego.NodeBalancer, *linodego.Firewall, func(), error) {
	t.Helper()
	client, nodebalancer, nodebalancerTeardown, err := setupNodeBalancer(t, fixturesYaml)
	device := linodego.DevicesCreationOptions{NodeBalancers: []int{nodebalancer.ID}}
	firewallModifiers = append(firewallModifiers,
		func(createOpts *linodego.FirewallCreateOptions) {
			createOpts.Devices = device
		})
	firewall, firewallTeardown, err := createFirewall(t, client, firewallModifiers...)

	teardown := func() {
		nodebalancerTeardown()
		firewallTeardown()
	}
	return client, nodebalancer, firewall, teardown, err
}
