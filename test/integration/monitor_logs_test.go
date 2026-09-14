package integration

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testCustomHTTPSEndpointURL        = "https://pl-labkrk-2-https-server-dev.cloud-streams-dev.akadns.net/stream-custom-https-server-api/logs/no-auth"
	testCustomHTTPSContentType        = "application/json"
	testCustomHTTPSUpdatedContentType = "application/json; charset=utf-8"
	testCustomHTTPSDataCompression    = "gzip"
	testCustomHTTPSHeaderName         = "X-Logs-Custom-Path"
	testCustomHTTPSHeaderValue        = "linodego-e2e-test"
)

// requireACLPLogsStreamTests skips the test if RUN_ACLP_LOGS_STREAM_TESTS is not set.
// Call this before creating a test client so the env check short-circuits early.
func requireACLPLogsStreamTests(t *testing.T) {
	if testingMode == recorder.ModeReplaying {
		return
	}
	t.Helper()
	val := os.Getenv("RUN_ACLP_LOGS_STREAM_TESTS")
	if val != "yes" && val != "true" {
		t.Skipf("RUN_ACLP_LOGS_STREAM_TESTS must be set to 'yes' or 'true' to run stream tests")
	}
}

// requireTrafficPeakDestinationTests skips TrafficPeak tests if
// RUN_TRAFFIC_PEAK_DESTINATION_TESTS is not set
// Unlike requireACLPLogsStreamTests, replay mode also
// requires the flag because TrafficPeak fixtures do not exist yet.
func requireTrafficPeakDestinationTests(t *testing.T) {
	t.Helper()

	val := os.Getenv("RUN_TRAFFIC_PEAK_DESTINATION_TESTS")
	if val != "yes" && val != "true" {
		t.Skipf("RUN_TRAFFIC_PEAK_DESTINATION_TESTS must be set to 'yes' or 'true' to run TrafficPeak tests")
	}
}

// requireCustomHTTPSDestinationTests skips live Custom HTTPS integration tests
// unless they are explicitly enabled.
func requireCustomHTTPSDestinationTests(t *testing.T) {
	t.Helper()

	if testingMode == recorder.ModeReplaying {
		return
	}

	val := os.Getenv("RUN_CUSTOM_HTTPS_DESTINATION_TESTS")
	if val != "yes" && val != "true" {
		t.Skipf("RUN_CUSTOM_HTTPS_DESTINATION_TESTS must be set to 'yes' or 'true' to run Custom HTTPS tests")
	}
}

// creates a object storage and access keys for use in tests
func setupObjectStorageForLogs(t *testing.T, client *linodego.Client) (*linodego.ObjectStorageBucket, *linodego.ObjectStorageKey, func()) {
	t.Helper()

	regions := getRegionsWithCaps(t, client, []linodego.RegionCapability{linodego.CapabilityObjectStorage})
	if len(regions) == 0 {
		t.Fatal("no region with Object Storage capability found")
	}

	bucket, err := client.CreateObjectStorageBucket(context.Background(), linodego.ObjectStorageBucketCreateOptions{
		Region:      regions[0],
		Label:       testLabel(),
		ACL:         "private",
		CorsEnabled: linodego.Pointer(false),
	})
	if err != nil {
		t.Fatalf("Error creating storage bucket, got error %v", err)
	}

	storageKey, err := client.CreateObjectStorageKey(context.Background(), linodego.ObjectStorageKeyCreateOptions{
		Label:   testLabel(),
		Regions: []string{regions[0]},
	})
	if err != nil {
		_ = client.DeleteObjectStorageBucket(context.Background(), bucket.Region, bucket.Label)
		t.Fatalf("Error creating storage key, got error %v", err)
	}

	teardown := func() {
		if terr := client.DeleteObjectStorageKey(context.Background(), storageKey.ID); terr != nil {
			t.Errorf("Expected to delete a storage key, but got %v", terr)
		}

		bucketObjects, terr := client.ListObjectStorageBucketContents(context.Background(), bucket.Region, bucket.Label, nil)
		if terr == nil {
			for _, obj := range bucketObjects.Data {
				url, err := client.CreateObjectStorageObjectURL(context.TODO(), bucket.Region, bucket.Label, linodego.ObjectStorageObjectURLCreateOptions{
					Name:      obj.Name,
					Method:    http.MethodDelete,
					ExpiresIn: &objectStorageObjectURLExpirySeconds,
				})
				if err != nil {
					t.Errorf("failed to get object DELETE url: %s", err)
					continue
				}

				if testingMode == recorder.ModeReplaying {
					continue
				}

				req, err := http.NewRequest(http.MethodDelete, url.URL, nil)
				if err != nil {
					t.Errorf("failed to build request: %s", err)
					continue
				}

				res, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Errorf("failed to delete object: %s", err)
					continue
				}
				_, _ = io.Copy(io.Discard, res.Body)
				res.Body.Close()
				if res.StatusCode != 204 {
					t.Errorf("expected status code to be 204; got %d", res.StatusCode)
				}
			}
		} else {
			t.Errorf("Expected to list objects in object storage, but got %v", terr)
		}

		if terr := client.DeleteObjectStorageBucket(context.Background(), bucket.Region, bucket.Label); terr != nil {
			t.Errorf("Expected to delete object storage bucket, but got %v", terr)
		}
	}

	return bucket, storageKey, teardown
}

