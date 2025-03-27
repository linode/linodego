package linodego

import (
	"context"
)

// MonitorServicesClient represents a MonitorServicesClient object
type MonitorServices struct {
	Label       string `json:"label"`
	ServiceType string `json:"service_type"`
}

// ListMonitorDashboards lists MonitorDashboards
func (c *Client) ListMonitorServices(ctx context.Context, opts *ListOptions) ([]MonitorServices, error) {
	return getPaginatedResults[MonitorServices](ctx, c, "monitor/services", opts)
}

func (c *Client) GetMonitorServiceByType(ctx context.Context, service_type string, opts *ListOptions) ([]MonitorServices, error) {
	e := formatAPIPath("monitor/services/%s", service_type)
	return getPaginatedResults[MonitorServices](ctx, c, e, opts)
}
