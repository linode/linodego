package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// creates a object storage and access keys for use in tests
func setupObjectStorageForLogs(t *testing.T, client *linodego.Client) (*linodego.ObjectStorageBucket, *linodego.ObjectStorageKey, func()) {
	t.Helper()

	bucket, err := client.CreateObjectStorageBucket(context.Background(), linodego.ObjectStorageBucketCreateOptions{
		Region:      "us-southeast",
		Label:       testLabel(),
		ACL:         "private",
		CorsEnabled: linodego.Pointer(false),
	})
	if err != nil {
		t.Fatalf("Error creating storage bucket, got error %v", err)
	}

	storageKey, err := client.CreateObjectStorageKey(context.Background(), linodego.ObjectStorageKeyCreateOptions{
		Label: testLabel(),
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
				url, err := client.CreateObjectStorageObjectURL(context.TODO(), bucket.Cluster, bucket.Label, linodego.ObjectStorageObjectURLCreateOptions{
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

func testLabel() string {
	return "go-test-logs-destination-" + getUniqueText()
}

func TestLogsDestination_List(t *testing.T) {
	client, dest, _, teardown := setupLogsDestination(t, "fixtures/TestLogsDestination_List")
	defer teardown()

	destinations, err := client.ListLogsDestinations(context.Background(), nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, destinations)

	for _, d := range destinations {
		assert.NotZero(t, d.ID)
		assert.NotEmpty(t, d.Label)
		assert.Equal(t, linodego.LogsDestinationTypeAkamaiObjectStorage, d.Type)
		assert.Contains(t,
			[]linodego.LogsDestinationStatus{
				linodego.LogsDestinationStatusActive,
				linodego.LogsDestinationStatusInactive,
			},
			d.Status,
		)
		assert.NotEmpty(t, d.Details.AccessKeyID)
		assert.NotEmpty(t, d.Details.BucketName)
		assert.NotEmpty(t, d.Details.Host)
	}

	ids := make([]int, len(destinations))
	for i, d := range destinations {
		ids[i] = d.ID
	}

	assert.Contains(t, ids, dest.ID)
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
	assert.Contains(t, apiErr.Message, "Invalid access key id or secret key")
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
