package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestSupportTicket_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("support_ticket_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("support/tickets", fixtureData)

	tickets, err := base.Client.ListTickets(context.Background(), nil)
	if err != nil {
		t.Fatalf("Error listing tickets: %v", err)
	}

	assert.Equal(t, 1, len(tickets))
	ticket := tickets[0]

	assert.Equal(t, []string{"screenshot.jpg", "screenshot.txt"}, ticket.Attachments)
	assert.Equal(t, false, ticket.Closeable)
	assert.Equal(t, "I am having trouble setting the root password on my Linode. I tried following the instructions but something is not working. Can you please help me figure out how I can reset it?", ticket.Description)
	assert.Equal(t, 10400, ticket.Entity.ID)
	assert.Equal(t, "linode123456", ticket.Entity.Label)
	assert.Equal(t, "linode", ticket.Entity.Type)
	assert.Equal(t, "/v4/linode/instances/123456", ticket.Entity.URL)
	assert.Equal(t, "474a1b7373ae0be4132649e69c36ce30", ticket.GravatarID)
	assert.Equal(t, 11223344, ticket.ID)
	assert.Equal(t, "some_user", ticket.OpenedBy)
	assert.Equal(t, linodego.TicketStatus("open"), ticket.Status)
	assert.Equal(t, "Having trouble resetting root password on my Linode", ticket.Summary)
	assert.Equal(t, "some_other_user", ticket.UpdatedBy)
}

func TestSupportTicket_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("support_ticket_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("support/tickets/11223344", fixtureData)

	ticket, err := base.Client.GetTicket(context.Background(), 11223344)
	if err != nil {
		t.Fatalf("Error getting ticket: %v", err)
	}

	assert.Equal(t, []string{"screenshot.jpg", "screenshot.txt"}, ticket.Attachments)
	assert.Equal(t, false, ticket.Closeable)
	assert.Equal(t, "I am having trouble setting the root password on my Linode. I tried following the instructions but something is not working. Can you please help me figure out how I can reset it?", ticket.Description)
	assert.Equal(t, 10400, ticket.Entity.ID)
	assert.Equal(t, "linode123456", ticket.Entity.Label)
	assert.Equal(t, "linode", ticket.Entity.Type)
	assert.Equal(t, "/v4/linode/instances/123456", ticket.Entity.URL)
	assert.Equal(t, "474a1b7373ae0be4132649e69c36ce30", ticket.GravatarID)
	assert.Equal(t, 11223344, ticket.ID)
	assert.Equal(t, "some_user", ticket.OpenedBy)
	assert.Equal(t, linodego.TicketStatus("open"), ticket.Status)
	assert.Equal(t, "Having trouble resetting root password on my Linode", ticket.Summary)
	assert.Equal(t, "some_other_user", ticket.UpdatedBy)
}
