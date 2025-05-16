package unit

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestListDatabasePostgreSQL_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("postgresql_databases_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/postgresql/instances", fixtureData)
	databases, err := base.Client.ListPostgresDatabases(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, databases, "Expected non-empty postgresql database list")
}

func TestDatabasePostgreSQL_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("postgresql_database_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/postgresql/instances/123", fixtureData)

	db, err := base.Client.GetPostgresDatabase(context.Background(), 123)
	assert.NoError(t, err)

	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, "postgresql", db.Engine)
	assert.Equal(t, 123, db.ID)
	assert.Equal(t, "example-db", db.Label)
	assert.Equal(t, linodego.DatabasePlatform("rdbms-default"), db.Platform)
	assert.Equal(t, 3306, db.Port)
	assert.Equal(t, "us-east", db.Region)
	assert.Equal(t, linodego.DatabaseStatus("active"), db.Status)
	assert.Equal(t, 15, db.TotalDiskSizeGB)
	assert.Equal(t, "g6-dedicated-2", db.Type)
	assert.Equal(t, 3, db.Updates.Duration)
	assert.Equal(t, linodego.DatabaseDayOfWeek(1), db.Updates.DayOfWeek)
	assert.Equal(t, linodego.DatabaseMaintenanceFrequency("weekly"), db.Updates.Frequency)
	assert.Equal(t, 0, db.Updates.HourOfDay)
	assert.Equal(t, 2, db.UsedDiskSizeGB)
	assert.Equal(t, "13.2", db.Version)

	assert.Equal(t, true, *db.EngineConfig.PGStatMonitorEnable)
	assert.Equal(t, int64(1000), *db.EngineConfig.PGLookout.MaxFailoverReplicationTimeLag)
	assert.Equal(t, 41.5, *db.EngineConfig.SharedBuffersPercentage)
	assert.Equal(t, 4, *db.EngineConfig.WorkMem)
	assert.Equal(t, 0.5, *db.EngineConfig.PG.AutovacuumAnalyzeScaleFactor)
	assert.Equal(t, int32(100), *db.EngineConfig.PG.AutovacuumAnalyzeThreshold)
	assert.Equal(t, 10, *db.EngineConfig.PG.AutovacuumMaxWorkers)
	assert.Equal(t, 100, *db.EngineConfig.PG.AutovacuumNaptime)
	assert.Equal(t, 50, *db.EngineConfig.PG.AutovacuumVacuumCostDelay)
	assert.Equal(t, 100, *db.EngineConfig.PG.AutovacuumVacuumCostLimit)
	assert.Equal(t, 0.5, *db.EngineConfig.PG.AutovacuumVacuumScaleFactor)
	assert.Equal(t, int32(100), *db.EngineConfig.PG.AutovacuumVacuumThreshold)
	assert.Equal(t, 200, *db.EngineConfig.PG.BGWriterDelay)
	assert.Equal(t, 512, *db.EngineConfig.PG.BGWriterFlushAfter)
	assert.Equal(t, 100, *db.EngineConfig.PG.BGWriterLRUMaxPages)
	assert.Equal(t, 2.0, *db.EngineConfig.PG.BGWriterLRUMultiplier)
	assert.Equal(t, 1000, *db.EngineConfig.PG.DeadlockTimeout)
	assert.Equal(t, "lz4", *db.EngineConfig.PG.DefaultToastCompression)
	assert.Equal(t, 100, *db.EngineConfig.PG.IdleInTransactionSessionTimeout)
	assert.Equal(t, true, *db.EngineConfig.PG.JIT)
	assert.Equal(t, 100, *db.EngineConfig.PG.MaxFilesPerProcess)
	assert.Equal(t, 100, *db.EngineConfig.PG.MaxLocksPerTransaction)
	assert.Equal(t, 32, *db.EngineConfig.PG.MaxLogicalReplicationWorkers)
	assert.Equal(t, 64, *db.EngineConfig.PG.MaxParallelWorkers)
	assert.Equal(t, 64, *db.EngineConfig.PG.MaxParallelWorkersPerGather)
	assert.Equal(t, 1000, *db.EngineConfig.PG.MaxPredLocksPerTransaction)
	assert.Equal(t, 32, *db.EngineConfig.PG.MaxReplicationSlots)
	assert.Equal(t, int32(100), *db.EngineConfig.PG.MaxSlotWALKeepSize)
	assert.Equal(t, 3507152, *db.EngineConfig.PG.MaxStackDepth)
	assert.Equal(t, 1000, *db.EngineConfig.PG.MaxStandbyArchiveDelay)
	assert.Equal(t, 1000, *db.EngineConfig.PG.MaxStandbyStreamingDelay)
	assert.Equal(t, 32, *db.EngineConfig.PG.MaxWALSenders)
	assert.Equal(t, 64, *db.EngineConfig.PG.MaxWorkerProcesses)
	assert.Equal(t, "scram-sha-256", *db.EngineConfig.PG.PasswordEncryption)
	assert.Equal(t, 3600, *db.EngineConfig.PG.PGPartmanBGWInterval)
	assert.Equal(t, "myrolename", *db.EngineConfig.PG.PGPartmanBGWRole)
	assert.Equal(t, false, *db.EngineConfig.PG.PGStatMonitorPGSMEnableQueryPlan)
	assert.Equal(t, 10, *db.EngineConfig.PG.PGStatMonitorPGSMMaxBuckets)
	assert.Equal(t, "top", *db.EngineConfig.PG.PGStatStatementsTrack)
	assert.Equal(t, int32(5000000), *db.EngineConfig.PG.TempFileLimit)
	assert.Equal(t, "Europe/Helsinki", *db.EngineConfig.PG.Timezone)
	assert.Equal(t, 1024, *db.EngineConfig.PG.TrackActivityQuerySize)
	assert.Equal(t, "off", *db.EngineConfig.PG.TrackCommitTimestamp)
	assert.Equal(t, "all", *db.EngineConfig.PG.TrackFunctions)
	assert.Equal(t, "off", *db.EngineConfig.PG.TrackIOTiming)
	assert.Equal(t, 60000, *db.EngineConfig.PG.WALSenderTimeout)
	assert.Equal(t, 50, *db.EngineConfig.PG.WALWriterDelay)
}

