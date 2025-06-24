package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestInstanceSnapshot_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_snapshot_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/backups/1", fixtureData)

	snapshot, err := base.Client.GetInstanceSnapshot(context.Background(), 123, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, snapshot.ID)
	assert.Equal(t, "snapshot-1", snapshot.Label)
	assert.Equal(t, linodego.SnapshotSuccessful, snapshot.Status)
	assert.True(t, snapshot.Available)
	assert.Len(t, snapshot.Configs, 2)
}

func TestInstanceSnapshot_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_snapshot_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("linode/instances/123/backups", fixtureData)

	opts := linodego.InstanceSnapshotCreateOptions{
		Label: "new-snapshot",
	}

	snapshot, err := base.Client.CreateInstanceSnapshot(context.Background(), 123, opts)
	assert.NoError(t, err)
	assert.Equal(t, "new-snapshot", snapshot.Label)
	assert.Equal(t, linodego.SnapshotPending, snapshot.Status)
}

func TestInstanceBackups_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_backups_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/backups", fixtureData)

	backups, err := base.Client.GetInstanceBackups(context.Background(), 123)
	assert.NoError(t, err)
	assert.NotNil(t, backups)
	assert.Len(t, backups.Automatic, 2)
	assert.NotNil(t, backups.Snapshot)
	assert.Equal(t, "auto-backup-1", backups.Automatic[0].Label)
}

func TestInstanceBackups_Enable(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("linode/instances/123/backups/enable", nil)

	err := base.Client.EnableInstanceBackups(context.Background(), 123)
	assert.NoError(t, err)
}

func TestInstanceBackups_Cancel(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("linode/instances/123/backups/cancel", nil)

	err := base.Client.CancelInstanceBackups(context.Background(), 123)
	assert.NoError(t, err)
}

func TestInstanceBackup_Restore(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	restoreOptions := linodego.RestoreInstanceOptions{
		LinodeID:  456,
		Overwrite: true,
	}

	base.MockPost("linode/instances/123/backups/1/restore", nil)

	err := base.Client.RestoreInstanceBackup(context.Background(), 123, 1, restoreOptions)
	assert.NoError(t, err)
}
