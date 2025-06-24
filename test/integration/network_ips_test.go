package integration

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/linode/linodego"
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

func TestIPAddress_GetFound_smoke(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestIPAddress_GetFound", true)
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

func TestIPAddresses_List_smoke(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestIPAddresses_List", true)
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

	// Set the RDNS for the IPv6 addresses
	for _, ip := range i {
		if ip.Type != "ipv6" || !ip.Public || ip.RDNS == "" {
			continue
		}

		rdns := fmt.Sprintf("%s.nip.io", ip.Address)
		_, err = client.UpdateIPAddressV2(context.Background(), ip.Address, IPAddressUpdateOptionsV2{
			RDNS: linodego.Pointer(linodego.Pointer(rdns)),
		})
		if err != nil {
			t.Fatalf("Failed to set RDNS for IPv6 address: %v", err)
		}
	}

	// Retrieve the networking info without IPv6 RDNS
	i, err = client.ListIPAddresses(context.Background(), &ListOptions{
		QueryParams: ListIPAddressesQuery{
			SkipIPv6RDNS: true,
		},
	})
	if err != nil {
		t.Errorf("Error listing ipaddresses, expected struct, got error %v", err)
	}

	// Validate that the IPv6 address RDNS field is not returned
	for _, ip := range i {
		if strings.Contains(string(ip.Type), "ipv6") && ip.RDNS != "" {
			t.Fatalf("expected empty rdns for ipv6 address; got %s", ip.RDNS)
		}
	}

	reservedFilter := "{\"reserved\":true}"
	reservedIpAddresses, err := client.ListIPAddresses(context.Background(), NewListOptions(0, reservedFilter))
	if err != nil {
		t.Errorf("Error listing ipaddresses, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of ipaddresses, but got none %v", reservedIpAddresses)
	}

	// Verify that all IPs in the reserved list are actually reserved
	for _, ip := range reservedIpAddresses {
		if !ip.Reserved {
			t.Errorf("IP %s is in the reserved list but has Reserved field set to false", ip.Address)
		}
	}

	unreservedFilter := "{\"reserved\":false}"
	unreservedIpAddresses, err := client.ListIPAddresses(context.Background(), NewListOptions(0, unreservedFilter))
	if err != nil {
		t.Errorf("Error listing ipaddresses, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of ipaddresses, but got none %v", unreservedIpAddresses)
	}

	// Verify that all IPs in the reserved list are actually unreserved
	for _, ip := range unreservedIpAddresses {
		if ip.Reserved {
			t.Errorf("IP %s is in the non-reserved list but has Reserved field set to true", ip.Address)
		}
	}
}

func TestIPAddresses_Instance_Get(t *testing.T) {
	client, _, _, instance, config, teardown := setupInstanceWith3Interfaces(t, "fixtures/TestIPAddresses_Instance_Get")
	defer teardown()

	i, err := client.GetInstanceIPAddresses(context.Background(), instance.ID)
	if err != nil {
		t.Errorf("Error listing ipaddresses, expected struct, got error %v", err)
	}
	if i.IPv4.Public[0].Address != instance.IPv4[0].String() {
		t.Errorf("Expected matching public IP ipaddresses with GetInstanceIPAddress Instance IPAddress but got %v", i)
	}
	if *i.IPv4.VPC[0].Address != *config.Interfaces[2].IPv4.VPC {
		t.Errorf("Expected matching VPC IP addresses with GetInstanceIPAddress Instance IPAddress but got %v", i)
	}
}

func TestIPAddress_Update(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestIPAddress_Update", true)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	reservedTrue := true
	reservedFalse := false

	address := instance.IPv4[0].String()

	i, err := client.GetInstanceIPAddresses(context.Background(), instance.ID)
	if err != nil {
		t.Errorf("Error getting ipaddress: %s", err)
	}

	originalRDNS := i.IPv4.Public[0].RDNS

	// Update RDNS to nip.io
	updateOpts := IPAddressUpdateOptionsV2{
		RDNS: linodego.Pointer(linodego.Pointer(fmt.Sprintf("%s.nip.io", i.IPv4.Public[0].Address))),
	}

	ip, err := client.UpdateIPAddressV2(context.Background(), address, updateOpts)
	require.NoError(t, err)

	// Update RDNS to default
	updateOpts = IPAddressUpdateOptionsV2{
		RDNS: linodego.Pointer[*string](nil),
	}
	ip, err = client.UpdateIPAddressV2(context.Background(), ip.Address, updateOpts)
	require.NoError(t, err)

	require.NotNil(t, ip)
	require.Equal(t, originalRDNS, ip.RDNS)

	createReservedIP := func() (string, error) {
		reservedIP, err := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{Region: instance.Region})
		if err != nil {
			return "", err
		}
		return reservedIP.Address, nil
	}

	// Scenario 1: Convert ephemeral IP to reserved IP

	ephemeralIP := instance.IPv4[0].String()
	updateOpts = IPAddressUpdateOptionsV2{
		Reserved: &reservedTrue,
	}
	updatedIP, err := client.UpdateIPAddressV2(context.Background(), ephemeralIP, updateOpts)
	if err != nil {
		t.Fatalf("Failed to convert ephemeral IP to reserved: %v", err)
	}
	if !updatedIP.Reserved {
		t.Errorf("Expected IP to be reserved, but it's not")
	}

	// Scenario 2: Convert reserved IP to reserved IP (no-op)

	reservedIP, err := createReservedIP()
	if err != nil {
		t.Fatalf("Failed to create reserved IP: %v", err)
	}
	defer client.DeleteReservedIPAddress(context.Background(), reservedIP)

	updateOpts = IPAddressUpdateOptionsV2{
		Reserved: &reservedTrue,
	}
	updatedIP, err = client.UpdateIPAddressV2(context.Background(), reservedIP, updateOpts)
	if err != nil {
		t.Fatalf("Failed to update reserved IP: %v", err)
	}
	if !updatedIP.Reserved {
		t.Errorf("Expected IP to remain reserved, but it's not")
	}

	// Scenario 3: Convert reserved to ephemeral

	ephemeralIP = instance.IPv4[0].String()
	updateOpts = IPAddressUpdateOptionsV2{
		Reserved: &reservedFalse,
	}
	updatedIP, err = client.UpdateIPAddressV2(context.Background(), ephemeralIP, updateOpts)
	if err != nil {
		t.Fatalf("Failed to update ephemeral IP: %v", err)
	}
	if updatedIP.Reserved {
		t.Errorf("Expected IP to remain ephemeral, but it's reserved")
	}

	// Scenario 4: Convert assigned reserved IP to ephemeral
	reservedIP, err = createReservedIP()
	if err != nil {
		t.Fatalf("Failed to create reserved IP: %v", err)
	}
	defer client.DeleteReservedIPAddress(context.Background(), reservedIP)

	// Assign the reserved IP to the instance
	assignOpts := LinodesAssignIPsOptions{
		Region: instance.Region,
		Assignments: []LinodeIPAssignment{
			{
				Address:  reservedIP,
				LinodeID: instance.ID,
			},
		},
	}
	err = client.InstancesAssignIPs(context.Background(), assignOpts)
	if err != nil {
		t.Fatalf("Failed to assign reserved IP: %v", err)
	}

	updateOpts = IPAddressUpdateOptionsV2{
		Reserved: &reservedFalse,
	}
	updatedIP, err = client.UpdateIPAddressV2(context.Background(), reservedIP, updateOpts)
	if err != nil {
		t.Fatalf("Failed to convert assigned reserved IP to ephemeral: %v", err)
	}
	if updatedIP.Reserved {
		t.Errorf("Expected IP to be converted to ephemeral, but it's still reserved")
	}

	// Scenario 5: Cannot set RDNS for unassigned reserved IP

	unassignedResIP, unassignedResIpErr := createReservedIP()
	if unassignedResIpErr != nil {
		t.Fatalf("Failed to create reserved IP: %v", unassignedResIpErr)
	}

	updateOpts = IPAddressUpdateOptionsV2{
		Reserved: &reservedTrue,
		RDNS:     linodego.Pointer(linodego.Pointer("sample rdns")),
	}
	_, err = client.UpdateIPAddressV2(context.Background(), unassignedResIP, updateOpts)
	if err == nil {
		t.Fatalf("Expected error when setting RDNS for unassigned reserved IP, but got none")
	}

	client.DeleteReservedIPAddress(context.Background(), unassignedResIP)

	// Scenario 6: Convert unassigned reserved IP to reserved (no-op)

	reservedIP, err = createReservedIP()
	if err != nil {
		t.Fatalf("Failed to create reserved IP: %v", err)
	}

	updateOpts = IPAddressUpdateOptionsV2{
		Reserved: &reservedTrue,
	}
	updatedIP, err = client.UpdateIPAddressV2(context.Background(), reservedIP, updateOpts)
	if err != nil {
		t.Fatalf("Failed to update unassigned reserved IP: %v", err)
	}
	if !updatedIP.Reserved || updatedIP.LinodeID != 0 {
		t.Errorf("Expected IP to remain unassigned reserved, but got: %+v", updatedIP)
	}

	client.DeleteReservedIPAddress(context.Background(), reservedIP)

	// Scenario 7: Convert unassigned reserved IP to unassigned (delete)

	reservedIP, err = createReservedIP()
	if err != nil {
		t.Fatalf("Failed to create reserved IP: %v", err)
	}

	updateOpts = IPAddressUpdateOptionsV2{
		Reserved: &reservedFalse,
	}
	_, err = client.UpdateIPAddressV2(context.Background(), reservedIP, updateOpts)
	if err != nil {
		t.Fatalf("Failed to convert unassigned reserved IP to unassigned: %v", err)
	}

	// Verify the IP has been deleted
	_, err = client.GetIPAddress(context.Background(), reservedIP)
	if err == nil {
		t.Errorf("Expected IP to be deleted, but it still exists")
	}

	// Scenario 10: Cannot convert non-owned reserved IP

	invalidResIp := "123.72.121.76"

	updateOpts = IPAddressUpdateOptionsV2{
		Reserved: &reservedFalse,
	}

	updatedIP, err = client.UpdateIPAddressV2(context.Background(), invalidResIp, updateOpts)
	if err == nil {
		t.Fatalf("Expected error indicating the IP address is invalid, got nil")
	}
}

