package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// MonitorDashboard represents an ACLP Dashboard object
type MonitorDashboard struct {
	ID          int               `json:"id"`
	Type        DashboardType     `json:"type"`
	ServiceType ServiceType       `json:"service_type"`
	Label       string            `json:"label"`
	Created     *time.Time        `json:"-"`
	Updated     *time.Time        `json:"-"`
	Widgets     []DashboardWidget `json:"widgets"`
}

// enum object for serviceType
type ServiceType string

const (
	Linode          ServiceType = "linode"
	LKE             ServiceType = "lke"
	DBaaS           ServiceType = "dbaas"
	ACLB            ServiceType = "aclb"
	Nodebalancer    ServiceType = "nodebalancer"
	ObjectStorage   ServiceType = "objectstorage"
	Vpc             ServiceType = "vpc"
	FirewallService ServiceType = "firewall"
)

// enum object for DashboardType
type DashboardType string

const (
	Standard DashboardType = "standard"
	Custom   DashboardType = "custom"
)

// DashboardWidget represents an ACLP DashboardWidget object
type DashboardWidget struct {
	Metric            string            `json:"metric"`
	Unit              string            `json:"unit"`
	Label             string            `json:"label"`
	Color             string            `json:"color"`
	Size              int               `json:"size"`
	ChartType         ChartType         `json:"chart_type"`
	YLabel            string            `json:"y_label"`
	AggregateFunction AggregateFunction `json:"aggregate_function"`
}

// Enum object for AggregateFunction
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

// Enum object for Chart type
type ChartType string

const (
	Line ChartType = "line"
	Area ChartType = "area"
)

// ListMonitorDashboards lists all the ACLP Monitor Dashboards
func (c *Client) ListMonitorDashboards(ctx context.Context, opts *ListOptions) ([]MonitorDashboard, error) {
	return getPaginatedResults[MonitorDashboard](ctx, c, "monitor/dashboards", opts)
}

// GetMonitorDashboard gets an ACLP Monitor Dashboard for a given dashboardID
func (c *Client) GetMonitorDashboard(ctx context.Context, dashboardID int) (*MonitorDashboard, error) {
	e := formatAPIPath("monitor/dashboards/%d", dashboardID)
	return doGETRequest[MonitorDashboard](ctx, c, e)
}

// ListMonitorDashboardsByServiceType lists ACLP Monitor Dashboards for a given serviceType
func (c *Client) ListMonitorDashboardsByServiceType(ctx context.Context, serviceType string, opts *ListOptions) ([]MonitorDashboard, error) {
	e := formatAPIPath("monitor/services/%s/dashboards", serviceType)
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