// creates a LogsDestination for use in tests
func setupLogsDestination(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.LogsDestination, *linodego.ObjectStorageKey, func()) {
	t.Helper()

	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	bucket, storageKey, storageTeardown := setupObjectStorageForLogs(t, client)

	dest, err := client.CreateLogsDestination(context.Background(), linodego.LogsDestinationCreateOptions{
		Label: testLabel(),
		Type:  linodego.LogsDestinationTypeAkamaiObjectStorage,
		Details: linodego.LogsDestinationDetailsCreateOptions{
			AccessKeyID:     storageKey.AccessKey,
			AccessKeySecret: storageKey.SecretKey,
			BucketName:      bucket.Label,
			Host:            bucket.Hostname,
		},
	})
	if err != nil {
		storageTeardown()
		fixtureTeardown()
		t.Fatalf("Error creating logs destination, got error %v", err)
	}

	teardown := func() {
		_, terr := client.GetLogsDestination(context.Background(), dest.ID)
		if apiErr, ok := terr.(*linodego.Error); ok && apiErr.Code == 404 {
			// Already gone — nothing to do.
		} else {
			if terr != nil {
				t.Errorf("Error while checking destination existence: %v", terr)
			}
			// Object exists or GET failed for another reason — try to delete anyway.
			if terr := client.DeleteLogsDestination(context.Background(), dest.ID); terr != nil {
				t.Errorf("Expected to delete a logs destination, but got %v", terr)
			}
		}
		storageTeardown()
		fixtureTeardown()
	}

	return client, dest, storageKey, teardown
}

// setupHTTPLogsDestination creates either a TrafficPeak or Custom HTTPS
// destination and returns a teardown that tolerates deletion by the test itself.
func setupHTTPLogsDestination(
	t *testing.T,
	fixturesYaml string,
	destinationType linodego.LogsDestinationType,
	details any,
) (*linodego.Client, *linodego.LogsDestination, func()) {
	t.Helper()

	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	destination, err := client.CreateLogsDestination(context.Background(), linodego.LogsDestinationCreateOptions{
		Label:   testLabel(),
		Type:    destinationType,
		Details: details,
	})
	if err != nil {
		fixtureTeardown()
	}
	require.NoError(t, err)
	require.NotNil(t, destination)

	teardown := func() {
		defer fixtureTeardown()

		_, err := client.GetLogsDestination(context.Background(), destination.ID)
		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code == http.StatusNotFound {
			return
		}
		if err != nil {
			t.Errorf("failed to check destination before cleanup: %v", err)
		}
		if err := client.DeleteLogsDestination(context.Background(), destination.ID); err != nil {
			t.Errorf("failed to delete logs destination: %v", err)
		}
	}

	return client, destination, teardown
}

func setupTrafficPeakLogsDestination(
	t *testing.T,
	fixturesYaml string,
) (*linodego.Client, *linodego.LogsDestination, func()) {
	t.Helper()
	requireTrafficPeakDestinationTests(t)

	return setupHTTPLogsDestination(
		t,
		fixturesYaml,
		linodego.LogsDestinationTypeTrafficPeak,
		linodego.LogsDestinationTrafficPeakDetailsCreateOptions{
			EndpointURL:     "https://example.com/",
			ContentType:     linodego.Pointer("application/json"),
			DataCompression: linodego.Pointer("None"),
			CustomHeaders: []linodego.LogsDestinationTrafficPeakHeader{
				{Name: "x-test-header", Value: "traffic-peak"},
			},
			Authentication: linodego.LogsDestinationTrafficPeakAuthDetails{
				Details: linodego.LogsDestinationTrafficPeakBasicAuthDetails{
					Username: "user",
					Password: "password",
				},
			},
		},
	)
}

func setupCustomHTTPSLogsDestination(
	t *testing.T,
	fixturesYaml string,
) (*linodego.Client, *linodego.LogsDestination, func()) {
	t.Helper()
	requireCustomHTTPSDestinationTests(t)

	return setupHTTPLogsDestination(
		t,
		fixturesYaml,
		linodego.LogsDestinationTypeCustomHTTPS,
		linodego.LogsDestinationCustomHTTPSDetailsCreateOptions{
			EndpointURL: testCustomHTTPSEndpointURL,
			Authentication: &linodego.LogsDestinationCustomHTTPSAuthDetails{
				Type: linodego.LogsDestinationCustomHTTPSAuthTypeNone,
			},
			ContentType:     testCustomHTTPSContentType,
			DataCompression: testCustomHTTPSDataCompression,
			CustomHeaders: []linodego.LogsDestinationCustomHTTPSHeader{
				{Name: testCustomHTTPSHeaderName, Value: testCustomHTTPSHeaderValue},
			},
		},
	)
}

func testLabel() string {
	return "go-test-logs-monitoring-" + getUniqueText()
}

