package integration

import (
	"context"
	"os"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseMySQL_EngineConfig_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestDatabaseMySQL_EngineConfig_Get")
	defer teardown()

	config, err := client.GetMySQLDatabaseConfig(context.Background())
	if err != nil {
		t.Fatalf("Error getting MySQL database config: %v", err)
	}

	assert.IsType(t, string("The number of seconds that the mysqld server waits for a connect packet before responding with Bad handshakes"),
		config.MySQL.ConnectTimeout.Description)
	assert.IsType(t, int(10), config.MySQL.ConnectTimeout.Example)
	assert.IsType(t, int(3600), config.MySQL.ConnectTimeout.Maximum)
	assert.IsType(t, int(2), config.MySQL.ConnectTimeout.Minimum)
	assert.IsType(t, int(0), config.MySQL.ConnectTimeout.Minimum)
	assert.IsType(t, bool(false), config.MySQL.ConnectTimeout.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.ConnectTimeout.Type)

	assert.IsType(t, string("Default server time zone as an offset from UTC (from -12:00 to +12:00), a time zone name, or 'SYSTEM' to use the MySQL server default."),
		config.MySQL.DefaultTimeZone.Description)
	assert.IsType(t, string("+03:00"), config.MySQL.DefaultTimeZone.Example)
	assert.IsType(t, int(100), config.MySQL.DefaultTimeZone.MaxLength)
	assert.IsType(t, int(2), config.MySQL.DefaultTimeZone.MinLength)
	assert.IsType(t, string("^([-+][\\d:]*|[\\w/]*)$"), config.MySQL.DefaultTimeZone.Pattern)
	assert.IsType(t, bool(false), config.MySQL.DefaultTimeZone.RequiresRestart)
	assert.Equal(t, "string", config.MySQL.DefaultTimeZone.Type) // keep this if the type value must be "string"

	assert.IsType(t, string("The maximum permitted result length in bytes for the GROUP_CONCAT() function."), config.MySQL.GroupConcatMaxLen.Description)
	assert.IsType(t, float64(1024), config.MySQL.GroupConcatMaxLen.Example)
	assert.IsType(t, float64(18446744073709551600), config.MySQL.GroupConcatMaxLen.Maximum)
	assert.IsType(t, float64(4), config.MySQL.GroupConcatMaxLen.Minimum)
	assert.IsType(t, bool(false), config.MySQL.GroupConcatMaxLen.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.GroupConcatMaxLen.Type)

	assert.IsType(t, string("The time, in seconds, before cached statistics expire"), config.MySQL.InformationSchemaStatsExpiry.Description)
	assert.IsType(t, int(86400), config.MySQL.InformationSchemaStatsExpiry.Example)
	assert.IsType(t, int(31536000), config.MySQL.InformationSchemaStatsExpiry.Maximum)
	assert.IsType(t, int(900), config.MySQL.InformationSchemaStatsExpiry.Minimum)
	assert.IsType(t, bool(false), config.MySQL.InformationSchemaStatsExpiry.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InformationSchemaStatsExpiry.Type)

	assert.IsType(t, string("Maximum size for the InnoDB change buffer, as a percentage of the total size of the buffer pool. Default is 25"), config.MySQL.InnoDBChangeBufferMaxSize.Description)
	assert.IsType(t, int(30), config.MySQL.InnoDBChangeBufferMaxSize.Example)
	assert.IsType(t, int(50), config.MySQL.InnoDBChangeBufferMaxSize.Maximum)
	assert.IsType(t, int(0), config.MySQL.InnoDBChangeBufferMaxSize.Minimum)
	assert.IsType(t, bool(false), config.MySQL.InnoDBChangeBufferMaxSize.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBChangeBufferMaxSize.Type)

	assert.IsType(t, string("Specifies whether flushing a page from the InnoDB buffer pool also flushes other dirty pages in the same extent (default is 1): 0 - dirty pages in the same extent are not flushed, 1 - flush contiguous dirty pages in the same extent, 2 - flush dirty pages in the same extent"), config.MySQL.InnoDBFlushNeighbors.Description)
	assert.IsType(t, int(0), config.MySQL.InnoDBFlushNeighbors.Example)
	assert.IsType(t, int(2), config.MySQL.InnoDBFlushNeighbors.Maximum)
	assert.IsType(t, int(0), config.MySQL.InnoDBFlushNeighbors.Minimum)
	assert.IsType(t, bool(false), config.MySQL.InnoDBFlushNeighbors.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBFlushNeighbors.Type)

	assert.IsType(t, string("Minimum length of words that are stored in an InnoDB FULLTEXT index. Changing this parameter will lead to a restart of the MySQL service."), config.MySQL.InnoDBFTMinTokenSize.Description)
	assert.IsType(t, int(3), config.MySQL.InnoDBFTMinTokenSize.Example)
	assert.IsType(t, int(16), config.MySQL.InnoDBFTMinTokenSize.Maximum)
	assert.IsType(t, int(0), config.MySQL.InnoDBFTMinTokenSize.Minimum)
	assert.IsType(t, bool(true), config.MySQL.InnoDBFTMinTokenSize.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBFTMinTokenSize.Type)

	assert.IsType(t, string("This option is used to specify your own InnoDB FULLTEXT index stopword list for all InnoDB tables."), config.MySQL.InnoDBFTServerStopwordTable.Description)
	assert.IsType(t, string("db_name/table_name"), config.MySQL.InnoDBFTServerStopwordTable.Example)
	assert.IsType(t, int(1024), config.MySQL.InnoDBFTServerStopwordTable.MaxLength)
	assert.IsType(t, string("^.+/.+$"), config.MySQL.InnoDBFTServerStopwordTable.Pattern)
	assert.IsType(t, bool(false), config.MySQL.InnoDBFTServerStopwordTable.RequiresRestart)
	assert.Equal(t, []string{"null", "string"}, config.MySQL.InnoDBFTServerStopwordTable.Type)

	assert.IsType(t, string("The length of time in seconds an InnoDB transaction waits for a row lock before giving up. Default is 120."), config.MySQL.InnoDBLockWaitTimeout.Description)
	assert.IsType(t, int(50), config.MySQL.InnoDBLockWaitTimeout.Example)
	assert.IsType(t, int(3600), config.MySQL.InnoDBLockWaitTimeout.Maximum)
	assert.IsType(t, int(1), config.MySQL.InnoDBLockWaitTimeout.Minimum)
	assert.IsType(t, bool(false), config.MySQL.InnoDBLockWaitTimeout.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBLockWaitTimeout.Type)

	assert.IsType(t, string("The size in bytes of the buffer that InnoDB uses to write to the log files on disk."), config.MySQL.InnoDBLogBufferSize.Description)
	assert.IsType(t, int(16777216), config.MySQL.InnoDBLogBufferSize.Example)
	assert.IsType(t, int(4294967295), config.MySQL.InnoDBLogBufferSize.Maximum)
	assert.IsType(t, int(1048576), config.MySQL.InnoDBLogBufferSize.Minimum)
	assert.IsType(t, bool(false), config.MySQL.InnoDBLogBufferSize.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBLogBufferSize.Type)

	assert.IsType(t, string("The upper limit in bytes on the size of the temporary log files used during online DDL operations for InnoDB tables."), config.MySQL.InnoDBOnlineAlterLogMaxSize.Description)
	assert.IsType(t, int(134217728), config.MySQL.InnoDBOnlineAlterLogMaxSize.Example)
	assert.IsType(t, int(1099511627776), config.MySQL.InnoDBOnlineAlterLogMaxSize.Maximum)
	assert.IsType(t, int(65536), config.MySQL.InnoDBOnlineAlterLogMaxSize.Minimum)
	assert.IsType(t, bool(false), config.MySQL.InnoDBOnlineAlterLogMaxSize.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBOnlineAlterLogMaxSize.Type)

	assert.IsType(t, string("The number of I/O threads for read operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service."), config.MySQL.InnoDBReadIOThreads.Description)
	assert.IsType(t, int(10), config.MySQL.InnoDBReadIOThreads.Example)
	assert.IsType(t, int(64), config.MySQL.InnoDBReadIOThreads.Maximum)
	assert.IsType(t, int(1), config.MySQL.InnoDBReadIOThreads.Minimum)
	assert.IsType(t, bool(true), config.MySQL.InnoDBReadIOThreads.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBReadIOThreads.Type)

	assert.IsType(t, string("When enabled a transaction timeout causes InnoDB to abort and roll back the entire transaction. Changing this parameter will lead to a restart of the MySQL service."), config.MySQL.InnoDBRollbackOnTimeout.Description)
	assert.IsType(t, bool(true), config.MySQL.InnoDBRollbackOnTimeout.Example)
	assert.IsType(t, bool(true), config.MySQL.InnoDBRollbackOnTimeout.RequiresRestart)
	assert.Equal(t, "boolean", config.MySQL.InnoDBRollbackOnTimeout.Type)

	assert.IsType(t, string("Defines the maximum number of threads permitted inside of InnoDB. Default is 0 (infinite concurrency - no limit)"), config.MySQL.InnoDBThreadConcurrency.Description)
	assert.IsType(t, int(10), config.MySQL.InnoDBThreadConcurrency.Example)
	assert.IsType(t, int(1000), config.MySQL.InnoDBThreadConcurrency.Maximum)
	assert.IsType(t, int(0), config.MySQL.InnoDBThreadConcurrency.Minimum)
	assert.IsType(t, bool(false), config.MySQL.InnoDBThreadConcurrency.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBThreadConcurrency.Type)

	assert.IsType(t, string("The number of I/O threads for write operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service."), config.MySQL.InnoDBWriteIOThreads.Description)
	assert.IsType(t, int(10), config.MySQL.InnoDBWriteIOThreads.Example)
	assert.IsType(t, int(64), config.MySQL.InnoDBWriteIOThreads.Maximum)
	assert.IsType(t, int(1), config.MySQL.InnoDBWriteIOThreads.Minimum)
	assert.IsType(t, bool(true), config.MySQL.InnoDBWriteIOThreads.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InnoDBWriteIOThreads.Type)

	assert.IsType(t, string("The number of seconds the server waits for activity on an interactive connection before closing it."), config.MySQL.InteractiveTimeout.Description)
	assert.IsType(t, int(3600), config.MySQL.InteractiveTimeout.Example)
	assert.IsType(t, int(604800), config.MySQL.InteractiveTimeout.Maximum)
	assert.IsType(t, int(30), config.MySQL.InteractiveTimeout.Minimum)
	assert.IsType(t, bool(false), config.MySQL.InteractiveTimeout.RequiresRestart)
	assert.Equal(t, "integer", config.MySQL.InteractiveTimeout.Type)

	assert.IsType(t, string("The storage engine for in-memory internal temporary tables."), config.MySQL.InternalTmpMemStorageEngine.Description)
	assert.IsType(t, string("TempTable"), config.MySQL.InternalTmpMemStorageEngine.Example)
	assert.Equal(t, []string{"TempTable", "MEMORY"}, config.MySQL.InternalTmpMemStorageEngine.Enum)
	assert.IsType(t, bool(false), config.MySQL.InternalTmpMemStorageEngine.RequiresRestart)
	assert.Equal(t, "string", config.MySQL.InternalTmpMemStorageEngine.Type)
}

