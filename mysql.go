package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

type MySQLDatabaseTarget string

type MySQLDatabaseMaintenanceWindow = DatabaseMaintenanceWindow

const (
	MySQLDatabaseTargetPrimary   MySQLDatabaseTarget = "primary"
	MySQLDatabaseTargetSecondary MySQLDatabaseTarget = "secondary"
)

// A MySQLDatabase is an instance of Linode MySQL Managed Databases
type MySQLDatabase struct {
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

	// Deprecated: ReplicationType is a deprecated property, as it is no longer supported in DBaaS V2.
	ReplicationType string `json:"replication_type"`
	// Deprecated: SSLConnection is a deprecated property, as it is no longer supported in DBaaS V2.
	SSLConnection bool `json:"ssl_connection"`
	// Deprecated: Encrypted is a deprecated property, as it is no longer supported in DBaaS V2.
	Encrypted bool `json:"encrypted"`

	AllowList         []string                  `json:"allow_list"`
	InstanceURI       string                    `json:"instance_uri"`
	Created           *time.Time                `json:"-"`
	Updated           *time.Time                `json:"-"`
	Updates           DatabaseMaintenanceWindow `json:"updates"`
	Fork              *DatabaseFork             `json:"fork"`
	OldestRestoreTime *time.Time                `json:"-"`
	UsedDiskSizeGB    int                       `json:"used_disk_size_gb"`
	TotalDiskSizeGB   int                       `json:"total_disk_size_gb"`
	Port              int                       `json:"port"`

	EngineConfig MySQLDatabaseEngineConfig `json:"engine_config"`
}

type MySQLDatabaseEngineConfig struct {
	MySQL                 *MySQLDatabaseEngineConfigMySQL `json:"mysql,omitempty"`
	BinlogRetentionPeriod *int                            `json:"binlog_retention_period,omitempty"`
}

type MySQLDatabaseEngineConfigMySQL struct {
	ConnectTimeout               *int     `json:"connect_timeout,omitempty"`
	DefaultTimeZone              *string  `json:"default_time_zone,omitempty"`
	GroupConcatMaxLen            *float64 `json:"group_concat_max_len,omitempty"`
	InformationSchemaStatsExpiry *int     `json:"information_schema_stats_expiry,omitempty"`
	InnoDBChangeBufferMaxSize    *int     `json:"innodb_change_buffer_max_size,omitempty"`
	InnoDBFlushNeighbors         *int     `json:"innodb_flush_neighbors,omitempty"`
	InnoDBFTMinTokenSize         *int     `json:"innodb_ft_min_token_size,omitempty"`
	InnoDBFTServerStopwordTable  **string `json:"innodb_ft_server_stopword_table,omitempty"`
	InnoDBLockWaitTimeout        *int     `json:"innodb_lock_wait_timeout,omitempty"`
	InnoDBLogBufferSize          *int     `json:"innodb_log_buffer_size,omitempty"`
	InnoDBOnlineAlterLogMaxSize  *int     `json:"innodb_online_alter_log_max_size,omitempty"`
	InnoDBReadIOThreads          *int     `json:"innodb_read_io_threads,omitempty"`
	InnoDBRollbackOnTimeout      *bool    `json:"innodb_rollback_on_timeout,omitempty"`
	InnoDBThreadConcurrency      *int     `json:"innodb_thread_concurrency,omitempty"`
	InnoDBWriteIOThreads         *int     `json:"innodb_write_io_threads,omitempty"`
	InteractiveTimeout           *int     `json:"interactive_timeout,omitempty"`
	InternalTmpMemStorageEngine  *string  `json:"internal_tmp_mem_storage_engine,omitempty"`
	MaxAllowedPacket             *int     `json:"max_allowed_packet,omitempty"`
	MaxHeapTableSize             *int     `json:"max_heap_table_size,omitempty"`
	NetBufferLength              *int     `json:"net_buffer_length,omitempty"`
	NetReadTimeout               *int     `json:"net_read_timeout,omitempty"`
	NetWriteTimeout              *int     `json:"net_write_timeout,omitempty"`
	SortBufferSize               *int     `json:"sort_buffer_size,omitempty"`
	SQLMode                      *string  `json:"sql_mode,omitempty"`
	SQLRequirePrimaryKey         *bool    `json:"sql_require_primary_key,omitempty"`
	TmpTableSize                 *int     `json:"tmp_table_size,omitempty"`
	WaitTimeout                  *int     `json:"wait_timeout,omitempty"`
}