// setupLogStream creates a test client, Object Storage resources, a LogsDestination,
// and a LogStream. The returned teardown cleans up all resources in reverse order.
func setupLogStream(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.Stream, *linodego.LogsDestination, func()) {
	t.Helper()

	client, dest, _, destTeardown := setupLogsDestination(t, fixturesYaml)
	requireNoExistingStreams(t, client)

	stream, err := client.CreateLogStream(context.Background(), linodego.StreamCreateOptions{
		Label:        fmt.Sprintf("go-test-log-stream-%d", time.Now().UnixNano()),
		Type:         linodego.StreamTypeAuditLogs,
		Destinations: []int{dest.ID},
	})
	if err != nil {
		destTeardown()
		t.Fatalf("Error creating log stream, got error %v", err)
	}

	teardown := func() {
		// Wait for stream to reach a stable (non-transitional) state before deleting;
		// deleting while in "deactivating"/"provisioning" silently fails and then
		// the destination delete returns 400 because the stream is still attached.
		if _, err := waitForLogStreamProvisioned(context.Background(), client, stream.ID, 60, 3600); err != nil {
			if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code == 404 {
				// Already deleted by the test itself — skip to resource cleanup.
				destTeardown()
				return
			}
			t.Logf("Warning: could not wait for stream %d to stabilize before deletion: %v", stream.ID, err)
		}
		if err := client.DeleteLogStream(context.Background(), stream.ID); err != nil {
			if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code == 404 {
				// Already gone — nothing to do.
			} else {
				t.Errorf("Expected to delete log stream %d, but got %v", stream.ID, err)
			}
		} else if err := waitForLogStreamDeleted(context.Background(), client, stream.ID, 60, 3600); err != nil {
			t.Logf("Warning: stream %d may not be fully deleted: %v", stream.ID, err)
		}
		destTeardown()
	}

	return client, stream, dest, teardown
}

// waitForLogStreamDeleted polls until the stream returns 404, indicating it has been fully removed.
func waitForLogStreamDeleted(
	ctx context.Context,
	client *linodego.Client,
	streamID int,
	pollIntervalSeconds int,
	timeoutSeconds int,
) error {
	deadline := time.Now().Add(time.Duration(timeoutSeconds) * time.Second)
	for {
		_, err := client.GetLogStream(ctx, streamID)
		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code == 404 {
			return nil
		}
		// In replay mode, if VCR has no matching interaction the fixture doesn't
		// record a post-delete GET — treat any error as "stream gone".
		if err != nil && testingMode == recorder.ModeReplaying {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for log stream %d to be deleted", streamID)
		}
		if testingMode != recorder.ModeReplaying {
			time.Sleep(time.Duration(pollIntervalSeconds) * time.Second)
		}
	}
}

// waitForLogStreamProvisioned polls until the stream leaves provisioning state.
// Stream provisioning can take up to 60 minutes; timeoutSeconds should be set accordingly.
func waitForLogStreamProvisioned(
	ctx context.Context,
	client *linodego.Client,
	streamID int,
	pollIntervalSeconds int,
	timeoutSeconds int,
) (*linodego.Stream, error) {
	deadline := time.Now().Add(time.Duration(timeoutSeconds) * time.Second)
	for {
		stream, err := client.GetLogStream(ctx, streamID)
		if err != nil {
			return nil, err
		}
		if stream.Status == linodego.StreamStatusActive || stream.Status == linodego.StreamStatusInactive {
			return stream, nil
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timed out waiting for log stream %d to leave provisioning state", streamID)
		}
		if testingMode != recorder.ModeReplaying {
			time.Sleep(time.Duration(pollIntervalSeconds) * time.Second)
		}
	}
}

// ---- Logs Destination tests ----

func TestLogsDestination_List(t *testing.T) {
	client, dest, _, teardown := setupLogsDestination(t, "fixtures/TestLogsDestination_List")
	defer teardown()

	destinations, err := client.ListLogsDestinations(context.Background(), nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, destinations)

	for _, d := range destinations {
		assert.NotZero(t, d.ID)
		assert.NotEmpty(t, d.Label)
	}

	var found *linodego.LogsDestination
	for i := range destinations {
		if destinations[i].ID == dest.ID {
			found = &destinations[i]
			break
		}
	}
	require.NotNil(t, found, "created destination not found in list")
	assert.Equal(t, linodego.LogsDestinationTypeAkamaiObjectStorage, found.Type)
	assert.Contains(
		t,
		[]linodego.LogsDestinationStatus{
			linodego.LogsDestinationStatusActive,
			linodego.LogsDestinationStatusInactive,
		},
		found.Status,
	)
	assert.NotEmpty(t, found.Details.AccessKeyID)
	assert.NotEmpty(t, found.Details.BucketName)
	assert.NotEmpty(t, found.Details.Host)
}

func TestLogsDestination_Delete(t *testing.T) {
	client, dest, _, teardown := setupLogsDestination(t, "fixtures/TestLogsDestination_Delete")
	defer teardown()

	err := client.DeleteLogsDestination(context.Background(), dest.ID)
	assert.NoError(t, err)

	// Verify it's gone
	_, err = client.GetLogsDestination(context.Background(), dest.ID)
	require.Error(t, err)

	apiErr, ok := err.(*linodego.Error)
	require.True(t, ok, "expected linodego.Error")
	assert.Equal(t, 404, apiErr.Code)
}

