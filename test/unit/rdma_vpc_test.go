package unit

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVPC_RDMA_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("vpc_rdma_get")
	require.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("vpcs/7", fixtureData)

	vpc, err := base.Client.GetVPC(context.Background(), 7)
	require.NoError(t, err)
	require.NotNil(t, vpc)

	assert.Equal(t, 7, vpc.ID)
	assert.Equal(t, "test-vpc-rdma", vpc.Label)
	assert.Equal(t, "fake-cph-5", vpc.Region)
	assert.Equal(t, linodego.VPCTypeRDMA, vpc.VPCType)
	assert.Equal(t, "RDMA VPC for GPUDirect", vpc.Description)
	assert.Empty(t, vpc.IPv6)

	// Subnet assertions
	require.Len(t, vpc.Subnets, 1)
	subnet := vpc.Subnets[0]
	assert.Equal(t, 8, subnet.ID)
	assert.Equal(t, "rdma-subnet", subnet.Label)
	assert.Equal(t, "10.0.0.0/8", subnet.IPv4)
	assert.Equal(t, linodego.VPCTypeRDMA, subnet.VPCType)
	assert.Empty(t, subnet.IPv6)

	// Subnet linode/interface assertions
	require.Len(t, subnet.Linodes, 1)
	assert.Equal(t, 506958, subnet.Linodes[0].ID)
	require.Len(t, subnet.Linodes[0].Interfaces, 1)
	assert.Equal(t, 10, subnet.Linodes[0].Interfaces[0].ID)
	assert.Nil(t, subnet.Linodes[0].Interfaces[0].ConfigID)
	assert.Equal(t, false, subnet.Linodes[0].Interfaces[0].Active)
}

func TestVPC_RDMA_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("vpc_rdma_create")
	require.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	createOptions := linodego.VPCCreateOptions{
		Label:       "new-rdma-vpc",
		Description: "A new RDMA VPC",
		Region:      "fake-cph-5",
		VPCType:     linodego.VPCTypeRDMA,
		Subnets: []linodego.VPCSubnetCreateOptions{
			{Label: "rdma-subnet-1", IPv4: "10.0.0.0/24"},
		},
	}

	httpmock.RegisterRegexpResponder(
		"POST",
		mockRequestURL(t, "/vpcs"),
		mockRequestBodyValidate(t, createOptions, fixtureData),
	)

	vpc, err := base.Client.CreateVPC(context.Background(), createOptions)
	require.NoError(t, err)
	require.NotNil(t, vpc)

	assert.Equal(t, 39, vpc.ID)
	assert.Equal(t, "new-rdma-vpc", vpc.Label)
	assert.Equal(t, linodego.VPCTypeRDMA, vpc.VPCType)
	assert.Equal(t, "fake-cph-5", vpc.Region)

	require.Len(t, vpc.Subnets, 1)
	assert.Equal(t, 40, vpc.Subnets[0].ID)
	assert.Equal(t, linodego.VPCTypeRDMA, vpc.Subnets[0].VPCType)
}

func TestVPC_Regular_VPCType(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("vpcs", linodego.VPC{
		ID:      100,
		Label:   "regular-vpc",
		Region:  "us-east",
		VPCType: linodego.VPCTypeRegular,
	})

	vpc, err := base.Client.CreateVPC(context.Background(), linodego.VPCCreateOptions{
		Label:   "regular-vpc",
		Region:  "us-east",
		VPCType: linodego.VPCTypeRegular,
	})
	require.NoError(t, err)
	assert.Equal(t, linodego.VPCTypeRegular, vpc.VPCType)
}

func TestVPC_VPCType_OmittedWhenEmpty(t *testing.T) {
	opts := linodego.VPCCreateOptions{
		Label:  "test",
		Region: "us-east",
	}

	data, err := json.Marshal(opts)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	_, exists := parsed["vpc_type"]
	assert.False(t, exists, "vpc_type should be omitted when empty")
}

// =============================================================================
// VPC Subnet RDMA Type Tests
// =============================================================================

func TestVPCSubnet_RDMA_VPCType(t *testing.T) {
	// Test that VPCSubnet properly unmarshals vpc_type field
	jsonData := `{
		"id": 8,
		"label": "rdma-subnet",
		"ipv4": "10.0.0.0/8",
		"vpc_type": "rdma",
		"ipv6": [],
		"linodes": [],
		"databases": [],
		"nodebalancers": [],
		"created": "2026-03-12T09:51:58",
		"updated": "2026-03-12T09:51:58"
	}`

	var subnet linodego.VPCSubnet
	err := json.Unmarshal([]byte(jsonData), &subnet)
	require.NoError(t, err)

	assert.Equal(t, 8, subnet.ID)
	assert.Equal(t, "rdma-subnet", subnet.Label)
	assert.Equal(t, linodego.VPCTypeRDMA, subnet.VPCType)
}

