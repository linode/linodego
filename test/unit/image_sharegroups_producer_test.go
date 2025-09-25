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

func TestImageShareGroup_Producer_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_producer_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.ImageShareGroupCreateOptions{
		Label:       "a-cool-label",
		Description: linodego.Pointer("This is the description."),
		Images: []linodego.ImageShareGroupImage{
			{
				ID:          "linode/debian11",
				Label:       linodego.Pointer("image-label"),
				Description: linodego.Pointer("A description."),
			},
		},
	}

	base.MockPost("images/sharegroups", fixtureData)

	imageShareGroup, err := base.Client.CreateImageShareGroup(context.Background(), requestData)
	assert.NoError(t, err)

	assert.Equal(t, 1234, imageShareGroup.ID)
	assert.Equal(t, "f47ac10b-58cc-4372-a567-0e02b2c3d479", imageShareGroup.UUID)
	assert.Equal(t, "a-cool-label", imageShareGroup.Label)
	assert.Equal(t, "This is the description.", imageShareGroup.Description)
	assert.Equal(t, false, imageShareGroup.IsSuspended)
	assert.Equal(t, 1, imageShareGroup.ImagesCount)
	assert.Equal(t, 0, imageShareGroup.MembersCount)
	assert.Equal(t, "2025-07-01T04:00:00Z", imageShareGroup.Created.Format(time.RFC3339))
	assert.Equal(t, "2025-07-02T04:00:00Z", imageShareGroup.Updated.Format(time.RFC3339))
	assert.Nil(t, imageShareGroup.Expiry)
}

func TestImageShareGroup_Producer_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_producer_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.ImageShareGroupUpdateOptions{
		Label:       linodego.Pointer("a-cool-updated-label"),
		Description: linodego.Pointer("This is the updated description."),
	}

	base.MockPut("images/sharegroups/1234", fixtureData)

	imageShareGroup, err := base.Client.UpdateImageShareGroup(context.Background(), 1234, requestData)
	assert.NoError(t, err)

	assert.Equal(t, 1234, imageShareGroup.ID)
	assert.Equal(t, "f47ac10b-58cc-4372-a567-0e02b2c3d479", imageShareGroup.UUID)
	assert.Equal(t, "a-cool-updated-label", imageShareGroup.Label)
	assert.Equal(t, "This is the updated description.", imageShareGroup.Description)
	assert.Equal(t, false, imageShareGroup.IsSuspended)
	assert.Equal(t, 1, imageShareGroup.ImagesCount)
	assert.Equal(t, 0, imageShareGroup.MembersCount)
	assert.Equal(t, "2025-07-01T04:00:00Z", imageShareGroup.Created.Format(time.RFC3339))
	assert.Equal(t, "2025-07-02T04:00:00Z", imageShareGroup.Updated.Format(time.RFC3339))
	assert.Nil(t, imageShareGroup.Expiry)
}

func TestImageShareGroup_Producer_AddImages(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_producer_add_images")

	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.ImageShareGroupAddImagesOptions{
		Images: []linodego.ImageShareGroupImage{
			{
				ID:          "linode/debian11",
				Label:       linodego.Pointer("image-label"),
				Description: linodego.Pointer("A description."),
			},
		},
	}

	base.MockPost("images/sharegroups/1234/images", fixtureData)

	images, err := base.Client.ImageShareGroupAddImages(context.Background(), 1234, requestData)
	assert.NoError(t, err)

	image := images[0]

	assert.Equal(t, "linode/debian11", image.ID)
	assert.Equal(t, "Debian 11", image.Label)
	assert.Equal(t, "Example image description.", image.Description)
	assert.Nil(t, image.Vendor)
	assert.Equal(t, true, image.IsPublic)
	assert.Equal(t, false, image.Deprecated)
	assert.Equal(t, "available", string(image.Status))
	assert.Equal(t, "2021-08-14T22:44:02Z", image.Created.Format(time.RFC3339))
	assert.Equal(t, "2021-08-14T22:44:02Z", image.Updated.Format(time.RFC3339))
	assert.Equal(t, "2026-07-01T04:00:00Z", image.EOL.Format(time.RFC3339))
	assert.Equal(t, 2500, image.Size)
	assert.Equal(t, 1234567, image.TotalSize)

	assert.Equal(t, 1234, image.ImageSharing.SharedBy.ShareGroupID)
	assert.Equal(t, "0ee8e1c1-b19b-4052-9487-e3b13faac111", image.ImageSharing.SharedBy.ShareGroupUUID)
	assert.Equal(t, "test-group-minecraft-1", image.ImageSharing.SharedBy.ShareGroupLabel)
	assert.Nil(t, image.ImageSharing.SharedBy.SourceImageID)
	assert.Nil(t, image.ImageSharing.SharedWith)

	assert.ElementsMatch(t, []string{"cloud-init", "distributed-sites"}, image.Capabilities)

	assert.Len(t, image.Regions, 1)
	assert.Equal(t, "us-iad", image.Regions[0].Region)
	assert.Equal(t, "available", string(image.Regions[0].Status))

	assert.ElementsMatch(t, []string{"repair-image", "fix-1"}, image.Tags)
}

