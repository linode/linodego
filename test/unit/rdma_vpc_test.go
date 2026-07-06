package unit

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
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

func TestInterface_GetRDMAVPC(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("interface_get_rdma_vpc")
	require.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/506958/interfaces/10", fixtureData)

	iface, err := base.Client.GetInterface(context.Background(), 506958, 10)
	require.NoError(t, err)
	require.NotNil(t, iface)

	assert.Equal(t, 10, iface.ID)
	assert.Equal(t, 1, iface.Version)
	assert.Equal(t, "22:00:f2:9e:d3:48", iface.MACAddress)
	assert.Equal(t, false, *iface.DefaultRoute.IPv4)
	assert.Equal(t, false, *iface.DefaultRoute.IPv6)

	// Non-RDMA fields should be nil
	assert.Nil(t, iface.Public)
	assert.Nil(t, iface.VPC)
	assert.Nil(t, iface.VLAN)

	// RDMA VPC assertions
	require.NotNil(t, iface.RDMAVPC)
	assert.Equal(t, 7, iface.RDMAVPC.VPCID)
	assert.Equal(t, 8, iface.RDMAVPC.SubnetID)

	require.Len(t, iface.RDMAVPC.IPv4.Addresses, 1)
	assert.Equal(t, "10.0.0.2", iface.RDMAVPC.IPv4.Addresses[0].Address)
	assert.Equal(t, true, iface.RDMAVPC.IPv4.Addresses[0].Primary)
}

func TestInterface_ListWithRDMA(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("interface_list_with_rdma")
	require.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/111/interfaces", fixtureData)

	ifaces, err := base.Client.ListInterfaces(context.Background(), 111, nil)
	require.NoError(t, err)
	require.Len(t, ifaces, 2)

	// First interface: regular VPC
	assert.Equal(t, 111111, ifaces[0].ID)
	assert.Nil(t, ifaces[0].RDMAVPC)
	require.NotNil(t, ifaces[0].VPC)
	assert.Equal(t, 123, ifaces[0].VPC.VPCID)
	assert.Equal(t, 456, ifaces[0].VPC.SubnetID)
	assert.True(t, *ifaces[0].DefaultRoute.IPv4)
	assert.Nil(t, ifaces[0].Public)

	// Second interface: RDMA VPC
	assert.Equal(t, 222222, ifaces[1].ID)
	require.NotNil(t, ifaces[1].RDMAVPC)
	assert.Equal(t, 123, ifaces[1].RDMAVPC.VPCID)
	assert.Equal(t, 456, ifaces[1].RDMAVPC.SubnetID)
	assert.Equal(t, "10.0.0.2", ifaces[1].RDMAVPC.IPv4.Addresses[0].Address)
	assert.True(t, ifaces[1].RDMAVPC.IPv4.Addresses[0].Primary)
	assert.Nil(t, ifaces[1].VPC)
	assert.Nil(t, ifaces[1].Public)
}

func TestInterface_UpdateRDMAVPC(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("interface_update_rdma_vpc")
	require.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("linode/instances/506958/interfaces/10", fixtureData)

	opts := linodego.LinodeInterfaceUpdateOptions{
		RDMAVPC: &linodego.RDMAVPCInterfaceUpdateOptions{
			SubnetID: 9,
			IPv4: linodego.RDMAVPCInterfaceIPv4Options{
				Addresses: []linodego.RDMAVPCInterfaceIPv4AddressOptions{
					{Address: "10.0.1.5", Primary: linodego.Pointer(true)},
				},
			},
		},
	}

	iface, err := base.Client.UpdateInterface(context.Background(), 506958, 10, opts)
	require.NoError(t, err)
	require.NotNil(t, iface)

	assert.Equal(t, 10, iface.ID)
	assert.Equal(t, 2, iface.Version)
	require.NotNil(t, iface.RDMAVPC)
	assert.Equal(t, 7, iface.RDMAVPC.VPCID)
	assert.Equal(t, 9, iface.RDMAVPC.SubnetID)
	assert.Equal(t, "10.0.1.5", iface.RDMAVPC.IPv4.Addresses[0].Address)
	assert.True(t, iface.RDMAVPC.IPv4.Addresses[0].Primary)
}

// =============================================================================
// Instance Create with RDMA Interfaces (LinodeInstanceInterfaces) Tests
// =============================================================================

func TestInstance_CreateWithRDMAInterfaces_MarshalJSON(t *testing.T) {
	createOptions := linodego.InstanceCreateOptions{
		Region:              "fake-cph-5",
		Type:                "g2-gpu-rdma-1",
		InterfaceGeneration: linodego.GenerationLinode,
		LinodeInstanceInterfaces: []linodego.LinodeInstanceInterfaceCreateOptions{
			{
				RDMAVPC: &linodego.RDMAVPCInterfaceCreateOptions{
					SubnetID: 1234,
					IPv4: linodego.RDMAVPCInterfaceIPv4Options{
						Addresses: []linodego.RDMAVPCInterfaceIPv4AddressOptions{
							{Address: "auto", Primary: linodego.Pointer(true)},
						},
					},
				},
			},
		},
	}

	data, err := json.Marshal(createOptions)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	// Verify interfaces array is correctly serialized
	ifaces, ok := parsed["interfaces"].([]interface{})
	require.True(t, ok, "expected interfaces to be an array")
	require.Len(t, ifaces, 1)

	ifaceMap, ok := ifaces[0].(map[string]interface{})
	require.True(t, ok)

	rdmaVPC, ok := ifaceMap["rdma_vpc"].(map[string]interface{})
	require.True(t, ok, "expected rdma_vpc key in interface")
	assert.Equal(t, float64(1234), rdmaVPC["subnet_id"])

	ipv4, ok := rdmaVPC["ipv4"].(map[string]interface{})
	require.True(t, ok)
	addresses, ok := ipv4["addresses"].([]interface{})
	require.True(t, ok)
	require.Len(t, addresses, 1)

	addr, ok := addresses[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "auto", addr["address"])
	assert.Equal(t, true, addr["primary"])
}