func TestLogsDestination_Get(t *testing.T) {
	client, dest, _, teardown := setupLogsDestination(t, "fixtures/TestLogsDestination_Get")
	defer teardown()

	fetched, err := client.GetLogsDestination(context.Background(), dest.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, dest.ID, fetched.ID)
	assert.Equal(t, dest.Label, fetched.Label)
	assert.Equal(t, dest.Type, fetched.Type)
}

func TestLogsDestination_UpdateAndHistory(t *testing.T) {
	client, dest, storageKey, teardown := setupLogsDestination(t, "fixtures/TestLogsDestination_UpdateAndHistory")
	defer teardown()

	newLabel := dest.Label + "-upd"
	newPath := "updated/logs/path/"

	// should update logs destination
	updated, err := client.UpdateLogsDestination(context.Background(), dest.ID, linodego.LogsDestinationUpdateOptions{
		Label: newLabel,
		Details: &linodego.LogsDestinationDetailsUpdateOptions{
			AccessKeyID:     dest.Details.AccessKeyID,
			AccessKeySecret: storageKey.SecretKey,
			BucketName:      dest.Details.BucketName,
			Host:            dest.Details.Host,
			Path:            &newPath,
		},
	})
	assert.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, newLabel, updated.Label)
	assert.Equal(t, newPath, updated.Details.Path)

	// history should contain both versions
	history, err := client.ListLogsDestinationHistory(context.Background(), dest.ID, nil)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(history), 2)

	var v1, v2 *linodego.LogsDestination
	for i := range history {
		switch history[i].Version {
		case 1:
			v1 = &history[i]
		case 2:
			v2 = &history[i]
		}
	}

	require.NotNil(t, v1, "expected version 1 in history")
	require.NotNil(t, v2, "expected version 2 in history")

	assert.Equal(t, dest.Label, v1.Label)
	assert.Equal(t, newLabel, v2.Label)
}

func TestLogsDestination_TrafficPeakCreate(t *testing.T) {
	_, destination, teardown := setupTrafficPeakLogsDestination(t, "fixtures/TestLogsDestination_TrafficPeakCreate")
	defer teardown()

	assert.Equal(t, linodego.LogsDestinationTypeTrafficPeak, destination.Type)
	assert.Equal(t, "https://example.com", destination.Details.EndpointURL)
	assert.Equal(t, "application/json", destination.Details.ContentType)
	assert.Equal(t, "None", destination.Details.DataCompression)
}

func TestLogsDestination_TrafficPeakGet(t *testing.T) {
	client, destination, teardown := setupTrafficPeakLogsDestination(t, "fixtures/TestLogsDestination_TrafficPeakGet")
	defer teardown()

	fetched, err := client.GetLogsDestination(context.Background(), destination.ID)
	require.NoError(t, err)
	assert.Equal(t, destination.ID, fetched.ID)
	assert.Equal(t, linodego.LogsDestinationTypeTrafficPeak, fetched.Type)
	assert.Equal(t, destination.Details.EndpointURL, fetched.Details.EndpointURL)
}

func TestLogsDestination_TrafficPeakList(t *testing.T) {
	client, destination, teardown := setupTrafficPeakLogsDestination(t, "fixtures/TestLogsDestination_TrafficPeakList")
	defer teardown()

	destinations, err := client.ListLogsDestinations(context.Background(), nil)
	require.NoError(t, err)
	var found *linodego.LogsDestination
	for i := range destinations {
		if destinations[i].ID == destination.ID {
			found = &destinations[i]
			break
		}
	}
	require.NotNil(t, found)
	assert.Equal(t, linodego.LogsDestinationTypeTrafficPeak, found.Type)
}

func TestLogsDestination_TrafficPeakUpdateAndHistory(t *testing.T) {
	client, destination, teardown := setupTrafficPeakLogsDestination(t, "fixtures/TestLogsDestination_TrafficPeakUpdateAndHistory")
	defer teardown()

	updatedLabel := destination.Label + "-updated"
	updated, err := client.UpdateLogsDestination(context.Background(), destination.ID, linodego.LogsDestinationUpdateOptions{
		Label: updatedLabel,
		Details: &linodego.LogsDestinationTrafficPeakDetailsUpdateOptions{
			ContentType: linodego.Pointer("application/x-ndjson"),
		},
	})
	require.NoError(t, err)
	assert.Equal(t, updatedLabel, updated.Label)
	assert.Equal(t, "application/x-ndjson", updated.Details.ContentType)

	history, err := client.ListLogsDestinationHistory(context.Background(), destination.ID, nil)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(history), 2)
}

func TestLogsDestination_TrafficPeakDelete(t *testing.T) {
	client, destination, teardown := setupTrafficPeakLogsDestination(t, "fixtures/TestLogsDestination_TrafficPeakDelete")
	defer teardown()

	require.NoError(t, client.DeleteLogsDestination(context.Background(), destination.ID))
	_, err := client.GetLogsDestination(context.Background(), destination.ID)
	require.Error(t, err)
	apiErr, ok := err.(*linodego.Error)
	require.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.Code)
}

