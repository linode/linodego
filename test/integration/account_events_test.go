package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/linode/linodego"
)

func TestListEvents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, instance, teardown, err := setupInstance(t, "fixtures/TestListEvents")
	defer teardown()
	if err != nil {
		t.Error(err)
	}
	configOpts := linodego.InstanceConfigCreateOptions{
		Label: "linodego-test-config",
	}
	_, err = client.CreateInstanceConfig(context.Background(), instance.ID, configOpts)
	if err != nil {
		t.Error(err)
	}

	filter := fmt.Sprintf("{\"entity.id\":%d, \"entity.type\": \"linode\"}", instance.ID)
	events, err := client.ListEvents(context.Background(), &linodego.ListOptions{Filter: filter})

	if err != nil {
		t.Errorf("Error getting Events, expected struct, got error %v", err)
	}

	if len(events) == 0 {
		t.Errorf("Expected to see at least one event")
	} else {
		assertDateSet(t, events[0].Created)
	}

}
