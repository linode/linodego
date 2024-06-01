package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/linode/linodego"
)

var testBasicObjectStorageKeyCreateOpts = ObjectStorageKeyCreateOptions{
	Label: label,
}

func TestObjectStorageKey_GetMissing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestObjectStorageKey_GetMissing")
	defer teardown()

	notfoundID := 123
	i, err := client.GetObjectStorageKey(context.Background(), notfoundID)
	if err == nil {
		t.Errorf("should have received an error requesting a missing object storage key, got %v", i)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing object storage key, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing object storage key, got %v", e.Code)
	}
}

func TestObjectStorageKey_GetFound(t *testing.T) {
	client, objectStorageKey, teardown, err := setupObjectStorageKey(t, testBasicObjectStorageKeyCreateOpts, "fixtures/TestObjectStorageKey_GetFound", nil, nil)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	i, err := client.GetObjectStorageKey(context.Background(), objectStorageKey.ID)
	if err != nil {
		t.Errorf("Error getting objectStorageKey, expected struct, got %v and error %v", i, err)
	}
	if i.ID != objectStorageKey.ID {
		t.Errorf("Expected objectStorageKey id %d, but got %d", i.ID, objectStorageKey.ID)
	}
	if testBasicObjectStorageKeyCreateOpts.Label != objectStorageKey.Label {
		t.Errorf("Expected objectStorageKey label '%s', but got '%s'", testBasicObjectStorageKeyCreateOpts.Label, objectStorageKey.Label)
	}
	if objectStorageKey.BucketAccess != nil || objectStorageKey.Limited {
		t.Errorf("Expected objectStorageKey to have full permissions, but got %v, %v", objectStorageKey.Limited, objectStorageKey.BucketAccess)
	}
}

func TestObjectStorageKey_Update(t *testing.T) {
	client, objectStorageKey, teardown, err := setupObjectStorageKey(t, testBasicObjectStorageKeyCreateOpts, "fixtures/TestObjectStorageKey_Update", nil, nil)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	renamedLabel := objectStorageKey.Label + "_r"
	updateOpts := ObjectStorageKeyUpdateOptions{
		Label: renamedLabel,
	}
	objectStorageKey, err = client.UpdateObjectStorageKey(context.Background(), objectStorageKey.ID, updateOpts)
	if err != nil {
		t.Errorf("Error renaming objectStorageKey, %s", err)
	}

	if !strings.Contains(objectStorageKey.Label, renamedLabel) {
		t.Errorf("objectStorageKey returned does not match objectStorageKey update request, %v", objectStorageKey)
	}
	if objectStorageKey.BucketAccess != nil || objectStorageKey.Limited {
		t.Errorf("Expected objectStorageKey to have full permissions, but got %v, %v", objectStorageKey.Limited, objectStorageKey.BucketAccess)
	}
}

func TestObjectStorageKeys_List(t *testing.T) {
	client, objkey, teardown, err := setupObjectStorageKey(t, testBasicObjectStorageKeyCreateOpts, "fixtures/TestObjectStorageKey_List", nil, nil)
	defer teardown()
	if err != nil {
		t.Error(err)
	}
	objectStorageKeys, err := client.ListObjectStorageKeys(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing objectStorageKeys, expected struct, got error %v", err)
	}
	if len(objectStorageKeys) == 0 {
		t.Errorf("Expected a list of objectStorageKeys, but got %v", objectStorageKeys)
	}

	notFound := true
	for i := range objectStorageKeys {
		if objectStorageKeys[i].Label == objkey.Label {
			notFound = false
			break
		}
	}
	if notFound {
		t.Errorf("Expected to find created objectStorageKey, but '%s' was not found", objkey.Label)
	}
}

func TestObjectStorageKeys_Limited(t *testing.T) {
	_, bucket, teardown, err := setupObjectStorageBucket(t, nil, "fixtures/TestObjectStorageKeys_Limited_Bucket", nil, nil)
	defer teardown()

	createOpts := testBasicObjectStorageKeyCreateOpts
	createOpts.BucketAccess = &[]ObjectStorageKeyBucketAccess{
		{
			Cluster:     "us-east-1",
			BucketName:  bucket.Label,
			Permissions: "read_only",
		},
		{
			Cluster:     "us-east-1",
			BucketName:  bucket.Label,
			Permissions: "read_write",
		},
	}

	_, objectStorageKey, teardown, err := setupObjectStorageKey(t, createOpts, "fixtures/TestObjectStorageKeys_Limited", nil, nil)
	defer teardown()
	if err != nil {
		t.Error(err)
	}
	if !objectStorageKey.Limited || !cmp.Equal(objectStorageKey.BucketAccess, createOpts.BucketAccess) {
		t.Errorf("objectStorageKey returned (%v) does not match objectStorageKey creation request (%v)", *objectStorageKey.BucketAccess, *createOpts.BucketAccess)
	}
}

