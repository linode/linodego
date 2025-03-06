package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/linode/linodego"
	. "github.com/linode/linodego"
)

// TestReservedIPAddresses_InsufficientPermissions tests the behavior when a user account
// doesn't have the permission to use the Reserved IP feature
func TestReservedIPAddresses_InsufficientPermissions(t *testing.T) {
	original := validTestAPIKey
	dummyToken := "badtoken"
	validTestAPIKey = dummyToken

	client, teardown := createTestClient(t, "fixtures/TestReservedIPAddresses_InsufficientPermissions")
	defer teardown()
	defer func() { validTestAPIKey = original }()

	filter := ""
	ips, listErr := client.ListReservedIPAddresses(context.Background(), NewListOptions(0, filter))
	if listErr == nil {
		t.Errorf("Expected error due to insufficient permissions, but got none %v", ips)
	} else {
		t.Logf("Correctly received error when listing IP addresses: %v", listErr)
	}

	if len(ips) != 0 {
		t.Errorf("Expected no IP addresses due to insufficient permissions, but got some: %v", ips)
	}

	// Attempt to reserve an IP address
	resIP, resErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{
		Region: "us-east",
	})
	if resErr == nil {
		t.Errorf("Expected error when reserving IP due to insufficient permissions, but got none")
	} else {
		t.Logf("Correctly received %v and error when reserving IP: %v", resIP, resErr)
	}

	// Attempt to get a reserved IP address
	address := "172.28.3.4"
	ip, getErr := client.GetReservedIPAddress(context.Background(), address)
	if getErr == nil {
		t.Errorf("Expected error when getting IP address due to insufficient permissions, but got none")
	} else {
		t.Logf("Correctly received %v for IP Address and error when getting IP address: %v", ip, getErr)
	}

	// Attempt to delete a reserved IP address
	delAddr := "172.28.3.4"
	delErr := client.DeleteReservedIPAddress(context.Background(), delAddr)
	if delErr == nil {
		t.Errorf("Expected error when deleting IP address due to insufficient permissions, but got none")
	} else {
		t.Logf("Correctly received error when deleting IP address: %v", delErr)
	}
}

// TestReservedIPAddresses_EndToEndTest performs an end-to-end test of the Reserved IP functionality
// for users with the can_reserve_ip flag enabled
func TestReservedIPAddresses_EndToEndTest(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestReservedIPAddresses_EndToEndTest")
	defer teardown()

	filter := ""

	ipList, err := client.ListReservedIPAddresses(context.Background(), NewListOptions(0, filter))
	if err != nil {
		t.Fatalf("Error listing IP addresses: %v", err)
	}

	initialCount := len(ipList)

	// Attempt to reserve an IP
	resIP, resErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{
		Region: "us-east",
	})

	if resErr != nil {
		t.Fatalf("Failed to reserve IP. This test expects the user to have 0 prior reservations and the ip_reservation_limit to be 2. Error from the API: %v", resErr)
	}

	t.Logf("Successfully reserved IP: %+v", resIP)

	// Fetch the reserved IP
	fetchedIP, fetchErr := client.GetReservedIPAddress(context.Background(), resIP.Address)
	if fetchErr != nil {
		t.Errorf("Error getting reserved IP address: %v", fetchErr)
	}

	if fetchedIP == nil {
		t.Errorf("Expected %s but got nil indicating a failure in fetching the reserved IP", resIP.Address)
	}

	// Verify the list of IPs has increased
	verifyList, verifyErr := client.ListReservedIPAddresses(context.Background(), NewListOptions(0, filter))
	if verifyErr != nil {
		t.Fatalf("Error listing IP addresses after reservation: %v", verifyErr)
	}

	if len(verifyList) != initialCount+1 {
		t.Errorf("Expected IP count to increase by 1, got %d, want %d", len(verifyList), initialCount+1)
	}

	// Delete the reserved IP
	delErr := client.DeleteReservedIPAddress(context.Background(), resIP.Address)
	if delErr != nil {
		t.Fatalf("Error deleting reserved IP address: %v", delErr)
	}

	// Verify the IP has been deleted
	_, fetchDelErr := client.GetReservedIPAddress(context.Background(), resIP.Address)
	if fetchDelErr == nil {
		t.Errorf("Expected error when fetching %s, got nil", resIP.Address)
	}

	verifyDelList, verifyDelErr := client.ListReservedIPAddresses(context.Background(), NewListOptions(0, filter))
	if verifyDelErr != nil {
		t.Fatalf("Error listing IP addresses after deletion: %v", verifyDelErr)
	}

	if len(verifyDelList) != initialCount {
		t.Errorf("Expected IP count to return to initial count, got %d, want %d", len(verifyDelList), initialCount)
	}
}