func TestImageShareGroup_Producer_UpdateImage(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_producer_update_image")

	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.ImageShareGroupUpdateImageOptions{
		Description: linodego.Pointer("Example updated image description."),
	}

	base.MockPut("images/sharegroups/1234/images/123", fixtureData)

	image, err := base.Client.ImageShareGroupUpdateImage(context.Background(), 1234, "123", requestData)
	assert.NoError(t, err)

	assert.Equal(t, "linode/debian11", image.ID)
	assert.Equal(t, "Debian 11", image.Label)
	assert.Equal(t, "Example updated image description.", image.Description)
	assert.Nil(t, image.Vendor)
	assert.Equal(t, true, image.IsPublic)
	assert.Equal(t, false, image.Deprecated)
	assert.Equal(t, "available", string(image.Status))
	assert.Equal(t, "2021-08-14T22:44:02Z", image.Created.Format(time.RFC3339))
	assert.Equal(t, "2021-08-14T22:44:02Z", image.Updated.Format(time.RFC3339))
	assert.Equal(t, "2026-07-01T04:00:00Z", image.EOL.Format(time.RFC3339))
	assert.Equal(t, 2500, image.Size)
	assert.Equal(t, 1234567, image.TotalSize)

	assert.Equal(t, 1234, image.ImageSharing.SharedBy.ShareGroupID)
	assert.Equal(t, "0ee8e1c1-b19b-4052-9487-e3b13faac111", image.ImageSharing.SharedBy.ShareGroupUUID)
	assert.Equal(t, "test-group-minecraft-1", image.ImageSharing.SharedBy.ShareGroupLabel)
	assert.Nil(t, image.ImageSharing.SharedBy.SourceImageID)
	assert.Nil(t, image.ImageSharing.SharedWith)

	assert.ElementsMatch(t, []string{"cloud-init", "distributed-sites"}, image.Capabilities)

	assert.Len(t, image.Regions, 1)
	assert.Equal(t, "us-iad", image.Regions[0].Region)
	assert.Equal(t, "available", string(image.Regions[0].Status))

	assert.ElementsMatch(t, []string{"repair-image", "fix-1"}, image.Tags)
}

func TestImageShareGroup_Producer_AddMember(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_producer_add_member")

	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.ImageShareGroupAddMemberOptions{
		Token: "my-token",
		Label: "CompanyEshop image sharing",
	}

	base.MockPost("images/sharegroups/1234/members", fixtureData)

	member, err := base.Client.ImageShareGroupAddMember(context.Background(), 1234, requestData)
	assert.NoError(t, err)

	assert.Equal(t, "24wef-243qg-45wgg-q343q", member.TokenUUID)
	assert.Equal(t, "active", member.Status)
	assert.Equal(t, "CompanyEshop image sharing", member.Label)
	assert.Equal(t, "2016-03-16T17:30:49Z", member.Created.Format(time.RFC3339))
	assert.Equal(t, "2016-03-18T17:30:49Z", member.Updated.Format(time.RFC3339))
	assert.Equal(t, "2016-03-18T17:30:50Z", member.Expiry.Format(time.RFC3339))
}