// TestIPAddress_Instance_Delete requires the customer account to have
// default_IPMax set to at least 2 and default_InterfaceMax set to 3.
func TestIPAddress_Instance_Delete(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestIPAddress_Instance_Delete", true)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	opts := linodego.InstanceIPAddOptions{
		Public: true,
	}
	ip, err := client.AddInstanceIPAddress(context.TODO(), instance.ID, opts)
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
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestIPAddress_Instance_Assign", true)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	newInstance, err := createInstance(t, client, true, func(client *Client, options *InstanceCreateOptions) {
		options.Label = linodego.Pointer("go-ins-test-assign")
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
		LinodeID:     &newInstance.ID,
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
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestIPAddress_Instance_Share", true)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	newInstance, err := createInstance(t, client, true, func(client *Client, options *InstanceCreateOptions) {
		options.Label = linodego.Pointer("go-ins-test-share")
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

	opts := linodego.InstanceIPAddOptions{
		Public: true,
	}
	ip, err := client.AddInstanceIPAddress(context.Background(), newInstance.ID, opts)
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

func TestIPAddress_Instance_Allocate(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestIPAddress_Instance_Allocate", true)
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	// Scenario 1: Valid request

	opts := AllocateReserveIPOptions{
		Type:     "ipv4",
		Public:   true,
		Reserved: linodego.Pointer(true),
		Region:   &instance.Region,
		LinodeID: &instance.ID,
	}
	validIp, err := client.AllocateReserveIP(context.Background(), opts)
	// defer cleanUpIPAllocation(t, client, validIp.Address)
	if err != nil {
		t.Fatalf("Expected successful IP reservation, got error: %v", err)
	}
	if !validIp.Reserved || validIp.LinodeID != instance.ID {
		t.Errorf("Unexpected IP reservation result: %+v", validIp)
	}

	// Scenario 2: Non-owned Linode
	nonOwnedLinodeOpts := AllocateReserveIPOptions{
		Type:     "ipv4",
		Public:   true,
		Reserved: linodego.Pointer(true),
		Region:   &instance.Region,
		LinodeID: linodego.Pointer(99999), // Assume this is a non-owned Linode ID
	}
	_, nonOwnedLinodeErr := client.AllocateReserveIP(context.Background(), nonOwnedLinodeOpts)
	if nonOwnedLinodeErr == nil {
		t.Fatal("Expected error for non-owned Linode, got nil")
	}

	// Scenario 3: Omit Linode ID

	omitLinodeIDOpts := AllocateReserveIPOptions{
		Type:     "ipv4",
		Public:   true,
		Reserved: linodego.Pointer(true),
		Region:   &instance.Region,
	}
	omitLinodeIDip, omitLinodeErr := client.AllocateReserveIP(context.Background(), omitLinodeIDOpts)
	if omitLinodeErr != nil {
		t.Fatalf("Expected successful unassigned IP reservation, got error: %v", omitLinodeErr)
	}
	if !omitLinodeIDip.Reserved || omitLinodeIDip.LinodeID != 0 || omitLinodeIDip.Region != instance.Region {
		t.Errorf("Unexpected unassigned IP reservation result: %+v", omitLinodeIDip)
	}

	// Scenario 4: Omit Region

	omitRegionOpts := AllocateReserveIPOptions{
		Type:     "ipv4",
		Public:   true,
		Reserved: linodego.Pointer(true),
		LinodeID: &instance.ID,
	}
	omitRegionip, omitRegionErr := client.AllocateReserveIP(context.Background(), omitRegionOpts)
	// defer cleanUpIPAllocation(t, client, omitRegionip.Address)
	if omitRegionErr != nil {
		t.Fatalf("Expected successful IP reservation without region, got error: %v", omitRegionErr)
	}
	if !omitRegionip.Reserved || omitRegionip.LinodeID != instance.ID {
		t.Errorf("Unexpected IP reservation result without region: %+v", omitRegionip)
	}

	cleanUpReserveIPAllocation(t, client, omitRegionip.Address)

	// Scenario 5: Omit both Region and Linode ID

	omitRegionAndLinodeIDopts := AllocateReserveIPOptions{
		Type:     "ipv4",
		Public:   true,
		Reserved: linodego.Pointer(true),
	}
	_, omitRegionAndLinodeIDerr := client.AllocateReserveIP(context.Background(), omitRegionAndLinodeIDopts)

	if omitRegionAndLinodeIDerr == nil {
		t.Fatal("Expected error when omitting both region and Linode ID, got nil")
	}

	// Scenario 6: Reserved true, Public false

	publicFalseOpts := AllocateReserveIPOptions{
		Type:     "ipv4",
		Public:   false,
		Reserved: linodego.Pointer(true),
		Region:   &instance.Region,
		LinodeID: &instance.ID,
	}
	_, publicFalseErr := client.AllocateReserveIP(context.Background(), publicFalseOpts)
	if publicFalseErr == nil {
		t.Fatal("Expected error for reserved true and public false, got nil")
	}

	// Scenario 7: Reserved false

	reservedFalseOpts := AllocateReserveIPOptions{
		Type:     "ipv4",
		Public:   true,
		Reserved: linodego.Pointer(false),
		Region:   &instance.Region,
		LinodeID: &instance.ID,
	}
	reservedFalseIp, reservedFalseErr := client.AllocateReserveIP(context.Background(), reservedFalseOpts)
	if reservedFalseErr != nil {
		t.Fatalf("Expected successful ephemeral IP assignment, got error: %v", reservedFalseErr)
	}
	if reservedFalseIp.Reserved || reservedFalseIp.LinodeID != instance.ID {
		t.Errorf("Unexpected ephemeral IP assignment result: %+v", reservedFalseIp)
	}

	cleanUpIPAllocation(t, client, instance.ID, reservedFalseIp.Address)

	// Scenario 8: Omit Reserved field

	omitReservedOpts := AllocateReserveIPOptions{
		Type:     "ipv4",
		Public:   true,
		Region:   &instance.Region,
		LinodeID: &instance.ID,
	}
	omitReservedip, omitReservedErr := client.AllocateReserveIP(context.Background(), omitReservedOpts)
	if omitReservedErr != nil {
		t.Fatalf("Expected successful IP assignment, got error: %v", omitReservedErr)
	}
	if omitReservedip.Reserved || omitReservedip.LinodeID != instance.ID {
		t.Errorf("Unexpected IP assignment result: %+v", omitReservedip)
	}

	cleanUpIPAllocation(t, client, instance.ID, omitReservedip.Address)

	// Scenario 9: Omit Linode ID, Reserved false

	omitOpts := AllocateReserveIPOptions{
		Type:     "ipv4",
		Public:   true,
		Reserved: linodego.Pointer(false),
		Region:   &instance.Region,
	}
	_, omitOptsErr := client.AllocateReserveIP(context.Background(), omitOpts)
	if omitOptsErr == nil {
		t.Fatal("Expected error when omitting Linode ID and setting reserved to false, got nil")
	}

	// Scenario 10: Omit Linode ID and Reserved fields

	omitIDResopts := AllocateReserveIPOptions{
		Type:   "ipv4",
		Public: true,
		Region: &instance.Region,
	}
	_, omitIDResErr := client.AllocateReserveIP(context.Background(), omitIDResopts)
	if omitIDResErr == nil {
		t.Fatal("Expected error when omitting Linode ID and reserved fields, got nil")
	}

	// Scenario 11: Reserved true, Type IPv6

	typeIPv6opts := AllocateReserveIPOptions{
		Type:     "ipv6",
		Public:   true,
		Reserved: linodego.Pointer(true),
		Region:   &instance.Region,
		LinodeID: &instance.ID,
	}
	_, typeIPv6Err := client.AllocateReserveIP(context.Background(), typeIPv6opts)
	if typeIPv6Err == nil {
		t.Fatal("Expected error for reserved true and type IPv6, got nil")
	}

	// Scenario 12: Reserved false, Type IPv6

	resFalseIPv6opts := AllocateReserveIPOptions{
		Type:     "ipv6",
		Public:   true,
		Reserved: linodego.Pointer(false),
		Region:   &instance.Region,
		LinodeID: &instance.ID,
	}
	_, resFalseIPv6Err := client.AllocateReserveIP(context.Background(), resFalseIPv6opts)
	if resFalseIPv6Err == nil {
		t.Fatalf("Expected unsuccessful IPv6 assignment, got nil")
	}

	// Scenario 13: Region mismatch

	regionMismatchOpts := AllocateReserveIPOptions{
		Type:     "ipv4",
		Public:   true,
		Reserved: linodego.Pointer(true),
		Region:   linodego.Pointer("us-west"), // Assume this is different from instance.Region
		LinodeID: &instance.ID,
	}
	_, regionMismatchErr := client.AllocateReserveIP(context.Background(), regionMismatchOpts)
	if regionMismatchErr == nil {
		t.Fatal("Expected error for region mismatch, got nil")
	}
}

