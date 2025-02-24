package unit

import (
	"context"
	"fmt"
	"github.com/jarcoal/httpmock"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestImage_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("images_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("images", fixtureData)

	images, err := base.Client.ListImages(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	image := images[0]
	assert.Equal(t, "linode/debian11", image.ID)
	assert.Equal(t, "Debian 11", image.Label)
	assert.Equal(t, "manual", image.Type)
	assert.Equal(t, "Example image description.", image.Description)
	assert.Equal(t, 2500, image.Size)
	assert.Equal(t, linodego.ImageStatus("available"), image.Status)
	assert.Equal(t, true, image.IsPublic)
	assert.Equal(t, "2026-07-01T04:00:00Z", image.EOL.Format(time.RFC3339))
	assert.ElementsMatch(t, []string{"repair-image", "fix-1"}, image.Tags)
	assert.Equal(t, "Debian", image.Vendor)
	assert.False(t, image.Deprecated)
	expectedCapabilities := []string{"cloud-init", "distributed-sites"}
	assert.ElementsMatch(t, expectedCapabilities, image.Capabilities)
	assert.Len(t, image.Regions, 1)
	assert.Equal(t, "us-iad", image.Regions[0].Region)
	assert.Equal(t, linodego.ImageRegionStatus("available"), image.Regions[0].Status)
}

func TestImage_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	imageID := "123"

	base.MockGet(formatMockAPIPath("images/%s", imageID), fixtureData)

	image, err := base.Client.GetImage(context.Background(), imageID)
	assert.NoError(t, err)

	assert.Equal(t, "linode/debian11", image.ID)
	assert.Equal(t, "Debian 11", image.Label)
	assert.Equal(t, "Example image description.", image.Description)
	assert.Equal(t, "Debian", image.Vendor)
	assert.Equal(t, true, image.IsPublic)
	assert.Equal(t, false, image.Deprecated)
	assert.Equal(t, "available", string(image.Status))
	assert.Equal(t, "2021-08-14T22:44:02Z", image.Created.Format(time.RFC3339))
	assert.Equal(t, "2021-08-14T22:44:02Z", image.Updated.Format(time.RFC3339))
	assert.Equal(t, "2026-07-01T04:00:00Z", image.EOL.Format(time.RFC3339))
	assert.Equal(t, 2500, image.Size)
	assert.Equal(t, 1234567, image.TotalSize)

	assert.ElementsMatch(t, []string{"cloud-init", "distributed-sites"}, image.Capabilities)

	assert.Len(t, image.Regions, 1)
	assert.Equal(t, "us-iad", image.Regions[0].Region)
	assert.Equal(t, "available", string(image.Regions[0].Status))

	assert.ElementsMatch(t, []string{"repair-image", "fix-1"}, image.Tags)
}

func TestImage_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.ImageCreateOptions{
		DiskID:      123456,
		Label:       "Debian 11",
		Description: "Example image description.",
		CloudInit:   true,
		Tags:        &[]string{"repair-image", "fix-1"},
	}

	base.MockPost("images", fixtureData)

	image, err := base.Client.CreateImage(context.Background(), requestData)
	assert.NoError(t, err)

	assert.Equal(t, "linode/debian11", image.ID)
	assert.Equal(t, "Debian 11", image.Label)
	assert.Equal(t, "Example image description.", image.Description)
	assert.Equal(t, "Debian", image.Vendor)
	assert.Equal(t, true, image.IsPublic)
	assert.Equal(t, false, image.Deprecated)
	assert.Equal(t, "available", string(image.Status))
	assert.Equal(t, "2021-08-14T22:44:02Z", image.Created.Format(time.RFC3339))
	assert.Equal(t, "2021-08-14T22:44:02Z", image.Updated.Format(time.RFC3339))
	assert.Equal(t, "2026-07-01T04:00:00Z", image.EOL.Format(time.RFC3339))
	assert.Equal(t, 2500, image.Size)
	assert.Equal(t, 1234567, image.TotalSize)

	assert.ElementsMatch(t, []string{"cloud-init", "distributed-sites"}, image.Capabilities)

	assert.Len(t, image.Regions, 1)
	assert.Equal(t, "us-iad", image.Regions[0].Region)
	assert.Equal(t, "available", string(image.Regions[0].Status))

	assert.ElementsMatch(t, []string{"repair-image", "fix-1"}, image.Tags)
}

