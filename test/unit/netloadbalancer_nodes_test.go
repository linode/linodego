package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestNetLoadBalancerNode_Create(t *testing.T) {
	tests := []struct {
		name              string
		netloadbalancerID int
		listenerID        int
		createOpts        linodego.NetLoadBalancerNodeCreateOptions
		fixture           string
		expectedLabel     string
		expectedAddressV6 string
		expectedWeight    int
	}{
		{
			name:              "basic node creation",
			netloadbalancerID: 123,
			listenerID:        456,
			createOpts: linodego.NetLoadBalancerNodeCreateOptions{
				Label:     "Web Server 1",
				AddressV6: "2600:3c03::f03c:91ff:fe24:abcd",
				Weight:    100,
			},
			fixture:           "netloadbalancer_node_create",
			expectedLabel:     "Web Server 1",
			expectedAddressV6: "2600:3c03::f03c:91ff:fe24:abcd",
			expectedWeight:    100,
		},
		{
			name:              "node creation with default weight",
			netloadbalancerID: 123,
			listenerID:        456,
			createOpts: linodego.NetLoadBalancerNodeCreateOptions{
				Label:     "Database Server",
				AddressV6: "2600:3c03::f03c:91ff:fe24:1234",
				Weight:    0, // Default weight
			},
			fixture:           "netloadbalancer_node_create_default_weight",
			expectedLabel:     "Database Server",
			expectedAddressV6: "2600:3c03::f03c:91ff:fe24:1234",
			expectedWeight:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixtureData, err := fixtures.GetFixture(tt.fixture)
			assert.NoError(t, err)

			var base ClientBaseCase
			base.SetUp(t)
			defer base.TearDown(t)

			expectedPath := "netloadbalancers/123/listeners/456/nodes"
			base.MockPost(expectedPath, fixtureData)

			node, err := base.Client.CreateNetLoadBalancerNode(
				context.Background(),
				tt.netloadbalancerID,
				tt.listenerID,
				tt.createOpts,
			)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedLabel, node.Label)
			assert.Equal(t, tt.expectedAddressV6, node.AddressV6)
			assert.Equal(t, tt.expectedWeight, node.Weight)
		})
	}
}

func TestNetLoadBalancerNode_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("netloadbalancer_node_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("netloadbalancers/123/listeners/456/nodes/789", fixtureData)

	node, err := base.Client.GetNetLoadBalancerNode(context.Background(), 123, 456, 789)
	assert.NoError(t, err)

	assert.Equal(t, 789, node.ID)
	assert.Equal(t, 12345, node.LinodeID)
	assert.Equal(t, "Production Web Server", node.Label)
	assert.Equal(t, "2600:3c03::f03c:91ff:fe24:prod", node.AddressV6)
	assert.Equal(t, 100, node.Weight)
}

func TestNetLoadBalancerNode_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("netloadbalancer_nodes_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("netloadbalancers/123/listeners/456/nodes", fixtureData)

	nodes, err := base.Client.ListNetLoadBalancerNodes(context.Background(), 123, 456, &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.Len(t, nodes, 2)

	// Verify details of the first node
	assert.Equal(t, 789, nodes[0].ID)
	assert.Equal(t, 11111, nodes[0].LinodeID)
	assert.Equal(t, "Web Server 1", nodes[0].Label)
	assert.Equal(t, "2600:3c03::f03c:91ff:fe24:0001", nodes[0].AddressV6)
	assert.Equal(t, 100, nodes[0].Weight)

	// Verify details of the second node
	assert.Equal(t, 790, nodes[1].ID)
	assert.Equal(t, 22222, nodes[1].LinodeID)
	assert.Equal(t, "Web Server 2", nodes[1].Label)
	assert.Equal(t, "2600:3c03::f03c:91ff:fe24:0002", nodes[1].AddressV6)
	assert.Equal(t, 50, nodes[1].Weight)
}

func TestNetLoadBalancerNode_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("netloadbalancer_node_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("netloadbalancers/123/listeners/456/nodes/789", fixtureData)

	updateOpts := linodego.NetLoadBalancerNodeLabelUpdateOptions{
		Label: "Updated Web Server",
	}

	node, err := base.Client.UpdateNetLoadBalancerNode(context.Background(), 123, 456, 789, updateOpts)
	assert.NoError(t, err)

	assert.Equal(t, 789, node.ID)
	assert.Equal(t, "Updated Web Server", node.Label)
}

func TestNetLoadBalancerNode_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("netloadbalancers/123/listeners/456/nodes/789", nil)

	err := base.Client.DeleteNetLoadBalancerNode(context.Background(), 123, 456, 789)
	assert.NoError(t, err)
}

