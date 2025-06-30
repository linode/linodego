package linodego

import (
	"context"
)

type MaintenancePolicy struct {
	Slug                  string `json:"slug"`
	Label                 string `json:"label"`
	Description           string `json:"description"`
	Type                  string `json:"type"`
	NotificationPeriodSec int    `json:"notification_period_sec"`
	IsDefault             bool   `json:"is_default"`
}

func (c *Client) ListMaintenancePolicies(ctx context.Context) ([]MaintenancePolicy, error) {
	response, err := doGETRequest[[]MaintenancePolicy](ctx, c, "maintenance/policies")
	if err != nil {
		return nil, err
	}

	return *response, nil
}
