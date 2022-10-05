package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestAccountEvents_List(t *testing.T) {
	client, instance, teardown, err := setupInstance(t, "fixtures/TestAccountEvents_List")
	defer teardown()
	if err != nil {
		t.Error(err)
	}
	configOpts := linodego.InstanceConfigCreateOptions{
		Label: "test-config",
	}
	instanceConfig, err := client.CreateInstanceConfig(context.Background(), instance.ID, configOpts)
	if err != nil {
		t.Error(err)
	}

	f := linodego.Filter{}
	f.AddField(linodego.Eq, "entity.id", instance.ID)
	f.AddField(linodego.Eq, "entity.type", "linode")
	f.AddField(linodego.Eq, "action", "linode_config_create")
	filter, err := f.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to marshal filter: %v", err)
	}
	events, err := client.ListEvents(context.Background(), &linodego.ListOptions{Filter: string(filter)})
	if err != nil {
		t.Errorf("Error getting Events, expected struct, got error %v", err)
	}

	if len(events) == 0 {
		t.Errorf("Expected to see at least one event")
	} else {
		event := events[0]
		assertDateSet(t, event.Created)
		if event.SecondaryEntity == nil {
			t.Errorf("Expected Secondary Entity to be set")
		} else if event.SecondaryEntity.Label != instanceConfig.Label {
			t.Errorf("Expected Secondary Entity label to be '%s', got '%s'", instanceConfig.Label, event.SecondaryEntity.Label)
		}
	}
}
