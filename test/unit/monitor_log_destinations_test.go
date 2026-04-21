package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

const (
    testLogsDestinationID = 12345
)

func TestCreateLogsDestination(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_destinations_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("monitor/streams/destinations", fixtureData)

	path := "audit-logs"
	opts := linodego.LogsDestinationCreateOptions{
		Label: "my-logs-destination",
		Type:  linodego.LogsDestinationTypeAkamaiObjectStorage,
		Details: linodego.LogsDestinationDetailsCreateOptions{
			AccessKeyID:     "123",
			BucketName:      "primary-bucket",
			Host:            "primary-bucket-1.us-east-12.linodeobjects.com",
			Path:            &path,
		},
	}

	dest, err := base.Client.CreateLogsDestination(context.Background(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, dest)
	assert.Equal(t, testLogsDestinationID, dest.ID)
	assert.Equal(t, "OBJ_logs_destination", dest.Label)
	assert.Equal(t, linodego.LogsDestinationStatusActive, dest.Status)
	assert.Equal(t, linodego.LogsDestinationTypeAkamaiObjectStorage, dest.Type)
	assert.Equal(t, "123", string(dest.Details.AccessKeyID))
	assert.Equal(t, "primary-bucket", dest.Details.BucketName)
	assert.Equal(t, "primary-bucket-1.us-iad-12.linodeobjects.com", dest.Details.Host)
	assert.Equal(t, "audit-logs", dest.Details.Path)
	assert.NotNil(t, dest.Created)
	assert.NotNil(t, dest.Updated)
}

func TestCreateLogsDestination_NoPath(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_destinations_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("monitor/streams/destinations", fixtureData)

	opts := linodego.LogsDestinationCreateOptions{
		Label: "my-logs-destination",
		Type:  linodego.LogsDestinationTypeAkamaiObjectStorage,
		Details: linodego.LogsDestinationDetailsCreateOptions{
			AccessKeyID:     "1ABCD23EFG4HIJKLMNO5",
			AccessKeySecret: "1aB2CD3e4fgHi5JK6lmnop7qR8STU9VxYzabcdefHh",
			BucketName:      "primary-bucket",
			Host:            "primary-bucket-1.us-east-12.linodeobjects.com",
			// Path intentionally omitted
		},
	}

	dest, err := base.Client.CreateLogsDestination(context.Background(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, dest)
	assert.Equal(t, testLogsDestinationID, dest.ID)
}

func TestGetLogsDestination(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_destinations_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/streams/destinations/12345", fixtureData)

	dest, err := base.Client.GetLogsDestination(context.Background(), testLogsDestinationID)
	assert.NoError(t, err)
	assert.NotNil(t, dest)
	assert.Equal(t, testLogsDestinationID, dest.ID)
	assert.Equal(t, "OBJ_logs_destination", dest.Label)
	assert.Equal(t, linodego.LogsDestinationStatusActive, dest.Status)
	assert.Equal(t, linodego.LogsDestinationTypeAkamaiObjectStorage, dest.Type)
	assert.Equal(t, "John Q. Linode", dest.CreatedBy)
	assert.Equal(t, "Jane Q. Linode", dest.UpdatedBy)
	assert.Equal(t, 1, dest.Version)
	assert.Equal(t, "123", string(dest.Details.AccessKeyID))
	assert.Equal(t, "primary-bucket", dest.Details.BucketName)
	assert.Equal(t, "primary-bucket-1.us-iad-12.linodeobjects.com", dest.Details.Host)
	assert.Equal(t, "audit-logs", dest.Details.Path)
	assert.NotNil(t, dest.Created)
	assert.NotNil(t, dest.Updated)
}

func TestListLogsDestinations(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_destinations_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/streams/destinations", fixtureData)

	dests, err := base.Client.ListLogsDestinations(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, dests, 1)
	assert.Equal(t, testLogsDestinationID, dests[0].ID)
	assert.Equal(t, "OBJ_logs_destination", dests[0].Label)
	assert.Equal(t, linodego.LogsDestinationStatusActive, dests[0].Status)
	assert.Equal(t, linodego.LogsDestinationTypeAkamaiObjectStorage, dests[0].Type)
	assert.Equal(t, "123", string(dests[0].Details.AccessKeyID))
}

func TestUpdateLogsDestination(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_destinations_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("monitor/streams/destinations/12345", fixtureData)

	opts := linodego.LogsDestinationUpdateOptions{
		Label: "my-logs-destination-renamed",
	}

	dest, err := base.Client.UpdateLogsDestination(context.Background(), testLogsDestinationID, opts)
	assert.NoError(t, err)
	assert.NotNil(t, dest)
}

func TestDeleteLogsDestination(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("monitor/streams/destinations/12345", nil)

	err := base.Client.DeleteLogsDestination(context.Background(), testLogsDestinationID)
	assert.NoError(t, err)
}

func TestListLogsDestinationHistory(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_destinations_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/streams/destinations/12345/history", fixtureData)

	history, err := base.Client.ListLogsDestinationHistory(context.Background(), testLogsDestinationID, nil)
	assert.NoError(t, err)
	assert.Len(t, history, 1)
	assert.Equal(t, testLogsDestinationID, history[0].ID)
	assert.Equal(t, "OBJ_logs_destination", history[0].Label)
	assert.Equal(t, linodego.LogsDestinationStatusActive, history[0].Status)
	assert.Equal(t, linodego.LogsDestinationTypeAkamaiObjectStorage, history[0].Type)
	assert.Equal(t, 1, history[0].Version)
	assert.Equal(t, "123", string(history[0].Details.AccessKeyID))
	assert.Equal(t, "primary-bucket", history[0].Details.BucketName)
	assert.NotNil(t, history[0].Created)
	assert.NotNil(t, history[0].Updated)
}