func TestInstance_CreateWithLinodeInterfaces_BackwardsCompatible(t *testing.T) {
	// Ensure existing LinodeInterfaces field still works
	createOptions := linodego.InstanceCreateOptions{
		Region:              "us-east",
		Type:                "g6-standard-1",
		InterfaceGeneration: linodego.GenerationLinode,
		LinodeInterfaces: []linodego.LinodeInterfaceCreateOptions{
			{
				VPC: &linodego.VPCInterfaceCreateOptions{
					SubnetID: 4,
					IPv4: &linodego.VPCInterfaceIPv4CreateOptions{
						Addresses: &[]linodego.VPCInterfaceIPv4AddressCreateOptions{
							{
								Address: linodego.Pointer("10.0.0.5"),
								Primary: linodego.Pointer(true),
							},
						},
					},
				},
			},
		},
	}

	data, err := json.Marshal(createOptions)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	// Verify interfaces array is correctly serialized
	ifaces, ok := parsed["interfaces"].([]interface{})
	require.True(t, ok, "expected interfaces to be an array")
	require.Len(t, ifaces, 1)

	ifaceMap, ok := ifaces[0].(map[string]interface{})
	require.True(t, ok)

	vpcData, ok := ifaceMap["vpc"].(map[string]interface{})
	require.True(t, ok, "expected vpc key in interface")
	assert.Equal(t, float64(4), vpcData["subnet_id"])
}

func TestInstance_Create_ConflictingInterfaceFields(t *testing.T) {
	t.Run("LinodeInterfaces_and_LinodeInstanceInterfaces", func(t *testing.T) {
		createOptions := linodego.InstanceCreateOptions{
			Region: "us-east",
			Type:   "g6-standard-1",
			LinodeInterfaces: []linodego.LinodeInterfaceCreateOptions{
				{},
			},
			LinodeInstanceInterfaces: []linodego.LinodeInstanceInterfaceCreateOptions{
				{},
			},
		}

		_, err := json.Marshal(createOptions)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "LinodeInterfaces and LinodeInstanceInterfaces")
	})

	t.Run("Interfaces_and_LinodeInstanceInterfaces", func(t *testing.T) {
		createOptions := linodego.InstanceCreateOptions{
			Region: "us-east",
			Type:   "g6-standard-1",
			Interfaces: []linodego.InstanceConfigInterfaceCreateOptions{
				{},
			},
			LinodeInstanceInterfaces: []linodego.LinodeInstanceInterfaceCreateOptions{
				{},
			},
		}

		_, err := json.Marshal(createOptions)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Interfaces and LinodeInstanceInterfaces")
	})
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

// =============================================================================
// RDMA VPC Interface Option Marshaling-Semantics Tests
// =============================================================================

// jsonToMap marshals v and unmarshals it into a generic map for key inspection.
func jsonToMap(t *testing.T, v any) map[string]any {
	t.Helper()

	data, err := json.Marshal(v)
	require.NoError(t, err)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(data, &parsed))

	return parsed
}

func TestRDMAVPCInterface_ExplicitPrimaryFalse(t *testing.T) {
	// A caller MUST be able to send "primary": false explicitly.
	// Primary is a pointer precisely so that false != unset.
	opts := linodego.RDMAVPCInterfaceCreateOptions{
		SubnetID: 1234,
		IPv4: linodego.RDMAVPCInterfaceIPv4Options{
			Addresses: []linodego.RDMAVPCInterfaceIPv4AddressOptions{
				{
					Address: "10.0.0.5",
					Primary: linodego.Pointer(false),
				},
			},
		},
	}

	parsed := jsonToMap(t, opts)

	ipv4 := parsed["ipv4"].(map[string]any)
	addresses := ipv4["addresses"].([]any)
	require.Len(t, addresses, 1)

	addr := addresses[0].(map[string]any)
	primary, exists := addr["primary"]
	require.True(t, exists, "primary must be present when explicitly set to false")
	assert.Equal(t, false, primary)
}

func TestRDMAVPCInterface_ExplicitEmptyAddresses(t *testing.T) {
	// On update, a caller MUST be able to send an explicit empty addresses list.
	// A non-nil empty slice is NOT the zero value, so omitzero still marshals it.
	opts := linodego.RDMAVPCInterfaceUpdateOptions{
		IPv4: linodego.RDMAVPCInterfaceIPv4Options{
			Addresses: []linodego.RDMAVPCInterfaceIPv4AddressOptions{},
		},
	}

	parsed := jsonToMap(t, opts)

	ipv4, exists := parsed["ipv4"]
	require.True(t, exists, "ipv4 must be present when it holds an explicit empty addresses list")

	addresses, exists := ipv4.(map[string]any)["addresses"]
	require.True(t, exists, "addresses must be present when explicitly set to an empty slice")

	addrSlice, ok := addresses.([]any)
	require.True(t, ok, "addresses should serialize as a JSON array")
	assert.Empty(t, addrSlice, "addresses should be an explicit empty array")
}