func TestDatabasePostgreSQL_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("postgresql_database_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.PostgresUpdateOptions{
		Label: "example-db-updated",
		EngineConfig: &linodego.PostgresDatabaseEngineConfig{
			PG: &linodego.PostgresDatabaseEngineConfigPG{
				AutovacuumMaxWorkers: linodego.Pointer(10),
			},
		},
	}

	base.MockPut("databases/postgresql/instances/123", fixtureData)

	db, err := base.Client.UpdatePostgresDatabase(context.Background(), 123, requestData)
	assert.NoError(t, err)

	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, "postgresql", db.Engine)
	assert.Equal(t, 123, db.ID)
	assert.Equal(t, "example-db-updated", db.Label)
	assert.Equal(t, linodego.DatabasePlatform("rdbms-default"), db.Platform)
	assert.Equal(t, 3306, db.Port)
	assert.Equal(t, "us-east", db.Region)
	assert.Equal(t, linodego.DatabaseStatus("active"), db.Status)
	assert.Equal(t, 15, db.TotalDiskSizeGB)
	assert.Equal(t, "g6-dedicated-2", db.Type)
	assert.Equal(t, 3, db.Updates.Duration)
	assert.Equal(t, linodego.DatabaseDayOfWeek(1), db.Updates.DayOfWeek)
	assert.Equal(t, linodego.DatabaseMaintenanceFrequency("weekly"), db.Updates.Frequency)
	assert.Equal(t, 0, db.Updates.HourOfDay)
	assert.Equal(t, 2, db.UsedDiskSizeGB)
	assert.Equal(t, "13.2", db.Version)

	assert.Equal(t, true, *db.EngineConfig.PGStatMonitorEnable)
	assert.Equal(t, int64(1000), *db.EngineConfig.PGLookout.MaxFailoverReplicationTimeLag)
	assert.Equal(t, 41.5, *db.EngineConfig.SharedBuffersPercentage)
	assert.Equal(t, 4, *db.EngineConfig.WorkMem)
	assert.Equal(t, 0.5, *db.EngineConfig.PG.AutovacuumAnalyzeScaleFactor)
	assert.Equal(t, int32(100), *db.EngineConfig.PG.AutovacuumAnalyzeThreshold)
	assert.Equal(t, 10, *db.EngineConfig.PG.AutovacuumMaxWorkers)
	assert.Equal(t, 100, *db.EngineConfig.PG.AutovacuumNaptime)
	assert.Equal(t, 50, *db.EngineConfig.PG.AutovacuumVacuumCostDelay)
	assert.Equal(t, 100, *db.EngineConfig.PG.AutovacuumVacuumCostLimit)
	assert.Equal(t, 0.5, *db.EngineConfig.PG.AutovacuumVacuumScaleFactor)
	assert.Equal(t, int32(100), *db.EngineConfig.PG.AutovacuumVacuumThreshold)
	assert.Equal(t, 200, *db.EngineConfig.PG.BGWriterDelay)
	assert.Equal(t, 512, *db.EngineConfig.PG.BGWriterFlushAfter)
	assert.Equal(t, 100, *db.EngineConfig.PG.BGWriterLRUMaxPages)
	assert.Equal(t, 2.0, *db.EngineConfig.PG.BGWriterLRUMultiplier)
	assert.Equal(t, 1000, *db.EngineConfig.PG.DeadlockTimeout)
	assert.Equal(t, "lz4", *db.EngineConfig.PG.DefaultToastCompression)
	assert.Equal(t, 100, *db.EngineConfig.PG.IdleInTransactionSessionTimeout)
	assert.Equal(t, true, *db.EngineConfig.PG.JIT)
	assert.Equal(t, 100, *db.EngineConfig.PG.MaxFilesPerProcess)
	assert.Equal(t, 100, *db.EngineConfig.PG.MaxLocksPerTransaction)
	assert.Equal(t, 32, *db.EngineConfig.PG.MaxLogicalReplicationWorkers)
	assert.Equal(t, 64, *db.EngineConfig.PG.MaxParallelWorkers)
	assert.Equal(t, 64, *db.EngineConfig.PG.MaxParallelWorkersPerGather)
	assert.Equal(t, 1000, *db.EngineConfig.PG.MaxPredLocksPerTransaction)
	assert.Equal(t, 32, *db.EngineConfig.PG.MaxReplicationSlots)
	assert.Equal(t, int32(100), *db.EngineConfig.PG.MaxSlotWALKeepSize)
	assert.Equal(t, 3507152, *db.EngineConfig.PG.MaxStackDepth)
	assert.Equal(t, 1000, *db.EngineConfig.PG.MaxStandbyArchiveDelay)
	assert.Equal(t, 1000, *db.EngineConfig.PG.MaxStandbyStreamingDelay)
	assert.Equal(t, 32, *db.EngineConfig.PG.MaxWALSenders)
	assert.Equal(t, 64, *db.EngineConfig.PG.MaxWorkerProcesses)
	assert.Equal(t, "scram-sha-256", *db.EngineConfig.PG.PasswordEncryption)
	assert.Equal(t, 3600, *db.EngineConfig.PG.PGPartmanBGWInterval)
	assert.Equal(t, "myrolename", *db.EngineConfig.PG.PGPartmanBGWRole)
	assert.Equal(t, false, *db.EngineConfig.PG.PGStatMonitorPGSMEnableQueryPlan)
	assert.Equal(t, 10, *db.EngineConfig.PG.PGStatMonitorPGSMMaxBuckets)
	assert.Equal(t, "top", *db.EngineConfig.PG.PGStatStatementsTrack)
	assert.Equal(t, int32(5000000), *db.EngineConfig.PG.TempFileLimit)
	assert.Equal(t, "Europe/Helsinki", *db.EngineConfig.PG.Timezone)
	assert.Equal(t, 1024, *db.EngineConfig.PG.TrackActivityQuerySize)
	assert.Equal(t, "off", *db.EngineConfig.PG.TrackCommitTimestamp)
	assert.Equal(t, "all", *db.EngineConfig.PG.TrackFunctions)
	assert.Equal(t, "off", *db.EngineConfig.PG.TrackIOTiming)
	assert.Equal(t, 60000, *db.EngineConfig.PG.WALSenderTimeout)
	assert.Equal(t, 50, *db.EngineConfig.PG.WALWriterDelay)
}

