package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestNetLoadBalancerListener_Create(t *testing.T) {
	tests := []struct {
		name              string
		netloadbalancerID int
		createOpts        linodego.NetLoadBalancerListenerCreateOptions
		fixture           string
		expectedProtocol  string
		expectedPort      int
		expectedLabel     string
	}{
		{
			name:              "basic listener creation",
			netloadbalancerID: 123,
			createOpts: linodego.NetLoadBalancerListenerCreateOptions{
				Protocol: "tcp",
				Port:     80,
				Label:    "HTTP Listener",
			},
			fixture:          "netloadbalancer_listener_create",
			expectedProtocol: "tcp",
			expectedPort:     80,
			expectedLabel:    "HTTP Listener",
		},
		{
			name:              "HTTPS listener creation",
			netloadbalancerID: 123,
			createOpts: linodego.NetLoadBalancerListenerCreateOptions{
				Protocol: "tcp",
				Port:     443,
				Label:    "HTTPS Listener",
			},
			fixture:          "netloadbalancer_listener_create_https",
			expectedProtocol: "tcp",
			expectedPort:     443,
			expectedLabel:    "HTTPS Listener",
		},
		{
			name:              "listener creation with nodes",
			netloadbalancerID: 456,
			createOpts: linodego.NetLoadBalancerListenerCreateOptions{
				Protocol: "tcp",
				Port:     8080,
				Label:    "App Listener",
				Nodes: []linodego.NetLoadBalancerNodeCreateOptions{
					{
						Label:     "App Server 1",
						AddressV6: "2600:3c03::f03c:91ff:fe24:app1",
						Weight:    100,
					},
					{
						Label:     "App Server 2",
						AddressV6: "2600:3c03::f03c:91ff:fe24:app2",
						Weight:    100,
					},
				},
			},
			fixture:          "netloadbalancer_listener_create_with_nodes",
			expectedProtocol: "tcp",
			expectedPort:     8080,
			expectedLabel:    "App Listener",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixtureData, err := fixtures.GetFixture(tt.fixture)
			assert.NoError(t, err)

			var base ClientBaseCase
			base.SetUp(t)
			defer base.TearDown(t)

			expectedPath := "netloadbalancers/123/listeners"
			if tt.netloadbalancerID != 123 {
				expectedPath = "netloadbalancers/456/listeners"
			}
			base.MockPost(expectedPath, fixtureData)

			listener, err := base.Client.CreateNetLoadBalancerListener(
				context.Background(),
				tt.netloadbalancerID,
				tt.createOpts,
			)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedProtocol, listener.Protocol)
			assert.Equal(t, tt.expectedPort, listener.Port)
			assert.Equal(t, tt.expectedLabel, listener.Label)
		})
	}
}

func TestNetLoadBalancerListener_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("netloadbalancer_listener_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("netloadbalancers/123/listeners/456", fixtureData)

	listener, err := base.Client.GetNetLoadBalancerListener(context.Background(), 123, 456)
	assert.NoError(t, err)

	assert.Equal(t, 456, listener.ID)
	assert.Equal(t, "tcp", listener.Protocol)
	assert.Equal(t, 80, listener.Port)
	assert.Equal(t, "Production HTTP Listener", listener.Label)
}

func TestNetLoadBalancerListener_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("netloadbalancer_listeners_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("netloadbalancers/123/listeners", fixtureData)

	listeners, err := base.Client.ListNetLoadBalancerListeners(context.Background(), 123, &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.Len(t, listeners, 2)

	// Verify details of the first listener
	assert.Equal(t, 456, listeners[0].ID)
	assert.Equal(t, "tcp", listeners[0].Protocol)
	assert.Equal(t, 80, listeners[0].Port)
	assert.Equal(t, "HTTP Listener", listeners[0].Label)

	// Verify details of the second listener
	assert.Equal(t, 457, listeners[1].ID)
	assert.Equal(t, "tcp", listeners[1].Protocol)
	assert.Equal(t, 443, listeners[1].Port)
	assert.Equal(t, "HTTPS Listener", listeners[1].Label)
}

func TestNetLoadBalancerListener_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("netloadbalancer_listener_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("netloadbalancers/123/listeners/456", fixtureData)

	updateOpts := linodego.NetLoadBalancerListenerUpdateOptions{
		Label: "Updated HTTP Listener",
	}

	listener, err := base.Client.UpdateNetLoadBalancerListener(context.Background(), 123, 456, updateOpts)
	assert.NoError(t, err)

	assert.Equal(t, 456, listener.ID)
	assert.Equal(t, "Updated HTTP Listener", listener.Label)
}

func TestNetLoadBalancerListener_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("netloadbalancers/123/listeners/456", nil)

	err := base.Client.DeleteNetLoadBalancerListener(context.Background(), 123, 456)
	assert.NoError(t, err)
}

