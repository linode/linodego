package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/linode/linodego"
)

func TestGetIPAddress_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetIPAddress_missing")
	defer teardown()

	doesNotExist := "010.020.030.040"
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

func TestGetIPAddress_found(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestGetIPAddress_found")
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

func TestListIPAddresses(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestListIPAddresses")
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

func TestGetInstanceIPAddresses(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestGetInstanceIPAddresses")
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

func TestUpdateIPAddress(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestUpdateIPAddress")
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

// TestDeleteInstanceIPAddress requires the customer account to have
// default_IPMax set to at least 2 and default_InterfaceMax set to 3.
func TestDeleteInstanceIPAddress(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestDeleteInstanceIPAddress")
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

func TestAssignInstancesIPs(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestAssignInstancesIPs")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	newInstance, err := createInstance(t, client, func(options *InstanceCreateOptions) {
		options.Label = fmt.Sprintf("linodego-%d", time.Now().UnixNano())
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
