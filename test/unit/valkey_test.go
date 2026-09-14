package unit

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func getFixtureBytes(t *testing.T, name string) []byte {
	t.Helper()

	fixtureData, err := fixtures.GetFixture(name)
	assert.NoError(t, err)

	var data []byte
	switch v := fixtureData.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	case map[string]interface{}:
		data, err = json.Marshal(v)
		assert.NoError(t, err, "Failed to marshal fixtureData")
	default:
		assert.Fail(t, "Unexpected fixtureData type")
	}

	return data
}

func TestUnmarshalValkeyDatabase(t *testing.T) {
	data := getFixtureBytes(t, "valkey_database_unmarshal")

	var db linodego.ValkeyDatabase
	err := json.Unmarshal(data, &db)
	assert.NoError(t, err)

	assert.Equal(t, 495395, db.ID)
	assert.Equal(t, linodego.DatabaseStatusActive, db.Status)
	assert.Equal(t, "valkey", db.Engine)
	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, linodego.DatabasePlatformRDBMSDefault, db.Platform)
	assert.NotNil(t, db.Created)
	assert.NotNil(t, db.Updated)
	assert.Equal(t, []time.Time{
		time.Date(2026, time.August, 28, 19, 37, 9, 0, time.UTC),
		time.Date(2026, time.September, 11, 7, 54, 23, 0, time.UTC),
	}, db.AvailableRestoreTimes)
	assert.Nil(t, db.Fork)
	assert.NotNil(t, db.PrivateNetwork)
	assert.Equal(t, 570237, db.PrivateNetwork.VPCID)

	assert.NotNil(t, db.EngineConfig.ValkeyMaxmemoryPolicy)
	assert.Equal(t, "noeviction", **db.EngineConfig.ValkeyMaxmemoryPolicy)
	assert.NotNil(t, db.EngineConfig.FrequentSnapshots)
	assert.True(t, *db.EngineConfig.FrequentSnapshots)
	assert.Nil(t, db.EngineConfig.BackupHour)
	assert.Nil(t, db.EngineConfig.BackupMinute)
}

func TestUnmarshalValkeyDatabase_EngineConfigEdgeCases(t *testing.T) {
	data := getFixtureBytes(t, "valkey_database_unmarshal_engine_config")

	var db linodego.ValkeyDatabase
	err := json.Unmarshal(data, &db)
	assert.NoError(t, err)

	assert.Empty(t, db.AvailableRestoreTimes)
	assert.Nil(t, db.PrivateNetwork)

	cfg := db.EngineConfig
	assert.Equal(t, 3, *cfg.BackupHour)
	assert.Equal(t, 30, *cfg.BackupMinute)
	assert.False(t, *cfg.FrequentSnapshots)
	assert.Equal(t, "resetchannels", *cfg.ValkeyACLChannelsDefault)
	assert.Equal(t, 5, *cfg.ValkeyActiveExpireEffort)
	assert.True(t, *cfg.ValkeyActiveDefrag)
	assert.Equal(t, 10, *cfg.ValkeyLFUDecayTime)
	assert.Equal(t, 20, *cfg.ValkeyLFULogFactor)
	// An explicit JSON null for a nullable field is indistinguishable from an
	// absent field once unmarshaled; both leave the outer pointer nil.
	assert.Nil(t, cfg.ValkeyMaxmemoryPolicy)
	assert.Equal(t, 32, *cfg.ValkeyNumberOfDatabases)
	assert.Equal(t, "off", *cfg.ValkeyPersistence)
	assert.Equal(t, 64, *cfg.ValkeyPubsubClientOutputBufferLimit)
	assert.Equal(t, 600, *cfg.ValkeyTimeout)
}

func TestUnmarshalValkeyDatabaseConfigInfo(t *testing.T) {
	payload := []byte(`{
		"backup_hour": {
			"description": "Hour of day for backups",
			"example": 3,
			"minimum": 0,
			"maximum": 23,
			"requires_restart": false,
			"type": ["integer", "null"]
		},
		"valkey_acl_channels_default": {
			"description": "Channel ACL default",
			"example": "allchannels",
			"enum": ["allchannels", "resetchannels"],
			"requires_restart": false,
			"type": "string"
		}
	}`)

	var config linodego.ValkeyDatabaseConfigInfo
	err := json.Unmarshal(payload, &config)
	assert.NoError(t, err)

	assert.Equal(t, linodego.ConfigParamType{"integer", "null"}, config.BackupHour.Type)
	assert.Equal(t, "string", config.ValkeyACLChannelsDefault.Type[0])
	assert.Equal(t, []string{"allchannels", "resetchannels"}, config.ValkeyACLChannelsDefault.Enum)
	assert.Equal(t, linodego.DatabaseEngineTypeValkey, linodego.DatabaseEngineTypeValkey)
}

