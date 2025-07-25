package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestObjectStorageEndpoint_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_endpoints_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("object-storage/endpoints", fixtureData)

	endpoints, err := base.Client.ListObjectStorageEndpoints(context.Background(), nil)
	if err != nil {
		t.Fatalf("Error getting endpoints: %v", err)
	}

	assert.Equal(t, 2, len(endpoints))
	assert.Equal(t, "us-east-1", endpoints[0].Region)
	assert.Equal(t, "https://s3.us-east-1.linodeobjects.com", *endpoints[0].S3Endpoint)
	assert.Equal(t, linodego.ObjectStorageEndpointE0, endpoints[0].EndpointType)
}
