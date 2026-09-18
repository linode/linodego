package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseValkey_EngineConfig_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestDatabaseValkey_EngineConfig_Get")
	defer teardown()

	config, err := client.GetValkeyDatabaseConfig(context.Background())
	if err != nil {
		t.Fatalf("error getting Valkey database config: %v", err)
	}

	assert.Equal(t, linodego.ConfigParamType{"integer", "null"}, config.BackupHour.Type)
	assert.Equal(t, linodego.ConfigParamType{"integer", "null"}, config.BackupMinute.Type)
	assert.Equal(t, linodego.ConfigParamType{"boolean"}, config.FrequentSnapshots.Type)
	assert.Equal(t, linodego.ConfigParamType{"string"}, config.ValkeyACLChannelsDefault.Type)
	assert.NotEmpty(t, config.ValkeyACLChannelsDefault.Enum)
	assert.Equal(t, linodego.ConfigParamType{"integer"}, config.ValkeyActiveExpireEffort.Type)
	assert.Equal(t, linodego.ConfigParamType{"boolean"}, config.ValkeyActiveDefrag.Type)
	assert.Equal(t, linodego.ConfigParamType{"integer"}, config.ValkeyLFUDecayTime.Type)
	assert.Equal(t, linodego.ConfigParamType{"integer"}, config.ValkeyLFULogFactor.Type)
	assert.Equal(t, linodego.ConfigParamType{"string", "null"}, config.ValkeyMaxmemoryPolicy.Type)
	assert.Equal(t, linodego.ConfigParamType{"integer"}, config.ValkeyNumberOfDatabases.Type)
	assert.Equal(t, linodego.ConfigParamType{"string"}, config.ValkeyPersistence.Type)
	assert.Equal(t, linodego.ConfigParamType{"integer"}, config.ValkeyPubsubClientOutputBufferLimit.Type)
	assert.Equal(t, linodego.ConfigParamType{"integer"}, config.ValkeyTimeout.Type)
}
