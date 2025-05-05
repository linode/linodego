package unit

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestListDatabaseMySQL_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("mysql_databases_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/mysql/instances", fixtureData)
	databases, err := base.Client.ListMySQLDatabases(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, databases, "Expected non-empty mysql database list")
}

func TestDatabaseMySQL_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("mysql_database_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/mysql/instances/123", fixtureData)

	db, err := base.Client.GetMySQLDatabase(context.Background(), 123)
	assert.NoError(t, err)

	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, "mysql", db.Engine)
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
	assert.Equal(t, "8.0.26", db.Version)

	assert.Equal(t, 600, *db.EngineConfig.BinlogRetentionPeriod)
	assert.Equal(t, 10, *db.EngineConfig.MySQL.ConnectTimeout)
	assert.Equal(t, "+03:00", *db.EngineConfig.MySQL.DefaultTimeZone)
	assert.Equal(t, float64(1024), *db.EngineConfig.MySQL.GroupConcatMaxLen)
	assert.Equal(t, 86400, *db.EngineConfig.MySQL.InformationSchemaStatsExpiry)
	assert.Equal(t, 30, *db.EngineConfig.MySQL.InnoDBChangeBufferMaxSize)
	assert.Equal(t, 0, *db.EngineConfig.MySQL.InnoDBFlushNeighbors)
	assert.Equal(t, 3, *db.EngineConfig.MySQL.InnoDBFTMinTokenSize)
	assert.Equal(t, "db_name/table_name", **db.EngineConfig.MySQL.InnoDBFTServerStopwordTable)
	assert.Equal(t, 50, *db.EngineConfig.MySQL.InnoDBLockWaitTimeout)
	assert.Equal(t, 16777216, *db.EngineConfig.MySQL.InnoDBLogBufferSize)
	assert.Equal(t, 134217728, *db.EngineConfig.MySQL.InnoDBOnlineAlterLogMaxSize)
	assert.Equal(t, 10, *db.EngineConfig.MySQL.InnoDBReadIOThreads)
	assert.Equal(t, true, *db.EngineConfig.MySQL.InnoDBRollbackOnTimeout)
	assert.Equal(t, 10, *db.EngineConfig.MySQL.InnoDBThreadConcurrency)
	assert.Equal(t, 10, *db.EngineConfig.MySQL.InnoDBWriteIOThreads)
	assert.Equal(t, 3600, *db.EngineConfig.MySQL.InteractiveTimeout)
	assert.Equal(t, "TempTable", *db.EngineConfig.MySQL.InternalTmpMemStorageEngine)
	assert.Equal(t, 67108864, *db.EngineConfig.MySQL.MaxAllowedPacket)
	assert.Equal(t, 16777216, *db.EngineConfig.MySQL.MaxHeapTableSize)
	assert.Equal(t, 16384, *db.EngineConfig.MySQL.NetBufferLength)
	assert.Equal(t, 30, *db.EngineConfig.MySQL.NetReadTimeout)
	assert.Equal(t, 30, *db.EngineConfig.MySQL.NetWriteTimeout)
	assert.Equal(t, 262144, *db.EngineConfig.MySQL.SortBufferSize)
	assert.Equal(t, "ANSI,TRADITIONAL", *db.EngineConfig.MySQL.SQLMode)
	assert.Equal(t, true, *db.EngineConfig.MySQL.SQLRequirePrimaryKey)
	assert.Equal(t, 16777216, *db.EngineConfig.MySQL.TmpTableSize)
	assert.Equal(t, 28800, *db.EngineConfig.MySQL.WaitTimeout)
}

