package integration

import (
	"context"
	"os"
	"testing"

	. "github.com/linode/linodego"
)

// TestReservedIPAddresses_InsufficientPermissions tests the behavior when a user account
// doesn't have the can_reserve_ip flag enabled
func TestReservedIPAddresses_InsufficientPermissions(t *testing.T) {
	original := os.Getenv("LINODE_TOKEN")
	dummyToken := "badtoken"
	os.Setenv("LINODE_TOKEN", dummyToken)

	client, teardown := createTestClient(t, "fixtures/TestReservedIpAddresses_InsufficientPermissions")
	defer teardown()
	defer os.Setenv("LINODE_TOKEN", original)

	filter := ""
	i, getReservedIpsError := client.GetReservedIPs(context.Background(), NewListOptions(0, filter))
	if getReservedIpsError != nil {
		t.Logf("Error listing ipaddresses, expected struct, got error %v", getReservedIpsError)
	}

	if len(i) == 0 {
		t.Logf("Expected a list of ipaddresses, but got none %v", i)
	}

	// Attempt to reserve an IP address
	reservedIp, reserveIpError := client.ReserveIPAddress(context.Background(), ReserveIPOptions{
		Region: "us-east",
	})
	if reserveIpError != nil {
		t.Logf("Failed to reserve IP: %v", reserveIpError)
	} else {
		t.Logf("Successfully reserved IP: %+v", reservedIp)
	}

	// Attempt to get a reserved IP address
	address := "172.28.3.4"
	ipaddress, getReservedIpError := client.GetReservedIPAddress(context.Background(), address)
	if getReservedIpError != nil {
		t.Logf("Error getting ipaddress, expected struct, got %v and error %v", ipaddress, getReservedIpError)
	}

	// Attempt to delete a reserved IP address
	addressToBeDeleted := "172.28.3.4"
	deleterr := client.DeleteReservedIPAddress(context.Background(), addressToBeDeleted)
	if deleterr != nil {
		t.Logf("Error deleting reserved IP address: %v", deleterr)
	} else {
		t.Logf("Successfully deleted reserved IP address: %s", addressToBeDeleted)
	}
}

// TestReservedIPAddresses_EndToEndTest performs an end-to-end test of the Reserved IP functionality
// for users with the can_reserve_ip flag enabled
func TestReservedIPAddresses_EndToEndTest(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestReservedIpAddresses_EndToEndTest")
	defer teardown()

	filter := ""
	ipList, listIpErr := client.GetReservedIPs(context.Background(), NewListOptions(0, filter))
	if listIpErr != nil {
		t.Logf("Error listing ipaddresses, expected struct, got error %v", listIpErr)
	} else {
		if len(ipList) == 0 {
			t.Logf("The customer has not reserved an IP %v", ipList)
		}
	}

	// Attempt to reserve an IP
	reserveIP, reserveIpErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{
		Region: "us-east",
	})
	if reserveIpErr != nil {
		t.Logf("Failed to reserve IP: %v", reserveIpErr)
	} else {
		t.Logf("Successfully reserved IP: %+v", reserveIP)
	}

	if reserveIP != nil {
		// Fetch the reserved IP
		reservedIP, fetchIpErr := client.GetReservedIPAddress(context.Background(), reserveIP.Address)
		if fetchIpErr != nil {
			t.Logf("Error getting ipaddress, expected struct, got %v and error %v", reservedIP, fetchIpErr)
		}

		// Verify the list of IPs has increased
		verifyList, verifyListIpErr := client.GetReservedIPs(context.Background(), NewListOptions(0, filter))
		if verifyListIpErr != nil {
			t.Logf("Error listing ipaddresses, expected struct, got error %v", verifyListIpErr)
		} else {
			if len(verifyList)-len(ipList) == 1 {
				t.Log("Increase in IP list confirmed", verifyList)
			} else {
				t.Errorf("Increase in IP list not confirmed")
			}
		}

		// Delete the reserved IP
		if reserveIP != nil {
			deleteErr := client.DeleteReservedIPAddress(context.Background(), reserveIP.Address)
			if deleteErr != nil {
				t.Logf("Error deleting reserved IP address: %v", deleteErr)
			} else {
				t.Logf("Successfully deleted reserved IP address: %s", reserveIP.Address)
			}
		}

		// Verify the IP has been deleted
		if reserveIP != nil {
			verifyDeletedIP, fetchDeletedIpErr := client.GetReservedIPAddress(context.Background(), reserveIP.Address)
			if fetchDeletedIpErr != nil {
				t.Logf("Error getting ipaddress, expected struct, got %v and error %v", verifyDeletedIP, fetchDeletedIpErr)
			}

			verifyDeletedFromList, verifyDeletedFromListIpErr := client.GetReservedIPs(context.Background(), NewListOptions(0, filter))
			if verifyDeletedFromListIpErr != nil {
				t.Logf("Error listing ipaddresses, expected struct, got error %v", verifyDeletedFromListIpErr)
			} else {
				if len(verifyDeletedFromList) < len(verifyList) {
					t.Log("IP address deletion confirmed", verifyDeletedFromList)
				} else {
					t.Errorf("Verification - Failed")
				}
			}
		}
	}
}

// TestReservedIPAddresses_ListIPAddressesVariants tests various scenarios for listing IP addresses
func TestReservedIPAddresses_ListIPAddressesVariants(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestReservedIPAddresses_ListIPAddressesVariants")
	defer teardown()

	filter := ""
	ipList, listIpErr := client.GetReservedIPs(context.Background(), NewListOptions(0, filter))
	if listIpErr != nil {
		t.Logf("Error listing ipaddresses, expected struct, got error %v", listIpErr)
	} else {
		if len(ipList) == 0 {
			t.Logf("The customer has not reserved an IP %v", ipList)
		}
	}
}

