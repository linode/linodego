package unit

import (
	"context"
	"fmt"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestListVolumes(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("volumes_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("volumes", fixtureData)

	volumes, err := base.Client.ListVolumes(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.NotEmpty(t, volumes, "Expected non-empty volumes list")

	// Assert specific volume details
	assert.Equal(t, 123, volumes[0].ID, "Expected volume ID to match")
	assert.Equal(t, "Test Volume", volumes[0].Label, "Expected volume label to match")
	assert.Equal(t, "active", string(volumes[0].Status), "Expected volume status to match")
	assert.Equal(t, "us-east", volumes[0].Region, "Expected volume region to match")
	assert.Equal(t, 20, volumes[0].Size, "Expected volume size to match")
	assert.Equal(t, "test", volumes[0].Tags[0], "Expected volume tag to match")
}

func TestGetVolume(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("volume_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	volumeID := 123
	base.MockGet(fmt.Sprintf("volumes/%d", volumeID), fixtureData)

	volume, err := base.Client.GetVolume(context.Background(), volumeID)
	assert.NoError(t, err)

	// Assert all fields
	assert.Equal(t, 123, volume.ID, "Expected volume ID to match")
	assert.Equal(t, "Test Volume", volume.Label, "Expected volume label to match")
	assert.Equal(t, "active", string(volume.Status), "Expected volume status to match")
	assert.Equal(t, "us-east", volume.Region, "Expected volume region to match")
	assert.Equal(t, 20, volume.Size, "Expected volume size to match")
	assert.Nil(t, volume.LinodeID, "Expected LinodeID to be nil")
	assert.Empty(t, volume.FilesystemPath, "Expected filesystem path to be empty")
	assert.Contains(t, volume.Tags, "test", "Expected tags to include 'test'")
	assert.Empty(t, volume.HardwareType, "Expected hardware type to be empty")
	assert.Empty(t, volume.LinodeLabel, "Expected Linode label to be empty")
}

func TestCreateVolume(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("volume_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("volumes", fixtureData)

	label := "new-volume"
	size := 20

	opts := linodego.VolumeCreateOptions{
		Label: &label,
		Size:  &size,
		Tags:  []string{"test"},
	}

	volume, err := base.Client.CreateVolume(context.Background(), opts)
	assert.NoError(t, err)

	// Assert all fields
	assert.Equal(t, 124, volume.ID, "Expected created volume ID to match")
	assert.Equal(t, "new-volume", volume.Label, "Expected created volume label to match")
	assert.Equal(t, "creating", string(volume.Status), "Expected created volume status to be 'creating'")
	assert.Equal(t, "us-east", volume.Region, "Expected created volume region to match")
	assert.Equal(t, 20, volume.Size, "Expected created volume size to match")
	assert.Nil(t, volume.LinodeID, "Expected LinodeID to be nil for newly created volume")
	assert.Contains(t, volume.Tags, "test", "Expected created volume tags to include 'test'")
}

func TestUpdateVolume(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("volume_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	volumeID := 123
	base.MockPut(fmt.Sprintf("volumes/%d", volumeID), fixtureData)

	label := "new-volume"

	opts := linodego.VolumeUpdateOptions{
		Label: &label,
		Tags:  []string{"updated"},
	}

	updatedVolume, err := base.Client.UpdateVolume(context.Background(), volumeID, opts)
	assert.NoError(t, err)

	// Assert all fields
	assert.Equal(t, 123, updatedVolume.ID, "Expected updated volume ID to match")
	assert.Equal(t, "updated-volume", updatedVolume.Label, "Expected updated volume label to match")
	assert.Equal(t, "active", string(updatedVolume.Status), "Expected updated volume status to match")
	assert.Contains(t, updatedVolume.Tags, "updated", "Expected updated volume tags to include 'updated'")
}

func TestDeleteVolume(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	volumeID := 123
	base.MockDelete(fmt.Sprintf("volumes/%d", volumeID), nil)

	err := base.Client.DeleteVolume(context.Background(), volumeID)
	assert.NoError(t, err, "Expected no error when deleting volume")
}

func TestAttachVolume(t *testing.T) {
	// Mock the API response for attaching a volume
	fixtureData, err := fixtures.GetFixture("volume_attach")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	volumeID := 123
	base.MockPost(fmt.Sprintf("volumes/%d/attach", volumeID), fixtureData)

	// Use direct pointer assignment for PersistAcrossBoots
	persistAcrossBoots := true
	opts := &linodego.VolumeAttachOptions{
		LinodeID:           456,
		PersistAcrossBoots: &persistAcrossBoots,
	}

	attachedVolume, err := base.Client.AttachVolume(context.Background(), volumeID, opts)
	assert.NoError(t, err, "Expected no error when attaching volume")

	// Verify the attached volume's LinodeID and filesystem path
	assert.Equal(t, 456, *attachedVolume.LinodeID, "Expected LinodeID to match input")
	assert.Equal(t, "/dev/disk/by-id/volume-123", attachedVolume.FilesystemPath, "Expected filesystem path to match fixture")
}

func TestDetachVolume(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	volumeID := 123
	base.MockPost(fmt.Sprintf("volumes/%d/detach", volumeID), nil)

	err := base.Client.DetachVolume(context.Background(), volumeID)
	assert.NoError(t, err, "Expected no error when detaching volume")
}

func TestResizeVolume(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	volumeID := 123
	base.MockPost(fmt.Sprintf("volumes/%d/resize", volumeID), nil)

	err := base.Client.ResizeVolume(context.Background(), volumeID, 50)
	assert.NoError(t, err, "Expected no error when resizing volume")
}
