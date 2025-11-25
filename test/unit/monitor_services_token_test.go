package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestCreateMonitorServicesToken(t *testing.T) {
	// Load the mock fixture for monitor services
	fixtureData, err := fixtures.GetFixture("monitor_service_token_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("monitor/services/dbaas/token", fixtureData)

	// Create request data for POST request
	opts := linodego.MonitorTokenCreateOptions{
		EntityIDs: []any{12345, "us-54321"},
	}

	token, err := base.Client.CreateMonitorServiceTokenForServiceType(context.Background(), "dbaas", opts)
	assert.NoError(t, err)
	assert.NotNil(t, token)
}