func TestLogsDestination_CustomHTTPSCreate(t *testing.T) {
	_, destination, teardown := setupCustomHTTPSLogsDestination(t, "fixtures/TestLogsDestination_CustomHTTPSCreate")
	defer teardown()

	assert.Equal(t, linodego.LogsDestinationTypeCustomHTTPS, destination.Type)
	assert.Equal(t, testCustomHTTPSEndpointURL, destination.Details.EndpointURL)
	assert.Equal(t, testCustomHTTPSContentType, destination.Details.ContentType)
	assert.Equal(t, testCustomHTTPSDataCompression, destination.Details.DataCompression)
	require.NotNil(t, destination.Details.Authentication)
	assert.Equal(t, linodego.LogsDestinationCustomHTTPSAuthTypeNone, destination.Details.Authentication.Type)
	assert.Nil(t, destination.Details.Authentication.Details)
	assert.Equal(t, []linodego.LogsDestinationCustomHTTPSHeader{
		{Name: testCustomHTTPSHeaderName, Value: testCustomHTTPSHeaderValue},
	}, destination.Details.CustomHeaders)
}

func TestLogsDestination_CustomHTTPSGet(t *testing.T) {
	client, destination, teardown := setupCustomHTTPSLogsDestination(t, "fixtures/TestLogsDestination_CustomHTTPSGet")
	defer teardown()

	fetched, err := client.GetLogsDestination(context.Background(), destination.ID)
	require.NoError(t, err)
	assert.Equal(t, destination.ID, fetched.ID)
	assert.Equal(t, linodego.LogsDestinationTypeCustomHTTPS, fetched.Type)
	assert.Equal(t, testCustomHTTPSEndpointURL, fetched.Details.EndpointURL)
	require.NotNil(t, fetched.Details.Authentication)
	assert.Equal(t, linodego.LogsDestinationCustomHTTPSAuthTypeNone, fetched.Details.Authentication.Type)
}

func TestLogsDestination_CustomHTTPSList(t *testing.T) {
	client, destination, teardown := setupCustomHTTPSLogsDestination(t, "fixtures/TestLogsDestination_CustomHTTPSList")
	defer teardown()

	destinations, err := client.ListLogsDestinations(context.Background(), &linodego.ListOptions{
		Filter: fmt.Sprintf(`{"id":%d}`, destination.ID),
	})
	require.NoError(t, err)
	require.Len(t, destinations, 1)
	assert.Equal(t, destination.ID, destinations[0].ID)
	assert.Equal(t, linodego.LogsDestinationTypeCustomHTTPS, destinations[0].Type)
	assert.Equal(t, testCustomHTTPSEndpointURL, destinations[0].Details.EndpointURL)
}

func TestLogsDestination_CustomHTTPSUpdateAndHistory(t *testing.T) {
	client, destination, teardown := setupCustomHTTPSLogsDestination(t, "fixtures/TestLogsDestination_CustomHTTPSUpdateAndHistory")
	defer teardown()

	updatedLabel := destination.Label + "-updated"
	updated, err := client.UpdateLogsDestination(context.Background(), destination.ID, linodego.LogsDestinationUpdateOptions{
		Label: updatedLabel,
		Details: &linodego.LogsDestinationCustomHTTPSDetailsUpdateOptions{
			ContentType: testCustomHTTPSUpdatedContentType,
			EndpointURL: testCustomHTTPSEndpointURL,
			Authentication: &linodego.LogsDestinationCustomHTTPSAuthDetails{
				Type: linodego.LogsDestinationCustomHTTPSAuthTypeNone,
			},
			DataCompression: testCustomHTTPSDataCompression,
			CustomHeaders: []linodego.LogsDestinationCustomHTTPSHeader{
				{Name: testCustomHTTPSHeaderName, Value: testCustomHTTPSHeaderValue},
			},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, updatedLabel, updated.Label)
	assert.Equal(t, testCustomHTTPSUpdatedContentType, updated.Details.ContentType)
	assert.Equal(t, testCustomHTTPSEndpointURL, updated.Details.EndpointURL)

	history, err := client.ListLogsDestinationHistory(context.Background(), destination.ID, nil)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(history), 2)
}

func TestLogsDestination_CustomHTTPSDelete(t *testing.T) {
	client, destination, teardown := setupCustomHTTPSLogsDestination(t, "fixtures/TestLogsDestination_CustomHTTPSDelete")
	defer teardown()

	require.NoError(t, client.DeleteLogsDestination(context.Background(), destination.ID))
	_, err := client.GetLogsDestination(context.Background(), destination.ID)
	require.Error(t, err)
	apiErr, ok := err.(*linodego.Error)
	require.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.Code)
}

type setupHTTPLogsDestinationFunc func(
	t *testing.T,
	fixturesYaml string,
) (*linodego.Client, *linodego.LogsDestination, func())

