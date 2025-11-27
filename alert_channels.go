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
	URL         string         `json:"url"`
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

// AlertingChannelCreateOptions are the options used to create a new Monitor Channel.
//
// Deprecated: AlertChannelCreateOptions should be used in all new implementations.
type AlertingChannelCreateOptions = AlertChannelCreateOptions

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

// GetAlertChannel gets an Alert Channel by ID.
func (c *Client) GetAlertChannel(ctx context.Context, channelID int) (*AlertChannel, error) {
	e := formatAPIPath("monitor/alert-channels/%d", channelID)
	return doGETRequest[AlertChannel](ctx, c, e)
}