func cleanUpReserveIPAllocation(t *testing.T, client *linodego.Client, address string) {
	err := client.DeleteReservedIPAddress(context.Background(), address)
	if err != nil {
		t.Logf("Failed to delete reserved IP %s: %v", address, err)
	}
}

func cleanUpIPAllocation(t *testing.T, client *linodego.Client, linodeID int, address string) {
	err := client.DeleteInstanceIPAddress(context.Background(), linodeID, address)
	if err != nil {
		t.Logf("Failed to delete reserved IP %s: %v", address, err)
	}
}

func TestIPAddress_Instance_ReserveIP_Assign(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestIPAddress_Instance_ReserveIP_Assign")
	defer teardown()

	// Create two Linodes for testing
	linode1, err := createInstance(t, client, true)
	if err != nil {
		t.Fatalf("Error creating first test Linode: %s", err)
	}
	defer func() {
		if err := client.DeleteInstance(context.Background(), linode1.ID); err != nil {
			t.Errorf("Error deleting first test Linode: %s", err)
		}
	}()

	linode2, err := createInstance(t, client, true)
	if err != nil {
		t.Fatalf("Error creating second test Linode: %s", err)
	}
	defer func() {
		if err := client.DeleteInstance(context.Background(), linode2.ID); err != nil {
			t.Errorf("Error deleting second test Linode: %s", err)
		}
	}()

	// Scenario 1: Assign unassigned reserved IP to existing Linode

	reservedIP, err := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{Region: linode1.Region})
	if err != nil {
		t.Fatalf("Failed to reserve IP: %v", err)
	}
	defer cleanUpReserveIPAllocation(t, client, reservedIP.Address)

	err = client.InstancesAssignIPs(context.Background(), LinodesAssignIPsOptions{
		Region: linode1.Region,
		Assignments: []LinodeIPAssignment{
			{
				Address:  reservedIP.Address,
				LinodeID: linode1.ID,
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to assign reserved IP: %v", err)
	}

	// Verify assignment
	ip, err := client.GetIPAddress(context.Background(), reservedIP.Address)
	if err != nil {
		t.Fatalf("Failed to get IP address info: %v", err)
	}
	if !ip.Reserved || ip.LinodeID != linode1.ID {
		t.Errorf("Unexpected IP assignment result: %+v", ip)
	}

	// Scenario 2: Reassign reserved IP to different Linode

	reassignIP, reassignIPerr := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{Region: linode1.Region})
	if reassignIPerr != nil {
		t.Fatalf("Failed to reserve IP: %v", err)
	}

	// Assign to first Linode
	assignErr := client.InstancesAssignIPs(context.Background(), LinodesAssignIPsOptions{
		Region: linode1.Region,
		Assignments: []LinodeIPAssignment{
			{
				Address:  reassignIP.Address,
				LinodeID: linode1.ID,
			},
		},
	})
	if assignErr != nil {
		t.Fatalf("Failed to assign reserved IP to first Linode: %v", assignErr)
	}

	// Reassign to second Linode
	reassignErr := client.InstancesAssignIPs(context.Background(), LinodesAssignIPsOptions{
		Region: linode2.Region,
		Assignments: []LinodeIPAssignment{
			{
				Address:  reassignIP.Address,
				LinodeID: linode2.ID,
			},
		},
	})
	if reassignErr != nil {
		t.Fatalf("Failed to reassign reserved IP: %v", reassignErr)
	}

	// Verify reassignment
	ipAddress, getIpErr := client.GetIPAddress(context.Background(), reassignIP.Address)
	if getIpErr != nil {
		t.Fatalf("Failed to get IP address info: %v", getIpErr)
	}
	if !ipAddress.Reserved || ipAddress.LinodeID != linode2.ID {
		t.Errorf("Unexpected IP reassignment result: %+v", ipAddress)
	}

	cleanUpReserveIPAllocation(t, client, reassignIP.Address)

	// Scenario 3: Attempt to assign non-owned reserved IP

	invalidIpErr := client.InstancesAssignIPs(context.Background(), LinodesAssignIPsOptions{
		Region: linode1.Region,
		Assignments: []LinodeIPAssignment{
			{
				Address:  "192.0.2.1", // Assume this is a non-owned IP
				LinodeID: linode1.ID,
			},
		},
	})
	if invalidIpErr == nil {
		t.Fatal("Expected error when assigning non-owned reserved IP, got nil")
	}

	// Scenario 4: Attempt to assign owned reserved IP to non-owned Linode

	validResIP, validResIPerr := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{Region: linode1.Region})
	if validResIPerr != nil {
		t.Fatalf("Failed to reserve IP: %v", validResIPerr)
	}
	defer cleanUpReserveIPAllocation(t, client, validResIP.Address)

	invalidLinodeErr := client.InstancesAssignIPs(context.Background(), LinodesAssignIPsOptions{
		Region: linode1.Region,
		Assignments: []LinodeIPAssignment{
			{
				Address:  validResIP.Address,
				LinodeID: 99999, // Assume this is a non-owned Linode ID
			},
		},
	})
	if invalidLinodeErr == nil {
		t.Fatal("Expected error when assigning to non-owned Linode, got nil")
	}
}
