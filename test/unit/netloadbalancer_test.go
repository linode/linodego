package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestNetLoadBalancer_Create(t *testing.T) {
	tests := []struct {
		name       string
		createOpts linodego.NetLoadBalancerCreateOptions
		fixture    string
	}{
		{
			name: "basic creation",
			createOpts: linodego.NetLoadBalancerCreateOptions{
				Label:  "Test NetLoadBalancer",
				Region: "us-east",
			},
			fixture: "netloadbalancer_create",
		},
		{
			name: "creation with listeners",
			createOpts: linodego.NetLoadBalancerCreateOptions{
				Label:  "NetLoadBalancer with Listeners",
				Region: "us-west",
				Listeners: []linodego.NetLoadBalancerListenerCreateOptions{
					{
						Protocol: "tcp",
						Port:     80,
						Label:    "HTTP Listener",
						Nodes: []linodego.NetLoadBalancerNodeCreateOptions{
							{
								Label:     "Web Server 1",
								AddressV6: "2600:3c03::f03c:91ff:fe24:abcd",
								Weight:    100,
							},
						},
					},
				},
			},
			fixture: "netloadbalancer_create_with_listeners",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixtureData, err := fixtures.GetFixture(tt.fixture)
			assert.NoError(t, err)

			var base ClientBaseCase
			base.SetUp(t)
			defer base.TearDown(t)

			base.MockPost("netloadbalancers", fixtureData)

			netloadbalancer, err := base.Client.CreateNetLoadBalancer(context.Background(), tt.createOpts)
			assert.NoError(t, err)

			assert.Equal(t, tt.createOpts.Label, netloadbalancer.Label)
			assert.Equal(t, tt.createOpts.Region, netloadbalancer.Region)
		})
	}
}

func TestNetLoadBalancer_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("netloadbalancer_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("netloadbalancers/123", fixtureData)

	netloadbalancer, err := base.Client.GetNetLoadBalancer(context.Background(), 123)
	assert.NoError(t, err)

	assert.Equal(t, 123, netloadbalancer.ID)
	assert.Equal(t, "Test NetLoadBalancer", netloadbalancer.Label)
	assert.Equal(t, "us-east", netloadbalancer.Region)
	assert.Equal(t, "active", netloadbalancer.Status)
	assert.Equal(t, "192.0.2.1", netloadbalancer.AddressV4)
	assert.Equal(t, "2600:3c03::f03c:91ff:fe24:1234", netloadbalancer.AddressV6)
}

func TestNetLoadBalancer_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("netloadbalancers_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("netloadbalancers", fixtureData)

	netloadbalancers, err := base.Client.ListNetLoadBalancers(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.Len(t, netloadbalancers, 2)

	// Verify details of the first NetLoadBalancer
	assert.Equal(t, 123, netloadbalancers[0].ID)
	assert.Equal(t, "NetLoadBalancer A", netloadbalancers[0].Label)
	assert.Equal(t, "us-east", netloadbalancers[0].Region)
	assert.Equal(t, "active", netloadbalancers[0].Status)

	// Verify details of the second NetLoadBalancer
	assert.Equal(t, 456, netloadbalancers[1].ID)
	assert.Equal(t, "NetLoadBalancer B", netloadbalancers[1].Label)
	assert.Equal(t, "us-west", netloadbalancers[1].Region)
	assert.Equal(t, "provisioning", netloadbalancers[1].Status)
}

func TestNetLoadBalancer_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("netloadbalancer_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("netloadbalancers/123", fixtureData)

	updateOpts := linodego.NetLoadBalancerUpdateOptions{
		Label: "Updated NetLoadBalancer",
	}
	netloadbalancer, err := base.Client.UpdateNetLoadBalancer(context.Background(), 123, updateOpts)
	assert.NoError(t, err)

	assert.Equal(t, 123, netloadbalancer.ID)
	assert.Equal(t, "Updated NetLoadBalancer", netloadbalancer.Label)
}

func TestNetLoadBalancer_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("netloadbalancers/123", nil)

	err := base.Client.DeleteNetLoadBalancer(context.Background(), 123)
	assert.NoError(t, err)
}

