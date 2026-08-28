package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestNATGateways_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nat_gateways_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("networking/natgateways", fixtureData)

	gateways, err := base.Client.ListNATGateways(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, gateways, 1)

	assert.Equal(t, 42, gateways[0].ID)
	assert.Equal(t, "gb-lon", gateways[0].Region)
	assert.Len(t, gateways[0].Addresses, 2)
	assert.Equal(t, "203.0.113.42", gateways[0].Addresses[0].Address)
	assert.Equal(t, "203.0.113.43", gateways[0].Addresses[1].Address)
	assert.Equal(t, 0, gateways[0].AddressAutoscaleMax)
	assert.Equal(t, 4096, gateways[0].DefaultPortsPerInterface)
	assert.Equal(t, "my-cloud-nat", gateways[0].Label)
	assert.Equal(t, 0, gateways[0].PortsetAssignments)
	assert.Equal(t, 126, gateways[0].PortsetCapacity)

	assert.Equal(t, 42, gateways[0].VPCSubnet.ID)
	assert.Equal(t, "vpc_subnet", gateways[0].VPCSubnet.Type)
	assert.Equal(t, "my-subnet", gateways[0].VPCSubnet.Label)
	assert.Equal(t, "https://www.something.com/vpcs/42/subnets/41", gateways[0].VPCSubnet.URL)
	assert.Equal(t, 13, gateways[0].VPCSubnet.VPCID)
	assert.Equal(t, "my-vpc", gateways[0].VPCSubnet.VPCLabel)
}

func TestNATGateways_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nat_gateways_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("networking/natgateways/42", fixtureData)

	gateway, err := base.Client.GetNATGateway(context.Background(), 42)
	assert.NoError(t, err)

	assert.Equal(t, 42, gateway.ID)
	assert.Equal(t, "gb-lon", gateway.Region)
	assert.Len(t, gateway.Addresses, 2)
	assert.Equal(t, "203.0.113.42", gateway.Addresses[0].Address)
	assert.Equal(t, "203.0.113.43", gateway.Addresses[1].Address)
	assert.Equal(t, 0, gateway.AddressAutoscaleMax)
	assert.Equal(t, 4096, gateway.DefaultPortsPerInterface)
	assert.Equal(t, "my-cloud-nat", gateway.Label)
	assert.Equal(t, 0, gateway.PortsetAssignments)
	assert.Equal(t, 126, gateway.PortsetCapacity)

	assert.Equal(t, 42, gateway.VPCSubnet.ID)
	assert.Equal(t, "vpc_subnet", gateway.VPCSubnet.Type)
	assert.Equal(t, "my-subnet", gateway.VPCSubnet.Label)
	assert.Equal(t, "https://www.something.com/vpcs/42/subnets/41", gateway.VPCSubnet.URL)
	assert.Equal(t, 13, gateway.VPCSubnet.VPCID)
	assert.Equal(t, "my-vpc", gateway.VPCSubnet.VPCLabel)
}

func TestNATGateways_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nat_gateways_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	createOptions := linodego.NATGatewayCreateOptions{
		Label:                    "my-cloud-nat",
		Region:                   "gb-lon",
		DefaultPortsPerInterface: linodego.Pointer(4096),
		UseAutoscaling:           linodego.Pointer(true),
		VPCSubnetID:              linodego.Pointer(42),
		Addresses: []linodego.NATGatewayAddressCreateOptions{
			{Address: "203.0.113.42"},
			{Address: "203.0.113.43"},
		},
	}

	base.MockPost("networking/natgateways", fixtureData)

	gateway, err := base.Client.CreateNATGateway(context.Background(), createOptions)
	assert.NoError(t, err)

	assert.Equal(t, 42, gateway.ID)
	assert.Equal(t, "gb-lon", gateway.Region)
	assert.Len(t, gateway.Addresses, 2)
	assert.Equal(t, "203.0.113.42", gateway.Addresses[0].Address)
	assert.Equal(t, "203.0.113.43", gateway.Addresses[1].Address)
	assert.Equal(t, 0, gateway.AddressAutoscaleMax)
	assert.Equal(t, 4096, gateway.DefaultPortsPerInterface)
	assert.Equal(t, "my-cloud-nat", gateway.Label)
	assert.Equal(t, 0, gateway.PortsetAssignments)
	assert.Equal(t, 126, gateway.PortsetCapacity)

	assert.Equal(t, 42, gateway.VPCSubnet.ID)
	assert.Equal(t, "vpc_subnet", gateway.VPCSubnet.Type)
	assert.Equal(t, "my-subnet", gateway.VPCSubnet.Label)
	assert.Equal(t, "https://www.something.com/vpcs/42/subnets/41", gateway.VPCSubnet.URL)
	assert.Equal(t, 13, gateway.VPCSubnet.VPCID)
	assert.Equal(t, "my-vpc", gateway.VPCSubnet.VPCLabel)
}

