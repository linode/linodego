package unit

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestListDatabases(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("databases_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/instances", fixtureData)
	databases, err := base.Client.ListDatabases(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, databases, "Expected non-empty database list")
}

func TestGetDatabaseEngine(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("database_engine_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	engineID := "mysql-8"
	base.MockGet(fmt.Sprintf("databases/engines/%s", engineID), fixtureData)
	databaseEngine, err := base.Client.GetDatabaseEngine(context.Background(), &linodego.ListOptions{}, engineID)
	assert.NoError(t, err)
	assert.NotNil(t, databaseEngine, "Expected database engine object to be returned")
	assert.Equal(t, engineID, databaseEngine.ID, "Expected correct database engine ID")
	assert.Equal(t, "mysql", databaseEngine.Engine, "Expected MySQL engine")
	assert.Equal(t, "8.0", databaseEngine.Version, "Expected MySQL 8.0 version")
}

func TestListDatabaseTypes(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("database_types_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/types", fixtureData)
	databaseTypes, err := base.Client.ListDatabaseTypes(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, databaseTypes, "Expected non-empty database types list")
}

func TestUnmarshalDatabase(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("database_unmarshal")
	assert.NoError(t, err)

	var data []byte
	switch v := fixtureData.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	case map[string]interface{}:
		data, err = json.Marshal(v) // Convert map to JSON string
		assert.NoError(t, err, "Failed to marshal fixtureData")
	default:
		assert.Fail(t, "Unexpected fixtureData type")
	}

	var db linodego.Database
	err = json.Unmarshal(data, &db)
	assert.NoError(t, err)
	assert.Equal(t, 123, db.ID, "Expected correct database ID")
	assert.Equal(t, "active", string(db.Status), "Expected active status")
	assert.Equal(t, "mysql", db.Engine, "Expected MySQL engine")
	assert.Equal(t, 3, db.ClusterSize, "Expected cluster size 3")
	assert.NotNil(t, db.Created, "Expected Created timestamp to be set")
}

func TestDatabaseMaintenanceWindowUnmarshal(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("database_maintenance_window")
	assert.NoError(t, err)

	var data []byte
	switch v := fixtureData.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	case map[string]interface{}:
		data, err = json.Marshal(v)
		assert.NoError(t, err, "Failed to marshal fixtureData")
	default:
		assert.Fail(t, "Unexpected fixtureData type")
	}

	var window linodego.DatabaseMaintenanceWindow
	err = json.Unmarshal(data, &window)
	assert.NoError(t, err)
	assert.Equal(t, linodego.DatabaseMaintenanceDayMonday, window.DayOfWeek, "Expected Monday as maintenance day")
	assert.Equal(t, 2, window.Duration, "Expected 2-hour maintenance window")
	assert.Equal(t, linodego.DatabaseMaintenanceFrequencyWeekly, window.Frequency, "Expected weekly frequency")
	assert.Equal(t, 3, window.HourOfDay, "Expected maintenance at 3 AM")
}

func TestDatabaseStatusAssertions(t *testing.T) {
	expectedStatuses := []string{
		"provisioning", "active", "deleting", "deleted",
		"suspending", "suspended", "resuming", "restoring",
		"failed", "degraded", "updating", "backing_up",
	}

	statuses := []linodego.DatabaseStatus{
		linodego.DatabaseStatusProvisioning, linodego.DatabaseStatusActive, linodego.DatabaseStatusDeleting,
		linodego.DatabaseStatusDeleted, linodego.DatabaseStatusSuspending, linodego.DatabaseStatusSuspended,
		linodego.DatabaseStatusResuming, linodego.DatabaseStatusRestoring, linodego.DatabaseStatusFailed,
		linodego.DatabaseStatusDegraded, linodego.DatabaseStatusUpdating, linodego.DatabaseStatusBackingUp,
	}

	for _, status := range statuses {
		t.Run(string(status), func(t *testing.T) {
			exists := slices.ContainsFunc(expectedStatuses, func(s string) bool {
				return s == string(status)
			})
			assert.True(t, exists, fmt.Sprintf("Expected status %s to exist", status))
		})
	}
}

func TestDatabaseMaintenanceFrequencyAssertions(t *testing.T) {
	expectedFrequencies := []string{"weekly", "monthly"}

	frequencies := []linodego.DatabaseMaintenanceFrequency{
		linodego.DatabaseMaintenanceFrequencyWeekly,
		linodego.DatabaseMaintenanceFrequencyMonthly,
	}

	for _, freq := range frequencies {
		t.Run(string(freq), func(t *testing.T) {
			exists := slices.ContainsFunc(expectedFrequencies, func(f string) bool {
				return f == string(freq)
			})
			assert.True(t, exists, fmt.Sprintf("Expected frequency %s to exist", freq))
		})
	}
}
