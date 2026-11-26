package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testLogsDestinationID = 12345
	testLogStreamID       = 456
	testLKELogStreamID    = 789
)

// ---- Logs Destination tests ----

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

// ---- Log Stream tests ----

func TestCreateLogStream(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_stream")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("monitor/streams", fixtureData)

	status := linodego.StreamStatusActive
	createOpts := linodego.StreamCreateOptions{
		Destinations: []int{testLogsDestinationID},
		Label:        "AuditLog-config",
		Type:         linodego.StreamTypeAuditLogs,
		Status:       &status,
	}

	stream, err := base.Client.CreateLogStream(context.Background(), createOpts)
	assert.NoError(t, err)
	assert.NotNil(t, stream)
	assert.Equal(t, testLogStreamID, stream.ID)
	assert.Equal(t, "AuditLog-config", stream.Label)
	assert.Equal(t, linodego.StreamStatusActive, stream.Status)
	assert.Equal(t, linodego.StreamTypeAuditLogs, stream.Type)
	assert.NotNil(t, stream.Created)
	assert.NotNil(t, stream.Updated)
	assert.Len(t, stream.Destinations, 1)
	assert.Equal(t, testLogsDestinationID, stream.Destinations[0].ID)
	assert.Equal(t, linodego.StreamDestinationTypeAkamaiObjectStorage, stream.Destinations[0].Type)
}

func TestListLogStreams(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_streams_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/streams", fixtureData)

	streams, err := base.Client.ListLogStreams(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.Len(t, streams, 1)
	assert.Equal(t, testLogStreamID, streams[0].ID)
	assert.Equal(t, "AuditLog-config", streams[0].Label)
	assert.Equal(t, linodego.StreamStatusActive, streams[0].Status)
	assert.Equal(t, 1, streams[0].Version)
}

func TestGetLogStream(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_stream")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/streams/456", fixtureData)

	stream, err := base.Client.GetLogStream(context.Background(), testLogStreamID)
	assert.NoError(t, err)
	assert.NotNil(t, stream)
	assert.Equal(t, testLogStreamID, stream.ID)
	assert.Equal(t, "AuditLog-config", stream.Label)
	assert.Equal(t, linodego.StreamStatusActive, stream.Status)
	assert.NotNil(t, stream.Created)
	assert.NotNil(t, stream.Updated)
}

func TestUpdateLogStream(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_stream")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("monitor/streams/456", fixtureData)

	label := "AuditLog-config"
	status := linodego.StreamStatusActive
	updateOpts := linodego.StreamUpdateOptions{
		Destinations: []int{testLogsDestinationID},
		Label:        &label,
		Status:       &status,
	}

	stream, err := base.Client.UpdateLogStream(context.Background(), testLogStreamID, updateOpts)
	assert.NoError(t, err)
	assert.NotNil(t, stream)
	assert.Equal(t, testLogStreamID, stream.ID)
	assert.Equal(t, "AuditLog-config", stream.Label)
	assert.Equal(t, linodego.StreamStatusActive, stream.Status)
	assert.Equal(t, 1, stream.Version)
}

func TestListLogStreamHistory(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_streams_history")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/streams/456/history", fixtureData)

	streams, err := base.Client.ListLogStreamHistory(context.Background(), testLogStreamID, nil)
	assert.NoError(t, err)
	assert.Len(t, streams, 2)
	assert.Equal(t, "AuditLog-config", streams[0].Label)
	assert.Equal(t, 1, streams[0].Version)
	assert.Equal(t, "AuditLog-config-updated", streams[1].Label)
	assert.Equal(t, 2, streams[1].Version)
	assert.Equal(t, linodego.StreamStatusInactive, streams[1].Status)
}

func TestGetLogStream_DestinationDetails(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_stream")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/streams/456", fixtureData)

	stream, err := base.Client.GetLogStream(context.Background(), testLogStreamID)
	assert.NoError(t, err)
	assert.Len(t, stream.Destinations, 1)

	dest := stream.Destinations[0]
	assert.Equal(t, testLogsDestinationID, dest.ID)
	assert.Equal(t, "OBJ_logs_destination", dest.Label)
	assert.Equal(t, linodego.StreamDestinationTypeAkamaiObjectStorage, dest.Type)
	assert.Equal(t, "123", dest.Details.AccessKeyID)
	assert.Equal(t, "primary-bucket", dest.Details.BucketName)
	assert.Equal(t, "primary-bucket-1.us-iad-12.linodeobjects.com", dest.Details.Host)
	assert.Equal(t, "audit-logs", dest.Details.Path)
}