func TestNATGateways_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nat_gateways_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	updateOptions := linodego.NATGatewayUpdateOptions{
		Label: linodego.Pointer("my-cloud-nat-updated"),
	}

	base.MockPut("networking/natgateways/42", fixtureData)

	gateway, err := base.Client.UpdateNATGateway(context.Background(), updateOptions, 42)
	assert.NoError(t, err)

	assert.Equal(t, "my-cloud-nat-updated", gateway.Label)
}

func TestNATGateways_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("networking/natgateways/42", nil)

	err := base.Client.DeleteNATGateway(context.Background(), 42)
	assert.NoError(t, err)
}

func TestNATGateways_AddAddress(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nat_gateways_add_address")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	addAddressesOptions := linodego.NATGatewayAddAddressOptions{
		Address: "203.0.113.42",
	}

	base.MockPost("networking/natgateways/42/addresses", fixtureData)

	address, err := base.Client.NATGatewayAddAddress(context.Background(), addAddressesOptions, 42)
	assert.NoError(t, err)

	assert.Equal(t, "203.0.113.43", address.Address)
	assert.Equal(t, true, address.InUse)
	assert.Equal(t, 6, address.InterfaceCount)
	assert.Equal(t, 6, address.PortsetAssignments)
	assert.Equal(t, 63, address.PortsetCapacity)
	assert.Equal(t, "https://www.something.com/interfaces/42", address.InterfaceURL)
}

func TestNATGateways_ListAddresses(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nat_gateways_list_addresses")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("networking/natgateways/42/addresses", fixtureData)

	addresses, err := base.Client.NATGatewayListAddresses(context.Background(), nil, 42)
	assert.NoError(t, err)
	assert.Len(t, addresses, 1)

	assert.Equal(t, "203.0.113.43", addresses[0].Address)
	assert.Equal(t, true, addresses[0].InUse)
	assert.Equal(t, 6, addresses[0].InterfaceCount)
	assert.Equal(t, 6, addresses[0].PortsetAssignments)
	assert.Equal(t, 63, addresses[0].PortsetCapacity)
	assert.Equal(t, "https://www.something.com/interfaces/42", addresses[0].InterfaceURL)
}

func TestNATGateways_GetAddress(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nat_gateways_get_address")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("networking/natgateways/42/addresses/203.0.113.43", fixtureData)

	address, err := base.Client.NATGatewayGetAddress(context.Background(), 42, "203.0.113.43")
	assert.NoError(t, err)

	assert.Equal(t, "203.0.113.43", address.Address)
	assert.Equal(t, true, address.InUse)
	assert.Equal(t, 6, address.InterfaceCount)
	assert.Equal(t, 6, address.PortsetAssignments)
	assert.Equal(t, 63, address.PortsetCapacity)
	assert.Equal(t, "https://www.something.com/interfaces/42", address.InterfaceURL)
}

func TestNATGateways_DeleteAddress(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("networking/natgateways/42/addresses/203.0.113.43", nil)

	err := base.Client.NATGatewayDeleteAddress(context.Background(), 42, "203.0.113.43")
	assert.NoError(t, err)
}

