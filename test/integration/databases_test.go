package integration

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/linode/linodego"
)

var ignoreDatabaseTimestampes = cmpopts.IgnoreFields(linodego.Database{}, "Created", "Updated")

func TestListDatabases(t *testing.T) {
	client, _, teardown, err := setupDatabase(t, []databaseModifier{
		func(createOpts *linodego.DatabaseCreateOptions) {
			createOpts.Label = randString(12, lowerBytes, upperBytes) + "-linodego-testing"
		},
	}, "fixtures/TestListDatabases")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	result, err := client.ListDatabases(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing Databases, expected struct, got error %v", err)
	}

	if len(result) == 0 {
		t.Errorf("Expected a list of Databases, but got none: %v", err)
	}
}

type databaseModifier func(*linodego.DatabaseCreateOptions)

func createDatabase(t *testing.T, client *linodego.Client, mods ...databaseModifier) (*linodego.Database, func(), error) {
	t.Helper()

	createOpts := testDatabaseCreateOpts
	for _, mod := range mods {
		mod(&createOpts)
	}

	database, err := client.CreateDatabase(context.Background(), createOpts)
	if err != nil {
		t.Errorf("failed to create database: %s", err)
	}

	teardown := func() {
		if err := client.DeleteDatabase(context.Background(), database.ID); err != nil {
			t.Errorf("failed to delete database: %s", err)
		}
	}
	return database, teardown, nil
}

func setupDatabase(t *testing.T, mods []databaseModifier, fixturesYaml string) (*linodego.Client, *linodego.Database, func(), error) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	database, databaseTeardown, err := createDatabase(t, client, mods...)

	teardown := func() {
		databaseTeardown()
		fixtureTeardown()
	}
	return client, firewall, teardown, err
}
