package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestListMonitorServices(t *testing.T) {
	// Load the mock fixture for monitor services
	fixtureData, err := fixtures.GetFixture("monitor_services")
	assert.NoError(t, err, "Expected no error when getting fixture")

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request for the monitor services endpoint
	base.MockGet("monitor/services", fixtureData)

	// Call the ListMonitorServices method
	clients, err := base.Client.ListMonitorServices(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err, "Expected no error when listing monitor services")
	assert.NotEmpty(t, clients, "Expected non-empty monitor services list")

	assert.Equal(t, "Databases", clients[0].Label, "Expected services label to match")
	assert.Equal(t, "dbaas", clients[0].ServiceType, "Expected service_type to match")
}

func TestListMonitorServicesByType(t *testing.T) {
	// Load the mock fixture for monitor services
	fixtureData, err := fixtures.GetFixture("monitor_services")
	assert.NoError(t, err, "Expected no error when getting fixture")

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request for the monitor services by type (dbaas)
	base.MockGet("monitor/services/dbaas", fixtureData)

	// Call the ListMonitorServiceByType method
	clients, err := base.Client.ListMonitorServiceByType(context.Background(), "dbaas", &linodego.ListOptions{})
	assert.NoError(t, err, "Expected no error when listing monitor services by type")
	assert.NotEmpty(t, clients, "Expected non-empty monitor services list")

	assert.Equal(t, "Databases", clients[0].Label, "Expected services label to match")
	assert.Equal(t, "dbaas", clients[0].ServiceType, "Expected service_type to match")
}
