package unit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestImageShareGroup_Consumer_CreateToken(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_consumer_create_token")

	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.ImageShareGroupCreateTokenOptions{
		ValidForShareGroupUUID: "e3407945-5946-40f9-9732-d3c58b131ec0",
		Label:                  linodego.Pointer("my_token"),
	}

	base.MockPost("images/sharegroups/tokens", fixtureData)

	token, err := base.Client.ImageShareGroupCreateToken(context.Background(), requestData)
	assert.NoError(t, err)

	assert.Equal(t, "singleusetoken", token.Token)
	assert.Equal(t, "18db04bf-fd0f-4bf6-944a-1fc2ae044dc6", token.TokenUUID)
	assert.Equal(t, "active", token.Status)
	assert.Equal(t, "my_token", token.Label)
	assert.Equal(t, "e3407945-5946-40f9-9732-d3c58b131ec0", token.ValidForShareGroupUUID)
	assert.Nil(t, token.ShareGroupUUID)
	assert.Nil(t, token.ShareGroupLabel)
	assert.Equal(t, "2025-07-01T04:00:00Z", token.Created.Format(time.RFC3339))
	assert.Equal(t, "2025-07-01T04:00:01Z", token.Updated.Format(time.RFC3339))
	assert.Equal(t, "2025-07-01T04:00:02Z", token.Expiry.Format(time.RFC3339))
}

func TestImageShareGroup_Consumer_UpdateToken(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_consumer_update_token")

	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.ImageShareGroupUpdateTokenOptions{
		Label: "my_updated_token",
	}

	base.MockPut("images/sharegroups/tokens/18db04bf-fd0f-4bf6-944a-1fc2ae044dc6", fixtureData)

	token, err := base.Client.ImageShareGroupUpdateToken(context.Background(), "18db04bf-fd0f-4bf6-944a-1fc2ae044dc6", requestData)
	assert.NoError(t, err)

	assert.Equal(t, "18db04bf-fd0f-4bf6-944a-1fc2ae044dc6", token.TokenUUID)
	assert.Equal(t, "active", token.Status)
	assert.Equal(t, "my_updated_token", token.Label)
	assert.Equal(t, "e3407945-5946-40f9-9732-d3c58b131ec0", token.ValidForShareGroupUUID)
	assert.Nil(t, token.ShareGroupUUID)
	assert.Nil(t, token.ShareGroupLabel)
	assert.Equal(t, "2025-07-01T04:00:00Z", token.Created.Format(time.RFC3339))
	assert.Equal(t, "2025-07-01T04:00:01Z", token.Updated.Format(time.RFC3339))
	assert.Equal(t, "2025-07-01T04:00:02Z", token.Expiry.Format(time.RFC3339))
}

