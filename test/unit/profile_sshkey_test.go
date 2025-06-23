package unit

import (
	"context"
	"testing"
	"time"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestProfileSSHKey_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_sshkey_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("profile/sshkeys/123", fixtureData)

	key, err := base.Client.GetSSHKey(context.Background(), 123)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	expectedTime, _ := time.Parse(time.RFC3339, "2024-03-10T12:00:00Z")

	if assert.NotNil(t, key.Created) {
		assert.Equal(t, expectedTime, *key.Created)
	}
}

func TestProfileSSHKeys_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_sshkeys_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("profile/sshkeys", fixtureData)

	keys, err := base.Client.ListSSHKeys(context.Background(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, keys)
	assert.Len(t, keys, 2)

	expectedTimes := []string{
		"2024-03-10T12:00:00Z",
		"2024-03-11T15:45:00Z",
	}

	for i, key := range keys {
		if assert.NotNil(t, key.Created) {
			expectedTime, _ := time.Parse(time.RFC3339, expectedTimes[i])
			assert.Equal(t, expectedTime, *key.Created)
		}
	}
}

func TestProfileSSHKey_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_sshkey_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("profile/sshkeys", fixtureData)

	opts := linodego.SSHKeyCreateOptions{
		Label:  "Test Key",
		SSHKey: "ssh-rsa AAAAB3...",
	}

	key, err := base.Client.CreateSSHKey(context.Background(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.Equal(t, opts.Label, key.Label)
	assert.Equal(t, opts.SSHKey, key.SSHKey)
}

func TestProfileSSHKey_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_sshkey_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("profile/sshkeys/123", fixtureData)

	opts := linodego.SSHKeyUpdateOptions{
		Label: "Updated Key",
	}

	key, err := base.Client.UpdateSSHKey(context.Background(), 123, opts)
	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.Equal(t, opts.Label, key.Label)
}

func TestProfileSSHKey_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("profile/sshkeys/123", nil)

	err := base.Client.DeleteSSHKey(context.Background(), 123)
	assert.NoError(t, err)
}
