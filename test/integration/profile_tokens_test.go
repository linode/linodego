package integration

import (
	"context"
	"time"

	. "github.com/linode/linodego"
	"github.com/linode/linodego/pkg/errors"

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
	e, ok := err.(*errors.Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing token, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing token, got %v", e.Code)
	}
}

func TestGetToken_found(t *testing.T) {
	tokenTTLSeconds := 120
	ttl := time.Now().UTC().Add(time.Second * time.Duration(tokenTTLSeconds))
	client, token, teardown, err := setupProfileToken(t, "fixtures/TestGetToken_found", &ttl)
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

	assertDateSet(t, i.Created)

	updateOpts := i.GetUpdateOptions()
	if updateOpts.Label != i.Label {
		t.Errorf("Expected matching Label from GetUpdateOptions, got: %v", updateOpts)
	}

	createOpts := i.GetCreateOptions()
	if createOpts.Expiry == nil {
		t.Errorf("Expected non-nil Expiry from GetCreateOptions, got: %v", createOpts)
	}
}

func TestGetToken_noexpiry(t *testing.T) {
	client, token, teardown, err := setupProfileToken(t, "fixtures/TestGetToken_noexpiry", nil)
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

	createOpts := i.GetCreateOptions()
	if createOpts.Expiry != nil && createOpts.Expiry.Year() != 2999 {
		t.Errorf("Expected \"never\" expiring timestamp from GetCreateOptions, got: %v", createOpts)
	}
}
func TestUpdateTokens(t *testing.T) {
	tokenTTLSeconds := 120
	ttl := time.Now().UTC().Add(time.Second * time.Duration(tokenTTLSeconds))
	client, token, teardown, err := setupProfileToken(t, "fixtures/TestUpdateToken", &ttl)
	defer teardown()
	if err != nil {
		t.Errorf("Error creating test token: %s", err)
	}

	createOpts := token.GetCreateOptions()
	if createOpts.Expiry == nil {
		t.Errorf("Expected non-nil Expiry from GetCreateOptions, got: %v", createOpts)
	}

	updateOpts := token.GetUpdateOptions()
	if updateOpts.Label != token.Label {
		t.Errorf("Expected matching Label from GetUpdateOptions, got: %v", updateOpts)
	}

	updateOpts.Label = updateOpts.Label + "_renamed"

	i, err := client.UpdateToken(context.Background(), token.ID, updateOpts)
	if err != nil {
		t.Errorf("Error updating token: %s", err)
	}
	if i.Label != updateOpts.Label {
		t.Errorf("Expected token label to be changed, but found %v", i)
	}
}
func TestListTokens(t *testing.T) {
	tokenTTLSeconds := 120
	ttl := time.Now().UTC().Add(time.Second * time.Duration(tokenTTLSeconds))
	client, _, teardown, err := setupProfileToken(t, "fixtures/TestListTokens", &ttl)
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

func setupProfileToken(t *testing.T, fixturesYaml string, ttl *time.Time) (*Client, *Token, func(), error) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	// This scope must be <= the scope used for testing
	limitedTestScope := "linodes:read_only"

	createOpts := TokenCreateOptions{
		Label:  "linodego-test-token",
		Expiry: ttl,
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
