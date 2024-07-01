package integration

import (
    "context"
    "testing"
)

func TestTicket_List(t *testing.T) {
    client, teardown := createTestClient(t, "fixtures/TestTicket_List")
    defer teardown()

    tickets, err := client.ListTickets(context.Background(), nil)
    if err != nil {
        t.Fatalf("Error getting Tickets, expected struct, got error %v", err)
    }

    if len(tickets) == 0 {
        t.Fatalf("Expected to see tickets returned.")
    }

    if tickets[0].ID != 123 {
        t.Fatalf("Expected ticket ID 123, got %d", tickets[0].ID)
    }
}

func TestTicket_Get(t *testing.T) {
    client, teardown := createTestClient(t, "fixtures/TestTicket_Get")
    defer teardown()

    ticket, err := client.GetTicket(context.Background(), 123)
    if err != nil {
        t.Fatalf("Error getting Ticket, expected struct, got error %v", err)
    }

    if ticket.ID != 123 {
        t.Fatalf("Expected ticket ID 123, got %d", ticket.ID)
    }

    if ticket.Description != "Test description" {
        t.Fatalf("Expected ticket description 'Test description', got %s", ticket.Description)
    }
}
