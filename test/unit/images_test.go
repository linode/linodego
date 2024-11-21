package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

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
