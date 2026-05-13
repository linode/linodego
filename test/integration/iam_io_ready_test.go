package integration

import (
	"context"
	"testing"
	"time"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func findVolumeByID(volumes []linodego.Volume, id int) linodego.Volume {
	for _, vol := range volumes {
		if vol.ID == id {
			return vol
		}
	}

	return linodego.Volume{}
}

func setupVolumeAttachedToLinode(
	t *testing.T,
	ctx context.Context,
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

	instance, err = client.WaitForInstanceStatus(ctx, instance.ID, linodego.InstanceRunning)
	require.NoErrorf(t, err, "Error waiting for instance to be running: %v", err)

	volume, teardownVolume, err := createVolume(t, client, func(l *linodego.Client, options *linodego.VolumeCreateOptions) {
		options.Region = instance.Region
		options.Size = 15
	})
	require.NoErrorf(t, err, "Error creating volume: %v", err)

	volume, err = client.AttachVolume(context.Background(), volume.ID, &linodego.VolumeAttachOptions{LinodeID: instance.ID})
	require.NoErrorf(t, err, "Error attaching volume to instance: %v", err)

	volume, err = client.WaitForVolumeIOReadyStatus(ctx, volume.ID, true)
	require.NoErrorf(t, err, "Error waiting for IO Ready status of attached volume: %v", err)

	teardown := func() {
		if detachVolume {
			// new a context because the context from t.Context() will be cancelled before cleanup run
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
			defer cancel()

			err = client.DetachVolume(cleanupCtx, volume.ID)
			require.NoErrorf(t, err, "Error detaching volume: %v", err)

			_, err = client.WaitForVolumeIOReadyStatus(cleanupCtx, volume.ID, false)
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
	ctx := waitContext(t, 30*time.Second)

	client, recordStopper := createTestClient(t, "fixtures/TestIAM_GetIOReadyForNotAttachedVolume")
	t.Cleanup(recordStopper)

	volume, teardown, err := createVolume(t, client, func(l *linodego.Client, options *linodego.VolumeCreateOptions) {
		options.Label = label
	})
	t.Cleanup(teardown)
	require.NoErrorf(t, err, "Error creating not attached volume: %v", err)

	volume, err = client.WaitForVolumeStatus(ctx, volume.ID, linodego.VolumeActive)
	require.NoErrorf(t, err, "Error waiting for volume to be active: %v", err)

	volumeList, err := client.ListVolumes(context.Background(), nil)
	require.NoErrorf(t, err, "Error listing volumes: %v", err)

	volumeFound := findVolumeByID(volumeList, volume.ID)
	require.NotEmpty(t, volumeFound, "Volume with ID %d not found in volumeList", volume.ID)
	assert.Equal(t, linodego.VolumeActive, volumeFound.Status)
	assert.Empty(t, volumeFound.LinodeID)
	assert.Empty(t, volumeFound.LinodeLabel)
	assert.False(t, volumeFound.IOReady)

	volume, err = client.GetVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error getting not attached volume: %v", err)
	assert.Equal(t, label, volume.Label)
	assert.Equal(t, linodego.VolumeActive, volume.Status)
	assert.Empty(t, volume.LinodeID)
	assert.Empty(t, volume.LinodeLabel)
	assert.False(t, volume.IOReady)
}

func TestIAM_GetIOReadyForAttachedDetachedVolume(t *testing.T) {
	ctx := waitContext(t, 270*time.Second)

	client, volume, instance, teardown := setupVolumeAttachedToLinode(t, ctx, "fixtures/TestIAM_GetIOReadyForAttachedDetachedVolume", false)
	t.Cleanup(teardown)
	requireSingleAttachedInstanceVolume(t, client, instance)

	volume, err := client.GetVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error getting attached volume: %v", err)
	assertVolumeAttachedToInstance(t, volume, instance)

	err = client.DetachVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error detaching volume: %v", err)

	_, err = client.WaitForVolumeIOReadyStatus(ctx, volume.ID, false)
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
	ctx := waitContext(t, 270*time.Second)

	client, volume, instance, teardown := setupVolumeAttachedToLinode(t, ctx, "fixtures/TestIAM_GetIOReadyForUpdatedVolume", true)
	t.Cleanup(teardown)
	assertVolumeAttachedToInstance(t, volume, instance)
	assert.NotContains(t, "-updated", volume.Label)
	assert.Empty(t, volume.Tags)

	labelUpdated := volume.Label + "-updated"
	tagsUpdated := []string{"updated"}

	updateOpts := linodego.VolumeUpdateOptions{
		Label: labelUpdated,
		Tags:  tagsUpdated,
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
	ctx := waitContext(t, 300*time.Second)

	client, volume, instance, teardown := setupVolumeAttachedToLinode(t, ctx, "fixtures/TestIAM_GetIOReadyForClonedVolume", true)
	t.Cleanup(teardown)
	requireSingleAttachedInstanceVolume(t, client, instance)

	labelCloned := volume.Label + "-cloned"
	volumeCloned, err := client.CloneVolume(context.Background(), volume.ID, labelCloned)
	t.Cleanup(func() {
		if err = client.DeleteVolume(context.Background(), volumeCloned.ID); err != nil {
			t.Errorf("Error deleting cloned volume: %v", err)
		}
	})
	require.NoErrorf(t, err, "Error cloning volume: %v", err)

	_, err = client.WaitForVolumeStatus(ctx, volumeCloned.ID, linodego.VolumeActive)
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
}

func TestIAM_GetIOReadyForResizedVolume(t *testing.T) {
	ctx := waitContext(t, 300*time.Second)

	client, volume, instance, teardown := setupVolumeAttachedToLinode(t, ctx, "fixtures/TestIAM_GetIOReadyForResizedVolume", true)
	t.Cleanup(teardown)
	assertVolumeAttachedToInstance(t, volume, instance)

	newSize := volume.Size + 10

	err := client.ResizeVolume(context.Background(), volume.ID, newSize)
	require.NoErrorf(t, err, "Error resizing volume: %v", err)

	_, err = client.WaitForVolumeStatus(ctx, volume.ID, linodego.VolumeActive)
	require.NoErrorf(t, err, "Error waiting for volume to be active: %v", err)

	instanceVolume := requireSingleAttachedInstanceVolume(t, client, instance)
	assert.Equal(t, newSize, instanceVolume.Size)

	volume, err = client.GetVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error getting updated volume: %v", err)
	assertVolumeAttachedToInstance(t, volume, instance)
	assert.Equal(t, newSize, volume.Size)
}
