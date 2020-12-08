package integration

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/linode/linodego"
)

var (
	objectStorageObjectURLExpirySeconds = 360
)

func putObjectStorageObject(t *testing.T, client *linodego.Client, bucket *linodego.ObjectStorageBucket, name, content string) {
	t.Helper()

	url, err := client.CreateObjectStorageObjectURL(context.TODO(), bucket.Cluster, bucket.Label, linodego.ObjectStorageObjectURLCreateOptions{
		Name:        name,
		Method:      http.MethodPut,
		ContentType: "text/plain",
		ExpiresIn:   &objectStorageObjectURLExpirySeconds,
	})
	if err != nil {
		t.Errorf("failed to get object PUT url: %s", err)
	}

	rec, teardownRecorder := testRecorder(t, "fixtures/TestGetObjectStorageObjectACLConfigBucketClientPut", testingMode, nil)
	defer teardownRecorder()

	httpClient := http.Client{Transport: rec}
	req, err := http.NewRequest(http.MethodPost, url.URL, bytes.NewReader([]byte(content)))
	if err != nil {
		t.Errorf("failed to build request: %s", err)
	}
	req.Method = http.MethodPut
	req.Header.Add("Content-Type", "text/plain")

	res, err := httpClient.Do(req)
	if err != nil {
		t.Errorf("failed to make request: %s", err)
	}

	if res.StatusCode != 200 {
		t.Errorf("expected status code to be 200; got %d", res.StatusCode)
	}
}

func deleteObjectStorageObject(t *testing.T, client *linodego.Client, bucket *linodego.ObjectStorageBucket, name string) {
	t.Helper()

	url, err := client.CreateObjectStorageObjectURL(context.TODO(), bucket.Cluster, bucket.Label, linodego.ObjectStorageObjectURLCreateOptions{
		Name:      name,
		Method:    http.MethodDelete,
		ExpiresIn: &objectStorageObjectURLExpirySeconds,
	})
	if err != nil {
		t.Errorf("failed to get object PUT url: %s", err)
	}

	rec, teardownRecorder := testRecorder(t, "fixtures/TestGetObjectStorageObjectACLConfigBucketClientDelete", testingMode, nil)
	defer teardownRecorder()

	httpClient := http.Client{Transport: rec}
	req, err := http.NewRequest(http.MethodPost, url.URL, nil)
	if err != nil {
		t.Errorf("failed to build request: %s", err)
	}
	req.Method = http.MethodDelete

	res, err := httpClient.Do(req)
	if res.StatusCode != 204 {
		t.Errorf("expected status code to be 204; got %d", res.StatusCode)
	}
}

func TestUpdateObjectStorageObjectACLConfig(t *testing.T) {
	client, bucket, teardown, err := setupObjectStorageBucket(t, "fixtures/TestGetObjectStorageObjectACLConfig")
	if err != nil {
		t.Fatalf("failed to create Object Storage Object: %s", err)
	}
	defer teardown()

	object := "test"
	putObjectStorageObject(t, client, bucket, object, "testing123")
	defer deleteObjectStorageObject(t, client, bucket, object)

	config, err := client.GetObjectStorageObjectACLConfig(context.TODO(), bucket.Cluster, bucket.Label, object)
	if err != nil {
		t.Errorf("failed to get ACL config: %s", err)
	}

	if config.ACL != "private" {
		t.Errorf("expected ACL to be private; got %s", config.ACL)
	}
	if config.ACLXML == "" {
		t.Error("expected ACL XML to be included")
	}

	updateOpts := linodego.ObjectStorageObjectACLConfigUpdateOptions{ACL: "public-read", Name: object}
	if _, err = client.UpdateObjectStorageObjectACLConfig(context.TODO(), bucket.Cluster, bucket.Label, updateOpts); err != nil {
		t.Errorf("failed to update ACL config: %s", err)
	}

	config, err = client.GetObjectStorageObjectACLConfig(context.TODO(), bucket.Cluster, bucket.Label, object)
	if err != nil {
		t.Errorf("failed to get updated ACL config: %s", err)
	}

	if config.ACL != updateOpts.ACL {
		t.Errorf("expected ACL config to be %s; got %s", updateOpts.ACL, config.ACL)
	}
	if config.ACLXML == "" {
		t.Error("expected ACL XML to be included")
	}
}
