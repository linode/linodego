package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/v2/internal/parseabletime"
)

// A ValkeyDatabase is an instance of Linode Valkey Managed Databases
type ValkeyDatabase struct {
	ID          int              `json:"id"`
	Status      DatabaseStatus   `json:"status"`
	Label       string           `json:"label"`
	Hosts       DatabaseHost     `json:"hosts"`
	Region      string           `json:"region"`
	Type        string           `json:"type"`
	Engine      string           `json:"engine"`
	Version     string           `json:"version"`
	ClusterSize int              `json:"cluster_size"`
	Platform    DatabasePlatform `json:"platform"`

	// Members has dynamic keys so it is a map
	Members map[string]DatabaseMemberType `json:"members"`

	SSLConnection         bool                      `json:"ssl_connection"`
	Encrypted             bool                      `json:"encrypted"`
	AllowList             []string                  `json:"allow_list"`
	Created               *time.Time                `json:"-"`
	Updated               *time.Time                `json:"-"`
	Updates               DatabaseMaintenanceWindow `json:"updates"`
	Fork                  *DatabaseFork             `json:"fork"`
	AvailableRestoreTimes []time.Time               `json:"-"`
	UsedDiskSizeGB        int                       `json:"used_disk_size_gb"`
	TotalDiskSizeGB       int                       `json:"total_disk_size_gb"`
	Port                  int                       `json:"port"`

	EngineConfig   ValkeyDatabaseEngineConfig `json:"engine_config"`
	PrivateNetwork *DatabasePrivateNetwork    `json:"private_network,omitzero"`
}

type ValkeyDatabaseEngineConfig struct {
	BackupHour                          *int     `json:"backup_hour,omitzero"`
	BackupMinute                        *int     `json:"backup_minute,omitzero"`
	FrequentSnapshots                   *bool    `json:"frequent_snapshots,omitzero"`
	ValkeyACLChannelsDefault            *string  `json:"valkey_acl_channels_default,omitzero"`
	ValkeyActiveExpireEffort            *int     `json:"valkey_active_expire_effort,omitzero"`
	ValkeyActiveDefrag                  *bool    `json:"valkey_activedefrag,omitzero"`
	ValkeyLFUDecayTime                  *int     `json:"valkey_lfu_decay_time,omitzero"`
	ValkeyLFULogFactor                  *int     `json:"valkey_lfu_log_factor,omitzero"`
	ValkeyMaxmemoryPolicy               **string `json:"valkey_maxmemory_policy,omitzero"`
	ValkeyNumberOfDatabases             *int     `json:"valkey_number_of_databases,omitzero"`
	ValkeyPersistence                   *string  `json:"valkey_persistence,omitzero"`
	ValkeyPubsubClientOutputBufferLimit *int     `json:"valkey_pubsub_client_output_buffer_limit,omitzero"`
	ValkeyTimeout                       *int     `json:"valkey_timeout,omitzero"`
}

