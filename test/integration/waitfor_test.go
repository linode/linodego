package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/linode/linodego"
)

func TestEventPoller_InstancePower(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestEventPoller_InstancePower")
	t.Cleanup(fixtureTeardown)

	p, err := client.NewEventPollerWithoutEntity(linodego.EntityLinode, linodego.ActionLinodeCreate)
	if err != nil {
		t.Fatalf("failed to initialize event poller: %s", err)
	}

	instance, err := client.CreateInstance(context.Background(), linodego.InstanceCreateOptions{
		Region:   getRegionsWithCaps(t, client, []string{"Linodes"})[0],
		Type:     "g6-nanode-1",
		Image:    linodego.Pointer("linode/ubuntu22.04"),
		RootPass: linodego.Pointer(randPassword()),
		Label:    linodego.Pointer("go-ins-poll-test"),
		Booted:   linodego.Pointer(false),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := client.DeleteInstance(context.Background(), instance.ID); err != nil {
			t.Error(err)
		}
	})

	p.EntityID = instance.ID

	if _, err := p.WaitForFinished(context.Background(), 120); err != nil {
		t.Fatal(err)
	}

	// Wait for the instance to be booted
	p, err = client.NewEventPoller(
		context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionLinodeBoot)
	if err != nil {
		t.Fatal(err)
	}

	if err := client.BootInstance(context.Background(), instance.ID, 0); err != nil {
		t.Fatal(err)
	}

	event, err := p.WaitForFinished(context.Background(), 200)
	if err != nil {
		t.Fatal(err)
	}

	inst, err := client.GetInstance(context.Background(), instance.ID)
	if err != nil {
		t.Fatal(err)
	}

	if inst.Status != linodego.InstanceRunning {
		t.Fatalf("instance expected to be running, got %s", inst.Status)
	}

	// Wait for the instance to be shut down
	p, err = client.NewEventPoller(
		context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionLinodeShutdown)
	if err != nil {
		t.Fatal(err)
	}

	if err := client.ShutdownInstance(context.Background(), instance.ID); err != nil {
		t.Fatal(err)
	}

	event, err = p.WaitForFinished(context.Background(), 200)
	if err != nil {
		t.Fatal(err)
	}

	if event.Action != linodego.ActionLinodeShutdown {
		t.Fatalf("action type mismatch: %s != %s", event.Action, linodego.ActionLinodeShutdown)
	}

	inst, err = client.GetInstance(context.Background(), instance.ID)
	if err != nil {
		t.Fatal(err)
	}

	if inst.Status != linodego.InstanceOffline {
		t.Fatalf("instance expected to be offline, got %s", inst.Status)
	}
}

func TestWaitForResourceFree(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestWaitForResourceFree")
	t.Cleanup(fixtureTeardown)

	// Create a booted instance
	instance, err := client.CreateInstance(context.Background(), linodego.InstanceCreateOptions{
		Region:   getRegionsWithCaps(t, client, []string{"Linodes"})[0],
		Type:     "g6-nanode-1",
		Image:    linodego.Pointer("linode/ubuntu22.04"),
		RootPass: linodego.Pointer(randPassword()),
		Label:    linodego.Pointer("linodego-waitforfree"),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := client.DeleteInstance(context.Background(), instance.ID); err != nil {
			t.Error(err)
		}
	})

	// Wait for the instance to be done with all events
	err = client.WaitForResourceFree(context.Background(), linodego.EntityLinode, instance.ID, 240)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the instance is no longer busy
	instance, err = client.GetInstance(context.Background(), instance.ID)
	if err != nil {
		t.Fatal(err)
	}

	if instance.Status == linodego.InstanceProvisioning {
		t.Fatalf("expected instance to not be provisioning, got %s", instance.Status)
	}
}

func TestEventPoller_Secondary(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestEventPoller_Secondary")
	defer fixtureTeardown()

	createPoller, err := client.NewEventPollerWithoutEntity(linodego.EntityLinode, linodego.ActionLinodeCreate)
	if err != nil {
		t.Fatalf("failed to initialize event poller: %s", err)
	}

	instance, err := client.CreateInstance(context.Background(), linodego.InstanceCreateOptions{
		Region: getRegionsWithCaps(t, client, []string{"Linodes"})[0],
		Type:   "g6-nanode-1",
		Label:  linodego.Pointer("go-ins-poll-test"),
		Booted: linodego.Pointer(false),
	})
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := client.DeleteInstance(context.Background(), instance.ID); err != nil {
			t.Error(err)
		}
	}()

	createPoller.EntityID = instance.ID

	if _, err := createPoller.WaitForFinished(context.Background(), 120); err != nil {
		t.Fatal(err)
	}

	// Create two instance disks
	disks := make([]*linodego.InstanceDisk, 2)

	for i := 0; i < 2; i++ {
		disk, err := client.CreateInstanceDisk(context.Background(), instance.ID, linodego.InstanceDiskCreateOptions{
			Label: fmt.Sprintf("test-disk-%d", i),
			Size:  512,
		})
		if err != nil {
			t.Fatalf("failed to create instance disk: %s", err)
		}

		disks[i] = disk
	}

	// Poll for the first disk to be deleted
	diskPoller, err := client.NewEventPollerWithSecondary(
		context.Background(),
		instance.ID,
		linodego.EntityLinode,
		disks[0].ID,
		linodego.ActionDiskDelete)
	if err != nil {
		t.Fatal(err)
	}

	// Poll for the second disk to be deleted
	diskPoller2, err := client.NewEventPollerWithSecondary(
		context.Background(),
		instance.ID,
		linodego.EntityLinode,
		disks[1].ID,
		linodego.ActionDiskDelete)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the first two disks
	for i := 0; i < 2; i++ {
		if err := client.DeleteInstanceDisk(context.Background(), instance.ID, disks[i].ID); err != nil {
			t.Fatal(err)
		}
	}

	// Wait for the first disk to be deleted
	deleteEvent, err := diskPoller.WaitForFinished(context.Background(), 60)
	if err != nil {
		t.Fatal(err)
	}

	// Poll for the second disk to be deleted.
	// This is necessary to cover an edge case that triggers a panic
	// when other deletion events with a null SecondaryEntity are discovered.
	if _, err := diskPoller2.WaitForFinished(context.Background(), 60); err != nil {
		t.Fatal(err)
	}

	entityID := deleteEvent.SecondaryEntity.ID

	// Sometimes the JSON unmarshaler will
	// parse IDs as floats rather than ints.
	if value, ok := entityID.(float64); ok {
		entityID = int(value)
	}

	if entityID != disks[0].ID {
		t.Fatalf("expected event and first deleteEvent id to match; got %v", deleteEvent.SecondaryEntity.ID)
	}
}
