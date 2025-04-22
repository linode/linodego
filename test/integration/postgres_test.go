package integration

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/google/go-cmp/cmp"
	"github.com/linode/linodego"
)

func TestDatabase_Postgres_Suite(t *testing.T) {
	client, database, teardown, err := setupPostgresDatabase(t, nil, "fixtures/TestDatabase_Postgres_Suite")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	postgresdbs, err := client.ListPostgresDatabases(context.Background(), nil)
	if len(postgresdbs) == 0 {
		t.Errorf("Expected a list of Databases, but got none: %v", err)
	}
	if err != nil {
		t.Errorf("Error listing Databases, expected struct, got error %v", err)
	}
	success := false
	for _, db := range postgresdbs {
		if db.ID == database.ID {
			success = true
		}
	}
	if !success {
		t.Error("database not in database list")
	}

	db, err := client.GetPostgresDatabase(context.Background(), database.ID)
	if err != nil {
		t.Errorf("Error viewing postgres database: %v", err)
	}
	if db.ID != database.ID {
		t.Errorf("got wrong db from GetPostgresDatabase: %v", db)
	}

	updatedWindow := linodego.DatabaseMaintenanceWindow{
		DayOfWeek: linodego.DatabaseMaintenanceDayWednesday,
		Duration:  4,
		Frequency: linodego.DatabaseMaintenanceFrequencyWeekly,
		HourOfDay: 8,
		Pending:   []linodego.DatabaseMaintenanceWindowPending{},
	}

	allowList := []string{"128.173.205.21", "123.177.200.20"}

	updatedLabel := database.Label + "-updated"
	opts := linodego.PostgresUpdateOptions{
		AllowList: &allowList,
		Label:     updatedLabel,
		Updates:   &updatedWindow,
	}
	db, err = client.UpdatePostgresDatabase(context.Background(), database.ID, opts)
	if err != nil {
		t.Errorf("failed to update db %d: %v", database.ID, err)
	}

	waitForDatabaseUpdated(t, client, db.ID, linodego.DatabaseEngineTypePostgres, db.Created)

	if db.ID != database.ID {
		t.Errorf("updated db does not match original id")
	}
	if db.Label != updatedLabel {
		t.Errorf("label not updated for db")
	}

	if !reflect.DeepEqual(db.Updates, updatedWindow) {
		t.Errorf("db maintenance window does not match update opts: %v", cmp.Diff(db.Updates, updatedWindow))
	}

	ssl, err := client.GetPostgresDatabaseSSL(context.Background(), database.ID)
	if err != nil {
		t.Errorf("failed to get ssl cert for db: %v", err)
	}
	if ssl == nil {
		t.Error("failed to get ssl cert for db")
	}

	creds, err := client.GetPostgresDatabaseCredentials(context.Background(), database.ID)
	if err != nil {
		t.Errorf("failed to get credentials for db: %v", err)
	}

	if testingMode == recorder.ModeRecording {
		time.Sleep(time.Minute * 5)
	}

	err = client.ResetPostgresDatabaseCredentials(context.Background(), database.ID)
	if err != nil {
		t.Errorf("failed to reset credentials for db: %v", err)
	}

	if testingMode == recorder.ModeRecording {
		time.Sleep(time.Second * 15)
	}

	newcreds, err := client.GetPostgresDatabaseCredentials(context.Background(), database.ID)
	if err != nil {
		t.Errorf("failed to get new credentials for db: %v", err)
	}
	if creds.Password == newcreds.Password {
		t.Error("credentials have not changed for db")
	}

	if testingMode == recorder.ModeRecording {
		time.Sleep(time.Minute * 5)
	}

	if err := client.PatchPostgresDatabase(context.Background(), database.ID); err != nil {
		t.Fatalf("failed to patch database: %s", err)
	}

	// Wait for the DB to enter updating status
	if err := client.WaitForDatabaseStatus(
		context.Background(), database.ID, linodego.DatabaseEngineTypePostgres,
		linodego.DatabaseStatusUpdating, 240); err != nil {
		t.Fatalf("failed to wait for database updating: %s", err)
	}

	// Wait for the DB to re-enter active status
	if err := client.WaitForDatabaseStatus(
		context.Background(), database.ID, linodego.DatabaseEngineTypePostgres,
		linodego.DatabaseStatusActive, 2400); err != nil {
		t.Fatalf("failed to wait for database updating: %s", err)
	}

	if err := client.SuspendPostgresDatabase(context.Background(), database.ID); err != nil {
		t.Fatalf("failed to suspend database: %s", err)
	}

	// Wait for the DB to enter suspended status
	if err := client.WaitForDatabaseStatus(
		context.Background(), database.ID, linodego.DatabaseEngineTypePostgres,
		linodego.DatabaseStatusSuspended, 2400); err != nil {
		t.Fatalf("failed to wait for database suspended: %s", err)
	}

	if err := client.ResumePostgresDatabase(context.Background(), database.ID); err != nil {
		t.Fatalf("failed to resume database: %s", err)
	}

	// Wait for the DB to re-enter active status
	if err := client.WaitForDatabaseStatus(
		context.Background(), database.ID, linodego.DatabaseEngineTypePostgres,
		linodego.DatabaseStatusActive, 2400); err != nil {
		t.Fatalf("failed to wait for database active: %s", err)
	}
}