func TestDatabaseMySQL_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("mysql_database_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.MySQLUpdateOptions{
		Label: "example-db-updated",
		EngineConfig: &linodego.MySQLDatabaseEngineConfig{
			MySQL: &linodego.MySQLDatabaseEngineConfigMySQL{
				ConnectTimeout: linodego.Pointer(20),
			},
		},
	}

	base.MockPut("databases/mysql/instances/123", fixtureData)

	db, err := base.Client.UpdateMySQLDatabase(context.Background(), 123, requestData)
	assert.NoError(t, err)

	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, "mysql", db.Engine)
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
	assert.Equal(t, "8.0.26", db.Version)

	assert.Equal(t, 600, *db.EngineConfig.BinlogRetentionPeriod)
	assert.Equal(t, 20, *db.EngineConfig.MySQL.ConnectTimeout)
	assert.Equal(t, "+03:00", *db.EngineConfig.MySQL.DefaultTimeZone)
	assert.Equal(t, float64(1024), *db.EngineConfig.MySQL.GroupConcatMaxLen)
	assert.Equal(t, 86400, *db.EngineConfig.MySQL.InformationSchemaStatsExpiry)
	assert.Equal(t, 30, *db.EngineConfig.MySQL.InnoDBChangeBufferMaxSize)
	assert.Equal(t, 0, *db.EngineConfig.MySQL.InnoDBFlushNeighbors)
	assert.Equal(t, 3, *db.EngineConfig.MySQL.InnoDBFTMinTokenSize)
	assert.Equal(t, "db_name/table_name", **db.EngineConfig.MySQL.InnoDBFTServerStopwordTable)
	assert.Equal(t, 50, *db.EngineConfig.MySQL.InnoDBLockWaitTimeout)
	assert.Equal(t, 16777216, *db.EngineConfig.MySQL.InnoDBLogBufferSize)
	assert.Equal(t, 134217728, *db.EngineConfig.MySQL.InnoDBOnlineAlterLogMaxSize)
	assert.Equal(t, 10, *db.EngineConfig.MySQL.InnoDBReadIOThreads)
	assert.Equal(t, true, *db.EngineConfig.MySQL.InnoDBRollbackOnTimeout)
	assert.Equal(t, 10, *db.EngineConfig.MySQL.InnoDBThreadConcurrency)
	assert.Equal(t, 10, *db.EngineConfig.MySQL.InnoDBWriteIOThreads)
	assert.Equal(t, 3600, *db.EngineConfig.MySQL.InteractiveTimeout)
	assert.Equal(t, "TempTable", *db.EngineConfig.MySQL.InternalTmpMemStorageEngine)
	assert.Equal(t, 67108864, *db.EngineConfig.MySQL.MaxAllowedPacket)
	assert.Equal(t, 16777216, *db.EngineConfig.MySQL.MaxHeapTableSize)
	assert.Equal(t, 16384, *db.EngineConfig.MySQL.NetBufferLength)
	assert.Equal(t, 30, *db.EngineConfig.MySQL.NetReadTimeout)
	assert.Equal(t, 30, *db.EngineConfig.MySQL.NetWriteTimeout)
	assert.Equal(t, 262144, *db.EngineConfig.MySQL.SortBufferSize)
	assert.Equal(t, "ANSI,TRADITIONAL", *db.EngineConfig.MySQL.SQLMode)
	assert.Equal(t, true, *db.EngineConfig.MySQL.SQLRequirePrimaryKey)
	assert.Equal(t, 16777216, *db.EngineConfig.MySQL.TmpTableSize)
	assert.Equal(t, 28800, *db.EngineConfig.MySQL.WaitTimeout)
}

