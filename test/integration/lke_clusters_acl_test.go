package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/linode/linodego"
)

func TestLKECluster_withACL(t *testing.T) {
	valueTrue := true

	client, cluster, teardown, err := setupLKECluster(
		t,
		[]clusterModifier{
			func(options *linodego.LKEClusterCreateOptions) {
				options.ControlPlane = &linodego.LKEClusterControlPlaneOptions{
					ACL: &linodego.LKEClusterControlPlaneACLOptions{
						Enabled: &valueTrue,
						Addresses: &linodego.LKEClusterControlPlaneACLAddressesOptions{
							IPv4: &[]string{"10.0.0.1/32"},
							IPv6: &[]string{"1234::5678"},
						},
					},
				}
			},
		},
		"fixtures/TestLKECluster_withACL",
	)
	require.NoError(t, err)
	defer teardown()

	acl, err := client.GetLKEClusterControlPlaneACL(context.Background(), cluster.ID)
	assert.NoError(t, err)

	require.Equal(t, true, acl.ACL.Enabled)
	require.Equal(t, "10.0.0.1/32", acl.ACL.Addresses.IPv4[0])
	require.Equal(t, "1234::5678/128", acl.ACL.Addresses.IPv6[0])

	testRevisionID := "test-revision-id"

	acl, err = client.UpdateLKEClusterControlPlaneACL(
		context.Background(),
		cluster.ID,
		linodego.LKEClusterControlPlaneACLUpdateOptions{
			ACL: linodego.LKEClusterControlPlaneACLOptions{
				Enabled: &valueTrue,
				Addresses: &linodego.LKEClusterControlPlaneACLAddressesOptions{
					IPv4: &[]string{"10.0.0.2/32"},
					IPv6: &[]string{},
				},
				RevisionID: testRevisionID,
			},
		},
	)
	require.NoError(t, err)

	require.Equal(t, true, acl.ACL.Enabled)
	require.Equal(t, "10.0.0.2/32", acl.ACL.Addresses.IPv4[0])
	require.Equal(t, 0, len(acl.ACL.Addresses.IPv6))
	assert.Equal(t, testRevisionID, acl.ACL.RevisionID)

	err = client.DeleteLKEClusterControlPlaneACL(context.Background(), cluster.ID)
	require.NoError(t, err)

	acl, err = client.GetLKEClusterControlPlaneACL(context.Background(), cluster.ID)
	assert.NoError(t, err)

	assert.Equal(t, false, acl.ACL.Enabled)
}
