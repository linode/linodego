package integration

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/linode/linodego"
)

type placementGroupModifier func(*linodego.Client, *linodego.PlacementGroupCreateOptions)

func TestPlacementGroup_basic(t *testing.T) {
	client, clientTeardown := createTestClient(t, "fixtures/TestPlacementGroup_basic")

	pg, pgTeardown, err := createPlacementGroup(t, client)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		pgTeardown()
		clientTeardown()
	}()

	require.NotEqual(t, pg.ID, 0)
	require.Contains(t, pg.Label, "linodego-test-")
	require.NotEmpty(t, pg.Label)
	require.Equal(t, pg.AffinityType, linodego.AffinityTypeAntiAffinityLocal)
	require.Equal(t, pg.IsStrict, false)

	pg, err = client.UpdatePlacementGroup(
		context.Background(),
		pg.ID,
		linodego.PlacementGroupUpdateOptions{
			Label: pg.Label + "-updated",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, t, pg.Label, pg.Label+"-updated")
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
