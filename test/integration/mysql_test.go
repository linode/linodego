package integration

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/google/go-cmp/cmp"
	"github.com/linode/linodego"
)

func TestDatabase_MySQL_Suite(t *testing.T) {
	client, database, teardown, err := setupMySQLDatabase(t, nil, "fixtures/TestDatabase_MySQL_Suite")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	mysqldbs, err := client.ListMySQLDatabases(context.Background(), nil)
	if len(mysqldbs) == 0 {
		t.Errorf("Expected a list of Databases, but got none: %v", err)
	}
	if err != nil {
		t.Errorf("Error listing Databases, expected struct, got error %v", err)
	}
	success := false
	for _, db := range mysqldbs {
		if db.ID == database.ID {
			success = true
		}
	}
	if !success {
		t.Error("database not in database list")
	}

	db, err := client.GetMySQLDatabase(context.Background(), database.ID)
	if err != nil {
		t.Errorf("Error viewing mysql database: %v", err)
	}
	if db.ID != database.ID {
		t.Errorf("got wrong db from GetMySQLDatabase: %v", db)
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
	opts := linodego.MySQLUpdateOptions{
		AllowList: &allowList,
		Label:     updatedLabel,
		Updates:   &updatedWindow,
	}
	db, err = client.UpdateMySQLDatabase(context.Background(), database.ID, opts)
	if err != nil {
		t.Errorf("failed to update db %d: %v", database.ID, err)
	}

	waitForDatabaseUpdated(t, client, db.ID, linodego.DatabaseEngineTypeMySQL, db.Created)

	if db.ID != database.ID {
		t.Errorf("updated db does not match original id")
	}
	if db.Label != updatedLabel {
		t.Errorf("label not updated for db")
	}

	if !reflect.DeepEqual(db.Updates, updatedWindow) {
		t.Errorf("db maintenance window does not match update opts: %v", cmp.Diff(db.Updates, updatedWindow))
	}

	ssl, err := client.GetMySQLDatabaseSSL(context.Background(), database.ID)
	if err != nil {
		t.Errorf("failed to get ssl cert for db: %v", err)
	}
	if ssl == nil {
		t.Error("failed to get ssl cert for db")
	}

	creds, err := client.GetMySQLDatabaseCredentials(context.Background(), database.ID)
	if err != nil {
		t.Errorf("failed to get credentials for db: %v", err)
	}

	if testingMode == recorder.ModeRecording {
		time.Sleep(time.Minute * 5)
	}

	err = client.ResetMySQLDatabaseCredentials(context.Background(), database.ID)
	if err != nil {
		t.Errorf("failed to reset credentials for db: %v", err)
	}

	if testingMode == recorder.ModeRecording {
		time.Sleep(time.Second * 15)
	}

	newcreds, err := client.GetMySQLDatabaseCredentials(context.Background(), database.ID)
	if err != nil {
		t.Errorf("failed to get new credentials for db: %v", err)
	}
	if creds.Password == newcreds.Password {
		t.Error("credentials have not changed for db")
	}

	if err := client.PatchMySQLDatabase(context.Background(), database.ID); err != nil {
		t.Fatalf("failed to patch database: %s", err)
	}

	// Wait for the DB to enter updating status
	if err := client.WaitForDatabaseStatus(
		context.Background(), database.ID, linodego.DatabaseEngineTypeMySQL,
		linodego.DatabaseStatusUpdating, 240); err != nil {
		t.Fatalf("failed to wait for database updating: %s", err)
	}

	// Wait for the DB to re-enter active status
	if err := client.WaitForDatabaseStatus(
		context.Background(), database.ID, linodego.DatabaseEngineTypeMySQL,
		linodego.DatabaseStatusActive, 2400); err != nil {
		t.Fatalf("failed to wait for database active: %s", err)
	}

	if err := client.SuspendMySQLDatabase(context.Background(), database.ID); err != nil {
		t.Fatalf("failed to suspend database: %s", err)
	}

	// Wait for the DB to enter suspended status
	if err := client.WaitForDatabaseStatus(
		context.Background(), database.ID, linodego.DatabaseEngineTypeMySQL,
		linodego.DatabaseStatusSuspended, 240); err != nil {
		t.Fatalf("failed to wait for database suspended: %s", err)
	}

	if err := client.ResumeMySQLDatabase(context.Background(), database.ID); err != nil {
		t.Fatalf("failed to resume database: %s", err)
	}

	// Wait for the DB to re-enter active status
	if err := client.WaitForDatabaseStatus(
		context.Background(), database.ID, linodego.DatabaseEngineTypeMySQL,
		linodego.DatabaseStatusActive, 2400); err != nil {
		t.Fatalf("failed to wait for database active: %s", err)
	}
}

type mysqlDatabaseModifier func(options *linodego.MySQLCreateOptions)

func createMySQLDatabase(t *testing.T, client *linodego.Client,
	databaseMofidiers []mysqlDatabaseModifier,
) (*linodego.MySQLDatabase, func(), error) {
	t.Helper()

	createOpts := linodego.MySQLCreateOptions{
		Label:       "go-mysql-test-def" + randLabel(),
		Region:      getRegionsWithCaps(t, client, []string{"Managed Databases"})[0],
		Type:        "g6-nanode-1",
		Engine:      "mysql/8",
		ClusterSize: 3,
		AllowList:   []string{"203.0.113.1", "192.0.1.0/24"},
	}

	for _, modifier := range databaseMofidiers {
		modifier(&createOpts)
	}

	database, err := client.CreateMySQLDatabase(context.Background(), createOpts)
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
				err := client.DeleteMySQLDatabase(ctx, database.ID)

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

func setupMySQLDatabase(t *testing.T, databaseMofidiers []mysqlDatabaseModifier,
	fixturesYaml string,
) (*linodego.Client, *linodego.MySQLDatabase, func(), error) {
	t.Helper()
	now := time.Now()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	database, databaseTeardown, err := createMySQLDatabase(t, client, databaseMofidiers)
	if err != nil {
		t.Fatalf("failed to create db: %s", err)
	}

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
