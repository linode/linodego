package unit

import (
	"context"
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestAccountOauthClient_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_oauth_client_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/oauth-clients", fixtureData)

	oauthClients, err := base.Client.ListOAuthClients(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	// Assertions on the returned data
	assert.Len(t, oauthClients, 1, "Expected one OAuth client")

	client := oauthClients[0]
	assert.Equal(t, "2737bf16b39ab5d7b4a1", client.ID, "Unexpected ID")
	assert.Equal(t, "Test_Client_1", client.Label, "Unexpected Label")
	assert.False(t, client.Public, "Unexpected Public value")
	assert.Equal(t, "https://example.org/oauth/callback", client.RedirectURI, "Unexpected Redirect URI")
	assert.Equal(t, linodego.OAuthClientStatus("active"), client.Status, "Unexpected Status")
	assert.Equal(t, "https://api.linode.com/v4/account/clients/2737bf16b39ab5d7b4a1/thumbnail", *client.ThumbnailURL, "Unexpected Thumbnail URL")
}

func TestAccountOauthClient_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_oauth_client_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	clientID := "2737bf16b39ab5d7b4a1"
	base.MockGet(fmt.Sprintf("account/oauth-clients/%s", clientID), fixtureData)

	oauthClient, err := base.Client.GetOAuthClient(context.Background(), clientID)
	assert.NoError(t, err)
	// Assertions on the returned data
	assert.Equal(t, "2737bf16b39ab5d7b4a1", oauthClient.ID, "Unexpected ID")
	assert.Equal(t, "Test_Client_1", oauthClient.Label, "Unexpected Label")
	assert.False(t, oauthClient.Public, "Unexpected Public value")
	assert.Equal(t, "https://example.org/oauth/callback", oauthClient.RedirectURI, "Unexpected Redirect URI")
	assert.Equal(t, linodego.OAuthClientStatus("active"), oauthClient.Status, "Unexpected Status")
	assert.Equal(t, "https://api.linode.com/v4/account/clients/2737bf16b39ab5d7b4a1/thumbnail", *oauthClient.ThumbnailURL, "Unexpected Thumbnail URL")
}

func TestAccountOauthClient_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_oauth_client_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.OAuthClientCreateOptions{
		Label:       "Test_Client_1",
		RedirectURI: "https://example.org/oauth/callback",
	}

	base.MockPost("account/oauth-clients", fixtureData)

	oauthClient, err := base.Client.CreateOAuthClient(context.Background(), requestData)
	assert.NoError(t, err)
	// Assertions on the returned data
	assert.Equal(t, "2737bf16b39ab5d7b4a1", oauthClient.ID, "Unexpected ID")
	assert.Equal(t, "Test_Client_1", oauthClient.Label, "Unexpected Label")
	assert.False(t, oauthClient.Public, "Unexpected Public value")
	assert.Equal(t, "https://example.org/oauth/callback", oauthClient.RedirectURI, "Unexpected Redirect URI")
	assert.Equal(t, linodego.OAuthClientStatus("active"), oauthClient.Status, "Unexpected Status")
	assert.Equal(t, "https://api.linode.com/v4/account/clients/2737bf16b39ab5d7b4a1/thumbnail", *oauthClient.ThumbnailURL, "Unexpected Thumbnail URL")
}

func TestAccountOauthClient_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_oauth_client_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.OAuthClientUpdateOptions{
		Label:       "Test_Client_1_Updated",
		RedirectURI: "https://example_updated.org/oauth/callback",
		Public:      true,
	}

	clientID := "2737bf16b39ab5d7b4a1"
	base.MockPut(fmt.Sprintf("account/oauth-clients/%s", clientID), fixtureData)

	oauthClient, err := base.Client.UpdateOAuthClient(context.Background(), clientID, requestData)
	assert.NoError(t, err)
	// Assertions on the updated data
	assert.Equal(t, "2737bf16b39ab5d7b4a1", oauthClient.ID, "Unexpected ID")
	assert.Equal(t, "Test_Client_1_Updated", oauthClient.Label, "Unexpected Label")
	assert.True(t, oauthClient.Public, "Unexpected Public value")
	assert.Equal(t, "https://example_updated.org/oauth/callback", oauthClient.RedirectURI, "Unexpected Redirect URI")
	assert.Equal(t, linodego.OAuthClientStatus("active"), oauthClient.Status, "Unexpected Status")
	assert.Equal(t, "https://api.linode.com/v4/account/clients/2737bf16b39ab5d7b4a1/thumbnail", *oauthClient.ThumbnailURL, "Unexpected Thumbnail URL")
}

func TestAccountOauthClient_Delete(t *testing.T) {
	client := createMockClient(t)

	clientID := "2737bf16b39ab5d7b4a1"

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, fmt.Sprintf("account/oauth-clients/%s", clientID)),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteOAuthClient(context.Background(), clientID); err != nil {
		t.Fatal(err)
	}
}

func TestAccountOauthClient_Reset(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_oauth_client_reset")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	clientID := "2737bf16b39ab5d7b4a1"
	base.MockPost(fmt.Sprintf("account/oauth-clients/%s/reset-secret", clientID), fixtureData)

	oauthClient, err := base.Client.ResetOAuthClientSecret(context.Background(), clientID)
	assert.NoError(t, err)

	assert.Equal(t, "2737bf16b39ab5d7b4a1", oauthClient.ID, "Unexpected ID")
	assert.Equal(t, "Test_Client_1", oauthClient.Label, "Unexpected Label")
	assert.False(t, oauthClient.Public, "Unexpected Public value")
	assert.Equal(t, "https://example.org/oauth/callback", oauthClient.RedirectURI, "Unexpected Redirect URI")
	assert.Equal(t, linodego.OAuthClientStatus("active"), oauthClient.Status, "Unexpected Status")
	assert.Equal(t, "https://api.linode.com/v4/account/clients/2737bf16b39ab5d7b4a1/thumbnail", *oauthClient.ThumbnailURL, "Unexpected Thumbnail URL")
	assert.Equal(t, "<REDACTED>", oauthClient.Secret, "Secret should have been reset")
	assert.NotEmpty(t, oauthClient.Secret, "Secret should not be empty after reset")
}
