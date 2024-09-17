package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestInstance_CreateWithReservedIPAddress(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestInstance_CreateWithReservedIPAddress")
	defer teardown()

	// Reserve an IP for testing
	reservedIP, err := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{Region: "us-east"})
	if err != nil {
		t.Fatalf("Failed to reserve IP: %v", err)
	}
	defer func() {
		err := client.DeleteReservedIPAddress(context.Background(), reservedIP.Address)
		if err != nil {
			t.Errorf("Failed to delete reserved IP: %v", err)
		}
	}()

	_, instanceTeardown, err := createInstanceWithReservedIP(t, client, reservedIP.Address)
	defer instanceTeardown()
	if err != nil {
		t.Fatalf("Error creating instance with reserved IP: %s", err)
	}

}

func createInstanceWithReservedIP(
	t *testing.T,
	client *linodego.Client,
	reservedIP string,
	modifiers ...instanceModifier,
) (*linodego.Instance, func(), error) {
	t.Helper()

	createOpts := linodego.InstanceCreateOptions{
		Label:    "go-test-ins-reserved-ip-" + randLabel(),
		Region:   "us-east",
		Type:     "g6-nanode-1",
		Booted:   linodego.Pointer(false),
		Image:    "linode/alpine3.17",
		RootPass: randPassword(),
		Interfaces: []linodego.InstanceConfigInterfaceCreateOptions{
			{
				Purpose:     linodego.InterfacePurposePublic,
				Label:       "",
				IPAMAddress: "",
			},
		},
		Ipv4: []string{reservedIP},
	}

	for _, modifier := range modifiers {
		modifier(client, &createOpts)
	}

	instance, err := client.CreateInstance(context.Background(), createOpts)
	if err != nil {
		return nil, func() {}, err
	}

	teardown := func() {
		if terr := client.DeleteInstance(context.Background(), instance.ID); terr != nil {
			t.Errorf("Error deleting test Instance: %s", terr)
		}
	}

	return instance, teardown, nil
}

func TestInstance_CreateWithOwnedNonAssignedReservedIP(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestInstance_CreateWithOwnedNonAssignedReservedIP")
	defer teardown()

	reservedIP, err := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{Region: "us-east"})
	if err != nil {
		t.Fatalf("Failed to reserve IP: %v", err)
	}
	defer func() {
		err := client.DeleteReservedIPAddress(context.Background(), reservedIP.Address)
		if err != nil {
			t.Errorf("Failed to delete reserved IP: %v", err)
		}
	}()

	_, instanceTeardown, err := createInstanceWithReservedIP(t, client, reservedIP.Address)
	defer instanceTeardown()
	if err != nil {
		t.Errorf("Unexpected error with owned non-assigned reserved IP: %v", err)
	}
}

func TestInstance_CreateWithAlreadyAssignedReservedIP(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestInstance_CreateWithAlreadyAssignedReservedIP")
	defer teardown()

	reservedIP, err := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{Region: "us-east"})
	if err != nil {
		t.Fatalf("Failed to reserve IP: %v", err)
	}
	defer func() {
		err := client.DeleteReservedIPAddress(context.Background(), reservedIP.Address)
		if err != nil {
			t.Errorf("Failed to delete reserved IP: %v", err)
		}
	}()

	// First, create an instance with the reserved IP
	_, instanceTeardown, err := createInstanceWithReservedIP(t, client, reservedIP.Address)
	defer instanceTeardown()
	if err != nil {
		t.Fatalf("Failed to create initial instance: %v", err)
	}

	// Now try to create another instance with the same IP
	_, secondInstanceTeardown, err := createInstanceWithReservedIP(t, client, reservedIP.Address)
	defer secondInstanceTeardown()
	if err == nil {
		t.Errorf("Expected error with already assigned reserved IP, but got none")
	}
}

func TestInstance_CreateWithNonReservedAddress(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestInstance_CreateWithNonReservedAddress")
	defer teardown()

	_, instanceTeardown, err := createInstanceWithReservedIP(t, client, "192.0.2.1")
	defer instanceTeardown()
	if err == nil {
		t.Errorf("Expected error with non-reserved address, but got none")
	}
}

func TestInstance_CreateWithNonOwnedReservedAddress(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestInstance_CreateWithNonOwnedReservedAddress")
	defer teardown()

	_, instanceTeardown, err := createInstanceWithReservedIP(t, client, "198.51.100.1")
	defer instanceTeardown()
	if err == nil {
		t.Errorf("Expected error with non-owned reserved address, but got none")
	}
}

func TestInstance_CreateWithEmptyIPAddress(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestInstance_CreateWithEmptyIPAddress")
	defer teardown()

	_, instanceTeardown, err := createInstanceWithReservedIP(t, client, "")
	defer instanceTeardown()
	if err == nil {
		t.Errorf("Expected error with empty IP address, but got none")
	}
}