func TestDatabaseMySQL_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("mysql_database_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.MySQLCreateOptions{
		Label:  "example-db-created",
		Region: "us-east",
		Type:   "g6-dedicated-2",
		Engine: "mysql",
		EngineConfig: &linodego.MySQLDatabaseEngineConfig{
			MySQL: &linodego.MySQLDatabaseEngineConfigMySQL{
				ConnectTimeout: linodego.Pointer(20),
			},
		},
	}

	base.MockPost("databases/mysql/instances", fixtureData)

	db, err := base.Client.CreateMySQLDatabase(context.Background(), requestData)
	assert.NoError(t, err)

	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, "mysql", db.Engine)
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
	assert.Equal(t, "8.0.26", db.Version)

	assert.Equal(t, 600, *db.EngineConfig.BinlogRetentionPeriod)
	assert.Equal(t, 20, *db.EngineConfig.MySQL.ConnectTimeout)
	assert.Equal(t, "+03:00", *db.EngineConfig.MySQL.DefaultTimeZone)
	assert.Equal(t, float64(1024), *db.EngineConfig.MySQL.GroupConcatMaxLen)
	assert.Equal(t, 86400, *db.EngineConfig.MySQL.InformationSchemaStatsExpiry)
	assert.Equal(t, 30, *db.EngineConfig.MySQL.InnoDBChangeBufferMaxSize)
	assert.Equal(t, 0, *db.EngineConfig.MySQL.InnoDBFlushNeighbors)
	assert.Equal(t, 3, *db.EngineConfig.MySQL.InnoDBFTMinTokenSize)
	assert.Equal(t, "db_name/table_name", **db.EngineConfig.MySQL.InnoDBFTServerStopwordTable)
	assert.Equal(t, 50, *db.EngineConfig.MySQL.InnoDBLockWaitTimeout)
	assert.Equal(t, 16777216, *db.EngineConfig.MySQL.InnoDBLogBufferSize)
	assert.Equal(t, 134217728, *db.EngineConfig.MySQL.InnoDBOnlineAlterLogMaxSize)
	assert.Equal(t, 10, *db.EngineConfig.MySQL.InnoDBReadIOThreads)
	assert.Equal(t, true, *db.EngineConfig.MySQL.InnoDBRollbackOnTimeout)
	assert.Equal(t, 10, *db.EngineConfig.MySQL.InnoDBThreadConcurrency)
	assert.Equal(t, 10, *db.EngineConfig.MySQL.InnoDBWriteIOThreads)
	assert.Equal(t, 3600, *db.EngineConfig.MySQL.InteractiveTimeout)
	assert.Equal(t, "TempTable", *db.EngineConfig.MySQL.InternalTmpMemStorageEngine)
	assert.Equal(t, 67108864, *db.EngineConfig.MySQL.MaxAllowedPacket)
	assert.Equal(t, 16777216, *db.EngineConfig.MySQL.MaxHeapTableSize)
	assert.Equal(t, 16384, *db.EngineConfig.MySQL.NetBufferLength)
	assert.Equal(t, 30, *db.EngineConfig.MySQL.NetReadTimeout)
	assert.Equal(t, 30, *db.EngineConfig.MySQL.NetWriteTimeout)
	assert.Equal(t, 262144, *db.EngineConfig.MySQL.SortBufferSize)
	assert.Equal(t, "ANSI,TRADITIONAL", *db.EngineConfig.MySQL.SQLMode)
	assert.Equal(t, true, *db.EngineConfig.MySQL.SQLRequirePrimaryKey)
	assert.Equal(t, 16777216, *db.EngineConfig.MySQL.TmpTableSize)
	assert.Equal(t, 28800, *db.EngineConfig.MySQL.WaitTimeout)
}

func TestDatabaseMySQL_Delete(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, "databases/mysql/instances/123"), httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteMySQLDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseMySQL_SSL_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("mysql_database_ssl_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/mysql/instances/123/ssl", fixtureData)

	ssl, err := base.Client.GetMySQLDatabaseSSL(context.Background(), 123)
	assert.NoError(t, err)

	expectedCACertificate := []byte("-----BEGIN CERTIFICATE-----\nThis is a test certificate\n-----END CERTIFICATE-----\n")

	assert.Equal(t, expectedCACertificate, ssl.CACertificate)
}

func TestDatabaseMySQL_Credentials_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("mysql_database_credentials_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/mysql/instances/123/credentials", fixtureData)

	creds, err := base.Client.GetMySQLDatabaseCredentials(context.Background(), 123)
	assert.NoError(t, err)

	assert.Equal(t, "linroot", creds.Username)
	assert.Equal(t, "s3cur3P@ssw0rd", creds.Password)
}