func TestListDatabaseValkey_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("valkey_databases_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/valkey/instances", fixtureData)
	databases, err := base.Client.ListValkeyDatabases(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, databases, "Expected non-empty valkey database list")
}

func TestDatabaseValkey_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("valkey_database_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/valkey/instances/123", fixtureData)

	db, err := base.Client.GetValkeyDatabase(context.Background(), 123)
	assert.NoError(t, err)

	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, "valkey", db.Engine)
	assert.Equal(t, 123, db.ID)
	assert.Equal(t, "example-db", db.Label)
	assert.Equal(t, linodego.DatabasePlatform("rdbms-default"), db.Platform)
	assert.Equal(t, 6379, db.Port)
	assert.Equal(t, "us-east", db.Region)
	assert.Equal(t, linodego.DatabaseStatus("active"), db.Status)
	assert.Equal(t, 15, db.TotalDiskSizeGB)
	assert.Equal(t, "g6-dedicated-2", db.Type)
	assert.Equal(t, 3, db.Updates.Duration)
	assert.Equal(t, linodego.DatabaseDayOfWeek(1), db.Updates.DayOfWeek)
	assert.Equal(t, linodego.DatabaseMaintenanceFrequency("weekly"), db.Updates.Frequency)
	assert.Equal(t, 0, db.Updates.HourOfDay)
	assert.Equal(t, 2, db.UsedDiskSizeGB)
	assert.Equal(t, "8.1.9", db.Version)
	assert.Equal(t, 300, *db.EngineConfig.ValkeyTimeout)
	assert.Equal(t, "resetchannels", *db.EngineConfig.ValkeyACLChannelsDefault)
	assert.Equal(t, true, *db.EngineConfig.FrequentSnapshots)
	assert.Equal(t, "noeviction", **db.EngineConfig.ValkeyMaxmemoryPolicy)
	assert.Equal(t, 1234, db.PrivateNetwork.VPCID)
	assert.Equal(t, 5678, db.PrivateNetwork.SubnetID)
	assert.Equal(t, true, db.PrivateNetwork.PublicAccess)
}

func TestDatabaseValkey_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("valkey_database_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.ValkeyUpdateOptions{
		Label: "example-db-updated",
		EngineConfig: &linodego.ValkeyDatabaseEngineConfig{
			ValkeyTimeout: linodego.Pointer(600),
		},
		PrivateNetwork: linodego.DoublePointer(
			linodego.DatabasePrivateNetwork{
				VPCID:        1234,
				SubnetID:     5678,
				PublicAccess: true,
			},
		),
	}

	base.MockPut("databases/valkey/instances/123", fixtureData)

	db, err := base.Client.UpdateValkeyDatabase(context.Background(), 123, requestData)
	assert.NoError(t, err)

	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, "valkey", db.Engine)
	assert.Equal(t, 123, db.ID)
	assert.Equal(t, "example-db-updated", db.Label)
	assert.Equal(t, linodego.DatabasePlatform("rdbms-default"), db.Platform)
	assert.Equal(t, 6379, db.Port)
	assert.Equal(t, "us-east", db.Region)
	assert.Equal(t, linodego.DatabaseStatus("active"), db.Status)
	assert.Equal(t, 15, db.TotalDiskSizeGB)
	assert.Equal(t, "g6-dedicated-2", db.Type)
	assert.Equal(t, 3, db.Updates.Duration)
	assert.Equal(t, linodego.DatabaseDayOfWeek(1), db.Updates.DayOfWeek)
	assert.Equal(t, linodego.DatabaseMaintenanceFrequency("weekly"), db.Updates.Frequency)
	assert.Equal(t, 0, db.Updates.HourOfDay)
	assert.Equal(t, 2, db.UsedDiskSizeGB)
	assert.Equal(t, "8.1.9", db.Version)
	assert.Equal(t, 300, *db.EngineConfig.ValkeyTimeout)
	assert.Equal(t, "resetchannels", *db.EngineConfig.ValkeyACLChannelsDefault)
	assert.Equal(t, true, *db.EngineConfig.FrequentSnapshots)
	assert.Equal(t, "noeviction", **db.EngineConfig.ValkeyMaxmemoryPolicy)
	assert.Equal(t, 1234, db.PrivateNetwork.VPCID)
	assert.Equal(t, 5678, db.PrivateNetwork.SubnetID)
	assert.Equal(t, true, db.PrivateNetwork.PublicAccess)
}

