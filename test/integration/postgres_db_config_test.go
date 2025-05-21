package integration

import (
	"context"
	"os"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestDatabasePostgres_EngineConfig_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestDatabasePostgres_EngineConfig_Get")
	defer teardown()

	config, err := client.GetPostgresDatabaseConfig(context.Background())
	if err != nil {
		t.Fatalf("Error getting PostgreSQL database config: %v", err)
	}

	// Assert that the config is not nil
	assert.NotNil(t, config, "PostgreSQL config should not be nil")

	assert.IsType(t, string(""), config.PG.AutovacuumAnalyzeScaleFactor.Description)
	assert.IsType(t, float64(1.0), config.PG.AutovacuumAnalyzeScaleFactor.Maximum)
	assert.IsType(t, float64(0.0), config.PG.AutovacuumAnalyzeScaleFactor.Minimum)
	assert.IsType(t, bool(false), config.PG.AutovacuumAnalyzeScaleFactor.RequiresRestart)
	assert.IsType(t, string(""), config.PG.AutovacuumAnalyzeScaleFactor.Type)

	assert.IsType(t, string(""), config.PG.AutovacuumAnalyzeThreshold.Description)
	assert.IsType(t, int32(2147483647), config.PG.AutovacuumAnalyzeThreshold.Maximum)
	assert.IsType(t, int32(0), config.PG.AutovacuumAnalyzeThreshold.Minimum)
	assert.IsType(t, bool(false), config.PG.AutovacuumAnalyzeThreshold.RequiresRestart)
	assert.IsType(t, string(""), config.PG.AutovacuumAnalyzeThreshold.Type)

	assert.IsType(t, string(""), config.PG.AutovacuumMaxWorkers.Description)
	assert.IsType(t, int(20), config.PG.AutovacuumMaxWorkers.Maximum)
	assert.IsType(t, int(1), config.PG.AutovacuumMaxWorkers.Minimum)
	assert.IsType(t, bool(false), config.PG.AutovacuumMaxWorkers.RequiresRestart)
	assert.IsType(t, string(""), config.PG.AutovacuumMaxWorkers.Type)

	assert.IsType(t, string(""), config.PG.AutovacuumNaptime.Description)
	assert.IsType(t, int(86400), config.PG.AutovacuumNaptime.Maximum)
	assert.IsType(t, int(1), config.PG.AutovacuumNaptime.Minimum)
	assert.IsType(t, bool(false), config.PG.AutovacuumNaptime.RequiresRestart)
	assert.IsType(t, string(""), config.PG.AutovacuumNaptime.Type)

	assert.IsType(t, string(""), config.PG.AutovacuumVacuumCostDelay.Description)
	assert.IsType(t, int(100), config.PG.AutovacuumVacuumCostDelay.Maximum)
	assert.IsType(t, int(-1), config.PG.AutovacuumVacuumCostDelay.Minimum)
	assert.IsType(t, bool(false), config.PG.AutovacuumVacuumCostDelay.RequiresRestart)
	assert.IsType(t, string(""), config.PG.AutovacuumVacuumCostDelay.Type)

	assert.IsType(t, string(""), config.PG.AutovacuumVacuumCostLimit.Description)
	assert.IsType(t, int(10000), config.PG.AutovacuumVacuumCostLimit.Maximum)
	assert.IsType(t, int(-1), config.PG.AutovacuumVacuumCostLimit.Minimum)
	assert.IsType(t, bool(false), config.PG.AutovacuumVacuumCostLimit.RequiresRestart)
	assert.IsType(t, string(""), config.PG.AutovacuumVacuumCostLimit.Type)

	assert.IsType(t, string(""), config.PG.AutovacuumVacuumScaleFactor.Description)
	assert.IsType(t, float64(1.0), config.PG.AutovacuumVacuumScaleFactor.Maximum)
	assert.IsType(t, float64(0.0), config.PG.AutovacuumVacuumScaleFactor.Minimum)
	assert.IsType(t, bool(false), config.PG.AutovacuumVacuumScaleFactor.RequiresRestart)
	assert.IsType(t, string(""), config.PG.AutovacuumVacuumScaleFactor.Type)

	assert.IsType(t, string(""), config.PG.AutovacuumVacuumThreshold.Description)
	assert.IsType(t, int32(2147483647), config.PG.AutovacuumVacuumThreshold.Maximum)
	assert.IsType(t, int32(0), config.PG.AutovacuumVacuumThreshold.Minimum)
	assert.IsType(t, bool(false), config.PG.AutovacuumVacuumThreshold.RequiresRestart)
	assert.IsType(t, string(""), config.PG.AutovacuumVacuumThreshold.Type)

	assert.IsType(t, string(""), config.PG.BGWriterDelay.Description)
	assert.IsType(t, int(200), config.PG.BGWriterDelay.Example)
	assert.IsType(t, int(10000), config.PG.BGWriterDelay.Maximum)
	assert.IsType(t, int(10), config.PG.BGWriterDelay.Minimum)
	assert.IsType(t, bool(false), config.PG.BGWriterDelay.RequiresRestart)
	assert.IsType(t, string(""), config.PG.BGWriterDelay.Type)

	assert.IsType(t, string(""), config.PG.BGWriterFlushAfter.Description)
	assert.IsType(t, int(512), config.PG.BGWriterFlushAfter.Example)
	assert.IsType(t, int(2048), config.PG.BGWriterFlushAfter.Maximum)
	assert.IsType(t, int(0), config.PG.BGWriterFlushAfter.Minimum)
	assert.IsType(t, bool(false), config.PG.BGWriterFlushAfter.RequiresRestart)
	assert.IsType(t, string(""), config.PG.BGWriterFlushAfter.Type)

	assert.IsType(t, string(""), config.PG.BGWriterLRUMaxPages.Description)
	assert.IsType(t, int(100), config.PG.BGWriterLRUMaxPages.Example)
	assert.IsType(t, int(1073741823), config.PG.BGWriterLRUMaxPages.Maximum)
	assert.IsType(t, int(0), config.PG.BGWriterLRUMaxPages.Minimum)
	assert.IsType(t, bool(false), config.PG.BGWriterLRUMaxPages.RequiresRestart)
	assert.IsType(t, string(""), config.PG.BGWriterLRUMaxPages.Type)

	assert.IsType(t, string(""), config.PG.BGWriterLRUMultiplier.Description)
	assert.IsType(t, float64(2.0), config.PG.BGWriterLRUMultiplier.Example)
	assert.IsType(t, float64(10.0), config.PG.BGWriterLRUMultiplier.Maximum)
	assert.IsType(t, float64(0.0), config.PG.BGWriterLRUMultiplier.Minimum)
	assert.IsType(t, bool(false), config.PG.BGWriterLRUMultiplier.RequiresRestart)
	assert.IsType(t, string(""), config.PG.BGWriterLRUMultiplier.Type)

	assert.IsType(t, string(""), config.PG.DeadlockTimeout.Description)
	assert.IsType(t, int(1000), config.PG.DeadlockTimeout.Example)
	assert.IsType(t, int(1800000), config.PG.DeadlockTimeout.Maximum)
	assert.IsType(t, int(500), config.PG.DeadlockTimeout.Minimum)
	assert.IsType(t, bool(false), config.PG.DeadlockTimeout.RequiresRestart)
	assert.IsType(t, string(""), config.PG.DeadlockTimeout.Type)

	assert.IsType(t, string(""), config.PG.DefaultToastCompression.Description)
	assert.IsType(t, []string{"lz4", "pglz"}, config.PG.DefaultToastCompression.Enum)
	assert.IsType(t, string("lz4"), config.PG.DefaultToastCompression.Example)
	assert.IsType(t, bool(false), config.PG.DefaultToastCompression.RequiresRestart)
	assert.IsType(t, string(""), config.PG.DefaultToastCompression.Type)

	assert.IsType(t, string(""), config.PG.IdleInTransactionSessionTimeout.Description)
	assert.IsType(t, int(0), config.PG.IdleInTransactionSessionTimeout.Maximum)
	assert.IsType(t, int(0), config.PG.IdleInTransactionSessionTimeout.Minimum)
	assert.IsType(t, false, config.PG.IdleInTransactionSessionTimeout.RequiresRestart)
	assert.IsType(t, string(""), config.PG.IdleInTransactionSessionTimeout.Type)

	assert.IsType(t, string(""), config.PG.JIT.Description)
	assert.IsType(t, true, config.PG.JIT.Example)
	assert.IsType(t, false, config.PG.JIT.RequiresRestart)
	assert.IsType(t, string(""), config.PG.JIT.Type)

	assert.IsType(t, string(""), config.PG.MaxFilesPerProcess.Description)
	assert.IsType(t, int(0), config.PG.MaxFilesPerProcess.Maximum)
	assert.IsType(t, int(0), config.PG.MaxFilesPerProcess.Minimum)
	assert.IsType(t, false, config.PG.MaxFilesPerProcess.RequiresRestart)
	assert.IsType(t, string(""), config.PG.MaxFilesPerProcess.Type)

	assert.IsType(t, string(""), config.PG.MaxLocksPerTransaction.Description)
	assert.IsType(t, int(0), config.PG.MaxLocksPerTransaction.Maximum)
	assert.IsType(t, int(0), config.PG.MaxLocksPerTransaction.Minimum)
	assert.IsType(t, false, config.PG.MaxLocksPerTransaction.RequiresRestart)
	assert.IsType(t, string(""), config.PG.MaxLocksPerTransaction.Type)

	assert.IsType(t, string(""), config.PG.MaxLogicalReplicationWorkers.Description)
	assert.IsType(t, int(0), config.PG.MaxLogicalReplicationWorkers.Maximum)
	assert.IsType(t, int(0), config.PG.MaxLogicalReplicationWorkers.Minimum)
	assert.IsType(t, false, config.PG.MaxLogicalReplicationWorkers.RequiresRestart)
	assert.IsType(t, string(""), config.PG.MaxLogicalReplicationWorkers.Type)

	assert.IsType(t, string(""), config.PG.MaxParallelWorkers.Description)
	assert.IsType(t, int(0), config.PG.MaxParallelWorkers.Maximum)
	assert.IsType(t, int(0), config.PG.MaxParallelWorkers.Minimum)
	assert.IsType(t, false, config.PG.MaxParallelWorkers.RequiresRestart)
	assert.IsType(t, string(""), config.PG.MaxParallelWorkers.Type)

	assert.IsType(t, string(""), config.PG.MaxParallelWorkersPerGather.Description)
	assert.IsType(t, int(0), config.PG.MaxParallelWorkersPerGather.Maximum)
	assert.IsType(t, int(0), config.PG.MaxParallelWorkersPerGather.Minimum)
	assert.IsType(t, false, config.PG.MaxParallelWorkersPerGather.RequiresRestart)
	assert.IsType(t, string(""), config.PG.MaxParallelWorkersPerGather.Type)

	assert.IsType(t, string(""), config.PG.MaxPredLocksPerTransaction.Description)
	assert.IsType(t, int(0), config.PG.MaxPredLocksPerTransaction.Maximum)
	assert.IsType(t, int(0), config.PG.MaxPredLocksPerTransaction.Minimum)
	assert.IsType(t, false, config.PG.MaxPredLocksPerTransaction.RequiresRestart)
	assert.IsType(t, string(""), config.PG.MaxPredLocksPerTransaction.Type)

	assert.IsType(t, string(""), config.PG.MaxReplicationSlots.Description)
	assert.IsType(t, int(0), config.PG.MaxReplicationSlots.Maximum)
	assert.IsType(t, int(0), config.PG.MaxReplicationSlots.Minimum)
	assert.IsType(t, false, config.PG.MaxReplicationSlots.RequiresRestart)
	assert.IsType(t, string(""), config.PG.MaxReplicationSlots.Type)

	assert.IsType(t, string(""), config.PG.MaxSlotWALKeepSize.Description)
	assert.IsType(t, int32(0), config.PG.MaxSlotWALKeepSize.Maximum)
	assert.IsType(t, int32(0), config.PG.MaxSlotWALKeepSize.Minimum)
	assert.IsType(t, false, config.PG.MaxSlotWALKeepSize.RequiresRestart)
	assert.IsType(t, string(""), config.PG.MaxSlotWALKeepSize.Type)

	assert.IsType(t, string(""), config.PG.MaxStackDepth.Description)
	assert.IsType(t, int(0), config.PG.MaxStackDepth.Maximum)
	assert.IsType(t, int(0), config.PG.MaxStackDepth.Minimum)
	assert.IsType(t, false, config.PG.MaxStackDepth.RequiresRestart)
	assert.IsType(t, string(""), config.PG.MaxStackDepth.Type)

	assert.IsType(t, string(""), config.PG.MaxStandbyArchiveDelay.Description)
	assert.IsType(t, int(0), config.PG.MaxStandbyArchiveDelay.Maximum)
	assert.IsType(t, int(0), config.PG.MaxStandbyArchiveDelay.Minimum)
	assert.IsType(t, false, config.PG.MaxStandbyArchiveDelay.RequiresRestart)
	assert.IsType(t, string(""), config.PG.MaxStandbyArchiveDelay.Type)

	assert.IsType(t, "Max standby streaming delay in milliseconds", config.PG.MaxStandbyStreamingDelay.Description)
	assert.IsType(t, int(43200000), config.PG.MaxStandbyStreamingDelay.Maximum)
	assert.IsType(t, int(1), config.PG.MaxStandbyStreamingDelay.Minimum)
	assert.IsType(t, false, config.PG.MaxStandbyStreamingDelay.RequiresRestart)
	assert.IsType(t, "integer", config.PG.MaxStandbyStreamingDelay.Type)

	assert.IsType(t, "PostgreSQL maximum WAL senders", config.PG.MaxWALSenders.Description)
	assert.IsType(t, int(64), config.PG.MaxWALSenders.Maximum)
	assert.IsType(t, int(20), config.PG.MaxWALSenders.Minimum)
	assert.IsType(t, false, config.PG.MaxWALSenders.RequiresRestart)
	assert.IsType(t, "integer", config.PG.MaxWALSenders.Type)

	assert.IsType(t, "Sets the maximum number of background processes that the system can support", config.PG.MaxWorkerProcesses.Description)
	assert.IsType(t, int(96), config.PG.MaxWorkerProcesses.Maximum)
	assert.IsType(t, int(8), config.PG.MaxWorkerProcesses.Minimum)
	assert.IsType(t, false, config.PG.MaxWorkerProcesses.RequiresRestart)
	assert.IsType(t, "integer", config.PG.MaxWorkerProcesses.Type)

	assert.IsType(t, "Chooses the algorithm for encrypting passwords.", config.PG.PasswordEncryption.Description)
	assert.IsType(t, []string{"md5", "scram-sha-256"}, config.PG.PasswordEncryption.Enum)
	assert.IsType(t, "scram-sha-256", config.PG.PasswordEncryption.Example)
	assert.IsType(t, false, config.PG.PasswordEncryption.RequiresRestart)
	assert.IsType(t, "string", config.PG.PasswordEncryption.Type)

	assert.IsType(t, "Sets the time interval to run pg_partman's scheduled tasks", config.PG.PGPartmanBGWInterval.Description)
	assert.IsType(t, int(3600), config.PG.PGPartmanBGWInterval.Example)
	assert.IsType(t, int(604800), config.PG.PGPartmanBGWInterval.Maximum)
	assert.IsType(t, int(3600), config.PG.PGPartmanBGWInterval.Minimum)
	assert.IsType(t, false, config.PG.PGPartmanBGWInterval.RequiresRestart)
	assert.IsType(t, "integer", config.PG.PGPartmanBGWInterval.Type)

	assert.IsType(t, "Controls which role to use for pg_partman's scheduled background tasks.", config.PG.PGPartmanBGWRole.Description)
	assert.IsType(t, "myrolename", config.PG.PGPartmanBGWRole.Example)
	assert.IsType(t, int(64), config.PG.PGPartmanBGWRole.MaxLength)
	assert.IsType(t, "^[_A-Za-z0-9][-._A-Za-z0-9]{0,63}$", config.PG.PGPartmanBGWRole.Pattern)
	assert.IsType(t, false, config.PG.PGPartmanBGWRole.RequiresRestart)
	assert.IsType(t, "string", config.PG.PGPartmanBGWRole.Type)

	assert.IsType(t, "Enables or disables query plan monitoring", config.PG.PGStatMonitorPGSMEnableQueryPlan.Description)
	assert.IsType(t, false, config.PG.PGStatMonitorPGSMEnableQueryPlan.Example)
	assert.IsType(t, false, config.PG.PGStatMonitorPGSMEnableQueryPlan.RequiresRestart)
	assert.IsType(t, "boolean", config.PG.PGStatMonitorPGSMEnableQueryPlan.Type)

	assert.IsType(t, "Sets the maximum number of buckets", config.PG.PGStatMonitorPGSMMaxBuckets.Description)
	assert.IsType(t, int(10), config.PG.PGStatMonitorPGSMMaxBuckets.Example)
	assert.IsType(t, int(10), config.PG.PGStatMonitorPGSMMaxBuckets.Maximum)
	assert.IsType(t, int(1), config.PG.PGStatMonitorPGSMMaxBuckets.Minimum)
	assert.IsType(t, false, config.PG.PGStatMonitorPGSMMaxBuckets.RequiresRestart)
	assert.IsType(t, "integer", config.PG.PGStatMonitorPGSMMaxBuckets.Type)

	assert.IsType(t, "Controls which statements are counted. Specify top to track top-level statements (those issued directly by clients), all to also track nested statements (such as statements invoked within functions), or none to disable statement statistics collection. The default value is top.", config.PG.PGStatStatementsTrack.Description)
	assert.IsType(t, []string{"all", "top", "none"}, config.PG.PGStatStatementsTrack.Enum)
	assert.IsType(t, false, config.PG.PGStatStatementsTrack.RequiresRestart)
	assert.IsType(t, "string", config.PG.PGStatStatementsTrack.Type)

	assert.IsType(t, "PostgreSQL temporary file limit in KiB, -1 for unlimited", config.PG.TempFileLimit.Description)
	assert.IsType(t, int32(5000000), config.PG.TempFileLimit.Example)
	assert.IsType(t, int32(2147483647), config.PG.TempFileLimit.Maximum)
	assert.IsType(t, int32(-1), config.PG.TempFileLimit.Minimum)
	assert.IsType(t, false, config.PG.TempFileLimit.RequiresRestart)
	assert.IsType(t, "integer", config.PG.TempFileLimit.Type)

	assert.IsType(t, "PostgreSQL service timezone", config.PG.Timezone.Description)
	assert.IsType(t, "Europe/Helsinki", config.PG.Timezone.Example)
	assert.IsType(t, int(64), config.PG.Timezone.MaxLength)
	assert.IsType(t, "^[\\w/]*$", config.PG.Timezone.Pattern)
	assert.IsType(t, false, config.PG.Timezone.RequiresRestart)
	assert.IsType(t, "string", config.PG.Timezone.Type)

	assert.IsType(t, "Specifies the number of bytes reserved to track the currently executing command for each active session.", config.PG.TrackActivityQuerySize.Description)
	assert.IsType(t, int(1024), config.PG.TrackActivityQuerySize.Example)
	assert.IsType(t, int(10240), config.PG.TrackActivityQuerySize.Maximum)
	assert.IsType(t, int(1024), config.PG.TrackActivityQuerySize.Minimum)
	assert.IsType(t, false, config.PG.TrackActivityQuerySize.RequiresRestart)
	assert.IsType(t, "integer", config.PG.TrackActivityQuerySize.Type)

	assert.IsType(t, "Record commit time of transactions.", config.PG.TrackCommitTimestamp.Description)
	assert.IsType(t, "off", config.PG.TrackCommitTimestamp.Example)
	assert.IsType(t, []string{"off", "on"}, config.PG.TrackCommitTimestamp.Enum)
	assert.IsType(t, false, config.PG.TrackCommitTimestamp.RequiresRestart)
	assert.IsType(t, "string", config.PG.TrackCommitTimestamp.Type)
}

