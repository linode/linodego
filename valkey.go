package linodego

import (
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

	d.AvailableRestoreTimes = make([]time.Time, len(p.AvailableRestoreTimes))
	for i, t := range p.AvailableRestoreTimes {
		d.AvailableRestoreTimes[i] = time.Time(t)
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
