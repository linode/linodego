package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestInstance_Disks_Clone(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_disks_clone")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("linode/instances/12345/disks/123/clone", fixtureData)

	disk, err := base.Client.CloneInstanceDisk(context.Background(), 12345, 123)
	assert.NoError(t, err)

	assert.Equal(t, linodego.DiskFilesystem("ext4"), disk.Filesystem)
	assert.Equal(t, 123, disk.ID)
	assert.Equal(t, "Debian 9 Disk", disk.Label)
	assert.Equal(t, 48640, disk.Size)
	assert.Equal(t, linodego.DiskStatus("ready"), disk.Status)
}

func TestInstanceDisk_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_disk_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/disks", fixtureData)

	disks, err := base.Client.ListInstanceDisks(context.Background(), 123, nil)
	assert.NoError(t, err)
	assert.Len(t, disks, 2)

	assert.Equal(t, 1, disks[0].ID)
	assert.Equal(t, "Disk 1", disks[0].Label)
	assert.Equal(t, linodego.DiskReady, disks[0].Status)
	assert.Equal(t, 20480, disks[0].Size)
	assert.Equal(t, linodego.FilesystemExt4, disks[0].Filesystem)
}

func TestInstanceDisk_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_disk_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/disks/1", fixtureData)

	disk, err := base.Client.GetInstanceDisk(context.Background(), 123, 1)
	assert.NoError(t, err)

	assert.Equal(t, 1, disk.ID)
	assert.Equal(t, "Disk 1", disk.Label)
	assert.Equal(t, linodego.DiskReady, disk.Status)
	assert.Equal(t, 20480, disk.Size)
	assert.Equal(t, linodego.FilesystemExt4, disk.Filesystem)
	assert.NotNil(t, disk.Created)
	assert.NotNil(t, disk.Updated)
}

func TestInstanceDisk_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_disk_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	createOptions := linodego.InstanceDiskCreateOptions{
		Label:      "New Disk",
		Size:       20480,
		Filesystem: "ext4",
	}

	base.MockPost("linode/instances/123/disks", fixtureData)

	disk, err := base.Client.CreateInstanceDisk(context.Background(), 123, createOptions)
	assert.NoError(t, err)

	assert.Equal(t, 3, disk.ID)
	assert.Equal(t, "New Disk", disk.Label)
	assert.Equal(t, linodego.DiskReady, disk.Status)
	assert.Equal(t, 20480, disk.Size)
	assert.Equal(t, linodego.FilesystemExt4, disk.Filesystem)
	assert.NotNil(t, disk.Created)
	assert.NotNil(t, disk.Updated)
}

func TestInstanceDisk_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_disk_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	updateOptions := linodego.InstanceDiskUpdateOptions{
		Label: "Updated Disk",
	}

	base.MockPut("linode/instances/123/disks/1", fixtureData)

	disk, err := base.Client.UpdateInstanceDisk(context.Background(), 123, 1, updateOptions)
	assert.NoError(t, err)

	assert.Equal(t, 1, disk.ID)
	assert.Equal(t, "Updated Disk", disk.Label)
	assert.Equal(t, linodego.DiskReady, disk.Status)
	assert.Equal(t, 20480, disk.Size)
	assert.Equal(t, linodego.FilesystemExt4, disk.Filesystem)
	assert.NotNil(t, disk.Created)
	assert.NotNil(t, disk.Updated)
}

func TestInstanceDisk_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("linode/instances/123/disks/1", nil)

	err := base.Client.DeleteInstanceDisk(context.Background(), 123, 1)
	assert.NoError(t, err)
}

func TestInstanceDisk_Resize(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("linode/instances/123/disks/1/resize", nil)

	err := base.Client.ResizeInstanceDisk(context.Background(), 123, 1, 40960)
	assert.NoError(t, err)
}

func TestInstanceDisk_PasswordReset(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("linode/instances/123/disks/1/password", nil)

	err := base.Client.PasswordResetInstanceDisk(context.Background(), 123, 1, "new-password")
	assert.NoError(t, err)
}
