package integration

import (
    "context"
    "testing"
)

func TestInvoice_List(t *testing.T) {
    client, teardown := createTestClient(t, "fixtures/TestInvoice_List")
    defer teardown()

    invoices, err := client.ListInvoices(context.Background(), nil)
    if err != nil {
        t.Fatalf("Error getting Invoices, expected struct, got error %v", err)
    }

    if len(invoices) == 0 {
        t.Fatalf("Expected to see invoices returned.")
    }
}

func TestInvoice_Get(t *testing.T) {
    client, teardown := createTestClient(t, "fixtures/TestInvoice_Get")
    defer teardown()

    invoice, err := client.GetInvoice(context.Background(), 123)
    if err != nil {
        t.Fatalf("Error getting Invoice, expected struct, got error %v", err)
    }

    if invoice.ID != 123 {
        t.Fatalf("Expected Invoice ID to be 123, got %v", invoice.ID)
    }

    if invoice.Label != "Invoice" {
        t.Fatalf("Expected Invoice Label to be 'Invoice', got %v", invoice.Label)
    }

    if invoice.Total != 132.5 {
        t.Fatalf("Expected Invoice Total to be 132.5, got %v", invoice.Total)
    }
}

func TestInvoiceItems_List(t *testing.T) {
    client, teardown := createTestClient(t, "fixtures/TestInvoiceItems_List")
    defer teardown()

    items, err := client.ListInvoiceItems(context.Background(), 123, nil)
    if err != nil {
        t.Fatalf("Error getting Invoice Items, expected struct, got error %v", err)
    }

    if len(items) == 0 {
        t.Fatalf("Expected to see invoice items returned.")
    }

    item := items[0]
    if item.Label != "Linode 2GB" {
        t.Fatalf("Expected item label to be 'Linode 2GB', got %v", item.Label)
    }

    if item.Amount != 10 {
        t.Fatalf("Expected item amount to be 10, got %v", item.Amount)
    }
}