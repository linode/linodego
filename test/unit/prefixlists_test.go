package unit

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestPrefixLists_List(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	response := map[string]any{
		"data": []map[string]any{
			{
				"id":                   321,
				"name":                 "pl:system:resolvers:us-iad:staging",
				"description":          "Resolver ACL",
				"visibility":           "restricted",
				"source_prefixlist_id": nil,
				"ipv4":                 []string{"139.144.192.62", "139.144.192.60"},
				"ipv6":                 []string{"2600:3c05:e001:bc::1", "2600:3c05:e001:bc::2"},
				"version":              4,
				"created":              "2018-01-01T00:01:01",
				"updated":              "2019-01-01T00:01:01",
				"deleted":              nil,
			},
		},
		"page":    1,
		"pages":   1,
		"results": 1,
	}

	base.MockGet("networking/prefixlists", response)

	prefixLists, err := base.Client.ListPrefixLists(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, prefixLists, 1)

	pl := prefixLists[0]
	assert.Equal(t, 321, pl.ID)
	assert.Equal(t, "pl:system:resolvers:us-iad:staging", pl.Name)
	assert.Equal(t, "restricted", pl.Visibility)
	if assert.NotNil(t, pl.IPv4) {
		assert.Equal(t, []string{"139.144.192.62", "139.144.192.60"}, *pl.IPv4)
	}
	if assert.NotNil(t, pl.IPv6) {
		assert.Equal(t, []string{"2600:3c05:e001:bc::1", "2600:3c05:e001:bc::2"}, *pl.IPv6)
	}
	if assert.NotNil(t, pl.Created) {
		assert.Equal(t, time.Date(2018, time.January, 1, 0, 1, 1, 0, time.UTC), pl.Created.UTC())
	}
	if assert.NotNil(t, pl.Updated) {
		assert.Equal(t, time.Date(2019, time.January, 1, 0, 1, 1, 0, time.UTC), pl.Updated.UTC())
	}
	assert.Nil(t, pl.Deleted)
}

func TestPrefixLists_Get(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	prefixListID := 654

	response := map[string]any{
		"id":                   prefixListID,
		"name":                 "pl::customer:example",
		"description":          "Example customer list",
		"visibility":           "account",
		"source_prefixlist_id": nil,
		"ipv4":                 []string{"198.51.100.0/24"},
		"ipv6":                 []string{"2001:db8::/32"},
		"version":              2,
		"created":              "2020-02-02T02:02:02",
		"updated":              "2020-03-03T03:03:03",
		"deleted":              nil,
	}

	base.MockGet(formatMockAPIPath("networking/prefixlists/%d", prefixListID), response)

	prefixList, err := base.Client.GetPrefixList(context.Background(), prefixListID)
	assert.NoError(t, err)
	assert.Equal(t, prefixListID, prefixList.ID)
	if assert.NotNil(t, prefixList.Created) {
		assert.Equal(t, time.Date(2020, time.February, 2, 2, 2, 2, 0, time.UTC), prefixList.Created.UTC())
	}
	if assert.NotNil(t, prefixList.Updated) {
		assert.Equal(t, time.Date(2020, time.March, 3, 3, 3, 3, 0, time.UTC), prefixList.Updated.UTC())
	}
}

func TestPrefixList_UnmarshalJSON(t *testing.T) {
	var prefixList linodego.PrefixList

	raw := []byte(`{
		"id": 1,
		"name": "pl:test",
		"visibility": "restricted",
		"ipv4": ["203.0.113.0/24"],
		"ipv6": ["2001:db8::/64"],
		"version": 5,
		"created": "2017-05-05T05:05:05",
		"updated": "2017-06-06T06:06:06",
		"deleted": null
	}`)

	assert.NoError(t, json.Unmarshal(raw, &prefixList))
	if assert.NotNil(t, prefixList.Created) {
		assert.Equal(t, time.Date(2017, time.May, 5, 5, 5, 5, 0, time.UTC), prefixList.Created.UTC())
	}
	if assert.NotNil(t, prefixList.Updated) {
		assert.Equal(t, time.Date(2017, time.June, 6, 6, 6, 6, 0, time.UTC), prefixList.Updated.UTC())
	}
	assert.Nil(t, prefixList.Deleted)
}
