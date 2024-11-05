package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountMaintenances_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_maintenance_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/maintenance", fixtureData)

	maintenances, err := base.Client.ListMaintenances(context.Background(), nil)
	if err != nil {
		t.Fatalf("Error listing maintenances: %v", err)
	}

	assert.Equal(t, 1, len(maintenances))
	maintenance := maintenances[0]
	assert.Equal(t, 1234, maintenance.Entity.ID)
	assert.Equal(t, "demo-linode", maintenance.Entity.Label)
	assert.Equal(t, "Linode", maintenance.Entity.Type)
	assert.Equal(t, "https://api.linode.com/v4/linode/instances/{linodeId}", maintenance.Entity.URL)
	assert.Equal(t, "This maintenance will allow us to update the BIOS on the host's motherboard.", maintenance.Reason)
	assert.Equal(t, "started", maintenance.Status)
	assert.Equal(t, "reboot", maintenance.Type)
}
