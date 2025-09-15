package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
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

	var entityIDs []any
	for _, db := range dbs {
		entityIDs = append(entityIDs, db.ID)
	}

	client, teardown = createTestClient(t, "fixtures/TestServiceToken_POST")
	defer teardown()

	// Create a JWE token for the given entity IDs
	createOpts := linodego.MonitorTokenCreateOptions{
		EntityIDs: entityIDs,
	}

	// Use the same context with timeout for the token creation
	token, createErr := client.CreateMonitorServiceTokenForServiceType(context.Background(), "dbaas", createOpts)
	if createErr != nil {
		t.Errorf("Error creating token : %s", createErr)
	}

	// Validate the token
	validateToken(t, *token)
}

func TestMonitorServicesTokenCreation_TryCreateTokenWithInvalidEntityID(t *testing.T) {
	var entityIDs []any
	entityIDs = append(entityIDs, 999999999)

	client, teardown := createTestClient(t, "fixtures/TestServiceToken_POST")
	defer teardown()
	createOptsWithInvalidID := linodego.MonitorTokenCreateOptions{
		EntityIDs: entityIDs,
	}
	_, createErr := client.CreateMonitorServiceTokenForServiceType(context.Background(), "dbaas", createOptsWithInvalidID)
	require.Error(t, createErr)
	assert.Contains(t, createErr.Error(), "[403]")
	assert.Contains(t, createErr.Error(), "The following entity_ids are not valid - [999999999]")
}

func validateToken(
	t *testing.T,
	token linodego.MonitorServiceToken,
) {
	require.NotEmpty(t, token.Token)
}
