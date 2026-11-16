package unit

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func Ptr[T any](v T) *T {
	return &v
}

func TestLKENodePool_Recycle(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "clusters/1234/pools/12345/recycle"), httpmock.NewStringResponder(200, "{}"))

	if err := client.RecycleLKENodePool(context.Background(), 1234, 12345); err != nil {
		t.Fatal(err)
	}
}

func TestLKENodePoolNode_Recycle(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "clusters/1234/nodes/abcde/recycle"), httpmock.NewStringResponder(200, "{}"))

	if err := client.RecycleLKENodePoolNode(context.Background(), 1234, "abcde"); err != nil {
		t.Fatal(err)
	}
}

func TestLKENodePoolNode_Get(t *testing.T) {
	fixtures := NewTestFixtures()

	fixtureData, err := fixtures.GetFixture("lke_node_pool_node_get")
	if err != nil {
		t.Fatalf("Failed to load fixture: %v", err)
	}

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("lke/clusters/1234/nodes/12345-abcde", fixtureData)

	node, err := base.Client.GetLKENodePoolNode(context.Background(), 1234, "12345-abcde")
	if err != nil {
		t.Fatalf("Error getting LKE node pool node: %v", err)
	}

	assert.Equal(t, "12345-abcde", node.ID)
	assert.Equal(t, 123456, node.InstanceID)
	assert.Equal(t, linodego.LKELinodeStatus("ready"), node.Status)
}

func TestLKENodePool_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_node_pool_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("lke/clusters/123/pools/456", fixtureData)

	nodePool, err := base.Client.GetLKENodePool(context.Background(), 123, 456)
	assert.NoError(t, err)
	assert.Equal(t, 456, nodePool.ID)
	assert.Equal(t, "g6-standard-2", nodePool.Type)
	assert.Equal(t, 3, nodePool.Count)
	assert.Equal(t, []string{"tag1", "tag2"}, nodePool.Tags)
	assert.Equal(t, []linodego.LockType{linodego.LockTypeCannotDelete}, nodePool.Locks)
	if assert.NotNil(t, nodePool.Isolation) {
		assert.True(t, nodePool.Isolation.PublicIPv4)
		assert.False(t, nodePool.Isolation.PublicIPv6)
	}
}

func TestLKENodePool_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_node_pool_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	label := "custom-label-create"

	createOptions := linodego.LKENodePoolCreateOptions{
		Count:  2,
		Type:   "g6-standard-2",
		Tags:   []string{"tag1"},
		Labels: map[string]string{"env": "dev"},
		Taints: []linodego.LKENodePoolTaint{
			{Key: "taintKey", Value: "taintValue", Effect: linodego.LKENodePoolTaintEffectNoSchedule},
		},
		Autoscaler: &linodego.LKENodePoolAutoscaler{
			Enabled: true,
			Min:     1,
			Max:     5,
		},
		Label: &label,
		Isolation: &linodego.LKENodePoolIsolationCreateOptions{
			PublicIPv4: Ptr(true),
			PublicIPv6: Ptr(false),
		},
	}

	base.MockPost("lke/clusters/123/pools", fixtureData)

	nodePool, err := base.Client.CreateLKENodePool(context.Background(), 123, createOptions)
	assert.NoError(t, err)
	assert.Equal(t, "g6-standard-2", nodePool.Type)
	assert.Equal(t, 2, nodePool.Count)
	assert.True(t, nodePool.Autoscaler.Enabled)
	assert.Equal(t, 1, nodePool.Autoscaler.Min)
	assert.Equal(t, 5, nodePool.Autoscaler.Max)
	assert.Equal(t, &label, nodePool.Label)
	if assert.NotNil(t, nodePool.Isolation) {
		assert.True(t, nodePool.Isolation.PublicIPv4)
		assert.False(t, nodePool.Isolation.PublicIPv6)
	}
}

func TestLKENodePool_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_node_pool_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	label := "custom-label-update"

	updateOptions := linodego.LKENodePoolUpdateOptions{
		Count:  5,
		Tags:   []string{"updated-tag"},
		Labels: Ptr(linodego.LKENodePoolLabels{"env": "prod"}),
		Autoscaler: &linodego.LKENodePoolAutoscaler{
			Enabled: true,
			Min:     2,
			Max:     8,
		},
		Label: &label,
	}

	base.MockPut("lke/clusters/123/pools/456", fixtureData)

	nodePool, err := base.Client.UpdateLKENodePool(context.Background(), 123, 456, updateOptions)
	assert.NoError(t, err)
	assert.Equal(t, 5, nodePool.Count)
	assert.Equal(t, []string{"updated-tag"}, nodePool.Tags)
	assert.True(t, nodePool.Autoscaler.Enabled)
	assert.Equal(t, 2, nodePool.Autoscaler.Min)
	assert.Equal(t, 8, nodePool.Autoscaler.Max)
	assert.Equal(t, &label, nodePool.Label)
}

func TestLKENodePool_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("lke/clusters/123/pools/456", nil)

	err := base.Client.DeleteLKENodePool(context.Background(), 123, 456)
	assert.NoError(t, err)
}

func TestLKENodePool_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_node_pool_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("lke/clusters/123/pools", fixtureData)

	nodePools, err := base.Client.ListLKENodePools(context.Background(), 123, nil)
	assert.NoError(t, err)
	assert.Len(t, nodePools, 2)
	assert.Equal(t, 456, nodePools[0].ID)
	assert.Equal(t, []linodego.LockType{linodego.LockTypeCannotDelete}, nodePools[0].Locks)
	assert.Equal(t, 789, nodePools[1].ID)
	assert.Empty(t, nodePools[1].Locks)
}