func (d *ValkeyDatabase) UnmarshalJSON(b []byte) error {
	type Mask ValkeyDatabase

	p := struct {
		*Mask

		Created               *parseabletime.ParseableTime  `json:"created"`
		Updated               *parseabletime.ParseableTime  `json:"updated"`
		AvailableRestoreTimes []parseabletime.ParseableTime `json:"available_restore_times"`
	}{
		Mask: (*Mask)(d),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	d.Created = (*time.Time)(p.Created)
	d.Updated = (*time.Time)(p.Updated)

	if p.AvailableRestoreTimes != nil {
		d.AvailableRestoreTimes = make([]time.Time, len(p.AvailableRestoreTimes))
		for i, t := range p.AvailableRestoreTimes {
			d.AvailableRestoreTimes[i] = time.Time(t)
		}
	}

	return nil
}

// ValkeyCreateOptions fields are used when creating a new Valkey Database
type ValkeyCreateOptions struct {
	Label       string   `json:"label"`
	Region      string   `json:"region"`
	Type        string   `json:"type"`
	Engine      string   `json:"engine"`
	AllowList   []string `json:"allow_list,omitzero"`
	ClusterSize int      `json:"cluster_size,omitzero"`

	Fork           *DatabaseFork               `json:"fork,omitzero"`
	EngineConfig   *ValkeyDatabaseEngineConfig `json:"engine_config,omitzero"`
	PrivateNetwork *DatabasePrivateNetwork     `json:"private_network,omitzero"`
}

// ValkeyUpdateOptions fields are used when altering the existing Valkey Database
type ValkeyUpdateOptions struct {
	Label          string                      `json:"label,omitzero"`
	Region         string                      `json:"region,omitzero"`
	AllowList      []string                    `json:"allow_list,omitzero"`
	Updates        *DatabaseMaintenanceWindow  `json:"updates,omitzero"`
	Type           string                      `json:"type,omitzero"`
	ClusterSize    int                         `json:"cluster_size,omitzero"`
	Version        string                      `json:"version,omitzero"`
	EngineConfig   *ValkeyDatabaseEngineConfig `json:"engine_config,omitzero"`
	PrivateNetwork **DatabasePrivateNetwork    `json:"private_network,omitzero"`
}

// ValkeyDatabaseCredential is the Root Credentials to access the Linode Managed Database
type ValkeyDatabaseCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ValkeyDatabaseSSL is the SSL Certificate to access the Linode Managed Valkey Database
type ValkeyDatabaseSSL struct {
	CACertificate []byte `json:"ca_certificate"`
}

// ConfigParamType handles API metadata "type" values that may be either a single
// string or an array of strings.
type ConfigParamType []string

func (t *ConfigParamType) UnmarshalJSON(b []byte) error {
	var single string
	if err := json.Unmarshal(b, &single); err == nil {
		*t = ConfigParamType{single}
		return nil
	}

	var multi []string
	if err := json.Unmarshal(b, &multi); err != nil {
		return err
	}

	*t = ConfigParamType(multi)

	return nil
}

type ValkeyDatabaseConfigInfo struct {
	BackupHour                          ValkeyDatabaseConfigInfoOption `json:"backup_hour"`
	BackupMinute                        ValkeyDatabaseConfigInfoOption `json:"backup_minute"`
	FrequentSnapshots                   ValkeyDatabaseConfigInfoOption `json:"frequent_snapshots"`
	ValkeyACLChannelsDefault            ValkeyDatabaseConfigInfoOption `json:"valkey_acl_channels_default"`
	ValkeyActiveExpireEffort            ValkeyDatabaseConfigInfoOption `json:"valkey_active_expire_effort"`
	ValkeyActiveDefrag                  ValkeyDatabaseConfigInfoOption `json:"valkey_activedefrag"`
	ValkeyLFUDecayTime                  ValkeyDatabaseConfigInfoOption `json:"valkey_lfu_decay_time"`
	ValkeyLFULogFactor                  ValkeyDatabaseConfigInfoOption `json:"valkey_lfu_log_factor"`
	ValkeyMaxmemoryPolicy               ValkeyDatabaseConfigInfoOption `json:"valkey_maxmemory_policy"`
	ValkeyNumberOfDatabases             ValkeyDatabaseConfigInfoOption `json:"valkey_number_of_databases"`
	ValkeyPersistence                   ValkeyDatabaseConfigInfoOption `json:"valkey_persistence"`
	ValkeyPubsubClientOutputBufferLimit ValkeyDatabaseConfigInfoOption `json:"valkey_pubsub_client_output_buffer_limit"`
	ValkeyTimeout                       ValkeyDatabaseConfigInfoOption `json:"valkey_timeout"`
}

type ValkeyDatabaseConfigInfoOption struct {
	Description     string          `json:"description"`
	Example         any             `json:"example,omitzero"`
	Maximum         *float64        `json:"maximum,omitzero"`
	Minimum         *float64        `json:"minimum,omitzero"`
	Default         any             `json:"default,omitzero"`
	Enum            []string        `json:"enum,omitzero"`
	RequiresRestart bool            `json:"requires_restart"`
	Type            ConfigParamType `json:"type"`
}

// GetValkeyDatabaseConfig returns the catalog of configuration options for Valkey databases.
func (c *Client) GetValkeyDatabaseConfig(ctx context.Context) (*ValkeyDatabaseConfigInfo, error) {
	return doGETRequest[ValkeyDatabaseConfigInfo](ctx, c, "databases/valkey/config")
}

// ListValkeyDatabases lists all Valkey Databases associated with the account
func (c *Client) ListValkeyDatabases(ctx context.Context, opts *ListOptions) ([]ValkeyDatabase, error) {
	return getPaginatedResults[ValkeyDatabase](ctx, c, "databases/valkey/instances", opts)
}

// GetValkeyDatabase returns a single Valkey Database matching the id
func (c *Client) GetValkeyDatabase(ctx context.Context, databaseID int) (*ValkeyDatabase, error) {
	e := formatAPIPath("databases/valkey/instances/%d", databaseID)
	return doGETRequest[ValkeyDatabase](ctx, c, e)
}

// CreateValkeyDatabase creates a new Valkey Database using the createOpts as configuration, returns the new Valkey Database
func (c *Client) CreateValkeyDatabase(ctx context.Context, opts ValkeyCreateOptions) (*ValkeyDatabase, error) {
	return doPOSTRequest[ValkeyDatabase](ctx, c, "databases/valkey/instances", opts)
}

// DeleteValkeyDatabase deletes an existing Valkey Database with the given id
func (c *Client) DeleteValkeyDatabase(ctx context.Context, databaseID int) error {
	e := formatAPIPath("databases/valkey/instances/%d", databaseID)
	return doDELETERequest(ctx, c, e)
}

// UpdateValkeyDatabase updates the given Valkey Database with the provided opts, returns the ValkeyDatabase with the new settings
func (c *Client) UpdateValkeyDatabase(ctx context.Context, databaseID int, opts ValkeyUpdateOptions) (*ValkeyDatabase, error) {
	e := formatAPIPath("databases/valkey/instances/%d", databaseID)
	return doPUTRequest[ValkeyDatabase](ctx, c, e, opts)
}

// GetValkeyDatabaseSSL returns the SSL Certificate for the given Valkey Database
func (c *Client) GetValkeyDatabaseSSL(ctx context.Context, databaseID int) (*ValkeyDatabaseSSL, error) {
	e := formatAPIPath("databases/valkey/instances/%d/ssl", databaseID)
	return doGETRequest[ValkeyDatabaseSSL](ctx, c, e)
}

// GetValkeyDatabaseCredentials returns the Root Credentials for the given Valkey Database
func (c *Client) GetValkeyDatabaseCredentials(ctx context.Context, databaseID int) (*ValkeyDatabaseCredential, error) {
	e := formatAPIPath("databases/valkey/instances/%d/credentials", databaseID)
	return doGETRequest[ValkeyDatabaseCredential](ctx, c, e)
}

// ResetValkeyDatabaseCredentials resets the Root Credentials for the given Valkey Database
func (c *Client) ResetValkeyDatabaseCredentials(ctx context.Context, databaseID int) error {
	e := formatAPIPath("databases/valkey/instances/%d/credentials/reset", databaseID)
	return doPOSTRequestNoRequestResponseBody(ctx, c, e)
}

// PatchValkeyDatabase applies security patches and updates to the underlying operating system of the Managed Valkey Database
func (c *Client) PatchValkeyDatabase(ctx context.Context, databaseID int) error {
	e := formatAPIPath("databases/valkey/instances/%d/patch", databaseID)
	return doPOSTRequestNoRequestResponseBody(ctx, c, e)
}

// SuspendValkeyDatabase suspends a Valkey Managed Database, releasing idle resources and keeping only necessary data.
// All service data is lost if there are no backups available.
func (c *Client) SuspendValkeyDatabase(ctx context.Context, databaseID int) error {
	e := formatAPIPath("databases/valkey/instances/%d/suspend", databaseID)
	return doPOSTRequestNoRequestResponseBody(ctx, c, e)
}

// ResumeValkeyDatabase resumes a suspended Valkey Managed Database
func (c *Client) ResumeValkeyDatabase(ctx context.Context, databaseID int) error {
	e := formatAPIPath("databases/valkey/instances/%d/resume", databaseID)
	return doPOSTRequestNoRequestResponseBody(ctx, c, e)
}
