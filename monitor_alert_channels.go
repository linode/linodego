package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/v2/internal/parseabletime"
)

type AlertNotificationType string

const (
	EmailAlertNotification AlertNotificationType = "email"
)

type AlertChannelType string

const (
	SystemAlertChannel AlertChannelType = "system"
	UserAlertChannel   AlertChannelType = "user"
)

// AlertChannel represents a Monitor Channel object.
type AlertChannel struct {
	Alerts      AlertsInfo            `json:"alerts"`
	ChannelType AlertNotificationType `json:"channel_type"`
	Details     ChannelDetails        `json:"details"`
	Created     *time.Time            `json:"-"`
	CreatedBy   string                `json:"created_by"`
	Updated     *time.Time            `json:"-"`
	UpdatedBy   string                `json:"updated_by"`
	ID          int                   `json:"id"`
	Label       string                `json:"label"`
	Type        AlertChannelType      `json:"type"`
}

// AlertsInfo represents alert information for a channel
type AlertsInfo struct {
	URL        string `json:"url"`
	Type       string `json:"type"`
	AlertCount int    `json:"alert_count"`
}

// ChannelDetails represents the details block for an AlertChannel
type ChannelDetails struct {
	Email *EmailChannelDetails `json:"email"`
	// Other channel types could be added here
}

// EmailChannelDetails represents email-specific details for a channel
type EmailChannelDetails struct {
	Usernames     []string `json:"usernames"`
	RecipientType string   `json:"recipient_type"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (a *AlertChannel) UnmarshalJSON(b []byte) error {
	type Mask AlertChannel

	p := struct {
		*Mask

		Updated *parseabletime.ParseableTime `json:"updated"`
		Created *parseabletime.ParseableTime `json:"created"`
	}{
		Mask: (*Mask)(a),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	a.Updated = (*time.Time)(p.Updated)
	a.Created = (*time.Time)(p.Created)

	return nil
}

// AlertChannelCreateOptions represents options for creating an alert notification channel.
type AlertChannelCreateOptions struct {
	ChannelType AlertNotificationType      `json:"channel_type"`
	Details     AlertChannelDetailsOptions `json:"details"`
	Label       *string                    `json:"label,omitzero"`
}

// AlertChannelDetailsOptions represents the details configuration for an alert channel.
type AlertChannelDetailsOptions struct {
	Email *EmailChannelCreateOptions `json:"email,omitzero"`
}

// EmailChannelCreateOptions represents email-specific configuration for an alert channel.
type EmailChannelCreateOptions struct {
	Usernames     []string `json:"usernames"`
	RecipientType *string  `json:"recipient_type,omitzero"`
}

// AlertChannelUpdateOptions represents options for updating an alert notification channel.
type AlertChannelUpdateOptions struct {
	Details *AlertChannelUpdateDetailsOptions `json:"details,omitzero"`
	Label   *string                           `json:"label,omitzero"`
}

// AlertChannelUpdateDetailsOptions represents update details for an alert channel.
type AlertChannelUpdateDetailsOptions struct {
	Email *EmailChannelUpdateOptions `json:"email,omitzero"`
}

// EmailChannelUpdateOptions represents email-specific update configuration for an alert channel.
type EmailChannelUpdateOptions struct {
	Usernames []string `json:"usernames,omitzero"`
}

// Alert represents an alert definition assigned to a notification channel.
type Alert struct {
	ID          int    `json:"id"`
	Label       string `json:"label"`
	ServiceType string `json:"service_type"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

// ListAlertChannels gets a paginated list of Alert Channels.
func (c *Client) ListAlertChannels(ctx context.Context, opts *ListOptions) ([]AlertChannel, error) {
	endpoint := formatAPIPath("monitor/alert-channels")
	return getPaginatedResults[AlertChannel](ctx, c, endpoint, opts)
}

// GetAlertChannel retrieves a single Alert Channel by ID.
func (c *Client) GetAlertChannel(ctx context.Context, channelID int) (*AlertChannel, error) {
	endpoint := formatAPIPath("monitor/alert-channels/%d", channelID)
	return doGETRequest[AlertChannel](ctx, c, endpoint)
}

// CreateAlertChannel creates a new alert notification channel.
func (c *Client) CreateAlertChannel(ctx context.Context, opts AlertChannelCreateOptions) (*AlertChannel, error) {
	endpoint := formatAPIPath("monitor/alert-channels")
	return doPOSTRequest[AlertChannel](ctx, c, endpoint, opts)
}

// UpdateAlertChannel updates an alert notification channel.
func (c *Client) UpdateAlertChannel(ctx context.Context, channelID int, opts AlertChannelUpdateOptions) (*AlertChannel, error) {
	endpoint := formatAPIPath("monitor/alert-channels/%d", channelID)
	return doPUTRequest[AlertChannel](ctx, c, endpoint, opts)
}

// DeleteAlertChannel deletes an alert notification channel.
func (c *Client) DeleteAlertChannel(ctx context.Context, channelID int) error {
	endpoint := formatAPIPath("monitor/alert-channels/%d", channelID)
	return doDELETERequest(ctx, c, endpoint)
}

// ListAlertsForChannel gets a paginated list of Alert Definitions for a specific alert notification channel.
func (c *Client) ListAlertsForChannel(ctx context.Context, channelID int, opts *ListOptions) ([]Alert, error) {
	endpoint := formatAPIPath("monitor/alert-channels/%d/alerts", channelID)
	return getPaginatedResults[Alert](ctx, c, endpoint, opts)
}