// TestReservedIPAddresses_GetIPAddressVariants tests various scenarios for getting a specific IP address
func TestReservedIPAddresses_GetIPAddressVariants(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestReservedIPAddresses_GetIPAddressVariants")
	defer teardown()

	// Reserve an IP for testing
	reserveIP, reserveIpErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{
		Region: "us-east",
	})
	if reserveIpErr != nil {
		t.Logf("Failed to reserve IP: %v", reserveIpErr)
	} else {
		t.Logf("Successfully reserved IP: %+v", reserveIP)
	}

	if reserveIP != nil {
		// Test getting a valid reserved IP
		validReservedIP, validFetchIpErr := client.GetReservedIPAddress(context.Background(), reserveIP.Address)
		if validFetchIpErr != nil {
			t.Logf("Error getting ipaddress, expected struct, got %v and error %v", validReservedIP, validFetchIpErr)
		}

		// Test getting an invalid IP
		invalidReservedIP, invalidFetchIpErr := client.GetReservedIPAddress(context.Background(), "INVALID_IP")
		if invalidFetchIpErr != nil {
			t.Logf("Error getting ipaddress, expected struct, got %v and error %v", invalidReservedIP, invalidFetchIpErr)
		}
	}
}

// TestReservedIPAddresses_ReserveIPVariants tests various scenarios for reserving an IP address
func TestReservedIPAddresses_ReserveIPVariants(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestReservedIPAddresses_ReserveIPVariants")
	defer teardown()

	// Test reserving IP with omitted region
	omitRegion, omitRegionErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{})
	if omitRegionErr != nil {
		t.Logf("Failed to reserve IP: %v", omitRegionErr)
	} else {
		t.Logf("Successfully reserved IP: %+v", omitRegion)
	}

	// Test reserving IP with invalid region
	invalidRegion, invalidRegionErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{
		Region: "us",
	})
	if invalidRegionErr != nil {
		t.Logf("Failed to reserve IP: %v", invalidRegionErr)
	} else {
		t.Logf("Successfully reserved IP: %+v", invalidRegion)
	}

	// Test reserving IP with empty region
	emptyRegion, emptyRegionErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{
		Region: "",
	})
	if emptyRegionErr != nil {
		t.Logf("Failed to reserve IP: %v", emptyRegionErr)
	} else {
		t.Logf("Successfully reserved IP: %+v", emptyRegion)
	}

	// Test valid IP reservation
	validReservation, validReservationErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{
		Region: "us-east",
	})
	if validReservationErr != nil {
		t.Logf("Failed to reserve IP: %v", validReservationErr)
	} else {
		t.Logf("Successfully reserved IP: %+v", validReservation)
	}

	// Test exceeding reservation limit
	exceedReservationLimit, exceedReservationErr := client.ReserveIPAddress(context.Background(), ReserveIPOptions{
		Region: "us-east",
	})
	if exceedReservationErr != nil {
		t.Logf("Failed to reserve IP: %v", exceedReservationErr)
	} else {
		t.Logf("Successfully reserved IP: %+v", exceedReservationLimit)
	}
}

// TestReservedIPAddresses_DeleteIPAddressVariants tests various scenarios for deleting a reserved IP address
func TestReservedIPAddresses_DeleteIPAddressVariants(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestReservedIPAddresses_DeleteIPAddressVariants")
	defer teardown()

	filter := ""
	ipList, listIpErr := client.GetReservedIPs(context.Background(), NewListOptions(0, filter))
	if listIpErr != nil {
		t.Logf("Error listing ipaddresses, expected struct, got error %v", listIpErr)
	} else {
		if len(ipList) == 0 {
			t.Logf("The customer has not reserved an IP %v", ipList)
		}
	}

	if len(ipList) > 0 {
		ipToBeDeleted := ipList[len(ipList)-1]
		deleteErr := client.DeleteReservedIPAddress(context.Background(), ipToBeDeleted.Address)
		if deleteErr != nil {
			t.Logf("Error deleting reserved IP address: %v", deleteErr)
		} else {
			t.Logf("Successfully deleted reserved IP address: %s", ipToBeDeleted.Address)
		}

		// Verify deletion
		verifyDeletedFromList, verifyDeletedFromListIpErr := client.GetReservedIPs(context.Background(), NewListOptions(0, filter))
		if verifyDeletedFromListIpErr != nil {
			t.Logf("Error listing ipaddresses, expected struct, got error %v", verifyDeletedFromListIpErr)
		} else {
			if len(verifyDeletedFromList) < len(ipList) {
				t.Log("IP address deletion confirmed", verifyDeletedFromList)
			} else {
				t.Errorf("Verification - Failed")
			}
		}

		verifyDeletedIP, fetchDeletedIpErr := client.GetReservedIPAddress(context.Background(), ipToBeDeleted.Address)
		if fetchDeletedIpErr != nil {
			t.Logf("IP address not found got %v and error %v", verifyDeletedIP, fetchDeletedIpErr)
		}

		// Test deleting an unowned IP
		deletingUnownedIPErr := client.DeleteReservedIPAddress(context.Background(), "255.255.255.4")
		if deletingUnownedIPErr != nil {
			t.Logf("Error deleting reserved IP address: %v", deletingUnownedIPErr)
		} else {
			t.Logf("Successfully deleted reserved IP address: %s", "255.255.255.4")
		}
	}
}
