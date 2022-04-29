package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/linode/linodego"
)

var testMySQLCreateOpts = linodego.MySQLCreateOptions{
	Label:           "basic-mysql1-linodego-testing",
	Region:          "us-east",
	Type:            "g6-nanode-1",
	Engine:          "mysql/8.0.26",
	Encrypted:       false,
	ClusterSize:     3,
	ReplicationType: "semi_synch",
	SSLConnection:   false,
	AllowList:       []string{"203.0.113.1", "192.0.1.0/24"},
}

var ignoreDatabaseTimestampes = cmpopts.IgnoreFields(linodego.Database{}, "Created", "Updated")

func TestDatabaseEngine(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestDatabaseEngine")
	defer teardown()

	engines, err := client.ListDatabaseEngines(context.Background(), nil)
	if err != nil {
		t.Errorf("Failed to get list of Database Engines: %v", err)
	}

	if len(engines) <= 0 {
		t.Fatal("failed to get list of database engines")
	}

	engine := engines[0]

	response, err := client.GetDatabaseEngine(context.Background(), nil, engine.ID)
	if err != nil {
		t.Errorf("Failed to get one database Engine: %v", err)
	}

	if engine.Engine != response.Engine {
		t.Fatal("recieved engine does not match source")
	}
}

func TestDatabase_Type(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestDatabase_Type")
	defer teardown()

	types, err := client.ListDatabaseTypes(context.Background(), nil)
	if err != nil {
		t.Errorf("Failed to get list of Database Types: %v", err)
	}

	if len(types) <= 0 {
		t.Fatal("failed to get list of database Types")
	}

	aType := types[0]

	response, err := client.GetDatabaseType(context.Background(), nil, aType.ID)
	if err != nil {
		t.Errorf("Failed to get one database Type: %v", err)
	}

	if aType.Label != response.Label {
		t.Fatal("recieved type does not match source")
	}

	if response.Engines.MySQL[0].Quantity != aType.Engines.MySQL[0].Quantity {
		t.Fatalf("mismatched type quantity: %d, %d", response.Engines.MySQL[0].Quantity, aType.Engines.MySQL[0].Quantity)
	}
}

func TestDatabaseSuite(t *testing.T) {
	client, database, teardown, err := setupDatabase(t, "fixtures/TestDatabaseSuite")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	dbs, err := client.ListDatabases(context.Background(), nil)
	if len(dbs) == 0 {
		t.Errorf("Expected a list of Databases, but got none: %v", err)
	}
	if err != nil {
		t.Errorf("Error listing Databases, expected struct, got error %v", err)
	}
	success := false
	for _, db := range dbs {
		if db.ID == database.ID {
			success = true
		}
	}
	if !success {
		t.Error("database not in database list")
	}

	mysqldbs, err := client.ListMySQLDatabases(context.Background(), nil)
	if len(dbs) == 0 {
		t.Errorf("Expected a list of Databases, but got none: %v", err)
	}
	if err != nil {
		t.Errorf("Error listing Databases, expected struct, got error %v", err)
	}
	success = false
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

	opts := linodego.MySQLUpdateOptions{
		AllowList: []string{"128.173.205.21", "123.177.200.20"},
		Label:     "updated-mysql1-linodego-testing",
	}
	db, err = client.UpdateMySQLDatabase(context.Background(), database.ID, opts)
	if err != nil {
		t.Errorf("failed to update db: %d", database.ID)
	}
	if db.ID != database.ID {
		t.Errorf("updated db does not match original id")
	}
	if db.Label != "updated-mysql1-linodego-testing" {
		t.Errorf("label not updated for db")
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
	time.Sleep(time.Minute * 5)
	err = client.ResetMySQLDatabaseCredentials(context.Background(), database.ID)
	if err != nil {
		t.Errorf("failed to reset credentials for db: %v", err)
	}
	time.Sleep(time.Second * 15)
	newcreds, err := client.GetMySQLDatabaseCredentials(context.Background(), database.ID)
	if err != nil {
		t.Errorf("failed to get new credentials for db: %v", err)
	}
	if creds.Password == newcreds.Password {
		t.Error("credentials have not changed for db")
	}

	backups, err := client.ListMySQLDatabaseBackups(context.Background(), database.ID, nil)
	if err != nil {
		t.Errorf("failed to get backups for db: %v", err)
	}

	if len(backups) > 0 {
		t.Errorf("expected 0 backups, recieved some: %v", backups)
	}

	// can't test get mysql/instances/{id}/backups/{id} until on demand backups
	// can't test post mysql/instances/{id}/backups/{id}/restore until on demand backups
}

func createDatabase(t *testing.T, client *linodego.Client) (*linodego.MySQLDatabase, func(), error) {
	t.Helper()

	createOpts := testMySQLCreateOpts
	database, err := client.CreateMySQLDatabase(context.Background(), createOpts)
	if err != nil {
		t.Fatalf("failed to create database: %s", err)
	}

	teardown := func() {
		if err := client.DeleteMySQLDatabase(context.Background(), database.ID); err != nil {
			t.Fatalf("failed to delete database: %s", err)
		}
	}
	return database, teardown, nil
}

func setupDatabase(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.MySQLDatabase, func(), error) {
	t.Helper()
	now := time.Now()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	database, databaseTeardown, err := createDatabase(t, client)
	if err != nil {
		t.Fatalf("failed to create db: %s", err)
	}

	_, err = client.WaitForEventFinished(context.Background(), database.ID, linodego.EntityDatabase, linodego.ActionDatabaseCreate, now, 3600)
	if err != nil {
		t.Fatalf("failed to wait for db create event: %s", err)
	}

	teardown := func() {
		databaseTeardown()
		fixtureTeardown()
	}
	return client, database, teardown, err
}
