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
	defer recordStopper()

	volume, teardown, err := createVolume(t, client, func(l *linodego.Client, options *linodego.VolumeCreateOptions) {
		options.Label = label
	})
	defer teardown()
	require.NoErrorf(t, err, "Error creating not attached volume: %v", err)

	volume, err = client.WaitForVolumeStatus(context.Background(), volume.ID, linodego.VolumeActive, 30)
	require.NoErrorf(t, err, "Error waiting for volume to be active: %v", err)

	volumeList, err := client.ListVolumes(context.Background(), nil)
	require.NoErrorf(t, err, "Error listing volumes: %v", err)
	for _, vol := range volumeList {
		if vol.ID == volume.ID {
			assert.Equal(t, linodego.VolumeActive, vol.Status)
			assert.Empty(t, vol.LinodeID)
			assert.Empty(t, vol.LinodeLabel)
			assert.False(t, vol.IOReady)
			break
		}
	}

	volume, err = client.GetVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error getting not attached volume: %v", err)
	assert.Equal(t, label, volume.Label)
	assert.Equal(t, linodego.VolumeActive, volume.Status)
	assert.Empty(t, volume.LinodeID)
	assert.Empty(t, volume.LinodeLabel)
	assert.False(t, volume.IOReady)
}

func TestIAM_GetIOReadyForAttachedDetachedVolume(t *testing.T) {
	client, volume, instance, teardown := setupVolumeAttachedToLinode(t, "fixtures/TestIAM_GetIOReadyForAttachedDetachedVolume", false)
	defer teardown()

	requireSingleAttachedInstanceVolume(t, client, instance)

	volume, err := client.GetVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error getting attached volume: %v", err)
	assertVolumeAttachedToInstance(t, volume, instance)

	err = client.DetachVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error detaching volume: %v", err)

	volume, err = client.WaitForVolumeIOReadyStatus(context.Background(), volume.ID, false, 45)
	require.NoErrorf(t, err, "Error waiting for IO Ready status of detached volume: %v", err)

	instanceVolumes, err := client.ListInstanceVolumes(context.Background(), instance.ID, nil)
	require.NoErrorf(t, err, "Error listing instance volumes after detach: %v", err)
	require.Len(t, instanceVolumes, 0, "Expected no volumes attached to instance, got %d", len(instanceVolumes))

	volume, err = client.GetVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error getting volume after detach: %v", err)
	assert.Empty(t, volume.LinodeID)
	assert.Empty(t, volume.LinodeLabel)
	assert.False(t, volume.IOReady)
}

func TestIAM_GetIOReadyForUpdatedVolume(t *testing.T) {
	client, volume, instance, teardown := setupVolumeAttachedToLinode(t, "fixtures/TestIAM_GetIOReadyForUpdatedVolume", true)
	defer teardown()
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
}

func TestIAM_GetIOReadyForClonedVolume(t *testing.T) {
	t.Skip("Skipping test due to possible API defect - JIRA ticket: STORENGSUP-855")

	client, volume, instance, teardown := setupVolumeAttachedToLinode(t, "fixtures/TestIAM_GetIOReadyForClonedVolume", true)
	defer teardown()
	assertVolumeAttachedToInstance(t, volume, instance)

	labelCloned := volume.Label + "-cloned"

	requireSingleAttachedInstanceVolume(t, client, instance)

	volumeCloned, err := client.CloneVolume(context.Background(), volume.ID, labelCloned)
	require.NoErrorf(t, err, "Error cloning volume: %v", err)

	instanceVolumes, err := client.ListInstanceVolumes(context.Background(), instance.ID, nil)
	require.NoErrorf(t, err, "Error listing instance volumes: %v", err)
	require.Len(t, instanceVolumes, 2, "Expected 2 volumes attached to instance, got %d", len(instanceVolumes))
	for _, vol := range instanceVolumes {
		assertVolumeAttachedToInstance(t, &vol, instance)
	}

	volumeCloned, err = client.GetVolume(context.Background(), volumeCloned.ID)
	require.NoErrorf(t, err, "Error getting cloned volume: %v", err)
	assert.Equal(t, labelCloned, volumeCloned.Label)
	assertVolumeAttachedToInstance(t, volumeCloned, instance)

	// Cleaning the cloned volume
	err = client.DetachVolume(context.Background(), volumeCloned.ID)
	require.NoErrorf(t, err, "Error detaching cloned volume: %v", err)

	_, err = client.WaitForVolumeIOReadyStatus(context.Background(), volumeCloned.ID, false, 15)
	require.NoErrorf(t, err, "Error waiting for IO Ready status of detached volume: %v", err)

	err = client.DeleteVolume(context.Background(), volumeCloned.ID)
	require.NoErrorf(t, err, "Error deleting cloned volume: %v", err)
}

func TestIAM_GetIOReadyForResizedVolume(t *testing.T) {
	client, volume, instance, teardown := setupVolumeAttachedToLinode(t, "fixtures/TestIAM_GetIOReadyForResizedVolume", true)
	defer teardown()
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
}
