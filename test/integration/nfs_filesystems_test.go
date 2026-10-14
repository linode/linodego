package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getNFSFilesystemCreateOptions(t *testing.T, client *linodego.Client) linodego.NFSFilesystemCreateOptions {
	return linodego.NFSFilesystemCreateOptions{
		Label:            "go-test-nfs-filesystem-" + randLabel(),
		Region:           getRegionsWithCaps(t, client, []linodego.RegionCapability{linodego.CapabilityNFSStorage})[0],
		MaxCapacityBytes: 1099511627776,
		Tags:             linodego.Pointer([]string{"testing"}),
		MaxFileCount:     linodego.Pointer(int64(1500)),
		ProtocolVersions: linodego.Pointer([]linodego.NFSProtocolVersion{linodego.NFSProtocolVersionV4}),
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
	requireActive bool,
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

	// Wait for NFS Filesystem status to be 'active' to read the details (if needed)
	if requireActive {
		ctx := waitContext(t, 240*time.Second)
		_, err = client.WaitForNFSFilesystemStatus(
			ctx,
			space.ID,
			filesystem.ID,
			linodego.NFSFilesystemStatusActive,
		)
		require.NoErrorf(t, err, "Failed to wait for Filesystem status to be active: %s", err)
	}

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

func setupNFSFilesystemCustomQuota(
	t *testing.T,
	client *linodego.Client,
	space *linodego.NFSSpace,
	filesystem *linodego.NFSFilesystem,
	modifiers ...func(opts *linodego.NFSQuotaCreateOptions),
) (*linodego.NFSQuota, linodego.NFSQuotaCreateOptions) {
	t.Helper()

	_, mountDir, isFound := strings.Cut(*filesystem.MountTargetFQDN, ":/")
	require.True(t, isFound, "Mount Target FQDN does not contain ':/'")

	createOpts := linodego.NFSQuotaCreateOptions{
		Path:             "/" + strings.Join([]string{mountDir, "testing-dx"}, "/"),
		MaxFileCount:     filesystem.MaxFileCount,
		MaxCapacityBytes: filesystem.MaxCapacityBytes,
	}

	for _, modifier := range modifiers {
		modifier(&createOpts)
	}

	quota, err := client.CreateNFSQuota(context.Background(), space.ID, filesystem.ID, createOpts)
	require.NoErrorf(t, err, "Error creating Filesystem Quota: %v", err)

	t.Cleanup(func() {
		if err = client.DeleteNFSQuota(context.Background(), space.ID, filesystem.ID, quota.ID); err != nil {
			t.Errorf("Error deleting NFS Filesystem Quota: %v", err)
		}
	})

	return quota, createOpts
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
	assert.EqualValues(t, createOpts.MaxFileCount, linodego.Pointer(filesystem.MaxFileCount))
	assert.NotEmpty(t, filesystem.MountTargetIPs)
	assert.NotNil(t, filesystem.MountTargetFQDN)
	assert.NotNil(t, filesystem.Stats)
	// TODO: Specified in API docs, but not available in API itself yet
	// assert.NotNil(t, filesystem.SnapshotUsageBytes)
	// assert.NotNil(t, filesystem.LDAPConfigID)
	// assert.NotNil(t, filesystem.SourceSnapshotID)
}

func verifyNFSFilesystemAccessPolicyUpdate(
	t *testing.T,
	filesystemAccPolicy *linodego.NFSFilesystemAccessPolicy,
	opts linodego.NFSFilesystemAccessPolicyUpdateOptions,
) {
	assert.True(t, filesystemAccPolicy.Enabled)
	assert.Len(t, filesystemAccPolicy.LinodeACL, 1)
	assert.Equal(t, (*opts.LinodeIDs)[0], filesystemAccPolicy.LinodeACL[0].ID)
	assert.Equal(t, opts.SquashPolicy, linodego.Pointer(filesystemAccPolicy.SquashPolicy))
	assertDateSet(t, filesystemAccPolicy.Updated)
}

func verifyNFSFilesystemQuota(
	t *testing.T,
	quota *linodego.NFSQuota,
	filesystem *linodego.NFSFilesystem,
) {
	assert.Equal(t, filesystem.ID, quota.FilesystemID)
	assert.Equal(t, filesystem.MaxFileCount, *quota.MaxFileCount)
	assert.Equal(t, filesystem.MaxCapacityBytes, *quota.MaxCapacityBytes)
	assert.Equal(t, linodego.NFSQuotaStatusActive, quota.Status)
	assertDateSet(t, quota.Created)
	assertDateSet(t, quota.Updated)
}

func TestNFSFilesystem_Create_smoke(t *testing.T) {
	_, space, filesystem, createOpts := setupNFSFilesystem(t, "fixtures/TestNFSFilesystem_Create", false)
	verifyNFSFilesystemBasics(t, space, filesystem, createOpts)
}

func TestNFSFilesystem_Get(t *testing.T) {
	client, space, filesystem, createOpts := setupNFSFilesystem(t, "fixtures/TestNFSFilesystem_Get", true)

	filesystemGet, err := client.GetNFSFilesystemInSpace(context.Background(), space.ID, filesystem.ID)
	require.NoErrorf(t, err, "Error retrieving NFS Filesystem in Space: %v", err)
	verifyNFSFilesystemBasics(t, space, filesystemGet, createOpts)
	verifyNFSFilesystemDetails(t, filesystemGet, createOpts)
}

func TestNFSFilesystem_GetWithoutNFSSpace(t *testing.T) {
	client, space, filesystem, createOpts := setupNFSFilesystem(t, "fixtures/TestNFSFilesystem_GetWithoutNFSSpace", false)

	filesystemGet, err := client.GetNFSFilesystem(context.Background(), filesystem.ID)
	require.NoErrorf(t, err, "Error retrieving NFS Filesystem with no Space: %v", err)
	verifyNFSFilesystemBasics(t, space, filesystemGet, createOpts)
}

func TestNFSFilesystem_List(t *testing.T) {
	client, space, filesystem, createOpts := setupNFSFilesystem(t, "fixtures/TestNFSFilesystem_List", false)

	f := linodego.Filter{}
	f.AddField(linodego.Eq, "label", filesystem.Label)
	filter, err := f.MarshalJSON()
	if err != nil {
		t.Fatalf("Failed to marshal filter: %v", err)
	}

	filesystemList, err := client.ListNFSFilesystems(context.Background(), space.ID, &linodego.ListOptions{Filter: string(filter)})
	require.NoErrorf(t, err, "Error listing NFS Spaces: %v", err)
	assert.Len(t, filesystemList, 1)
	verifyNFSFilesystemBasics(t, space, linodego.Pointer(filesystemList[0]), createOpts)
}

func TestNFSFilesystem_Update(t *testing.T) {
	client, space, filesystem, createOpts := setupNFSFilesystem(t, "fixtures/TestNFSFilesystem_Update", false)

	updateOpts := filesystem.GetUpdateOptions()
	updateOpts.Label = linodego.Pointer(filesystem.Label + "-updated")
	updateOpts.Tags = linodego.Pointer([]string{"updated"})
	// updateOpts.MaxCapacityBytes = linodego.Pointer(int64(1073741824)) // TODO: Currently API says it is not editable field
	// updateOpts.MaxCountFile = linodego.Pointer(int64(2500)) // TODO: Currently API says it is not editable field

	filesystemUpdate, err := client.UpdateNFSFilesystem(context.Background(), space.ID, filesystem.ID, updateOpts)
	require.NoErrorf(t, err, "Error updating NFS Space: %v", err)

	verifyOpts := linodego.NFSFilesystemCreateOptions{
		Label:  *updateOpts.Label,
		Region: createOpts.Region,
		Tags:   updateOpts.Tags,
		// MaxCapacityBytes: *updateOpts.MaxCapacityBytes, // TODO: Currently API says it is not editable field
		// MaxFileCount: *updateOpts.MaxCountFile, // TODO: Currently API says it is not editable field
	}
	verifyNFSFilesystemBasics(t, space, filesystemUpdate, verifyOpts)
}

func TestNFSFilesystem_GetAccessPolicy_smoke(t *testing.T) {
	client, space, filesystem, _ := setupNFSFilesystem(t, "fixtures/TestNFSFilesystem_GetAccessPolicy", false)

	filesystemAccPolicy, err := client.GetNFSFilesystemAccessPolicy(context.Background(), space.ID, filesystem.ID)
	require.NoErrorf(t, err, "Error retrieving NFS Filesystem Access Policy: %v", err)
	assert.Equal(t, filesystem.ID, filesystemAccPolicy.FilesystemID)
	assert.Equal(t, filesystem.Label, filesystemAccPolicy.Label)
	assert.False(t, filesystemAccPolicy.Enabled)
	// assert.False(t, filesystemAccPolicy.Inherit) // TODO: No Inherit field in API spec, but returned by API response
	assert.NotEmpty(t, filesystemAccPolicy.Protocols)
	assert.Empty(t, filesystemAccPolicy.LinodeACL)
	assert.NotNil(t, filesystemAccPolicy.SquashPolicy)
	assert.Equal(t, linodego.NFSAccessPolicyStatusActive, filesystemAccPolicy.Status)
	assertDateSet(t, filesystemAccPolicy.Created)
	assert.Nil(t, filesystemAccPolicy.Updated)
}

func TestNFSFilesystem_UpdateAccessPolicy(t *testing.T) {
	client, space, filesystem, _ := setupNFSFilesystem(t, "fixtures/TestNFSFilesystem_UpdateAccessPolicy", false)

	linode, err := createInstance(
		t,
		client,
		true,
		func(client *linodego.Client, opts *linodego.InstanceCreateOptions) {
			opts.Region = filesystem.Region
			opts.Booted = linodego.Pointer(true)
		},
	)
	require.NoErrorf(t, err, "Error creating Linode for NFS Filesystem Access Policy: %v", err)

	t.Cleanup(func() {
		if err := client.DeleteInstance(context.Background(), linode.ID); err != nil {
			if t != nil {
				t.Errorf("Error deleting Linode instance: %s", err)
			}
		}
	})

	updateOpts := linodego.NFSFilesystemAccessPolicyUpdateOptions{
		Enabled:      linodego.Pointer(true),
		LinodeIDs:    linodego.Pointer([]int{linode.ID}),
		SquashPolicy: linodego.Pointer(linodego.NFSSquashPolicyAllSquash),
	}

	filesystemAccPolicyUpdate, err := client.UpdateNFSFilesystemAccessPolicy(
		context.Background(),
		space.ID,
		filesystem.ID,
		updateOpts,
	)
	require.NoErrorf(t, err, "Error updating NFS Filesystem Access Policy: %v", err)
	verifyNFSFilesystemAccessPolicyUpdate(t, filesystemAccPolicyUpdate, updateOpts)

	// Wait for NFS Filesystem Access Policy status to be 'active' to read the details
	ctx := waitContext(t, 180*time.Second)
	_, err = client.WaitForNFSFilesystemAccessPolicyStatus(
		ctx,
		space.ID,
		filesystem.ID,
		linodego.NFSAccessPolicyStatusActive,
	)
	require.NoErrorf(t, err, "Failed to wait for Filesystem Access Policy status: %s", err)

	filesystemAccPolicyGet, err := client.GetNFSFilesystemAccessPolicy(context.Background(), space.ID, filesystem.ID)
	require.NoErrorf(t, err, "Error getting NFS Filesystem Access Policy: %v", err)
	verifyNFSFilesystemAccessPolicyUpdate(t, filesystemAccPolicyGet, updateOpts)
}

func TestNFSFilesystem_GetRootQuota_smoke(t *testing.T) {
	client, space, filesystem, _ := setupNFSFilesystem(t, "fixtures/TestNFSFilesystem_GetRootQuota_smoke", true)

	filesystem, err := client.GetNFSFilesystemInSpace(context.Background(), space.ID, filesystem.ID)
	require.NoErrorf(t, err, "Error retrieving NFS Filesystem in Space: %v", err)

	_, mountDir, isFound := strings.Cut(*filesystem.MountTargetFQDN, ":/")
	require.True(t, isFound, "Mount Target FQDN does not contain ':/'")

	quotas, err := client.ListNFSQuotas(context.Background(), space.ID, filesystem.ID, nil)
	require.NoErrorf(t, err, "Error listing Filesystem Quotas: %v", err)
	assert.Len(t, quotas, 1)
	assert.Contains(t, quotas[0].Path, mountDir)
	assert.True(t, quotas[0].IsRoot)
	verifyNFSFilesystemQuota(t, linodego.Pointer(quotas[0]), filesystem)
}

func TestNFSFilesystem_ManageCustomQuota(t *testing.T) {
	t.Skip("Removal of Filesystem with custom Quota requires manual action on mounted Linode instance")
	client, space, filesystem, _ := setupNFSFilesystem(t, "fixtures/TestNFSFilesystem_ManageCustomQuota", true)

	filesystem, err := client.GetNFSFilesystemInSpace(context.Background(), space.ID, filesystem.ID)
	require.NoErrorf(t, err, "Error retrieving NFS Filesystem in Space: %v", err)

	quota, createOpts := setupNFSFilesystemCustomQuota(t, client, space, filesystem)
	assert.Equal(t, createOpts.Path, quota.Path)
	assert.False(t, quota.IsRoot)
	verifyNFSFilesystemQuota(t, quota, filesystem)

	quota, err = client.GetNFSQuota(context.Background(), space.ID, filesystem.ID, quota.ID)
	require.NoErrorf(t, err, "Error retrieving Filesystem Quota: %v", err)
	assert.Equal(t, createOpts.Path, quota.Path)
	assert.False(t, quota.IsRoot)
	verifyNFSFilesystemQuota(t, quota, filesystem)

	quotas, err := client.ListNFSQuotas(context.Background(), space.ID, filesystem.ID, nil)
	require.NoErrorf(t, err, "Error listing Filesystem Quotas: %v", err)
	assert.Len(t, quotas, 2)
}
