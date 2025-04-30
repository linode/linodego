package unit

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstances_List(t *testing.T) {
	fixtures := NewTestFixtures()

	fixtureData, err := fixtures.GetFixture("linodes_list")
	if err != nil {
		t.Fatalf("Failed to load fixture: %v", err)
	}

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances", fixtureData)

	instances, err := base.Client.ListInstances(context.Background(), nil)
	if err != nil {
		t.Fatalf("Error listing instances: %v", err)
	}

	assert.Equal(t, 1, len(instances))
	linode := instances[0]
	assert.Equal(t, 123, linode.ID)
	assert.Equal(t, "linode123", linode.Label)
	assert.Equal(t, "running", string(linode.Status))
	assert.Equal(t, "203.0.113.1", linode.IPv4[0].String())
	assert.Equal(t, "g6-standard-1", linode.Type)
	assert.Equal(t, "us-east", linode.Region)
	assert.Equal(t, 4096, linode.Specs.Memory)
	assert.Equal(t, "2018-01-01 00:01:01 +0000 UTC", linode.Backups.LastSuccessful.String())
	require.NotNil(t, linode.PlacementGroup.MigratingTo)
	assert.Equal(t, 2468, *linode.PlacementGroup.MigratingTo)
}

func TestInstance_Get(t *testing.T) {
	fixtures := NewTestFixtures()

	fixtureData, err := fixtures.GetFixture("instance_get")
	if err != nil {
		t.Fatalf("Failed to load fixture: %v", err)
	}

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	instanceID := 123
	base.MockGet(fmt.Sprintf("linode/instances/%d", instanceID), fixtureData)

	instance, err := base.Client.GetInstance(context.Background(), instanceID)
	if err != nil {
		t.Fatalf("Error fetching instance: %v", err)
	}

	assert.Equal(t, 123, instance.ID)
	assert.Equal(t, "linode123", instance.Label)
	assert.Equal(t, "running", string(instance.Status))
	assert.Equal(t, "203.0.113.1", instance.IPv4[0].String())
	assert.Equal(t, "g6-standard-1", instance.Type)
	assert.Equal(t, "us-east", instance.Region)
	assert.Equal(t, 4096, instance.Specs.Memory)
	assert.Equal(t, "2018-01-01 00:01:01 +0000 UTC", instance.Backups.LastSuccessful.String())
	require.NotNil(t, instance.PlacementGroup.MigratingTo)
	assert.Equal(t, 2468, *instance.PlacementGroup.MigratingTo)
}

func TestInstance_Migrate(t *testing.T) {
	client := createMockClient(t)

	upgrade := false

	requestData := linodego.InstanceMigrateOptions{
		Type:    "cold",
		Region:  "us-west",
		Upgrade: &upgrade,
	}

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "/linode/instances/123456/migrate"),
		mockRequestBodyValidate(t, requestData, nil))

	if err := client.MigrateInstance(context.Background(), 123456, requestData); err != nil {
		t.Fatal(err)
	}
}

func TestInstance_ResetPassword(t *testing.T) {
	client := createMockClient(t)

	requestData := linodego.InstancePasswordResetOptions{
		RootPass: "@v3ry53cu3eP@s5w0rd",
	}

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "linode/instances/123456/password"),
		mockRequestBodyValidate(t, requestData, nil))

	if err := client.ResetInstancePassword(context.Background(), 123456, requestData); err != nil {
		t.Fatal(err)
	}
}

func TestInstance_Get_MonthlyTransfer(t *testing.T) {
	if strconv.IntSize < 64 {
		t.Skip("V1 monthly transfer doesn't work on 32 or lower bits system")
	}
	fixtureData, err := fixtures.GetFixture("instance_monthly_transfer_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/12345/transfer/2024/11", fixtureData)

	stats, err := base.Client.GetInstanceTransferMonthly(context.Background(), 12345, 2024, 11)
	assert.NoError(t, err)

	assert.Equal(t, 30471077120, stats.BytesIn)
	assert.Equal(t, 22956600198, stats.BytesOut)
	assert.Equal(t, 53427677318, stats.BytesTotal)
}