func TestDatabaseMySQL_EngineConfig_Suite(t *testing.T) {
	databaseModifiers := []mysqlDatabaseModifier{
		createMySQLOptionsModifier(),
	}

	client, database, teardown, err := setupMySQLDatabase(
		t,
		databaseModifiers,
		"fixtures/TestDatabaseMySQL_EngineConfig_Suite",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	assertSQLDatabaseBasics(t, database)

	expected := newExpectedSQLEngineConfig()
	assertMySQLEngineConfigEqual(t, database.EngineConfig.MySQL, expected)
	assert.Equal(t, 600, *database.EngineConfig.BinlogRetentionPeriod)
	assert.ElementsMatch(t, []string{"192.0.1.0/24", "203.0.113.1/32"}, database.AllowList)

	fetchedDB, err := client.GetMySQLDatabase(context.Background(), database.ID)
	assert.NoError(t, err)
	assertMySQLEngineConfigEqual(t, fetchedDB.EngineConfig.MySQL, expected)

	updateOptions := linodego.MySQLUpdateOptions{
		Label: "db-engine-config-updated",
		EngineConfig: &linodego.MySQLDatabaseEngineConfig{
			MySQL: &linodego.MySQLDatabaseEngineConfigMySQL{
				ConnectTimeout:              linodego.Pointer(20),
				InnoDBLockWaitTimeout:       linodego.Pointer(60),
				NetReadTimeout:              linodego.Pointer(40),
				InnoDBFTServerStopwordTable: linodego.Pointer[*string](nil),
			},
		},
	}

	updatedDB, err := client.UpdateMySQLDatabase(context.Background(), database.ID, updateOptions)
	if err != nil {
		t.Errorf("failed to update db %d: %v", database.ID, err)
	}

	waitForDatabaseUpdated(t, client, updatedDB.ID, linodego.DatabaseEngineTypeMySQL, updatedDB.Created)

	assertUpdatedSQLFields(t, updatedDB.EngineConfig.MySQL)
}

func TestDatabaseMySQL_EngineConfig_Create_NullableFieldAsNilValue(t *testing.T) {
	databaseModifiers := []mysqlDatabaseModifier{
		createMySQLOptionsModifierNullableField(),
	}

	client, fixtureTeardown := createTestClient(t, "fixtures/TestDatabaseMySQL_EngineConfig_Create_NullableFieldAsNilValue")

	database, databaseTeardown, err := createMySQLDatabase(t, client, databaseModifiers)
	if err != nil {
		t.Fatalf("failed to create db: %s", err)
	}

	defer func() {
		databaseTeardown()
		fixtureTeardown()
	}()

	assert.Nil(t, database.EngineConfig.MySQL.InnoDBFTServerStopwordTable)
}

func TestDatabaseMySQL_EngineConfig_Create_Fails_EmptyDoublePointerValue(t *testing.T) {
	if os.Getenv("LINODE_FIXTURE_MODE") == "play" {
		t.Skip("Skipping negative test scenario: LINODE_FIXTURE_MODE is 'play'")
	}

	invalidRequestData := linodego.MySQLCreateOptions{
		Label:  "db-with-engine-config",
		Region: "us-east",
		Type:   "g6-dedicated-2",
		Engine: "mysql/8",
		EngineConfig: &linodego.MySQLDatabaseEngineConfig{
			MySQL: &linodego.MySQLDatabaseEngineConfigMySQL{
				InnoDBFTServerStopwordTable: DoublePointer(linodego.Pointer("")),
			},
		},
	}

	client, _ := createTestClient(t, "")

	_, err := client.CreateMySQLDatabase(context.Background(), invalidRequestData)

	assert.Contains(t, err.Error(), "Invalid format: must match pattern ^.+/.+$")
}

func DoublePointer[T any](v *T) **T {
	return &v
}

func createMySQLOptionsModifierNullableField() mysqlDatabaseModifier {
	return func(options *linodego.MySQLCreateOptions) {
		options.Label = "example-db-created-with-config"
		options.Region = "us-east"
		options.Type = "g6-dedicated-2"
		options.Engine = "mysql/8"
		options.EngineConfig = &linodego.MySQLDatabaseEngineConfig{
			MySQL: &linodego.MySQLDatabaseEngineConfigMySQL{
				InnoDBFTServerStopwordTable: nil,
			},
		}
	}
}

func createMySQLOptionsModifier() mysqlDatabaseModifier {
	return func(options *linodego.MySQLCreateOptions) {
		options.Label = "example-db-created-with-config"
		options.Region = "us-east"
		options.Type = "g6-dedicated-2"
		options.Engine = "mysql/8"
		options.EngineConfig = &linodego.MySQLDatabaseEngineConfig{
			MySQL: &linodego.MySQLDatabaseEngineConfigMySQL{
				ConnectTimeout:               linodego.Pointer(10),
				DefaultTimeZone:              linodego.Pointer("+03:00"),
				GroupConcatMaxLen:            linodego.Pointer(1024.0),
				InformationSchemaStatsExpiry: linodego.Pointer(86400),
				InnoDBChangeBufferMaxSize:    linodego.Pointer(30),
				InnoDBFlushNeighbors:         linodego.Pointer(1),
				InnoDBFTMinTokenSize:         linodego.Pointer(3),
				InnoDBFTServerStopwordTable:  DoublePointer(linodego.Pointer("mydb/stopwords")),
				InnoDBLockWaitTimeout:        linodego.Pointer(50),
				InnoDBLogBufferSize:          linodego.Pointer(16777216),
				InnoDBOnlineAlterLogMaxSize:  linodego.Pointer(134217728),
				InnoDBReadIOThreads:          linodego.Pointer(10),
				InnoDBRollbackOnTimeout:      linodego.Pointer(true),
				InnoDBThreadConcurrency:      linodego.Pointer(10),
				InnoDBWriteIOThreads:         linodego.Pointer(10),
				InteractiveTimeout:           linodego.Pointer(3600),
				InternalTmpMemStorageEngine:  linodego.Pointer("TempTable"),
				MaxAllowedPacket:             linodego.Pointer(67108864),
				MaxHeapTableSize:             linodego.Pointer(16777216),
				NetBufferLength:              linodego.Pointer(16384),
				NetReadTimeout:               linodego.Pointer(30),
				NetWriteTimeout:              linodego.Pointer(30),
				SortBufferSize:               linodego.Pointer(262144),
				SQLMode:                      linodego.Pointer("TRADITIONAL"),
				SQLRequirePrimaryKey:         linodego.Pointer(true),
				TmpTableSize:                 linodego.Pointer(16777216),
				WaitTimeout:                  linodego.Pointer(28800),
			},
			BinlogRetentionPeriod: linodego.Pointer(600),
		}
	}
}

func assertSQLDatabaseBasics(t *testing.T, db *linodego.MySQLDatabase) {
	assert.IsType(t, int(16494), db.ID)
	assert.Equal(t, "example-db-created-with-config", db.Label)
	assert.Equal(t, "us-east", db.Region)
	assert.Equal(t, "mysql", db.Engine)
	assert.Equal(t, "8", db.Version)
	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, linodego.DatabasePlatform("rdbms-default"), db.Platform)
	assert.Equal(t, "g6-dedicated-2", db.Type)
	assert.IsType(t, int(58), db.TotalDiskSizeGB)
	assert.Equal(t, 0, db.UsedDiskSizeGB)
	assert.IsType(t, int(25698), db.Port)
	assert.NotNil(t, db.EngineConfig)
	assert.IsType(t, &linodego.MySQLDatabaseEngineConfigMySQL{}, db.EngineConfig.MySQL)
}

