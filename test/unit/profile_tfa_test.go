package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
)

func TestTwoFactor_CreateSecret_smoke(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_two_factor_secret_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("profile/tfa-enable", fixtureData)

	secret, err := base.Client.CreateTwoFactorSecret(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, time.Time(time.Date(2018, time.March, 1, 0, 1, 1, 0, time.UTC)), *secret.Expiry)
	assert.Equal(t, "5FXX6KLACOC33GTC", secret.Secret)
}

func TestTwoFactor_Disable(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "/profile/tfa-disable"),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.DisableTwoFactor(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestTwoFactor_Confirm(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_two_factor_enable")
	assert.NoError(t, err)

	request := linodego.ConfirmTwoFactorOptions{TFACode: "reallycoolandlegittfacode"}

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("profile/tfa-enable-confirm", fixtureData)

	response, err := base.Client.ConfirmTwoFactor(context.Background(), request)
	assert.NoError(t, err)

	assert.Equal(t, "reallycoolandlegittfacode", response.Scratch)
}
