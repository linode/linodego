package integration

import (
	"context"
	"slices"
	"testing"

	"github.com/linode/linodego"
	. "github.com/linode/linodego"
)

var objectStorageBucketTestLabel = "go-bucket-test-def"

var testObjectStorageBucketCreateOpts = ObjectStorageBucketCreateOptions{
	Cluster: "us-east-1",
	Label:   objectStorageBucketTestLabel,
}

var testRegionalObjectStorageBucketCreateOpts = ObjectStorageBucketCreateOptions{
	Region: "us-east",
	Label:  objectStorageBucketTestLabel,
}

func TestObjectStorageBucket_Create_smoke(t *testing.T) {
	_, bucket, teardown, err := setupObjectStorageBucket(t,
		nil,
		"fixtures/TestObjectStorageBucket_Create", nil, nil, nil)
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

func TestObjectStorageBucket_Regional(t *testing.T) {
	// t.Skip("skipping region test before GA")
	client, teardown := createTestClient(t, "fixtures/TestObjectStorageBucket_Regional")
	regions := getRegionsWithCaps(t, client, []string{"Object Storage"})
	if len(regions) < 1 {
		t.Fatal("Can't get region with Object Storage capability")
	}
	region := regions[0]

	client, bucket, teardown, err := setupObjectStorageBucket(t,
		[]objectStorageBucketModifier{
			func(opts *ObjectStorageBucketCreateOptions) {
				opts.Cluster = ""
				opts.Region = region
			},
		},
		"fixtures/TestObjectStorageBucket_Regional",
		client, teardown, nil,
	)
	defer teardown()

	if err != nil {
		t.Errorf("Error creating Object Storage Bucket, got error %v", err)
	}

	expected := testObjectStorageBucketCreateOpts

	// when comparing fixtures to random value Label will differ, compare the known prefix
	if bucket.Label != expected.Label ||
		bucket.Region != region {
		t.Errorf("Object Storage Bucket did not match CreateOptions")
	}

	assertDateSet(t, bucket.Created)

	bucket, err = client.GetObjectStorageBucket(context.Background(), region, expected.Label)
	if err != nil {
		t.Errorf("Error getting Object Storage Bucket, %v", err)
	}
}

func TestObjectStorageBucket_GetMissing(t *testing.T) {
	client, bucket, teardown, err := setupObjectStorageBucket(t,
		nil,
		"fixtures/TestObjectStorageBucket_GetMissing", nil, nil, nil)
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
		"fixtures/TestObjectStorageBucket_GetFound", nil, nil, nil)
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

func TestObjectStorageBuckets_List_smoke(t *testing.T) {
	client, _, teardown, err := setupObjectStorageBucket(t,
		nil,
		"fixtures/TestObjectStorageBuckets_List", nil, nil, nil)
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

func TestObjectStorageBucketsInCluster_List(t *testing.T) {
	client, bucket, teardown, err := setupObjectStorageBucket(t,
		nil,
		"fixtures/TestObjectStorageBucketsInCluster_List", nil, nil, nil)
	defer teardown()

	i, err := client.ListObjectStorageBucketsInCluster(context.Background(), nil, bucket.Cluster)
	if err != nil {
		t.Errorf("Error listing ObjectStorageBucketsInCluster, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of ObjectStorageBucketsInCluster, but got none %v", i)
	} else if i[0].Label == "" ||
		i[0].Cluster == "" {
		t.Errorf("Listed Object Storage Bucket in Cluster did not have attribuets %v", i)
	}
}

func TestObjectStorageBucket_Access_Get(t *testing.T) {
	corsEnabled := false

	createOpts := ObjectStorageBucketCreateOptions{
		ACL:         ACLAuthenticatedRead,
		CorsEnabled: &corsEnabled,
	}
	endpointType := linodego.ObjectStorageEndpointE1
	client, bucket, teardown, err := setupObjectStorageBucket(t,
		[]objectStorageBucketModifier{
			func(opts *ObjectStorageBucketCreateOptions) {
				opts.ACL = createOpts.ACL
				opts.CorsEnabled = createOpts.CorsEnabled
			},
		},
		"fixtures/TestObjectStorageBucket_Access_Get", nil, nil,
		&endpointType,
	)
	defer teardown()

	newBucket, err := client.GetObjectStorageBucketAccess(context.Background(), bucket.Region, bucket.Label)
	if err != nil {
		t.Errorf("Error getting ObjectStorageBucket access, got error %s", err)
	}

	newBucketv2, err := client.GetObjectStorageBucketAccessV2(context.Background(), bucket.Region, bucket.Label)
	if err != nil {
		t.Errorf("Error getting ObjectStorageBucket access, got error %s", err)
	}

	if newBucket.CorsEnabled != corsEnabled {
		t.Errorf("ObjectStorageBucket access CORS does not match update, expected %t, got %t", corsEnabled, newBucket.CorsEnabled)
	}

	if newBucketv2.CorsEnabled == nil {
		t.Errorf("ObjectStorageBucket access CORS does not match update, expected %t, got nil", corsEnabled)
	}

	if newBucketv2.CorsEnabled != nil && *newBucketv2.CorsEnabled != corsEnabled {
		t.Errorf("ObjectStorageBucket access CORS does not match update, expected %t, got %t", corsEnabled, *newBucketv2.CorsEnabled)
	}

	if newBucket.ACL != createOpts.ACL {
		t.Errorf("ObjectStorageBucket access ACL does not match update, expected %s, got %s",
			createOpts.ACL,
			newBucket.ACL)
	}
	if newBucketv2.ACL != createOpts.ACL {
		t.Errorf("ObjectStorageBucket access ACL does not match update, expected %s, got %s",
			createOpts.ACL,
			newBucketv2.ACL)
	}
}

func TestObjectStorageBucket_Access_Update(t *testing.T) {
	endpointType := linodego.ObjectStorageEndpointE1
	client, bucket, teardown, err := setupObjectStorageBucket(t,
		nil,
		"fixtures/TestObjectStorageBucket_Access_Update",
		nil, nil, &endpointType,
	)
	defer teardown()

	corsEnabled := false

	opts := ObjectStorageBucketUpdateAccessOptions{
		ACL:         ACLPrivate,
		CorsEnabled: &corsEnabled,
	}

	err = client.UpdateObjectStorageBucketAccess(context.Background(), bucket.Region, bucket.Label, opts)
	if err != nil {
		t.Errorf("Error updating ObjectStorageBucket access, got error %s", err)
	}

	newBucket, err := client.GetObjectStorageBucketAccess(context.Background(), bucket.Region, bucket.Label)
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

func setupObjectStorageBucket(
	t *testing.T,
	bucketModifiers []objectStorageBucketModifier,
	fixturesYaml string,
	client *Client,
	teardown func(),
	endpointType *linodego.ObjectStorageEndpointType,
) (*Client, *ObjectStorageBucket, func(), error) {
	t.Helper()

	if (client == nil) != (teardown == nil) {
		t.Fatalf(
			"The client and fixtureTeardown variables must either both be nil or both " +
				"have a value. They cannot have one set to nil and the other set to a non-nil value.",
		)
	}

	if client == nil {
		client, teardown = createTestClient(t, fixturesYaml)
	}

	createOpts := testRegionalObjectStorageBucketCreateOpts

	if endpointType != nil {
		endpoints, err := client.ListObjectStorageEndpoints(context.Background(), nil)
		if err != nil {
			t.Fatalf("Error listing endpoints: %s", err)
		} else {
			selectedEndpoint := endpoints[slices.IndexFunc(endpoints, func(e linodego.ObjectStorageEndpoint) bool {
				return e.EndpointType == linodego.ObjectStorageEndpointE1
			})]
			createOpts.Region = selectedEndpoint.Region
			createOpts.EndpointType = selectedEndpoint.EndpointType
		}

	}

	for _, modifier := range bucketModifiers {
		modifier(&createOpts)
	}

	bucket, err := client.CreateObjectStorageBucket(context.Background(), createOpts)
	if err != nil {
		t.Fatalf("Error creating test Bucket: %s", err)
	}

	newTeardown := func() {
		if err := client.DeleteObjectStorageBucket(context.Background(), bucket.Cluster, bucket.Label); err != nil {
			if t != nil {
				t.Errorf("Error deleting test Bucket: %s", err)
			}
		}
		teardown()
	}

	return client, bucket, newTeardown, err
}
