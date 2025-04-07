package linodego

import (
	"context"
)

// MonitorMetricsDefinition represents a MonitorMetricsDefinition object
type MonitorMetricsDefinition struct {
	AvailableAggregateFunctions []AggregateFunction `json:"available_aggregate_functions"`
	Dimensions                  []MonitorDimensions `json:"dimensions"`
	IsAlertable                 bool                `json:"is_alertable"`
	Label                       string              `json:"label"`
	Metric                      string              `json:"metric"`
	MetricType                  MetricType          `json:"metric_type"`
	ScrapeInterval              string              `json:"scrape_interval"`
	Unit                        Unit                `json:"unit"`
}

type MetricType string

const (
	Counter   MetricType = "counter"
	Histogram MetricType = "histogram"
	Gauge     MetricType = "gauge"
	Summary   MetricType = "summary"
)

type Unit string

const (
	CountUnit      Unit = "count"
	Percent        Unit = "percent"
	Byte           Unit = "byte"
	Second         Unit = "second"
	BitsPerSecond  Unit = "bits_per_second"
	Millisecond    Unit = "millisecond"
	KB             Unit = "KB"
	MB             Unit = "MB"
	GB             Unit = "GB"
	RateUnit       Unit = "rate"
	BytesPerSecond Unit = "bytes_per_second"
	Percentile     Unit = "percentile"
	Ratio          Unit = "ratio"
	OpsPerSecond   Unit = "ops_per_second"
	Iops           Unit = "iops"
)

type MonitorDimensions struct {
	DimensionLabel string   `json:"dimension_label"`
	Label          string   `json:"label"`
	Values         []string `json:"values"`
}

// ListMonitorMetricsDefinitionByServiceType lists metric definitions
func (c *Client) ListMonitorMetricsDefinitionByServiceType(ctx context.Context, service_type string, opts *ListOptions) ([]MonitorMetricsDefinition, error) {
	e := formatAPIPath("monitor/services/%s/metric-definitions", service_type)
	return getPaginatedResults[MonitorMetricsDefinition](ctx, c, e, opts)
}