func testLogStreamHTTPDestination(
	t *testing.T,
	fixturesYaml string,
	setupDestination setupHTTPLogsDestinationFunc,
	expectedType linodego.StreamDestinationType,
) {
	t.Helper()
	requireACLPLogsStreamTests(t)

	client, destination, destinationTeardown := setupDestination(t, fixturesYaml)
	defer destinationTeardown()

	requireNoExistingStreams(t, client)

	stream, err := client.CreateLogStream(context.Background(), linodego.StreamCreateOptions{
		Destinations: []int{destination.ID},
		Label:        fmt.Sprintf("go-test-http-destination-stream-%d", time.Now().UnixNano()),
		Type:         linodego.StreamTypeAuditLogs,
	})
	require.NoError(t, err)
	require.NotNil(t, stream)

	defer func() {
		if err := client.DeleteLogStream(context.Background(), stream.ID); err != nil {
			t.Errorf("failed to delete HTTP destination stream: %v", err)
			return
		}
		if err := waitForLogStreamDeleted(context.Background(), client, stream.ID, 60, 3600); err != nil {
			t.Errorf("failed waiting for HTTP destination stream deletion: %v", err)
		}
	}()

	provisioned, err := waitForLogStreamProvisioned(context.Background(), client, stream.ID, 60, 3600)
	require.NoError(t, err)
	require.Len(t, provisioned.Destinations, 1)
	assert.Equal(t, destination.ID, provisioned.Destinations[0].ID)
	assert.Equal(t, expectedType, provisioned.Destinations[0].Type)
	assert.Equal(t, destination.Details.EndpointURL, provisioned.Destinations[0].Details.EndpointURL)
}

func TestLogStream_TrafficPeakDestination(t *testing.T) {
	testLogStreamHTTPDestination(
		t,
		"fixtures/TestLogStream_TrafficPeakDestination",
		setupTrafficPeakLogsDestination,
		linodego.StreamDestinationTypeTrafficPeak,
	)
}

func TestLogStream_CustomHTTPSDestination(t *testing.T) {
	testLogStreamHTTPDestination(
		t,
		"fixtures/TestLogStream_CustomHTTPSDestination",
		setupCustomHTTPSLogsDestination,
		linodego.StreamDestinationTypeCustomHTTPS,
	)
}

func TestLogsDestination_Create_InvalidSecret(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLogsDestination_Create_InvalidSecret")
	defer teardown()

	bucket, _, storageTeardown := setupObjectStorageForLogs(t, client)
	defer storageTeardown()

	_, err := client.CreateLogsDestination(context.Background(), linodego.LogsDestinationCreateOptions{
		Label: testLabel(),
		Type:  linodego.LogsDestinationTypeAkamaiObjectStorage,
		Details: linodego.LogsDestinationDetailsCreateOptions{
			AccessKeyID:     "1",
			AccessKeySecret: "1",
			BucketName:      bucket.Label,
			Host:            bucket.Hostname,
		},
	})
	require.Error(t, err)

	apiErr, ok := err.(*linodego.Error)
	require.True(t, ok, "expected linodego.Error")
	assert.Equal(t, 400, apiErr.Code)
	assert.Contains(t, apiErr.Message, "An error occurred")
}

func TestLogsDestination_Create_InvalidType(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLogsDestination_Create_InvalidType")
	defer teardown()

	bucket, storageKey, storageTeardown := setupObjectStorageForLogs(t, client)
	defer storageTeardown()

	_, err := client.CreateLogsDestination(context.Background(), linodego.LogsDestinationCreateOptions{
		Label: testLabel(),
		Type:  "invalid_type",
		Details: linodego.LogsDestinationDetailsCreateOptions{
			AccessKeyID:     storageKey.AccessKey,
			AccessKeySecret: storageKey.SecretKey,
			BucketName:      bucket.Label,
			Host:            bucket.Hostname,
		},
	})
	require.Error(t, err)

	apiErr, ok := err.(*linodego.Error)
	require.True(t, ok, "expected linodego.Error")
	assert.Equal(t, 400, apiErr.Code)
	assert.Contains(t, apiErr.Message, "[type] Must be one of akamai_object_storage, custom_https")
}

// ---- Log Stream tests ----

// requireNoExistingStreams skips the test if a stream already exists on the account,
// since only one stream is allowed per account.
func requireNoExistingStreams(t *testing.T, client *linodego.Client) {
	t.Helper()
	existing, err := client.ListLogStreams(context.Background(), nil)
	require.NoError(t, err)
	if len(existing) > 0 {
		ids := make([]int, len(existing))
		for i, s := range existing {
			ids[i] = s.ID
		}
		t.Skipf("existing stream(s) on account (IDs: %v); only one stream allowed per account", ids)
	}
}

func TestLogStream_Create_InvalidDestination(t *testing.T) {
	requireACLPLogsStreamTests(t)

	client, teardown := createTestClient(t, "fixtures/TestLogStream_Create_InvalidDestination")
	defer teardown()
	requireNoExistingStreams(t, client)

	_, err := client.CreateLogStream(context.Background(), linodego.StreamCreateOptions{
		Label:        fmt.Sprintf("go-test-invalid-dest-%d", time.Now().UnixNano()),
		Type:         linodego.StreamTypeAuditLogs,
		Destinations: []int{999999999},
	})
	require.Error(t, err)

	apiErr, ok := err.(*linodego.Error)
	require.True(t, ok, "expected *linodego.Error")
	assert.Equal(t, 400, apiErr.Code)
	assert.Contains(t, apiErr.Message, "[destination] Destination not found")
}

