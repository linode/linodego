package integration

import (
	"context"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/linode/linodego/v2"
)

func TestDatabase_Valkey_Suite(t *testing.T) {
	ctx := waitContext(t, 12000*time.Second)
	client, database, teardown, err := setupValkeyDatabase(t, ctx, "fixtures/TestDatabase_Valkey_Suite")
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	databases, err := client.ListValkeyDatabases(context.Background(), nil)
	if err != nil {
		t.Fatalf("error listing Valkey databases: %v", err)
	}
	if len(databases) == 0 {
		t.Fatal("expected Valkey database list to be non-empty")
	}

	db, err := client.GetValkeyDatabase(context.Background(), database.ID)
	if err != nil {
		t.Fatalf("error getting Valkey database: %v", err)
	}
	if db.ID != database.ID {
		t.Fatalf("got database %d, expected %d", db.ID, database.ID)
	}

	updatedLabel := database.Label + "-updated"
	db, err = client.UpdateValkeyDatabase(context.Background(), database.ID, linodego.ValkeyUpdateOptions{
		Label: updatedLabel,
	})
	if err != nil {
		t.Fatalf("error updating Valkey database: %v", err)
	}
	waitForDatabaseUpdated(t, ctx, client, db.ID, linodego.DatabaseEngineTypeValkey, db.Created)
	if db.Label != updatedLabel {
		t.Fatalf("got label %q, expected %q", db.Label, updatedLabel)
	}

	if _, err := client.GetValkeyDatabaseSSL(context.Background(), database.ID); err != nil {
		t.Fatalf("error getting Valkey SSL certificate: %v", err)
	}
	credentials, err := client.GetValkeyDatabaseCredentials(context.Background(), database.ID)
	if err != nil {
		t.Fatalf("error getting Valkey credentials: %v", err)
	}
	if credentials.Username == "" || credentials.Password == "" {
		t.Fatal("Valkey credentials were empty")
	}

	if testingMode == recorder.ModeRecording {
		time.Sleep(5 * time.Minute)
	}
	if err := client.ResetValkeyDatabaseCredentials(context.Background(), database.ID); err != nil {
		t.Fatalf("error resetting Valkey credentials: %v", err)
	}
	if testingMode == recorder.ModeRecording {
		time.Sleep(15 * time.Second)
	}
	newCredentials, err := client.GetValkeyDatabaseCredentials(context.Background(), database.ID)
	if err != nil {
		t.Fatalf("error getting reset Valkey credentials: %v", err)
	}
	if credentials.Password == newCredentials.Password {
		t.Fatal("Valkey credentials did not change after reset")
	}

	if err := client.PatchValkeyDatabase(context.Background(), database.ID); err != nil {
		t.Fatalf("error patching Valkey database: %v", err)
	}
	waitForValkeyStatus(t, ctx, client, database.ID, linodego.DatabaseStatusUpdating)
	waitForValkeyStatus(t, ctx, client, database.ID, linodego.DatabaseStatusActive)

	if err := client.SuspendValkeyDatabase(context.Background(), database.ID); err != nil {
		t.Fatalf("error suspending Valkey database: %v", err)
	}
	waitForValkeyStatus(t, ctx, client, database.ID, linodego.DatabaseStatusSuspended)

	if err := client.ResumeValkeyDatabase(context.Background(), database.ID); err != nil {
		t.Fatalf("error resuming Valkey database: %v", err)
	}
	waitForValkeyStatus(t, ctx, client, database.ID, linodego.DatabaseStatusActive)
}

func setupValkeyDatabase(t *testing.T, ctx context.Context, fixturesYAML string) (*linodego.Client, *linodego.ValkeyDatabase, func(), error) {
	t.Helper()
	now := time.Now()
	client, fixtureTeardown := createTestClient(t, fixturesYAML)
	regions := getRegionsWithCaps(t, client, []linodego.RegionCapability{linodego.CapabilityDBAAS})
	database, err := client.CreateValkeyDatabase(context.Background(), linodego.ValkeyCreateOptions{
		Label:       "go-valkey-test-" + randLabel(),
		Region:      regions[0],
		Type:        "g7-dedicated-4-2",
		Engine:      "valkey/8.1",
		ClusterSize: 1,
		AllowList:   []string{"203.0.113.1"},
	})
	if err != nil {
		fixtureTeardown()
		return nil, nil, nil, err
	}

	if _, err := client.WaitForEventFinished(ctx, database.ID, linodego.EntityDatabase, linodego.ActionDatabaseCreate, now); err != nil {
		deleteCtx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
		defer cancel()
		_ = client.DeleteValkeyDatabase(deleteCtx, database.ID)
		fixtureTeardown()
		return nil, nil, nil, err
	}

	if err := client.WaitForDatabaseStatus(ctx, database.ID, linodego.DatabaseEngineTypeValkey, linodego.DatabaseStatusActive); err != nil {
		deleteCtx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
		defer cancel()
		_ = client.DeleteValkeyDatabase(deleteCtx, database.ID)
		fixtureTeardown()
		return nil, nil, nil, err
	}

	teardown := func() {
		deleteCtx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
		defer cancel()
		if err := client.DeleteValkeyDatabase(deleteCtx, database.ID); err != nil {
			t.Errorf("failed to delete Valkey database: %v", err)
		}
		fixtureTeardown()
	}
	return client, database, teardown, nil
}

func waitForValkeyStatus(t *testing.T, ctx context.Context, client *linodego.Client, databaseID int, status linodego.DatabaseStatus) {
	t.Helper()
	if err := client.WaitForDatabaseStatus(ctx, databaseID, linodego.DatabaseEngineTypeValkey, status); err != nil {
		t.Fatalf("failed waiting for Valkey status %q: %v", status, err)
	}
}
