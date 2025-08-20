package unit

import (
	"context"
	"fmt"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestFetchEntityMetrics(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_api_entity_metrics")
	assert.NoError(t, err)

	var base MonitorClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	serviceType := "dbaas"
	base.MockPost(fmt.Sprintf("monitor/services/%s/metrics", serviceType), fixtureData)

	opts := linodego.EntityMetricsFetchOptions{
		EntityIDs: []any{13217, 13316},
		Metrics: []linodego.EntityMetric{
			{
				Name:              "read_iops",
				AggregateFunction: linodego.AggregateFunctionAvg,
			},
			{
				Name:              "cpu_usage",
				AggregateFunction: linodego.AggregateFunctionAvg,
			},
		},
		RelativeTimeDuration: &linodego.MetricRelativeTimeDuration{
			Unit:  linodego.MetricTimeUnitHr,
			Value: 1,
		},
	}

	metrics, err := base.MonitorClient.FetchEntityMetrics(context.Background(), serviceType, &opts)
	assert.NoError(t, err, "Expected no error when getting the entity metrics")

	// get the first metric and assert its values
	metric := metrics.Data.Result[0]

	assert.Equal(t, metric.Metric["entity_id"], float64(13316))
	assert.Equal(t, metric.Metric["metric_name"], "avg_read_iops")
	assert.Equal(t, metric.Metric["node_id"], "primary-9")
	assert.Equal(t, metric.Values[0][0], float64(1728996500))
	assert.Equal(t, metric.Values[0][1], "90.55555555555556")
	assert.False(t, metrics.IsPartial)
	assert.Equal(t, metrics.Status, "success")
	assert.Equal(t, metrics.Stats.ExecutionTimeMsec, 21)
	assert.Equal(t, metrics.Stats.SeriesFetched, "2")
}
