package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"

	"github.com/jarcoal/httpmock"
)

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
