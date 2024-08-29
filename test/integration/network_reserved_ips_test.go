package integration

import (
	"context"
	"fmt"
	"strings"

	"testing"

	"github.com/linode/linodego"
	. "github.com/linode/linodego"
)

// TestReservedIPAddresses_InsufficientPermissions tests the behavior when a user account
// doesn't have the can_reserve_ip flag enabled
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
	fmt.Println(ipList)

	if err != nil {
		t.Fatalf("Error listing IP addresses: %v", err)
	}

	initialCount := len(ipList)

	// Attempt to reserve an IP
	resIP, resErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{
		Region: "us-east",
	})

	if resErr != nil {
		t.Fatalf("Failed to reserve IP. This test should start with 0 reservations or reservations < limit. Error from the API: %v", resErr)
	}

	if resIP == nil {
		t.Fatalf("Reserved IP is nil")
	}

	t.Logf("Successfully reserved IP: %+v", resIP)

	// Fetch the reserved IP
	fetchedIP, fetchErr := client.GetReservedIPAddress(context.Background(), resIP.Address)
	if fetchErr != nil {
		t.Errorf("Error getting reserved IP address: %v", fetchErr)
	}

	if fetchedIP == nil {
		t.Errorf("Retrieved reserved IP is nil")
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
		t.Errorf("Expected error when fetching deleted IP, got nil")
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

	for _, ip := range ipList {
		if !ip.Reserved {
			t.Errorf("Expected all IPs to be reserved, but found non-reserved IP: %s", ip.Address)
		}
		if ip.Region != "us-east" {
			t.Errorf("Expected all IPs to be in us-east region, but found IP in %s region: %s", ip.Region, ip.Address)
		}
	}

	if len(ipList) == 0 {
		t.Log("No reserved IPs found in us-east region")
	} else {
		t.Logf("First reserved IP in us-east: %+v", ipList[0])
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
		t.Fatalf("Failed to reserve IP. This test should start with 0 reservations or reservations < limit. Error from the API: %v", resErr)
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

	// Test valid IP reservations until limit is reached
	for {
		reservation, err := client.ReserveIPAddress(context.Background(), ReserveIPOptions{Region: "us-east"})
		if err != nil {
			if strings.Contains(err.Error(), "Additional Reserved IPv4 addresses require technical justification. Please contact support describing your requirement.") {
				t.Logf("Reservation limit reached after %d reservations", len(reservedIPs))
				break
			}
			t.Fatalf("Unexpected error when reserving IP: %v", err)
		}

		if reservation == nil {
			t.Fatalf("Valid reservation returned nil IP")
		}

		reservedIPs = append(reservedIPs, reservation.Address)
		t.Logf("Successfully reserved IP: %s", reservation.Address)
	}

	// Verify that we can't reserve more IPs
	_, exceedErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{Region: "us-east"})
	if exceedErr == nil {
		t.Errorf("Expected error when exceeding reservation limit, got nil")
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
		t.Skip("No reserved IPs available for testing deletion")
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
