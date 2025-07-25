package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/linode/linodego"
)

func TestNodeBalancerConfigs_UDP(t *testing.T) {
	createFixture, err := fixtures.GetFixture("nodebalancers_configs_create_udp")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("nodebalancers/12345/configs", createFixture)

	opts := linodego.NodeBalancerConfigCreateOptions{
		Protocol:     linodego.Pointer(linodego.ProtocolUDP),
		Port:         1234,
		Algorithm:    linodego.Pointer(linodego.AlgorithmRingHash),
		Stickiness:   linodego.Pointer(linodego.StickinessSourceIP),
		UDPCheckPort: linodego.Pointer(12345),
	}

	config, err := base.Client.CreateNodeBalancerConfig(context.Background(), 12345, opts)
	require.NoError(t, err)

	require.Equal(t, linodego.AlgorithmRingHash, config.Algorithm)
	require.Equal(t, linodego.CheckHTTPBody, config.Check)
	require.Equal(t, 3, config.CheckAttempts)
	require.Equal(t, "it works", config.CheckBody)
	require.Equal(t, 90, config.CheckInterval)
	require.Equal(t, "/test", config.CheckPath)
	require.Equal(t, 10, config.CheckTimeout)
	require.Equal(t, linodego.CipherRecommended, config.CipherSuite)
	require.Equal(t, 4567, config.ID)
	require.Equal(t, 1234, config.NodeBalancerID)
	require.Equal(t, 0, config.NodesStatus.Down)
	require.Equal(t, 4, config.NodesStatus.Up)
	require.Equal(t, 1234, config.Port)
	require.Equal(t, linodego.ProtocolUDP, config.Protocol)
	require.Equal(t, "www.example.com", config.SSLCommonName)
	require.Equal(t, "00:01:02:03:04:05:06:07:08:09:0A:0B:0C:0D:0E:0F:10:11:12:13", config.SSLFingerprint)
	require.Equal(t, "<REDACTED>", config.SSLKey)
	require.Equal(t, linodego.StickinessSourceIP, config.Stickiness)
	require.Equal(t, 12345, config.UDPCheckPort)
	require.Equal(t, 10, config.UDPSessionTimeout)
}
