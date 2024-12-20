package unit

import (
	"context"
	"fmt"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLinodeTypes_List(t *testing.T) {
	// Load the fixture data for types
	fixtureData, err := fixtures.GetFixture("linode_types_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/types", fixtureData)

	types, err := base.Client.ListTypes(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	// Verify a specific type exists in the list
	var nanodeType *linodego.LinodeType
	for _, t := range types {
		if t.ID == "g6-nanode-1" {
			nanodeType = &t
			break
		}
	}

	if nanodeType == nil {
		t.Errorf("Expected type 'g6-nanode-1' to be in the response, but it was not found")
	} else {
		assert.Equal(t, "nanode", string(nanodeType.Class), "Expected class to be 'nanode'")
		assert.Equal(t, 1, nanodeType.VCPUs, "Expected VCPUs for 'g6-nanode-1' to be 1")
		assert.Equal(t, 250, nanodeType.Transfer, "Expected transfer for 'g6-nanode-1' to be 250GB")
		assert.NotNil(t, nanodeType.Price, "Expected 'g6-nanode-1' to have a price object")
		if nanodeType.Price != nil {
			assert.Equal(t, float32(5), nanodeType.Price.Monthly, "Expected monthly price for 'g6-nanode-1' to be $5")
		}
	}
}

func TestLinodeType_Get(t *testing.T) {
	// Load the fixture data for a specific type
	fixtureData, err := fixtures.GetFixture("linode_type_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	typeID := "g6-standard-2"
	base.MockGet(fmt.Sprintf("linode/types/%s", typeID), fixtureData)

	typeObj, err := base.Client.GetType(context.Background(), typeID)
	assert.NoError(t, err)

	assert.Equal(t, typeID, typeObj.ID, "Expected type ID to match")
	assert.Equal(t, "standard", string(typeObj.Class), "Expected class to be 'standard'")
	assert.Equal(t, 2, typeObj.VCPUs, "Expected VCPUs to be 2")
	assert.Equal(t, 4000, typeObj.Disk, "Expected disk to be 4000MB")
	assert.Equal(t, 4000, typeObj.Memory, "Expected memory to be 4000MB")
	assert.NotNil(t, typeObj.Price, "Expected type to have a price object")
	if typeObj.Price != nil {
		assert.Equal(t, float32(10), typeObj.Price.Monthly, "Expected monthly price to be $10")
	}

	assert.NotNil(t, typeObj.Addons, "Expected type to have addons")
	if typeObj.Addons != nil && typeObj.Addons.Backups != nil {
		assert.NotNil(t, typeObj.Addons.Backups.Price, "Expected backups to have a price object")
	}
}
