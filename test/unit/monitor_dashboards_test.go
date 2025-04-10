package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestListMonitorDashboards(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_dashboards")
	assert.NoError(t, err, "Expected no error when getting fixture")

	// fmt.Printf("[DEBUG] fixtureData = %+v\n", fixtureData)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/dashboards", fixtureData)

	clients, err := base.Client.ListMonitorDashboards(context.Background(), &linodego.ListOptions{})
	// fmt.Printf("[DEBUG] CLIENT = %+v\n", clients)
	assert.NoError(t, err, "Expected no error when listing monitor dashboards")
	assert.NotEmpty(t, clients, "Expected non-empty monitor dashboard list")

	assert.Equal(t, linodego.DashboardType("standard"), clients[0].Type, "Expected dashboard type to match")
	assert.Equal(t, linodego.ServiceType("dbaas"), clients[0].ServiceType, "Expected service_type to match")
}

func TestListMonitorDashboardsByID(t *testing.T) {
	// Load the mock fixture for monitor dashboard
	fixtureData, err := fixtures.GetFixture("monitor_dashboard_by_id")
	assert.NoError(t, err, "Expected no error when getting fixture")

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request for the monitor dashboard by ID (1)
	base.MockGet("monitor/dashboards/1", fixtureData)

	// Call the GetMonitorDashboard method
	clients, err := base.Client.GetMonitorDashboard(context.Background(), 1)
	assert.NoError(t, err, "Expected no error when listing monitor dashboard by type")
	assert.NotEmpty(t, clients, "Expected non-empty monitor dashboard list")

	assert.Equal(t, linodego.DashboardType("standard"), clients.Type, "Expected dashboard type to match")
	assert.Equal(t, linodego.ServiceType("dbaas"), clients.ServiceType, "Expected service_type to match")
}

// monitor_dashboard_by_service_type
func TestListMonitorDashboardsByServiceType(t *testing.T) {
	// Load the mock fixture for monitor dashboard
	fixtureData, err := fixtures.GetFixture("monitor_dashboard_by_service_type")
	assert.NoError(t, err, "Expected no error when getting fixture")

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request for the monitor dashboard by type (dbaas)
	base.MockGet("monitor/services/dbaas/dashboards", fixtureData)

	// Call the ListMonitorDashcdboardsByServiceType method
	clients, err := base.Client.ListMonitorDashboardsByServiceType(context.Background(), "dbaas", nil)
	assert.NoError(t, err, "Expected no error when listing monitor dashboard by type")
	assert.NotEmpty(t, clients, "Expected non-empty monitor dashboard list")

	assert.Equal(t, linodego.DashboardType("standard"), clients[0].Type, "Expected dashboard type to match")
	assert.Equal(t, linodego.ServiceType("dbaas"), clients[0].ServiceType, "Expected service_type to match")
}
