package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestListMonitorMetricDefinitionsByType(t *testing.T) {
	// Load the mock fixture for monitor_service_metrics
	fixtureData, err := fixtures.GetFixture("monitor_service_metrics")
	assert.NoError(t, err, "Expected no error when getting fixture")

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request for the monitor metric-definitions by type (dbaas)
	base.MockGet("monitor/services/dbaas/metric-definitions", fixtureData)

	// Call the ListMonitorMetricsDefinitionByServiceType method
	clients, err := base.Client.ListMonitorMetricsDefinitionByServiceType(context.Background(), "dbaas", &linodego.ListOptions{})
	assert.NoError(t, err, "Expected no error when listing monitor metric-definitions by type")
	assert.NotEmpty(t, clients, "Expected non-empty monitor metric-definitions list")

	assert.Equal(t, "cpu_usage", clients[0].Metric, "Expected Metric to match")
	assert.Equal(t, linodego.MetricType("gauge"), clients[0].MetricType, "Expected MetricType to match")
}
