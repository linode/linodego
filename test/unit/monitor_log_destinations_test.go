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
			AccessKeyID:     "1ABCD23EFG4HIJKLMNO5",
			AccessKeySecret: "1aB2CD3e4fgHi5JK6lmnop7qR8STU9VxYzabcdefHh",
			BucketName:      "primary-bucket",
			Host:            "primary-bucket-1.us-iad-12.linodeobjects.com",
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
	assert.Equal(t, "1ABCD23EFG4HIJKLMNO5", string(dest.Details.AccessKeyID))
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
			Host:            "primary-bucket-1.us-iad-12.linodeobjects.com",
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
	assert.Equal(t, "user", dest.CreatedBy)
	assert.Equal(t, "user", dest.UpdatedBy)
	assert.Equal(t, 1, dest.Version)
	assert.Equal(t, "1ABCD23EFG4HIJKLMNO5", string(dest.Details.AccessKeyID))
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
	assert.Len(t, dests, 2)

	// First destination: akamai_object_storage
	assert.Equal(t, testLogsDestinationID, dests[0].ID)
	assert.Equal(t, "OBJ_logs_destination", dests[0].Label)
	assert.Equal(t, linodego.LogsDestinationStatusActive, dests[0].Status)
	assert.Equal(t, linodego.LogsDestinationTypeAkamaiObjectStorage, dests[0].Type)
	assert.Equal(t, "123", dests[0].Details.AccessKeyID)
	assert.Equal(t, "primary-bucket", dests[0].Details.BucketName)
	assert.Equal(t, "primary-bucket-1.us-iad-12.linodeobjects.com", dests[0].Details.Host)
	assert.Equal(t, "audit-logs", dests[0].Details.Path)

	// Second destination: custom_https
	assert.Equal(t, 67890, dests[1].ID)
	assert.Equal(t, "HTTPS_logs_destination", dests[1].Label)
	assert.Equal(t, linodego.LogsDestinationStatusActive, dests[1].Status)
	assert.Equal(t, linodego.LogsDestinationTypeCustomHTTPS, dests[1].Type)
	assert.Equal(t, "https://my-site.com/log-storage/database-info", dests[1].Details.EndpointURL)
	assert.NotNil(t, dests[1].Details.Authentication)
	assert.Equal(t, linodego.LogsDestinationCustomHTTPSAuthTypeBasic, dests[1].Details.Authentication.Type)
	assert.NotNil(t, dests[1].Details.Authentication.Details)
	assert.Equal(t, "John_Q", dests[1].Details.Authentication.Details.Username)
	assert.Equal(t, "application/json", dests[1].Details.ContentType)
	assert.Len(t, dests[1].Details.CustomHeaders, 1)
	assert.Equal(t, "Cache-Control", dests[1].Details.CustomHeaders[0].Name)
	assert.Equal(t, "max-age=0", dests[1].Details.CustomHeaders[0].Value)
	assert.Equal(t, "gzip", dests[1].Details.DataCompression)
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

func TestCreateLogsDestination_CustomHTTPS(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_destinations_custom_https_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("monitor/streams/destinations", fixtureData)

	opts := linodego.LogsDestinationCreateOptions{
		Label: "HTTPS_logs_destination",
		Type:  linodego.LogsDestinationTypeCustomHTTPS,
		Details: linodego.LogsDestinationCustomHTTPSDetailsCreateOptions{
			EndpointURL: "https://my-site.com/log-storage/database-info",
			Authentication: &linodego.LogsDestinationCustomHTTPSAuthDetails{
				Type: linodego.LogsDestinationCustomHTTPSAuthTypeBasic,
				Details: &linodego.LogsDestinationCustomHTTPSBasicAuthDetails{
					Username: "John_Q",
					Password: "p@$$w0Rd",
				},
			},
			ContentType:     "application/json",
			DataCompression: "gzip",
			CustomHeaders: []linodego.LogsDestinationCustomHTTPSHeader{
				{Name: "Cache-Control", Value: "max-age=0"},
			},
		},
	}

	dest, err := base.Client.CreateLogsDestination(context.Background(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, dest)
	assert.Equal(t, 67890, dest.ID)
	assert.Equal(t, "HTTPS_logs_destination", dest.Label)
	assert.Equal(t, linodego.LogsDestinationStatusActive, dest.Status)
	assert.Equal(t, linodego.LogsDestinationTypeCustomHTTPS, dest.Type)
	assert.Equal(t, "https://my-site.com/log-storage/database-info", dest.Details.EndpointURL)
	assert.NotNil(t, dest.Details.Authentication)
	assert.Equal(t, linodego.LogsDestinationCustomHTTPSAuthTypeBasic, dest.Details.Authentication.Type)
	assert.Equal(t, "application/json", dest.Details.ContentType)
	assert.Equal(t, "gzip", dest.Details.DataCompression)
	assert.Len(t, dest.Details.CustomHeaders, 1)
	assert.Equal(t, "Cache-Control", dest.Details.CustomHeaders[0].Name)
	assert.NotNil(t, dest.Created)
	assert.NotNil(t, dest.Updated)
}

func TestUpdateLogsDestination_CustomHTTPS(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_destinations_custom_https_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("monitor/streams/destinations/67890", fixtureData)

	newURL := "https://my-site.com/log-storage/v2"
	opts := linodego.LogsDestinationUpdateOptions{
		Label: "HTTPS_logs_destination_renamed",
		Details: &linodego.LogsDestinationCustomHTTPSDetailsUpdateOptions{
			EndpointURL: newURL,
			Authentication: &linodego.LogsDestinationCustomHTTPSAuthDetails{
				Type: linodego.LogsDestinationCustomHTTPSAuthTypeBasic,
				Details: &linodego.LogsDestinationCustomHTTPSBasicAuthDetails{
					Username: "John_Q",
					Password: "newpassword",
				},
			},
		},
	}

	dest, err := base.Client.UpdateLogsDestination(context.Background(), 67890, opts)
	assert.NoError(t, err)
	assert.NotNil(t, dest)
	assert.Equal(t, 67890, dest.ID)
	assert.Equal(t, linodego.LogsDestinationTypeCustomHTTPS, dest.Type)
}

func TestListLogsDestinationHistory(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_destinations_history_list")
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
	assert.Equal(t, "123", history[0].Details.AccessKeyID)
	assert.Equal(t, "primary-bucket", history[0].Details.BucketName)
	assert.NotNil(t, history[0].Created)
	assert.NotNil(t, history[0].Updated)
}
