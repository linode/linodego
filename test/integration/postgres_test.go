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

const (
	testPostgresBackupLabel = "reallycoolbackup"
	testPostgresDBLabel     = "linodego-test-postgres-database"
)

var testPostgresCreateOpts = linodego.PostgresCreateOptions{
	Label:           testPostgresDBLabel,
	Region:          "us-east",
	Type:            "g6-nanode-1",
	Engine:          "postgresql/10.14",
	Encrypted:       false,
	SSLConnection:   false,
	ClusterSize:     3,
	ReplicationType: linodego.PostgresReplicationAsynch,
	AllowList:       []string{"203.0.113.1", "192.0.1.0/24"},
}

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

	week := 3

	updatedWindow := linodego.DatabaseMaintenanceWindow{
		DayOfWeek:   linodego.DatabaseMaintenanceDayWednesday,
		Duration:    1,
		Frequency:   linodego.DatabaseMaintenanceFrequencyMonthly,
		HourOfDay:   8,
		WeekOfMonth: &week,
	}

	opts := linodego.PostgresUpdateOptions{
		AllowList: []string{"128.173.205.21", "123.177.200.20"},
		Label:     fmt.Sprintf("%s-updated", database.Label),
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
	if db.Label != fmt.Sprintf("%s-updated", database.Label) {
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

	backupOptions := linodego.PostgresBackupCreateOptions{
		Label:  testPostgresBackupLabel,
		Target: linodego.PostgresDatabaseTargetPrimary,
	}

	if err := client.CreatePostgresDatabaseBackup(context.Background(), database.ID, backupOptions); err != nil {
		t.Errorf("failed to create db backup: %v", err)
	}

	backup, err := client.WaitForPostgresDatabaseBackup(context.Background(), database.ID, testPostgresBackupLabel, 1200)
	if err != nil {
		t.Fatalf("failed to wait for backup: %s", err)
	}

	if backup.Label != testPostgresBackupLabel {
		t.Fatalf("backup label mismatch: %v != %v", testPostgresBackupLabel, backup.Label)
	}

	backup, err = client.GetPostgresDatabaseBackup(context.Background(), database.ID, backup.ID)
	if err != nil {
		t.Errorf("failed to get backup %d for db: %v", backup.ID, err)
	}

	if backup.Label != testPostgresBackupLabel {
		t.Fatalf("backup label mismatch: %v != %v", testPostgresBackupLabel, backup.Label)
	}

	// Wait for the DB to re-enter active status before final deletion
	if err := client.WaitForDatabaseStatus(
		context.Background(), database.ID, linodego.DatabaseEngineTypePostgres,
		linodego.DatabaseStatusActive, 2400); err != nil {
		t.Fatalf("failed to wait for database updating: %s", err)
	}
}

type postgresDatabaseModifier func(options *linodego.PostgresCreateOptions)

func createPostgresDatabase(t *testing.T, client *linodego.Client,
	databaseMofidiers []postgresDatabaseModifier,
) (*linodego.PostgresDatabase, func(), error) {
	t.Helper()

	createOpts := testPostgresCreateOpts
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

		ticker := time.NewTicker(client.GetPollDelay() * time.Millisecond)
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

	_, err = client.WaitForEventFinished(context.Background(), database.ID, linodego.EntityDatabase,
		linodego.ActionDatabaseCreate, now, 3600)
	if err != nil {
		t.Fatalf("failed to wait for db create event: %s", err)
	}

	teardown := func() {
		databaseTeardown()
		fixtureTeardown()
	}
	return client, database, teardown, err
}
