package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// MonitorDashboardsClient represents a MonitorDashboardsClient object
type MonitorDashboards struct {
	ID          int                `json:"id"`
	Type        string             `json:"type"`
	ServiceType string             `json:"service_type"`
	Label       string             `json:"label"`
	Created     *time.Time         `json:"-"`
	Updated     *time.Time         `json:"-"`
	Widgets     []DashboardWidgets `json:"widgets"`
}

type DashboardWidgets struct {
	Metric            string `json:"metric"`
	Unit              string `json:"unit"`
	Label             string `json:"label"`
	Color             string `json:"color"`
	Size              int    `json:"size"`
	ChartType         string `json:"chart_type"`
	YLabel            string `json:"y_label"`
	AggregateFunction string `json:"aggregate_function"`
}

// ListMonitorDashboards lists MonitorDashboards
func (c *Client) ListMonitorDashboards(ctx context.Context, opts *ListOptions) ([]MonitorDashboards, error) {
	return getPaginatedResults[MonitorDashboards](ctx, c, "monitor/dashboards", opts)
}

func (c *Client) GetMonitorDashboardsByID(ctx context.Context, dashboard_id int) (*MonitorDashboards, error) {
	e := formatAPIPath("monitor/dashboards/%d", dashboard_id)
	return doGETRequest[MonitorDashboards](ctx, c, e)
}

func (c *Client) GetMonitorDashboardsByServiceType(ctx context.Context, service_type string, opts *ListOptions) ([]MonitorDashboards, error) {
	e := formatAPIPath("monitor/services/%s/dashboards", service_type)
	return getPaginatedResults[MonitorDashboards](ctx, c, e, opts)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *MonitorDashboards) UnmarshalJSON(b []byte) error {
	type Mask MonitorDashboards

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Created = (*time.Time)(p.Created)
	i.Updated = (*time.Time)(p.Updated)

	return nil
}