func TestDatabaseMySQL_Credentials_Reset(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/mysql/instances/123/credentials/reset"), httpmock.NewStringResponder(200, "{}"))

	if err := client.ResetMySQLDatabaseCredentials(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseMySQL_Patch(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/mysql/instances/123/patch"), httpmock.NewStringResponder(200, "{}"))

	if err := client.PatchMySQLDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseMySQL_Suspend(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/mysql/instances/123/suspend"), httpmock.NewStringResponder(200, "{}"))

	if err := client.SuspendMySQLDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseMySQL_Resume(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/mysql/instances/123/resume"), httpmock.NewStringResponder(200, "{}"))

	if err := client.ResumeMySQLDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseMySQLConfig_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("mysql_database_config_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/mysql/config", fixtureData)

	config, err := base.Client.GetMySQLDatabaseConfig(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, "The number of seconds that the mysqld server waits for a connect packet before responding with Bad handshake",
		config.MySQL.ConnectTimeout.Description)
	assert.Equal(t, 10, config.MySQL.ConnectTimeout.Example)
	assert.Equal(t, 3600, config.MySQL.ConnectTimeout.Maximum)
	assert.Equal(t, 2, config.MySQL.ConnectTimeout.Minimum)
	assert.False(t, config.MySQL.ConnectTimeout.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.ConnectTimeout.Type)

	assert.Equal(t, "Default server time zone as an offset from UTC (from -12:00 to +12:00), a time zone name, or 'SYSTEM' to use the MySQL server default.",
		config.MySQL.DefaultTimeZone.Description)
	assert.Equal(t, "+03:00", config.MySQL.DefaultTimeZone.Example)
	assert.Equal(t, 100, config.MySQL.DefaultTimeZone.MaxLength)
	assert.Equal(t, 2, config.MySQL.DefaultTimeZone.MinLength)
	assert.Equal(t, "^([-+][\\d:]*|[\\w/]*)$", config.MySQL.DefaultTimeZone.Pattern)
	assert.False(t, config.MySQL.DefaultTimeZone.RequiresRestart)
	assert.Equal(t, "string", config.MySQL.DefaultTimeZone.Type)

	assert.Equal(t, "The maximum permitted result length in bytes for the GROUP_CONCAT() function.",
		config.MySQL.GroupConcatMaxLen.Description)
	assert.Equal(t, float64(1024), config.MySQL.GroupConcatMaxLen.Example)
	assert.Equal(t, float64(18446744073709551600), config.MySQL.GroupConcatMaxLen.Maximum)
	assert.Equal(t, float64(4), config.MySQL.GroupConcatMaxLen.Minimum)
	assert.False(t, config.MySQL.GroupConcatMaxLen.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.GroupConcatMaxLen.Type)

	assert.Equal(t, "The time, in seconds, before cached statistics expire",
		config.MySQL.InformationSchemaStatsExpiry.Description)
	assert.Equal(t, 86400, config.MySQL.InformationSchemaStatsExpiry.Example)
	assert.Equal(t, 31536000, config.MySQL.InformationSchemaStatsExpiry.Maximum)
	assert.Equal(t, 900, config.MySQL.InformationSchemaStatsExpiry.Minimum)
	assert.False(t, config.MySQL.InformationSchemaStatsExpiry.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InformationSchemaStatsExpiry.Type)

	assert.Equal(t, "Maximum size for the InnoDB change buffer, as a percentage of the total size of the buffer pool. Default is 25",
		config.MySQL.InnoDBChangeBufferMaxSize.Description)
	assert.Equal(t, 30, config.MySQL.InnoDBChangeBufferMaxSize.Example)
	assert.Equal(t, 50, config.MySQL.InnoDBChangeBufferMaxSize.Maximum)
	assert.Equal(t, 0, config.MySQL.InnoDBChangeBufferMaxSize.Minimum)
	assert.False(t, config.MySQL.InnoDBChangeBufferMaxSize.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBChangeBufferMaxSize.Type)

	assert.Equal(t, "Specifies whether flushing a page from the InnoDB buffer pool also flushes other dirty pages in the same extent (default is 1): 0 - dirty pages in the same extent are not flushed, 1 - flush contiguous dirty pages in the same extent, 2 - flush dirty pages in the same extent",
		config.MySQL.InnoDBFlushNeighbors.Description)
	assert.Equal(t, 0, config.MySQL.InnoDBFlushNeighbors.Example)
	assert.Equal(t, 2, config.MySQL.InnoDBFlushNeighbors.Maximum)
	assert.Equal(t, 0, config.MySQL.InnoDBFlushNeighbors.Minimum)
	assert.False(t, config.MySQL.InnoDBFlushNeighbors.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBFlushNeighbors.Type)

	assert.Equal(t, "Minimum length of words that are stored in an InnoDB FULLTEXT index. Changing this parameter will lead to a restart of the MySQL service.",
		config.MySQL.InnoDBFTMinTokenSize.Description)
	assert.Equal(t, 3, config.MySQL.InnoDBFTMinTokenSize.Example)
	assert.Equal(t, 16, config.MySQL.InnoDBFTMinTokenSize.Maximum)
	assert.Equal(t, 0, config.MySQL.InnoDBFTMinTokenSize.Minimum)
	assert.True(t, config.MySQL.InnoDBFTMinTokenSize.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBFTMinTokenSize.Type)

	assert.Equal(t, "This option is used to specify your own InnoDB FULLTEXT index stopword list for all InnoDB tables.",
		config.MySQL.InnoDBFTServerStopwordTable.Description)
	assert.Equal(t, "db_name/table_name", config.MySQL.InnoDBFTServerStopwordTable.Example)
	assert.Equal(t, 1024, config.MySQL.InnoDBFTServerStopwordTable.MaxLength)
	assert.Equal(t, "^.+/.+$", config.MySQL.InnoDBFTServerStopwordTable.Pattern)
	assert.False(t, config.MySQL.InnoDBFTServerStopwordTable.RequiresRestart)
	assert.Equal(t, []string{"null", "string"}, config.MySQL.InnoDBFTServerStopwordTable.Type)

	assert.Equal(t, "The length of time in seconds an InnoDB transaction waits for a row lock before giving up. Default is 120.",
		config.MySQL.InnoDBLockWaitTimeout.Description)
	assert.Equal(t, 50, config.MySQL.InnoDBLockWaitTimeout.Example)
	assert.Equal(t, 3600, config.MySQL.InnoDBLockWaitTimeout.Maximum)
	assert.Equal(t, 1, config.MySQL.InnoDBLockWaitTimeout.Minimum)
	assert.False(t, config.MySQL.InnoDBLockWaitTimeout.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBLockWaitTimeout.Type)

	assert.Equal(t, "The size in bytes of the buffer that InnoDB uses to write to the log files on disk.",
		config.MySQL.InnoDBLogBufferSize.Description)
	assert.Equal(t, 16777216, config.MySQL.InnoDBLogBufferSize.Example)
	assert.Equal(t, 4294967295, config.MySQL.InnoDBLogBufferSize.Maximum)
	assert.Equal(t, 1048576, config.MySQL.InnoDBLogBufferSize.Minimum)
	assert.False(t, config.MySQL.InnoDBLogBufferSize.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBLogBufferSize.Type)

	assert.Equal(t, "The upper limit in bytes on the size of the temporary log files used during online DDL operations for InnoDB tables.",
		config.MySQL.InnoDBOnlineAlterLogMaxSize.Description)
	assert.Equal(t, 134217728, config.MySQL.InnoDBOnlineAlterLogMaxSize.Example)
	assert.Equal(t, 1099511627776, config.MySQL.InnoDBOnlineAlterLogMaxSize.Maximum)
	assert.Equal(t, 65536, config.MySQL.InnoDBOnlineAlterLogMaxSize.Minimum)
	assert.False(t, config.MySQL.InnoDBOnlineAlterLogMaxSize.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBOnlineAlterLogMaxSize.Type)

	assert.Equal(t, "The number of I/O threads for read operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.",
		config.MySQL.InnoDBReadIOThreads.Description)
	assert.Equal(t, 10, config.MySQL.InnoDBReadIOThreads.Example)
	assert.Equal(t, 64, config.MySQL.InnoDBReadIOThreads.Maximum)
	assert.Equal(t, 1, config.MySQL.InnoDBReadIOThreads.Minimum)
	assert.True(t, config.MySQL.InnoDBReadIOThreads.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBReadIOThreads.Type)

	assert.Equal(t, "When enabled a transaction timeout causes InnoDB to abort and roll back the entire transaction. Changing this parameter will lead to a restart of the MySQL service.",
		config.MySQL.InnoDBRollbackOnTimeout.Description)
	assert.Equal(t, true, config.MySQL.InnoDBRollbackOnTimeout.Example)
	assert.True(t, config.MySQL.InnoDBRollbackOnTimeout.RequiresRestart)
	assert.Equal(t, "boolean", config.MySQL.InnoDBRollbackOnTimeout.Type)

	assert.Equal(t, "Defines the maximum number of threads permitted inside of InnoDB. Default is 0 (infinite concurrency - no limit)",
		config.MySQL.InnoDBThreadConcurrency.Description)
	assert.Equal(t, 10, config.MySQL.InnoDBThreadConcurrency.Example)
	assert.Equal(t, 1000, config.MySQL.InnoDBThreadConcurrency.Maximum)
	assert.Equal(t, 0, config.MySQL.InnoDBThreadConcurrency.Minimum)
	assert.False(t, config.MySQL.InnoDBThreadConcurrency.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBThreadConcurrency.Type)

	assert.Equal(t, "The number of I/O threads for write operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.",
		config.MySQL.InnoDBWriteIOThreads.Description)
	assert.Equal(t, 10, config.MySQL.InnoDBWriteIOThreads.Example)
	assert.Equal(t, 64, config.MySQL.InnoDBWriteIOThreads.Maximum)
	assert.Equal(t, 1, config.MySQL.InnoDBWriteIOThreads.Minimum)
	assert.True(t, config.MySQL.InnoDBWriteIOThreads.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBWriteIOThreads.Type)

	assert.Equal(t, "The number of seconds the server waits for activity on an interactive connection before closing it.",
		config.MySQL.InteractiveTimeout.Description)
	assert.Equal(t, 3600, config.MySQL.InteractiveTimeout.Example)
	assert.Equal(t, 604800, config.MySQL.InteractiveTimeout.Maximum)
	assert.Equal(t, 30, config.MySQL.InteractiveTimeout.Minimum)
	assert.False(t, config.MySQL.InteractiveTimeout.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InteractiveTimeout.Type)

	assert.Equal(t, "The storage engine for in-memory internal temporary tables.",
		config.MySQL.InternalTmpMemStorageEngine.Description)
	assert.Equal(t, "TempTable", config.MySQL.InternalTmpMemStorageEngine.Example)
	assert.Equal(t, []string{"TempTable", "MEMORY"}, config.MySQL.InternalTmpMemStorageEngine.Enum)
	assert.False(t, config.MySQL.InternalTmpMemStorageEngine.RequiresRestart)
	assert.Equal(t, "string", config.MySQL.InternalTmpMemStorageEngine.Type)

	assert.Equal(t, "Size of the largest message in bytes that can be received by the server. Default is 67108864 (64M)",
		config.MySQL.MaxAllowedPacket.Description)
	assert.Equal(t, 67108864, config.MySQL.MaxAllowedPacket.Example)
	assert.Equal(t, 1073741824, config.MySQL.MaxAllowedPacket.Maximum)
	assert.Equal(t, 102400, config.MySQL.MaxAllowedPacket.Minimum)
	assert.False(t, config.MySQL.MaxAllowedPacket.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.MaxAllowedPacket.Type)

	assert.Equal(t, "Limits the size of internal in-memory tables. Also set tmp_table_size. Default is 16777216 (16M)",
		config.MySQL.MaxHeapTableSize.Description)
	assert.Equal(t, 16777216, config.MySQL.MaxHeapTableSize.Example)
	assert.Equal(t, 1073741824, config.MySQL.MaxHeapTableSize.Maximum)
	assert.Equal(t, 1048576, config.MySQL.MaxHeapTableSize.Minimum)
	assert.False(t, config.MySQL.MaxHeapTableSize.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.MaxHeapTableSize.Type)

	assert.Equal(t, "Start sizes of connection buffer and result buffer. Default is 16384 (16K). Changing this parameter will lead to a restart of the MySQL service.",
		config.MySQL.NetBufferLength.Description)
	assert.Equal(t, 16384, config.MySQL.NetBufferLength.Example)
	assert.Equal(t, 1048576, config.MySQL.NetBufferLength.Maximum)
	assert.Equal(t, 1024, config.MySQL.NetBufferLength.Minimum)
	assert.True(t, config.MySQL.NetBufferLength.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.NetBufferLength.Type)

	assert.Equal(t, "The number of seconds to wait for more data from a connection before aborting the read.",
		config.MySQL.NetReadTimeout.Description)
	assert.Equal(t, 30, config.MySQL.NetReadTimeout.Example)
	assert.Equal(t, 3600, config.MySQL.NetReadTimeout.Maximum)
	assert.Equal(t, 1, config.MySQL.NetReadTimeout.Minimum)
	assert.False(t, config.MySQL.NetReadTimeout.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.NetReadTimeout.Type)

	assert.Equal(t, "The number of seconds to wait for a block to be written to a connection before aborting the write.",
		config.MySQL.NetWriteTimeout.Description)
	assert.Equal(t, 30, config.MySQL.NetWriteTimeout.Example)
	assert.Equal(t, 3600, config.MySQL.NetWriteTimeout.Maximum)
	assert.Equal(t, 1, config.MySQL.NetWriteTimeout.Minimum)
	assert.False(t, config.MySQL.NetWriteTimeout.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.NetWriteTimeout.Type)

	assert.Equal(t, "Sort buffer size in bytes for ORDER BY optimization. Default is 262144 (256K)",
		config.MySQL.SortBufferSize.Description)
	assert.Equal(t, 262144, config.MySQL.SortBufferSize.Example)
	assert.Equal(t, 1073741824, config.MySQL.SortBufferSize.Maximum)
	assert.Equal(t, 32768, config.MySQL.SortBufferSize.Minimum)
	assert.False(t, config.MySQL.SortBufferSize.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.SortBufferSize.Type)

	assert.Equal(t, "Global SQL mode. Set to empty to use MySQL server defaults. When creating a new service and not setting this field Akamai default SQL mode (strict, SQL standard compliant) will be assigned.",
		config.MySQL.SQLMode.Description)
	assert.Equal(t, "ANSI,TRADITIONAL", config.MySQL.SQLMode.Example)
	assert.Equal(t, 1024, config.MySQL.SQLMode.MaxLength)
	assert.Equal(t, "^[A-Z_]*(,[A-Z_]+)*$", config.MySQL.SQLMode.Pattern)
	assert.False(t, config.MySQL.SQLMode.RequiresRestart)
	assert.Equal(t, "string", config.MySQL.SQLMode.Type)

	assert.Equal(t, "Require primary key to be defined for new tables or old tables modified with ALTER TABLE and fail if missing. It is recommended to always have primary keys because various functionality may break if any large table is missing them.",
		config.MySQL.SQLRequirePrimaryKey.Description)
	assert.Equal(t, true, config.MySQL.SQLRequirePrimaryKey.Example)
	assert.False(t, config.MySQL.SQLRequirePrimaryKey.RequiresRestart)
	assert.Equal(t, "boolean", config.MySQL.SQLRequirePrimaryKey.Type)

	assert.Equal(t, "Limits the size of internal in-memory tables. Also set max_heap_table_size. Default is 16777216 (16M)",
		config.MySQL.TmpTableSize.Description)
	assert.Equal(t, 16777216, config.MySQL.TmpTableSize.Example)
	assert.Equal(t, 1073741824, config.MySQL.TmpTableSize.Maximum)
	assert.Equal(t, 1048576, config.MySQL.TmpTableSize.Minimum)
	assert.False(t, config.MySQL.TmpTableSize.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.TmpTableSize.Type)

	assert.Equal(t, "The number of seconds the server waits for activity on a noninteractive connection before closing it.",
		config.MySQL.WaitTimeout.Description)
	assert.Equal(t, 28800, config.MySQL.WaitTimeout.Example)
	assert.Equal(t, 2147483, config.MySQL.WaitTimeout.Maximum)
	assert.Equal(t, 1, config.MySQL.WaitTimeout.Minimum)
	assert.False(t, config.MySQL.WaitTimeout.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.WaitTimeout.Type)

	assert.Equal(t, "The minimum amount of time in seconds to keep binlog entries before deletion. This may be extended for services that require binlog entries for longer than the default for example if using the MySQL Debezium Kafka connector.",
		config.BinlogRetentionPeriod.Description)
	assert.Equal(t, 600, config.BinlogRetentionPeriod.Example)
	assert.Equal(t, 86400, config.BinlogRetentionPeriod.Maximum)
	assert.Equal(t, 600, config.BinlogRetentionPeriod.Minimum)
	assert.False(t, config.BinlogRetentionPeriod.RequiresRestart)
	assert.Equal(t, "integer", config.BinlogRetentionPeriod.Type)
}