func TestNATGateways_ListInterfaces(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nat_gateways_list_interfaces")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("networking/natgateways/42/interfaces", fixtureData)

	interfaces, err := base.Client.NATGatewayListInterfaces(context.Background(), nil, 42)
	assert.NoError(t, err)

	assert.Len(t, interfaces, 2)

	assert.Equal(t, 142, interfaces[0].ID)
	assert.Equal(t, 1001, interfaces[0].Linode.ID)
	assert.Equal(t, "linode1001", interfaces[0].Linode.Label)
	assert.Equal(t, "linode", interfaces[0].Linode.Type)
	assert.Equal(t, "https://www.something.com/1", interfaces[0].Linode.URL)
	assert.Equal(t, "172.24.213.144", interfaces[0].Addresses[0])
	assert.Equal(t, "203.0.113.42", interfaces[0].Portsets[0].Address)
	assert.Equal(t, 2048, interfaces[0].Portsets[0].Ports[0].Start)
	assert.Equal(t, 3071, interfaces[0].Portsets[0].Ports[0].End)

	assert.Equal(t, 143, interfaces[1].ID)
	assert.Equal(t, 1002, interfaces[1].Linode.ID)
	assert.Equal(t, "linode1002", interfaces[1].Linode.Label)
	assert.Equal(t, "linode", interfaces[1].Linode.Type)
	assert.Equal(t, "https://www.something.com/2", interfaces[1].Linode.URL)
	assert.Equal(t, "172.24.213.144", interfaces[1].Addresses[0])
	assert.Equal(t, "203.0.113.43", interfaces[1].Portsets[0].Address)
	assert.Equal(t, 2048, interfaces[1].Portsets[0].Ports[0].Start)
	assert.Equal(t, 3071, interfaces[1].Portsets[0].Ports[0].End)
}

func TestNATGateways_ListAddressInterfaces(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nat_gateways_list_address_interfaces")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("networking/natgateways/42/addresses/172.24.213.144/interfaces", fixtureData)

	interfaces, err := base.Client.NATGatewayListAddressInterfaces(context.Background(), nil, 42, "172.24.213.144")
	assert.NoError(t, err)

	assert.Len(t, interfaces, 1)

	assert.Equal(t, 142, interfaces[0].ID)
	assert.Equal(t, 1001, interfaces[0].Linode.ID)
	assert.Equal(t, "linode1001", interfaces[0].Linode.Label)
	assert.Equal(t, "linode", interfaces[0].Linode.Type)
	assert.Equal(t, "https://www.something.com/1", interfaces[0].Linode.URL)
	assert.Equal(t, "172.24.213.144", interfaces[0].Addresses[0])
	assert.Equal(t, "203.0.113.42", interfaces[0].Portsets[0].Address)
	assert.Equal(t, 2048, interfaces[0].Portsets[0].Ports[0].Start)
	assert.Equal(t, 3071, interfaces[0].Portsets[0].Ports[0].End)
}

func TestNATGateways_GetTypes(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nat_gateways_get_types")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("networking/natgateways/types", fixtureData)

	types, err := base.Client.NATGatewayGetTypes(context.Background(), nil)
	assert.NoError(t, err)

	assert.Equal(t, "g1-natgateway", types[0].ID)
	assert.Equal(t, "NAT Gateway", types[0].Label)
	assert.Equal(t, 0.035, types[0].Price.Hourly)
	assert.Equal(t, 25.0, types[0].Price.Monthly)
}

func TestNATGateways_GetSettings(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nat_gateways_get_settings")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("networking/natgateways/settings", fixtureData)

	settings, err := base.Client.NATGatewayGetSettings(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, 100, settings.MaximumReservedAddressesPerNATGateway)
	assert.Equal(t, 100, settings.MaximumAutoscalingAddressesPerNATGateway)
	assert.Len(t, settings.AllowedPortsPerInterface, 3)
	assert.Equal(t, 4096, settings.AllowedPortsPerInterface[0])
	assert.Equal(t, 8192, settings.AllowedPortsPerInterface[1])
	assert.Equal(t, 16384, settings.AllowedPortsPerInterface[2])
}
