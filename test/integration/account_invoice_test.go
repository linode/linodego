package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvoice_List(t *testing.T) {
	warnSensitiveTest(t)
	client, teardown := createTestClient(t, "fixtures/TestInvoice_List")
	defer teardown()

	invoices, err := client.ListInvoices(context.Background(), nil)
	require.NoError(t, err, "Error getting Invoices, expected struct")
	require.NotEmpty(t, invoices, "Expected to see invoices returned")
}

func TestInvoice_Get(t *testing.T) {
	warnSensitiveTest(t)
	client, teardown := createTestClient(t, "fixtures/TestInvoice_Get")
	defer teardown()

	invoice, err := client.GetInvoice(context.Background(), 123)
	require.NoError(t, err, "Error getting Invoice, expected struct")
	require.Equal(t, 123, invoice.ID, "Expected Invoice ID to be 123")
	require.Equal(t, "Invoice", invoice.Label, "Expected Invoice Label to be 'Invoice'")
	require.Equal(t, 132.5, float64(invoice.Total), "Expected Invoice Total to be 132.5")
}

func TestInvoiceItems_List(t *testing.T) {
	warnSensitiveTest(t)
	client, teardown := createTestClient(t, "fixtures/TestInvoiceItems_List")
	defer teardown()

	items, err := client.ListInvoiceItems(context.Background(), 123, nil)
	require.NoError(t, err, "Error getting Invoice Items, expected struct")
	require.NotEmpty(t, items, "Expected to see invoice items returned")

	item := items[0]
	require.Equal(t, "Linode 2GB", item.Label, "Expected item label to be 'Linode 2GB'")
	require.Equal(t, 10.0, float64(item.Amount), "Expected item amount to be 10")
}
