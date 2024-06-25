package unit

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestImage_Replicate(t *testing.T) {
	client := createMockClient(t)

	requestData := linodego.ImageReplicateOptions{
		Regions: []string{
			"us-mia",
			"us-ord",
		},
	}

	responseData := linodego.Image{
		ID:    "private/1234",
		Label: "test",
		Regions: []linodego.ImageRegion{
			{
				Region: "us-iad",
				Status: linodego.ImageStatusAvailable,
			},
			{
				Region: "us-mia",
				Status: linodego.ImageStatusReplicating,
			},
			{
				Region: "us-ord",
				Status: linodego.ImageStatusPendingReplication,
			},
		},
	}

	httpmock.RegisterRegexpResponder(
		"POST",
		mockRequestURL(t, "images/private%2F1234/regions"),
		mockRequestBodyValidate(t, requestData, responseData),
	)

	image, err := client.ReplicateImage(context.Background(), "private/1234", requestData)
	require.NoError(t, err)

	require.Equal(t, "private/1234", image.ID)
	require.Equal(t, "test", image.Label)

	require.EqualValues(t, "us-iad", image.Regions[0].Region)
	require.EqualValues(t, linodego.ImageStatusAvailable, image.Regions[0].Status)

	require.EqualValues(t, "us-mia", image.Regions[1].Region)
	require.EqualValues(t, linodego.ImageStatusReplicating, image.Regions[1].Status)

	require.EqualValues(t, "us-ord", image.Regions[2].Region)
	require.EqualValues(t, linodego.ImageStatusPendingReplication, image.Regions[2].Status)
}
