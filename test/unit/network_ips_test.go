package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestIPUpdateAddressV2(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	ip := "192.168.1.1"

	// Mock API response
	base.MockPut("networking/ips/"+ip, linodego.InstanceIP{
		Address:  ip,
		Reserved: true,
	})

	updatedIP, err := base.Client.UpdateIPAddressV2(context.Background(), ip, linodego.IPAddressUpdateOptionsV2{
		Reserved: linodego.Pointer(true),
	})
	assert.NoError(t, err, "Expected no error when updating IP address")
	assert.NotNil(t, updatedIP, "Expected non-nil updated IP address")
	assert.Equal(t, ip, updatedIP.Address, "Expected updated IP address to match")
	assert.True(t, updatedIP.Reserved, "Expected Reserved to be true")
}

func TestIPAllocateReserve(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock API response
	base.MockPost("networking/ips", linodego.InstanceIP{
		Address: "192.168.1.3",
		Region:  "us-east",
		Public:  true,
	})

	ip, err := base.Client.AllocateReserveIP(context.Background(), linodego.AllocateReserveIPOptions{
		Type:     "ipv4",
		Public:   true,
		Region:   linodego.Pointer("us-east"),
		LinodeID: linodego.Pointer(12345),
	})
	assert.NoError(t, err, "Expected no error when allocating reserve IP")
	assert.NotNil(t, ip, "Expected non-nil allocated IP")
	assert.Equal(t, "192.168.1.3", ip.Address, "Expected allocated IP address to match")
	assert.Equal(t, "us-east", ip.Region, "Expected Region to match")
	assert.True(t, ip.Public, "Expected Public to be true")
}

func TestIPAssignInstances(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock API response
	base.MockPost("networking/ips/assign", nil)

	err := base.Client.InstancesAssignIPs(context.Background(), linodego.LinodesAssignIPsOptions{
		Region: "us-east",
		Assignments: []linodego.LinodeIPAssignment{
			{Address: "192.168.1.10", LinodeID: 123},
		},
	})
	assert.NoError(t, err, "Expected no error when assigning IPs to instances")
}

func TestIPShareAddresses(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock API response
	base.MockPost("networking/ips/share", nil)

	err := base.Client.ShareIPAddresses(context.Background(), linodego.IPAddressesShareOptions{
		IPs:      []string{"192.168.1.20"},
		LinodeID: 456,
	})
	assert.NoError(t, err, "Expected no error when sharing IP addresses")
}