func TestDatabasePostgres_EngineConfig_Suite(t *testing.T) {
	databaseModifiers := []postgresDatabaseModifier{
		createPostgresOptionsModifier(),
	}

	client, database, teardown, err := setupPostgresDatabase(t, databaseModifiers, "fixtures/TestDatabasePostgres_EngineConfig_Suite")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	assertPostgresDatabaseBasics(t, database)

	expected := newExpectedPostgresEngineConfig()

	assertPostgresEngineConfigEqual(t, database.EngineConfig.PG, expected)

	fetchedDB, err := client.GetPostgresDatabase(context.Background(), database.ID)
	assert.NoError(t, err)
	assertPostgresEngineConfigEqual(t, fetchedDB.EngineConfig.PG, expected)

	updateOptions := linodego.PostgresUpdateOptions{
		Label: "example-db-updated",
		EngineConfig: &linodego.PostgresDatabaseEngineConfig{
			PG: &linodego.PostgresDatabaseEngineConfigPG{
				AutovacuumVacuumThreshold: linodego.Pointer(int32(500)),
				DeadlockTimeout:           linodego.Pointer(3000),
			},
		},
	}

	updatedDB, err := client.UpdatePostgresDatabase(context.Background(), database.ID, updateOptions)
	if err != nil {
		t.Errorf("failed to update db %d: %v", database.ID, err)
	}

	waitForDatabaseUpdated(t, client, updatedDB.ID, linodego.DatabaseEngineTypePostgres, updatedDB.Created)

	assertUpdatedPostgresFields(t, updatedDB.EngineConfig.PG)
}

