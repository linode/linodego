package integration

import (
	"context"
	"testing"
	"time"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getNFSFilesystemCreateOptions(t *testing.T, client *linodego.Client) linodego.NFSFilesystemCreateOptions {
	return linodego.NFSFilesystemCreateOptions{
		Label:            "go-test-nfs-filesystem-mw-" + randLabel(),
		Region:           getRegionsWithCaps(t, client, []linodego.RegionCapability{linodego.CapabilityNFSStorage})[0],
		MaxCapacityBytes: 1099511627776,
		Tags:             linodego.Pointer([]string{"testing"}),
		// MaxFileCount:	  1000000,
		// ProtocolVersions: linodego.Pointer([]linodego.NFSProtocolVersion{linodego.NFSProtocolVersionV4}),
	}
}

func waitForNFSFilesystemDeleted(
	ctx context.Context,
	client *linodego.Client,
	spaceID int,
) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			filesystems, err := client.ListNFSFilesystems(ctx, spaceID, nil)

			if err == nil && len(filesystems) != 0 {
				continue
			}

			// if API returns [404] Not found, deletion is done
			if linodego.IsNotFound(err) {
				return nil
			}

			return err
		}
	}
}

func setupNFSFilesystem(
	t *testing.T,
	fixtureYaml string,
	modifiers ...func(opts *linodego.NFSFilesystemCreateOptions),
) (*linodego.Client, *linodego.NFSSpace, *linodego.NFSFilesystem, linodego.NFSFilesystemCreateOptions) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixtureYaml)
	createOpts := getNFSFilesystemCreateOptions(t, client)

	for _, modifier := range modifiers {
		modifier(&createOpts)
	}

	spaceCreateOpts := getNFSSpaceCreateOptions()
	space, err := client.CreateNFSSpace(context.Background(), spaceCreateOpts)
	require.NoErrorf(t, err, "Error creating NFS Space: %v", err)

	filesystem, err := client.CreateNFSFilesystem(context.Background(), space.ID, createOpts)
	require.NoErrorf(t, err, "Error creating NFS Filesystem: %v", err)

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
		defer cancel()

		if err = client.DeleteNFSFilesystem(context.Background(), space.ID, filesystem.ID); err != nil {
			t.Errorf("Error deleting NFS Filesystem: %v", err)
		}

		err = waitForNFSFilesystemDeleted(ctx, client, space.ID)
		require.NoErrorf(t, err, "Error waiting for NFS Filesystem to be deleted: %v", err)

		if err = client.DeleteNFSSpace(context.Background(), space.ID); err != nil {
			t.Errorf("Error deleting NFS Space: %v", err)
		}

		fixtureTeardown()
	})

	return client, space, filesystem, createOpts
}

func verifyNFSFilesystemBasics(t *testing.T, space *linodego.NFSSpace, filesystem *linodego.NFSFilesystem, createOpts linodego.NFSFilesystemCreateOptions) {
	t.Helper()
	assert.Equal(t, space.ID, filesystem.SpaceID)
	assert.Equal(t, createOpts.Label, filesystem.Label)
	assert.Equal(t, createOpts.Region, filesystem.Region)
	assert.Equal(t, createOpts.Tags, linodego.Pointer(filesystem.Tags))
	assert.NotNil(t, filesystem.Status)
	assertDateSet(t, filesystem.Created)
	assertDateSet(t, filesystem.Updated)
}

func verifyNFSFilesystemDetails(t *testing.T, filesystem *linodego.NFSFilesystem, createOpts linodego.NFSFilesystemCreateOptions) {
	t.Helper()
	assert.NotEmpty(t, filesystem.ProtocolVersions)
	assert.EqualValues(t, createOpts.MaxCapacityBytes, filesystem.MaxCapacityBytes)
	assert.NotEmpty(t, filesystem.MountTargetIPs)
	assert.NotNil(t, filesystem.MountTargetFQDN)
	assert.NotNil(t, filesystem.Stats)
	//assert.NotNil(t, filesystem.SnapshotUsageBytes)
	//assert.NotNil(t, filesystem.LDAPConfigID)
	//assert.NotNil(t, filesystem.SourceSnapshotID)
}

func TestNFSFilesystem_Create_smoke(t *testing.T) {
	_, space, filesystem, createOpts := setupNFSFilesystem(t, "fixtures/TestNFSFilesystem_Create")
	verifyNFSFilesystemBasics(t, space, filesystem, createOpts)
}

func TestNFSFilesystem_Get(t *testing.T) {
	ctx := waitContext(t, 180*time.Second)
	client, space, filesystem, createOpts := setupNFSFilesystem(t, "fixtures/TestNFSFilesystem_Get")
	verifyNFSFilesystemBasics(t, space, filesystem, createOpts)

	// Wait for NFS Filesystem status to be 'active'
	_, err := client.WaitForNFSFilesystemStatus(
		ctx,
		space.ID,
		filesystem.ID,
		linodego.NFSFilesystemStatusActive,
	)
	require.NoErrorf(t, err, "Failed to wait for Filesystem status to be active: %s", err)

	filesystem, err = client.GetNFSFilesystemInSpace(context.Background(), space.ID, filesystem.ID)
	require.NoErrorf(t, err, "Error retrieving NFS Filesystem in Space: %v", err)
	verifyNFSFilesystemDetails(t, filesystem, createOpts)
}