func TestDatabasePostgreSQL_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("postgresql_database_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.PostgresCreateOptions{
		Label:  "example-db-created",
		Region: "us-east",
		Type:   "g6-dedicated-2",
		Engine: "postgresql",
		EngineConfig: &linodego.PostgresDatabaseEngineConfig{
			PG: &linodego.PostgresDatabaseEngineConfigPG{
				AutovacuumMaxWorkers: linodego.Pointer(10),
			},
		},
	}

	base.MockPost("databases/postgresql/instances", fixtureData)

	db, err := base.Client.CreatePostgresDatabase(context.Background(), requestData)
	assert.NoError(t, err)

	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, "postgresql", db.Engine)
	assert.Equal(t, 123, db.ID)
	assert.Equal(t, "example-db-created", db.Label)
	assert.Equal(t, linodego.DatabasePlatform("rdbms-default"), db.Platform)
	assert.Equal(t, 3306, db.Port)
	assert.Equal(t, "us-east", db.Region)
	assert.Equal(t, linodego.DatabaseStatus("active"), db.Status)
	assert.Equal(t, 15, db.TotalDiskSizeGB)
	assert.Equal(t, "g6-dedicated-2", db.Type)
	assert.Equal(t, 3, db.Updates.Duration)
	assert.Equal(t, linodego.DatabaseDayOfWeek(1), db.Updates.DayOfWeek)
	assert.Equal(t, linodego.DatabaseMaintenanceFrequency("weekly"), db.Updates.Frequency)
	assert.Equal(t, 0, db.Updates.HourOfDay)
	assert.Equal(t, 2, db.UsedDiskSizeGB)
	assert.Equal(t, "13.2", db.Version)

	assert.Equal(t, true, *db.EngineConfig.PGStatMonitorEnable)
	assert.Equal(t, int64(1000), *db.EngineConfig.PGLookout.MaxFailoverReplicationTimeLag)
	assert.Equal(t, 41.5, *db.EngineConfig.SharedBuffersPercentage)
	assert.Equal(t, 4, *db.EngineConfig.WorkMem)
	assert.Equal(t, 0.5, *db.EngineConfig.PG.AutovacuumAnalyzeScaleFactor)
	assert.Equal(t, int32(100), *db.EngineConfig.PG.AutovacuumAnalyzeThreshold)
	assert.Equal(t, 10, *db.EngineConfig.PG.AutovacuumMaxWorkers)
	assert.Equal(t, 100, *db.EngineConfig.PG.AutovacuumNaptime)
	assert.Equal(t, 50, *db.EngineConfig.PG.AutovacuumVacuumCostDelay)
	assert.Equal(t, 100, *db.EngineConfig.PG.AutovacuumVacuumCostLimit)
	assert.Equal(t, 0.5, *db.EngineConfig.PG.AutovacuumVacuumScaleFactor)
	assert.Equal(t, int32(100), *db.EngineConfig.PG.AutovacuumVacuumThreshold)
	assert.Equal(t, 200, *db.EngineConfig.PG.BGWriterDelay)
	assert.Equal(t, 512, *db.EngineConfig.PG.BGWriterFlushAfter)
	assert.Equal(t, 100, *db.EngineConfig.PG.BGWriterLRUMaxPages)
	assert.Equal(t, 2.0, *db.EngineConfig.PG.BGWriterLRUMultiplier)
	assert.Equal(t, 1000, *db.EngineConfig.PG.DeadlockTimeout)
	assert.Equal(t, "lz4", *db.EngineConfig.PG.DefaultToastCompression)
	assert.Equal(t, 100, *db.EngineConfig.PG.IdleInTransactionSessionTimeout)
	assert.Equal(t, true, *db.EngineConfig.PG.JIT)
	assert.Equal(t, 100, *db.EngineConfig.PG.MaxFilesPerProcess)
	assert.Equal(t, 100, *db.EngineConfig.PG.MaxLocksPerTransaction)
	assert.Equal(t, 32, *db.EngineConfig.PG.MaxLogicalReplicationWorkers)
	assert.Equal(t, 64, *db.EngineConfig.PG.MaxParallelWorkers)
	assert.Equal(t, 64, *db.EngineConfig.PG.MaxParallelWorkersPerGather)
	assert.Equal(t, 1000, *db.EngineConfig.PG.MaxPredLocksPerTransaction)
	assert.Equal(t, 32, *db.EngineConfig.PG.MaxReplicationSlots)
	assert.Equal(t, int32(100), *db.EngineConfig.PG.MaxSlotWALKeepSize)
	assert.Equal(t, 3507152, *db.EngineConfig.PG.MaxStackDepth)
	assert.Equal(t, 1000, *db.EngineConfig.PG.MaxStandbyArchiveDelay)
	assert.Equal(t, 1000, *db.EngineConfig.PG.MaxStandbyStreamingDelay)
	assert.Equal(t, 32, *db.EngineConfig.PG.MaxWALSenders)
	assert.Equal(t, 64, *db.EngineConfig.PG.MaxWorkerProcesses)
	assert.Equal(t, "scram-sha-256", *db.EngineConfig.PG.PasswordEncryption)
	assert.Equal(t, 3600, *db.EngineConfig.PG.PGPartmanBGWInterval)
	assert.Equal(t, "myrolename", *db.EngineConfig.PG.PGPartmanBGWRole)
	assert.Equal(t, false, *db.EngineConfig.PG.PGStatMonitorPGSMEnableQueryPlan)
	assert.Equal(t, 10, *db.EngineConfig.PG.PGStatMonitorPGSMMaxBuckets)
	assert.Equal(t, "top", *db.EngineConfig.PG.PGStatStatementsTrack)
	assert.Equal(t, int32(5000000), *db.EngineConfig.PG.TempFileLimit)
	assert.Equal(t, "Europe/Helsinki", *db.EngineConfig.PG.Timezone)
	assert.Equal(t, 1024, *db.EngineConfig.PG.TrackActivityQuerySize)
	assert.Equal(t, "off", *db.EngineConfig.PG.TrackCommitTimestamp)
	assert.Equal(t, "all", *db.EngineConfig.PG.TrackFunctions)
	assert.Equal(t, "off", *db.EngineConfig.PG.TrackIOTiming)
	assert.Equal(t, 60000, *db.EngineConfig.PG.WALSenderTimeout)
	assert.Equal(t, 50, *db.EngineConfig.PG.WALWriterDelay)
}

