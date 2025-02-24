package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestLKEClusterPool_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_cluster_pool_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("lke/clusters/123/pools", fixtureData)

	pools, err := base.Client.ListLKEClusterPools(context.Background(), 123, nil)
	assert.NoError(t, err)
	assert.Len(t, pools, 2)

	assert.Equal(t, 1, pools[0].ID)
	assert.Equal(t, "pool-1", pools[0].Type)
}

func TestLKEClusterPool_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_cluster_pool_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("lke/clusters/123/pools/1", fixtureData)

	pool, err := base.Client.GetLKEClusterPool(context.Background(), 123, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, pool.ID)
	assert.Equal(t, "pool-1", pool.Type)
}

func TestLKEClusterPool_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_cluster_pool_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	createOptions := linodego.LKEClusterPoolCreateOptions{
		Type:  "g6-standard-2",
		Count: 3,
	}

	base.MockPost("lke/clusters/123/pools", fixtureData)

	pool, err := base.Client.CreateLKEClusterPool(context.Background(), 123, createOptions)
	assert.NoError(t, err)
	assert.Equal(t, "g6-standard-2", pool.Type)
	assert.Equal(t, 3, pool.Count)
}

func TestLKEClusterPool_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_cluster_pool_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	updateOptions := linodego.LKEClusterPoolUpdateOptions{
		Count: 5,
	}

	base.MockPut("lke/clusters/123/pools/1", fixtureData)

	pool, err := base.Client.UpdateLKEClusterPool(context.Background(), 123, 1, updateOptions)
	assert.NoError(t, err)
	assert.Equal(t, 5, pool.Count)
}

func TestLKEClusterPool_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("lke/clusters/123/pools/1", nil)

	err := base.Client.DeleteLKEClusterPool(context.Background(), 123, 1)
	assert.NoError(t, err)
}

func TestLKEClusterPool_DeleteNode(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("lke/clusters/123/nodes/abc123", nil)

	err := base.Client.DeleteLKEClusterPoolNode(context.Background(), 123, "abc123")
	assert.NoError(t, err)
}
