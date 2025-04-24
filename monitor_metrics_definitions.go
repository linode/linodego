package linodego

import (
	"context"
)

// MonitorMetricsDefinition represents an ACLP MetricsDefinition object
type MonitorMetricsDefinition struct {
	AvailableAggregateFunctions []AggregateFunction `json:"available_aggregate_functions"`
	Dimensions                  []MonitorDimension  `json:"dimensions"`
	IsAlertable                 bool                `json:"is_alertable"`
	Label                       string              `json:"label"`
	Metric                      string              `json:"metric"`
	MetricType                  MetricType          `json:"metric_type"`
	ScrapeInterval              string              `json:"scrape_interval"`
	Unit                        MetricUnit          `json:"unit"`
}

// Enum object for MetricType
type MetricType string

const (
	MetricTypeCounter   MetricType = "counter"
	MetricTypeHistogram MetricType = "histogram"
	MetricTypeGauge     MetricType = "gauge"
	MetricTypeSummary   MetricType = "summary"
)

// Enum object for Unit
type MetricUnit string

const (
	UnitCount          MetricUnit = "count"
	UnitPercent        MetricUnit = "percent"
	UnitByte           MetricUnit = "byte"
	UnitSecond         MetricUnit = "second"
	UnitBitsPerSecond  MetricUnit = "bits_per_second"
	UnitMillisecond    MetricUnit = "millisecond"
	UnitKB             MetricUnit = "KB"
	UnitMB             MetricUnit = "MB"
	UnitGB             MetricUnit = "GB"
	UnitRate           MetricUnit = "rate"
	UnitBytesPerSecond MetricUnit = "bytes_per_second"
	UnitPercentile     MetricUnit = "percentile"
	UnitRatio          MetricUnit = "ratio"
	UnitOpsPerSecond   MetricUnit = "ops_per_second"
	UnitIops           MetricUnit = "iops"
)

// MonitorDimension represents an ACLP MonitorDimension object
type MonitorDimension struct {
	DimensionLabel string   `json:"dimension_label"`
	Label          string   `json:"label"`
	Values         []string `json:"values"`
}

// ListMonitorMetricsDefinitionByServiceType lists metric definitions
func (c *Client) ListMonitorMetricsDefinitionByServiceType(ctx context.Context, serviceType string, opts *ListOptions) ([]MonitorMetricsDefinition, error) {
	e := formatAPIPath("monitor/services/%s/metric-definitions", serviceType)
	return getPaginatedResults[MonitorMetricsDefinition](ctx, c, e, opts)
}
