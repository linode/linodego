package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestObjectStorageBucketCert_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_bucket_cert")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	clusterID := "us-east-1"
	bucketName := "my-bucket"

	base.MockGet("object-storage/buckets/"+clusterID+"/"+bucketName+"/ssl", fixtureData)

	cert, err := base.Client.GetObjectStorageBucketCertV2(context.Background(), clusterID, bucketName)
	assert.NoError(t, err)
	assert.NotNil(t, cert)
	assert.True(t, *cert.SSL)
}

func TestObjectStorageBucketCert_Upload(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_bucket_cert")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	clusterID := "us-east-1"
	bucketName := "my-bucket"

	uploadOpts := linodego.ObjectStorageBucketCertUploadOptions{
		Certificate: "mock-cert",
		PrivateKey:  "mock-key",
	}

	base.MockPost("object-storage/buckets/"+clusterID+"/"+bucketName+"/ssl", fixtureData)

	uploadedCert, err := base.Client.UploadObjectStorageBucketCertV2(context.Background(), clusterID, bucketName, uploadOpts)
	assert.NoError(t, err)
	assert.NotNil(t, uploadedCert)
	assert.True(t, *uploadedCert.SSL)
}

func TestObjectStorageBucketCertV2_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_bucket_cert")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	clusterID := "us-east-1"
	bucketName := "my-bucket"

	base.MockGet("object-storage/buckets/"+clusterID+"/"+bucketName+"/ssl", fixtureData)

	cert, err := base.Client.GetObjectStorageBucketCertV2(context.Background(), clusterID, bucketName)
	assert.NoError(t, err)
	assert.NotNil(t, cert)
	assert.NotNil(t, cert.SSL)
	assert.True(t, *cert.SSL)
}

func TestObjectStorageBucketCertV2_Upload(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_bucket_cert")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	clusterID := "us-east-1"
	bucketName := "my-bucket"

	uploadOpts := linodego.ObjectStorageBucketCertUploadOptions{
		Certificate: "mock-cert",
		PrivateKey:  "mock-key",
	}

	base.MockPost("object-storage/buckets/"+clusterID+"/"+bucketName+"/ssl", fixtureData)

	uploadedCert, err := base.Client.UploadObjectStorageBucketCertV2(context.Background(), clusterID, bucketName, uploadOpts)
	assert.NoError(t, err)
	assert.NotNil(t, uploadedCert)
	assert.NotNil(t, uploadedCert.SSL)
	assert.True(t, *uploadedCert.SSL)
}

func TestObjectStorageBucketCert_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	clusterID := "us-east-1"
	bucketName := "my-bucket"

	base.MockDelete("object-storage/buckets/"+clusterID+"/"+bucketName+"/ssl", nil)

	err := base.Client.DeleteObjectStorageBucketCert(context.Background(), clusterID, bucketName)
	assert.NoError(t, err)
}