func TestDatabasePostgreSQL_Delete(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, "databases/postgresql/instances/123"), httpmock.NewStringResponder(200, "{}"))

	if err := client.DeletePostgresDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabasePostgreSQL_SSL_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("postgresql_database_ssl_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/postgresql/instances/123/ssl", fixtureData)

	ssl, err := base.Client.GetPostgresDatabaseSSL(context.Background(), 123)
	assert.NoError(t, err)

	expectedCACertificate := []byte("-----BEGIN CERTIFICATE-----\nThis is a test certificate\n-----END CERTIFICATE-----\n")

	assert.Equal(t, expectedCACertificate, ssl.CACertificate)
}

func TestDatabasePostgreSQL_Credentials_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("postgresql_database_credentials_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/postgresql/instances/123/credentials", fixtureData)

	creds, err := base.Client.GetPostgresDatabaseCredentials(context.Background(), 123)
	assert.NoError(t, err)

	assert.Equal(t, "linroot", creds.Username)
	assert.Equal(t, "s3cur3P@ssw0rd", creds.Password)
}

func TestDatabasePostgreSQL_Credentials_Reset(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/postgresql/instances/123/credentials/reset"), httpmock.NewStringResponder(200, "{}"))

	if err := client.ResetPostgresDatabaseCredentials(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabasePostgreSQL_Patch(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/postgresql/instances/123/patch"), httpmock.NewStringResponder(200, "{}"))

	if err := client.PatchPostgresDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabasePostgreSQL_Suspend(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/postgresql/instances/123/suspend"), httpmock.NewStringResponder(200, "{}"))

	if err := client.SuspendPostgresDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabasePostgreSQL_Resume(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/postgresql/instances/123/resume"), httpmock.NewStringResponder(200, "{}"))

	if err := client.ResumePostgresDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabasePostgreSQLConfig_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("postgresql_database_config_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/postgresql/config", fixtureData)

	config, err := base.Client.GetPostgresDatabaseConfig(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, "Specifies a fraction of the table size to add to autovacuum_analyze_threshold when "+
		"deciding whether to trigger an ANALYZE. The default is 0.2 (20% of table size)",
		config.PG.AutovacuumAnalyzeScaleFactor.Description)
	assert.Equal(t, 1.0, config.PG.AutovacuumAnalyzeScaleFactor.Maximum)
	assert.Equal(t, 0.0, config.PG.AutovacuumAnalyzeScaleFactor.Minimum)
	assert.False(t, config.PG.AutovacuumAnalyzeScaleFactor.RequiresRestart)
	assert.Equal(t, "number", config.PG.AutovacuumAnalyzeScaleFactor.Type)

	assert.Equal(t, "Specifies the minimum number of inserted, updated or deleted tuples needed to trigger an ANALYZE in any one table. The default is 50 tuples.",
		config.PG.AutovacuumAnalyzeThreshold.Description)
	assert.Equal(t, int32(2147483647), config.PG.AutovacuumAnalyzeThreshold.Maximum)
	assert.Equal(t, int32(0), config.PG.AutovacuumAnalyzeThreshold.Minimum)
	assert.False(t, config.PG.AutovacuumAnalyzeThreshold.RequiresRestart)
	assert.Equal(t, "integer", config.PG.AutovacuumAnalyzeThreshold.Type)

	assert.Equal(t, "Specifies the maximum number of autovacuum processes (other than the autovacuum launcher) that may be running at any one time. The default is three. This parameter can only be set at server start.",
		config.PG.AutovacuumMaxWorkers.Description)
	assert.Equal(t, 20, config.PG.AutovacuumMaxWorkers.Maximum)
	assert.Equal(t, 1, config.PG.AutovacuumMaxWorkers.Minimum)
	assert.False(t, config.PG.AutovacuumMaxWorkers.RequiresRestart)
	assert.Equal(t, "integer", config.PG.AutovacuumMaxWorkers.Type)

	assert.Equal(t, "Specifies the minimum delay between autovacuum runs on any given database. The delay is measured in seconds, and the default is one minute",
		config.PG.AutovacuumNaptime.Description)
	assert.Equal(t, 86400, config.PG.AutovacuumNaptime.Maximum)
	assert.Equal(t, 1, config.PG.AutovacuumNaptime.Minimum)
	assert.False(t, config.PG.AutovacuumNaptime.RequiresRestart)
	assert.Equal(t, "integer", config.PG.AutovacuumNaptime.Type)

	assert.Equal(t, "Specifies the cost delay value that will be used in automatic VACUUM operations. If -1 is specified, the regular vacuum_cost_delay value will be used. The default value is 20 milliseconds",
		config.PG.AutovacuumVacuumCostDelay.Description)
	assert.Equal(t, 100, config.PG.AutovacuumVacuumCostDelay.Maximum)
	assert.Equal(t, -1, config.PG.AutovacuumVacuumCostDelay.Minimum)
	assert.False(t, config.PG.AutovacuumVacuumCostDelay.RequiresRestart)
	assert.Equal(t, "integer", config.PG.AutovacuumVacuumCostDelay.Type)

	assert.Equal(t, "Specifies the cost limit value that will be used in automatic VACUUM operations. If -1 is specified (which is the default), the regular vacuum_cost_limit value will be used.",
		config.PG.AutovacuumVacuumCostLimit.Description)
	assert.Equal(t, 10000, config.PG.AutovacuumVacuumCostLimit.Maximum)
	assert.Equal(t, -1, config.PG.AutovacuumVacuumCostLimit.Minimum)
	assert.False(t, config.PG.AutovacuumVacuumCostLimit.RequiresRestart)
	assert.Equal(t, "integer", config.PG.AutovacuumVacuumCostLimit.Type)

	assert.Equal(t, "Specifies a fraction of the table size to add to autovacuum_vacuum_threshold when deciding whether to trigger a VACUUM. The default is 0.2 (20% of table size)",
		config.PG.AutovacuumVacuumScaleFactor.Description)
	assert.Equal(t, 1.0, config.PG.AutovacuumVacuumScaleFactor.Maximum)
	assert.Equal(t, 0.0, config.PG.AutovacuumVacuumScaleFactor.Minimum)
	assert.False(t, config.PG.AutovacuumVacuumScaleFactor.RequiresRestart)
	assert.Equal(t, "number", config.PG.AutovacuumVacuumScaleFactor.Type)

	assert.Equal(t, "Specifies the minimum number of updated or deleted tuples needed to trigger a VACUUM in any one table. The default is 50 tuples",
		config.PG.AutovacuumVacuumThreshold.Description)
	assert.Equal(t, int32(2147483647), config.PG.AutovacuumVacuumThreshold.Maximum)
	assert.Equal(t, int32(0), config.PG.AutovacuumVacuumThreshold.Minimum)
	assert.False(t, config.PG.AutovacuumVacuumThreshold.RequiresRestart)
	assert.Equal(t, "integer", config.PG.AutovacuumVacuumThreshold.Type)

	assert.Equal(t, "Specifies the delay between activity rounds for the background writer in milliseconds. Default is 200.",
		config.PG.BGWriterDelay.Description)
	assert.Equal(t, 200, config.PG.BGWriterDelay.Example)
	assert.Equal(t, 10000, config.PG.BGWriterDelay.Maximum)
	assert.Equal(t, 10, config.PG.BGWriterDelay.Minimum)
	assert.False(t, config.PG.BGWriterDelay.RequiresRestart)
	assert.Equal(t, "integer", config.PG.BGWriterDelay.Type)

	assert.Equal(t, "Whenever more than bgwriter_flush_after bytes have been written by the background writer, attempt to force the OS to issue these writes to the underlying storage. Specified in kilobytes, default is 512. Setting of 0 disables forced writeback.",
		config.PG.BGWriterFlushAfter.Description)
	assert.Equal(t, 512, config.PG.BGWriterFlushAfter.Example)
	assert.Equal(t, 2048, config.PG.BGWriterFlushAfter.Maximum)
	assert.Equal(t, 0, config.PG.BGWriterFlushAfter.Minimum)
	assert.False(t, config.PG.BGWriterFlushAfter.RequiresRestart)
	assert.Equal(t, "integer", config.PG.BGWriterFlushAfter.Type)

	assert.Equal(t, "In each round, no more than this many buffers will be written by the background writer. Setting this to zero disables background writing. Default is 100.",
		config.PG.BGWriterLRUMaxPages.Description)
	assert.Equal(t, 100, config.PG.BGWriterLRUMaxPages.Example)
	assert.Equal(t, 1073741823, config.PG.BGWriterLRUMaxPages.Maximum)
	assert.Equal(t, 0, config.PG.BGWriterLRUMaxPages.Minimum)
	assert.False(t, config.PG.BGWriterLRUMaxPages.RequiresRestart)
	assert.Equal(t, "integer", config.PG.BGWriterLRUMaxPages.Type)

	assert.Equal(t, "The average recent need for new buffers is multiplied by bgwriter_lru_multiplier to arrive at an estimate of the number that will be needed during the next round, (up to bgwriter_lru_maxpages). 1.0 represents a “just in time” policy of writing exactly the number of buffers predicted to be needed. Larger values provide some cushion against spikes in demand, while smaller values intentionally leave writes to be done by server processes. The default is 2.0.",
		config.PG.BGWriterLRUMultiplier.Description)
	assert.Equal(t, 2.0, config.PG.BGWriterLRUMultiplier.Example)
	assert.Equal(t, 10.0, config.PG.BGWriterLRUMultiplier.Maximum)
	assert.Equal(t, 0.0, config.PG.BGWriterLRUMultiplier.Minimum)
	assert.False(t, config.PG.BGWriterLRUMultiplier.RequiresRestart)
	assert.Equal(t, "number", config.PG.BGWriterLRUMultiplier.Type)

	assert.Equal(t, "This is the amount of time, in milliseconds, to wait on a lock before checking to see if there is a deadlock condition.",
		config.PG.DeadlockTimeout.Description)
	assert.Equal(t, 1000, config.PG.DeadlockTimeout.Example)
	assert.Equal(t, 1800000, config.PG.DeadlockTimeout.Maximum)
	assert.Equal(t, 500, config.PG.DeadlockTimeout.Minimum)
	assert.False(t, config.PG.DeadlockTimeout.RequiresRestart)
	assert.Equal(t, "integer", config.PG.DeadlockTimeout.Type)

	assert.Equal(t, "Specifies the default TOAST compression method for values of compressible columns (the default is lz4).",
		config.PG.DefaultToastCompression.Description)
	assert.ElementsMatch(t, []string{"lz4", "pglz"}, config.PG.DefaultToastCompression.Enum)
	assert.Equal(t, "lz4", config.PG.DefaultToastCompression.Example)
	assert.False(t, config.PG.DefaultToastCompression.RequiresRestart)
	assert.Equal(t, "string", config.PG.DefaultToastCompression.Type)

	assert.Equal(t, "Time out sessions with open transactions after this number of milliseconds",
		config.PG.IdleInTransactionSessionTimeout.Description)
	assert.Equal(t, 604800000, config.PG.IdleInTransactionSessionTimeout.Maximum)
	assert.Equal(t, 0, config.PG.IdleInTransactionSessionTimeout.Minimum)
	assert.False(t, config.PG.IdleInTransactionSessionTimeout.RequiresRestart)
	assert.Equal(t, "integer", config.PG.IdleInTransactionSessionTimeout.Type)

	assert.Equal(t, "Controls system-wide use of Just-in-Time Compilation (JIT).",
		config.PG.JIT.Description)
	assert.Equal(t, true, config.PG.JIT.Example)
	assert.False(t, config.PG.JIT.RequiresRestart)
	assert.Equal(t, "boolean", config.PG.JIT.Type)

	assert.Equal(t, "PostgreSQL maximum number of files that can be open per process",
		config.PG.MaxFilesPerProcess.Description)
	assert.Equal(t, 4096, config.PG.MaxFilesPerProcess.Maximum)
	assert.Equal(t, 1000, config.PG.MaxFilesPerProcess.Minimum)
	assert.False(t, config.PG.MaxFilesPerProcess.RequiresRestart)
	assert.Equal(t, "integer", config.PG.MaxFilesPerProcess.Type)

	assert.Equal(t, "PostgreSQL maximum locks per transaction",
		config.PG.MaxLocksPerTransaction.Description)
	assert.Equal(t, 6400, config.PG.MaxLocksPerTransaction.Maximum)
	assert.Equal(t, 64, config.PG.MaxLocksPerTransaction.Minimum)
	assert.False(t, config.PG.MaxLocksPerTransaction.RequiresRestart)
	assert.Equal(t, "integer", config.PG.MaxLocksPerTransaction.Type)

	assert.Equal(t, "PostgreSQL maximum logical replication workers (taken from the pool of max_parallel_workers)",
		config.PG.MaxLogicalReplicationWorkers.Description)
	assert.Equal(t, 64, config.PG.MaxLogicalReplicationWorkers.Maximum)
	assert.Equal(t, 4, config.PG.MaxLogicalReplicationWorkers.Minimum)
	assert.False(t, config.PG.MaxLogicalReplicationWorkers.RequiresRestart)
	assert.Equal(t, "integer", config.PG.MaxLogicalReplicationWorkers.Type)

	assert.Equal(t, "Sets the maximum number of workers that the system can support for parallel queries",
		config.PG.MaxParallelWorkers.Description)
	assert.Equal(t, 96, config.PG.MaxParallelWorkers.Maximum)
	assert.Equal(t, 0, config.PG.MaxParallelWorkers.Minimum)
	assert.False(t, config.PG.MaxParallelWorkers.RequiresRestart)
	assert.Equal(t, "integer", config.PG.MaxParallelWorkers.Type)

	assert.Equal(t, "Sets the maximum number of workers that can be started by a single Gather or Gather Merge node",
		config.PG.MaxParallelWorkersPerGather.Description)
	assert.Equal(t, 96, config.PG.MaxParallelWorkersPerGather.Maximum)
	assert.Equal(t, 0, config.PG.MaxParallelWorkersPerGather.Minimum)
	assert.False(t, config.PG.MaxParallelWorkersPerGather.RequiresRestart)
	assert.Equal(t, "integer", config.PG.MaxParallelWorkersPerGather.Type)

	assert.Equal(t, "PostgreSQL maximum predicate locks per transaction",
		config.PG.MaxPredLocksPerTransaction.Description)
	assert.Equal(t, 5120, config.PG.MaxPredLocksPerTransaction.Maximum)
	assert.Equal(t, 64, config.PG.MaxPredLocksPerTransaction.Minimum)
	assert.False(t, config.PG.MaxPredLocksPerTransaction.RequiresRestart)
	assert.Equal(t, "integer", config.PG.MaxPredLocksPerTransaction.Type)

	assert.Equal(t, "PostgreSQL maximum replication slots",
		config.PG.MaxReplicationSlots.Description)
	assert.Equal(t, 64, config.PG.MaxReplicationSlots.Maximum)
	assert.Equal(t, 8, config.PG.MaxReplicationSlots.Minimum)
	assert.False(t, config.PG.MaxReplicationSlots.RequiresRestart)
	assert.Equal(t, "integer", config.PG.MaxReplicationSlots.Type)

	assert.Equal(t, "PostgreSQL maximum WAL size (MB) reserved for replication slots. Default is -1 (unlimited). wal_keep_size minimum WAL size setting takes precedence over this.",
		config.PG.MaxSlotWALKeepSize.Description)
	assert.Equal(t, int32(2147483647), config.PG.MaxSlotWALKeepSize.Maximum)
	assert.Equal(t, int32(-1), config.PG.MaxSlotWALKeepSize.Minimum)
	assert.False(t, config.PG.MaxSlotWALKeepSize.RequiresRestart)
	assert.Equal(t, "integer", config.PG.MaxSlotWALKeepSize.Type)

	assert.Equal(t, "Maximum depth of the stack in bytes",
		config.PG.MaxStackDepth.Description)
	assert.Equal(t, 6291456, config.PG.MaxStackDepth.Maximum)
	assert.Equal(t, 2097152, config.PG.MaxStackDepth.Minimum)
	assert.False(t, config.PG.MaxStackDepth.RequiresRestart)
	assert.Equal(t, "integer", config.PG.MaxStackDepth.Type)

	assert.Equal(t, "Max standby archive delay in milliseconds",
		config.PG.MaxStandbyArchiveDelay.Description)
	assert.Equal(t, 43200000, config.PG.MaxStandbyArchiveDelay.Maximum)
	assert.Equal(t, 1, config.PG.MaxStandbyArchiveDelay.Minimum)
	assert.False(t, config.PG.MaxStandbyArchiveDelay.RequiresRestart)
	assert.Equal(t, "integer", config.PG.MaxStandbyArchiveDelay.Type)

	assert.Equal(t, "Max standby streaming delay in milliseconds",
		config.PG.MaxStandbyStreamingDelay.Description)
	assert.Equal(t, 43200000, config.PG.MaxStandbyStreamingDelay.Maximum)
	assert.Equal(t, 1, config.PG.MaxStandbyStreamingDelay.Minimum)
	assert.False(t, config.PG.MaxStandbyStreamingDelay.RequiresRestart)
	assert.Equal(t, "integer", config.PG.MaxStandbyStreamingDelay.Type)

	assert.Equal(t, "PostgreSQL maximum WAL senders",
		config.PG.MaxWALSenders.Description)
	assert.Equal(t, 64, config.PG.MaxWALSenders.Maximum)
	assert.Equal(t, 20, config.PG.MaxWALSenders.Minimum)
	assert.False(t, config.PG.MaxWALSenders.RequiresRestart)
	assert.Equal(t, "integer", config.PG.MaxWALSenders.Type)

	assert.Equal(t, "Sets the maximum number of background processes that the system can support",
		config.PG.MaxWorkerProcesses.Description)
	assert.Equal(t, 96, config.PG.MaxWorkerProcesses.Maximum)
	assert.Equal(t, 8, config.PG.MaxWorkerProcesses.Minimum)
	assert.False(t, config.PG.MaxWorkerProcesses.RequiresRestart)
	assert.Equal(t, "integer", config.PG.MaxWorkerProcesses.Type)

	assert.Equal(t, "Chooses the algorithm for encrypting passwords.",
		config.PG.PasswordEncryption.Description)
	assert.Equal(t, []string{"md5", "scram-sha-256"}, config.PG.PasswordEncryption.Enum)
	assert.Equal(t, "scram-sha-256", config.PG.PasswordEncryption.Example)
	assert.False(t, config.PG.PasswordEncryption.RequiresRestart)
	assert.Equal(t, "string", config.PG.PasswordEncryption.Type)

	assert.Equal(t, "Sets the time interval to run pg_partman's scheduled tasks",
		config.PG.PGPartmanBGWInterval.Description)
	assert.Equal(t, 3600, config.PG.PGPartmanBGWInterval.Example)
	assert.Equal(t, 604800, config.PG.PGPartmanBGWInterval.Maximum)
	assert.Equal(t, 3600, config.PG.PGPartmanBGWInterval.Minimum)
	assert.False(t, config.PG.PGPartmanBGWInterval.RequiresRestart)
	assert.Equal(t, "integer", config.PG.PGPartmanBGWInterval.Type)

	assert.Equal(t, "Controls which role to use for pg_partman's scheduled background tasks.",
		config.PG.PGPartmanBGWRole.Description)
	assert.Equal(t, "myrolename", config.PG.PGPartmanBGWRole.Example)
	assert.Equal(t, 64, config.PG.PGPartmanBGWRole.MaxLength)
	assert.Equal(t, "^[_A-Za-z0-9][-._A-Za-z0-9]{0,63}$", config.PG.PGPartmanBGWRole.Pattern)
	assert.False(t, config.PG.PGPartmanBGWRole.RequiresRestart)
	assert.Equal(t, "string", config.PG.PGPartmanBGWRole.Type)

	assert.Equal(t, "Enables or disables query plan monitoring",
		config.PG.PGStatMonitorPGSMEnableQueryPlan.Description)
	assert.Equal(t, false, config.PG.PGStatMonitorPGSMEnableQueryPlan.Example)
	assert.False(t, config.PG.PGStatMonitorPGSMEnableQueryPlan.RequiresRestart)
	assert.Equal(t, "boolean", config.PG.PGStatMonitorPGSMEnableQueryPlan.Type)

	assert.Equal(t, "Sets the maximum number of buckets",
		config.PG.PGStatMonitorPGSMMaxBuckets.Description)
	assert.Equal(t, 10, config.PG.PGStatMonitorPGSMMaxBuckets.Example)
	assert.Equal(t, 10, config.PG.PGStatMonitorPGSMMaxBuckets.Maximum)
	assert.Equal(t, 1, config.PG.PGStatMonitorPGSMMaxBuckets.Minimum)
	assert.False(t, config.PG.PGStatMonitorPGSMMaxBuckets.RequiresRestart)
	assert.Equal(t, "integer", config.PG.PGStatMonitorPGSMMaxBuckets.Type)

	assert.Equal(t, "Controls which statements are counted. Specify top to track top-level statements (those issued directly by clients), all to also track nested statements (such as statements invoked within functions), or none to disable statement statistics collection. The default value is top.",
		config.PG.PGStatStatementsTrack.Description)
	assert.Equal(t, []string{"all", "top", "none"}, config.PG.PGStatStatementsTrack.Enum)
	assert.False(t, config.PG.PGStatStatementsTrack.RequiresRestart)
	assert.Equal(t, "string", config.PG.PGStatStatementsTrack.Type)

	assert.Equal(t, "PostgreSQL temporary file limit in KiB, -1 for unlimited",
		config.PG.TempFileLimit.Description)
	assert.Equal(t, int32(5000000), config.PG.TempFileLimit.Example)
	assert.Equal(t, int32(2147483647), config.PG.TempFileLimit.Maximum)
	assert.Equal(t, int32(-1), config.PG.TempFileLimit.Minimum)
	assert.False(t, config.PG.TempFileLimit.RequiresRestart)
	assert.Equal(t, "integer", config.PG.TempFileLimit.Type)

	assert.Equal(t, "PostgreSQL service timezone",
		config.PG.Timezone.Description)
	assert.Equal(t, "Europe/Helsinki", config.PG.Timezone.Example)
	assert.Equal(t, 64, config.PG.Timezone.MaxLength)
	assert.Equal(t, "^[\\w/]*$", config.PG.Timezone.Pattern)
	assert.False(t, config.PG.Timezone.RequiresRestart)
	assert.Equal(t, "string", config.PG.Timezone.Type)

	assert.Equal(t, "Specifies the number of bytes reserved to track the currently executing command for each active session.",
		config.PG.TrackActivityQuerySize.Description)
	assert.Equal(t, 1024, config.PG.TrackActivityQuerySize.Example)
	assert.Equal(t, 10240, config.PG.TrackActivityQuerySize.Maximum)
	assert.Equal(t, 1024, config.PG.TrackActivityQuerySize.Minimum)
	assert.False(t, config.PG.TrackActivityQuerySize.RequiresRestart)
	assert.Equal(t, "integer", config.PG.TrackActivityQuerySize.Type)

	assert.Equal(t, "Record commit time of transactions.",
		config.PG.TrackCommitTimestamp.Description)
	assert.Equal(t, "off", config.PG.TrackCommitTimestamp.Example)
	assert.Equal(t, []string{"off", "on"}, config.PG.TrackCommitTimestamp.Enum)
	assert.False(t, config.PG.TrackCommitTimestamp.RequiresRestart)
	assert.Equal(t, "string", config.PG.TrackCommitTimestamp.Type)

	assert.Equal(t, "Enables tracking of function call counts and time used.",
		config.PG.TrackFunctions.Description)
	assert.Equal(t, []string{"all", "pl", "none"}, config.PG.TrackFunctions.Enum)
	assert.False(t, config.PG.TrackFunctions.RequiresRestart)
	assert.Equal(t, "string", config.PG.TrackFunctions.Type)

	assert.Equal(t, "Enables timing of database I/O calls. This parameter is off by default, because it will repeatedly query the operating system for the current time, which may cause significant overhead on some platforms.",
		config.PG.TrackIOTiming.Description)
	assert.Equal(t, "off", config.PG.TrackIOTiming.Example)
	assert.Equal(t, []string{"off", "on"}, config.PG.TrackIOTiming.Enum)
	assert.False(t, config.PG.TrackIOTiming.RequiresRestart)
	assert.Equal(t, "string", config.PG.TrackIOTiming.Type)

	assert.Equal(t, "Terminate replication connections that are inactive for longer than this amount of time, in milliseconds. Setting this value to zero disables the timeout.",
		config.PG.WALSenderTimeout.Description)
	assert.Equal(t, 60000, config.PG.WALSenderTimeout.Example)
	assert.False(t, config.PG.WALSenderTimeout.RequiresRestart)
	assert.Equal(t, "integer", config.PG.WALSenderTimeout.Type)

	assert.Equal(t, "WAL flush interval in milliseconds. Note that setting this value to lower than the default 200ms may negatively impact performance",
		config.PG.WALWriterDelay.Description)
	assert.Equal(t, 50, config.PG.WALWriterDelay.Example)
	assert.Equal(t, 200, config.PG.WALWriterDelay.Maximum)
	assert.Equal(t, 10, config.PG.WALWriterDelay.Minimum)
	assert.False(t, config.PG.WALWriterDelay.RequiresRestart)
	assert.Equal(t, "integer", config.PG.WALWriterDelay.Type)

	assert.Equal(t, "Enable the pg_stat_monitor extension. Enabling this extension will cause the cluster to be restarted.When this extension is enabled, pg_stat_statements results for utility commands are unreliable",
		config.PGStatMonitorEnable.Description)
	assert.True(t, config.PGStatMonitorEnable.RequiresRestart)
	assert.Equal(t, "boolean", config.PGStatMonitorEnable.Type)

	assert.Equal(t, "Number of seconds of master unavailability before triggering database failover to standby",
		config.PGLookout.PGLookoutMaxFailoverReplicationTimeLag.Description)
	assert.Equal(t, int64(9223372036854775000), config.PGLookout.PGLookoutMaxFailoverReplicationTimeLag.Maximum)
	assert.Equal(t, int64(10), config.PGLookout.PGLookoutMaxFailoverReplicationTimeLag.Minimum)
	assert.False(t, config.PGLookout.PGLookoutMaxFailoverReplicationTimeLag.RequiresRestart)
	assert.Equal(t, "integer", config.PGLookout.PGLookoutMaxFailoverReplicationTimeLag.Type)

	assert.Equal(t, "Percentage of total RAM that the database server uses for shared memory buffers. Valid range is 20-60 (float), which corresponds to 20% - 60%. This setting adjusts the shared_buffers configuration value.",
		config.SharedBuffersPercentage.Description)
	assert.Equal(t, 41.5, config.SharedBuffersPercentage.Example)
	assert.Equal(t, 60.0, config.SharedBuffersPercentage.Maximum)
	assert.Equal(t, 20.0, config.SharedBuffersPercentage.Minimum)
	assert.False(t, config.SharedBuffersPercentage.RequiresRestart)
	assert.Equal(t, "number", config.SharedBuffersPercentage.Type)

	assert.Equal(t, "Sets the maximum amount of memory to be used by a query operation (such as a sort or hash table) before writing to temporary disk files, in MB. Default is 1MB + 0.075% of total RAM (up to 32MB).",
		config.WorkMem.Description)
	assert.Equal(t, 4, config.WorkMem.Example)
	assert.Equal(t, 1024, config.WorkMem.Maximum)
	assert.Equal(t, 1, config.WorkMem.Minimum)
	assert.False(t, config.WorkMem.RequiresRestart)
	assert.Equal(t, "integer", config.WorkMem.Type)
}