func TestLogStream_Create_EmptyDestinations(t *testing.T) {
	requireACLPLogsStreamTests(t)

	client, teardown := createTestClient(t, "fixtures/TestLogStream_Create_EmptyDestinations")
	defer teardown()
	requireNoExistingStreams(t, client)

	_, err := client.CreateLogStream(context.Background(), linodego.StreamCreateOptions{
		Label:        fmt.Sprintf("go-test-empty-dest-%d", time.Now().UnixNano()),
		Type:         linodego.StreamTypeAuditLogs,
		Destinations: []int{},
	})
	require.Error(t, err)

	apiErr, ok := err.(*linodego.Error)
	require.True(t, ok, "expected *linodego.Error")
	assert.Equal(t, 400, apiErr.Code)
	assert.Contains(t, apiErr.Message, "[destinations] Must contain one item")
}

func TestLogStream_Create_TwoDestinations(t *testing.T) {
	requireACLPLogsStreamTests(t)

	client, teardown := createTestClient(t, "fixtures/TestLogStream_Create_TwoDestinations")
	defer teardown()
	requireNoExistingStreams(t, client)

	bucket1, storageKey1, storageTeardown1 := setupObjectStorageForLogs(t, client)
	defer storageTeardown1()
	bucket2, storageKey2, storageTeardown2 := setupObjectStorageForLogs(t, client)
	defer storageTeardown2()

	dest1, err := client.CreateLogsDestination(context.Background(), linodego.LogsDestinationCreateOptions{
		Label: testLabel(),
		Type:  linodego.LogsDestinationTypeAkamaiObjectStorage,
		Details: linodego.LogsDestinationDetailsCreateOptions{
			AccessKeyID:     storageKey1.AccessKey,
			AccessKeySecret: storageKey1.SecretKey,
			BucketName:      bucket1.Label,
			Host:            bucket1.Hostname,
		},
	})
	require.NoError(t, err)
	defer func() { _ = client.DeleteLogsDestination(context.Background(), dest1.ID) }()

	dest2, err := client.CreateLogsDestination(context.Background(), linodego.LogsDestinationCreateOptions{
		Label: testLabel(),
		Type:  linodego.LogsDestinationTypeAkamaiObjectStorage,
		Details: linodego.LogsDestinationDetailsCreateOptions{
			AccessKeyID:     storageKey2.AccessKey,
			AccessKeySecret: storageKey2.SecretKey,
			BucketName:      bucket2.Label,
			Host:            bucket2.Hostname,
		},
	})
	require.NoError(t, err)
	defer func() { _ = client.DeleteLogsDestination(context.Background(), dest2.ID) }()

	_, err = client.CreateLogStream(context.Background(), linodego.StreamCreateOptions{
		Label:        fmt.Sprintf("go-test-two-dest-%d", time.Now().UnixNano()),
		Type:         linodego.StreamTypeAuditLogs,
		Destinations: []int{dest1.ID, dest2.ID},
	})
	require.Error(t, err)

	apiErr, ok := err.(*linodego.Error)
	require.True(t, ok, "expected *linodego.Error")
	assert.Equal(t, 400, apiErr.Code)
	assert.Contains(t, apiErr.Message, "[destinations] Must contain one item")
}

func TestLogStream_Delete(t *testing.T) {
	requireACLPLogsStreamTests(t)

	client, stream, _, teardown := setupLogStream(t, "fixtures/TestLogStream_Delete")
	defer teardown()

	provisioned, err := waitForLogStreamProvisioned(context.Background(), client, stream.ID, 60, 3600)
	require.NoError(t, err)

	err = client.DeleteLogStream(context.Background(), provisioned.ID)
	require.NoError(t, err)

	_, err = client.GetLogStream(context.Background(), provisioned.ID)
	require.Error(t, err)

	apiErr, ok := err.(*linodego.Error)
	require.True(t, ok, "expected *linodego.Error")
	assert.Equal(t, 404, apiErr.Code)
}

func TestLogStream_List(t *testing.T) {
	requireACLPLogsStreamTests(t)

	client, stream, _, teardown := setupLogStream(t, "fixtures/TestLogStream_List")
	defer teardown()

	provisioned, err := waitForLogStreamProvisioned(context.Background(), client, stream.ID, 60, 3600)
	require.NoError(t, err)

	streams, err := client.ListLogStreams(context.Background(), nil)
	require.NoError(t, err)
	assert.NotEmpty(t, streams)

	found := false
	for _, s := range streams {
		if s.ID == provisioned.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "created stream not found in list")
}

