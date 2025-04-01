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
