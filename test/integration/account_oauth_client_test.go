package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
	. "github.com/linode/linodego"
)

func TestOAuthClient_GetMissing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestOAuthClient_GetMissing")
	defer teardown()

	i, err := client.GetOAuthClient(context.Background(), "does-not-exist")
	if err == nil {
		t.Errorf("should have received an error requesting a missing oauthClient, got %v", i)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing oauthClient, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing oauthClient, got %v", e.Code)
	}
}

func TestOAuthClient_GetFound(t *testing.T) {
	createOpts := linodego.OAuthClientCreateOptions{
		Public:      true,
		RedirectURI: "https://example.com",
		Label:       "go-client-test",
	}

	client, oauthClient, teardown, err := setupOAuthClient(t, createOpts, "fixtures/TestOAuthClient_GetFound")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	i, err := client.GetOAuthClient(context.Background(), oauthClient.ID)
	if err != nil {
		t.Errorf("Error getting oauthClient, expected struct, got %v and error %v", i, err)
	}
	if i.Label != createOpts.Label {
		t.Errorf("Expected a specific oauthClient, but got a different one %v", i)
	}
}

func TestOAuthClients_List(t *testing.T) {
	createOpts := linodego.OAuthClientCreateOptions{
		Public:      true,
		RedirectURI: "https://example.com",
		Label:       "go-client-test",
	}
	client, _, teardown, err := setupOAuthClient(t, createOpts, "fixtures/TestOAuthClients_List")
	defer teardown()
	if err != nil {
		t.Error(err)
	}
	i, err := client.ListOAuthClients(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing oauthClients, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of oauthClients, but got none %v", i)
	}
}

func TestOAuthClients_Reset(t *testing.T) {
	createOpts := linodego.OAuthClientCreateOptions{
		Public:      true,
		RedirectURI: "https://example.com",
		Label:       "go-client-test",
	}
	client, oauthClient, teardown, err := setupOAuthClient(t, createOpts, "fixtures/TestOAuthClients_Reset")
	defer teardown()
	if err != nil {
		t.Error(err)
	}
	oauthClientAfterReset, err := client.ResetOAuthClientSecret(context.Background(), oauthClient.ID)
	if err != nil {
		t.Errorf("Error resetting oauthClient secret, expected struct, got error %v", err)
	}

	assert.NotEqual(t, oauthClient.Secret, oauthClientAfterReset.Secret, "Secret should have been reset")
}

func setupOAuthClient(t *testing.T, createOpts linodego.OAuthClientCreateOptions, fixturesYaml string) (*linodego.Client, *linodego.OAuthClient, func(), error) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	oauthClient, err := client.CreateOAuthClient(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error creating test OAuthClient: %s", err)
	}

	teardown := func() {
		if err := client.DeleteOAuthClient(context.Background(), oauthClient.ID); err != nil {
			t.Errorf("Error deleting test OAuthClient: %s", err)
		}
		fixtureTeardown()
	}
	return client, oauthClient, teardown, err
}