func TestDatabasePostgres_EngineConfig_Create_PasswordEncryption_DefaultsToMD5(t *testing.T) {
	databaseModifiers := []postgresDatabaseModifier{
		createPostgresOptionsModifierWithNullField(),
	}
	client, fixtureTeardown := createTestClient(t, "fixtures/TestDatabasePostgres_EngineConfig_Create_PasswordEncryption_DefaultsToMD5")

	database, databaseTeardown, err := createPostgresDatabase(t, client, databaseModifiers)
	if err != nil {
		t.Fatalf("failed to create db: %s", err)
	}

	defer func() {
		databaseTeardown()
		fixtureTeardown()
	}()

	// Password Encryption Value will default to md5 if initial input is null
	assert.Contains(t, *database.EngineConfig.PG.PasswordEncryption, "md5")
}

func TestDatabasePostgres_EngineConfig_Create_Fails_EmptyDoublePointerValue(t *testing.T) {
	if os.Getenv("LINODE_FIXTURE_MODE") == "play" {
		t.Skip("Skipping negative test scenario: LINODE_FIXTURE_MODE is 'play'")
	}

	invalidRequestData := linodego.PostgresCreateOptions{
		Label:  "example-db-created-fails",
		Region: "us-east",
		Type:   "g6-dedicated-2",
		Engine: "postgresql/14",
		EngineConfig: &linodego.PostgresDatabaseEngineConfig{
			PG: &linodego.PostgresDatabaseEngineConfigPG{
				PasswordEncryption: linodego.Pointer(""),
			},
		},
	}

	client, _ := createTestClient(t, "")

	_, err := client.CreatePostgresDatabase(context.Background(), invalidRequestData)

	assert.Contains(t, err.Error(), "Invalid value: expected one of ['md5', 'scram-sha-256']")
}

