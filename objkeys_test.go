package linodego_test

import (
	"context"
	. "github.com/linode/linodego"
	"strings"
	"testing"
)

var (
	testObjKeyCreateOpts = ObjKeyCreateOptions{
		Label: label,
	}
)

func TestGetObjKey_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetObjKey_missing")
	defer teardown()

	notfoundID := 123
	i, err := client.GetObjKey(context.Background(), notfoundID)
	if err == nil {
		t.Errorf("should have received an error requesting a missing objkey, got %v", i)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing objkey, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing objkey, got %v", e.Code)
	}
}

func TestGetObjKey_found(t *testing.T) {
	client, objkey, teardown, err := setupObjKey(t, "fixtures/TestGetObjKey_found")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	i, err := client.GetObjKey(context.Background(), objkey.ID)
	if err != nil {
		t.Errorf("Error getting objkey, expected struct, got %v and error %v", i, err)
	}
	if i.ID != objkey.ID {
		t.Errorf("Expected objkey id %d, but got %d", i.ID, objkey.ID)
	}
	if testObjKeyCreateOpts.Label != objkey.Label {
		t.Errorf("Expected objkey label '%s', but got '%s'", testObjKeyCreateOpts.Label, objkey.Label)
	}
}

func TestUpdateObjKey(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, objkey, teardown, err := setupObjKey(t, "fixtures/TestUpdateObjKey")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	renamedLabel := objkey.Label + "_r"
	updateOpts := ObjKeyUpdateOptions{
		Label: renamedLabel,
	}
	objkey, err = client.UpdateObjKey(context.Background(), objkey.ID, updateOpts)

	if err != nil {
		t.Errorf("Error renaming objkey, %s", err)
	}

	if !strings.Contains(objkey.Label, "-linodego-testing_r") {
		t.Errorf("objkey returned does not match objkey update request")
	}
}

func TestListObjKeys(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, objkey, teardown, err := setupObjKey(t, "fixtures/TestListObjKey")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	objkeys, err := client.ListObjKeys(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing objkeys, expected struct, got error %v", err)
	}
	if len(objkeys) == 0 {
		t.Errorf("Expected a list of objkeys, but got %v", objkeys)
	}
	if !strings.Contains(objkey.Label, "-linodego-testing") {
		t.Errorf("objkey returned does not match objkey update request")
	}
}

func setupObjKey(t *testing.T, fixturesYaml string) (*Client, *ObjKey, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	createOpts := testObjKeyCreateOpts
	objkey, err := client.CreateObjKey(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error listing objkeys, expected struct, got error %v", err)
	}

	teardown := func() {
		if err := client.DeleteObjKey(context.Background(), objkey.ID); err != nil {
			t.Errorf("Expected to delete a objkey, but got %v", err)
		}
		fixtureTeardown()
	}
	return client, objkey, teardown, err
}