func TestDatabaseValkey_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("valkey_database_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.ValkeyCreateOptions{
		Label:  "example-db-created",
		Region: "us-east",
		Type:   "g6-dedicated-2",
		Engine: "valkey",
		EngineConfig: &linodego.ValkeyDatabaseEngineConfig{
			ValkeyTimeout: linodego.Pointer(300),
		},
		PrivateNetwork: &linodego.DatabasePrivateNetwork{
			VPCID:        1234,
			SubnetID:     5678,
			PublicAccess: true,
		},
	}

	base.MockPost("databases/valkey/instances", fixtureData)

	db, err := base.Client.CreateValkeyDatabase(context.Background(), requestData)
	assert.NoError(t, err)

	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, "valkey", db.Engine)
	assert.Equal(t, 123, db.ID)
	assert.Equal(t, "example-db-created", db.Label)
	assert.Equal(t, linodego.DatabasePlatform("rdbms-default"), db.Platform)
	assert.Equal(t, 6379, db.Port)
	assert.Equal(t, "us-east", db.Region)
	assert.Equal(t, linodego.DatabaseStatus("active"), db.Status)
	assert.Equal(t, 15, db.TotalDiskSizeGB)
	assert.Equal(t, "g6-dedicated-2", db.Type)
	assert.Equal(t, 3, db.Updates.Duration)
	assert.Equal(t, linodego.DatabaseDayOfWeek(1), db.Updates.DayOfWeek)
	assert.Equal(t, linodego.DatabaseMaintenanceFrequency("weekly"), db.Updates.Frequency)
	assert.Equal(t, 0, db.Updates.HourOfDay)
	assert.Equal(t, 2, db.UsedDiskSizeGB)
	assert.Equal(t, "8.1.9", db.Version)
	assert.Equal(t, 300, *db.EngineConfig.ValkeyTimeout)
	assert.Equal(t, "resetchannels", *db.EngineConfig.ValkeyACLChannelsDefault)
	assert.Equal(t, true, *db.EngineConfig.FrequentSnapshots)
	assert.Equal(t, "noeviction", **db.EngineConfig.ValkeyMaxmemoryPolicy)
	assert.Equal(t, 1234, db.PrivateNetwork.VPCID)
	assert.Equal(t, 5678, db.PrivateNetwork.SubnetID)
	assert.Equal(t, true, db.PrivateNetwork.PublicAccess)
}

func TestDatabaseValkey_Delete(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, "databases/valkey/instances/123"), httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteValkeyDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseValkey_SSL_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("valkey_database_ssl_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/valkey/instances/123/ssl", fixtureData)

	ssl, err := base.Client.GetValkeyDatabaseSSL(context.Background(), 123)
	if assert.NoError(t, err) {
		expectedCACertificate := []byte("-----BEGIN CERTIFICATE-----\nThis is a test certificate\n-----END CERTIFICATE-----\n")
		assert.Equal(t, expectedCACertificate, ssl.CACertificate)
	}
}

func TestDatabaseValkey_Credentials_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("valkey_database_credentials_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/valkey/instances/123/credentials", fixtureData)

	creds, err := base.Client.GetValkeyDatabaseCredentials(context.Background(), 123)
	assert.NoError(t, err)
	assert.Equal(t, "linroot", creds.Username)
	assert.Equal(t, "s3cur3P@ssw0rd", creds.Password)
}

