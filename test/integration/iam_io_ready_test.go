package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupVolumeAttachedToLinode(
	t *testing.T,
	fixturesYaml string,
	detachVolume bool,
) (*linodego.Client, *linodego.Volume, *linodego.Instance, func()) {
	t.Helper()
	var fixtureTeardown func()

	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	instance, err := createInstance(t, client, true, func(l *linodego.Client, options *linodego.InstanceCreateOptions) {
		options.Booted = linodego.Pointer(true)
	})
	require.NoErrorf(t, err, "Error setting up Linode instance: %v", err)

	teardownInstance := func() {
		err = client.DeleteInstance(context.Background(), instance.ID)
		require.NoErrorf(t, err, "Error deleting instance: %v", err)
	}

	instance, err = client.WaitForInstanceStatus(context.Background(), instance.ID, linodego.InstanceRunning, 180)
	require.NoErrorf(t, err, "Error waiting for instance to be running: %v", err)

	volume, teardownVolume, err := createVolume(t, client, func(l *linodego.Client, options *linodego.VolumeCreateOptions) {
		options.Region = instance.Region
		options.Size = 15
	})
	require.NoErrorf(t, err, "Error creating volume: %v", err)

	volume, err = client.AttachVolume(context.Background(), volume.ID, &linodego.VolumeAttachOptions{LinodeID: instance.ID})
	require.NoErrorf(t, err, "Error attaching volume to instance: %v", err)

	volume, err = client.WaitForVolumeIOReadyStatus(context.Background(), volume.ID, true, 45)
	require.NoErrorf(t, err, "Error waiting for IO Ready status of attached volume: %v", err)

	teardown := func() {
		if detachVolume {
			err = client.DetachVolume(context.Background(), volume.ID)
			require.NoErrorf(t, err, "Error detaching volume: %v", err)

			_, err = client.WaitForVolumeIOReadyStatus(context.Background(), volume.ID, false, 45)
			require.NoErrorf(t, err, "Error waiting for IO Ready status of detached volume: %v", err)
		}
		teardownVolume()
		teardownInstance()
		fixtureTeardown()
	}

	return client, volume, instance, teardown
}

func assertVolumeAttachedToInstance(t *testing.T, volume *linodego.Volume, instance *linodego.Instance) {
	t.Helper()

	require.NotNil(t, volume.LinodeID)
	assert.Equal(t, instance.ID, *volume.LinodeID)
	assert.Equal(t, instance.Label, volume.LinodeLabel)
	assert.True(t, volume.IOReady)
}

func requireSingleAttachedInstanceVolume(t *testing.T, client *linodego.Client, instance *linodego.Instance) linodego.Volume {
	t.Helper()

	instanceVolumes, err := client.ListInstanceVolumes(context.Background(), instance.ID, nil)
	require.NoErrorf(t, err, "Error listing instance volumes: %v", err)
	require.Len(t, instanceVolumes, 1, "Expected 1 volume attached to instance, got %d", len(instanceVolumes))

	volume := instanceVolumes[0]
	require.NotNil(t, volume.LinodeID)
	assert.Equal(t, instance.ID, *volume.LinodeID)
	assert.Equal(t, instance.Label, volume.LinodeLabel)
	assert.True(t, volume.IOReady)

	return volume
}

func TestIAM_GetIOReadyForNotAttachedVolume(t *testing.T) {
	client, recordStopper := createTestClient(t, "fixtures/TestIAM_GetIOReadyForNotAttachedVolume")

	volume, teardown, err := createVolume(t, client, func(l *linodego.Client, options *linodego.VolumeCreateOptions) {
		options.Label = label
	})
	require.NoErrorf(t, err, "Error creating not attached volume: %v", err)

	volume, err = client.WaitForVolumeStatus(context.Background(), volume.ID, linodego.VolumeActive, 30)
	require.NoErrorf(t, err, "Error waiting for volume to be active: %v", err)

	volumeList, err := client.ListVolumes(context.Background(), nil)
	require.NoErrorf(t, err, "Error listing volumes: %v", err)
	volumeFound := false

	for _, vol := range volumeList {
		if vol.ID == volume.ID {
			volumeFound = true
			assert.Equal(t, linodego.VolumeActive, vol.Status)
			assert.Empty(t, vol.LinodeID)
			assert.Empty(t, vol.LinodeLabel)
			assert.False(t, vol.IOReady)
			break
		}
	}
	require.True(t, volumeFound, "Volume with ID %d not found in volumeList", volume.ID)

	volume, err = client.GetVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error getting not attached volume: %v", err)
	assert.Equal(t, label, volume.Label)
	assert.Equal(t, linodego.VolumeActive, volume.Status)
	assert.Empty(t, volume.LinodeID)
	assert.Empty(t, volume.LinodeLabel)
	assert.False(t, volume.IOReady)

	t.Cleanup(func() {
		teardown()
		recordStopper()
	})
}

