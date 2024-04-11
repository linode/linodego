package integration

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"

	"github.com/linode/linodego"
)

type placementGroupModifier func(*linodego.Client, *linodego.PlacementGroupCreateOptions)

func TestPlacementGroup_basic(t *testing.T) {
	client, clientTeardown := createTestClient(t, "fixtures/TestPlacementGroup_basic")

	pg, pgTeardown, err := createPlacementGroup(t, client)
	require.NoError(t, err)

	defer func() {
		pgTeardown()
		clientTeardown()
	}()

	require.NotEqual(t, pg.ID, 0)
	require.Contains(t, pg.Label, "linodego-test-")
	require.NotEmpty(t, pg.Label)
	require.Equal(t, pg.AffinityType, linodego.AffinityTypeAntiAffinityLocal)
	require.Equal(t, pg.IsStrict, false)

	updatedLabel := pg.Label + "-updated"

	pg, err = client.UpdatePlacementGroup(
		context.Background(),
		pg.ID,
		linodego.PlacementGroupUpdateOptions{
			Label: updatedLabel,
		},
	)
	require.NoError(t, err)
	require.Equal(t, pg.Label, updatedLabel)

	refreshedPG, err := client.GetPlacementGroup(context.Background(), pg.ID)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(refreshedPG, pg))

	listedPGs, err := client.ListPlacementGroups(context.Background(), &linodego.ListOptions{
		Filter: fmt.Sprintf("{\"id\": %d}", pg.ID),
	})
	require.NoError(t, err)
	require.NotEmpty(t, listedPGs)
	require.True(t, reflect.DeepEqual(listedPGs[0], *pg))
}

func createPlacementGroup(
	t *testing.T,
	client *linodego.Client,
	pgModifier ...placementGroupModifier,
) (*linodego.PlacementGroup, func(), error) {
	t.Helper()
	createOpts := linodego.PlacementGroupCreateOptions{
		Label:        "linodego-test-" + getUniqueText(),
		Region:       getRegionsWithCaps(t, client, []string{"Placement Group"})[0],
		AffinityType: linodego.AffinityTypeAntiAffinityLocal,
		IsStrict:     false,
	}

	for _, mod := range pgModifier {
		mod(client, &createOpts)
	}

	pg, err := client.CreatePlacementGroup(context.Background(), createOpts)
	if err != nil {
		t.Fatalf("failed to create placement group: %s", err)
	}

	teardown := func() {
		if err := client.DeletePlacementGroup(context.Background(), pg.ID); err != nil {
			t.Errorf("failed to delete placement group: %s", err)
		}
	}
	return pg, teardown, err
}
