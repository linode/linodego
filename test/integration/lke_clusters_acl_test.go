package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/linode/linodego"
)

func TestLKECluster_withACL(t *testing.T) {
	client, cluster, teardown, err := setupLKECluster(
		t,
		[]clusterModifier{
			func(options *linodego.LKEClusterCreateOptions) {
				options.ControlPlane = &linodego.LKEClusterControlPlane{
					ACL: &linodego.LKEClusterControlPlaneACL{
						Enabled: true,
						Addresses: &linodego.LKEClusterControlPlaneACLAddresses{
							IPv4: []string{"10.0.0.1/32"},
							IPv6: []string{"1234::5678"},
						},
					},
				}
			},
		},
		"fixtures/TestLKECluster_withACL",
	)
	require.NoError(t, err)
	defer teardown()

	// TODO: Not currently populated in response, uncomment when available
	// require.Equal(t, true, cluster.ControlPlane.ACL.Enabled)
	// require.Equal(t, "10.0.0.1/32", cluster.ControlPlane.ACL.Addresses.IPv4[0])
	// require.Equal(t, "1234::5678", cluster.ControlPlane.ACL.Addresses.IPv6[0])

	acl, err := client.GetLKEClusterControlPlaneACL(context.Background(), cluster.ID)
	assert.NoError(t, err)

	require.Equal(t, true, acl.ACL.Enabled)
	require.Equal(t, "10.0.0.1/32", acl.ACL.Addresses.IPv4[0])
	require.Equal(t, "1234::5678/128", acl.ACL.Addresses.IPv6[0])

	acl, err = client.UpdateLKEClusterControlPlaneACL(
		context.Background(),
		cluster.ID,
		linodego.LKEClusterControlPlaneACLUpdateOptions{
			ACL: linodego.LKEClusterControlPlaneACL{
				Enabled: true,
				Addresses: linodego.LKEClusterControlPlaneACLAddresses{
					IPv4: []string{"10.0.0.2/32"},
					IPv6: []string{},
				},
			},
		},
	)
	require.NoError(t, err)

	require.Equal(t, true, acl.ACL.Enabled)
	require.Equal(t, "10.0.0.2/32", acl.ACL.Addresses.IPv4[0])
	require.Equal(t, 0, len(acl.ACL.Addresses.IPv6))

	err = client.DeleteLKEClusterControlPlaneACL(context.Background(), cluster.ID)
	require.NoError(t, err)

	acl, err = client.GetLKEClusterControlPlaneACL(context.Background(), cluster.ID)
	assert.NoError(t, err)

	assert.Equal(t, false, acl.ACL.Enabled)
}
