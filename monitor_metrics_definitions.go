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
	Unit                        Unit                `json:"unit"`
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
type Unit string

const (
	UnitCount          Unit = "count"
	UnitPercent        Unit = "percent"
	UnitByte           Unit = "byte"
	UnitSecond         Unit = "second"
	UnitBitsPerSecond  Unit = "bits_per_second"
	UnitMillisecond    Unit = "millisecond"
	UnitKB             Unit = "KB"
	UnitMB             Unit = "MB"
	UnitGB             Unit = "GB"
	UnitRate           Unit = "rate"
	UnitBytesPerSecond Unit = "bytes_per_second"
	UnitPercentile     Unit = "percentile"
	UnitRatio          Unit = "ratio"
	UnitOpsPerSecond   Unit = "ops_per_second"
	UnitIops           Unit = "iops"
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