func createPostgresOptionsModifier() postgresDatabaseModifier {
	return func(options *linodego.PostgresCreateOptions) {
		options.Label = "postgres-db-created-with-config"
		options.Region = "us-east"
		options.Type = "g6-dedicated-2"
		options.Engine = "postgresql/17"
		options.EngineConfig = &linodego.PostgresDatabaseEngineConfig{
			PG: &linodego.PostgresDatabaseEngineConfigPG{
				AutovacuumAnalyzeScaleFactor:     linodego.Pointer(0.1),
				AutovacuumAnalyzeThreshold:       linodego.Pointer(int32(500)),
				AutovacuumMaxWorkers:             linodego.Pointer(3),
				AutovacuumNaptime:                linodego.Pointer(60),
				AutovacuumVacuumCostDelay:        linodego.Pointer(100), // Reduced to <= 100
				AutovacuumVacuumCostLimit:        linodego.Pointer(2000),
				AutovacuumVacuumScaleFactor:      linodego.Pointer(0.2),
				AutovacuumVacuumThreshold:        linodego.Pointer(int32(500)),
				BGWriterDelay:                    linodego.Pointer(100),
				BGWriterFlushAfter:               linodego.Pointer(1000),
				BGWriterLRUMaxPages:              linodego.Pointer(100),
				BGWriterLRUMultiplier:            linodego.Pointer(2.0),
				DeadlockTimeout:                  linodego.Pointer(1000),
				DefaultToastCompression:          linodego.Pointer("lz4"),
				IdleInTransactionSessionTimeout:  linodego.Pointer(600),
				JIT:                              linodego.Pointer(true),
				MaxFilesPerProcess:               linodego.Pointer(1000),
				MaxLocksPerTransaction:           linodego.Pointer(64),
				MaxLogicalReplicationWorkers:     linodego.Pointer(4),
				MaxParallelWorkers:               linodego.Pointer(8),
				MaxParallelWorkersPerGather:      linodego.Pointer(2),
				MaxPredLocksPerTransaction:       linodego.Pointer(64),
				MaxReplicationSlots:              linodego.Pointer(8), // Adjusted to >= 8
				MaxSlotWALKeepSize:               linodego.Pointer(int32(512)),
				MaxStackDepth:                    linodego.Pointer(2097152), // Adjusted to >= 2MB
				MaxStandbyArchiveDelay:           linodego.Pointer(30000),
				MaxStandbyStreamingDelay:         linodego.Pointer(30000),
				MaxWALSenders:                    linodego.Pointer(20), // Adjusted to >= 20
				MaxWorkerProcesses:               linodego.Pointer(8),
				PasswordEncryption:               linodego.Pointer("scram-sha-256"),
				PGPartmanBGWInterval:             linodego.Pointer(3600),
				PGPartmanBGWRole:                 linodego.Pointer("pg_partman_bgw"),
				PGStatMonitorPGSMEnableQueryPlan: linodego.Pointer(true),
				PGStatMonitorPGSMMaxBuckets:      linodego.Pointer(10), // Adjusted to <= 10
				PGStatStatementsTrack:            linodego.Pointer("top"),
				TempFileLimit:                    linodego.Pointer(int32(1000)),
				Timezone:                         linodego.Pointer("UTC"),
				TrackActivityQuerySize:           linodego.Pointer(1024),
				TrackCommitTimestamp:             linodego.Pointer("on"),
				TrackFunctions:                   linodego.Pointer("all"), // Adjusted to valid value
				TrackIOTiming:                    linodego.Pointer("on"),
				WALSenderTimeout:                 linodego.Pointer(60000),
				WALWriterDelay:                   linodego.Pointer(200), // Adjusted to <= 200
			},
			PGStatMonitorEnable:     linodego.Pointer(true),
			PGLookout:               &linodego.PostgresDatabaseEngineConfigPGLookout{},
			SharedBuffersPercentage: linodego.Pointer(25.0),
			WorkMem:                 linodego.Pointer(1024), // Adjusted to <= 1024
		}
	}
}

