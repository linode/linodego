package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestObjectStorageObjectURL_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_object_url_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	createOptions := linodego.ObjectStorageObjectURLCreateOptions{
		Name:        "test-object",
		Method:      "GET",
		ExpiresIn:   nil,
		ContentType: "application/json",
	}

	base.MockPost("object-storage/buckets/my-bucket/test-label/object-url", fixtureData)

	urlResponse, err := base.Client.CreateObjectStorageObjectURL(context.Background(), "my-bucket", "test-label", createOptions)
	assert.NoError(t, err)
	assert.Equal(t, "https://s3.example.com/my-bucket/test-object", urlResponse.URL)
	assert.True(t, urlResponse.Exists)
}

func TestObjectStorageObjectACLConfigV2_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_object_acl_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("object-storage/buckets/my-bucket/test-label/object-acl?name=test-object", fixtureData)

	aclConfig, err := base.Client.GetObjectStorageObjectACLConfigV2(context.Background(), "my-bucket", "test-label", "test-object")
	assert.NoError(t, err)
	assert.NotNil(t, aclConfig.ACL)
	assert.Equal(t, "public-read", *aclConfig.ACL)
}

func TestObjectStorageObjectACLConfigV2_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_object_acl_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	updateOptions := linodego.ObjectStorageObjectACLConfigUpdateOptions{
		Name: "test-object",
		ACL:  "private",
	}

	base.MockPut("object-storage/buckets/my-bucket/test-label/object-acl", fixtureData)

	updatedACLConfig, err := base.Client.UpdateObjectStorageObjectACLConfigV2(context.Background(), "my-bucket", "test-label", updateOptions)
	assert.NoError(t, err)
	assert.NotNil(t, updatedACLConfig.ACL)
	assert.Equal(t, "private", *updatedACLConfig.ACL)
}