func TestInstance_CreateWithNullIPAddress(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestInstance_CreateWithNullIPAddress")
	defer teardown()

	_, instanceTeardown, err := createInstanceWithReservedIP(t, client, "", func(client *linodego.Client, opts *linodego.InstanceCreateOptions) {
		opts.Ipv4 = nil
	})
	defer instanceTeardown()
	if err != nil {
		t.Errorf("Unexpected error with null IP address: %v", err)
	}
}

func TestInstance_CreateWithMultipleIPAddresses(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestInstance_CreateWithMultipleIPAddresses")
	defer teardown()

	reservedIP, err := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{Region: "us-east"})
	if err != nil {
		t.Fatalf("Failed to reserve IP: %v", err)
	}
	defer func() {
		err := client.DeleteReservedIPAddress(context.Background(), reservedIP.Address)
		if err != nil {
			t.Errorf("Failed to delete reserved IP: %v", err)
		}
	}()

	_, instanceTeardown, err := createInstanceWithReservedIP(t, client, "", func(client *linodego.Client, opts *linodego.InstanceCreateOptions) {
		opts.Ipv4 = []string{reservedIP.Address, "192.0.2.2"}
	})
	defer instanceTeardown()
	if err == nil {
		t.Errorf("Expected error with multiple IP addresses, but got none")
	}
}

func TestInstance_CreateWithoutIPv4Field(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestInstance_CreateWithoutIPv4Field")
	defer teardown()

	_, instanceTeardown, err := createInstanceWithReservedIP(t, client, "", func(client *linodego.Client, opts *linodego.InstanceCreateOptions) {
		opts.Ipv4 = nil
	})
	defer instanceTeardown()
	if err != nil {
		t.Errorf("Unexpected error when omitting IPv4 field: %v", err)
	}
}

func TestInstance_AddReservedIPToInstance(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestInstance_AddReservedIPToInstance")
	defer teardown()

	// Create a test Linode instance
	instance, err := client.CreateInstance(context.Background(), linodego.InstanceCreateOptions{
		Region:   "us-east",
		Type:     "g6-nanode-1",
		Label:    "test-instance-for-ip-reservation",
		RootPass: randPassword(),
	})
	if err != nil {
		t.Fatalf("Error creating test instance: %v", err)
	}
	defer func() {
		if err := client.DeleteInstance(context.Background(), instance.ID); err != nil {
			t.Errorf("Error deleting test instance: %v", err)
		}
	}()

	// Reserve an IP address
	reservedIP, err := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{
		Region: "us-east",
	})
	if err != nil {
		t.Fatalf("Error reserving IP address: %v", err)
	}
	defer func() {
		if err := client.DeleteReservedIPAddress(context.Background(), reservedIP.Address); err != nil {
			t.Errorf("Error deleting reserved IP: %v", err)
		}
	}()

	// Add the reserved IP to the instance
	opts := linodego.InstanceReserveIPOptions{
		Type:    "ipv4",
		Public:  true,
		Address: reservedIP.Address,
	}
	_, err = client.AddReservedIPToInstance(context.Background(), instance.ID, opts)
	if err != nil {
		t.Fatalf("Error adding reserved IP to instance: %v", err)
	}

	// Verify the IP was added to the instance
	ips, err := client.GetInstanceIPAddresses(context.Background(), instance.ID)
	if err != nil {
		t.Fatalf("Error getting instance IP addresses: %v", err)
	}

	found := false
	for _, ip := range ips.IPv4.Public {
		if ip.Address == reservedIP.Address {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Reserved IP %s was not found in instance's IP addresses", reservedIP.Address)
	}
}

