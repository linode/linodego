package unit

import (
	"context"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestNFSFilesystem_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_filesystems_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nfs/spaces/123/filesystems", fixtureData)

	filesystems, err := base.Client.ListNFSFilesystems(context.Background(), 123, &linodego.ListOptions{})
	assert.NoError(t, err)
	if assert.Len(t, filesystems, 1) {
		assertNFSFilesystem(t, &filesystems[0])
		assert.NotNil(t, filesystems[0].Stats.CollectedAt)
	}
}

func TestNFSFilesystem_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_filesystem")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nfs/spaces/123/filesystems/456", fixtureData)

	filesystem, err := base.Client.GetNFSFilesystem(context.Background(), 123, 456)
	assert.NoError(t, err)
	assertNFSFilesystem(t, filesystem)
	assert.Nil(t, filesystem.SnapshotUsageBytes)
	assert.Nil(t, filesystem.Stats.CollectedAt)
}

func TestNFSFilesystem_GetByID(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_filesystem")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nfs/filesystems/456", fixtureData)

	filesystem, err := base.Client.GetNFSFilesystemByID(context.Background(), 456)
	assert.NoError(t, err)
	assertNFSFilesystem(t, filesystem)
}

func TestNFSFilesystem_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_filesystem")
	assert.NoError(t, err)

	client := createMockClient(t)
	opts := linodego.NFSFilesystemCreateOptions{
		Label:            "training-data",
		Region:           "us-east",
		MaxCapacityBytes: 1099511627776,
		ProtocolVersions: linodego.Pointer([]linodego.NFSProtocolVersion{linodego.NFSProtocolVersionV4}),
		Tags:             linodego.Pointer([]string{"production"}),
	}

	httpmock.RegisterRegexpResponder(
		"POST",
		mockRequestURL(t, "nfs/spaces/123/filesystems"),
		mockRequestBodyValidate(t, opts, fixtureData),
	)

	filesystem, err := client.CreateNFSFilesystem(context.Background(), 123, opts)
	assert.NoError(t, err)
	assertNFSFilesystem(t, filesystem)
}

func TestNFSFilesystem_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_filesystem")
	assert.NoError(t, err)

	client := createMockClient(t)
	opts := linodego.NFSFilesystemUpdateOptions{
		Label:            linodego.Pointer("training-data"),
		MaxCapacityBytes: linodego.Pointer[int64](1099511627776),
		Tags:             linodego.Pointer([]string{"production"}),
	}

	httpmock.RegisterRegexpResponder(
		"PUT",
		mockRequestURL(t, "nfs/spaces/123/filesystems/456"),
		mockRequestBodyValidate(t, opts, fixtureData),
	)

	filesystem, err := client.UpdateNFSFilesystem(context.Background(), 123, 456, opts)
	assert.NoError(t, err)
	assertNFSFilesystem(t, filesystem)
}

func TestNFSFilesystem_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("nfs/spaces/123/filesystems/456", nil)

	assert.NoError(t, base.Client.DeleteNFSFilesystem(context.Background(), 123, 456))
}

func assertNFSFilesystem(t *testing.T, filesystem *linodego.NFSFilesystem) {
	t.Helper()
	if !assert.NotNil(t, filesystem) {
		return
	}

	assert.Equal(t, 456, filesystem.ID)
	assert.Equal(t, 123, filesystem.SpaceID)
	assert.Equal(t, "training-data", filesystem.Label)
	assert.Equal(t, "us-east", filesystem.Region)
	assert.Equal(t, []linodego.NFSProtocolVersion{linodego.NFSProtocolVersionV4}, filesystem.ProtocolVersions)
	assert.Equal(t, linodego.NFSFilesystemStatusActive, filesystem.Status)
	if assert.NotNil(t, filesystem.MountTargetFQDN) {
		assert.Equal(t, "production-space-7b.nfs.us-east.linode.com:/training-data-1c8", *filesystem.MountTargetFQDN)
	}
	assert.Equal(t, int64(1099511627776), filesystem.MaxCapacityBytes)
	assert.Equal(t, []string{"production"}, filesystem.Tags)
	assert.Equal(t, "2026-05-01T14:00:00Z", filesystem.Created.Format(time.RFC3339))
}
