package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestObjectStorageKey_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_key_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("object-storage/keys", fixtureData)

	keys, err := base.Client.ListObjectStorageKeys(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, keys, 1)

	assert.Equal(t, "my-key", keys[0].Label)
	assert.Equal(t, "my-access-key", keys[0].AccessKey)
	assert.Equal(t, "my-secret-key", keys[0].SecretKey)
	assert.True(t, keys[0].Limited)
}

func TestObjectStorageKey_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_key_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	createOptions := linodego.ObjectStorageKeyCreateOptions{
		Label:   "new-key",
		Regions: []string{"us-east-1"},
	}

	base.MockPost("object-storage/keys", fixtureData)

	key, err := base.Client.CreateObjectStorageKey(context.Background(), createOptions)
	assert.NoError(t, err)
	assert.Equal(t, "new-key", key.Label)
	assert.Equal(t, "new-access-key", key.AccessKey)
	assert.Equal(t, "new-secret-key", key.SecretKey)
}

func TestObjectStorageKey_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_key_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("object-storage/keys/123", fixtureData)

	key, err := base.Client.GetObjectStorageKey(context.Background(), 123)
	assert.NoError(t, err)
	assert.Equal(t, "my-key", key.Label)
	assert.Equal(t, "my-access-key", key.AccessKey)
	assert.Equal(t, "my-secret-key", key.SecretKey)
}

func TestObjectStorageKey_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_key_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	updateOptions := linodego.ObjectStorageKeyUpdateOptions{
		Label:   "updated-key",
		Regions: []string{"us-west-1"},
	}

	base.MockPut("object-storage/keys/123", fixtureData)

	key, err := base.Client.UpdateObjectStorageKey(context.Background(), 123, updateOptions)
	assert.NoError(t, err)
	assert.Equal(t, "updated-key", key.Label)
	assert.Equal(t, "updated-access-key", key.AccessKey)
}

func TestObjectStorageKey_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("object-storage/keys/123", nil)

	err := base.Client.DeleteObjectStorageKey(context.Background(), 123)
	assert.NoError(t, err)
}

func TestObjectStorageKey_ListByRegion(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_key_list_by_region")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("object-storage/keys", fixtureData)

	keys, err := base.Client.ListObjectStorageKeys(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, keys, 2)

	var filteredKeys []linodego.ObjectStorageKey
	for _, key := range keys {
		for _, region := range key.Regions {
			if region.ID == "us-east-1" {
				filteredKeys = append(filteredKeys, key)
			}
		}
	}

	assert.Len(t, filteredKeys, 2)
	assert.Equal(t, "us-east-1", filteredKeys[0].Regions[0].ID)
}