func TestObjectStorageKeys_Limited_NoAccess(t *testing.T) {
	createOpts := testBasicObjectStorageKeyCreateOpts
	createOpts.BucketAccess = &[]ObjectStorageKeyBucketAccess{}

	_, objectStorageKey, teardown, err := setupObjectStorageKey(t, createOpts, "fixtures/TestObjectStorageKeys_Limited_NoAccess", nil, nil)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	if !objectStorageKey.Limited || objectStorageKey.BucketAccess == nil || len(*objectStorageKey.BucketAccess) != 0 {
		t.Errorf("objectStorageKey returned access, %v, %v", objectStorageKey.Limited, objectStorageKey.BucketAccess)
	}
}

func TestObjectStorageKeys_Regional_Limited(t *testing.T) {
	// t.Skip("skipping region test before GA")
	client, teardown := createTestClient(t, "fixtures/TestObjectStorageKeys_Regional_Limited")
	regions := getRegionsWithCaps(t, client, []string{"Object Storage"})
	if len(regions) < 1 {
		t.Fatal("Can't get region with Object Storage capability")
	}
	region := regions[0]

	client, bucket, teardown, err := setupObjectStorageBucket(t, []objectStorageBucketModifier{
		func(createOpts *ObjectStorageBucketCreateOptions) {
			createOpts.Cluster = ""
			createOpts.Region = region
		},
	}, "fixtures/TestObjectStorageKeys_Regional_Limited", client, teardown)
	if err != nil {
		t.Error(err)
	}

	createOpts := testBasicObjectStorageKeyCreateOpts
	createOpts.BucketAccess = &[]ObjectStorageKeyBucketAccess{
		{
			Region:      region,
			BucketName:  bucket.Label,
			Permissions: "read_only",
		},
	}
	initialRegion := "us-east"
	createOpts.Regions = []string{initialRegion}
	_, key, teardown, err := setupObjectStorageKey(t, createOpts, "fixtures/TestObjectStorageKeys_Regional_Limited", client, teardown)
	defer teardown()
	if err != nil {
		t.Fatalf("error creating the obj regional key: %v", err)
	}

	if !key.Limited || key.BucketAccess == nil || len(*key.BucketAccess) == 0 {
		t.Errorf("Regional limited Object Storage key returned access, %v, %v", key.Limited, key.BucketAccess)
	}

	if len(key.Regions) == 0 || key.Regions[0].ID != initialRegion {
		t.Errorf("Unexpected key regions, expected regions: %v, actual regions: %v", createOpts.Regions, key.Regions)
	}

	updatedRegion := "us-east"
	updateOpts := ObjectStorageKeyUpdateOptions{
		Regions: []string{updatedRegion},
	}
	key, err = client.UpdateObjectStorageKey(context.Background(), key.ID, updateOpts)
	if err != nil {
		t.Fatalf("error updating the obj regional key: %v", err)
	}

	if len(key.Regions) == 0 || key.Regions[0].ID != updatedRegion {
		t.Errorf("Unexpected key regions, expected regions: %v, actual regions: %v", updatedRegion, key.Regions[0].ID)
	}
}

func setupObjectStorageKey(t *testing.T, createOpts ObjectStorageKeyCreateOptions, fixturesYaml string, client *Client, teardown func()) (*Client, *ObjectStorageKey, func(), error) {
	t.Helper()

	if (client == nil) != (teardown == nil) {
		t.Error(
			"The client and fixtureTeardown variables must either both be nil or both " +
				"have a value. They cannot have one set to nil and the other set to a non-nil value.",
		)
	}

	if client == nil {
		client, teardown = createTestClient(t, fixturesYaml)
	}

	objectStorageKey, err := client.CreateObjectStorageKey(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error creating ObjectStorageKey: %v", err)
	}

	newTeardown := func() {
		if err := client.DeleteObjectStorageKey(context.Background(), objectStorageKey.ID); err != nil {
			t.Errorf("Expected to delete a objectStorageKey, but got %v", err)
		}
		teardown()
	}
	return client, objectStorageKey, newTeardown, err
}