// TestReservedIPAddresses_ListIPAddressesVariants tests filters for listing IP addresses
func TestReservedIPAddresses_ListIPAddressesVariants(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestReservedIPAddresses_ListIPAddressesVariants")
	defer teardown()

	expected_ips := 2

	// Reserve two IP addresses in us-east region
	reservedIPs := make([]string, expected_ips)
	for i := 0; i < expected_ips; i++ {
		reserveIP, err := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{
			Region: "us-east",
		})
		if err != nil {
			t.Fatalf("Failed to reserve IP %d: %v", i+1, err)
		}
		reservedIPs[i] = reserveIP.Address
		t.Logf("Successfully reserved IP %d: %s", i+1, reserveIP.Address)
	}

	// Defer cleanup of reserved IPs
	defer func() {
		for _, ip := range reservedIPs {
			err := client.DeleteReservedIPAddress(context.Background(), ip)
			if err != nil {
				t.Errorf("Failed to delete reserved IP %s: %v", ip, err)
			}
		}
	}()

	// Create ListOptions with the filter for reserved IPs in us-east region
	listOptions := linodego.ListOptions{
		PageOptions: &linodego.PageOptions{
			Page: 0,
		},
		Filter: "{\"reserved\":true,\"region\":\"us-east\"}",
	}

	ipList, err := client.ListIPAddresses(context.Background(), &listOptions)
	if err != nil {
		t.Fatalf("Error listing reserved IP addresses in us-east: %v", err)
	}

	t.Logf("Retrieved %d reserved IP addresses in us-east", len(ipList))

	// Check if at least the two reserved IPs are in the list
	foundReservedIPs := 0
	for _, ip := range ipList {
		if !ip.Reserved {
			t.Errorf("Expected all IPs to be reserved, but found non-reserved IP: %s", ip.Address)
		}
		if ip.Region != "us-east" {
			t.Errorf("Expected all IPs to be in us-east region, but found IP in %s region: %s", ip.Region, ip.Address)
		}
		for _, reservedIP := range reservedIPs {
			if ip.Address == reservedIP {
				foundReservedIPs++
				break
			}
		}
	}
	if foundReservedIPs != expected_ips {
		t.Errorf("Expected %d but found %d while listing reserved IP addresses", expected_ips, foundReservedIPs)
	}
}

// TestReservedIPAddresses_GetIPAddressVariants tests various scenarios for getting a specific IP address
func TestReservedIPAddresses_GetIPAddressVariants(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestReservedIPAddresses_GetIPAddressVariants")
	defer teardown()

	// Reserve an IP for testing
	resIP, resErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{
		Region: "us-east",
	})

	if resErr != nil {
		t.Fatalf("Failed to reserve IP. This test expects the user to have 0 prior reservations and the ip_reservation_limit to be 2. Error from the API: %v", resErr)
	}

	if resIP == nil {
		t.Fatalf("Reserved IP is nil")
	}

	t.Logf("Successfully reserved IP: %+v", resIP)

	// Test getting a valid reserved IP
	validIP, fetchErr := client.GetReservedIPAddress(context.Background(), resIP.Address)
	if fetchErr != nil {
		t.Errorf("Error getting valid reserved IP address: %v", fetchErr)
	}

	if validIP == nil {
		t.Errorf("Retrieved valid reserved IP is nil")
	} else {
		if validIP.Address != resIP.Address {
			t.Errorf("Retrieved IP address does not match reserved IP address. Got %s, want %s", validIP.Address, resIP.Address)
		}
	}

	// Test getting an invalid IP
	invalidIP := "999.999.999.999"
	_, invalidFetchErr := client.GetReservedIPAddress(context.Background(), invalidIP)
	if invalidFetchErr == nil {
		t.Errorf("Expected error when fetching invalid IP, got nil")
	}

	// Clean up: Delete the reserved IP
	delErr := client.DeleteReservedIPAddress(context.Background(), resIP.Address)
	if delErr != nil {
		t.Errorf("Failed to delete reserved IP: %v", delErr)
	}
}