type MySQLDatabaseConfigInfo struct {
	MySQL                 MySQLDatabaseConfigInfoMySQL                 `json:"mysql"`
	BinlogRetentionPeriod MySQLDatabaseConfigInfoBinlogRetentionPeriod `json:"binlog_retention_period"`
}

type MySQLDatabaseConfigInfoMySQL struct {
	ConnectTimeout               ConnectTimeout               `json:"connect_timeout"`
	DefaultTimeZone              DefaultTimeZone              `json:"default_time_zone"`
	GroupConcatMaxLen            GroupConcatMaxLen            `json:"group_concat_max_len"`
	InformationSchemaStatsExpiry InformationSchemaStatsExpiry `json:"information_schema_stats_expiry"`
	InnoDBChangeBufferMaxSize    InnoDBChangeBufferMaxSize    `json:"innodb_change_buffer_max_size"`
	InnoDBFlushNeighbors         InnoDBFlushNeighbors         `json:"innodb_flush_neighbors"`
	InnoDBFTMinTokenSize         InnoDBFTMinTokenSize         `json:"innodb_ft_min_token_size"`
	InnoDBFTServerStopwordTable  InnoDBFTServerStopwordTable  `json:"innodb_ft_server_stopword_table"`
	InnoDBLockWaitTimeout        InnoDBLockWaitTimeout        `json:"innodb_lock_wait_timeout"`
	InnoDBLogBufferSize          InnoDBLogBufferSize          `json:"innodb_log_buffer_size"`
	InnoDBOnlineAlterLogMaxSize  InnoDBOnlineAlterLogMaxSize  `json:"innodb_online_alter_log_max_size"`
	InnoDBReadIOThreads          InnoDBReadIOThreads          `json:"innodb_read_io_threads"`
	InnoDBRollbackOnTimeout      InnoDBRollbackOnTimeout      `json:"innodb_rollback_on_timeout"`
	InnoDBThreadConcurrency      InnoDBThreadConcurrency      `json:"innodb_thread_concurrency"`
	InnoDBWriteIOThreads         InnoDBWriteIOThreads         `json:"innodb_write_io_threads"`
	InteractiveTimeout           InteractiveTimeout           `json:"interactive_timeout"`
	InternalTmpMemStorageEngine  InternalTmpMemStorageEngine  `json:"internal_tmp_mem_storage_engine"`
	MaxAllowedPacket             MaxAllowedPacket             `json:"max_allowed_packet"`
	MaxHeapTableSize             MaxHeapTableSize             `json:"max_heap_table_size"`
	NetBufferLength              NetBufferLength              `json:"net_buffer_length"`
	NetReadTimeout               NetReadTimeout               `json:"net_read_timeout"`
	NetWriteTimeout              NetWriteTimeout              `json:"net_write_timeout"`
	SortBufferSize               SortBufferSize               `json:"sort_buffer_size"`
	SQLMode                      SQLMode                      `json:"sql_mode"`
	SQLRequirePrimaryKey         SQLRequirePrimaryKey         `json:"sql_require_primary_key"`
	TmpTableSize                 TmpTableSize                 `json:"tmp_table_size"`
	WaitTimeout                  WaitTimeout                  `json:"wait_timeout"`
}