func TestUpdateLogStream_DestinationsOnly(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_stream")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("monitor/streams/456", fixtureData)

	stream, err := base.Client.UpdateLogStream(context.Background(), testLogStreamID, linodego.StreamUpdateOptions{
		Destinations: []int{testLogsDestinationID},
	})
	assert.NoError(t, err)
	assert.NotNil(t, stream)
	assert.Len(t, stream.Destinations, 1)
	assert.Equal(t, testLogsDestinationID, stream.Destinations[0].ID)
}

func TestUpdateLogStream_LabelAndStatus(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_stream")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("monitor/streams/456", fixtureData)

	label := "AuditLog-config"
	status := linodego.StreamStatusActive
	stream, err := base.Client.UpdateLogStream(context.Background(), testLogStreamID, linodego.StreamUpdateOptions{
		Label:  &label,
		Status: &status,
	})
	assert.NoError(t, err)
	assert.Equal(t, "AuditLog-config", stream.Label)
	assert.Equal(t, linodego.StreamStatusActive, stream.Status)
}

func TestDeleteLogStream(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("monitor/streams/456", nil)

	err := base.Client.DeleteLogStream(context.Background(), testLogStreamID)
	assert.NoError(t, err)
}

func TestStreamTypeLKEAuditLogs_Constant(t *testing.T) {
	assert.Equal(t, linodego.StreamType("audit_logs"), linodego.StreamTypeAuditLogs)
	assert.Equal(t, linodego.StreamType("lke_audit_logs"), linodego.StreamTypeLKEAuditLogs)
}

func TestGetLogStream_LKEAuditLogs(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_stream_lke")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/streams/789", fixtureData)

	stream, err := base.Client.GetLogStream(context.Background(), testLKELogStreamID)
	assert.NoError(t, err)
	assert.NotNil(t, stream)
	assert.Equal(t, testLKELogStreamID, stream.ID)
	assert.Equal(t, linodego.StreamTypeLKEAuditLogs, stream.Type)
	assert.NotNil(t, stream.Details)
	assert.Equal(t, []int{1234, 5678}, stream.Details.ClusterIDs)
	assert.False(t, stream.Details.IsAutoAddAllClustersEnabled)
}

func TestGetLogStream_AuditLogs_NoDetails(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_stream")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/streams/456", fixtureData)

	stream, err := base.Client.GetLogStream(context.Background(), testLogStreamID)
	assert.NoError(t, err)
	assert.Nil(t, stream.Details)
}

func TestCreateLogStream_LKEAuditLogs(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_stream_lke")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("monitor/streams", fixtureData)

	createOpts := linodego.StreamCreateOptions{
		Destinations: []int{testLogsDestinationID},
		Label:        "LKEAuditLog-config",
		Type:         linodego.StreamTypeLKEAuditLogs,
		Details: &linodego.StreamDetails{
			ClusterIDs:                  []int{1111, 2222},
			IsAutoAddAllClustersEnabled: false,
		},
	}

	stream, err := base.Client.CreateLogStream(context.Background(), createOpts)
	assert.NoError(t, err)
	assert.NotNil(t, stream)
	assert.Equal(t, testLKELogStreamID, stream.ID)
	assert.Equal(t, linodego.StreamTypeLKEAuditLogs, stream.Type)
	assert.NotNil(t, stream.Details)
	assert.Equal(t, []int{1234, 5678}, stream.Details.ClusterIDs)
}

func TestCreateLogStream_AuditLogs_OmitsDetails(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_stream")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("monitor/streams", fixtureData)

	createOpts := linodego.StreamCreateOptions{
		Destinations: []int{testLogsDestinationID},
		Label:        "AuditLog-config",
		Type:         linodego.StreamTypeAuditLogs,
	}

	stream, err := base.Client.CreateLogStream(context.Background(), createOpts)
	assert.NoError(t, err)
	assert.NotNil(t, stream)
	assert.Nil(t, stream.Details)
}