func TestImageShareGroup_Producer_List_ContainingPrivateImage(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroups_producer_list_containing_private_image")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("images/private%2F1234/sharegroups", fixtureData)

	imageShareGroups, err := base.Client.ListImageShareGroupsContainingPrivateImage(context.Background(), "private/1234", &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.Len(t, imageShareGroups, 2)

	// First share group
	assert.Equal(t, 456, imageShareGroups[0].ID)
	assert.Equal(t, "07358d70-0f95-4bd1-b2fd-fd7ab052baf4", imageShareGroups[0].UUID)
	assert.Equal(t, "some share group label", imageShareGroups[0].Label)
	assert.Equal(t, "some desc", imageShareGroups[0].Description)
	assert.False(t, imageShareGroups[0].IsSuspended)
	assert.Equal(t, 12, imageShareGroups[0].MembersCount)
	assert.Equal(t, 3, imageShareGroups[0].ImagesCount)
	assert.Equal(t, "2026-03-16T17:30:49Z", imageShareGroups[0].Created.Format(time.RFC3339))
	assert.Equal(t, "2026-04-16T17:30:49Z", imageShareGroups[0].Updated.Format(time.RFC3339))
	assert.Nil(t, imageShareGroups[0].Expiry)

	// Second share group
	assert.Equal(t, 457, imageShareGroups[1].ID)
	assert.Equal(t, "eb5b9f0f-2e70-46a6-aee4-a081a2b99699", imageShareGroups[1].UUID)
	assert.Equal(t, "some share other group label", imageShareGroups[1].Label)
	assert.Equal(t, "some_desc", imageShareGroups[1].Description)
	assert.False(t, imageShareGroups[1].IsSuspended)
	assert.Equal(t, 1, imageShareGroups[1].MembersCount)
	assert.Equal(t, 1, imageShareGroups[1].ImagesCount)
	assert.Equal(t, "2026-01-16T17:30:49Z", imageShareGroups[1].Created.Format(time.RFC3339))
	assert.Equal(t, "2026-09-16T17:30:49Z", imageShareGroups[1].Updated.Format(time.RFC3339))
	assert.Nil(t, imageShareGroups[1].Expiry)
}

func TestImageShareGroup_Producer_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroups_producer_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("images/sharegroups", fixtureData)

	imageShareGroups, err := base.Client.ListImageShareGroups(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.Len(t, imageShareGroups, 2)

	// First share group
	assert.Equal(t, 456, imageShareGroups[0].ID)
	assert.Equal(t, "07358d70-0f95-4bd1-b2fd-fd7ab052baf4", imageShareGroups[0].UUID)
	assert.Equal(t, "some share group label", imageShareGroups[0].Label)
	assert.Equal(t, "some desc", imageShareGroups[0].Description)
	assert.False(t, imageShareGroups[0].IsSuspended)
	assert.Equal(t, 12, imageShareGroups[0].MembersCount)
	assert.Equal(t, 3, imageShareGroups[0].ImagesCount)
	assert.Equal(t, "2026-03-16T17:30:49Z", imageShareGroups[0].Created.Format(time.RFC3339))
	assert.Equal(t, "2026-04-16T17:30:49Z", imageShareGroups[0].Updated.Format(time.RFC3339))
	assert.Nil(t, imageShareGroups[0].Expiry)

	// Second share group
	assert.Equal(t, 457, imageShareGroups[1].ID)
	assert.Equal(t, "eb5b9f0f-2e70-46a6-aee4-a081a2b99699", imageShareGroups[1].UUID)
	assert.Equal(t, "some share other group label", imageShareGroups[1].Label)
	assert.Equal(t, "some_desc", imageShareGroups[1].Description)
	assert.False(t, imageShareGroups[1].IsSuspended)
	assert.Equal(t, 1, imageShareGroups[1].MembersCount)
	assert.Equal(t, 1, imageShareGroups[1].ImagesCount)
	assert.Equal(t, "2026-01-16T17:30:49Z", imageShareGroups[1].Created.Format(time.RFC3339))
	assert.Equal(t, "2026-09-16T17:30:49Z", imageShareGroups[1].Updated.Format(time.RFC3339))
	assert.Nil(t, imageShareGroups[1].Expiry)
}

func TestImageShareGroup_Producer_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_producer_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet(formatMockAPIPath("images/sharegroups/%d", 457), fixtureData)

	imageShareGroup, err := base.Client.GetImageShareGroup(context.Background(), 457)
	assert.NoError(t, err)

	assert.Equal(t, 457, imageShareGroup.ID)
	assert.Equal(t, "eb5b9f0f-2e70-46a6-aee4-a081a2b99699", imageShareGroup.UUID)
	assert.Equal(t, "some share other group label", imageShareGroup.Label)
	assert.Equal(t, "some larger text", imageShareGroup.Description)
	assert.False(t, imageShareGroup.IsSuspended)
	assert.Equal(t, 1, imageShareGroup.MembersCount)
	assert.Equal(t, 1, imageShareGroup.ImagesCount)
	assert.Equal(t, "2026-01-16T17:30:49Z", imageShareGroup.Created.Format(time.RFC3339))
	assert.Equal(t, "2026-09-16T17:30:49Z", imageShareGroup.Updated.Format(time.RFC3339))
	assert.Nil(t, imageShareGroup.Expiry)
}

