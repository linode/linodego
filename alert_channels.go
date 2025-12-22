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
	ID          int            `json:"id"`
	Label       string         `json:"label"`
	ChannelType string         `json:"channel_type"`
	Content     ChannelContent `json:"content"`
	Created     string         `json:"created"`
	CreatedBy   string         `json:"created_by"`
	Updated     string         `json:"updated"`
	UpdatedBy   string         `json:"updated_by"`
}

// AlertChannelDetailOptions are the options used to create the details of a new Monitor Channel.
type AlertChannelDetailOptions struct {
	To string `json:"to,omitempty"`
}

type EmailChannelContent struct {
	EmailAddresses []string `json:"email_addresses"`
}

// ChannelContent represents the content block for an AlertChannel, which varies by channel type.
type ChannelContent struct {
	Email *EmailChannelContent `json:"email,omitempty"`
	// Other channel types like 'webhook', 'slack' could be added here as optional fields.
}

// ListAlertChannels gets a paginated list of Alert Channels.
func (c *Client) ListAlertChannels(ctx context.Context, opts *ListOptions) ([]AlertChannel, error) {
	endpoint := formatAPIPath("monitor/alert-channels")
	return getPaginatedResults[AlertChannel](ctx, c, endpoint, opts)
}
