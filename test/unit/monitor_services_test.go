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

	// Assert alert details for the first service
	assert.NotNil(t, clients[0].Alert, "Expected alert to be present")
	assert.NotEmpty(t, clients[0].Alert.PollingIntervalSeconds, "Expected polling_interval_seconds to be present")
	assert.NotEmpty(t, clients[0].Alert.EvaluationPeriodSeconds, "Expected evaluation_period_seconds to be present")
	assert.NotEmpty(t, clients[0].Alert.Scope, "Expected scope to be present")
}

func TestListMonitorServicesByType(t *testing.T) {
	// Load the mock fixture for monitor services
	fixtureData, err := fixtures.GetFixture("monitor_services_dbaas")
	assert.NoError(t, err, "Expected no error when getting fixture")

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request for the monitor services by type (dbaas)
	base.MockGet("monitor/services/dbaas", fixtureData)

	// Call the ListMonitorServiceByType method
	client, err := base.Client.ListMonitorServiceByType(context.Background(), "dbaas", &linodego.ListOptions{})
	assert.NoError(t, err, "Expected no error when listing monitor services by type")
	// No need to check NotEmpty, just check fields directly
	assert.Equal(t, "Databases", client.Label, "Expected services label to match")
	assert.Equal(t, "dbaas", client.ServiceType, "Expected service_type to match")

	// Assert alert details for the first service
	assert.NotNil(t, client.Alert, "Expected alert to be present")
	assert.NotEmpty(t, client.Alert.PollingIntervalSeconds, "Expected polling_interval_seconds to be present")
	assert.NotEmpty(t, client.Alert.EvaluationPeriodSeconds, "Expected evaluation_period_seconds to be present")
	assert.NotEmpty(t, client.Alert.Scope, "Expected scope to be present")
}
