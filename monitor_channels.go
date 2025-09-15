package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// MonitorChannel represents a Monitor Channel object.
type MonitorChannel struct {
	ID      int                  `json:"id"`
	Label   string               `json:"label"`
	Type    string               `json:"type"`
	Details MonitorChannelDetail `json:"details"`
	Created *time.Time           `json:"-"`
	Updated *time.Time           `json:"-"`
}

// MonitorChannelDetail represents the details of a Monitor Channel.
type MonitorChannelDetail struct {
	To    string `json:"to,omitempty"`
	From  string `json:"from,omitempty"`
	User  string `json:"user,omitempty"`
	Token string `json:"token,omitempty"`
	URL   string `json:"url,omitempty"`
}

// MonitorChannelCreateOptions are the options used to create a new Monitor Channel.
type MonitorChannelCreateOptions struct {
	Label   string                      `json:"label"`
	Type    string                      `json:"type"`
	Details MonitorChannelDetailOptions `json:"details"`
}

// MonitorChannelDetailOptions are the options used to create the details of a new Monitor Channel.
type MonitorChannelDetailOptions struct {
	To string `json:"to,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (i *MonitorChannel) UnmarshalJSON(b []byte) error {
	type Mask MonitorChannel

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

// CreateMonitorChannel creates a new Monitor Channel.
func (c *Client) CreateMonitorChannel(ctx context.Context, opts MonitorChannelCreateOptions) (*MonitorChannel, error) {
	e := "v4beta/monitor/channels"
	return doPOSTRequest[MonitorChannel, MonitorChannelCreateOptions](ctx, c, e, opts)
}

// DeleteMonitorChannel deletes a Monitor Channel.
func (c *Client) DeleteMonitorChannel(ctx context.Context, channelID int) error {
	e := formatAPIV4BetaPath("monitor/channels/%d", channelID)
	return doDELETERequest(ctx, c, e)
}

// ListMonitorChannels gets a paginated list of Monitor Channels.
func (c *Client) ListMonitorChannels(ctx context.Context, opts *ListOptions) ([]MonitorChannel, error) {
	endpoint := formatAPIV4BetaPath("monitor/channels")
	return getPaginatedResults[MonitorChannel](ctx, c, endpoint, opts)
}

// GetMonitorChannel gets a Monitor Channel by ID.
func (c *Client) GetMonitorChannel(ctx context.Context, channelID int) (*MonitorChannel, error) {
	e := formatAPIV4BetaPath("monitor/channels/%d", channelID)
	return doGETRequest[MonitorChannel](ctx, c, e)
}

func (c *Client) GetAlertChannels(ctx context.Context) (*MonitorChannel, error) {
	e := formatAPIV4BetaPath("monitor/alert-channels/")
	return doGETRequest[MonitorChannel](ctx, c, e)
}