func TestNetLoadBalancer_GetCreateOptions(t *testing.T) {
	tests := []struct {
		name            string
		netloadbalancer linodego.NetLoadBalancer
		expected        linodego.NetLoadBalancerCreateOptions
	}{
		{
			name: "basic conversion",
			netloadbalancer: linodego.NetLoadBalancer{
				ID:     123,
				Label:  "Test NetLoadBalancer",
				Region: "us-east",
				Listeners: []linodego.NetLoadBalancerListener{
					{
						ID:       1,
						Protocol: "tcp",
						Port:     80,
						Label:    "HTTP Listener",
					},
				},
			},
			expected: linodego.NetLoadBalancerCreateOptions{
				Label:  "Test NetLoadBalancer",
				Region: "us-east",
				Listeners: []linodego.NetLoadBalancerListenerCreateOptions{
					{
						Protocol: "tcp",
						Port:     80,
						Label:    "HTTP Listener",
					},
				},
			},
		},
		{
			name: "conversion without listeners",
			netloadbalancer: linodego.NetLoadBalancer{
				ID:        456,
				Label:     "Simple NetLoadBalancer",
				Region:    "us-west",
				Listeners: []linodego.NetLoadBalancerListener{},
			},
			expected: linodego.NetLoadBalancerCreateOptions{
				Label:     "Simple NetLoadBalancer",
				Region:    "us-west",
				Listeners: []linodego.NetLoadBalancerListenerCreateOptions{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.netloadbalancer.GetCreateOptions()

			assert.Equal(t, tt.expected.Label, result.Label)
			assert.Equal(t, tt.expected.Region, result.Region)
			assert.Len(t, result.Listeners, len(tt.expected.Listeners))

			for i, expectedListener := range tt.expected.Listeners {
				assert.Equal(t, expectedListener.Protocol, result.Listeners[i].Protocol)
				assert.Equal(t, expectedListener.Port, result.Listeners[i].Port)
				assert.Equal(t, expectedListener.Label, result.Listeners[i].Label)
			}
		})
	}
}

func TestNetLoadBalancer_GetUpdateOptions(t *testing.T) {
	tests := []struct {
		name            string
		netloadbalancer linodego.NetLoadBalancer
		expected        linodego.NetLoadBalancerUpdateOptions
	}{
		{
			name: "basic conversion",
			netloadbalancer: linodego.NetLoadBalancer{
				ID:     123,
				Label:  "Updated NetLoadBalancer",
				Region: "us-east",
				Listeners: []linodego.NetLoadBalancerListener{
					{
						ID:       1,
						Protocol: "tcp",
						Port:     8080,
						Label:    "Updated Listener",
					},
				},
			},
			expected: linodego.NetLoadBalancerUpdateOptions{
				Label: "Updated NetLoadBalancer",
				Listeners: []linodego.NetLoadBalancerListenerUpdateOptions{
					{
						Protocol: "tcp",
						Port:     8080,
						Label:    "Updated Listener",
					},
				},
			},
		},
		{
			name: "conversion without listeners",
			netloadbalancer: linodego.NetLoadBalancer{
				ID:        456,
				Label:     "Simple Updated NetLoadBalancer",
				Region:    "us-west",
				Listeners: []linodego.NetLoadBalancerListener{},
			},
			expected: linodego.NetLoadBalancerUpdateOptions{
				Label:     "Simple Updated NetLoadBalancer",
				Listeners: []linodego.NetLoadBalancerListenerUpdateOptions{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.netloadbalancer.GetUpdateOptions()

			assert.Equal(t, tt.expected.Label, result.Label)
			assert.Len(t, result.Listeners, len(tt.expected.Listeners))

			for i, expectedListener := range tt.expected.Listeners {
				assert.Equal(t, expectedListener.Label, result.Listeners[i].Label)
			}
		})
	}
}

func TestNetLoadBalancer_ListWithOptions(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("netloadbalancers_list_filtered")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("netloadbalancers", fixtureData)

	listOpts := &linodego.ListOptions{
		PageOptions: &linodego.PageOptions{
			Page:    1,
			Pages:   1,
			Results: 1,
		},
		Filter: "{\"region\":\"us-east\"}",
	}

	netloadbalancers, err := base.Client.ListNetLoadBalancers(context.Background(), listOpts)
	assert.NoError(t, err)

	assert.Len(t, netloadbalancers, 1)
	assert.Equal(t, "us-east", netloadbalancers[0].Region)
}

func TestNetLoadBalancer_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected linodego.NetLoadBalancer
		wantErr  bool
	}{
		{
			name: "valid JSON with all fields",
			jsonData: `{
				"id": 789,
				"label": "Test NetLoadBalancer",
				"region": "us-southeast",
				"address_v4": "192.0.2.10",
				"address_v6": "2600:3c03::f03c:91ff:fe24:5678",
				"status": "active",
				"listeners": [],
				"created": "2025-03-01T10:00:00",
				"updated": "2025-03-15T14:30:00",
				"last_composite_updated": "2025-03-15T15:00:00"
			}`,
			expected: linodego.NetLoadBalancer{
				ID:        789,
				Label:     "Test NetLoadBalancer",
				Region:    "us-southeast",
				AddressV4: "192.0.2.10",
				AddressV6: "2600:3c03::f03c:91ff:fe24:5678",
				Status:    "active",
				Listeners: []linodego.NetLoadBalancerListener{},
			},
			wantErr: false,
		},
		{
			name: "valid JSON with minimal fields",
			jsonData: `{
				"id": 456,
				"label": "Minimal NetLoadBalancer",
				"region": "us-west",
				"address_v4": "",
				"address_v6": "",
				"status": "provisioning",
				"listeners": []
			}`,
			expected: linodego.NetLoadBalancer{
				ID:        456,
				Label:     "Minimal NetLoadBalancer",
				Region:    "us-west",
				AddressV4: "",
				AddressV6: "",
				Status:    "provisioning",
				Listeners: []linodego.NetLoadBalancerListener{},
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
			var nlb linodego.NetLoadBalancer
			err := nlb.UnmarshalJSON([]byte(tt.jsonData))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.ID, nlb.ID)
			assert.Equal(t, tt.expected.Label, nlb.Label)
			assert.Equal(t, tt.expected.Region, nlb.Region)
			assert.Equal(t, tt.expected.AddressV4, nlb.AddressV4)
			assert.Equal(t, tt.expected.AddressV6, nlb.AddressV6)
			assert.Equal(t, tt.expected.Status, nlb.Status)
			assert.Equal(t, tt.expected.Listeners, nlb.Listeners)

			// Test time fields if they were set in the JSON
			if tt.expected.ID == 789 {
				assert.NotNil(t, nlb.Created)
				assert.NotNil(t, nlb.Updated)
				assert.NotNil(t, nlb.LastCompositeUpdated)
			}
		})
	}
}

