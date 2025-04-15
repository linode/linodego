package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// AccountMaintenance represents a Maintenance object for any entity a user has permissions to view
type AccountMaintenance struct {
	Body                 string     `json:"body"`
	Entity               *Entity    `json:"entity"`
	Label                string     `json:"label"`
	Message              string     `json:"message"`
	Severity             string     `json:"severity"`
	Type                 string     `json:"type"`
	EventType            string     `json:"event_type"`
	MaintenancePolicySet string     `json:"maintenance_policy_set"`
	Description          string     `json:"description"`
	Source               string     `json:"source"`
	NotBefore            *time.Time `json:"-"`
	StartTime            *time.Time `json:"-"`
	CompleteTime         *time.Time `json:"-"`
	Status               string     `json:"status"`
	When                 *time.Time `json:"-"`
	Until                *time.Time `json:"-"`
}

// The entity being affected by maintenance
type Entity struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (accountMaintenance *AccountMaintenance) UnmarshalJSON(b []byte) error {
	type Mask AccountMaintenance

	p := struct {
		*Mask
		NotBefore    *parseabletime.ParseableTime `json:"not_before"`
		StartTime    *parseabletime.ParseableTime `json:"start_time"`
		CompleteTime *parseabletime.ParseableTime `json:"complete_time"`
		When         *parseabletime.ParseableTime `json:"when"`
		Until        *parseabletime.ParseableTime `json:"until"`
	}{
		Mask: (*Mask)(accountMaintenance),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	accountMaintenance.NotBefore = (*time.Time)(p.NotBefore)
	accountMaintenance.StartTime = (*time.Time)(p.StartTime)
	accountMaintenance.CompleteTime = (*time.Time)(p.CompleteTime)
	accountMaintenance.When = (*time.Time)(p.When)
	accountMaintenance.Until = (*time.Time)(p.Until)

	return nil
}

// ListMaintenances lists Account Maintenance objects for any entity a user has permissions to view
func (c *Client) ListMaintenances(ctx context.Context, opts *ListOptions) ([]AccountMaintenance, error) {
	return getPaginatedResults[AccountMaintenance](ctx, c, "account/maintenance", opts)
}
