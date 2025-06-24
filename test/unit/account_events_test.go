package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestAccountEvents_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_events_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/events", fixtureData)

	events, err := base.Client.ListEvents(context.Background(), nil)
	if err != nil {
		t.Fatalf("Error listing events: %v", err)
	}

	assert.Equal(t, 1, len(events))
	event := events[0]
	assert.Equal(t, linodego.EventAction("ticket_create"), event.Action)
	assert.Equal(t, 300.56, event.Duration)
	assert.Equal(t, float64(11111), event.Entity.ID)
	assert.Equal(t, "Problem booting my Linode", event.Entity.Label)
	assert.Equal(t, linodego.EntityType("ticket"), event.Entity.Type)
	assert.Equal(t, "/v4/support/tickets/11111", event.Entity.URL)
	assert.Equal(t, 123, event.ID)
	assert.Equal(t, "None", event.Message)
	assert.Equal(t, true, event.Read)
	assert.Equal(t, true, event.Seen)
	assert.Equal(t, linodego.EventStatus("failed"), event.Status)
	assert.Equal(t, "exampleUser", event.Username)
	assert.Equal(t, "linode/debian9", event.SecondaryEntity.ID)
	assert.Equal(t, "linode1234", event.SecondaryEntity.Label)
	assert.Equal(t, linodego.EntityType("linode"), event.SecondaryEntity.Type)
	assert.Equal(t, "/v4/linode/instances/1234", event.SecondaryEntity.URL)
}

func TestAccountEvents_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_events_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/events/11111", fixtureData)

	event, err := base.Client.GetEvent(context.Background(), 11111)
	if err != nil {
		t.Fatalf("Error getting event: %v", err)
	}

	assert.Equal(t, linodego.EventAction("ticket_create"), event.Action)
	assert.Equal(t, 300.56, event.Duration)
	assert.Equal(t, float64(11111), event.Entity.ID)
	assert.Equal(t, "Problem booting my Linode", event.Entity.Label)
	assert.Equal(t, linodego.EntityType("ticket"), event.Entity.Type)
	assert.Equal(t, "/v4/support/tickets/11111", event.Entity.URL)
	assert.Equal(t, 123, event.ID)
	assert.Equal(t, "None", event.Message)
	assert.Equal(t, true, event.Read)
	assert.Equal(t, true, event.Seen)
	assert.Equal(t, linodego.EventStatus("failed"), event.Status)
	assert.Equal(t, "exampleUser", event.Username)
	assert.Equal(t, "linode/debian9", event.SecondaryEntity.ID)
	assert.Equal(t, "linode1234", event.SecondaryEntity.Label)
	assert.Equal(t, linodego.EntityType("linode"), event.SecondaryEntity.Type)
	assert.Equal(t, "/v4/linode/instances/1234", event.SecondaryEntity.URL)
}