type postgresDatabaseModifier func(options *linodego.PostgresCreateOptions)

func createPostgresDatabase(t *testing.T, client *linodego.Client,
	databaseMofidiers []postgresDatabaseModifier,
) (*linodego.PostgresDatabase, func(), error) {
	t.Helper()

	createOpts := linodego.PostgresCreateOptions{
		Label:       "go-postgres-testing-def" + randLabel(),
		Region:      getRegionsWithCaps(t, client, []string{"Managed Databases"})[0],
		Type:        "g6-nanode-1",
		Engine:      "postgresql/14",
		ClusterSize: 3,
		AllowList:   []string{"203.0.113.1", "192.0.1.0/24"},
	}
	for _, modifier := range databaseMofidiers {
		modifier(&createOpts)
	}

	database, err := client.CreatePostgresDatabase(context.Background(), createOpts)
	if err != nil {
		t.Fatalf("failed to create database: %s", err)
	}

	// We should retry on db cleanup
	teardown := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
		defer cancel()

		ticker := time.NewTicker(client.GetPollDelay())
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := client.DeletePostgresDatabase(ctx, database.ID)

				if err == nil {
					return
				}

				if lErr, ok := err.(*linodego.Error); ok && lErr.Code == 500 {
					continue
				}

				t.Fatalf("failed to delete database: %s", err)

			case <-ctx.Done():
				t.Fatalf("failed to retry database deletion: %s", ctx.Err())
			}
		}
	}
	return database, teardown, nil
}