func createPostgresOptionsModifierWithNullField() postgresDatabaseModifier {
	return func(options *linodego.PostgresCreateOptions) {
		options.Label = "postgres-db-created-with-config"
		options.Region = "us-east"
		options.Type = "g6-dedicated-2"
		options.Engine = "postgresql/17"
		options.EngineConfig = &linodego.PostgresDatabaseEngineConfig{
			PG: &linodego.PostgresDatabaseEngineConfigPG{
				PasswordEncryption: nil,
			},
		}
	}
}

func newExpectedPostgresEngineConfig() map[string]any {
	return map[string]any{
		"AutovacuumAnalyzeScaleFactor":     0.1,
		"AutovacuumAnalyzeThreshold":       int32(500),
		"AutovacuumMaxWorkers":             3,
		"AutovacuumNaptime":                60,
		"AutovacuumVacuumCostDelay":        100,
		"AutovacuumVacuumCostLimit":        2000,
		"AutovacuumVacuumScaleFactor":      0.2,
		"AutovacuumVacuumThreshold":        int32(500),
		"BGWriterDelay":                    100,
		"BGWriterFlushAfter":               1000,
		"BGWriterLRUMaxPages":              100,
		"BGWriterLRUMultiplier":            2.0,
		"DeadlockTimeout":                  1000,
		"DefaultToastCompression":          "lz4",
		"IdleInTransactionSessionTimeout":  600,
		"JIT":                              true,
		"MaxFilesPerProcess":               1000,
		"MaxLocksPerTransaction":           64,
		"MaxLogicalReplicationWorkers":     4,
		"MaxParallelWorkers":               8,
		"MaxParallelWorkersPerGather":      2,
		"MaxPredLocksPerTransaction":       64,
		"MaxReplicationSlots":              8,
		"MaxSlotWALKeepSize":               int32(512),
		"MaxStackDepth":                    2097152,
		"MaxStandbyArchiveDelay":           30000,
		"MaxStandbyStreamingDelay":         30000,
		"MaxWALSenders":                    20,
		"MaxWorkerProcesses":               8,
		"PasswordEncryption":               "scram-sha-256",
		"PGPartmanBGWInterval":             3600,
		"PGPartmanBGWRole":                 "pg_partman_bgw",
		"PGStatMonitorPGSMEnableQueryPlan": true,
		"PGStatMonitorPGSMMaxBuckets":      10,
		"PGStatStatementsTrack":            "top",
		"TempFileLimit":                    int32(1000),
		"Timezone":                         "UTC",
		"TrackActivityQuerySize":           1024,
		"TrackCommitTimestamp":             "on",
		"TrackFunctions":                   "all",
		"TrackIOTiming":                    "on",
		"WALSenderTimeout":                 60000,
		"WALWriterDelay":                   200,
		"PGStatMonitorEnable":              true,
		"PGLookout":                        map[string]any{},
		"SharedBuffersPercentage":          25.0,
		"WorkMem":                          1024,
	}
}

