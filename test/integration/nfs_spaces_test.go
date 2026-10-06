package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var TestNFSSpaceCreateOptions = linodego.NFSSpaceCreateOptions{
	Label:       "go-test-nfs-space-" + randLabel(),
	Description: linodego.Pointer("Test description"),
	Tags:        linodego.Pointer([]string{"testing"}),
}

func createNFSSpace(
	t *testing.T,
	fixtureYaml string,
	modifiers ...func(opts *linodego.NFSSpaceCreateOptions),
) (*linodego.Client, *linodego.NFSSpace) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixtureYaml)
	createOpts := TestNFSSpaceCreateOptions

	for _, modifier := range modifiers {
		modifier(&createOpts)
	}

	space, err := client.CreateNFSSpace(context.Background(), createOpts)
	require.NoErrorf(t, err, "Error creating NFS Space: %v", err)

	t.Cleanup(func() {
		if err := client.DeleteNFSSpace(context.Background(), space.ID); err != nil {
			t.Errorf("Error deleting NFS Space: %v", err)
		}
		fixtureTeardown()
	})

	return client, space
}

func verifyNFSSpace(t *testing.T, space *linodego.NFSSpace) {
	t.Helper()
	assert.Equal(t, TestNFSSpaceCreateOptions.Label, space.Label)
	assert.Equal(t, TestNFSSpaceCreateOptions.Description, space.Description)
	assert.Equal(t, TestNFSSpaceCreateOptions.Tags, linodego.Pointer(space.Tags))
	assert.NotNil(t, space.Status)
	assertDateSet(t, space.Created)
	assertDateSet(t, space.Updated)
}

func TestNFSSpace_Create_smoke(t *testing.T) {
	client, space := createNFSSpace(t, "fixtures/TestNFSSpace_Create")
	verifyNFSSpace(t, space)

	spaceGet, err := client.GetNFSSpace(context.Background(), space.ID)
	require.NoErrorf(t, err, "Error retrieving NFS Space: %v", err)
	verifyNFSSpace(t, spaceGet)
}

func TestNFSSpace_List(t *testing.T) {
	client, space := createNFSSpace(t, "fixtures/TestNFSSpace_List")

	f := linodego.Filter{}
	f.AddField(linodego.Eq, "label", space.Label)
	filter, err := f.MarshalJSON()
	if err != nil {
		t.Fatalf("Failed to marshal filter: %v", err)
	}

	spaceList, err := client.ListNFSSpaces(context.Background(), &linodego.ListOptions{Filter: string(filter)})
	require.NoErrorf(t, err, "Error listing NFS Spaces: %v", err)
	assert.Len(t, spaceList, 1)
	verifyNFSSpace(t, linodego.Pointer(spaceList[0]))
}

func TestNFSSpace_Update(t *testing.T) {
	client, space := createNFSSpace(t, "fixtures/TestNFSSpace_Update")
	updateOpts := linodego.NFSSpaceUpdateOptions{
		Label:       linodego.Pointer(space.Label + "-updated"),
		Description: linodego.Pointer("Description updated"),
		Tags:        linodego.Pointer([]string{"updated"}),
	}

	spaceUpdate, err := client.UpdateNFSSpace(context.Background(), space.ID, updateOpts)
	require.NoErrorf(t, err, "Error updating NFS Space: %v", err)
	assert.Equal(t, updateOpts.Label, linodego.Pointer(spaceUpdate.Label))
	assert.Equal(t, updateOpts.Description, spaceUpdate.Description)
	assert.Equal(t, updateOpts.Tags, linodego.Pointer(spaceUpdate.Tags))
}
