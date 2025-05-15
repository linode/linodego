package unit

import (
	"context"
	"testing"
	"time"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestProfileTokens_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_tokens_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("profile/tokens", fixtureData)

	tokens, err := base.Client.ListTokens(context.Background(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, tokens)
	assert.Len(t, tokens, 2)

	expectedTimes := []string{
		"2024-03-10T12:00:00Z",
		"2024-03-11T15:45:00Z",
	}

	for i, token := range tokens {
		if assert.NotNil(t, token.Created) {
			expectedTime, _ := time.Parse(time.RFC3339, expectedTimes[i])
			assert.Equal(t, expectedTime, *token.Created)
		}
	}
}

func TestProfileToken_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_token_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("profile/tokens/123", fixtureData)

	token, err := base.Client.GetToken(context.Background(), 123)
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, "My API Token", token.Label)
}

func TestProfileToken_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_token_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("profile/tokens", fixtureData)

	opts := linodego.TokenCreateOptions{
		Label:  "New API Token",
		Scopes: "read_write",
	}

	token, err := base.Client.CreateToken(context.Background(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, "New API Token", token.Label)
}

func TestProfileToken_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_token_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("profile/tokens/123", fixtureData)

	opts := linodego.TokenUpdateOptions{
		Label: "Updated API Token",
	}

	token, err := base.Client.UpdateToken(context.Background(), 123, opts)
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, "Updated API Token", token.Label)
}

func TestProfileToken_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("profile/tokens/123", nil)

	err := base.Client.DeleteToken(context.Background(), 123)
	assert.NoError(t, err)
}
