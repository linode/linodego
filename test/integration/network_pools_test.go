package integration

import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestIPv6Pool_List(t *testing.T) {
    client, teardown := createTestClient(t, "fixtures/TestIPv6Pool_List")
    defer teardown()

    ipv6Pools, err := client.ListIPv6Pools(context.Background(), nil)
    require.NoError(t, err, "Error getting IPv6 Pools, expected struct")

    require.NotEmpty(t, ipv6Pools, "Expected to see IPv6 pools returned")

    require.Equal(t, "2600:3c00::/32", ipv6Pools[0].Range, "Expected IPv6 pool range '2600:3c00::/32'")
}

func TestIPv6Pool_Get(t *testing.T) {
    client, teardown := createTestClient(t, "fixtures/TestIPv6Pool_Get")
    defer teardown()

    ipv6Pool, err := client.GetIPv6Pool(context.Background(), "2600:3c00::/32")
    require.NoError(t, err, "Error getting IPv6 Pool, expected struct")

    require.Equal(t, "2600:3c00::/32", ipv6Pool.Range, "Expected IPv6 pool range '2600:3c00::/32'")
    require.Equal(t, "us-east", ipv6Pool.Region, "Expected IPv6 pool region 'us-east'")
}
