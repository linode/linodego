package integration

import (
	"context"
	"testing"
)

func TestProfile_Get_smoke(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestProfile_Get")
	defer teardown()

	i, err := client.GetProfile(context.Background())
	if err != nil {
		t.Errorf("Error getting profile: %s", err)
	}
	if len(i.Email) < 1 {
		t.Errorf("Expected profile email")
	}
}

func TestProfile_Update(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestProfile_Update")
	defer teardown()

	profile, err := client.GetProfile(context.Background())
	if err != nil {
		t.Errorf("Error getting profile: %s", err)
	}

	updateOpts := profile.GetUpdateOptions()
	if updateOpts.Email != profile.Email {
		t.Errorf("Expected matching Username from GetUpdateOptions, got: %v", updateOpts)
	}

	i, err := client.UpdateProfile(context.Background(), updateOpts)
	if err != nil {
		t.Errorf("Error updating profile: %s", err)
	}
	if i.Email != updateOpts.Email {
		t.Errorf("Expected profile email to be changed, but found %v", i)
	}
}
