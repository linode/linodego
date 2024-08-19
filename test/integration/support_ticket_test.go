package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTicket_List(t *testing.T) {
	warnSensitiveTest(t)
	client, teardown := createTestClient(t, "fixtures/TestTicket_List")
	defer teardown()

	tickets, err := client.ListTickets(context.Background(), nil)
	require.NoError(t, err, "Error getting Tickets, expected struct")

	require.NotEmpty(t, tickets, "Expected to see tickets returned")

	require.Equal(t, 123, tickets[0].ID, "Expected ticket ID 123")
}

func TestTicket_Get(t *testing.T) {
	warnSensitiveTest(t)
	client, teardown := createTestClient(t, "fixtures/TestTicket_Get")
	defer teardown()

	ticket, err := client.GetTicket(context.Background(), 123)
	require.NoError(t, err, "Error getting Ticket, expected struct")

	require.Equal(t, 123, ticket.ID, "Expected ticket ID 123")
	require.Equal(t, "Test description", ticket.Description, "Expected ticket description 'Test description'")
}