func TestImage_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	desc := "Example image description."
	requestData := linodego.ImageUpdateOptions{
		Label:       "Debian 11",
		Description: &desc,
		Tags:        &[]string{"repair-image", "fix-1"},
	}

	imageID := "123"

	base.MockPut(formatMockAPIPath("images/%s", imageID), fixtureData)

	image, err := base.Client.UpdateImage(context.Background(), imageID, requestData)
	assert.NoError(t, err)

	assert.Equal(t, "linode/debian11", image.ID)
	assert.Equal(t, "Debian 11", image.Label)
	assert.Equal(t, "Example image description.", image.Description)
	assert.Equal(t, "Debian", image.Vendor)
	assert.Equal(t, true, image.IsPublic)
	assert.Equal(t, false, image.Deprecated)
	assert.Equal(t, "available", string(image.Status))
	assert.Equal(t, "2021-08-14T22:44:02Z", image.Created.Format(time.RFC3339))
	assert.Equal(t, "2021-08-14T22:44:02Z", image.Updated.Format(time.RFC3339))
	assert.Equal(t, "2026-07-01T04:00:00Z", image.EOL.Format(time.RFC3339))
	assert.Equal(t, 2500, image.Size)
	assert.Equal(t, 1234567, image.TotalSize)

	assert.ElementsMatch(t, []string{"cloud-init", "distributed-sites"}, image.Capabilities)

	assert.Len(t, image.Regions, 1)
	assert.Equal(t, "us-iad", image.Regions[0].Region)
	assert.Equal(t, "available", string(image.Regions[0].Status))

	assert.ElementsMatch(t, []string{"repair-image", "fix-1"}, image.Tags)
}

func TestImage_Upload(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_upload")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.ImageUploadOptions{
		Region:      "us-iad",
		Label:       "Debian 11",
		Description: "Example image description.",
		CloudInit:   true,
		Tags:        &[]string{"repair-image", "fix-1"},
		Image:       strings.NewReader("mock image data"),
	}

	base.MockPost("images/upload", fixtureData)

	image, err := base.Client.UploadImage(context.Background(), requestData)
	assert.NoError(t, err)

	assert.Equal(t, "linode/debian11", image.ID)
	assert.Equal(t, "Debian 11", image.Label)
	assert.Equal(t, "manual", image.Type)
	assert.Equal(t, "Example image description.", image.Description)
	assert.Equal(t, 2500, image.Size)
	assert.Equal(t, linodego.ImageStatus("available"), image.Status)
	assert.Equal(t, true, image.IsPublic)
	assert.Equal(t, "2026-07-01T04:00:00Z", image.EOL.Format(time.RFC3339))
	assert.ElementsMatch(t, []string{"repair-image", "fix-1"}, image.Tags)
	assert.Equal(t, "Debian", image.Vendor)
	assert.False(t, image.Deprecated)
	expectedCapabilities := []string{"cloud-init", "distributed-sites"}
	assert.ElementsMatch(t, expectedCapabilities, image.Capabilities)
	assert.Len(t, image.Regions, 1)
	assert.Equal(t, "us-iad", image.Regions[0].Region)
	assert.Equal(t, linodego.ImageRegionStatus("available"), image.Regions[0].Status)

	// Ensure total_size is set correctly
	assert.Equal(t, 1234567, image.TotalSize)
}

func TestImage_Delete(t *testing.T) {
	client := createMockClient(t)

	imageID := "123"

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, fmt.Sprintf("images/%s", imageID)),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteImage(context.Background(), imageID); err != nil {
		t.Fatal(err)
	}
}

func TestImage_Replicate(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("image_replicate")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.ImageReplicateOptions{
		Regions: []string{
			"us-mia",
			"us-ord",
		},
	}

	base.MockPost("images/private%2F1234/regions", fixtureData)

	image, err := base.Client.ReplicateImage(context.Background(), "private/1234", requestData)
	assert.NoError(t, err)

	assert.Equal(t, "linode/debian11", image.ID)
	assert.Equal(t, "Debian 11", image.Label)
	assert.Equal(t, "Example image description.", image.Description)
	assert.Equal(t, "Debian", image.Vendor)
	assert.Equal(t, true, image.IsPublic)
	assert.Equal(t, false, image.Deprecated)
	assert.Equal(t, "available", string(image.Status))
	assert.Equal(t, "2021-08-14T22:44:02Z", image.Created.Format(time.RFC3339))
	assert.Equal(t, "2021-08-14T22:44:02Z", image.Updated.Format(time.RFC3339))
	assert.Equal(t, "2026-07-01T04:00:00Z", image.EOL.Format(time.RFC3339))
	assert.Equal(t, 2500, image.Size)
	assert.Equal(t, 1234567, image.TotalSize)

	assert.ElementsMatch(t, []string{"cloud-init", "distributed-sites"}, image.Capabilities)

	assert.Len(t, image.Regions, 1)
	assert.Equal(t, "us-iad", image.Regions[0].Region)
	assert.Equal(t, "available", string(image.Regions[0].Status))

	assert.ElementsMatch(t, []string{"repair-image", "fix-1"}, image.Tags)
}