func TestNetLoadBalancerListener_GetCreateOptions(t *testing.T) {
	tests := []struct {
		name     string
		listener linodego.NetLoadBalancerListener
		expected linodego.NetLoadBalancerListenerCreateOptions
	}{
		{
			name: "basic conversion",
			listener: linodego.NetLoadBalancerListener{
				ID:       456,
				Protocol: "tcp",
				Port:     80,
				Label:    "Test Listener",
			},
			expected: linodego.NetLoadBalancerListenerCreateOptions{
				Protocol: "tcp",
				Port:     80,
				Label:    "Test Listener",
			},
		},
		{
			name: "HTTPS listener conversion",
			listener: linodego.NetLoadBalancerListener{
				ID:       789,
				Protocol: "tcp",
				Port:     443,
				Label:    "SSL Listener",
			},
			expected: linodego.NetLoadBalancerListenerCreateOptions{
				Protocol: "tcp",
				Port:     443,
				Label:    "SSL Listener",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.listener.GetCreateOptions()

			assert.Equal(t, tt.expected.Protocol, result.Protocol)
			assert.Equal(t, tt.expected.Port, result.Port)
			assert.Equal(t, tt.expected.Label, result.Label)
			// Note: GetCreateOptions doesn't include nodes as per the comment in the source
			assert.Nil(t, result.Nodes)
		})
	}
}

func TestNetLoadBalancerListener_GetUpdateOptions(t *testing.T) {
	tests := []struct {
		name     string
		listener linodego.NetLoadBalancerListener
		expected linodego.NetLoadBalancerListenerUpdateOptions
	}{
		{
			name: "basic conversion",
			listener: linodego.NetLoadBalancerListener{
				ID:       456,
				Protocol: "tcp",
				Port:     80,
				Label:    "Updated Listener",
			},
			expected: linodego.NetLoadBalancerListenerUpdateOptions{
				Label: "Updated Listener",
			},
		},
		{
			name: "conversion with different label",
			listener: linodego.NetLoadBalancerListener{
				ID:       789,
				Protocol: "tcp",
				Port:     8080,
				Label:    "Custom App Listener",
			},
			expected: linodego.NetLoadBalancerListenerUpdateOptions{
				Label: "Custom App Listener",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.listener.GetUpdateOptions()

			assert.Equal(t, tt.expected.Label, result.Label)
			// Note: GetUpdateOptions only includes the label as per the source implementation
		})
	}
}

func TestNetLoadBalancerListener_UpdateNodeWeights(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the POST request for node weights update (no response body expected)
	base.MockPost("netloadbalancers/123/listeners/456/node-weights", nil)

	updateOpts := linodego.NetLoadBalancerListenerNodeWeightsUpdateOptions{
		Nodes: []linodego.NetLoadBalancerListenerNodeWeightUpdateOptions{
			{
				ID:     789,
				Weight: 150,
			},
			{
				ID:     790,
				Weight: 75,
			},
		},
	}

	err := base.Client.UpdateNetLoadBalancerListenerNodeWeights(context.Background(), 123, 456, updateOpts)
	assert.NoError(t, err)
}

func TestNetLoadBalancerListener_ListWithOptions(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("netloadbalancer_listeners_list_filtered")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("netloadbalancers/123/listeners", fixtureData)

	listOpts := &linodego.ListOptions{
		PageOptions: &linodego.PageOptions{
			Page:    1,
			Pages:   1,
			Results: 1,
		},
		Filter: "{\"port\":80}",
	}

	listeners, err := base.Client.ListNetLoadBalancerListeners(context.Background(), 123, listOpts)
	assert.NoError(t, err)

	assert.Len(t, listeners, 1)
	assert.Equal(t, 80, listeners[0].Port)
}

