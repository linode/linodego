package integration

import (
	"context"
	"fmt"
	"testing"

	. "github.com/linode/linodego"
)

func TestIPAddress_GetMissing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestIPAddress_GetMissing")
	defer teardown()

	doesNotExist := "10.0.0.1"
	i, err := client.GetIPAddress(context.Background(), doesNotExist)
	if err == nil {
		t.Errorf("should have received an error requesting a missing ipaddress, got %v", i)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing ipaddress, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing ipaddress, got %v", e.Code)
	}
}

func TestIPAddress_GetFound(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestIPAddress_GetFound")
	defer teardown()
	if err != nil {
		t.Errorf("Error creating IPAddress test Instance, got error %v", err)
	}

	address := instance.IPv4[0].String()
	i, err := client.GetIPAddress(context.Background(), address)
	if err != nil {
		t.Errorf("Error getting ipaddress, expected struct, got %v and error %v", i, err)
	}
	if i.Address != address {
		t.Errorf("Expected a specific ipaddress, but got a different one %v", i)
	}
}

func TestIPAddresses_List(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestIPAddresses_List")
	defer teardown()
	if err != nil {
		t.Errorf("Error creating IPAddress test Instance, got error %v", err)
	}

	filter := fmt.Sprintf("{\"linode_id\":%d}", instance.ID)
	i, err := client.ListIPAddresses(context.Background(), NewListOptions(0, filter))
	if err != nil {
		t.Errorf("Error listing ipaddresses, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of ipaddresses, but got none %v", i)
	}
}

func TestIPAddresses_Instance_Get(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestIPAddresses_Instance_Get")
	defer teardown()
	if err != nil {
		t.Errorf("Error creating IPAddress test Instance, got error %v", err)
	}

	i, err := client.GetInstanceIPAddresses(context.Background(), instance.ID)
	if err != nil {
		t.Errorf("Error listing ipaddresses, expected struct, got error %v", err)
	}
	if i.IPv4.Public[0].Address != instance.IPv4[0].String() {
		t.Errorf("Expected matching ipaddresses with GetInstanceIPAddress Instance IPAddress but got %v", i)
	}
}

func TestIPAddress_Update(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestIPAddress_Update")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	address := instance.IPv4[0].String()
	i, err := client.GetInstanceIPAddresses(context.Background(), instance.ID)
	if err != nil {
		t.Errorf("Error getting ipaddress: %s", err)
	}
	rdns := i.IPv4.Public[0].RDNS

	updateOpts := IPAddressUpdateOptions{
		RDNS: &rdns,
	}

	_, err = client.UpdateIPAddress(context.Background(), address, updateOpts)
	if err != nil {
		t.Error(err)
	}
}

// TestIPAddress_Instance_Delete requires the customer account to have
// default_IPMax set to at least 2 and default_InterfaceMax set to 3.
func TestIPAddress_Instance_Delete(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestIPAddress_Instance_Delete")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	ip, err := client.AddInstanceIPAddress(context.TODO(), instance.ID, true)
	if err != nil {
		t.Fatalf("failed to allocate public IPv4 for instance (%d): %s", instance.ID, err)
	}

	i, err := client.GetInstanceIPAddresses(context.TODO(), instance.ID)
	if err != nil {
		t.Fatalf("failed to get instance (%d) IP addresses: %s", instance.ID, err)
	}
	if len(i.IPv4.Public) != 2 {
		t.Errorf("expected instance (%d) to have 2 public IPv4 addresses; got %d", instance.ID, len(i.IPv4.Public))
	}

	if err := client.DeleteInstanceIPAddress(context.TODO(), instance.ID, ip.Address); err != nil {
		t.Fatalf("failed to delete instance (%d) public IPv4 address (%s): %s", instance.ID, ip.Address, err)
	}

	if i, err = client.GetInstanceIPAddresses(context.TODO(), instance.ID); err != nil {
		t.Fatalf("failed to get instance (%d) IP addresses: %s", instance.ID, err)
	}
	if len(i.IPv4.Public) != 1 {
		t.Errorf("expected instance (%d) to have 1 public IPv4 address; got %d", instance.ID, len(i.IPv4.Public))
	}
}

func TestIPAddress_Instance_Assign(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestIPAddress_Instance_Assign")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	newInstance, err := createInstance(t, client, func(client *Client, options *InstanceCreateOptions) {
		options.Label = "go-ins-test-assign"
		options.Region = instance.Region
	})

	defer func() {
		if err := client.DeleteInstance(context.Background(), newInstance.ID); err != nil {
			if t != nil {
				t.Errorf("Error deleting test Instance: %s", err)
			}
		}
	}()

	if err != nil {
		t.Error(err)
	}

	ipRange, err := client.CreateIPv6Range(context.Background(), IPv6RangeCreateOptions{
		PrefixLength: 64,
		LinodeID:     newInstance.ID,
	})
	if err != nil {
		t.Error(err)
	}

	// IP reassignment
	err = client.InstancesAssignIPs(context.Background(), LinodesAssignIPsOptions{
		Region: instance.Region,
		Assignments: []LinodeIPAssignment{
			{
				LinodeID: instance.ID,
				Address:  ipRange.Range,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}

	ips, err := client.GetInstanceIPAddresses(context.Background(), instance.ID)
	if err != nil {
		t.Error(err)
	}

	for _, r := range ips.IPv6.Global {
		if fmt.Sprintf("%s/%d", r.Range, r.Prefix) == ipRange.Range {
			return
		}
	}

	t.Errorf("failed to find assigned ip")
}

func TestIPAddress_Instance_Share(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestIPAddress_Instance_Share", func(client *Client, options *InstanceCreateOptions) {
		// This should stay hardcoded at the moment as the
		// IP sharing rollout does not have a corresponding capability.
		options.Region = "us-west"
	})
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	newInstance, err := createInstance(t, client, func(client *Client, options *InstanceCreateOptions) {
		options.Label = "go-ins-test-share"
		options.Region = instance.Region
	})

	defer func() {
		if err := client.DeleteInstance(context.Background(), newInstance.ID); err != nil {
			if t != nil {
				t.Errorf("Error deleting test Instance: %s", err)
			}
		}
	}()

	if err != nil {
		t.Error(err)
	}

	ip, err := client.AddInstanceIPAddress(context.Background(), newInstance.ID, true)
	if err != nil {
		t.Error(err)
	}

	// IP sharing
	err = client.ShareIPAddresses(context.Background(), IPAddressesShareOptions{
		LinodeID: instance.ID,
		IPs: []string{
			ip.Address,
		},
	})
	if err != nil {
		t.Error(err)
	}

	ips, err := client.GetInstanceIPAddresses(context.Background(), instance.ID)
	if err != nil {
		t.Error(err)
	}

	for _, r := range ips.IPv4.Shared {
		if r.Address == ip.Address {
			return
		}
	}

	t.Errorf("failed to find assigned ip")
}
