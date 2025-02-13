package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/linode/linodego"
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
			},
			{
				Address:  "192.168.1.20",
				Region:   "us-west",
				LinodeID: 67890,
				Reserved: true,
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
	}

	base.MockGet("networking/reserved/ips/"+ip, mockResponse)

	reservedIP, err := base.Client.GetReservedIPAddress(context.Background(), ip)

	assert.NoError(t, err, "Expected no error when getting reserved IP address")
	assert.NotNil(t, reservedIP, "Expected non-nil reserved IP address")
	assert.Equal(t, ip, reservedIP.Address, "Expected reserved IP address to match")
	assert.Equal(t, "us-east", reservedIP.Region, "Expected region to match")
	assert.Equal(t, 12345, reservedIP.LinodeID, "Expected Linode ID to match")
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
	}

	base.MockPost("networking/reserved/ips", mockResponse)

	reservedIP, err := base.Client.ReserveIPAddress(context.Background(), createOpts)

	assert.NoError(t, err, "Expected no error when reserving IP address")
	assert.NotNil(t, reservedIP, "Expected non-nil reserved IP address")
	assert.Equal(t, "192.168.1.30", reservedIP.Address, "Expected reserved IP address to match")
	assert.Equal(t, "us-west", reservedIP.Region, "Expected region to match")
	assert.Equal(t, 13579, reservedIP.LinodeID, "Expected Linode ID to match")
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
