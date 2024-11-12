package unit

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountInvoices_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_invoices_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/invoices", fixtureData)

	invoices, err := base.Client.ListInvoices(context.Background(), nil)
	if err != nil {
		t.Fatalf("Error listing invoices: %v", err)
	}

	assert.Equal(t, 1, len(invoices))
	invoice := invoices[0]
	assert.Equal(t, "linode", invoice.BillingSource)
	assert.Equal(t, 123, invoice.ID)
	assert.Equal(t, "Invoice", invoice.Label)
	assert.Equal(t, float32(120.25), invoice.Subtotal)
	assert.Equal(t, float32(12.25), invoice.Tax)
	assert.Equal(t, "PA STATE TAX", invoice.TaxSummary[0].Name)
	assert.Equal(t, float32(12.25), invoice.TaxSummary[0].Tax)
	assert.Equal(t, float32(132.5), invoice.Total)
}

func TestAccountInvoices_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_invoices_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/invoices/123", fixtureData)

	invoice, err := base.Client.GetInvoice(context.Background(), 123)
	if err != nil {
		t.Fatalf("Error getting invoice: %v", err)
	}

	assert.Equal(t, "linode", invoice.BillingSource)
	assert.Equal(t, 123, invoice.ID)
	assert.Equal(t, "Invoice", invoice.Label)
	assert.Equal(t, float32(120.25), invoice.Subtotal)
	assert.Equal(t, float32(12.25), invoice.Tax)
	assert.Equal(t, "PA STATE TAX", invoice.TaxSummary[0].Name)
	assert.Equal(t, float32(12.25), invoice.TaxSummary[0].Tax)
	assert.Equal(t, float32(132.5), invoice.Total)
}

func TestAccountInvoiceItems_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_invoice_items_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/invoices/123/items", fixtureData)

	invoiceitems, err := base.Client.ListInvoiceItems(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("Error listing invoice items: %v", err)
	}

	assert.Equal(t, 1, len(invoiceitems))
	invoiceItem := invoiceitems[0]
	assert.Equal(t, float32(20.2), invoiceItem.Amount)
	assert.Equal(t, "Linode 123", invoiceItem.Label)
	assert.Equal(t, 4, invoiceItem.Quantity)
	assert.Equal(t, "us-west", *invoiceItem.Region)
	assert.Equal(t, float32(1.25), invoiceItem.Tax)
	assert.Equal(t, float32(21.45), invoiceItem.Total)
	assert.Equal(t, "hourly", invoiceItem.Type)
	assert.Equal(t, float32(5.05), invoiceItem.UnitPrice)
}