// TestReservedIPAddresses_ReserveIPVariants tests various scenarios for reserving an IP address
func TestReservedIPAddresses_ReserveIPAddressVariants(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestReservedIPAddresses_ReserveIPVariants")
	defer teardown()

	// Slice to keep track of all reserved IPs
	var reservedIPs []string

	// Helper function to clean up reserved IPs
	cleanupIPs := func() {
		for _, ip := range reservedIPs {
			err := client.DeleteReservedIPAddress(context.Background(), ip)
			if err != nil {
				t.Errorf("Failed to delete reserved IP %s: %v", ip, err)
			}
		}
	}
	defer cleanupIPs()

	// Test reserving IP with omitted region
	_, omitErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{})
	if omitErr == nil {
		t.Errorf("Expected error when reserving IP with omitted region, got nil")
	}

	// Test reserving IP with invalid region
	_, invalidErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{Region: "us"})
	if invalidErr == nil {
		t.Errorf("Expected error when reserving IP with invalid region, got nil")
	}

	// Test reserving IP with empty region
	_, emptyErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{Region: ""})
	if emptyErr == nil {
		t.Errorf("Expected error when reserving IP with empty region, got nil")
	}

	// Make 2 valid IP Reservations
	for i := 0; i < 2; i++ {
		reserveIP, err := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{
			Region: "us-east",
		})
		if err != nil {
			t.Fatalf("Failed to reserve IP %d: %v", i+1, err)
		}
		reservedIPs = append(reservedIPs, reserveIP.Address)
		t.Logf("Successfully reserved IP %d: %s", i+1, reserveIP.Address)
	}
}

func TestReservedIPAddresses_ExceedLimit(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestReservedIPAddresses_ExceedLimit")
	defer teardown()

	// Slice to keep track of all reserved IPs
	var reservedIPs []string

	// Helper function to clean up reserved IPs
	cleanupIPs := func() {
		for _, ip := range reservedIPs {
			err := client.DeleteReservedIPAddress(context.Background(), ip)
			if err != nil {
				t.Errorf("Failed to delete reserved IP %s: %v", ip, err)
			}
		}
	}
	defer cleanupIPs()

	// Reserve IPs until the limit is reached and assert the error message
	for i := 0; i < 100; i++ {
		reservedIP, err := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{
			Region: "us-east",
		})
		if err != nil {
			expectedErrorMessage := "[400] Additional Reserved IPv4 addresses require technical justification."
			if !strings.Contains(err.Error(), expectedErrorMessage) {
				t.Errorf("Expected error message to contain '%s', but got: %v", expectedErrorMessage, err)
			} else {
				t.Logf("Failed to reserve IP %d as expected: %v", i+1, err)
			}
			break
		}

		reservedIPs = append(reservedIPs, reservedIP.Address)

		if i == 99 {
			t.Errorf("Expected to hit reservation limit, but did not reach it after 100 attempts")
		}
	}
}

