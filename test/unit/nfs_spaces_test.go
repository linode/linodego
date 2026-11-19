package unit

import (
	"context"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestNFSSpace_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_spaces_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nfs/spaces", fixtureData)

	spaces, err := base.Client.ListNFSSpaces(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	if assert.Len(t, spaces, 1) {
		assertNFSSpace(t, &spaces[0])
	}
}

func TestNFSSpace_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_space")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nfs/spaces/123", fixtureData)

	space, err := base.Client.GetNFSSpace(context.Background(), 123)
	assert.NoError(t, err)
	assertNFSSpace(t, space)
}

func TestNFSSpace_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_space")
	assert.NoError(t, err)

	client := createMockClient(t)
	opts := linodego.NFSSpaceCreateOptions{
		Label:       "production-space",
		Description: linodego.Pointer("Test NFS Space"),
		Tags:        linodego.Pointer([]string{"production"}),
	}

	httpmock.RegisterRegexpResponder(
		"POST",
		mockRequestURL(t, "nfs/spaces"),
		mockRequestBodyValidate(t, opts, fixtureData),
	)

	space, err := client.CreateNFSSpace(context.Background(), opts)
	assert.NoError(t, err)
	assertNFSSpace(t, space)
}

func TestNFSSpace_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_space")
	assert.NoError(t, err)

	client := createMockClient(t)
	opts := linodego.NFSSpaceUpdateOptions{
		Label:       linodego.Pointer("production-space"),
		Description: linodego.Pointer("Test NFS Space"),
		Tags:        linodego.Pointer([]string{"production"}),
	}

	httpmock.RegisterRegexpResponder(
		"PUT",
		mockRequestURL(t, "nfs/spaces/123"),
		mockRequestBodyValidate(t, opts, fixtureData),
	)

	space, err := client.UpdateNFSSpace(context.Background(), 123, opts)
	assert.NoError(t, err)
	assertNFSSpace(t, space)
}

func TestNFSSpace_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("nfs/spaces/123", nil)

	assert.NoError(t, base.Client.DeleteNFSSpace(context.Background(), 123))
}

func assertNFSSpace(t *testing.T, space *linodego.NFSSpace) {
	t.Helper()
	if !assert.NotNil(t, space) {
		return
	}

	assert.Equal(t, 123, space.ID)
	assert.Equal(t, "production-space", space.Label)
	if assert.NotNil(t, space.Description) {
		assert.Equal(t, "Test NFS Space", *space.Description)
	}
	assert.Equal(t, linodego.NFSSpaceStatusActive, space.Status)
	assert.Equal(t, int64(1099511627776), space.MaxCapacityBytes)
	assert.Equal(t, int64(107374182400), space.UsedCapacityBytes)
	assert.Equal(t, []string{"production"}, space.Tags)
	if assert.NotNil(t, space.Created) {
		assert.Equal(t, "2026-05-01T14:00:00Z", space.Created.Format(time.RFC3339))
	}
	if assert.NotNil(t, space.Updated) {
		assert.Equal(t, "2026-05-10T10:00:00Z", space.Updated.Format(time.RFC3339))
	}
}
