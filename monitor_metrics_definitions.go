package linodego

import (
	"context"
)

// MonitorMetricsDefinition represents a MonitorMetricsDefinition object
type MonitorMetricsDefinition struct {
	AvailableAggregateFunctions []string            `json:"available_aggregate_functions"`
	Dimensions                  []MonitorDimensions `json:"dimensions"`
	IsAlertable                 bool                `json:"is_alertable"`
	Label                       string              `json:"label"`
	Metric                      string              `json:"metric"`
	MetricType                  string              `json:"metric_type"`
	ScrapeInterval              string              `json:"scrape_interval"`
	Unit                        string              `json:"unit"`
}

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