type ConnectTimeout struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type DefaultTimeZone struct {
	Description     string `json:"description"`
	Example         string `json:"example"`
	MaxLength       int    `json:"maxLength"`
	MinLength       int    `json:"minLength"`
	Pattern         string `json:"pattern"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type GroupConcatMaxLen struct {
	Description     string  `json:"description"`
	Example         float64 `json:"example"`
	Maximum         float64 `json:"maximum"`
	Minimum         float64 `json:"minimum"`
	RequiresRestart bool    `json:"requires_restart"`
	Type            string  `json:"type"`
}

type InformationSchemaStatsExpiry struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type InnoDBChangeBufferMaxSize struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type InnoDBFlushNeighbors struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type InnoDBFTMinTokenSize struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type InnoDBFTServerStopwordTable struct {
	Description     string   `json:"description"`
	Example         string   `json:"example"`
	MaxLength       int      `json:"maxLength"`
	Pattern         string   `json:"pattern"`
	RequiresRestart bool     `json:"requires_restart"`
	Type            []string `json:"type"`
}

type InnoDBLockWaitTimeout struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type InnoDBLogBufferSize struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type InnoDBOnlineAlterLogMaxSize struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type InnoDBReadIOThreads struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type InnoDBRollbackOnTimeout struct {
	Description     string `json:"description"`
	Example         bool   `json:"example"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type InnoDBThreadConcurrency struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type InnoDBWriteIOThreads struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type InteractiveTimeout struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type InternalTmpMemStorageEngine struct {
	Description     string   `json:"description"`
	Enum            []string `json:"enum"`
	Example         string   `json:"example"`
	RequiresRestart bool     `json:"requires_restart"`
	Type            string   `json:"type"`
}

type MaxAllowedPacket struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type MaxHeapTableSize struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type NetBufferLength struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type NetReadTimeout struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type NetWriteTimeout struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type SortBufferSize struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type SQLMode struct {
	Description     string `json:"description"`
	Example         string `json:"example"`
	MaxLength       int    `json:"maxLength"`
	Pattern         string `json:"pattern"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type SQLRequirePrimaryKey struct {
	Description     string `json:"description"`
	Example         bool   `json:"example"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type TmpTableSize struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type WaitTimeout struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

type MySQLDatabaseConfigInfoBinlogRetentionPeriod struct {
	Description     string `json:"description"`
	Example         int    `json:"example"`
	Maximum         int    `json:"maximum"`
	Minimum         int    `json:"minimum"`
	RequiresRestart bool   `json:"requires_restart"`
	Type            string `json:"type"`
}

func (d *MySQLDatabase) UnmarshalJSON(b []byte) error {
	type Mask MySQLDatabase

	p := struct {
		*Mask
		Created           *parseabletime.ParseableTime `json:"created"`
		Updated           *parseabletime.ParseableTime `json:"updated"`
		OldestRestoreTime *parseabletime.ParseableTime `json:"oldest_restore_time"`
	}{
		Mask: (*Mask)(d),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	d.Created = (*time.Time)(p.Created)
	d.Updated = (*time.Time)(p.Updated)
	d.OldestRestoreTime = (*time.Time)(p.OldestRestoreTime)
	return nil
}

// MySQLCreateOptions fields are used when creating a new MySQL Database
type MySQLCreateOptions struct {
	Label       string   `json:"label"`
	Region      string   `json:"region"`
	Type        string   `json:"type"`
	Engine      string   `json:"engine"`
	AllowList   []string `json:"allow_list,omitempty"`
	ClusterSize int      `json:"cluster_size,omitempty"`

	// Deprecated: ReplicationType is a deprecated property, as it is no longer supported in DBaaS V2.
	ReplicationType string `json:"replication_type,omitempty"`
	// Deprecated: Encrypted is a deprecated property, as it is no longer supported in DBaaS V2.
	Encrypted bool `json:"encrypted,omitempty"`
	// Deprecated: SSLConnection is a deprecated property, as it is no longer supported in DBaaS V2.
	SSLConnection bool `json:"ssl_connection,omitempty"`

	Fork         *DatabaseFork              `json:"fork,omitempty"`
	EngineConfig *MySQLDatabaseEngineConfig `json:"engine_config,omitempty"`
}

// MySQLUpdateOptions fields are used when altering the existing MySQL Database
type MySQLUpdateOptions struct {
	Label        string                     `json:"label,omitempty"`
	AllowList    *[]string                  `json:"allow_list,omitempty"`
	Updates      *DatabaseMaintenanceWindow `json:"updates,omitempty"`
	Type         string                     `json:"type,omitempty"`
	ClusterSize  int                        `json:"cluster_size,omitempty"`
	Version      string                     `json:"version,omitempty"`
	EngineConfig *MySQLDatabaseEngineConfig `json:"engine_config,omitempty"`
}

// MySQLDatabaseBackup is information for interacting with a backup for the existing MySQL Database
// Deprecated: MySQLDatabaseBackup is a deprecated struct, as the backup endpoints are no longer supported in DBaaS V2.
// In DBaaS V2, databases can be backed up via database forking.
type MySQLDatabaseBackup struct {
	ID      int        `json:"id"`
	Label   string     `json:"label"`
	Type    string     `json:"type"`
	Created *time.Time `json:"-"`
}

// MySQLBackupCreateOptions are options used for CreateMySQLDatabaseBackup(...)
// Deprecated: MySQLBackupCreateOptions is a deprecated struct, as the backup endpoints are no longer supported in DBaaS V2.
// In DBaaS V2, databases can be backed up via database forking.
type MySQLBackupCreateOptions struct {
	Label  string              `json:"label"`
	Target MySQLDatabaseTarget `json:"target"`
}

func (d *MySQLDatabaseBackup) UnmarshalJSON(b []byte) error {
	type Mask MySQLDatabaseBackup

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
	}{
		Mask: (*Mask)(d),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	d.Created = (*time.Time)(p.Created)
	return nil
}

// MySQLDatabaseCredential is the Root Credentials to access the Linode Managed Database
type MySQLDatabaseCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// MySQLDatabaseSSL is the SSL Certificate to access the Linode Managed MySQL Database
type MySQLDatabaseSSL struct {
	CACertificate []byte `json:"ca_certificate"`
}

// ListMySQLDatabases lists all MySQL Databases associated with the account
func (c *Client) ListMySQLDatabases(ctx context.Context, opts *ListOptions) ([]MySQLDatabase, error) {
	return getPaginatedResults[MySQLDatabase](ctx, c, "databases/mysql/instances", opts)
}

// ListMySQLDatabaseBackups lists all MySQL Database Backups associated with the given MySQL Database
// Deprecated: ListMySQLDatabaseBackups is a deprecated method, as the backup endpoints are no longer supported in DBaaS V2.
// In DBaaS V2, databases can be backed up via database forking.
func (c *Client) ListMySQLDatabaseBackups(ctx context.Context, databaseID int, opts *ListOptions) ([]MySQLDatabaseBackup, error) {
	return getPaginatedResults[MySQLDatabaseBackup](ctx, c, formatAPIPath("databases/mysql/instances/%d/backups", databaseID), opts)
}

// GetMySQLDatabase returns a single MySQL Database matching the id
func (c *Client) GetMySQLDatabase(ctx context.Context, databaseID int) (*MySQLDatabase, error) {
	e := formatAPIPath("databases/mysql/instances/%d", databaseID)
	return doGETRequest[MySQLDatabase](ctx, c, e)
}

// CreateMySQLDatabase creates a new MySQL Database using the createOpts as configuration, returns the new MySQL Database
func (c *Client) CreateMySQLDatabase(ctx context.Context, opts MySQLCreateOptions) (*MySQLDatabase, error) {
	return doPOSTRequest[MySQLDatabase](ctx, c, "databases/mysql/instances", opts)
}

// DeleteMySQLDatabase deletes an existing MySQL Database with the given id
func (c *Client) DeleteMySQLDatabase(ctx context.Context, databaseID int) error {
	e := formatAPIPath("databases/mysql/instances/%d", databaseID)
	return doDELETERequest(ctx, c, e)
}

// UpdateMySQLDatabase updates the given MySQL Database with the provided opts, returns the MySQLDatabase with the new settings
func (c *Client) UpdateMySQLDatabase(ctx context.Context, databaseID int, opts MySQLUpdateOptions) (*MySQLDatabase, error) {
	e := formatAPIPath("databases/mysql/instances/%d", databaseID)
	return doPUTRequest[MySQLDatabase](ctx, c, e, opts)
}

// GetMySQLDatabaseSSL returns the SSL Certificate for the given MySQL Database
func (c *Client) GetMySQLDatabaseSSL(ctx context.Context, databaseID int) (*MySQLDatabaseSSL, error) {
	e := formatAPIPath("databases/mysql/instances/%d/ssl", databaseID)
	return doGETRequest[MySQLDatabaseSSL](ctx, c, e)
}

// GetMySQLDatabaseCredentials returns the Root Credentials for the given MySQL Database
func (c *Client) GetMySQLDatabaseCredentials(ctx context.Context, databaseID int) (*MySQLDatabaseCredential, error) {
	e := formatAPIPath("databases/mysql/instances/%d/credentials", databaseID)
	return doGETRequest[MySQLDatabaseCredential](ctx, c, e)
}

// ResetMySQLDatabaseCredentials returns the Root Credentials for the given MySQL Database (may take a few seconds to work)
func (c *Client) ResetMySQLDatabaseCredentials(ctx context.Context, databaseID int) error {
	e := formatAPIPath("databases/mysql/instances/%d/credentials/reset", databaseID)
	return doPOSTRequestNoRequestResponseBody(ctx, c, e)
}

// GetMySQLDatabaseBackup returns a specific MySQL Database Backup with the given ids
// Deprecated: GetMySQLDatabaseBackup is a deprecated method, as the backup endpoints are no longer supported in DBaaS V2.
// In DBaaS V2, databases can be backed up via database forking.
func (c *Client) GetMySQLDatabaseBackup(ctx context.Context, databaseID int, backupID int) (*MySQLDatabaseBackup, error) {
	e := formatAPIPath("databases/mysql/instances/%d/backups/%d", databaseID, backupID)
	return doGETRequest[MySQLDatabaseBackup](ctx, c, e)
}

// RestoreMySQLDatabaseBackup returns the given MySQL Database with the given Backup
// Deprecated: RestoreMySQLDatabaseBackup is a deprecated method, as the backup endpoints are no longer supported in DBaaS V2.
// In DBaaS V2, databases can be backed up via database forking.
func (c *Client) RestoreMySQLDatabaseBackup(ctx context.Context, databaseID int, backupID int) error {
	e := formatAPIPath("databases/mysql/instances/%d/backups/%d/restore", databaseID, backupID)
	return doPOSTRequestNoRequestResponseBody(ctx, c, e)
}

// CreateMySQLDatabaseBackup creates a snapshot for the given MySQL database
// Deprecated: CreateMySQLDatabaseBackup is a deprecated method, as the backup endpoints are no longer supported in DBaaS V2.
// In DBaaS V2, databases can be backed up via database forking.
func (c *Client) CreateMySQLDatabaseBackup(ctx context.Context, databaseID int, opts MySQLBackupCreateOptions) error {
	e := formatAPIPath("databases/mysql/instances/%d/backups", databaseID)
	return doPOSTRequestNoResponseBody(ctx, c, e, opts)
}

// PatchMySQLDatabase applies security patches and updates to the underlying operating system of the Managed MySQL Database
func (c *Client) PatchMySQLDatabase(ctx context.Context, databaseID int) error {
	e := formatAPIPath("databases/mysql/instances/%d/patch", databaseID)
	return doPOSTRequestNoRequestResponseBody(ctx, c, e)
}

// SuspendMySQLDatabase suspends a MySQL Managed Database, releasing idle resources and keeping only necessary data.
// All service data is lost if there are no backups available.
func (c *Client) SuspendMySQLDatabase(ctx context.Context, databaseID int) error {
	e := formatAPIPath("databases/mysql/instances/%d/suspend", databaseID)
	return doPOSTRequestNoRequestResponseBody(ctx, c, e)
}

// ResumeMySQLDatabase resumes a suspended MySQL Managed Database
func (c *Client) ResumeMySQLDatabase(ctx context.Context, databaseID int) error {
	e := formatAPIPath("databases/mysql/instances/%d/resume", databaseID)
	return doPOSTRequestNoRequestResponseBody(ctx, c, e)
}

// GetMySQLDatabaseConfig returns a detailed list of all the configuration options for MySQL Databases
func (c *Client) GetMySQLDatabaseConfig(ctx context.Context) (*MySQLDatabaseConfigInfo, error) {
	return doGETRequest[MySQLDatabaseConfigInfo](ctx, c, "databases/mysql/config")
}