func newExpectedSQLEngineConfig() map[string]any {
	return map[string]any{
		"ConnectTimeout":               10,
		"DefaultTimeZone":              "+03:00",
		"GroupConcatMaxLen":            1024.0,
		"InformationSchemaStatsExpiry": 86400,
		"InnoDBChangeBufferMaxSize":    30,
		"InnoDBFlushNeighbors":         1,
		"InnoDBFTMinTokenSize":         3,
		"InnoDBFTServerStopwordTable":  "mydb/stopwords",
		"InnoDBLockWaitTimeout":        50,
		"InnoDBLogBufferSize":          16777216,
		"InnoDBOnlineAlterLogMaxSize":  134217728,
		"InnoDBReadIOThreads":          10,
		"InnoDBRollbackOnTimeout":      true,
		"InnoDBThreadConcurrency":      10,
		"InnoDBWriteIOThreads":         10,
		"InteractiveTimeout":           3600,
		"InternalTmpMemStorageEngine":  "TempTable",
		"MaxAllowedPacket":             67108864,
		"MaxHeapTableSize":             16777216,
		"NetBufferLength":              16384,
		"NetReadTimeout":               30,
		"NetWriteTimeout":              30,
		"SortBufferSize":               262144,
		"SQLMode":                      "TRADITIONAL",
		"SQLRequirePrimaryKey":         true,
		"TmpTableSize":                 16777216,
		"WaitTimeout":                  28800,
	}
}

