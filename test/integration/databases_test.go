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
}

func createDatabase(t *testing.T, client *linodego.Client) (*linodego.Database, func(), error) {
	t.Helper()

	createOpts := testMySQLCreateOpts
	database, err := client.CreateMySQL(context.Background(), createOpts)
	if err != nil {
		t.Fatalf("failed to create database: %s", err)
	}

	teardown := func() {
		if err := client.DeleteMySQL(context.Background(), database.ID); err != nil {
			t.Fatalf("failed to delete database: %s", err)
		}
	}
	return database, teardown, nil
}

func setupDatabase(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.Database, func(), error) {
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
