package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// LogsDestinationType represents the type of a logs destination.
type LogsDestinationType string

const (
	LogsDestinationTypeAkamaiObjectStorage LogsDestinationType = "akamai_object_storage"
)

// LogsDestinationStatus represents the status of a logs destination.
type LogsDestinationStatus string

const (
	LogsDestinationStatusActive   LogsDestinationStatus = "active"
	LogsDestinationStatusInactive LogsDestinationStatus = "inactive"
)

// LogsDestinationDetails represents the details block returned in a LogsDestination response.
type LogsDestinationDetails struct {
	AccessKeyID string `json:"access_key_id"`
	BucketName  string `json:"bucket_name"`
	Host        string `json:"host"`
	Path        string `json:"path"`
}

// LogsDestinationDetailsCreateOptions represents the details block used when creating a LogsDestination.
type LogsDestinationDetailsCreateOptions struct {
	AccessKeyID     string  `json:"access_key_id"`
	AccessKeySecret string  `json:"access_key_secret"`
	BucketName      string  `json:"bucket_name"`
	Host            string  `json:"host"`
	Path            *string `json:"path,omitempty"`
}

// LogsDestination represents a logs destination object.
type LogsDestination struct {
	Created   *time.Time             `json:"-"`
	CreatedBy string                 `json:"created_by"`
	Details   LogsDestinationDetails `json:"details"`
	ID        int                    `json:"id"`
	Label     string                 `json:"label"`
	Status    LogsDestinationStatus  `json:"status"`
	Type      LogsDestinationType    `json:"type"`
	Updated   *time.Time             `json:"-"`
	UpdatedBy string                 `json:"updated_by"`
	Version   int                    `json:"version"`
}

// UnmarshalJSON implements the json.Unmarshaler interface for LogsDestination.
func (i *LogsDestination) UnmarshalJSON(b []byte) error {
	type Mask LogsDestination

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

// LogsDestinationCreateOptions are the options used to create a new logs destination.
type LogsDestinationCreateOptions struct {
	Label   string                              `json:"label"`
	Type    LogsDestinationType                 `json:"type"`
	Details LogsDestinationDetailsCreateOptions `json:"details"`
}

// LogsDestinationUpdateOptions are the options used to update a logs destination.
type LogsDestinationUpdateOptions struct {
	Label   string                  `json:"label,omitempty"`
	Details *LogsDestinationDetails `json:"details,omitempty"`
}

// ListLogsDestinations returns a paginated list of logs destinations.
func (c *Client) ListLogsDestinations(ctx context.Context, opts *ListOptions) ([]LogsDestination, error) {
	return getPaginatedResults[LogsDestination](ctx, c, "monitor/streams/destinations", opts)
}

// GetLogsDestination gets a single logs destination by ID.
func (c *Client) GetLogsDestination(ctx context.Context, destinationID int) (*LogsDestination, error) {
	e := formatAPIPath("monitor/streams/destinations/%d", destinationID)
	return doGETRequest[LogsDestination](ctx, c, e)
}

// CreateLogsDestination creates a new logs destination.
func (c *Client) CreateLogsDestination(ctx context.Context, opts LogsDestinationCreateOptions) (*LogsDestination, error) {
	return doPOSTRequest[LogsDestination](ctx, c, "monitor/streams/destinations", opts)
}

// UpdateLogsDestination updates a logs destination.
func (c *Client) UpdateLogsDestination(ctx context.Context, destinationID int, opts LogsDestinationUpdateOptions) (*LogsDestination, error) {
	e := formatAPIPath("monitor/streams/destinations/%d", destinationID)
	return doPUTRequest[LogsDestination](ctx, c, e, opts)
}

// DeleteLogsDestination deletes a logs destination.
func (c *Client) DeleteLogsDestination(ctx context.Context, destinationID int) error {
	e := formatAPIPath("monitor/streams/destinations/%d", destinationID)
	return doDELETERequest(ctx, c, e)
}

// ListLogsDestinationHistory returns the version history for a logs destination.
func (c *Client) ListLogsDestinationHistory(ctx context.Context, destinationID int, opts *ListOptions) ([]LogsDestination, error) {
	e := formatAPIPath("monitor/streams/destinations/%d/history", destinationID)
	return getPaginatedResults[LogsDestination](ctx, c, e, opts)
}