func TestImageShareGroup_Consumer_List_Tokens(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_consumer_list_tokens")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("images/sharegroups/tokens", fixtureData)

	tokens, err := base.Client.ImageShareGroupListTokens(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	token := tokens[0]

	assert.Len(t, tokens, 1)

	assert.Equal(t, "18db04bf-fd0f-4bf6-944a-1fc2ae044dc6", token.TokenUUID)
	assert.Equal(t, "active", token.Status)
	assert.Equal(t, "my_token", token.Label)
	assert.Equal(t, "e3407945-5946-40f9-9732-d3c58b131ec0", token.ValidForShareGroupUUID)
	assert.Equal(t, "e3407945-5946-40f9-9732-d3c58b131ec0", *token.ShareGroupUUID)
	assert.Equal(t, "my_sharegroup", *token.ShareGroupLabel)
	assert.Equal(t, "2025-07-01T04:00:00Z", token.Created.Format(time.RFC3339))
	assert.Equal(t, "2025-07-01T04:00:01Z", token.Updated.Format(time.RFC3339))
	assert.Equal(t, "2025-07-01T04:00:02Z", token.Expiry.Format(time.RFC3339))
}

func TestImageShareGroup_Consumer_Get_Token(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_consumer_get_token")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("images/sharegroups/tokens/18db04bf-fd0f-4bf6-944a-1fc2ae044dc6", fixtureData)

	token, err := base.Client.ImageShareGroupGetToken(context.Background(), "18db04bf-fd0f-4bf6-944a-1fc2ae044dc6")
	assert.NoError(t, err)

	assert.Equal(t, "18db04bf-fd0f-4bf6-944a-1fc2ae044dc6", token.TokenUUID)
	assert.Equal(t, "active", token.Status)
	assert.Equal(t, "my_token", token.Label)
	assert.Equal(t, "e3407945-5946-40f9-9732-d3c58b131ec0", token.ValidForShareGroupUUID)
	assert.Equal(t, "e3407945-5946-40f9-9732-d3c58b131ec0", *token.ShareGroupUUID)
	assert.Equal(t, "my_sharegroup", *token.ShareGroupLabel)
	assert.Equal(t, "2025-07-01T04:00:00Z", token.Created.Format(time.RFC3339))
	assert.Equal(t, "2025-07-01T04:00:01Z", token.Updated.Format(time.RFC3339))
	assert.Equal(t, "2025-07-01T04:00:02Z", token.Expiry.Format(time.RFC3339))
}

func TestImageShareGroup_Consumer_Get_ShareGroup_ByToken(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_consumer_get_sharegroup_by_token")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("images/sharegroups/tokens/18db04bf-fd0f-4bf6-944a-1fc2ae044dc6/sharegroup", fixtureData)

	sharegroup, err := base.Client.ImageShareGroupGetByToken(context.Background(), "18db04bf-fd0f-4bf6-944a-1fc2ae044dc6")
	assert.NoError(t, err)

	assert.Equal(t, 1, sharegroup.ID)
	assert.Equal(t, "967913ad-9379-4039-b166-31b6b1440019", sharegroup.UUID)
	assert.Equal(t, "new_sharegroup_for_testing", sharegroup.Label)
	assert.Equal(t, "my description.", sharegroup.Description)
	assert.Equal(t, true, sharegroup.IsSuspended)
	assert.Equal(t, "2025-07-21T20:18:37Z", sharegroup.Created.Format(time.RFC3339))
	assert.Equal(t, "2025-07-22T18:09:07Z", sharegroup.Updated.Format(time.RFC3339))
}

func TestImageShareGroup_Consumer_Get_ShareGroup_Images_ByToken(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_consumer_get_sharegroup_images_by_token")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("images/sharegroups/tokens/18db04bf-fd0f-4bf6-944a-1fc2ae044dc6/sharegroup/images", fixtureData)

	images, err := base.Client.ImageShareGroupGetImagesByToken(context.Background(), "18db04bf-fd0f-4bf6-944a-1fc2ae044dc6", &linodego.ListOptions{})
	assert.NoError(t, err)

	image := images[0]

	assert.Len(t, images, 1)

	assert.Equal(t, "shared/123", image.ID)
	assert.Equal(t, "producer defined share label", image.Label)
	assert.Equal(t, "some blob of legal text", image.Description)
	assert.Equal(t, "2024-12-03T01:51:24Z", image.Created.Format(time.RFC3339))
	assert.Equal(t, "2024-12-03T01:51:24Z", image.Updated.Format(time.RFC3339))
	assert.Equal(t, 1761, image.Size)
	assert.Nil(t, image.CreatedBy)
	assert.Equal(t, "manual", image.Type)
	assert.Empty(t, image.Tags)
	assert.False(t, image.IsPublic)
	assert.Nil(t, image.IsShared)
	assert.False(t, image.Deprecated)
	assert.Nil(t, image.Vendor)
	assert.Nil(t, image.Expiry)
	assert.Nil(t, image.EOL)
	assert.Equal(t, linodego.ImageStatus("available"), image.Status)
	assert.Equal(t, []string{"distributed-sites"}, image.Capabilities)
	assert.Equal(t, "us-ord", image.Regions[0].Region)
	assert.Equal(t, linodego.ImageRegionStatus("available"), image.Regions[0].Status)
	assert.Equal(t, 222, image.TotalSize)
	assert.Nil(t, image.ImageSharing.SharedWith)
	assert.Equal(t, 1234, image.ImageSharing.SharedBy.ShareGroupID)
	assert.Equal(t, "0ee8e1c1-b19b-4052-9487-e3b13faac111", image.ImageSharing.SharedBy.ShareGroupUUID)
	assert.Equal(t, "test-group-minecraft-1", image.ImageSharing.SharedBy.ShareGroupLabel)
	assert.Nil(t, image.ImageSharing.SharedBy.SourceImageID)
}

func TestImageShareGroup_Consumer_RemoveToken(t *testing.T) {
	client := createMockClient(t)

	tokenUUID := "18db04bf-fd0f-4bf6-944a-1fc2ae044dc6"

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, fmt.Sprintf("images/sharegroups/tokens/%s", tokenUUID)),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.ImageShareGroupRemoveToken(context.Background(), tokenUUID); err != nil {
		t.Fatal(err)
	}
}