func TestNetLoadBalancer_CreateWithComplexListeners(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("netloadbalancer_create_complex")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("netloadbalancers", fixtureData)

	createOpts := linodego.NetLoadBalancerCreateOptions{
		Label:  "Complex NetLoadBalancer",
		Region: "us-central",
		Listeners: []linodego.NetLoadBalancerListenerCreateOptions{
			{
				Protocol: "tcp",
				Port:     80,
				Label:    "HTTP Listener",
				Nodes: []linodego.NetLoadBalancerNodeCreateOptions{
					{
						Label:     "Web Server 1",
						AddressV6: "2600:3c03::f03c:91ff:fe24:0001",
						Weight:    100,
					},
					{
						Label:     "Web Server 2",
						AddressV6: "2600:3c03::f03c:91ff:fe24:0002",
						Weight:    100,
					},
				},
			},
			{
				Protocol: "tcp",
				Port:     443,
				Label:    "HTTPS Listener",
				Nodes: []linodego.NetLoadBalancerNodeCreateOptions{
					{
						Label:     "SSL Server 1",
						AddressV6: "2600:3c03::f03c:91ff:fe24:0003",
						Weight:    50,
					},
				},
			},
		},
	}

	netloadbalancer, err := base.Client.CreateNetLoadBalancer(context.Background(), createOpts)
	assert.NoError(t, err)

	assert.Equal(t, "Complex NetLoadBalancer", netloadbalancer.Label)
	assert.Equal(t, "us-central", netloadbalancer.Region)
	assert.Len(t, netloadbalancer.Listeners, 2)
}
