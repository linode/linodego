package unit

import (
	"encoding/json"
	"testing"
	"time"

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