func TestInstance_Get_MonthlyTransferV2(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_monthly_transfer_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/12345/transfer/2024/11", fixtureData)

	stats, err := base.Client.GetInstanceTransferMonthlyV2(context.Background(), 12345, 2024, 11)
	assert.NoError(t, err)

	assert.Equal(t, uint64(30471077120), stats.BytesIn)
	assert.Equal(t, uint64(22956600198), stats.BytesOut)
	assert.Equal(t, uint64(53427677318), stats.BytesTotal)
}

func TestInstance_Upgrade(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("linode/instances/12345/mutate", nil)

	err := base.Client.UpgradeInstance(context.Background(), 12345, linodego.InstanceUpgradeOptions{
		AllowAutoDiskResize: true,
	})
	assert.NoError(t, err)
}

func TestInstance_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	createOptions := linodego.InstanceCreateOptions{
		Region:   "us-east",
		Type:     "g6-standard-1",
		Label:    "new-instance",
		Image:    "linode/ubuntu22.04",
		RootPass: "securepassword",
	}

	base.MockPost("linode/instances", fixtureData)

	instance, err := base.Client.CreateInstance(context.Background(), createOptions)
	assert.NoError(t, err)
	assert.Equal(t, "new-instance", instance.Label)
}

func TestInstance_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	updateOptions := linodego.InstanceUpdateOptions{
		Label: "updated-instance",
	}

	base.MockPut("linode/instances/123", fixtureData)

	instance, err := base.Client.UpdateInstance(context.Background(), 123, updateOptions)
	assert.NoError(t, err)
	assert.Equal(t, "updated-instance", instance.Label)
}

func TestInstance_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("linode/instances/123", nil)

	err := base.Client.DeleteInstance(context.Background(), 123)
	assert.NoError(t, err)
}

func TestInstance_Boot(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("linode/instances/123/boot", nil)

	err := base.Client.BootInstance(context.Background(), 123, 0)
	assert.NoError(t, err)
}

func TestInstance_Reboot(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("linode/instances/123/reboot", nil)

	err := base.Client.RebootInstance(context.Background(), 123, 0)
	assert.NoError(t, err)
}

func TestInstance_Clone(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_clone")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	cloneOptions := linodego.InstanceCloneOptions{
		Region: "us-east",
		Type:   "g6-standard-1",
		Label:  "cloned-instance",
	}

	base.MockPost("linode/instances/123/clone", fixtureData)

	instance, err := base.Client.CloneInstance(context.Background(), 123, cloneOptions)
	assert.NoError(t, err)
	assert.Equal(t, "cloned-instance", instance.Label)
}

func TestInstance_Resize(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	resizeOptions := linodego.InstanceResizeOptions{
		Type: "g6-standard-2",
	}

	base.MockPost("linode/instances/123/resize", "{}")

	err := base.Client.ResizeInstance(context.Background(), 123, resizeOptions)
	assert.NoError(t, err)
}

func TestInstance_Rescue(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	rescueOptions := linodego.InstanceRescueOptions{}

	base.MockPost("linode/instances/123/rescue", nil)

	err := base.Client.RescueInstance(context.Background(), 123, rescueOptions)
	assert.NoError(t, err)
}

func TestInstance_Rebuild(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_rebuild")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	rebuildOptions := linodego.InstanceRebuildOptions{
		Image: "linode/ubuntu22.04",
	}

	base.MockPost("linode/instances/123/rebuild", fixtureData)

	instance, err := base.Client.RebuildInstance(context.Background(), 123, rebuildOptions)
	assert.NoError(t, err)
	assert.Equal(t, "linode/ubuntu22.04", instance.Image)
}
