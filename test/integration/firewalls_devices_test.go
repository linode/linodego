package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/linode/linodego"
	"github.com/linode/linodego/pkg/errors"
)

func TestListFirewallDevices(t *testing.T) {
	client, instance, teardown, err := setupInstance(t, "fixtures/TestListFirewallDevices")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	firewall, teardownFirewall, err := createFirewall(t, client, func(opts *linodego.FirewallCreateOptions) {
		opts.Devices.Linodes = []int{instance.ID}
	})
	if err != nil {
		t.Error(err)
	}
	defer teardownFirewall()

	firewallDevices, err := client.ListFirewallDevices(context.Background(), firewall.ID, nil)
	if err != nil {
		t.Error(err)
	}

	if len(firewallDevices) != 1 {
		t.Errorf("expected 1 firewall device but got %d", len(firewallDevices))
	}
}

func TestGetFirewallDevice(t *testing.T) {
	client, instance, teardown, err := setupInstance(t, "fixtures/TestGetFirewallDevice")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	firewall, teardownFirewall, err := createFirewall(t, client)
	if err != nil {
		t.Error(err)
	}
	defer teardownFirewall()

	firewallDevice, err := client.CreateFirewallDevice(context.Background(), firewall.ID, linodego.FirewallDeviceCreateOptions{
		Type: linodego.FirewallDeviceLinode,
		ID:   instance.ID,
	})
	if err != nil {
		t.Error(err)
	}

	if device, err := client.GetFirewallDevice(context.Background(), firewall.ID, firewallDevice.ID); err != nil {
		t.Error(err)
	} else if !cmp.Equal(device, firewallDevice) {
		t.Errorf("expected device to match create result but got diffs: %s", cmp.Diff(device, firewallDevice))
	}
}

func TestDeleteFirewallDevice(t *testing.T) {
	client, instance, teardown, err := setupInstance(t, "fixtures/TestDeleteFirewallDevice")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	firewall, teardownFirewall, err := createFirewall(t, client)
	if err != nil {
		t.Error(err)
	}
	defer teardownFirewall()

	firewallDevice, err := client.CreateFirewallDevice(context.Background(), firewall.ID, linodego.FirewallDeviceCreateOptions{
		Type: linodego.FirewallDeviceLinode,
		ID:   instance.ID,
	})
	if err != nil {
		t.Error(err)
	}

	assertDateSet(t, firewallDevice.Created)
	assertDateSet(t, firewallDevice.Updated)

	if err := client.DeleteFirewallDevice(context.Background(), firewall.ID, firewallDevice.ID); err != nil {
		t.Error(err)
	}

	if _, getErr := client.GetFirewallDevice(context.Background(), firewall.ID, firewallDevice.ID); err != nil {
		t.Error("expected fetching firewall device to fail")
	} else if apiError, ok := getErr.(*errors.Error); !ok || apiError.Code != http.StatusNotFound {
		t.Errorf("expected fetching firewall device to throw Not Found but got: %s", getErr)
	}
}
