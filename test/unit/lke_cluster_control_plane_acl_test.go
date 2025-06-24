package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func BoolToPtr(b bool) *bool {
	return &b
}

func TestLKEClusterControlPlaneACL_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_cluster_control_plane_acl_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("lke/clusters/123/control_plane_acl", fixtureData)

	acl, err := base.Client.GetLKEClusterControlPlaneACL(context.Background(), 123)
	assert.NoError(t, err)
	assert.True(t, acl.ACL.Enabled)
	assert.Equal(t, "rev-abc123", acl.ACL.RevisionID)
	assert.Equal(t, []string{"192.168.1.1/32"}, acl.ACL.Addresses.IPv4)
	assert.Equal(t, []string{"2001:db8::/32"}, acl.ACL.Addresses.IPv6)
}

func TestLKEClusterControlPlaneACL_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_cluster_control_plane_acl_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	updateOptions := linodego.LKEClusterControlPlaneACLUpdateOptions{
		ACL: linodego.LKEClusterControlPlaneACLOptions{
			Enabled:    BoolToPtr(true),
			RevisionID: "rev-abc123",
			Addresses: &linodego.LKEClusterControlPlaneACLAddressesOptions{
				IPv4: []string{"10.0.0.1/32"},
				IPv6: []string{"2001:db8::/64"},
			},
		},
	}

	base.MockPut("lke/clusters/123/control_plane_acl", fixtureData)

	acl, err := base.Client.UpdateLKEClusterControlPlaneACL(context.Background(), 123, updateOptions)
	assert.NoError(t, err)
	assert.True(t, acl.ACL.Enabled)
	assert.Equal(t, "rev-abc124", acl.ACL.RevisionID)
	assert.Equal(t, []string{"10.0.0.1/32"}, acl.ACL.Addresses.IPv4)
	assert.Equal(t, []string{"2001:db8::/64"}, acl.ACL.Addresses.IPv6)
}

func TestLKEClusterControlPlaneACL_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("lke/clusters/123/control_plane_acl", nil)

	err := base.Client.DeleteLKEClusterControlPlaneACL(context.Background(), 123)
	assert.NoError(t, err)
}
