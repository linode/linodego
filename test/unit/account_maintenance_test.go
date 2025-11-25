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

	assert.Equal(t, 2, len(maintenances))

	assert.Equal(t, 1234, maintenances[0].Entity.ID)
	assert.Equal(t, "Linode #1234", maintenances[0].Entity.Label)
	assert.Equal(t, "linode", maintenances[0].Entity.Type)
	assert.Equal(t, "/linodes/1234", maintenances[0].Entity.URL)
	assert.Equal(t, "Scheduled upgrade to faster NVMe hardware.", maintenances[0].Reason)
	assert.Equal(t, "linode_migrate", maintenances[0].Type)
	assert.Equal(t, "Power on/off", maintenances[0].MaintenancePolicySet)
	assert.Equal(t, "Scheduled Maintenance", maintenances[0].Description)
	assert.Equal(t, "platform", maintenances[0].Source)
	assert.Equal(t, "2025-03-25T10:00:00Z", maintenances[0].NotBefore.Format(time.RFC3339))
	assert.Equal(t, "2025-03-25T12:00:00Z", maintenances[0].StartTime.Format(time.RFC3339))
	assert.Equal(t, "2025-03-25T14:00:00Z", maintenances[0].CompleteTime.Format(time.RFC3339))
	assert.Equal(t, "scheduled", maintenances[0].Status)

	assert.Equal(t, 1234, maintenances[1].Entity.ID)
	assert.Equal(t, "Linode #1234", maintenances[1].Entity.Label)
	assert.Equal(t, "linode", maintenances[1].Entity.Type)
	assert.Equal(t, "/linodes/1234", maintenances[1].Entity.URL)
	assert.Equal(t, "Pending migration of Linode #1234 to a new host.", maintenances[1].Reason)
	assert.Equal(t, "linode_migrate", maintenances[1].Type)
	assert.Equal(t, "Migrate", maintenances[1].MaintenancePolicySet)
	assert.Equal(t, "Emergency Maintenance", maintenances[1].Description)
	assert.Equal(t, "user", maintenances[1].Source)
	assert.Equal(t, "2025-03-25T10:00:00Z", maintenances[1].NotBefore.Format(time.RFC3339))
	assert.Equal(t, "2025-03-25T12:00:00Z", maintenances[1].StartTime.Format(time.RFC3339))
	assert.Equal(t, "2025-03-25T14:00:00Z", maintenances[1].CompleteTime.Format(time.RFC3339))
	assert.Equal(t, "in-progress", maintenances[1].Status)
}
