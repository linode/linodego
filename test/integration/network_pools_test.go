package integration

import (
    "context"
    "testing"
)

func TestIPv6Pool_List(t *testing.T) {
    client, teardown := createTestClient(t, "fixtures/TestIPv6Pool_List")
    defer teardown()

    ipv6Pools, err := client.ListIPv6Pools(context.Background(), nil)
    if err != nil {
        t.Fatalf("Error getting IPv6 Pools, expected struct, got error %v", err)
    }

    if len(ipv6Pools) == 0 {
        t.Fatalf("Expected to see IPv6 pools returned.")
    }

    if ipv6Pools[0].Range != "2600:3c00::/32" {
        t.Fatalf("Expected IPv6 pool range '2600:3c00::/32', got %s", ipv6Pools[0].Range)
    }
}

func TestIPv6Pool_Get(t *testing.T) {
    client, teardown := createTestClient(t, "fixtures/TestIPv6Pool_Get")
    defer teardown()

    ipv6Pool, err := client.GetIPv6Pool(context.Background(), "2600:3c00::/32")
    if err != nil {
        t.Fatalf("Error getting IPv6 Pool, expected struct, got error %v", err)
    }

    if ipv6Pool.Range != "2600:3c00::/32" {
        t.Fatalf("Expected IPv6 pool range '2600:3c00::/32', got %s", ipv6Pool.Range)
    }

    if ipv6Pool.Region != "us-east" {
        t.Fatalf("Expected IPv6 pool region 'us-east', got %s", ipv6Pool.Region)
    }
}
