package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestVPC_Create(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("vpcs", linodego.VPC{
		ID:          123,
		Label:       "test-vpc",
		Description: "Test VPC description",
		Region:      "us-east",
		IPv6: []linodego.VPCIPv6Range{
			{
				Range: "fd71:1140:a9d0::/52",
			},
		},
		Subnets: []linodego.VPCSubnet{
			{ID: 1, Label: "subnet-1"},
			{ID: 2, Label: "subnet-2"},
		},
		Created: linodego.Pointer(time.Now()),
		Updated: linodego.Pointer(time.Now()),
	})

	vpc, err := base.Client.CreateVPC(context.Background(), linodego.VPCCreateOptions{
		Label:       "test-vpc",
		Description: "Test VPC description",
		Region:      "us-east",
		IPv6: []linodego.VPCCreateOptionsIPv6{
			{
				Range:           linodego.Pointer("/52"),
				AllocationClass: linodego.Pointer("test"),
			},
		},
		Subnets: []linodego.VPCSubnetCreateOptions{
			{Label: "subnet-1"},
			{Label: "subnet-2"},
		},
	})
	assert.NoError(t, err, "Expected no error when creating VPC")
	assert.NotNil(t, vpc, "Expected VPC to be created")
	assert.Equal(t, "test-vpc", vpc.Label, "Expected VPC label to match")
	assert.Equal(t, "us-east", vpc.Region, "Expected VPC region to match")
	assert.Len(t, vpc.Subnets, 2, "Expected VPC to have 2 subnets")
}

func TestVPC_Get(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	fixtureData, err := fixtures.GetFixture("vpc_get")
	require.NoError(t, err)

	base.MockGet("vpcs/123", fixtureData)

	vpc, err := base.Client.GetVPC(context.Background(), 123)
	assert.NoError(t, err, "Expected no error when getting VPC")
	assert.NotNil(t, vpc, "Expected non-nil VPC")

	assertJSONObjectsSimilar(t, vpc, vpc.GetCreateOptions())
	assertJSONObjectsSimilar(t, vpc, vpc.GetUpdateOptions())

	assert.Equal(t, 123, vpc.ID, "Expected VPC ID to match")
	assert.Equal(t, "cool-vpc", vpc.Label, "Expected VPC label to match")
	assert.Equal(t, "A description of my VPC.", vpc.Description)
	assert.Equal(t, "us-east", vpc.Region)

	assert.Equal(t, "fd71:1140:a9d0::/52", vpc.IPv6[0].Range)

	assert.NotEmpty(t, vpc.Subnets)
	assert.Equal(t, 456, vpc.Subnets[0].ID)
	assert.Equal(t, "subnet-1", vpc.Subnets[0].Label)
	assert.Equal(t, "192.168.1.0/24", vpc.Subnets[0].IPv4)
	assert.Equal(t, 111, vpc.Subnets[0].Linodes[0].ID)
	assert.Equal(t, true, vpc.Subnets[0].Linodes[0].Interfaces[0].Active)
	assert.Equal(t, 4567, *vpc.Subnets[0].Linodes[0].Interfaces[0].ConfigID)
	assert.Equal(t, 421, vpc.Subnets[0].Linodes[0].Interfaces[0].ID)
}

func TestVPC_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("vpc_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("vpcs", fixtureData)

	vpcs, err := base.Client.ListVPCs(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.NotEmpty(t, vpcs, "Expected non-empty VPC list")

	vpc := vpcs[0]

	assertJSONObjectsSimilar(t, vpc, vpc.GetCreateOptions())
	assertJSONObjectsSimilar(t, vpc, vpc.GetUpdateOptions())

	assert.Equal(t, 123, vpc.ID, "Expected VPC ID to match")
	assert.Equal(t, "cool-vpc", vpc.Label, "Expected VPC label to match")
	assert.Equal(t, "A description of my VPC.", vpc.Description, "Expected VPC description to match")
	assert.Equal(t, "us-east", vpc.Region, "Expected VPC region to match")

	assert.Equal(t, "fd71:1140:a9d0::/52", vpc.IPv6[0].Range)

	assert.NotEmpty(t, vpc.Subnets, "Expected VPC to have subnets")
	assert.Equal(t, 456, vpc.Subnets[0].ID, "Expected subnet ID to match")
	assert.Equal(t, "subnet-1", vpc.Subnets[0].Label, "Expected subnet label to match")
	assert.Equal(t, "192.0.2.210/24", vpc.Subnets[0].IPv4, "Expected subnet IPv4 to match")

	assert.Equal(t, 111, vpc.Subnets[0].Linodes[0].ID)
	assert.Equal(t, true, vpc.Subnets[0].Linodes[0].Interfaces[0].Active)
	assert.Nil(t, vpc.Subnets[0].Linodes[0].Interfaces[0].ConfigID)
	assert.Equal(t, 421, vpc.Subnets[0].Linodes[0].Interfaces[0].ID)
}

func TestVPC_Update(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	updatedMockVPC := linodego.VPC{
		ID:          123,
		Label:       "updated-vpc",
		Description: "Updated description",
	}

	base.MockPut("vpcs/123", updatedMockVPC)

	opts := linodego.VPCUpdateOptions{
		Label:       "updated-vpc",
		Description: "Updated description",
	}

	vpc, err := base.Client.UpdateVPC(context.Background(), 123, opts)
	assert.NoError(t, err, "Expected no error when updating VPC")
	assert.NotNil(t, vpc, "Expected non-nil updated VPC")

	assertJSONObjectsSimilar(t, vpc, vpc.GetCreateOptions())
	assertJSONObjectsSimilar(t, vpc, vpc.GetUpdateOptions())

	assert.Equal(t, "updated-vpc", vpc.Label, "Expected VPC label to match")
	assert.Equal(t, "Updated description", vpc.Description, "Expected VPC description to match")
}

func TestVPC_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("vpcs/123", nil)

	err := base.Client.DeleteVPC(context.Background(), 123)
	assert.NoError(t, err, "Expected no error when deleting VPC")
}