func TestNetLoadBalancerListener_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected linodego.NetLoadBalancerListener
		wantErr  bool
	}{
		{
			name: "valid JSON with all fields",
			jsonData: `{
				"id": 999,
				"protocol": "tcp",
				"port": 8080,
				"label": "Test JSON Listener",
				"created": "2025-03-01T10:00:00",
				"updated": "2025-03-15T14:30:00"
			}`,
			expected: linodego.NetLoadBalancerListener{
				ID:       999,
				Protocol: "tcp",
				Port:     8080,
				Label:    "Test JSON Listener",
			},
			wantErr: false,
		},
		{
			name: "valid JSON with minimal fields",
			jsonData: `{
				"id": 111,
				"protocol": "udp",
				"port": 53,
				"label": ""
			}`,
			expected: linodego.NetLoadBalancerListener{
				ID:       111,
				Protocol: "udp",
				Port:     53,
				Label:    "",
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
			var listener linodego.NetLoadBalancerListener
			err := listener.UnmarshalJSON([]byte(tt.jsonData))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.ID, listener.ID)
			assert.Equal(t, tt.expected.Protocol, listener.Protocol)
			assert.Equal(t, tt.expected.Port, listener.Port)
			assert.Equal(t, tt.expected.Label, listener.Label)

			// Test time fields if they were set in the JSON
			if tt.expected.ID == 999 {
				assert.NotNil(t, listener.Created)
				assert.NotNil(t, listener.Updated)
			}
		})
	}
}

func TestNetLoadBalancerListener_EdgeCases(t *testing.T) {
	t.Run("create listener with UDP protocol", func(t *testing.T) {
		fixtureData, err := fixtures.GetFixture("netloadbalancer_listener_create_udp")
		assert.NoError(t, err)

		var base ClientBaseCase
		base.SetUp(t)
		defer base.TearDown(t)

		base.MockPost("netloadbalancers/123/listeners", fixtureData)

		createOpts := linodego.NetLoadBalancerListenerCreateOptions{
			Protocol: "udp",
			Port:     53,
			Label:    "DNS Listener",
		}

		listener, err := base.Client.CreateNetLoadBalancerListener(context.Background(), 123, createOpts)
		assert.NoError(t, err)
		assert.Equal(t, "udp", listener.Protocol)
		assert.Equal(t, 53, listener.Port)
	})

	t.Run("create listener with high port number", func(t *testing.T) {
		fixtureData, err := fixtures.GetFixture("netloadbalancer_listener_create_high_port")
		assert.NoError(t, err)

		var base ClientBaseCase
		base.SetUp(t)
		defer base.TearDown(t)

		base.MockPost("netloadbalancers/123/listeners", fixtureData)

		createOpts := linodego.NetLoadBalancerListenerCreateOptions{
			Protocol: "tcp",
			Port:     65535,
			Label:    "High Port Listener",
		}

		listener, err := base.Client.CreateNetLoadBalancerListener(context.Background(), 123, createOpts)
		assert.NoError(t, err)
		assert.Equal(t, 65535, listener.Port)
	})

	t.Run("update multiple node weights", func(t *testing.T) {
		var base ClientBaseCase
		base.SetUp(t)
		defer base.TearDown(t)

		base.MockPost("netloadbalancers/123/listeners/456/node-weights", nil)

		updateOpts := linodego.NetLoadBalancerListenerNodeWeightsUpdateOptions{
			Nodes: []linodego.NetLoadBalancerListenerNodeWeightUpdateOptions{
				{ID: 1, Weight: 100},
				{ID: 2, Weight: 200},
				{ID: 3, Weight: 50},
				{ID: 4, Weight: 0}, // Zero weight to disable
			},
		}

		err := base.Client.UpdateNetLoadBalancerListenerNodeWeights(context.Background(), 123, 456, updateOpts)
		assert.NoError(t, err)
	})

	t.Run("update with complex listener options", func(t *testing.T) {
		fixtureData, err := fixtures.GetFixture("netloadbalancer_listener_update_complex")
		assert.NoError(t, err)

		var base ClientBaseCase
		base.SetUp(t)
		defer base.TearDown(t)

		base.MockPut("netloadbalancers/123/listeners/456", fixtureData)

		updateOpts := linodego.NetLoadBalancerListenerUpdateOptions{
			Protocol: "tcp",
			Port:     8080,
			Label:    "Complex Updated Listener",
			Nodes: []linodego.NetLoadBalancerNodeUpdateOptions{
				{
					Label:     "Updated Node 1",
					AddressV6: "2600:3c03::f03c:91ff:fe24:upd1",
					Weight:    150,
				},
			},
		}

		listener, err := base.Client.UpdateNetLoadBalancerListener(context.Background(), 123, 456, updateOpts)
		assert.NoError(t, err)
		assert.Equal(t, "Complex Updated Listener", listener.Label)
	})
}