func setupPostgresDatabase(t *testing.T, databaseMofidiers []postgresDatabaseModifier,
	fixturesYaml string,
) (*linodego.Client, *linodego.PostgresDatabase, func(), error) {
	t.Helper()
	now := time.Now()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	database, databaseTeardown, err := createPostgresDatabase(t, client, databaseMofidiers)
	if err != nil {
		t.Fatalf("failed to create db: %s", err)
	}

	// After creating the database using the modifier
	fmt.Printf("Database: %+v\n", database)
	fmt.Printf("Database Region: %+v\n", database.Region)
	fmt.Printf("Database Type: %+v\n", database.Type)
	fmt.Printf("Database Engine: %+v\n", database.Engine)
	fmt.Printf("Database EngineConfig: %+v\n", database.EngineConfig)

	// Print PG Engine Config
	fmt.Printf("PG EngineConfig: %+v\n", *database.EngineConfig.PG)
	fmt.Printf("AutovacuumAnalyzeScaleFactor: %+v\n", *database.EngineConfig.PG.AutovacuumAnalyzeScaleFactor)
	fmt.Printf("AutovacuumAnalyzeThreshold: %+v\n", *database.EngineConfig.PG.AutovacuumAnalyzeThreshold)
	fmt.Printf("AutovacuumFreezeMaxAge: %+v\n", *database.EngineConfig.PG.AutovacuumFreezeMaxAge)
	fmt.Printf("AutovacuumMaxWorkers: %+v\n", *database.EngineConfig.PG.AutovacuumMaxWorkers)
	fmt.Printf("AutovacuumNaptime: %+v\n", *database.EngineConfig.PG.AutovacuumNaptime)
	fmt.Printf("AutovacuumVacuumCostDelay: %+v\n", *database.EngineConfig.PG.AutovacuumVacuumCostDelay)
	fmt.Printf("AutovacuumVacuumCostLimit: %+v\n", *database.EngineConfig.PG.AutovacuumVacuumCostLimit)
	fmt.Printf("AutovacuumVacuumScaleFactor: %+v\n", *database.EngineConfig.PG.AutovacuumVacuumScaleFactor)
	fmt.Printf("AutovacuumVacuumThreshold: %+v\n", *database.EngineConfig.PG.AutovacuumVacuumThreshold)
	fmt.Printf("BGWriterDelay: %+v\n", *database.EngineConfig.PG.BGWriterDelay)
	fmt.Printf("BGWriterFlushAfter: %+v\n", *database.EngineConfig.PG.BGWriterFlushAfter)
	fmt.Printf("BGWriterLRUMaxPages: %+v\n", *database.EngineConfig.PG.BGWriterLRUMaxPages)
	fmt.Printf("BGWriterLRUMultiplier: %+v\n", *database.EngineConfig.PG.BGWriterLRUMultiplier)
	fmt.Printf("DeadlockTimeout: %+v\n", *database.EngineConfig.PG.DeadlockTimeout)
	fmt.Printf("DefaultToastCompression: %+v\n", *database.EngineConfig.PG.DefaultToastCompression)
	fmt.Printf("IdleInTransactionSessionTimeout: %+v\n", *database.EngineConfig.PG.IdleInTransactionSessionTimeout)
	fmt.Printf("JIT: %+v\n", *database.EngineConfig.PG.JIT)
	fmt.Printf("MaxFilesPerProcess: %+v\n", *database.EngineConfig.PG.MaxFilesPerProcess)
	fmt.Printf("MaxLocksPerTransaction: %+v\n", *database.EngineConfig.PG.MaxLocksPerTransaction)
	fmt.Printf("MaxLogicalReplicationWorkers: %+v\n", *database.EngineConfig.PG.MaxLogicalReplicationWorkers)
	fmt.Printf("MaxParallelWorkers: %+v\n", *database.EngineConfig.PG.MaxParallelWorkers)
	fmt.Printf("MaxParallelWorkersPerGather: %+v\n", *database.EngineConfig.PG.MaxParallelWorkersPerGather)
	fmt.Printf("MaxPredLocksPerTransaction: %+v\n", *database.EngineConfig.PG.MaxPredLocksPerTransaction)
	fmt.Printf("MaxWALSenders: %+v\n", *database.EngineConfig.PG.MaxWALSenders)
	fmt.Printf("MaxWorkerProcesses: %+v\n", *database.EngineConfig.PG.MaxWorkerProcesses)
	fmt.Printf("PasswordEncryption: %+v\n", *database.EngineConfig.PG.PasswordEncryption)
	fmt.Printf("PGPartmanBGWInterval: %+v\n", *database.EngineConfig.PG.PGPartmanBGWInterval)
	fmt.Printf("PGPartmanBGWRole: %+v\n", *database.EngineConfig.PG.PGPartmanBGWRole)
	fmt.Printf("PGStatMonitorPGSMEnableQueryPlan: %+v\n", *database.EngineConfig.PG.PGStatMonitorPGSMEnableQueryPlan)
	fmt.Printf("PGStatMonitorPGSMMaxBuckets: %+v\n", *database.EngineConfig.PG.PGStatMonitorPGSMMaxBuckets)
	fmt.Printf("PGStatStatementsTrack: %+v\n", *database.EngineConfig.PG.PGStatStatementsTrack)
	fmt.Printf("TempFileLimit: %+v\n", *database.EngineConfig.PG.TempFileLimit)
	fmt.Printf("Timezone: %+v\n", *database.EngineConfig.PG.Timezone)
	fmt.Printf("TrackActivityQuerySize: %+v\n", *database.EngineConfig.PG.TrackActivityQuerySize)
	fmt.Printf("TrackCommitTimestamp: %+v\n", *database.EngineConfig.PG.TrackCommitTimestamp)
	fmt.Printf("TrackIOTiming: %+v\n", *database.EngineConfig.PG.TrackIOTiming)
	fmt.Printf("WALSenderTimeout: %+v\n", *database.EngineConfig.PG.WALSenderTimeout)

	// Removed print statements for the following parameters:
	// LogAutovacuumMinDuration
	// LogErrorVerbosity
	// LogLinePrefix
	// LogMinDurationStatement
	// LogTempFiles
	// MaxPreparedTransactions
	// MaxReplicationSlots
	// MaxStackDepth
	// TrackFunctions
	// WALWriterDelay
	// SynchronousReplication
	// WorkMem

	// Print additional engine config
	fmt.Printf("PGStatMonitorEnable: %+v\n", *database.EngineConfig.PGStatMonitorEnable)
	fmt.Printf("PGLookout: %+v\n", *database.EngineConfig.PGLookout.MaxFailoverReplicationTimeLag)
	//fmt.Printf("ServiceLog: %+v\n", *database.EngineConfig.ServiceLog)
	fmt.Printf("SharedBuffersPercentage: %+v\n", *database.EngineConfig.SharedBuffersPercentage)
	fmt.Printf("WorkMem: %+v\n", *database.EngineConfig.WorkMem)

	_, err = client.WaitForEventFinished(context.Background(), database.ID, linodego.EntityDatabase,
		linodego.ActionDatabaseCreate, now, 5400)
	if err != nil {
		t.Fatalf("failed to wait for db create event: %s", err)
	}

	teardown := func() {
		databaseTeardown()
		fixtureTeardown()
	}
	return client, database, teardown, err
}
