package integration

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocks(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestLocks", false)
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(teardown)

	createOpts := linodego.LockCreateOptions{
		EntityType: linodego.EntityLinode,
		EntityID:   instance.ID,
		LockType:   linodego.LockTypeCannotDelete,
	}

	createdLock, err := client.CreateLock(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error creating lock: %s", err)
	}

	t.Cleanup(func() {
		if err := client.DeleteLock(context.Background(), createdLock.ID); err != nil {
			t.Errorf("Error deleting lock: %s", err)
		}
	})

	// Test Get
	lock, err := client.GetLock(context.Background(), createdLock.ID)
	if err != nil {
		t.Errorf("Error getting lock: %s", err)
	}

	if !cmp.Equal(lock, createdLock) {
		t.Errorf("Expected lock to match created lock but got diffs: %s", cmp.Diff(lock, createdLock))
	}

	// Test List
	locks, err := client.ListLocks(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing locks: %s", err)
	}

	if len(locks) == 0 {
		t.Error("Expected at least one lock")
	}

	// Since this is a dedicated test instance, we can assume the first lock is ours
	if locks[0].ID != createdLock.ID {
		t.Errorf("Expected first lock ID to be %d, got %d", createdLock.ID, locks[0].ID)
	}

	if !cmp.Equal(&locks[0], createdLock) {
		t.Errorf("Expected lock to match created lock but got diffs: %s", cmp.Diff(&locks[0], createdLock))
	}

	// Test Instance locks field
	refreshedInstance, err := client.GetInstance(context.Background(), instance.ID)
	if err != nil {
		t.Errorf("Error getting instance: %s", err)
	}

	if len(refreshedInstance.Locks) == 0 {
		t.Error("Expected instance to have locks")
	}

	if refreshedInstance.Locks[0] != linodego.LockTypeCannotDelete {
		t.Errorf("Expected instance to have %s lock, got %s", linodego.LockTypeCannotDelete, refreshedInstance.Locks[0])
	}
}

func TestTryToLockTwoResourcesWithTheSameType(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t,
		"fixtures/TestTryToLockTwoResourcesWithTheSameType", false)
	require.NoError(t, err)
	t.Cleanup(teardown)
	createOpts := linodego.LockCreateOptions{
		EntityType: linodego.EntityLinode,
		EntityID:   instance.ID,
		LockType:   linodego.LockTypeCannotDelete,
	}

	createdLock, err := client.CreateLock(context.Background(), createOpts)
	require.NoError(t, err)
	t.Cleanup(func() {
		client.DeleteLock(context.Background(), createdLock.ID)
	})

	createOpts.LockType = linodego.LockTypeCannotDeleteWithSubresources
	_, errConflictingLock := client.CreateLock(context.Background(), createOpts)
	require.Error(t, errConflictingLock)
}

func TestTryToCreateWithInvalidData(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestTryToCreateWithInvalidData")
	t.Cleanup(teardown)

	createOpts := linodego.LockCreateOptions{
		EntityType: linodego.EntityLinode,
		EntityID:   -99999876,
		LockType:   linodego.LockTypeCannotDeleteWithSubresources,
	}

	_, createLockErr := client.CreateLock(context.Background(), createOpts)
	require.Error(t, createLockErr)
	assert.Equal(t, "[400] [entity_id] entity_id is not valid", createLockErr.Error())
}
