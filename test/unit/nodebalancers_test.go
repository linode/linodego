package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestNodeBalancers_UDP(t *testing.T) {
	createFixture, err := fixtures.GetFixture("nodebalancers_create_udp")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("nodebalancers", createFixture)

	opts := linodego.NodeBalancerCreateOptions{
		Label:                 linodego.Pointer("foobar"),
		Region:                linodego.Pointer("us-mia"),
		ClientUDPSessThrottle: linodego.Pointer(5),
		Configs: []*linodego.NodeBalancerConfigCreateOptions{
			{
				Protocol:     linodego.Pointer(linodego.ProtocolUDP),
				Port:         1234,
				Algorithm:    linodego.Pointer(linodego.AlgorithmRingHash),
				Stickiness:   linodego.Pointer(linodego.StickinessSourceIP),
				UDPCheckPort: linodego.Pointer(80),
			},
		},
	}

	nb, err := base.Client.CreateNodeBalancer(context.Background(), opts)
	require.NoError(t, err)

	require.Equal(t, 0, nb.ClientConnThrottle)
	require.Equal(t, 10, nb.ClientUDPSessThrottle)

	require.Equal(t, "192.0.2.1.ip.linodeusercontent.com", *nb.Hostname)
	require.Equal(t, 12345, nb.ID)
	require.Equal(t, "203.0.113.1", *nb.IPv4)
	require.Nil(t, nb.IPv6)
	require.Equal(t, "balancer12345", *nb.Label)
	require.Equal(t, "us-mia", nb.Region)
	require.Equal(t, "example tag", nb.Tags[0])
	require.Equal(t, "another example", nb.Tags[1])
	require.NotZero(t, nb.Transfer.In)
	require.NotZero(t, nb.Transfer.Out)
	require.NotZero(t, nb.Transfer.Total)

	require.Equal(
		t,
		time.Date(2018, 1, 1, 0, 1, 1, 0, time.UTC),
		*nb.Created,
	)

	require.Equal(
		t,
		time.Date(2018, 3, 1, 0, 1, 1, 0, time.UTC),
		*nb.Updated,
	)
}
