package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, 3, len(endpoints))
	assert.Equal(t, "us-east-1", endpoints[0].Region)
	assert.Equal(t, "https://s3.us-east-1.linodeobjects.com", *endpoints[0].S3Endpoint)
	assert.Equal(t, linodego.ObjectStorageEndpointE0, endpoints[0].EndpointType)
	assert.Equal(t, []linodego.ObjectStorageKeyClass{linodego.ObjectStorageKeyClassE}, endpoints[0].SupportedKeyClasses)
}

func TestObjectStorageEndpoint_ListIA(t *testing.T) {
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

	ia := endpoints[2]
	assert.Equal(t, "us-iad-obj", ia.Region)
	assert.Equal(t, linodego.ObjectStorageEndpointIA0, ia.EndpointType)
	assert.Equal(t, []linodego.ObjectStorageKeyClass{linodego.ObjectStorageKeyClassIA}, ia.SupportedKeyClasses)
	assert.Equal(t, "us", ia.Country)
	assert.Equal(t, "Washington, DC", ia.RegionLabel)
}
