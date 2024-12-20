package unit

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestPlacementGroups_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("placement_groups_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("placement/groups", fixtureData)

	pgs, err := base.Client.ListPlacementGroups(context.Background(), nil)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(pgs))
	pg := pgs[0]

	assert.Equal(t, 528, pg.ID)
	assert.Equal(t, true, pg.IsCompliant)
	assert.Equal(t, "PG_Miami_failover", pg.Label)
	assert.Equal(t, true, pg.Members[0].IsCompliant)
	assert.Equal(t, 123, pg.Members[0].LinodeID)
	assert.Equal(t, 123, pg.Migrations.Inbound[0].LinodeID)
	assert.Equal(t, 456, pg.Migrations.Outbound[0].LinodeID)
	assert.Equal(t, linodego.PlacementGroupPolicy("strict"), pg.PlacementGroupPolicy)
	assert.Equal(t, linodego.PlacementGroupType("anti-affinity:local"), pg.PlacementGroupType)
	assert.Equal(t, "us-mia", pg.Region)
}

func TestPlacementGroups_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("placement_groups_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("placement/groups/528", fixtureData)

	pg, err := base.Client.GetPlacementGroup(context.Background(), 528)
	assert.NoError(t, err)

	assert.Equal(t, 528, pg.ID)
	assert.Equal(t, true, pg.IsCompliant)
	assert.Equal(t, "PG_Miami_failover", pg.Label)
	assert.Equal(t, true, pg.Members[0].IsCompliant)
	assert.Equal(t, 123, pg.Members[0].LinodeID)
	assert.Equal(t, 123, pg.Migrations.Inbound[0].LinodeID)
	assert.Equal(t, 456, pg.Migrations.Outbound[0].LinodeID)
	assert.Equal(t, linodego.PlacementGroupPolicy("strict"), pg.PlacementGroupPolicy)
	assert.Equal(t, linodego.PlacementGroupType("anti-affinity:local"), pg.PlacementGroupType)
	assert.Equal(t, "us-mia", pg.Region)
}

func TestPlacementGroups_Delete(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, "placement/groups/528"), httpmock.NewStringResponder(200, "{}"))

	if err := client.DeletePlacementGroup(context.Background(), 528); err != nil {
		t.Fatal(err)
	}
}

func TestPlacementGroups_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("placement_group_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.PlacementGroupUpdateOptions{
		Label: "PG_Miami_failover_new",
	}

	base.MockPut("placement/groups/528", fixtureData)

	pg, err := base.Client.UpdatePlacementGroup(context.Background(), 528, requestData)
	assert.NoError(t, err)

	assert.Equal(t, 528, pg.ID)
	assert.Equal(t, true, pg.IsCompliant)
	assert.Equal(t, "PG_Miami_failover_new", pg.Label)
	assert.Equal(t, true, pg.Members[0].IsCompliant)
	assert.Equal(t, 123, pg.Members[0].LinodeID)
	assert.Equal(t, linodego.PlacementGroupPolicy("strict"), pg.PlacementGroupPolicy)
	assert.Equal(t, linodego.PlacementGroupType("anti-affinity:local"), pg.PlacementGroupType)
	assert.Equal(t, "us-mia", pg.Region)
}

func TestPlacementGroups_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("placement_group_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.PlacementGroupCreateOptions{
		Label:                "PG_Miami_failover",
		Region:               "us-mia",
		PlacementGroupType:   linodego.PlacementGroupType("anti-affinity:local"),
		PlacementGroupPolicy: linodego.PlacementGroupPolicy("strict"),
	}

	base.MockPost("placement/groups", fixtureData)

	pg, err := base.Client.CreatePlacementGroup(context.Background(), requestData)
	assert.NoError(t, err)

	assert.Equal(t, 528, pg.ID)
	assert.Equal(t, true, pg.IsCompliant)
	assert.Equal(t, "PG_Miami_failover", pg.Label)
	assert.Equal(t, true, pg.Members[0].IsCompliant)
	assert.Equal(t, 123, pg.Members[0].LinodeID)
	assert.Equal(t, linodego.PlacementGroupPolicy("strict"), pg.PlacementGroupPolicy)
	assert.Equal(t, linodego.PlacementGroupType("anti-affinity:local"), pg.PlacementGroupType)
	assert.Equal(t, "us-mia", pg.Region)
}

func TestPlacementGroups_Assign(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("placement_group_assign")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.PlacementGroupAssignOptions{
		Linodes: []int{456},
	}

	base.MockPost("placement/groups/528/assign", fixtureData)

	pg, err := base.Client.AssignPlacementGroupLinodes(context.Background(), 528, requestData)
	assert.NoError(t, err)

	assert.Equal(t, 528, pg.ID)
	assert.Equal(t, true, pg.IsCompliant)
	assert.Equal(t, "PG_Miami_failover", pg.Label)
	assert.Equal(t, true, pg.Members[0].IsCompliant)
	assert.Equal(t, 123, pg.Members[0].LinodeID)
	assert.Equal(t, true, pg.Members[1].IsCompliant)
	assert.Equal(t, 456, pg.Members[1].LinodeID)
	assert.Equal(t, linodego.PlacementGroupPolicy("strict"), pg.PlacementGroupPolicy)
	assert.Equal(t, linodego.PlacementGroupType("anti-affinity:local"), pg.PlacementGroupType)
	assert.Equal(t, "us-mia", pg.Region)
}

func TestPlacementGroups_Unassign(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("placement_group_unassign")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.PlacementGroupUnAssignOptions{
		Linodes: []int{456},
	}

	base.MockPost("placement/groups/528/unassign", fixtureData)

	pg, err := base.Client.UnassignPlacementGroupLinodes(context.Background(), 528, requestData)
	assert.NoError(t, err)

	assert.Equal(t, 528, pg.ID)
	assert.Equal(t, true, pg.IsCompliant)
	assert.Equal(t, "PG_Miami_failover", pg.Label)
	assert.Equal(t, true, pg.Members[0].IsCompliant)
	assert.Equal(t, 123, pg.Members[0].LinodeID)
	assert.Equal(t, linodego.PlacementGroupPolicy("strict"), pg.PlacementGroupPolicy)
	assert.Equal(t, linodego.PlacementGroupType("anti-affinity:local"), pg.PlacementGroupType)
	assert.Equal(t, "us-mia", pg.Region)
}