func TestLKENodePoolNode_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("lke/clusters/123/nodes/abc123", nil)

	err := base.Client.DeleteLKENodePoolNode(context.Background(), 123, "abc123")
	assert.NoError(t, err)
}

func TestLKEEnterpriseNodePool_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_e_node_pool_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	createOptions := linodego.LKENodePoolCreateOptions{
		Count:          2,
		Type:           "g6-standard-2",
		Tags:           []string{"testing"},
		K8sVersion:     linodego.Pointer("v1.31.1+lke1"),
		UpdateStrategy: linodego.Pointer(linodego.LKENodePoolOnRecycle),
		Label:          linodego.Pointer("custom-label-create"),
		FirewallID:     linodego.Pointer(12345),
	}

	base.MockPost("lke/clusters/123/pools", fixtureData)

	nodePool, err := base.Client.CreateLKENodePool(context.Background(), 123, createOptions)
	assert.NoError(t, err)
	assert.Equal(t, "g6-standard-2", nodePool.Type)
	assert.Equal(t, 2, nodePool.Count)
	assert.Equal(t, "v1.31.1+lke1", *nodePool.K8sVersion)
	assert.Equal(t, "on_recycle", string(*nodePool.UpdateStrategy))
	assert.Equal(t, linodego.Pointer("custom-label-create"), nodePool.Label)
	assert.Equal(t, 12345, *nodePool.FirewallID)
}

func TestLKENodePoolRDMA_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_node_pool_rdma_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	createOptions := linodego.LKENodePoolCreateOptions{
		Count:          6,
		Type:           "g3-gpu-rtxpro6000-blackwell-rdma-2",
		K8sVersion:     linodego.Pointer("v1.31.8+lke3"),
		UpdateStrategy: linodego.Pointer(linodego.LKENodePoolRollingUpdate),
		Label:          linodego.Pointer("rdma-pool"),
		FirewallID:     linodego.Pointer(12345),
		Networking: &linodego.LKENodePoolNetworking{
			Interfaces: []linodego.LKENodePoolInterface{
				{Type: linodego.LKENodePoolInterfaceTypeVPC, VPCID: 10001, SubnetID: 20001},
				{Type: linodego.LKENodePoolInterfaceTypeRDMAVPC, VPCID: 12345, SubnetID: 45678},
				{Type: linodego.LKENodePoolInterfaceTypeRDMAVPC, VPCID: 12345, SubnetID: 91012},
			},
		},
	}

	base.MockPost("lke/clusters/123/pools", fixtureData)

	nodePool, err := base.Client.CreateLKENodePool(context.Background(), 123, createOptions)
	assert.NoError(t, err)
	assert.Equal(t, "g3-gpu-rtxpro6000-blackwell-rdma-2", nodePool.Type)
	assert.Equal(t, 6, nodePool.Count)
	if assert.NotNil(t, nodePool.Networking) {
		ifaces := nodePool.Networking.Interfaces
		assert.Len(t, ifaces, 3)
		assert.Equal(t, linodego.LKENodePoolInterfaceTypeVPC, ifaces[0].Type)
		assert.Equal(t, 10001, ifaces[0].VPCID)
		assert.Equal(t, 20001, ifaces[0].SubnetID)
		assert.Equal(t, linodego.LKENodePoolInterfaceTypeRDMAVPC, ifaces[1].Type)
		assert.Equal(t, 12345, ifaces[1].VPCID)
		assert.Equal(t, 45678, ifaces[1].SubnetID)
		assert.Equal(t, linodego.LKENodePoolInterfaceTypeRDMAVPC, ifaces[2].Type)
		assert.Equal(t, 91012, ifaces[2].SubnetID)
	}
}

func TestLKENodePoolRDMA_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_node_pool_rdma_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("lke/clusters/123/pools/99001", fixtureData)

	nodePool, err := base.Client.GetLKENodePool(context.Background(), 123, 99001)
	assert.NoError(t, err)
	assert.Equal(t, 99001, nodePool.ID)
	assert.Equal(t, "g3-gpu-rtxpro6000-blackwell-rdma-2", nodePool.Type)
	if assert.NotNil(t, nodePool.Networking) {
		ifaces := nodePool.Networking.Interfaces
		assert.Len(t, ifaces, 3)
		assert.Equal(t, linodego.LKENodePoolInterfaceTypeVPC, ifaces[0].Type)
		assert.Equal(t, linodego.LKENodePoolInterfaceTypeRDMAVPC, ifaces[1].Type)
		assert.Equal(t, linodego.LKENodePoolInterfaceTypeRDMAVPC, ifaces[2].Type)
	}
}

func TestLKEEnterpriseNodePool_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_e_node_pool_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	updateOptions := linodego.LKENodePoolUpdateOptions{
		K8sVersion:     linodego.Pointer("v1.31.1+lke1"),
		UpdateStrategy: linodego.Pointer(linodego.LKENodePoolRollingUpdate),
		Label:          linodego.Pointer("custom-label-update"),
		FirewallID:     linodego.Pointer(12345),
	}

	base.MockPut("lke/clusters/123/pools/12345", fixtureData)

	nodePool, err := base.Client.UpdateLKENodePool(context.Background(), 123, 12345, updateOptions)
	assert.NoError(t, err)
	assert.Equal(t, "v1.31.1+lke1", *nodePool.K8sVersion)
	assert.Equal(t, "rolling_update", string(*nodePool.UpdateStrategy))
	assert.Equal(t, linodego.Pointer("custom-label-update"), nodePool.Label)
	assert.Equal(t, 12345, *nodePool.FirewallID)
}
