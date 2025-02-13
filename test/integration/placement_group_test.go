package integration

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/linode/linodego"
)

type placementGroupModifier func(*linodego.Client, *linodego.PlacementGroupCreateOptions)

func TestPlacementGroup_basic_smoke(t *testing.T) {
	client, clientTeardown := createTestClient(t, "fixtures/TestPlacementGroup_basic")

	// Create a PG
	pg, pgTeardown, err := createPlacementGroup(t, client)
	require.NoError(t, err)

	defer func() {
		pgTeardown()
		clientTeardown()
	}()

	require.NotEqual(t, pg.ID, 0)
	require.Contains(t, pg.Label, "linodego-test-")
	require.NotEmpty(t, pg.Label)
	require.Equal(t, pg.PlacementGroupType, linodego.PlacementGroupTypeAntiAffinityLocal)
	require.Equal(t, pg.PlacementGroupPolicy, linodego.PlacementGroupPolicyFlexible)
	require.Len(t, pg.Members, 0)

	updatedLabel := pg.Label + "-updated"

	// Test that the PG can be updated
	pg, err = client.UpdatePlacementGroup(
		context.Background(),
		pg.ID,
		linodego.PlacementGroupUpdateOptions{
			Label: &updatedLabel,
		},
	)
	require.NoError(t, err)
	require.Equal(t, pg.Label, updatedLabel)

	// Test that the PG can be retrieved from the get endpoint
	refreshedPG, err := client.GetPlacementGroup(context.Background(), pg.ID)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(refreshedPG, pg))

	// Test that the PG can be retrieved from the list endpoint
	listedPGs, err := client.ListPlacementGroups(context.Background(), &linodego.ListOptions{
		Filter: fmt.Sprintf("{\"id\": %d}", pg.ID),
	})
	require.NoError(t, err)
	require.NotEmpty(t, listedPGs)
	require.True(t, reflect.DeepEqual(listedPGs[0], *pg))
}

func TestPlacementGroup_assignment(t *testing.T) {
	client, clientTeardown := createTestClient(t, "fixtures/TestPlacementGroup_assignment")

	pg, pgTeardown, err := createPlacementGroup(t, client)
	require.NoError(t, err)

	// Create an instance to assign to the PG
	inst, err := createInstance(t, client, true, func(client *linodego.Client, options *linodego.InstanceCreateOptions) {
		options.Region = pg.Region
	})
	require.NoError(t, err)

	defer func() {
		// client.DeleteInstance(context.Background(), inst.ID)
		pgTeardown()
		clientTeardown()
	}()

	// Ensure assignment works as expected
	pg, err = client.AssignPlacementGroupLinodes(
		context.Background(),
		pg.ID,
		linodego.PlacementGroupAssignOptions{
			Linodes: []int{
				inst.ID,
			},
		},
	)
	require.NoError(t, err)
	require.Len(t, pg.Members, 1)
	require.Equal(t, pg.Members[0].LinodeID, inst.ID)

	// Refresh the instance to ensure the assignment has completed
	inst, err = client.GetInstance(context.Background(), inst.ID)
	require.NoError(t, err)
	require.NotNil(t, inst.PlacementGroup)
	require.Equal(t, inst.PlacementGroup.ID, pg.ID)
	require.Equal(t, inst.PlacementGroup.Label, pg.Label)
	require.Equal(t, inst.PlacementGroup.PlacementGroupPolicy, pg.PlacementGroupPolicy)
	require.Equal(t, inst.PlacementGroup.PlacementGroupType, pg.PlacementGroupType)

	// Ensure unassignment works as expected
	pg, err = client.UnassignPlacementGroupLinodes(
		context.Background(),
		pg.ID,
		linodego.PlacementGroupUnAssignOptions{
			Linodes: []int{
				inst.ID,
			},
		},
	)
	require.NoError(t, err)
	require.Len(t, pg.Members, 0)
}

func createPlacementGroup(
	t *testing.T,
	client *linodego.Client,
	pgModifier ...placementGroupModifier,
) (*linodego.PlacementGroup, func(), error) {
	t.Helper()
	createOpts := linodego.PlacementGroupCreateOptions{
		Label:                "linodego-test-" + getUniqueText(),
		Region:               getRegionsWithCaps(t, client, []string{"Placement Group"})[0],
		PlacementGroupType:   linodego.PlacementGroupTypeAntiAffinityLocal,
		PlacementGroupPolicy: linodego.PlacementGroupPolicyFlexible,
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