func assertPostgresDatabaseBasics(t *testing.T, db *linodego.PostgresDatabase) {
	// Assert basic fields with types and values
	assert.IsType(t, int(0), db.ID) // Assuming ID is an int32, use the type in your assert accordingly
	assert.Equal(t, "postgres-db-created-with-config", db.Label)
	assert.Equal(t, "us-east", db.Region)
	assert.Equal(t, "g6-dedicated-2", db.Type)
	assert.Equal(t, "postgresql", db.Engine)
	assert.Equal(t, "17", db.Version)
	assert.NotEmpty(t, db.AllowList)
	assert.IsType(t, int(25698), db.Port)
	assert.IsType(t, int(3), db.ClusterSize)
	assert.Equal(t, linodego.DatabasePlatform("rdbms-default"), db.Platform)

	// Assert Hosts structure
	assert.NotEmpty(t, db.Hosts.Primary)
	assert.Empty(t, db.Hosts.Secondary)

	// Assert EngineConfig structure
	assert.NotNil(t, db.EngineConfig.PG)
	assert.True(t, *db.EngineConfig.PGStatMonitorEnable)
	assert.NotNil(t, db.EngineConfig.PGLookout)
	assert.NotNil(t, db.EngineConfig.SharedBuffersPercentage)
	assert.NotNil(t, db.EngineConfig.WorkMem)
}

