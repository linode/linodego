package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestMonitorServicesTokenCreation_Get_smoke(t *testing.T) {

	client, _, teardown, err := setupPostgresDatabase(t, nil, "fixtures/TestDatabaseACLP_List")
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

	// Extract IDs from dbs
	var entityIDs []int
	for _, db := range dbs {
		entityIDs = append(entityIDs, db.ID)
	}

	client1, teardown1 := createTestClient(t, "fixtures/TestMonitorServiceToken_POST")
	defer teardown1()

	// Create a JWE token for the given entity IDs
	createOpts := linodego.MonitorTokenCreateOptions{
		EntityIDs: entityIDs,
	}

	// Use the same context with timeout for the token creation
	token, createErr := client1.CreateMonitorServiceTokenForServiceType(context.Background(), "dbaas", createOpts)
	if createErr != nil {
		t.Errorf("Error creating token : %s", createErr)
	}

	// Validate the token
	validateToken(t, *token)
}

func validateToken(
	t *testing.T,
	token linodego.MonitorServiceToken,
) {
	require.NotEmpty(t, token.Token)
}
