package unit

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestObjectStorage_Cancel(t *testing.T) {
	client := createMockClient(t)

	desiredResponse := make(map[string]interface{})

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "/object-storage/cancel"),
		httpmock.NewJsonResponderOrPanic(200, &desiredResponse))

	err := client.CancelObjectStorage(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestObjectStorage_ObjectList(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_buckets_object_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("object-storage/buckets/us-east/bucket-name/object-list", fixtureData)

	content, err := base.Client.ListObjectStorageBucketContents(context.Background(), "us-east", "bucket-name", nil)
	if err != nil {
		t.Fatalf("Error getting content: %v", err)
	}

	assert.Equal(t, false, content.IsTruncated)
	assert.Nil(t, content.NextMarker)
	assert.Equal(t, "example", content.Data[0].Name)
	assert.Equal(t, "bfc70ab2-e3d4-42a4-ad55-83921822270c", content.Data[0].Owner)
	assert.Equal(t, 123, content.Data[0].Size)
	assert.Equal(t, "9f254c71e28e033bf9e0e5262e3e72ab", content.Data[0].Etag)
}