// TestReservedIPAddresses_DeleteIPAddressVariants tests various scenarios for deleting a reserved IP address
func TestReservedIPAddresses_DeleteIPAddressVariants(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestReservedIPAddresses_DeleteIPAddressVariants")
	defer teardown()

	validRes, validErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{Region: "us-east"})
	if validErr != nil {
		t.Fatalf("Failed to reserve IP. This test should start with 0 reservations or reservations < limit. Error from the API: %v", validErr)
	}

	if validRes == nil {
		t.Fatalf("Valid reservation returned nil IP")
	}

	t.Logf("Successfully reserved IP: %+v", validRes)

	filter := ""
	ipList, listErr := client.ListReservedIPAddresses(context.Background(), NewListOptions(0, filter))
	if listErr != nil {
		t.Fatalf("Error listing IP addresses: %v", listErr)
	}

	if len(ipList) == 0 {
		t.Fatalf("No reserved IPs available for testing deletion")
	}

	delErr := client.DeleteReservedIPAddress(context.Background(), validRes.Address)
	if delErr != nil {
		t.Fatalf("Failed to delete reserved IP address: %v", delErr)
	}

	// Verify deletion
	verifyDelList, verifyDelErr := client.ListReservedIPAddresses(context.Background(), NewListOptions(0, filter))
	if verifyDelErr != nil {
		t.Fatalf("Error listing IP addresses after deletion: %v", verifyDelErr)
	}

	if len(verifyDelList) >= len(ipList) {
		t.Errorf("IP address deletion not confirmed. Expected count < %d, got %d", len(ipList), len(verifyDelList))
	}

	_, fetchDelErr := client.GetReservedIPAddress(context.Background(), validRes.Address)
	if fetchDelErr == nil {
		t.Errorf("Expected error when fetching deleted IP, got nil")
	}

	// Test deleting an unowned IP
	unownedIP := "255.255.255.4"
	delUnownedErr := client.DeleteReservedIPAddress(context.Background(), unownedIP)
	if delUnownedErr == nil {
		t.Errorf("Expected error when deleting unowned IP, got nil")
	}
}

func TestReservedIPAddresses_GetIPReservationStatus(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestReservedIPAddresses_GetInstanceIPReservationStatus")
	defer teardown()

	// Create a Linode with a reserved IP
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

	instanceWithReservedIP, instanceTeardown, err := createInstanceWithReservedIP(t, client, reservedIP.Address)
	if err != nil {
		t.Fatalf("Error creating instance with reserved IP: %s", err)
	}
	defer instanceTeardown()

	// Make GET request for the Linode with reserved IP
	instanceAddresses, err := client.GetInstanceIPAddresses(context.Background(), instanceWithReservedIP.ID)
	if err != nil {
		t.Fatalf("Failed to get instance info for Linode with reserved IP: %v", err)
	}

	// Check if the 'reserved' field is set to true
	foundReserved := false
	for _, ip := range instanceAddresses.IPv4.Public {
		if ip.Address == reservedIP.Address {
			if !ip.Reserved {
				t.Errorf("Expected 'Reserved' field to be true for reserved IP %s, but it was false", ip.Address)
			}
			foundReserved = true
			break
		}
	}
	if !foundReserved {
		t.Errorf("Reserved IP %s not found in instance's public IP addresses", reservedIP.Address)
	}

	// Create a Linode with an ephemeral IP
	instanceWithEphemeralIP, err := client.CreateInstance(context.Background(), linodego.InstanceCreateOptions{
		Region:   "us-east",
		Type:     "g6-nanode-1",
		Label:    linodego.Pointer("test-instance-ephemeral-ip"),
		RootPass: linodego.Pointer(randPassword()),
	})
	if err != nil {
		t.Fatalf("Failed to create Linode with ephemeral IP: %v", err)
	}
	defer func() {
		if err := client.DeleteInstance(context.Background(), instanceWithEphemeralIP.ID); err != nil {
			t.Errorf("Error deleting test Instance with ephemeral IP: %s", err)
		}
	}()

	// Make GET request for the Linode with ephemeral IP
	ephemeralInstanceAddresses, err := client.GetInstanceIPAddresses(context.Background(), instanceWithEphemeralIP.ID)
	if err != nil {
		t.Fatalf("Failed to get instance IP addresses for Linode with ephemeral IP: %v", err)
	}

	// Check that all public IPs have 'Reserved' field set to false
	for _, ip := range ephemeralInstanceAddresses.IPv4.Public {
		if ip.Reserved {
			t.Errorf("Expected 'Reserved' field to be false for ephemeral IP %s, but it was true", ip.Address)
		}
	}
}
