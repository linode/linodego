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

	_, instance, teardownInstance, err := setupInstance(
		t, "fixtures/TestIAM_GetIOReadyForAttachedVolume", true, func(l *linodego.Client, options *linodego.InstanceCreateOptions) {
			options.Booted = linodego.Pointer(true)
		})
	require.NoErrorf(t, err, "Error setting up Linode instance: %v", err)

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
			require.NoErrorf(t, err, "Error waiting for IO Ready status of attached volume: %v", err)
		}
		teardownVolume()
		teardownInstance()
		fixtureTeardown()
	}

	return client, volume, instance, teardown
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
	assert.Equal(t, label, volumeList[0].Label)
	assert.Equal(t, linodego.VolumeActive, volumeList[0].Status)
	assert.Empty(t, volumeList[0].LinodeID)
	assert.Empty(t, volumeList[0].LinodeLabel)
	assert.False(t, volumeList[0].IOReady)

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

	instanceVolumes, err := client.ListInstanceVolumes(context.Background(), instance.ID, nil)
	require.NoErrorf(t, err, "Error listing instance volumes: %v", err)
	require.Len(t, instanceVolumes, 1, "Expected 1 volume attached to instance, got %d", len(instanceVolumes))
	assert.Equal(t, *instanceVolumes[0].LinodeID, *volume.LinodeID)
	assert.Equal(t, instanceVolumes[0].LinodeLabel, volume.LinodeLabel)
	assert.True(t, instanceVolumes[0].IOReady)

	volume, err = client.GetVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error getting volume after detaching from instance: %v", err)
	assert.Equal(t, instance.ID, *volume.LinodeID)
	assert.Equal(t, instance.Label, volume.LinodeLabel)
	assert.True(t, volume.IOReady)

	err = client.DetachVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error detaching volume: %v", err)

	volume, err = client.WaitForVolumeIOReadyStatus(context.Background(), volume.ID, false, 45)
	require.NoErrorf(t, err, "Error waiting for IO Ready status of attached volume: %v", err)

	instanceVolumes, err = client.ListInstanceVolumes(context.Background(), instance.ID, nil)
	require.NoErrorf(t, err, "Error listing instance volumes: %v", err)
	require.Len(t, instanceVolumes, 0, "Expected no volumes attached to instance, got %d", len(instanceVolumes))

	volume, err = client.GetVolume(context.Background(), volume.ID)
	require.NoErrorf(t, err, "Error getting volume after detaching from instance: %v", err)
	assert.Empty(t, volume.LinodeID)
	assert.Empty(t, volume.LinodeLabel)
	assert.False(t, volume.IOReady)
}
