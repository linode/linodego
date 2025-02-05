package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/linode/linodego"
)

func TestIPListIPv6Pools(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	mockResponse := struct {
		Data []linodego.IPv6Range `json:"data"`
	}{
		Data: []linodego.IPv6Range{
			{
				Range:       "2600:3c00::/64",
				Region:      "us-east",
				Prefix:      64,
				RouteTarget: "2600:3c00::1",
				IsBGP:       true,
				Linodes:     []int{12345, 67890},
			},
		},
	}

	base.MockGet("networking/ipv6/pools", mockResponse)

	pools, err := base.Client.ListIPv6Pools(context.Background(), nil)

	assert.NoError(t, err, "Expected no error when listing IPv6 pools")
	assert.NotNil(t, pools, "Expected non-nil IPv6 pools response")
	assert.Len(t, pools, 1, "Expected one IPv6 pool in response")
	assert.Equal(t, "2600:3c00::/64", pools[0].Range, "Expected matching IPv6 range")
	assert.Equal(t, "us-east", pools[0].Region, "Expected matching region")
	assert.Equal(t, 64, pools[0].Prefix, "Expected matching prefix length")
	assert.Equal(t, "2600:3c00::1", pools[0].RouteTarget, "Expected matching route target")
	assert.True(t, pools[0].IsBGP, "Expected IsBGP to be true")
	assert.ElementsMatch(t, []int{12345, 67890}, pools[0].Linodes, "Expected matching Linodes list")
}

func TestIPGetIPv6Pool(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	id := "1"
	mockResponse := linodego.IPv6Range{
		Range:       "2600:3c00::/64",
		Region:      "us-east",
		Prefix:      64,
		RouteTarget: "2600:3c00::1",
		IsBGP:       false,
		Linodes:     []int{54321},
	}
	base.MockGet("networking/ipv6/pools/"+id, mockResponse)

	pool, err := base.Client.GetIPv6Pool(context.Background(), id)
	assert.NoError(t, err, "Expected no error when getting IPv6 pool")
	assert.NotNil(t, pool, "Expected non-nil IPv6 pool response")
	assert.Equal(t, "2600:3c00::/64", pool.Range, "Expected matching IPv6 range")
	assert.Equal(t, "us-east", pool.Region, "Expected matching region")
	assert.Equal(t, 64, pool.Prefix, "Expected matching prefix length")
	assert.Equal(t, "2600:3c00::1", pool.RouteTarget, "Expected matching route target")
	assert.False(t, pool.IsBGP, "Expected IsBGP to be false")
	assert.ElementsMatch(t, []int{54321}, pool.Linodes, "Expected matching Linodes list")
}