func TestImageShareGroup_Producer_List_Images(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_producer_list_images")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("images/sharegroups/1234/images", fixtureData)

	images, err := base.Client.ImageShareGroupListImages(context.Background(), 1234, &linodego.ListOptions{})
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

func TestImageShareGroup_Producer_List_Members(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_producer_list_members")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("images/sharegroups/1234/members", fixtureData)

	members, err := base.Client.ImageShareGroupListMembers(context.Background(), 1234, &linodego.ListOptions{})
	assert.NoError(t, err)

	member := members[0]

	assert.Len(t, members, 1)

	assert.Equal(t, "24wef-243qg-45wgg-q343q", member.TokenUUID)
	assert.Equal(t, "active", member.Status)
	assert.Equal(t, "CompanyEshop image sharing", member.Label)
	assert.Equal(t, "2016-03-16T17:30:49Z", member.Created.Format(time.RFC3339))
	assert.Equal(t, "2016-03-18T17:30:49Z", member.Updated.Format(time.RFC3339))
	assert.Equal(t, "2016-03-18T17:30:50Z", member.Expiry.Format(time.RFC3339))
}

func TestImageShareGroup_Producer_Get_Member(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_producer_get_member")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("images/sharegroups/1234/members/24wef-243qg-45wgg-q343q", fixtureData)

	member, err := base.Client.ImageShareGroupGetMember(context.Background(), 1234, "24wef-243qg-45wgg-q343q")
	assert.NoError(t, err)

	assert.Equal(t, "24wef-243qg-45wgg-q343q", member.TokenUUID)
	assert.Equal(t, "active", member.Status)
	assert.Equal(t, "CompanyEshop image sharing", member.Label)
	assert.Equal(t, "2016-03-16T17:30:49Z", member.Created.Format(time.RFC3339))
	assert.Equal(t, "2016-03-18T17:30:49Z", member.Updated.Format(time.RFC3339))
	assert.Equal(t, "2016-03-18T17:30:50Z", member.Expiry.Format(time.RFC3339))
}

func TestImageShareGroup_Producer_UpdateMember(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_sharegroup_producer_update_member")

	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.ImageShareGroupUpdateMemberOptions{
		Label: "CompanyEshop image sharing updated",
	}

	base.MockPut("images/sharegroups/1234/members/24wef-243qg-45wgg-q343q", fixtureData)

	member, err := base.Client.ImageShareGroupUpdateMember(context.Background(), 1234, "24wef-243qg-45wgg-q343q", requestData)
	assert.NoError(t, err)

	assert.Equal(t, "24wef-243qg-45wgg-q343q", member.TokenUUID)
	assert.Equal(t, "active", member.Status)
	assert.Equal(t, "CompanyEshop image sharing updated", member.Label)
	assert.Equal(t, "2016-03-16T17:30:49Z", member.Created.Format(time.RFC3339))
	assert.Equal(t, "2016-03-18T17:30:49Z", member.Updated.Format(time.RFC3339))
	assert.Equal(t, "2016-03-18T17:30:50Z", member.Expiry.Format(time.RFC3339))
}

func TestImageShareGroup_Producer_Delete(t *testing.T) {
	client := createMockClient(t)

	imageShareGroupID := 123

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, fmt.Sprintf("images/sharegroups/%d", imageShareGroupID)),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteImageShareGroup(context.Background(), imageShareGroupID); err != nil {
		t.Fatal(err)
	}
}

func TestImageShareGroup_Producer_RemoveImage(t *testing.T) {
	client := createMockClient(t)

	imageShareGroupID := 123
	imageID := "123"

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, fmt.Sprintf("images/sharegroups/%d/images/%s", imageShareGroupID, imageID)),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.ImageShareGroupRemoveImage(context.Background(), imageShareGroupID, imageID); err != nil {
		t.Fatal(err)
	}
}

func TestImageShareGroup_Producer_RemoveMember(t *testing.T) {
	client := createMockClient(t)

	imageShareGroupID := 123
	tokenUUID := "18db04bf-fd0f-4bf6-944a-1fc2ae044dc6"

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, fmt.Sprintf("images/sharegroups/%d/members/%s", imageShareGroupID, tokenUUID)),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.ImageShareGroupRemoveMember(context.Background(), imageShareGroupID, tokenUUID); err != nil {
		t.Fatal(err)
	}
}
