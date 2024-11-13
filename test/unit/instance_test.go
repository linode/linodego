package unit

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"

	"github.com/stretchr/testify/assert"
)

func TestListInstances(t *testing.T) {
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
