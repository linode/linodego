package integration

import (
	"context"
	"github.com/linode/linodego"
	"testing"
)

func TestEventPoller_InstancePower(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestEventPoller_InstancePower")
	t.Cleanup(fixtureTeardown)

	p, err := client.NewEventPollerWithoutEntity(linodego.EntityLinode, linodego.ActionLinodeCreate)
	if err != nil {
		t.Fatalf("failed to initialize event poller: %s", err)
	}

	booted := false

	instance, err := client.CreateInstance(context.Background(), linodego.InstanceCreateOptions{
		Region:   "us-southeast",
		Type:     "g6-nanode-1",
		Image:    "linode/ubuntu22.04",
		RootPass: "c00lp@ss!",
		Label:    "linodego-test-eventpoller-instancepower",
		Booted:   &booted,
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