func TestNetLoadBalancerNode_GetCreateOptions(t *testing.T) {
	tests := []struct {
		name     string
		node     linodego.NetLoadBalancerNode
		expected linodego.NetLoadBalancerNodeCreateOptions
	}{
		{
			name: "basic conversion",
			node: linodego.NetLoadBalancerNode{
				ID:        789,
				LinodeID:  12345,
				Label:     "Test Node",
				AddressV6: "2600:3c03::f03c:91ff:fe24:test",
				Weight:    75,
			},
			expected: linodego.NetLoadBalancerNodeCreateOptions{
				Label:     "Test Node",
				AddressV6: "2600:3c03::f03c:91ff:fe24:test",
				Weight:    75,
			},
		},
		{
			name: "conversion with zero weight",
			node: linodego.NetLoadBalancerNode{
				ID:        456,
				LinodeID:  67890,
				Label:     "Zero Weight Node",
				AddressV6: "2600:3c03::f03c:91ff:fe24:zero",
				Weight:    0,
			},
			expected: linodego.NetLoadBalancerNodeCreateOptions{
				Label:     "Zero Weight Node",
				AddressV6: "2600:3c03::f03c:91ff:fe24:zero",
				Weight:    0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.node.GetCreateOptions()

			assert.Equal(t, tt.expected.Label, result.Label)
			assert.Equal(t, tt.expected.AddressV6, result.AddressV6)
			assert.Equal(t, tt.expected.Weight, result.Weight)
		})
	}
}

func TestNetLoadBalancerNode_GetUpdateOptions(t *testing.T) {
	tests := []struct {
		name     string
		node     linodego.NetLoadBalancerNode
		expected linodego.NetLoadBalancerNodeUpdateOptions
	}{
		{
			name: "basic conversion",
			node: linodego.NetLoadBalancerNode{
				ID:        789,
				LinodeID:  12345,
				Label:     "Updated Node",
				AddressV6: "2600:3c03::f03c:91ff:fe24:updt",
				Weight:    150,
			},
			expected: linodego.NetLoadBalancerNodeUpdateOptions{
				Label:     "Updated Node",
				AddressV6: "2600:3c03::f03c:91ff:fe24:updt",
				Weight:    150,
			},
		},
		{
			name: "conversion with changed weight",
			node: linodego.NetLoadBalancerNode{
				ID:        321,
				LinodeID:  54321,
				Label:     "Reweighted Node",
				AddressV6: "2600:3c03::f03c:91ff:fe24:rwgt",
				Weight:    25,
			},
			expected: linodego.NetLoadBalancerNodeUpdateOptions{
				Label:     "Reweighted Node",
				AddressV6: "2600:3c03::f03c:91ff:fe24:rwgt",
				Weight:    25,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.node.GetUpdateOptions()

			assert.Equal(t, tt.expected.Label, result.Label)
			assert.Equal(t, tt.expected.AddressV6, result.AddressV6)
			assert.Equal(t, tt.expected.Weight, result.Weight)
		})
	}
}

func TestNetLoadBalancerNode_ListWithOptions(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("netloadbalancer_nodes_list_filtered")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("netloadbalancers/123/listeners/456/nodes", fixtureData)

	listOpts := &linodego.ListOptions{
		PageOptions: &linodego.PageOptions{
			Page:    1,
			Pages:   1,
			Results: 1,
		},
		Filter: "{\"weight\":{\"gt\":50}}",
	}

	nodes, err := base.Client.ListNetLoadBalancerNodes(context.Background(), 123, 456, listOpts)
	assert.NoError(t, err)

	assert.Len(t, nodes, 1)
	assert.Greater(t, nodes[0].Weight, 50)
}