func assertMySQLEngineConfigEqual(t *testing.T, cfg *linodego.MySQLDatabaseEngineConfigMySQL, expected map[string]any) {
	assert.Equal(t, expected["ConnectTimeout"], *cfg.ConnectTimeout)
	assert.Equal(t, expected["DefaultTimeZone"], *cfg.DefaultTimeZone)
	assert.Equal(t, expected["GroupConcatMaxLen"], *cfg.GroupConcatMaxLen)
	assert.Equal(t, expected["InformationSchemaStatsExpiry"], *cfg.InformationSchemaStatsExpiry)
	assert.Equal(t, expected["InnoDBChangeBufferMaxSize"], *cfg.InnoDBChangeBufferMaxSize)
	assert.Equal(t, expected["InnoDBFlushNeighbors"], *cfg.InnoDBFlushNeighbors)
	assert.Equal(t, expected["InnoDBFTMinTokenSize"], *cfg.InnoDBFTMinTokenSize)
	assert.Equal(t, expected["InnoDBFTServerStopwordTable"], **cfg.InnoDBFTServerStopwordTable)
	assert.Equal(t, expected["InnoDBLockWaitTimeout"], *cfg.InnoDBLockWaitTimeout)
	assert.Equal(t, expected["InnoDBLogBufferSize"], *cfg.InnoDBLogBufferSize)
	assert.Equal(t, expected["InnoDBOnlineAlterLogMaxSize"], *cfg.InnoDBOnlineAlterLogMaxSize)
	assert.Equal(t, expected["InnoDBReadIOThreads"], *cfg.InnoDBReadIOThreads)
	assert.Equal(t, expected["InnoDBRollbackOnTimeout"], *cfg.InnoDBRollbackOnTimeout)
	assert.Equal(t, expected["InnoDBThreadConcurrency"], *cfg.InnoDBThreadConcurrency)
	assert.Equal(t, expected["InnoDBWriteIOThreads"], *cfg.InnoDBWriteIOThreads)
	assert.Equal(t, expected["InteractiveTimeout"], *cfg.InteractiveTimeout)
	assert.Equal(t, expected["InternalTmpMemStorageEngine"], *cfg.InternalTmpMemStorageEngine)
	assert.Equal(t, expected["MaxAllowedPacket"], *cfg.MaxAllowedPacket)
	assert.Equal(t, expected["MaxHeapTableSize"], *cfg.MaxHeapTableSize)
	assert.Equal(t, expected["NetBufferLength"], *cfg.NetBufferLength)
	assert.Equal(t, expected["NetReadTimeout"], *cfg.NetReadTimeout)
	assert.Equal(t, expected["NetWriteTimeout"], *cfg.NetWriteTimeout)
	assert.Equal(t, expected["SortBufferSize"], *cfg.SortBufferSize)
	assert.Equal(t, expected["SQLMode"], *cfg.SQLMode)
	assert.Equal(t, expected["SQLRequirePrimaryKey"], *cfg.SQLRequirePrimaryKey)
	assert.Equal(t, expected["TmpTableSize"], *cfg.TmpTableSize)
	assert.Equal(t, expected["WaitTimeout"], *cfg.WaitTimeout)
}

func assertUpdatedSQLFields(t *testing.T, cfg *linodego.MySQLDatabaseEngineConfigMySQL) {
	assert.Equal(t, 20, *cfg.ConnectTimeout)
	assert.Equal(t, 60, *cfg.InnoDBLockWaitTimeout)
	assert.Equal(t, 40, *cfg.NetReadTimeout)
	assert.Nil(t, cfg.InnoDBFTServerStopwordTable)
}
