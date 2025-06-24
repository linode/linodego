package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProfileLogin_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_login_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("profile/logins/456", fixtureData)

	login, err := base.Client.GetProfileLogin(context.Background(), 456)
	assert.NoError(t, err)
	assert.NotNil(t, login)

	expectedTime, _ := time.Parse(time.RFC3339, "2024-01-15T14:30:00Z")

	if assert.NotNil(t, login.Datetime) {
		assert.Equal(t, expectedTime, *login.Datetime)
	}
}

func TestProfileLogins_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_logins_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("profile/logins", fixtureData)

	logins, err := base.Client.ListProfileLogins(context.Background(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, logins)
	assert.Len(t, logins, 2)

	expectedTimes := []string{
		"2024-02-01T10:15:30Z",
		"2024-02-02T18:45:00Z",
	}

	for i, login := range logins {
		if assert.NotNil(t, login.Datetime) {
			expectedTime, _ := time.Parse(time.RFC3339, expectedTimes[i])
			assert.Equal(t, expectedTime, *login.Datetime)
		}
	}
}
