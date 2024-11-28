package unit

import (
	"context"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProfile_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("profile", fixtureData)

	profile, err := base.Client.GetProfile(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, "password", profile.AuthenticationType)
	assert.Equal(t, "example-user@gmail.com", profile.Email)
	assert.Equal(t, true, profile.EmailNotifications)
	assert.Equal(t, linodego.LishAuthMethod("keys_only"), profile.LishAuthMethod)
	assert.Equal(t, "871be32f49c1411b14f29f618aaf0c14637fb8d3", profile.Referrals.Code)
	assert.Equal(t, 0, profile.Referrals.Completed)
	assert.Equal(t, float64(0), profile.Referrals.Credit)
	assert.Equal(t, 0, profile.Referrals.Pending)
	assert.Equal(t, 0, profile.Referrals.Total)
	assert.Equal(t, "https://www.linode.com/?r=871be32f49c1411b14f29f618aaf0c14637fb8d3", profile.Referrals.URL)
	assert.Equal(t, false, profile.Restricted)
	assert.Equal(t, "US/Eastern", profile.Timezone)
	assert.Equal(t, true, profile.TwoFactorAuth)
	assert.Equal(t, 1234, profile.UID)
	assert.Equal(t, "exampleUser", profile.Username)
	assert.Equal(t, "+5555555555", profile.VerifiedPhoneNumber)
}

func TestProfile_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.ProfileUpdateOptions{
		Email: "example-user-new@gmail.com",
	}

	base.MockPut("profile", fixtureData)

	profile, err := base.Client.UpdateProfile(context.Background(), requestData)
	assert.NoError(t, err)

	assert.Equal(t, "password", profile.AuthenticationType)
	assert.Equal(t, "example-user-new@gmail.com", profile.Email)
	assert.Equal(t, true, profile.EmailNotifications)
	assert.Equal(t, linodego.LishAuthMethod("keys_only"), profile.LishAuthMethod)
	assert.Equal(t, "871be32f49c1411b14f29f618aaf0c14637fb8d3", profile.Referrals.Code)
	assert.Equal(t, 0, profile.Referrals.Completed)
	assert.Equal(t, float64(0), profile.Referrals.Credit)
	assert.Equal(t, 0, profile.Referrals.Pending)
	assert.Equal(t, 0, profile.Referrals.Total)
	assert.Equal(t, "https://www.linode.com/?r=871be32f49c1411b14f29f618aaf0c14637fb8d3", profile.Referrals.URL)
	assert.Equal(t, false, profile.Restricted)
	assert.Equal(t, "US/Eastern", profile.Timezone)
	assert.Equal(t, true, profile.TwoFactorAuth)
	assert.Equal(t, 1234, profile.UID)
	assert.Equal(t, "exampleUser", profile.Username)
	assert.Equal(t, "+5555555555", profile.VerifiedPhoneNumber)
}
