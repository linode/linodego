package linodego_test

import (
	"context"
	"time"

	. "github.com/linode/linodego"

	"testing"
)

func TestGetToken_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetToken_missing")
	defer teardown()

	var doesNotExist = 123
	i, err := client.GetToken(context.Background(), doesNotExist)
	if err == nil {
		t.Errorf("should have received an error requesting a missing token, got %v", i)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing token, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing token, got %v", e.Code)
	}
}

func TestGetToken_found(t *testing.T) {
	client, token, teardown, err := setupProfileToken(t, "fixtures/TestGetToken_found")
	defer teardown()
	if err != nil {
		t.Errorf("Error creating test token: %s", err)
	}

	i, err := client.GetToken(context.Background(), token.ID)
	if err != nil {
		t.Errorf("Error getting token, expected struct, got %v and error %v", i, err)
	}
	if i.ID != token.ID {
		t.Errorf("Expected a specific token, but got a different one %v", i)
	}
}
func TestListTokens(t *testing.T) {
	client, _, teardown, err := setupProfileToken(t, "fixtures/TestListTokens")
	defer teardown()
	if err != nil {
		t.Errorf("Error creating test token: %s", err)
	}

	i, err := client.ListTokens(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing tokens, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of tokens, but got none %v", i)
	}
}

func setupProfileToken(t *testing.T, fixturesYaml string) (*Client, *Token, func(), error) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	tokenTTLSeconds := 120
	// This scope must be <= the scope used for testing
	limitedTestScope := "linodes:read_only"

	ttl := time.Now().UTC().Add(time.Second * time.Duration(tokenTTLSeconds))

	createOpts := TokenCreateOptions{
		Label:  "linodego-test-token",
		Expiry: &ttl,
		Scopes: limitedTestScope,
	}
	token, err := client.CreateToken(context.Background(), createOpts)
	if err != nil {
		t.Fatalf("Error creating test Token: %s", err)
	}

	teardown := func() {
		if token != nil {
			if err := client.DeleteToken(context.Background(), token.ID); err != nil {
				t.Errorf("Error deleting test Token: %s", err)
			}
		}
		fixtureTeardown()
	}
	return client, token, teardown, err
}
