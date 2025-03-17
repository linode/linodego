package linodego

import (
	"context"
)

// MonitorServicesClient represents a MonitorServicesClient object
type MonitorServices struct {
	Label       string `json:"label"`
	ServiceType string `json:"service_type"`
}

// ListLongviewClients lists LongviewClients
func (c *Client) ListMonitorServices(ctx context.Context, opts *ListOptions) ([]MonitorServices, error) {
	return getPaginatedResults[MonitorServices](ctx, c, "monitor/services", opts)
}