func TestNetLoadBalancerNode_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected linodego.NetLoadBalancerNode
		wantErr  bool
	}{
		{
			name: "valid JSON with all fields",
			jsonData: `{
				"id": 999,
				"linode_id": 88888,
				"label": "Test JSON Node",
				"address_v6": "2600:3c03::f03c:91ff:fe24:json",
				"weight": 125,
				"created": "2025-03-01T10:00:00",
				"updated": "2025-03-15T14:30:00",
				"weight_updated": "2025-03-20T09:15:00"
			}`,
			expected: linodego.NetLoadBalancerNode{
				ID:        999,
				LinodeID:  88888,
				Label:     "Test JSON Node",
				AddressV6: "2600:3c03::f03c:91ff:fe24:json",
				Weight:    125,
			},
			wantErr: false,
		},
		{
			name: "valid JSON with minimal fields",
			jsonData: `{
				"id": 111,
				"linode_id": 0,
				"label": "",
				"address_v6": "",
				"weight": 0
			}`,
			expected: linodego.NetLoadBalancerNode{
				ID:        111,
				LinodeID:  0,
				Label:     "",
				AddressV6: "",
				Weight:    0,
			},
			wantErr: false,
		},
		{
			name:     "invalid JSON",
			jsonData: `{"id": "not-a-number"}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var node linodego.NetLoadBalancerNode
			err := node.UnmarshalJSON([]byte(tt.jsonData))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.ID, node.ID)
			assert.Equal(t, tt.expected.LinodeID, node.LinodeID)
			assert.Equal(t, tt.expected.Label, node.Label)
			assert.Equal(t, tt.expected.AddressV6, node.AddressV6)
			assert.Equal(t, tt.expected.Weight, node.Weight)

			// Test time fields if they were set in the JSON
			if tt.expected.ID == 999 {
				assert.NotNil(t, node.Created)
				assert.NotNil(t, node.Updated)
				assert.NotNil(t, node.WeightUpdated)
			}
		})
	}
}

func TestNetLoadBalancerNode_EdgeCases(t *testing.T) {
	t.Run("create node with high weight", func(t *testing.T) {
		fixtureData, err := fixtures.GetFixture("netloadbalancer_node_create_high_weight")
		assert.NoError(t, err)

		var base ClientBaseCase
		base.SetUp(t)
		defer base.TearDown(t)

		base.MockPost("netloadbalancers/123/listeners/456/nodes", fixtureData)

		createOpts := linodego.NetLoadBalancerNodeCreateOptions{
			Label:     "High Weight Node",
			AddressV6: "2600:3c03::f03c:91ff:fe24:high",
			Weight:    999,
		}

		node, err := base.Client.CreateNetLoadBalancerNode(context.Background(), 123, 456, createOpts)
		assert.NoError(t, err)
		assert.Equal(t, 999, node.Weight)
	})

	t.Run("update node with empty label should fail validation", func(t *testing.T) {
		var base ClientBaseCase
		base.SetUp(t)
		defer base.TearDown(t)

		// This should still work at the client level, validation happens server-side
		base.MockPut("netloadbalancers/123/listeners/456/nodes/789", map[string]interface{}{
			"id":         789,
			"label":      "",
			"linode_id":  12345,
			"address_v6": "2600:3c03::f03c:91ff:fe24:empty",
			"weight":     100,
		})

		updateOpts := linodego.NetLoadBalancerNodeLabelUpdateOptions{
			Label: "", // Empty label
		}

		node, err := base.Client.UpdateNetLoadBalancerNode(context.Background(), 123, 456, 789, updateOpts)
		assert.NoError(t, err) // Client doesn't validate, server would
		assert.Equal(t, "", node.Label)
	})
}
