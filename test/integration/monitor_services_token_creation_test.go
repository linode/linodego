package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestMonitorServicesTokenCreation_Get_smoke(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestMonitorServiceToken_POST")
	defer teardown()

	// Create a JWE token for the given entity IDs
	createOpts := linodego.MonitorTokenCreateOptions{
		EntityIds: []int{187468, 188020},
	}

	token, getErr := client.CreateMonitorServiceTokenForServiceType(context.Background(), "dbaas", createOpts)
	if getErr != nil {
		t.Errorf("Error creating token : %s", getErr)
	}

	require.NotEmpty(t, token.Token)
}
