package integration

import (
	"context"
	"testing"

	. "github.com/linode/linodego"
)

var testObjectStorageBucketCreateOpts = ObjectStorageBucketCreateOptions{
	Cluster: "us-east-1",
	Label:   "go-bucket-test-def",
}

func TestObjectStorageBucket_Create(t *testing.T) {
	_, bucket, teardown, err := setupObjectStorageBucket(t,
		nil,
		"fixtures/TestObjectStorageBucket_Create")
	defer teardown()

	if err != nil {
		t.Errorf("Error creating Object Storage Bucket, got error %v", err)
	}

	expected := testObjectStorageBucketCreateOpts

	// when comparing fixtures to random value Label will differ, compare the known prefix
	if bucket.Label != expected.Label ||
		bucket.Cluster != expected.Cluster {
		t.Errorf("Object Storage Bucket did not match CreateOptions")
	}

	assertDateSet(t, bucket.Created)
}

func TestObjectStorageBucket_GetMissing(t *testing.T) {
	client, bucket, teardown, err := setupObjectStorageBucket(t,
		nil,
		"fixtures/TestObjectStorageBucket_GetMissing")
	defer teardown()

	sameLabel := bucket.Label
	differentCluster := "us-west-1"

	i, err := client.GetObjectStorageBucket(context.Background(), differentCluster, sameLabel)
	if err == nil {
		t.Errorf("should have received an error requesting a missing ObjectStorageBucket, got %v", i)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing ObjectStorageBucket, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing ObjectStorageBucket, got %v", e.Code)
	}
}

func TestObjectStorageBucket_GetFound(t *testing.T) {
	client, bucket, teardown, err := setupObjectStorageBucket(t,
		nil,
		"fixtures/TestObjectStorageBucket_GetFound")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	i, err := client.GetObjectStorageBucket(context.Background(), bucket.Cluster, bucket.Label)
	if err != nil {
		t.Errorf("Error getting ObjectStorageBucket, expected struct, got %v and error %v", i, err)
	}
	if i.Label != bucket.Label {
		t.Errorf("Expected a specific ObjectStorageBucket, but got a different one %v", i)
	}
	expected := testObjectStorageBucketCreateOpts

	// when comparing fixtures to random value Label will differ, compare the known prefix
	if bucket.Label != expected.Label ||
		bucket.Cluster != expected.Cluster {
		t.Errorf("Object Storage Bucket did not match CreateOptions")
	}
}

func TestObjectStorageBuckets_List(t *testing.T) {
	client, _, teardown, err := setupObjectStorageBucket(t,
		nil,
		"fixtures/TestObjectStorageBuckets_List")
	defer teardown()

	i, err := client.ListObjectStorageBuckets(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing ObjectStorageBuckets, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of ObjectStorageBuckets, but got none %v", i)
	} else if i[0].Label == "" ||
		i[0].Cluster == "" {
		t.Errorf("Listed Object Storage Bucket did not have attribuets %v", i)
	}
}

func TestObjectStorageBucket_Access_Get(t *testing.T) {
	corsEnabled := false

	createOpts := ObjectStorageBucketCreateOptions{
		ACL:         ACLAuthenticatedRead,
		CorsEnabled: &corsEnabled,
	}

	client, bucket, teardown, err := setupObjectStorageBucket(t,
		[]objectStorageBucketModifier{
			func(opts *ObjectStorageBucketCreateOptions) {
				opts.ACL = createOpts.ACL
				opts.CorsEnabled = createOpts.CorsEnabled
			},
		},
		"fixtures/TestObjectStorageBucket_Access_Get")
	defer teardown()

	newBucket, err := client.GetObjectStorageBucketAccess(context.Background(), bucket.Cluster, bucket.Label)
	if err != nil {
		t.Errorf("Error getting ObjectStorageBucket access, got error %s", err)
	}

	if newBucket.CorsEnabled != corsEnabled {
		t.Errorf("ObjectStorageBucket access CORS does not match update, expected %t, got %t", corsEnabled, newBucket.CorsEnabled)
	}

	if newBucket.ACL != createOpts.ACL {
		t.Errorf("ObjectStorageBucket access ACL does not match update, expected %s, got %s",
			createOpts.ACL,
			newBucket.ACL)
	}
}

func TestObjectStorageBucket_Access_Update(t *testing.T) {
	client, bucket, teardown, err := setupObjectStorageBucket(t,
		nil,
		"fixtures/TestObjectStorageBucket_Access_Update")
	defer teardown()

	corsEnabled := false

	opts := ObjectStorageBucketUpdateAccessOptions{
		ACL:         ACLPrivate,
		CorsEnabled: &corsEnabled,
	}

	err = client.UpdateObjectStorageBucketAccess(context.Background(), bucket.Cluster, bucket.Label, opts)
	if err != nil {
		t.Errorf("Error updating ObjectStorageBucket access, got error %s", err)
	}

	newBucket, err := client.GetObjectStorageBucketAccess(context.Background(), bucket.Cluster, bucket.Label)
	if err != nil {
		t.Errorf("Error getting ObjectStorageBucket access, got error %s", err)
	}

	if newBucket.CorsEnabled != corsEnabled {
		t.Errorf("ObjectStorageBucket access CORS does not match update, expected %t, got %t", corsEnabled, newBucket.CorsEnabled)
	}

	if newBucket.ACL != opts.ACL {
		t.Errorf("ObjectStorageBucket access ACL does not match update, expected %s, got %s", opts.ACL, newBucket.ACL)
	}
}

type objectStorageBucketModifier func(*ObjectStorageBucketCreateOptions)

func setupObjectStorageBucket(t *testing.T, bucketModifiers []objectStorageBucketModifier, fixturesYaml string) (*Client, *ObjectStorageBucket, func(), error) {
	t.Helper()

	createOpts := testObjectStorageBucketCreateOpts

	for _, modifier := range bucketModifiers {
		modifier(&createOpts)
	}

	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	bucket, err := client.CreateObjectStorageBucket(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error creating test Bucket: %s", err)
	}

	teardown := func() {
		if err := client.DeleteObjectStorageBucket(context.Background(), bucket.Cluster, bucket.Label); err != nil {
			if t != nil {
				t.Errorf("Error deleting test Bucket: %s", err)
			}
		}
		fixtureTeardown()
	}

	return client, bucket, teardown, err
}