func TestDatabaseValkey_Credentials_Reset(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/valkey/instances/123/credentials/reset"), httpmock.NewStringResponder(200, "{}"))

	if err := client.ResetValkeyDatabaseCredentials(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseValkey_Patch(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/valkey/instances/123/patch"), httpmock.NewStringResponder(200, "{}"))

	if err := client.PatchValkeyDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseValkey_Suspend(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/valkey/instances/123/suspend"), httpmock.NewStringResponder(200, "{}"))

	if err := client.SuspendValkeyDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseValkey_Resume(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/valkey/instances/123/resume"), httpmock.NewStringResponder(200, "{}"))

	if err := client.ResumeValkeyDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseValkeyConfig_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("valkey_database_config_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/valkey/config", fixtureData)

	config, err := base.Client.GetValkeyDatabaseConfig(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, "Hour of day for backups.", config.BackupHour.Description)
	assert.Equal(t, float64(3), config.BackupHour.Example)
	assert.Equal(t, 23.0, *config.BackupHour.Maximum)
	assert.Equal(t, 0.0, *config.BackupHour.Minimum)
	assert.False(t, config.BackupHour.RequiresRestart)
	assert.Equal(t, linodego.ConfigParamType{"integer", "null"}, config.BackupHour.Type)

	assert.Equal(t, "Default ACL channels policy.", config.ValkeyACLChannelsDefault.Description)
	assert.Equal(t, "resetchannels", config.ValkeyACLChannelsDefault.Example)
	assert.Equal(t, []string{"allchannels", "resetchannels"}, config.ValkeyACLChannelsDefault.Enum)
	assert.False(t, config.ValkeyACLChannelsDefault.RequiresRestart)
	assert.Equal(t, linodego.ConfigParamType{"string"}, config.ValkeyACLChannelsDefault.Type)

	assert.Equal(t, "Whether to take frequent snapshots.", config.FrequentSnapshots.Description)
	assert.Equal(t, true, config.FrequentSnapshots.Example)
	assert.False(t, config.FrequentSnapshots.RequiresRestart)
	assert.Equal(t, linodego.ConfigParamType{"boolean"}, config.FrequentSnapshots.Type)

	assert.Equal(t, "Maxmemory eviction policy.", config.ValkeyMaxmemoryPolicy.Description)
	assert.Equal(t, "noeviction", config.ValkeyMaxmemoryPolicy.Default)
	assert.Equal(t, "noeviction", config.ValkeyMaxmemoryPolicy.Example)
	assert.False(t, config.ValkeyMaxmemoryPolicy.RequiresRestart)
	assert.Equal(t, linodego.ConfigParamType{"string"}, config.ValkeyMaxmemoryPolicy.Type)
}

func TestMarshalValkeyCreateOptions(t *testing.T) {
	opts := linodego.ValkeyCreateOptions{
		Label:  "example-db",
		Region: "us-east",
		Type:   "g6-standard-2",
		Engine: "valkey/8",
	}

	data, err := json.Marshal(opts)
	assert.NoError(t, err)

	var m map[string]interface{}
	assert.NoError(t, json.Unmarshal(data, &m))

	assert.Equal(t, "example-db", m["label"])
	assert.Equal(t, "us-east", m["region"])
	assert.Equal(t, "g6-standard-2", m["type"])
	assert.Equal(t, "valkey/8", m["engine"])

	// Optional, unset fields should be omitted entirely
	assert.NotContains(t, m, "allow_list")
	assert.NotContains(t, m, "cluster_size")
	assert.NotContains(t, m, "fork")
	assert.NotContains(t, m, "engine_config")
	assert.NotContains(t, m, "private_network")
}

func TestMarshalValkeyUpdateOptions_ClearsPrivateNetwork(t *testing.T) {
	var nilNetwork *linodego.DatabasePrivateNetwork
	opts := linodego.ValkeyUpdateOptions{
		Label:          "renamed-db",
		PrivateNetwork: &nilNetwork,
	}

	data, err := json.Marshal(opts)
	assert.NoError(t, err)

	var m map[string]interface{}
	assert.NoError(t, json.Unmarshal(data, &m))

	assert.Equal(t, "renamed-db", m["label"])
	assert.Contains(t, m, "private_network")
	assert.Nil(t, m["private_network"])

	// Unset optional fields should still be omitted
	assert.NotContains(t, m, "region")
	assert.NotContains(t, m, "type")
	assert.NotContains(t, m, "version")
	assert.NotContains(t, m, "cluster_size")
	assert.NotContains(t, m, "updates")
	assert.NotContains(t, m, "engine_config")
	assert.NotContains(t, m, "allow_list")
}
