package linodego

// @TODO: disable alert channel for now because API made breaking changes.
// type AlertNotificationType string
//
// const (
//	EmailAlertNotification AlertNotificationType = "email"
//)
//
// type AlertChannelType string
//
// const (
//	SystemAlertChannel AlertChannelType = "system"
//	UserAlertChannel   AlertChannelType = "user"
//)
//
//// AlertChannelEnvelope represents a single alert channel entry returned inside alert definition
// type AlertChannelEnvelope struct {
//	ID    int    `json:"id"`
//	Label string `json:"label"`
//	Type  string `json:"type"`
//	URL   string `json:"url"`
//}
//
//// AlertChannel represents a Monitor Channel object.
// type AlertChannel struct {
//	Alerts      []AlertChannelEnvelope `json:"alerts"`
//	ChannelType AlertNotificationType  `json:"channel_type"`
//	Content     ChannelContent         `json:"content"`
//	Created     *time.Time             `json:"-"`
//	CreatedBy   string                 `json:"created_by"`
//	Updated     *time.Time             `json:"-"`
//	UpdatedBy   string                 `json:"updated_by"`
//	ID          int                    `json:"id"`
//	Label       string                 `json:"label"`
//	Type        AlertChannelType       `json:"type"`
// }
//
// type EmailChannelContent struct {
//	EmailAddresses []string `json:"email_addresses"`
// }
//
//// ChannelContent represents the content block for an AlertChannel, which varies by channel type.
// type ChannelContent struct {
//	Email *EmailChannelContent `json:"email"`
//	// Other channel types like 'webhook', 'slack' could be added here as optional fields.
// }
//
//// UnmarshalJSON implements the json.Unmarshaler interface
// func (a *AlertChannel) UnmarshalJSON(b []byte) error {
//	type Mask AlertChannel
//
//	p := struct {
//		*Mask
//
//		Updated *parseabletime.ParseableTime `json:"updated"`
//		Created *parseabletime.ParseableTime `json:"created"`
//	}{
//		Mask: (*Mask)(a),
//	}
//
//	if err := json.Unmarshal(b, &p); err != nil {
//		return err
//	}
//
//	a.Updated = (*time.Time)(p.Updated)
//	a.Created = (*time.Time)(p.Created)
//
//	return nil
// }

//// ListAlertChannels gets a paginated list of Alert Channels.
// func (c *Client) ListAlertChannels(ctx context.Context, opts *ListOptions) ([]AlertChannel, error) {
//	endpoint := formatAPIPath("monitor/alert-channels")
//	return getPaginatedResults[AlertChannel](ctx, c, endpoint, opts)
// }
