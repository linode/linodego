package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestReservedIPAddresses_List(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock API response with additional attributes
	mockResponse := struct {
		Data []linodego.InstanceIP `json:"data"`
	}{
		Data: []linodego.InstanceIP{
			{
				Address:  "192.168.1.10",
				Region:   "us-east",
				LinodeID: 12345,
				Reserved: true,
				Tags:     []string{"lb"},
			},
			{
				Address:  "192.168.1.20",
				Region:   "us-west",
				LinodeID: 67890,
				Reserved: true,
				Tags:     []string{},
			},
		},
	}

	base.MockGet("networking/reserved/ips", mockResponse)

	reservedIPs, err := base.Client.ListReservedIPAddresses(context.Background(), nil)

	assert.NoError(t, err, "Expected no error when listing reserved IP addresses")
	assert.NotNil(t, reservedIPs, "Expected non-nil reserved IP addresses")
	assert.Len(t, reservedIPs, 2, "Expected two reserved IP addresses")
	assert.Equal(t, "192.168.1.10", reservedIPs[0].Address, "Expected first reserved IP address to match")
	assert.Equal(t, "us-east", reservedIPs[0].Region, "Expected region to match")
	assert.Equal(t, 12345, reservedIPs[0].LinodeID, "Expected Linode ID to match")
	assert.True(t, reservedIPs[0].Reserved, "Expected first IP to be reserved")
	assert.Equal(t, []string{"lb"}, reservedIPs[0].Tags, "Expected tags to match")
}

func TestReservedIPAddress_Get(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	ip := "192.168.1.10"

	// Mock response with necessary attributes
	mockResponse := linodego.InstanceIP{
		Address:  ip,
		Region:   "us-east",
		LinodeID: 12345,
		Reserved: true,
		Tags:     []string{"lb"},
	}

	base.MockGet("networking/reserved/ips/"+ip, mockResponse)

	reservedIP, err := base.Client.GetReservedIPAddress(context.Background(), ip)

	assert.NoError(t, err, "Expected no error when getting reserved IP address")
	assert.NotNil(t, reservedIP, "Expected non-nil reserved IP address")
	assert.Equal(t, ip, reservedIP.Address, "Expected reserved IP address to match")
	assert.Equal(t, "us-east", reservedIP.Region, "Expected region to match")
	assert.Equal(t, 12345, reservedIP.LinodeID, "Expected Linode ID to match")
	assert.True(t, reservedIP.Reserved, "Expected IP to be reserved")
	assert.Equal(t, []string{"lb"}, reservedIP.Tags, "Expected tags to match")
}

func TestIPReserveIPAddress(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	createOpts := linodego.ReserveIPOptions{
		Region: "us-west",
	}

	// Mock the POST request for reserving an IP
	mockResponse := linodego.InstanceIP{
		Address:  "192.168.1.30",
		Region:   "us-west",
		LinodeID: 13579,
		Reserved: true,
		Tags:     []string{"env:staging"},
	}

	base.MockPost("networking/reserved/ips", mockResponse)

	reservedIP, err := base.Client.ReserveIPAddress(context.Background(), createOpts)

	assert.NoError(t, err, "Expected no error when reserving IP address")
	assert.NotNil(t, reservedIP, "Expected non-nil reserved IP address")
	assert.Equal(t, "192.168.1.30", reservedIP.Address, "Expected reserved IP address to match")
	assert.Equal(t, "us-west", reservedIP.Region, "Expected region to match")
	assert.Equal(t, 13579, reservedIP.LinodeID, "Expected Linode ID to match")
	assert.Equal(t, []string{"env:staging"}, reservedIP.Tags, "Expected tags to match")
}

func TestReservedIPAddress_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	ip := "192.168.1.10"

	// Mock the DELETE request for deleting the reserved IP
	base.MockDelete("networking/reserved/ips/"+ip, nil)

	err := base.Client.DeleteReservedIPAddress(context.Background(), ip)

	assert.NoError(t, err, "Expected no error when deleting reserved IP address")
}

func TestUpdateReservedIPAddress(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("network_reserved_ip_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	ip := "192.168.1.10"
	base.MockPut("networking/reserved/ips/"+ip, fixtureData)

	updateOpts := linodego.UpdateReservedIPOptions{
		Tags: []string{"lb", "team:infra"},
	}

	updated, err := base.Client.UpdateReservedIPAddress(context.Background(), ip, updateOpts)

	assert.NoError(t, err, "Expected no error when updating reserved IP address")
	assert.NotNil(t, updated)
	assert.Equal(t, ip, updated.Address)
	assert.True(t, updated.Reserved)
	assert.Equal(t, []string{"lb", "team:infra"}, updated.Tags)
}

func TestListReservedIPTypes(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("network_reserved_ip_types_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("networking/reserved/ips/types", fixtureData)

	types, err := base.Client.ListReservedIPTypes(context.Background(), nil)

	assert.NoError(t, err, "Expected no error when listing reserved IP types")
	assert.Len(t, types, 1)
	assert.Equal(t, "ipv4_address", types[0].ID)
	assert.Equal(t, "IPv4 Address", types[0].Label)
	assert.Equal(t, 0.005, types[0].Price.Hourly)
	assert.Equal(t, 2.00, types[0].Price.Monthly)
	assert.Len(t, types[0].RegionPrices, 1)
	assert.Equal(t, "us-east", types[0].RegionPrices[0].ID)
	assert.Equal(t, 0.006, types[0].RegionPrices[0].Hourly)
}
