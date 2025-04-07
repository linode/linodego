package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestMonitorMetricDefinitions_Get_smoke(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestMonitorMetricDefinitions_Get")
	defer teardown()

	// Get the metric-definitions by serviceType
	monitorMetricDefinitionsClientList, listErr := client.ListMonitorMetricsDefinitionByServiceType(context.Background(), "dbaas", nil)
	if listErr != nil {
		t.Errorf("Error listing monitor metrics:%s", listErr)
	}

	for _, metrics_def := range monitorMetricDefinitionsClientList {
		validateMetricDefinitions(t, metrics_def)
	}

	// Get the metric-definitions by serviceType for the filter "is_alertable":false
	monitorMetricDefinitionsClientListFilter, listErr := client.ListMonitorMetricsDefinitionByServiceType(context.Background(), "dbaas", linodego.NewListOptions(0, "{\"is_alertable\":false}"))
	if listErr != nil {
		t.Errorf("Error listing monitor metrics:%s", listErr)
	}

	for _, metrics_def := range monitorMetricDefinitionsClientListFilter {
		validateMetricDefinitionsFilters(t, metrics_def)
	}
}

func validateMetricDefinitions(
	t *testing.T,
	metrics_def linodego.MonitorMetricsDefinition,
) {
	require.NotEmpty(t, metrics_def.AvailableAggregateFunctions)
	require.NotEmpty(t, metrics_def.Dimensions)
	require.NotEmpty(t, metrics_def.Label)
	require.NotEmpty(t, metrics_def.Metric)
	require.NotEmpty(t, metrics_def.MetricType)
	require.NotEmpty(t, metrics_def.ScrapeInterval)
	require.NotEmpty(t, metrics_def.Unit)

	require.True(t, metrics_def.IsAlertable || !metrics_def.IsAlertable, "IsAlertable should be true or false")
}

// Validation function for filter "is_alertable":false
func validateMetricDefinitionsFilters(
	t *testing.T,
	metrics_def linodego.MonitorMetricsDefinition,
) {
	require.NotEmpty(t, metrics_def.AvailableAggregateFunctions)
	require.NotEmpty(t, metrics_def.Dimensions)
	require.NotEmpty(t, metrics_def.Label)
	require.NotEmpty(t, metrics_def.Metric)
	require.NotEmpty(t, metrics_def.MetricType)
	require.NotEmpty(t, metrics_def.ScrapeInterval)
	require.NotEmpty(t, metrics_def.Unit)
	require.False(t, metrics_def.IsAlertable, "IsAlertable should be false for the given filter")
}