func TestLogStream_Get(t *testing.T) {
	requireACLPLogsStreamTests(t)

	client, stream, _, teardown := setupLogStream(t, "fixtures/TestLogStream_Get")
	defer teardown()

	provisioned, err := waitForLogStreamProvisioned(context.Background(), client, stream.ID, 60, 3600)
	require.NoError(t, err)

	fetched, err := client.GetLogStream(context.Background(), provisioned.ID)
	require.NoError(t, err)
	assert.Equal(t, provisioned.ID, fetched.ID)
	assert.Equal(t, provisioned.Label, fetched.Label)
	assert.Equal(t, provisioned.Status, fetched.Status)
	assert.Len(t, fetched.Destinations, 1)
}

func TestLogStream_Update_LabelAndStatus(t *testing.T) {
	requireACLPLogsStreamTests(t)

	client, stream, _, teardown := setupLogStream(t, "fixtures/TestLogStream_Update_LabelAndStatus")
	defer teardown()

	provisioned, err := waitForLogStreamProvisioned(context.Background(), client, stream.ID, 60, 3600)
	require.NoError(t, err)

	originalLabel := provisioned.Label
	originalStatus := provisioned.Status
	versionBefore := provisioned.Version

	newLabel := originalLabel + "-upd"
	// When active, requesting inactive triggers "deactivating" first; vice versa triggers "provisioning".
	// Accept both the requested status and its transitional counterpart.
	var newStatus linodego.StreamStatus
	var expectedStatuses []linodego.StreamStatus
	if originalStatus == linodego.StreamStatusInactive {
		newStatus = linodego.StreamStatusActive
		expectedStatuses = []linodego.StreamStatus{linodego.StreamStatusActive, linodego.StreamStatusProvisioning}
	} else {
		newStatus = linodego.StreamStatusInactive
		expectedStatuses = []linodego.StreamStatus{linodego.StreamStatusInactive, linodego.StreamStatusDeactivating}
	}

	_, err = client.UpdateLogStream(context.Background(), provisioned.ID, linodego.StreamUpdateOptions{
		Label:  &newLabel,
		Status: &newStatus,
	})
	require.NoError(t, err)

	updated, err := client.GetLogStream(context.Background(), provisioned.ID)
	require.NoError(t, err)
	assert.Equal(t, newLabel, updated.Label)
	assert.Contains(t, expectedStatuses, updated.Status)

	history, err := client.ListLogStreamHistory(context.Background(), provisioned.ID, nil)
	require.NoError(t, err)

	var snapOriginal, snapUpdated *linodego.Stream
	for i := range history {
		switch history[i].Version {
		case versionBefore:
			snapOriginal = &history[i]
		case updated.Version:
			snapUpdated = &history[i]
		}
	}
	require.NotNil(t, snapOriginal, "original version not found in history")
	require.NotNil(t, snapUpdated, "updated version not found in history")
	assert.Equal(t, originalLabel, snapOriginal.Label)
	assert.Equal(t, newLabel, snapUpdated.Label)
	assert.Equal(t, provisioned.ID, snapUpdated.ID)
}

func TestLogStream_Update_Destinations(t *testing.T) {
	requireACLPLogsStreamTests(t)

	client, stream, dest, teardown := setupLogStream(t, "fixtures/TestLogStream_Update_Destinations")
	defer teardown()

	objBucket2, objKey2, objTeardown2 := setupObjectStorageForLogs(t, client)
	defer objTeardown2()
	secondaryDest, err := client.CreateLogsDestination(context.Background(), linodego.LogsDestinationCreateOptions{
		Label: testLabel(),
		Type:  linodego.LogsDestinationTypeAkamaiObjectStorage,
		Details: linodego.LogsDestinationDetailsCreateOptions{
			AccessKeyID:     objKey2.AccessKey,
			AccessKeySecret: objKey2.SecretKey,
			BucketName:      objBucket2.Label,
			Host:            objBucket2.Hostname,
		},
	})
	require.NoError(t, err)
	defer func() {
		_ = client.DeleteLogsDestination(context.Background(), secondaryDest.ID)
	}()

	provisioned, err := waitForLogStreamProvisioned(context.Background(), client, stream.ID, 60, 3600)
	require.NoError(t, err)

	versionBefore := provisioned.Version

	updated, err := client.UpdateLogStream(context.Background(), provisioned.ID, linodego.StreamUpdateOptions{
		Destinations: []int{secondaryDest.ID},
	})
	require.NoError(t, err)
	require.Len(t, updated.Destinations, 1)
	assert.Equal(t, secondaryDest.ID, updated.Destinations[0].ID)

	defer func() {
		_, _ = client.UpdateLogStream(context.Background(), provisioned.ID, linodego.StreamUpdateOptions{
			Destinations: []int{dest.ID},
		})
	}()

	history, err := client.ListLogStreamHistory(context.Background(), provisioned.ID, nil)
	require.NoError(t, err)

	var snapOriginal, snapUpdated *linodego.Stream
	for i := range history {
		switch history[i].Version {
		case versionBefore:
			snapOriginal = &history[i]
		case updated.Version:
			snapUpdated = &history[i]
		}
	}
	require.NotNil(t, snapOriginal, "original version not found in history")
	require.NotNil(t, snapUpdated, "updated version not found in history")
	assert.Equal(t, dest.ID, snapOriginal.Destinations[0].ID)
	assert.Equal(t, secondaryDest.ID, snapUpdated.Destinations[0].ID)
}
