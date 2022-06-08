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
	testMongoBackupLabel = "reallycoolbackup"
	testMongoDBLabel     = "linodego-test-mongo-database"
)

var testMongoCreateOpts = linodego.MongoCreateOptions{
	Label:         testMongoDBLabel,
	Region:        "us-east",
	Type:          "g6-nanode-1",
	Engine:        "mongodb/4.4.10",
	Encrypted:     false,
	ClusterSize:   3,
	SSLConnection: false,
	AllowList:     []string{"203.0.113.1", "192.0.1.0/24"},
}

func TestDatabase_Mongo_Suite(t *testing.T) {
	client, database, teardown, err := setupMongoDatabase(t, nil, "fixtures/TestDatabase_Mongo_Suite")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	mongodbs, err := client.ListMongoDatabases(context.Background(), nil)
	if len(mongodbs) == 0 {
		t.Errorf("Expected a list of Databases, but got none: %v", err)
	}
	if err != nil {
		t.Errorf("Error listing Databases, expected struct, got error %v", err)
	}
	success := false
	for _, db := range mongodbs {
		if db.ID == database.ID {
			success = true
		}
	}
	if !success {
		t.Error("database not in database list")
	}

	db, err := client.GetMongoDatabase(context.Background(), database.ID)
	if err != nil {
		t.Errorf("Error viewing mongo database: %v", err)
	}
	if db.ID != database.ID {
		t.Errorf("got wrong db from GetMongoDatabase: %v", db)
	}

	week := 3

	updatedWindow := linodego.DatabaseMaintenanceWindow{
		DayOfWeek:   linodego.DatabaseMaintenanceDayWednesday,
		Duration:    1,
		Frequency:   linodego.DatabaseMaintenanceFrequencyMonthly,
		HourOfDay:   8,
		WeekOfMonth: &week,
	}

	allowList := []string{"128.173.205.21", "123.177.200.20"}

	opts := linodego.MongoUpdateOptions{
		AllowList: &allowList,
		Label:     fmt.Sprintf("%s-updated", database.Label),
		Updates:   &updatedWindow,
	}
	db, err = client.UpdateMongoDatabase(context.Background(), database.ID, opts)
	if err != nil {
		t.Errorf("failed to update db %d: %v", database.ID, err)
	}

	if db.ID != database.ID {
		t.Errorf("updated db does not match original id")
	}
	if db.Label != fmt.Sprintf("%s-updated", database.Label) {
		t.Errorf("label not updated for db")
	}

	if !reflect.DeepEqual(db.Updates, updatedWindow) {
		t.Errorf("db maintenance window does not match update opts: %v", cmp.Diff(db.Updates, updatedWindow))
	}

	ssl, err := client.GetMongoDatabaseSSL(context.Background(), database.ID)
	if err != nil {
		t.Errorf("failed to get ssl cert for db: %v", err)
	}
	if ssl == nil {
		t.Error("failed to get ssl cert for db")
	}

	creds, err := client.GetMongoDatabaseCredentials(context.Background(), database.ID)
	if err != nil {
		t.Errorf("failed to get credentials for db: %v", err)
	}

	if testingMode == recorder.ModeRecording {
		time.Sleep(time.Minute * 5)
	}

	err = client.ResetMongoDatabaseCredentials(context.Background(), database.ID)
	if err != nil {
		t.Errorf("failed to reset credentials for db: %v", err)
	}

	if testingMode == recorder.ModeRecording {
		time.Sleep(time.Second * 15)
	}

	newcreds, err := client.GetMongoDatabaseCredentials(context.Background(), database.ID)
	if err != nil {
		t.Errorf("failed to get new credentials for db: %v", err)
	}
	if creds.Password == newcreds.Password {
		t.Error("credentials have not changed for db")
	}

	if testingMode == recorder.ModeRecording {
		time.Sleep(time.Minute * 5)
	}

	if err := client.PatchMongoDatabase(context.Background(), database.ID); err != nil {
		t.Fatalf("failed to patch database: %s", err)
	}

	// Wait for the DB to enter updating status
	if err := client.WaitForDatabaseStatus(
		context.Background(), database.ID, linodego.DatabaseEngineTypeMongo,
		linodego.DatabaseStatusUpdating, 240); err != nil {
		t.Fatalf("failed to wait for database updating: %s", err)
	}

	// Wait for the DB to re-enter active status
	if err := client.WaitForDatabaseStatus(
		context.Background(), database.ID, linodego.DatabaseEngineTypeMongo,
		linodego.DatabaseStatusActive, 2400); err != nil {
		t.Fatalf("failed to wait for database updating: %s", err)
	}

	backupOptions := linodego.MongoBackupCreateOptions{
		Label:  testMongoBackupLabel,
		Target: linodego.MongoDatabaseTargetPrimary,
	}

	if err := client.CreateMongoDatabaseBackup(context.Background(), database.ID, backupOptions); err != nil {
		t.Errorf("failed to create db backup: %v", err)
	}

	backup, err := client.WaitForMongoDatabaseBackup(context.Background(), database.ID, testMongoBackupLabel, 1200)
	if err != nil {
		t.Fatalf("failed to wait for backup: %s", err)
	}

	if backup.Label != testMongoBackupLabel {
		t.Fatalf("backup label mismatch: %v != %v", testMongoBackupLabel, backup.Label)
	}

	backup, err = client.GetMongoDatabaseBackup(context.Background(), database.ID, backup.ID)
	if err != nil {
		t.Errorf("failed to get backup %d for db: %v", backup.ID, err)
	}

	if backup.Label != testMongoBackupLabel {
		t.Fatalf("backup label mismatch: %v != %v", testMongoBackupLabel, backup.Label)
	}

	// Wait for the DB to re-enter active status before final deletion
	if err := client.WaitForDatabaseStatus(
		context.Background(), database.ID, linodego.DatabaseEngineTypeMongo,
		linodego.DatabaseStatusActive, 2400); err != nil {
		t.Fatalf("failed to wait for database updating: %s", err)
	}
}

type mongoDatabaseModifier func(options *linodego.MongoCreateOptions)

func createMongoDatabase(t *testing.T, client *linodego.Client,
	databaseMofidiers []mongoDatabaseModifier,
) (*linodego.MongoDatabase, func(), error) {
	t.Helper()

	createOpts := testMongoCreateOpts
	for _, modifier := range databaseMofidiers {
		modifier(&createOpts)
	}

	database, err := client.CreateMongoDatabase(context.Background(), createOpts)
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
				err := client.DeleteMongoDatabase(ctx, database.ID)

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

func setupMongoDatabase(t *testing.T, databaseMofidiers []mongoDatabaseModifier,
	fixturesYaml string,
) (*linodego.Client, *linodego.MongoDatabase, func(), error) {
	t.Helper()
	now := time.Now()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	database, databaseTeardown, err := createMongoDatabase(t, client, databaseMofidiers)
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
