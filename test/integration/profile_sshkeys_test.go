package integration

import (
	"context"
	"strings"
	"testing"

	. "github.com/linode/linodego"
)

var testSSHKeyCreateOpts = SSHKeyCreateOptions{
	Label:  label,
	SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDYlv4Ns3tY2NEseuuMXEz1sLzO9sGC0cwaT2ECbWFyrsn1Fg5ISdkaJD8LiuhZ41/1Mh0Sq49wY89yLkmw+Ukrd+thFbhUqTzjL09U89kn3Ds/ajVJgwnJ4pXmBqhq0/3pmO/UkYIBi5ErTnPWL+yHAoQ1HsVetxYUmY2SPaT0pduDIrvNZRvWn3Nvn9qsUVfthWiGc8oHWE5xyd7+3UPLHSMkE4rZd2k6e7bJWCM/VJ7ZrJQ6UVTDXjBCkkT12WsOWxcEuL36RUGgGa4h5M4IY0SkgQSKHer01dJSj3c6OBzj2CRDZFoM8f/YC66s0+ZQ9cE/aADDycMIvqOJBI6X " + label,
}

func TestSSHKey_GetMissing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestSSHKey_GetMissing")
	defer teardown()

	notfoundID := 123
	i, err := client.GetSSHKey(context.Background(), notfoundID)
	if err == nil {
		t.Errorf("should have received an error requesting a missing sshkey, got %v", i)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing sshkey, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing sshkey, got %v", e.Code)
	}
}

func TestSSHKey_GetFound(t *testing.T) {
	client, sshkey, teardown, err := setupSSHKey(t, "fixtures/TestSSHKey_GetFound")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	i, err := client.GetSSHKey(context.Background(), sshkey.ID)
	if err != nil {
		t.Errorf("Error getting sshkey, expected struct, got %v and error %v", i, err)
	}
	if i.ID != sshkey.ID {
		t.Errorf("Expected sshkey id %d, but got %d", i.ID, sshkey.ID)
	}
	if testSSHKeyCreateOpts.Label != sshkey.Label {
		t.Errorf("Expected sshkey label '%s', but got '%s'", testSSHKeyCreateOpts.Label, sshkey.Label)
	}
	if testSSHKeyCreateOpts.SSHKey != sshkey.SSHKey {
		t.Errorf("Expected sshkey sshkey, but got a different one")
	}

	assertDateSet(t, sshkey.Created)
}

func TestSSHKey_Update(t *testing.T) {
	client, sshkey, teardown, err := setupSSHKey(t, "fixtures/TestSSHKey_Update")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	renamedLabel := sshkey.Label + "_r"
	updateOpts := SSHKeyUpdateOptions{
		Label: renamedLabel,
	}
	sshkey, err = client.UpdateSSHKey(context.Background(), sshkey.ID, updateOpts)

	if err != nil {
		t.Errorf("Error renaming sshkey, %s", err)
	}

	if !strings.Contains(sshkey.Label, "-linodego-testing_r") {
		t.Errorf("sshkey returned does not match sshkey update request")
	}
}

func TestSSHKeys_List(t *testing.T) {
	client, sshkey, teardown, err := setupSSHKey(t, "fixtures/TestSSHKeys_List")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	sshkeys, err := client.ListSSHKeys(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing sshkeys, expected struct, got error %v", err)
	}
	if len(sshkeys) == 0 {
		t.Errorf("Expected a list of sshkeys, but got %v", sshkeys)
	}
	notFound := true
	for i := range sshkeys {
		if sshkeys[i].Label == sshkey.Label {
			notFound = false

			if sshkeys[i].Created == nil {
				t.Errorf("Expected listed sshkeys to have parsed Created")
			}
			assertDateSet(t, sshkeys[i].Created)
			break
		}
	}
	if notFound {
		t.Errorf("Expected to find created sshkey, but '%s' was not found", sshkey.Label)
	}
}

func setupSSHKey(t *testing.T, fixturesYaml string) (*Client, *SSHKey, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	createOpts := testSSHKeyCreateOpts
	sshkey, err := client.CreateSSHKey(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error listing sshkeys, expected struct, got error %v", err)
	}

	teardown := func() {
		if err := client.DeleteSSHKey(context.Background(), sshkey.ID); err != nil {
			t.Errorf("Expected to delete a sshkey, but got %v", err)
		}
		fixtureTeardown()
	}
	return client, sshkey, teardown, err
}