func TestInstance_AddReservedIPToInstanceVariants(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestInstance_AddReservedIPToInstanceVariants")
	defer teardown()

	// Create a test Linode instance
	instance, err := client.CreateInstance(context.Background(), linodego.InstanceCreateOptions{
		Region:   "us-east",
		Type:     "g6-nanode-1",
		Label:    "test-instance-for-ip-reservation",
		RootPass: randPassword(),
	})
	if err != nil {
		t.Fatalf("Error creating test instance: %v", err)
	}
	defer func() {
		if err := client.DeleteInstance(context.Background(), instance.ID); err != nil {
			t.Errorf("Error deleting test instance: %v", err)
		}
	}()

	// Reserve an IP address
	reservedIP, err := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{
		Region: "us-east",
	})
	if err != nil {
		t.Fatalf("Error reserving IP address: %v", err)
	}
	defer func() {
		if err := client.DeleteReservedIPAddress(context.Background(), reservedIP.Address); err != nil {
			t.Errorf("Error deleting reserved IP: %v", err)
		}
	}()

	// Test: Add reserved IP to instance with valid parameters
	opts := linodego.InstanceReserveIPOptions{
		Type:    "ipv4",
		Public:  true,
		Address: reservedIP.Address,
	}
	_, err = client.AddReservedIPToInstance(context.Background(), instance.ID, opts)
	if err != nil {
		t.Fatalf("Error adding reserved IP to instance: %v", err)
	}

	// Test: Omit public field
	omitPublicOpts := linodego.InstanceReserveIPOptions{
		Type:    "ipv4",
		Address: reservedIP.Address,
		// Public field is omitted here
	}
	_, err = client.AddReservedIPToInstance(context.Background(), instance.ID, omitPublicOpts)
	if err == nil {
		t.Fatalf("Expected error when adding reserved IP with omitted public field, but got none")
	}

	// Assume we have a Linode that has been created without a reserved IP address and IPMAX set to 1
	linodeID := 63510870

	// Reserve IP address
	resIP, err := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{
		Region: "us-east",
	})
	if err != nil {
		t.Fatalf("Failed to reserve IP: %v", err)
	}

	//  Add IP address to the Linode
	_, err = client.AddReservedIPToInstance(context.Background(), linodeID, linodego.InstanceReserveIPOptions{
		Type:    "ipv4",
		Public:  true,
		Address: resIP.Address,
	})
	if err == nil {
		t.Errorf("Expected error when adding reserved IP to a Linode at its IPMAX limit, but got none")
	}

	// Delete the reserved IP Address

	if err := client.DeleteReservedIPAddress(context.Background(), resIP.Address); err != nil {
		t.Errorf("Failed to delete first reserved IP: %v", err)
	}

	// Test: Non-owned Linode ID
	nonOwnedInstanceID := 888888 // Replace with an actual non-owned Linode ID
	_, err = client.AddReservedIPToInstance(context.Background(), nonOwnedInstanceID, opts)
	if err == nil {
		t.Errorf("Expected error when adding reserved IP to non-owned Linode, but got none")
	}

	// Test: Already assigned reserved IP
	_, err = client.AddReservedIPToInstance(context.Background(), instance.ID, opts)
	if err == nil {
		t.Errorf("Expected error when adding already assigned reserved IP, but got none")
	}

	// Test: Non-owned reserved IP
	_, err = client.AddReservedIPToInstance(context.Background(), instance.ID, linodego.InstanceReserveIPOptions{
		Type:    "ipv4",
		Public:  true,
		Address: "198.51.100.1", // Assume this is a non-owned reserved IP
	})
	if err == nil {
		t.Errorf("Expected error when adding non-owned reserved IP, but got none")
	}

	// Test: Reserved IP in different datacenter
	// Reserve an IP address
	diffDataCentreIP, err := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{
		Region: "ca-central",
	})
	if err != nil {
		t.Fatalf("Error reserving IP address: %v", err)
	}
	defer func() {
		if err := client.DeleteReservedIPAddress(context.Background(), diffDataCentreIP.Address); err != nil {
			t.Errorf("Error deleting reserved IP: %v", err)
		}
	}()
	_, err = client.AddReservedIPToInstance(context.Background(), instance.ID, linodego.InstanceReserveIPOptions{
		Type:    "ipv4",
		Public:  true,
		Address: diffDataCentreIP.Address, // Assume this IP is in a different datacenter
	})
	if err == nil {
		t.Errorf("Expected error when adding reserved IP in different datacenter, but got none")
	}

	// Test: IPv6 type
	_, err = client.AddReservedIPToInstance(context.Background(), instance.ID, linodego.InstanceReserveIPOptions{
		Type:    "ipv6",
		Public:  true,
		Address: reservedIP.Address,
	})
	if err == nil {
		t.Errorf("Expected error when adding reserved IP with type ipv6, but got none")
	}

	// Test: Public field set to false
	opts.Public = false
	_, err = client.AddReservedIPToInstance(context.Background(), instance.ID, opts)
	if err == nil {
		t.Errorf("Expected error when adding reserved IP with public field set to false, but got none")
	}

	// Test: Integer as address
	_, err = client.AddReservedIPToInstance(context.Background(), instance.ID, linodego.InstanceReserveIPOptions{
		Type:    "ipv4",
		Public:  true,
		Address: "12345", // Invalid IP format
	})
	if err == nil {
		t.Errorf("Expected error when adding reserved IP with integer as address, but got none")
	}

	// Test: Empty address
	_, err = client.AddReservedIPToInstance(context.Background(), instance.ID, linodego.InstanceReserveIPOptions{
		Type:    "ipv4",
		Public:  true,
		Address: "",
	})
	if err == nil {
		t.Errorf("Expected error when adding reserved IP with empty address, but got none")
	}

	// Test: Null address
	_, err = client.AddReservedIPToInstance(context.Background(), instance.ID, linodego.InstanceReserveIPOptions{
		Type:   "ipv4",
		Public: true,
	})
	if err == nil {
		t.Errorf("Expected error when adding reserved IP with null address, but got none")
	}

	// Test: Omit address field
	_, err = client.AddReservedIPToInstance(context.Background(), instance.ID, linodego.InstanceReserveIPOptions{
		Type:   "ipv4",
		Public: true,
	})
	if err == nil {
		t.Errorf("Expected error when omitting address field, but got none")
	}
}
