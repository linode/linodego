package unit

import (
	"context"
	"fmt"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices" 
	"testing"
)

func TestVLAN_List(t *testing.T) {
	// Load the fixture data for VLANs
	fixtureData, err := fixtures.GetFixture("vlans_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request
	base.MockGet("networking/vlans", fixtureData)

	vlans, err := base.Client.ListVLANs(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, vlans, "Expected non-empty VLAN list")

	// Use slices.IndexFunc to find the index of the specific VLAN
	index := slices.IndexFunc(vlans, func(v linodego.VLAN) bool {
		return v.Label == "test-vlan"
	})

	if index == -1 {
		t.Errorf("Expected VLAN 'test-vlan' to be in the response, but it was not found")
	} else {
		testVLAN := vlans[index]
		assert.Equal(t, "us-east", testVLAN.Region, "Expected region to be 'us-east'")
		assert.Contains(t, testVLAN.Linodes, 12345, "Expected Linodes to include 12345")
		assert.NotNil(t, testVLAN.Created, "Expected 'test-vlan' to have a created timestamp")
	}
}

func TestVLAN_GetIPAMAddress(t *testing.T) {
	// Load the fixture data for VLAN IPAM address
	fixtureData, err := fixtures.GetFixture("vlan_get_ipam_address")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	linodeID := 12345
	vlanLabel := "test-vlan"
	// Mock the GET request
	base.MockGet(fmt.Sprintf("linode/instances/%d/configs", linodeID), fixtureData)

	ipamAddress, err := base.Client.GetVLANIPAMAddress(context.Background(), linodeID, vlanLabel)
	assert.NoError(t, err)
	assert.NotEmpty(t, ipamAddress, "Expected non-empty IPAM address")

	// Verify the returned IPAM address
	assert.Equal(t, "10.0.0.1/24", ipamAddress, "Expected IPAM address to be '10.0.0.1/24'")
}
