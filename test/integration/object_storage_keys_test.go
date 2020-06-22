package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/linode/linodego"
	"github.com/linode/linodego/pkg/errors"
)

var (
	testObjectStorageKeyCreateOpts = linodego.ObjectStorageKeyCreateOptions{
		Label: label,
	}
)

func TestGetObjectStorageKey_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetObjectStorageKey_missing")
	defer teardown()

	notfoundID := 123
	i, err := client.GetObjectStorageKey(context.Background(), notfoundID)
	if err == nil {
		t.Errorf("should have received an error requesting a missing object storage key, got %v", i)
	}
	e, ok := err.(*errors.Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing object storage key, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing object storage key, got %v", e.Code)
	}
}

func TestGetObjectStorageKey_found(t *testing.T) {
	client, objectStorageKey, teardown, err := setupObjectStorageKey(t, "fixtures/TestGetObjectStorageKey_found")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	i, err := client.GetObjectStorageKey(context.Background(), objectStorageKey.ID)
	if err != nil {
		t.Errorf("Error getting objectStorageKey, expected struct, got %v and error %v", i, err)
	}
	if i.ID != objectStorageKey.ID {
		t.Errorf("Expected objectStorageKey id %d, but got %d", i.ID, objectStorageKey.ID)
	}
	if testObjectStorageKeyCreateOpts.Label != objectStorageKey.Label {
		t.Errorf("Expected objectStorageKey label '%s', but got '%s'", testObjectStorageKeyCreateOpts.Label, objectStorageKey.Label)
	}
}

func TestUpdateObjectStorageKey(t *testing.T) {
	client, objectStorageKey, teardown, err := setupObjectStorageKey(t, "fixtures/TestUpdateObjectStorageKey")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	renamedLabel := objectStorageKey.Label + "_r"
	updateOpts := linodego.ObjectStorageKeyUpdateOptions{
		Label: renamedLabel,
	}
	objectStorageKey, err = client.UpdateObjectStorageKey(context.Background(), objectStorageKey.ID, updateOpts)

	if err != nil {
		t.Errorf("Error renaming objectStorageKey, %s", err)
	}

	if !strings.Contains(objectStorageKey.Label, "-linodego-testing_r") {
		t.Errorf("objectStorageKey returned does not match objectStorageKey update request")
	}
}

func TestListObjectStorageKeys(t *testing.T) {
	client, objkey, teardown, err := setupObjectStorageKey(t, "fixtures/TestListObjectStorageKey")
	defer teardown()
	if err != nil {
		t.Error(err)
	}
	objectStorageKeys, err := client.ListObjectStorageKeys(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing objectStorageKeys, expected struct, got error %v", err)
	}
	if len(objectStorageKeys) == 0 {
		t.Errorf("Expected a list of objectStorageKeys, but got %v", objectStorageKeys)
	}

	notFound := true
	for i := range objectStorageKeys {
		if objectStorageKeys[i].Label == objkey.Label {
			notFound = false
			break
		}
	}
	if notFound {
		t.Errorf("Expected to find created objectStorageKey, but '%s' was not found", objkey.Label)
	}
}

func setupObjectStorageKey(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.ObjectStorageKey, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	createOpts := testObjectStorageKeyCreateOpts
	objectStorageKey, err := client.CreateObjectStorageKey(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error listing objectStorageKey, expected struct, got error %v", err)
	}

	teardown := func() {
		if err := client.DeleteObjectStorageKey(context.Background(), objectStorageKey.ID); err != nil {
			t.Errorf("Expected to delete a objectStorageKey, but got %v", err)
		}
		fixtureTeardown()
	}
	return client, objectStorageKey, teardown, err
}
