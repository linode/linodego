package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestInstanceVolumes_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_volume_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/volumes", fixtureData)

	volumes, err := base.Client.ListInstanceVolumes(context.Background(), 123, nil)
	assert.NoError(t, err)
	assert.Len(t, volumes, 2)

	// Validate first volume
	assert.Equal(t, 1001, volumes[0].ID)
	assert.Equal(t, "volume-1", volumes[0].Label)
	assert.Equal(t, 50, volumes[0].Size)
	assert.Equal(t, linodego.VolumeStatus("available"), volumes[0].Status)

	// Validate second volume
	assert.Equal(t, 1002, volumes[1].ID)
	assert.Equal(t, "volume-2", volumes[1].Label)
	assert.Equal(t, 100, volumes[1].Size)
	assert.Equal(t, linodego.VolumeStatus("resizing"), volumes[1].Status)
}
