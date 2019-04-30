package linodego_test

import (
	"context"

	. "github.com/linode/linodego"

	"testing"
)

func TestGetOAuthClient_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetOAuthClient_missing")
	defer teardown()

	i, err := client.GetOAuthClient(context.Background(), 0)
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

func TestGetOAuthClient_found(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetOAuthClient_found")
	defer teardown()

	i, err := client.GetOAuthClient(context.Background(), "linode/ubuntu16.04lts")
	if err != nil {
		t.Errorf("Error getting oauthClient, expected struct, got %v and error %v", i, err)
	}
	if i.ID != "linode/ubuntu16.04lts" {
		t.Errorf("Expected a specific oauthClient, but got a different one %v", i)
	}
}
func TestListOAuthClients(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestListOAuthClients")
	defer teardown()

	i, err := client.ListOAuthClients(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing oauthClients, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of oauthClients, but got none %v", i)
	}
}