func TestListLogStreamQuotas(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("monitor_log_streams_quotas_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/streams/quotas", fixtureData)

	quotas, err := base.Client.ListLogStreamQuotas(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.Len(t, quotas, 2)

	assert.Equal(t, "aclp_audit_logs_streams", quotas[0].QuotaID)
	assert.Equal(t, "Number of Audit Logs Streams", quotas[0].QuotaName)
	assert.Equal(t, "Current number of audit logs streams per account", quotas[0].Description)
	assert.Equal(t, 5, quotas[0].QuotaLimit)
	assert.Equal(t, "logs_streams_aclp_audit_logs_streams", quotas[0].QuotaType)

	assert.Equal(t, "aclp_lke_audit_logs_streams", quotas[1].QuotaID)
	assert.Equal(t, "Number of LKE Audit Logs Streams", quotas[1].QuotaName)
	assert.Equal(t, "Current number of LKE audit logs streams per account", quotas[1].Description)
	assert.Equal(t, 10, quotas[1].QuotaLimit)
	assert.Equal(t, "logs_streams_aclp_lke_audit_logs_streams", quotas[1].QuotaType)
}

func TestStreamDetailsMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		details linodego.StreamDetails
		want    string
		wantErr bool
	}{
		{
			name:    "audit subtypes",
			details: linodego.StreamDetails{SubTypes: []string{linodego.AuditLogSubTypeManagementEvents, linodego.AuditLogSubTypeLoginEvents, linodego.AuditLogSubTypeDataEvents}},
			want:    `{"sub_types":["management_events","login_events","data_events"]}`,
		},
		{
			name:    "audit empty subtypes remain explicit",
			details: linodego.StreamDetails{SubTypes: []string{}},
			want:    `{"sub_types":[]}`,
		},
		{
			name:    "audit subtypes with cluster IDs",
			details: linodego.StreamDetails{SubTypes: []string{"data_events"}, ClusterIDs: []int{50}},
			wantErr: true,
		},
		{
			name:    "audit subtypes with auto add",
			details: linodego.StreamDetails{SubTypes: []string{"data_events"}, IsAutoAddAllClustersEnabled: true},
			wantErr: true,
		},
		{
			name:    "empty details retain LKE explicit false",
			details: linodego.StreamDetails{},
			want:    `{"is_auto_add_all_clusters_enabled":false}`,
		},
		{
			name:    "LKE explicit false",
			details: linodego.StreamDetails{ClusterIDs: []int{50, 60}},
			want:    `{"cluster_ids":[50,60],"is_auto_add_all_clusters_enabled":false}`,
		},
		{
			name:    "LKE auto add",
			details: linodego.StreamDetails{IsAutoAddAllClustersEnabled: true},
			want:    `{"is_auto_add_all_clusters_enabled":true}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.details)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(data))
		})
	}
}

func TestLogStreamsAuditSubTypes(t *testing.T) {
	tests := []struct {
		name            string
		subTypes        []string
		createdSubTypes []string
		updatedSubTypes []string
	}{
		{name: "management_events", subTypes: []string{"management_events"}, createdSubTypes: []string{"management_events"}, updatedSubTypes: []string{"management_events"}},
		{name: "login_events", subTypes: []string{"login_events"}, createdSubTypes: []string{"login_events"}, updatedSubTypes: []string{"login_events"}},
		{name: "data events", subTypes: []string{"data_events"}, createdSubTypes: []string{"data_events"}, updatedSubTypes: []string{"data_events"}},
		{name: "sorted and deduplicated subtypes", subTypes: []string{"login_events", "data_events", "management_events", "login_events"}, createdSubTypes: []string{"data_events", "login_events", "management_events"}, updatedSubTypes: []string{"data_events", "login_events", "management_events"}},
		{name: "empty subtypes", subTypes: []string{}, createdSubTypes: []string{"management_events", "login_events"}, updatedSubTypes: []string{"data_events"}},
		{name: "omitted details", createdSubTypes: []string{"management_events", "login_events"}, updatedSubTypes: []string{"data_events"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var base ClientBaseCase
			base.SetUp(t)
			defer base.TearDown(t)

			fixture, err := fixtures.GetFixture("monitor_log_stream_audit_sub_types")
			require.NoError(t, err)
			fixtureJSON, err := json.Marshal(fixture)
			require.NoError(t, err)
			var response map[string]interface{}
			require.NoError(t, json.Unmarshal(fixtureJSON, &response))

			checkStream := func(stream *linodego.Stream, expectedSubTypes []string) {
				t.Helper()
				require.NotNil(t, stream)
				assert.Equal(t, linodego.StreamTypeAuditLogs, stream.Type)
				require.NotNil(t, stream.Details)
				assert.Equal(t, expectedSubTypes, stream.Details.SubTypes)
				assert.Nil(t, stream.Details.ClusterIDs)
			}
			for _, method := range []string{http.MethodPost, http.MethodPut} {
				path := "monitor/streams"
				if method == http.MethodPut {
					path += "/456"
				}
				httpmock.RegisterResponder(method, base.BaseURL+path, func(req *http.Request) (*http.Response, error) {
					var body map[string]json.RawMessage
					require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
					if tt.subTypes == nil {
						assert.NotContains(t, body, "details")
					} else {
						expected, err := json.Marshal(map[string]interface{}{"sub_types": tt.subTypes})
						require.NoError(t, err)
						assert.JSONEq(t, string(expected), string(body["details"]))
					}
					if req.Method == http.MethodPost {
						assert.JSONEq(t, `"audit_logs"`, string(body["type"]))
						response["details"] = map[string]interface{}{"sub_types": tt.createdSubTypes}
					} else {
						assert.NotContains(t, body, "type")
						response["details"] = map[string]interface{}{"sub_types": tt.updatedSubTypes}
					}
					return httpmock.NewJsonResponse(http.StatusOK, response)
				})
			}

			var details *linodego.StreamDetails
			if tt.subTypes != nil {
				details = &linodego.StreamDetails{SubTypes: tt.subTypes}
			}
			created, err := base.Client.CreateLogStream(context.Background(), linodego.StreamCreateOptions{
				Label: "Audit logs", Type: linodego.StreamTypeAuditLogs, Destinations: []int{1}, Details: details,
			})
			require.NoError(t, err)
			checkStream(created, tt.createdSubTypes)

			updated, err := base.Client.UpdateLogStream(context.Background(), testLogStreamID, linodego.StreamUpdateOptions{Details: details})
			require.NoError(t, err)
			checkStream(updated, tt.updatedSubTypes)

			base.MockGet("monitor/streams/456", response)
			stream, err := base.Client.GetLogStream(context.Background(), testLogStreamID)
			require.NoError(t, err)
			checkStream(stream, tt.updatedSubTypes)

			page := map[string]interface{}{
				"data": []interface{}{response}, "page": 1, "pages": 1, "results": 1,
			}
			base.MockGet("monitor/streams", page)
			streams, err := base.Client.ListLogStreams(context.Background(), nil)
			require.NoError(t, err)
			require.Len(t, streams, 1)
			checkStream(&streams[0], tt.updatedSubTypes)

			base.MockGet("monitor/streams/456/history", page)
			history, err := base.Client.ListLogStreamHistory(context.Background(), testLogStreamID, nil)
			require.NoError(t, err)
			require.Len(t, history, 1)
			checkStream(&history[0], tt.updatedSubTypes)
		})
	}
}

func TestCreateLogStreamAuditSubTypesAPIResponse(t *testing.T) {
	fixture, err := fixtures.GetFixture("monitor_log_stream_audit_sub_types")
	require.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	httpmock.RegisterResponder(http.MethodPost, base.BaseURL+"monitor/streams", func(req *http.Request) (*http.Response, error) {
		var body map[string]json.RawMessage
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		assert.JSONEq(t, `{"sub_types":["login_events","data_events","login_events"]}`, string(body["details"]))
		return httpmock.NewJsonResponse(http.StatusCreated, fixture)
	})

	stream, err := base.Client.CreateLogStream(context.Background(), linodego.StreamCreateOptions{
		Label:        "Audit logs with subtypes",
		Type:         linodego.StreamTypeAuditLogs,
		Destinations: []int{1},
		Details: &linodego.StreamDetails{SubTypes: []string{
			linodego.AuditLogSubTypeLoginEvents,
			linodego.AuditLogSubTypeDataEvents,
			linodego.AuditLogSubTypeLoginEvents,
		}},
	})
	require.NoError(t, err)
	require.NotNil(t, stream.Details)
	assert.Equal(t, []string{"data_events", "login_events"}, stream.Details.SubTypes)
	assert.Equal(t, linodego.StreamStatusProvisioning, stream.Status)
	require.Len(t, stream.Destinations, 1)
	assert.Equal(t, 1, stream.Destinations[0].ID)
}

func TestCreateLogStreamAuditSubTypeDisabled(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	httpmock.RegisterResponder(http.MethodPost, base.BaseURL+"monitor/streams",
		httpmock.NewStringResponder(http.StatusBadRequest,
			`{"errors":[{"field":"details.sub_types","reason":"Must be of type management_events, login_events"}]}`,
		).HeaderSet(http.Header{"Content-Type": []string{"application/json"}}))

	_, err := base.Client.CreateLogStream(context.Background(), linodego.StreamCreateOptions{
		Label:        "Audit logs",
		Type:         linodego.StreamTypeAuditLogs,
		Destinations: []int{1},
		Details:      &linodego.StreamDetails{SubTypes: []string{linodego.AuditLogSubTypeDataEvents}},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "details.sub_types")
	assert.Contains(t, err.Error(), "Must be of type management_events, login_events")
}
