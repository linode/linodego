package unit

import (
	"context"
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestListTags(t *testing.T) {
	// Load the fixture data for tags
	fixtureData, err := fixtures.GetFixture("tags_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("tags", fixtureData)

	tags, err := base.Client.ListTags(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.NotEmpty(t, tags, "Expected non-empty tag list")

	// Check if a specific tag exists using slices.ContainsFunc
	exists := slices.ContainsFunc(tags, func(tag linodego.Tag) bool {
		return tag.Label == "example-tag"
	})

	assert.True(t, exists, "Expected tag list to contain 'example-tag'")
}

func TestCreateTag(t *testing.T) {
	// Load the fixture data for tag creation
	fixtureData, err := fixtures.GetFixture("tag_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("tags", fixtureData)

	opts := linodego.TagCreateOptions{
		Label: "new-tag",
	}

	tag, err := base.Client.CreateTag(context.Background(), opts)
	assert.NoError(t, err, "Expected no error when creating tag")

	// Verify the created tag's label
	assert.Equal(t, "new-tag", tag.Label, "Expected created tag label to match input")
}

func TestDeleteTag(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	tagLabel := "delete-tag"
	base.MockDelete(fmt.Sprintf("tags/%s", tagLabel), nil)

	err := base.Client.DeleteTag(context.Background(), tagLabel)
	assert.NoError(t, err, "Expected no error when deleting tag")
}

func TestListTaggedObjects(t *testing.T) {
	// Load the fixture data for tagged objects
	fixtureData, err := fixtures.GetFixture("tagged_objects_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	tagLabel := "example-tag"
	base.MockGet(fmt.Sprintf("tags/%s", tagLabel), fixtureData)

	taggedObjects, err := base.Client.ListTaggedObjects(context.Background(), tagLabel, &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.NotEmpty(t, taggedObjects, "Expected non-empty tagged objects list")

	// Find the specific tagged object using slices.IndexFunc
	index := slices.IndexFunc(taggedObjects, func(obj linodego.TaggedObject) bool {
		return obj.Type == "linode"
	})

	assert.NotEqual(t, -1, index, "Expected to find a tagged object of type 'linode'")
	if index != -1 {
		assert.Equal(t, "linode", taggedObjects[index].Type, "Expected tagged object type to be 'linode'")
	}
}

func TestSortedObjects(t *testing.T) {
	// Load the fixture data for tagged objects
	fixtureData, err := fixtures.GetFixture("tagged_objects_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	tagLabel := "example-tag"
	base.MockGet(fmt.Sprintf("tags/%s", tagLabel), fixtureData)

	taggedObjects, err := base.Client.ListTaggedObjects(context.Background(), tagLabel, &linodego.ListOptions{})
	assert.NoError(t, err)

	sortedObjects, err := taggedObjects.SortedObjects()
	assert.NoError(t, err)

	assert.NotEmpty(t, sortedObjects.Instances, "Expected non-empty instances list in sorted objects")
	assert.Equal(t, "example-instance", sortedObjects.Instances[0].Label, "Expected instance label to be 'example-instance'")
}

func TestCreateTagWithReservedIPv4Addresses(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("tag_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("tags", fixtureData)

	opts := linodego.TagCreateOptions{
		Label:                 "new-tag",
		ReservedIPv4Addresses: []string{"192.168.1.10", "192.168.1.20"},
	}

	tag, err := base.Client.CreateTag(context.Background(), opts)
	assert.NoError(t, err, "Expected no error when creating tag with reserved IPv4 addresses")
	assert.Equal(t, "new-tag", tag.Label)
}

func TestListTaggedObjectsWithReservedIPv4Address(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("tagged_objects_reserved_ip_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	tagLabel := "lb"
	base.MockGet(fmt.Sprintf("tags/%s", tagLabel), fixtureData)

	taggedObjects, err := base.Client.ListTaggedObjects(context.Background(), tagLabel, &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.Len(t, taggedObjects, 1)

	assert.Equal(t, "reserved_ipv4_address", taggedObjects[0].Type)

	ip, ok := taggedObjects[0].Data.(linodego.InstanceIP)
	assert.True(t, ok, "Expected Data to be InstanceIP")
	assert.Equal(t, "192.168.1.10", ip.Address)
	assert.Equal(t, "192.0.2.1", ip.Gateway)
	assert.Equal(t, 24, ip.Prefix)
	assert.True(t, ip.Public)
	assert.Equal(t, "", ip.RDNS)
	assert.Equal(t, "us-east", ip.Region)
	assert.Equal(t, "255.255.255.0", ip.SubnetMask)
	assert.Equal(t, linodego.InstanceIPType("ipv4"), ip.Type)
	assert.Nil(t, ip.VPCNAT1To1)
	assert.True(t, ip.Reserved)
	assert.Equal(t, []string{"lb"}, ip.Tags)
}

func TestCreateTag_ReservedIPv4AddressesRequestBody(t *testing.T) {
	client := createMockClient(t)

	opts := linodego.TagCreateOptions{
		Label:                 "new-tag",
		ReservedIPv4Addresses: []string{"192.0.2.141"},
	}

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "/tags"),
		mockRequestBodyValidate(t, opts, nil))

	if _, err := client.CreateTag(context.Background(), opts); err != nil {
		t.Fatal(err)
	}
}