func assertPostgresEngineConfigEqual(t *testing.T, cfg *linodego.PostgresDatabaseEngineConfigPG, expected map[string]any) {
	// Compare all fields of the struct using the expected map

	// Autovacuum
	assert.Equal(t, expected["AutovacuumAnalyzeScaleFactor"], *cfg.AutovacuumAnalyzeScaleFactor)
	assert.Equal(t, expected["AutovacuumAnalyzeThreshold"], *cfg.AutovacuumAnalyzeThreshold)
	assert.Equal(t, expected["AutovacuumMaxWorkers"], *cfg.AutovacuumMaxWorkers)
	assert.Equal(t, expected["AutovacuumNaptime"], *cfg.AutovacuumNaptime)
	assert.Equal(t, expected["AutovacuumVacuumCostDelay"], *cfg.AutovacuumVacuumCostDelay)
	assert.Equal(t, expected["AutovacuumVacuumCostLimit"], *cfg.AutovacuumVacuumCostLimit)
	assert.Equal(t, expected["AutovacuumVacuumScaleFactor"], *cfg.AutovacuumVacuumScaleFactor)
	assert.Equal(t, expected["AutovacuumVacuumThreshold"], *cfg.AutovacuumVacuumThreshold)

	// BGWriter
	assert.Equal(t, expected["BGWriterDelay"], *cfg.BGWriterDelay)
	assert.Equal(t, expected["BGWriterFlushAfter"], *cfg.BGWriterFlushAfter)
	assert.Equal(t, expected["BGWriterLRUMaxPages"], *cfg.BGWriterLRUMaxPages)
	assert.Equal(t, expected["BGWriterLRUMultiplier"], *cfg.BGWriterLRUMultiplier)

	// DeadlockTimeout
	assert.Equal(t, expected["DeadlockTimeout"], *cfg.DeadlockTimeout)

	// DefaultToastCompression
	assert.Equal(t, expected["DefaultToastCompression"], *cfg.DefaultToastCompression)

	// IdleInTransactionSessionTimeout
	assert.Equal(t, expected["IdleInTransactionSessionTimeout"], *cfg.IdleInTransactionSessionTimeout)

	// JIT
	assert.Equal(t, expected["JIT"], *cfg.JIT)

	// Max files and locks
	assert.Equal(t, expected["MaxFilesPerProcess"], *cfg.MaxFilesPerProcess)
	assert.Equal(t, expected["MaxLocksPerTransaction"], *cfg.MaxLocksPerTransaction)
	assert.Equal(t, expected["MaxLogicalReplicationWorkers"], *cfg.MaxLogicalReplicationWorkers)
	assert.Equal(t, expected["MaxParallelWorkers"], *cfg.MaxParallelWorkers)
	assert.Equal(t, expected["MaxParallelWorkersPerGather"], *cfg.MaxParallelWorkersPerGather)
	assert.Equal(t, expected["MaxPredLocksPerTransaction"], *cfg.MaxPredLocksPerTransaction)

	// MaxReplicationSlots and MaxSlotWALKeepSize
	assert.Equal(t, expected["MaxReplicationSlots"], *cfg.MaxReplicationSlots)
	assert.Equal(t, expected["MaxSlotWALKeepSize"], *cfg.MaxSlotWALKeepSize)

	// MaxStandby
	assert.Equal(t, expected["MaxStandbyArchiveDelay"], *cfg.MaxStandbyArchiveDelay)
	assert.Equal(t, expected["MaxStandbyStreamingDelay"], *cfg.MaxStandbyStreamingDelay)

	// MaxWALSenders and MaxWorkerProcesses
	assert.Equal(t, expected["MaxWALSenders"], *cfg.MaxWALSenders)
	assert.Equal(t, expected["MaxWorkerProcesses"], *cfg.MaxWorkerProcesses)

	// PasswordEncryption
	assert.Equal(t, expected["PasswordEncryption"], *cfg.PasswordEncryption)

	// PGPartman settings
	assert.Equal(t, expected["PGPartmanBGWInterval"], *cfg.PGPartmanBGWInterval)
	assert.Equal(t, expected["PGPartmanBGWRole"], *cfg.PGPartmanBGWRole)

	// PGStatMonitor
	assert.Equal(t, expected["PGStatMonitorPGSMEnableQueryPlan"], *cfg.PGStatMonitorPGSMEnableQueryPlan)
	assert.Equal(t, expected["PGStatMonitorPGSMMaxBuckets"], *cfg.PGStatMonitorPGSMMaxBuckets)
	assert.Equal(t, expected["PGStatStatementsTrack"], *cfg.PGStatStatementsTrack)

	// TempFileLimit and Timezone
	assert.Equal(t, expected["TempFileLimit"], *cfg.TempFileLimit)
	assert.Equal(t, expected["Timezone"], *cfg.Timezone)

	// TrackActivityQuerySize and TrackCommitTimestamp
	assert.Equal(t, expected["TrackActivityQuerySize"], *cfg.TrackActivityQuerySize)
	assert.Equal(t, expected["TrackCommitTimestamp"], *cfg.TrackCommitTimestamp)

	// TrackFunctions and TrackIOTiming
	assert.Equal(t, expected["TrackFunctions"], *cfg.TrackFunctions)
	assert.Equal(t, expected["TrackIOTiming"], *cfg.TrackIOTiming)

	// WALSenderTimeout
	assert.Equal(t, expected["WALSenderTimeout"], *cfg.WALSenderTimeout)

	// WALWriterDelay adjusted to <= 200
	assert.Equal(t, expected["WALWriterDelay"], *cfg.WALWriterDelay)
}

func assertUpdatedPostgresFields(t *testing.T, updatedConfig *linodego.PostgresDatabaseEngineConfigPG) {
	assert.Equal(t, int32(500), *updatedConfig.AutovacuumVacuumThreshold)
	assert.Equal(t, int(3000), *updatedConfig.DeadlockTimeout)
}
