package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// MonitorDashboard represents a MonitorDashboard object
type MonitorDashboard struct {
	ID          int                `json:"id"`
	Type        DashboardType      `json:"type"`
	ServiceType ServiceType        `json:"service_type"`
	Label       string             `json:"label"`
	Created     *time.Time         `json:"-"`
	Updated     *time.Time         `json:"-"`
	Widgets     []DashboardWidgets `json:"widgets"`
}

type ServiceType string

const (
	Linode          ServiceType = "linode"
	LKE             ServiceType = "lke"
	Dbaas           ServiceType = "dbaas"
	ACLB            ServiceType = "aclb"
	nodebalancer    ServiceType = "nodebalancer"
	objectstorage   ServiceType = "objectstorage"
	Vpc             ServiceType = "vpc"
	FirewallService ServiceType = "firewall"
)

type DashboardType string

const (
	Standard DashboardType = "standard"
	Custom   DashboardType = "custom"
)

type DashboardWidgets struct {
	Metric            string            `json:"metric"`
	Unit              string            `json:"unit"`
	Label             string            `json:"label"`
	Color             string            `json:"color"`
	Size              int               `json:"size"`
	ChartType         ChartType         `json:"chart_type"`
	YLabel            string            `json:"y_label"`
	AggregateFunction AggregateFunction `json:"aggregate_function"`
}

type AggregateFunction string

const (
	Min      AggregateFunction = "min"
	Max      AggregateFunction = "max"
	Avg      AggregateFunction = "avg"
	Sum      AggregateFunction = "sum"
	Rate     AggregateFunction = "rate"
	Increase AggregateFunction = "increase"
	Count    AggregateFunction = "count"
	Last     AggregateFunction = "last"
)

type ChartType string

const (
	Line ChartType = "line"
	Area ChartType = "area"
)

// ListMonitorDashboards lists MonitorDashboard
func (c *Client) ListMonitorDashboards(ctx context.Context, opts *ListOptions) ([]MonitorDashboard, error) {
	return getPaginatedResults[MonitorDashboard](ctx, c, "monitor/dashboards", opts)
}

func (c *Client) GetMonitorDashboardsByID(ctx context.Context, dashboard_id int) (*MonitorDashboard, error) {
	e := formatAPIPath("monitor/dashboards/%d", dashboard_id)
	return doGETRequest[MonitorDashboard](ctx, c, e)
}

func (c *Client) GetMonitorDashboardsByServiceType(ctx context.Context, service_type string, opts *ListOptions) ([]MonitorDashboard, error) {
	e := formatAPIPath("monitor/services/%s/dashboards", service_type)
	return getPaginatedResults[MonitorDashboard](ctx, c, e, opts)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *MonitorDashboard) UnmarshalJSON(b []byte) error {
	type Mask MonitorDashboard

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
