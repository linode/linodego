package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestIPv6Ranges_List(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock response for the GET request
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

	base.MockGet("networking/ipv6/ranges", mockResponse)

	ranges, err := base.Client.ListIPv6Ranges(context.Background(), nil)

	assert.NoError(t, err, "Expected no error when listing IPv6 ranges")
	assert.NotNil(t, ranges, "Expected non-nil IPv6 ranges response")
	assert.Len(t, ranges, 1, "Expected one IPv6 range in response")
	assert.Equal(t, "2600:3c00::/64", ranges[0].Range, "Expected matching IPv6 range")
}

func TestIPv6Range_Get(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	id := "2600:3c00::/64"

	// Mock response for the GET request
	mockResponse := struct {
		Range       string `json:"range"`
		Region      string `json:"region"`
		Prefix      int    `json:"prefix"`
		RouteTarget string `json:"route_target"`
		IsBGP       bool   `json:"is_bgp"`
		Linodes     []int  `json:"linodes"`
	}{
		Range:       id,
		Region:      "us-east",
		Prefix:      64,
		RouteTarget: "2600:3c00::1",
		IsBGP:       false,
		Linodes:     []int{54321},
	}

	base.MockGet("networking/ipv6/ranges/2600:3c00::%2F64", mockResponse)

	rangeObj, err := base.Client.GetIPv6Range(context.Background(), id)

	assert.NoError(t, err, "Expected no error when getting IPv6 range")
	assert.NotNil(t, rangeObj, "Expected non-nil IPv6 range response")
	assert.Equal(t, id, rangeObj.Range, "Expected matching IPv6 range")
	assert.Equal(t, "us-east", rangeObj.Region, "Expected matching region")
	assert.Equal(t, 64, rangeObj.Prefix, "Expected matching prefix length")
	assert.Equal(t, "2600:3c00::1", rangeObj.RouteTarget, "Expected matching route target")
	assert.False(t, rangeObj.IsBGP, "Expected IsBGP to be false")
	assert.ElementsMatch(t, []int{54321}, rangeObj.Linodes, "Expected matching Linodes list")
}

func TestIPv6Range_Create(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	createOpts := linodego.IPv6RangeCreateOptions{
		LinodeID:     linodego.Pointer(12345),
		PrefixLength: 64,
		RouteTarget:  linodego.Pointer("2600:3c00::1"),
	}

	// Mock the POST request
	base.MockPost("networking/ipv6/ranges", createOpts)

	createdRange, err := base.Client.CreateIPv6Range(context.Background(), createOpts)

	assert.NoError(t, err, "Expected no error when creating IPv6 range")
	assert.NotNil(t, createdRange, "Expected non-nil IPv6 range response")
	assert.Equal(t, *createOpts.RouteTarget, createdRange.RouteTarget, "Expected matching route target")
}

func TestIPv6Range_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	id := "2600:3c00::/64"

	base.MockDelete("networking/ipv6/ranges/2600:3c00::%2F64", nil)

	err := base.Client.DeleteIPv6Range(context.Background(), id)

	assert.NoError(t, err, "Expected no error when deleting IPv6 range")
}
