package linodego

import (
	"context"
)

// AlertChannelEnvelope represents a single alert channel entry returned inside alert definition
type AlertChannelEnvelope struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

// AlertChannel represents a Monitor Channel object.
type AlertChannel struct {
	ID           int                    `json:"id"`
	Alerts       []AlertChannelEnvelope `json:"alerts"`
	Label        string                 `json:"label"`
	Channel_type string                 `json:"channel_type"`
	Content      ChannelContent         `json:"content"`
	Type         AlertType              `json:"type"`
	Details      AlertChannelDetail     `json:"details"`
	Created      string                 `json:"created"`
	Created_by   string                 `json:"created_by"`
	Updated      string                 `json:"updated"`
	Updated_by   string                 `json:"updated_by"`
}

// AlertChannelDetail represents the details of a Monitor Channel.
type AlertChannelDetail struct {
	To    string `json:"to,omitempty"`
	From  string `json:"from,omitempty"`
	User  string `json:"user,omitempty"`
	Token string `json:"token,omitempty"`
	URL   string `json:"url,omitempty"`
}

// AlertChannelCreateOptions are the options used to create a new Monitor Channel.
type AlertChannelCreateOptions struct {
	Label   string                    `json:"label"`
	Type    string                    `json:"type"`
	Details AlertChannelDetailOptions `json:"details"`
}

// AlertChannelDetailOptions are the options used to create the details of a new Monitor Channel.
type AlertChannelDetailOptions struct {
	To string `json:"to,omitempty"`
}

// Backwards-compat alias for older name
type AlertingChannelCreateOptions = AlertChannelCreateOptions

type EmailChannelContent struct {
	EmailAddresses []string `json:"email_addresses"`
}

// ChannelContent represents the content block for an AlertChannel, which varies by channel type.
type ChannelContent struct {
	Email *EmailChannelContent `json:"email,omitempty"`
	// Other channel types like 'webhook', 'slack' could be added here as optional fields.
}

// ListMonitorChannels gets a paginated list of Monitor Channels.
func (c *Client) ListAlertChannels(ctx context.Context, opts *ListOptions) ([]AlertChannel, error) {
	endpoint := formatAPIV4BetaPath("monitor/channels")
	return getPaginatedResults[AlertChannel](ctx, c, endpoint, opts)
}

// GetMonitorChannel gets a Monitor Channel by ID.
func (c *Client) GetAlertChannel(ctx context.Context, channelID int) (*AlertChannel, error) {
	e := formatAPIV4BetaPath("monitor/channels/%d", channelID)
	return doGETRequest[AlertChannel](ctx, c, e)
}

func (c *Client) GetAlertChannels(ctx context.Context) (*AlertChannel, error) {
	e := formatAPIV4BetaPath("monitor/alert-channels/")
	return doGETRequest[AlertChannel](ctx, c, e)
}
