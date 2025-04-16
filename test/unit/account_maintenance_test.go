package unit

import (
	"context"
	"testing"
	"time"

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
	assert.Equal(t, "Scheduled upgrade to faster NVMe hardware. This will affect Linode #1234.", maintenance.Body)
	assert.Equal(t, 1234, maintenance.Entity.ID)
	assert.Equal(t, "Linode #1234", maintenance.Entity.Label)
	assert.Equal(t, "linode", maintenance.Entity.Type)
	assert.Equal(t, "/linodes/1234", maintenance.Entity.URL)
	assert.Equal(t, "Scheduled Maintenance for Linode #1234", maintenance.Label)
	assert.Equal(t, "Scheduled upgrade to faster NVMe hardware.", maintenance.Message)
	assert.Equal(t, "major", maintenance.Severity)
	assert.Equal(t, "maintenance_scheduled", maintenance.Type)
	assert.Equal(t, "linode_migrate", maintenance.EventType)
	assert.Equal(t, "Power on/off", maintenance.MaintenancePolicySet)
	assert.Equal(t, "Scheduled Maintenance", maintenance.Description)
	assert.Equal(t, "platform", maintenance.Source)
	assert.Equal(t, "2025-03-25T10:00:00Z", maintenance.NotBefore.Format(time.RFC3339))
	assert.Equal(t, "2025-03-25T12:00:00Z", maintenance.StartTime.Format(time.RFC3339))
	assert.Equal(t, "2025-03-25T14:00:00Z", maintenance.CompleteTime.Format(time.RFC3339))
	assert.Equal(t, "scheduled", maintenance.Status)
	assert.Equal(t, "2025-03-25T12:00:00Z", maintenance.When.Format(time.RFC3339))
	assert.Equal(t, "2025-03-25T14:00:00Z", maintenance.Until.Format(time.RFC3339))
}
