package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestProfilePreferences_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_preferences_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("profile/preferences", fixtureData)

	preferences, err := base.Client.GetProfilePreferences(context.Background())
	assert.NoError(t, err)

	expectedPreferences := linodego.ProfilePreferences{
		"key1": "value1",
		"key2": "value2",
	}
	assert.Equal(t, expectedPreferences, *preferences)
}

func TestProfilePreferences_update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_preferences_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.ProfilePreferences{
		"key1": "value1_new",
		"key2": "value2_new",
	}

	base.MockPut("profile/preferences", fixtureData)

	preferences, err := base.Client.UpdateProfilePreferences(context.Background(), requestData)
	assert.NoError(t, err)

	expectedPreferences := linodego.ProfilePreferences{
		"key1": "value1_new",
		"key2": "value2_new",
	}
	assert.Equal(t, expectedPreferences, *preferences)
}
