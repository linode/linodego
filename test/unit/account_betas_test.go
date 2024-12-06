package unit

import (
	"context"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountBetaProgram_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_beta_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/betas", fixtureData)

	betaPrograms, err := base.Client.ListAccountBetaPrograms(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.Len(t, betaPrograms, 1, "Expected exactly 1 beta program")

	betaProgram := betaPrograms[0]

	assert.Equal(t, "example_open", betaProgram.ID, "Expected beta program ID to be 'example_open'")
	assert.Equal(t, "Example Open Beta", betaProgram.Label, "Expected beta program label to be 'Example Open Beta'")
	assert.Equal(t, "This is an open public beta for an example feature.", betaProgram.Description, "Beta program description does not match")
	assert.Equal(t, "2023-07-11 00:00:00 +0000 UTC", betaProgram.Started.String(), "Expected beta program started date to be '2023-07-11T00:00:00'")
	assert.Equal(t, "2023-09-11 00:00:00 +0000 UTC", betaProgram.Enrolled.String(), "Expected beta program enrolled date to be '2023-09-11T00:00:00'")
	assert.Nil(t, betaProgram.Ended, "Expected beta program ended date to be nil")
}

func TestAccountBetaProgram_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_beta_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	betaId := "example_open"

	base.MockGet(fmt.Sprintf("account/betas/%s", betaId), fixtureData)

	betaProgram, err := base.Client.GetAccountBetaProgram(context.Background(), betaId)
	assert.NoError(t, err)

	assert.Equal(t, "example_open", betaProgram.ID, "Expected beta program ID to be 'example_open'")
	assert.Equal(t, "Example Open Beta", betaProgram.Label, "Expected beta program label to be 'Example Open Beta'.")
	assert.Equal(t, "This is an open public beta for an example feature.", betaProgram.Description, "Beta program description does not match")
	assert.Equal(t, "2023-07-11 00:00:00 +0000 UTC", betaProgram.Started.String(), "Expected beta program started date to be '2023-07-11T00:00:00'")
	assert.Equal(t, "2023-09-11 00:00:00 +0000 UTC", betaProgram.Enrolled.String(), "Expected beta program enrolled date to be '2023-09-11T00:00:00'")
	assert.Nil(t, betaProgram.Ended, "Expected beta program ended date to be nil")
}

func TestAccountBetaProgram_Join(t *testing.T) {
	client := createMockClient(t)

	betaId := "global_load_balancer_beta"

	requestData := linodego.AccountBetaProgramCreateOpts{
		ID: betaId,
	}

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "account/betas"),
		mockRequestBodyValidate(t, requestData, nil))

	if _, err := client.JoinBetaProgram(context.Background(), requestData); err != nil {
		t.Fatal(err)
	}
}
