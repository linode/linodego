package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectStorageCluster_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_cluster_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("object-storage/clusters", fixtureData)

	clusters, err := base.Client.ListObjectStorageClusters(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, clusters, 1)

	assert.Equal(t, "us-east-1", clusters[0].Region)
	assert.Equal(t, "active", clusters[0].Status)
	assert.Equal(t, "my-cluster-id", clusters[0].ID)
	assert.Equal(t, "example.com", clusters[0].Domain)
	assert.Equal(t, "static.example.com", clusters[0].StaticSiteDomain)
}

func TestObjectStorageCluster_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_cluster_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("object-storage/clusters/my-cluster-id", fixtureData)

	cluster, err := base.Client.GetObjectStorageCluster(context.Background(), "my-cluster-id")
	assert.NoError(t, err)
	assert.NotNil(t, cluster)

	assert.Equal(t, "us-east-1", cluster.Region)
	assert.Equal(t, "active", cluster.Status)
	assert.Equal(t, "my-cluster-id", cluster.ID)
	assert.Equal(t, "example.com", cluster.Domain)
	assert.Equal(t, "static.example.com", cluster.StaticSiteDomain)
}
