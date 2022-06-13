package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/linode/linodego"
)

var ignoreDatabaseTimestampes = cmpopts.IgnoreFields(linodego.Database{}, "Created", "Updated")

func TestDatabase_Engine(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestDatabase_Engine")
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

func TestDatabase_List(t *testing.T) {
	client, database, teardown, err := setupMongoDatabase(t, nil, "fixtures/TestDatabase_List")
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

func waitForDatabaseUpdated(t *testing.T, client *linodego.Client, dbID int,
	dbType linodego.DatabaseEngineType, minStart *time.Time) {
	_, err := client.WaitForEventFinished(context.Background(), dbID, linodego.EntityDatabase,
		linodego.ActionDatabaseUpdate, *minStart, 1200)
	if err != nil {
		t.Fatalf("failed to wait for database update: %s", err)
	}

	// Sometimes the event has finished but the status hasn't caught up
	err = client.WaitForDatabaseStatus(context.Background(), dbID, dbType,
		linodego.DatabaseStatusActive, 120)
	if err != nil {
		t.Fatalf("failed to wait for database active: %s", err)
	}
}
