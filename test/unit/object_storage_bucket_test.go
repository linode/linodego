package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestObjectStorageBucket_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_bucket_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("object-storage/buckets", fixtureData)

	buckets, err := base.Client.ListObjectStorageBuckets(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, buckets, 1)

	assert.Equal(t, "my-bucket", buckets[0].Label)
	assert.Equal(t, "us-east-1", buckets[0].Region)
	assert.Equal(t, "https://s3.us-east-1.linodeobjects.com", buckets[0].S3Endpoint)
	assert.Equal(t, linodego.ObjectStorageEndpointType("public"), buckets[0].EndpointType)
	assert.Equal(t, "my-bucket.us-east-1.linodeobjects.com", buckets[0].Hostname)
	assert.Equal(t, 5, buckets[0].Objects)
	assert.Equal(t, 10240, buckets[0].Size)
}

//func TestObjectStorageBucket_ListInRegion(t *testing.T) {
//	fixtureData, err := fixtures.GetFixture("object_storage_bucket_list")
//	assert.NoError(t, err)
//
//	var base ClientBaseCase
//	base.SetUp(t)
//	defer base.TearDown(t)
//
//	regionID := "us-east"
//	base.MockGet("object-storage/buckets/"+regionID, fixtureData)
//
//	buckets, err := base.Client.ListObjectStorageBucketsInRegion(context.Background(), nil, regionID)
//	assert.NoError(t, err)
//	assert.NotEmpty(t, buckets)
//}

func TestObjectStorageBucket_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_bucket_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	clusterID := "us-east-1"
	bucketLabel := "my-bucket"

	base.MockGet("object-storage/buckets/"+clusterID+"/"+bucketLabel, fixtureData)

	bucket, err := base.Client.GetObjectStorageBucket(context.Background(), clusterID, bucketLabel)
	assert.NoError(t, err)
	assert.NotNil(t, bucket)
	assert.Equal(t, bucketLabel, bucket.Label)
}

func TestObjectStorageBucket_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_bucket_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	createOpts := linodego.ObjectStorageBucketCreateOptions{
		Region: linodego.Pointer("us-east"),
		Label:  "new-bucket",
	}

	base.MockPost("object-storage/buckets", fixtureData)

	bucket, err := base.Client.CreateObjectStorageBucket(context.Background(), createOpts)
	assert.NoError(t, err)
	assert.NotNil(t, bucket)
	assert.Equal(t, createOpts.Label, bucket.Label)
}

func TestObjectStorageBucket_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	clusterID := "us-east-1"
	bucketLabel := "my-bucket"

	base.MockDelete("object-storage/buckets/"+clusterID+"/"+bucketLabel, nil)

	err := base.Client.DeleteObjectStorageBucket(context.Background(), clusterID, bucketLabel)
	assert.NoError(t, err)
}

func TestObjectStorageBucket_GetAccess(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_bucket_access_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	clusterID := "us-east-1"
	bucketLabel := "my-bucket"

	base.MockGet("object-storage/buckets/"+clusterID+"/"+bucketLabel+"/access", fixtureData)

	access, err := base.Client.GetObjectStorageBucketAccessV2(context.Background(), clusterID, bucketLabel)
	assert.NoError(t, err)
	assert.NotNil(t, access)
	assert.Equal(t, linodego.ACLPublicRead, access.ACL)
}

func TestObjectStorageBucket_UpdateAccess(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	clusterID := "us-east-1"
	bucketLabel := "my-bucket"

	updateOpts := linodego.ObjectStorageBucketUpdateAccessOptions{
		ACL: linodego.Pointer(linodego.ACLPrivate),
	}

	base.MockPost("object-storage/buckets/"+clusterID+"/"+bucketLabel+"/access", nil)

	err := base.Client.UpdateObjectStorageBucketAccess(context.Background(), clusterID, bucketLabel, updateOpts)
	assert.NoError(t, err)
}

func TestObjectStorageBucket_ListContents(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_bucket_contents")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	clusterID := "us-east-1"
	bucketLabel := "my-bucket"

	base.MockGet("object-storage/buckets/"+clusterID+"/"+bucketLabel+"/object-list", fixtureData)

	contents, err := base.Client.ListObjectStorageBucketContents(context.Background(), clusterID, bucketLabel, nil)
	assert.NoError(t, err)
	assert.NotNil(t, contents)
	assert.True(t, contents.IsTruncated)
}