func TestIAM_GetIOReadyForAttachedDetachedVolume(t *testing.T) {
	client, volume, instance, teardown := setupVolumeAttachedToLinode(t, "fixtures/TestIAM_GetIOReadyForAttachedDetachedVolume", false)
	requireSingleAttachedInstanceVolume(t, client, instance)

	volume, err := client.GetVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error getting attached volume: %v", err)
	assertVolumeAttachedToInstance(t, volume, instance)

	err = client.DetachVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error detaching volume: %v", err)

	_, err = client.WaitForVolumeIOReadyStatus(context.Background(), volume.ID, false, 45)
	require.NoErrorf(t, err, "Error waiting for IO Ready status of detached volume: %v", err)

	instanceVolumes, err := client.ListInstanceVolumes(context.Background(), instance.ID, nil)
	require.NoErrorf(t, err, "Error listing instance volumes after detach: %v", err)
	require.Len(t, instanceVolumes, 0, "Expected no volumes attached to instance, got %d", len(instanceVolumes))

	volume, err = client.GetVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error getting volume after detach: %v", err)
	assert.Empty(t, volume.LinodeID)
	assert.Empty(t, volume.LinodeLabel)
	assert.False(t, volume.IOReady)

	t.Cleanup(func() { teardown() })
}

func TestIAM_GetIOReadyForUpdatedVolume(t *testing.T) {
	client, volume, instance, teardown := setupVolumeAttachedToLinode(t, "fixtures/TestIAM_GetIOReadyForUpdatedVolume", true)
	assertVolumeAttachedToInstance(t, volume, instance)
	assert.NotContains(t, "-updated", volume.Label)
	assert.Empty(t, volume.Tags)

	labelUpdated := volume.Label + "-updated"
	tagsUpdated := []string{"updated"}

	updateOpts := linodego.VolumeUpdateOptions{
		Label: labelUpdated,
		Tags:  &tagsUpdated,
	}
	volume, err := client.UpdateVolume(context.Background(), volume.ID, updateOpts)
	require.NoErrorf(t, err, "Error updating volume: %v", err)

	instanceVolume := requireSingleAttachedInstanceVolume(t, client, instance)
	assert.Equal(t, labelUpdated, instanceVolume.Label)
	assert.Equal(t, tagsUpdated, instanceVolume.Tags)

	volume, err = client.GetVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error getting updated volume: %v", err)
	assertVolumeAttachedToInstance(t, volume, instance)
	assert.Equal(t, labelUpdated, volume.Label)
	assert.Equal(t, tagsUpdated, volume.Tags)

	t.Cleanup(func() { teardown() })
}

func TestIAM_GetIOReadyForClonedVolume(t *testing.T) {
	client, volume, instance, teardown := setupVolumeAttachedToLinode(t, "fixtures/TestIAM_GetIOReadyForClonedVolume", true)
	requireSingleAttachedInstanceVolume(t, client, instance)

	labelCloned := volume.Label + "-cloned"
	volumeCloned, err := client.CloneVolume(context.Background(), volume.ID, labelCloned)
	require.NoErrorf(t, err, "Error cloning volume: %v", err)

	_, err = client.WaitForVolumeStatus(context.Background(), volumeCloned.ID, linodego.VolumeActive, 30)
	require.NoErrorf(t, err, "Error waiting for IO Ready status of attached volume: %v", err)

	volumeCloned, err = client.GetVolume(context.Background(), volumeCloned.ID)
	require.NoErrorf(t, err, "Error getting cloned volume: %v", err)
	assert.Equal(t, labelCloned, volumeCloned.Label)
	assert.Equal(t, linodego.VolumeActive, volumeCloned.Status)
	assert.Equal(t, volume.Region, volumeCloned.Region)
	assert.Equal(t, volume.Size, volumeCloned.Size)
	// Cloned volume should not be attached to instance automatically
	requireSingleAttachedInstanceVolume(t, client, instance)
	assert.Empty(t, volumeCloned.LinodeID)
	assert.Empty(t, volumeCloned.LinodeLabel)
	assert.False(t, volumeCloned.IOReady)

	t.Cleanup(func() {
		err = client.DeleteVolume(context.Background(), volumeCloned.ID)
		require.NoErrorf(t, err, "Error deleting cloned volume: %v", err)

		teardown()
	})
}

func TestIAM_GetIOReadyForResizedVolume(t *testing.T) {
	client, volume, instance, teardown := setupVolumeAttachedToLinode(t, "fixtures/TestIAM_GetIOReadyForResizedVolume", true)
	assertVolumeAttachedToInstance(t, volume, instance)

	newSize := volume.Size + 10

	err := client.ResizeVolume(context.Background(), volume.ID, newSize)
	require.NoErrorf(t, err, "Error resizing volume: %v", err)

	_, err = client.WaitForVolumeStatus(context.Background(), volume.ID, linodego.VolumeActive, 30)
	require.NoErrorf(t, err, "Error waiting for volume to be active: %v", err)

	instanceVolume := requireSingleAttachedInstanceVolume(t, client, instance)
	assert.Equal(t, newSize, instanceVolume.Size)

	volume, err = client.GetVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error getting updated volume: %v", err)
	assertVolumeAttachedToInstance(t, volume, instance)
	assert.Equal(t, newSize, volume.Size)

	t.Cleanup(func() { teardown() })
}
