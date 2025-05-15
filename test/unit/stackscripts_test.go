package unit

import (
	"context"
	"fmt"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestListStackscripts(t *testing.T) {
	// Mock the API response to match the expected structure for a paginated response
	fixtureData, err := fixtures.GetFixture("stackscripts_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the request with a correct paginated structure
	base.MockGet("linode/stackscripts", fixtureData)

	stackscripts, err := base.Client.ListStackscripts(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.NotEmpty(t, stackscripts, "Expected non-empty stackscripts list")

	// Check if a specific stackscript exists using slices.ContainsFunc
	exists := slices.ContainsFunc(stackscripts, func(stackscript linodego.Stackscript) bool {
		return stackscript.Label == "Test Stackscript"
	})

	assert.True(t, exists, "Expected stackscripts list to contain 'Test Stackscript'")
}

func TestCreateStackscript(t *testing.T) {
	// Load the fixture data for stackscript creation
	fixtureData, err := fixtures.GetFixture("stackscript_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("linode/stackscripts", fixtureData)

	opts := linodego.StackscriptCreateOptions{
		Label:       "new-stackscript",
		Description: "A new stackscript",
		Images:      []string{"linode/ubuntu20.04"},
		IsPublic:    true,
		RevNote:     "Initial revision",
		Script:      "#!/bin/bash\necho Hello",
	}

	stackscript, err := base.Client.CreateStackscript(context.Background(), opts)
	assert.NoError(t, err, "Expected no error when creating stackscript")

	// Verify the created stackscript's label
	assert.Equal(t, "new-stackscript", stackscript.Label, "Expected created stackscript label to match input")
}

func TestDeleteStackscript(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	stackscriptID := 123
	base.MockDelete(fmt.Sprintf("linode/stackscripts/%d", stackscriptID), nil)

	err := base.Client.DeleteStackscript(context.Background(), stackscriptID)
	assert.NoError(t, err, "Expected no error when deleting stackscript")
}

func TestGetStackscript(t *testing.T) {
	// Load the fixture data for a single stackscript
	fixtureData, err := fixtures.GetFixture("stackscript_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	stackscriptID := 123
	base.MockGet(fmt.Sprintf("linode/stackscripts/%d", stackscriptID), fixtureData)

	stackscript, err := base.Client.GetStackscript(context.Background(), stackscriptID)
	assert.NoError(t, err)

	// Verify the stackscript's label
	assert.Equal(t, "new-stackscript", stackscript.Label, "Expected stackscript label to match fixture")
}

func TestUpdateStackscript(t *testing.T) {
	// Load the fixture data for stackscript update
	fixtureData, err := fixtures.GetFixture("stackscript_revision")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	stackscriptID := 123
	base.MockPut(fmt.Sprintf("linode/stackscripts/%d", stackscriptID), fixtureData)

	opts := linodego.StackscriptUpdateOptions{
		Label:       "Updated Stackscript",
		Description: "Updated description",
		Images:      []string{"linode/ubuntu20.04"},
		IsPublic:    false,
		RevNote:     "Updated revision",
		Script:      "#!/bin/bash\necho Hello Updated",
	}

	updatedStackscript, err := base.Client.UpdateStackscript(context.Background(), stackscriptID, opts)
	assert.NoError(t, err)

	// Verify the updated stackscript's label
	assert.Equal(t, "Updated Stackscript", updatedStackscript.Label, "Expected updated stackscript label to match input")
}
