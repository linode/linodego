package unit

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestAccountPayments_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_payment_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("account/payments", fixtureData)

	cvv := "123"

	requestData := linodego.PaymentCreateOptions{
		CVV: &cvv,
		USD: json.Number("120.50"),
	}

	payment, err := base.Client.CreatePayment(context.Background(), requestData)
	if err != nil {
		t.Fatalf("Error creating payment: %v", err)
	}

	assert.Equal(t, 123, payment.ID)
	assert.Equal(t, json.Number("120.50"), payment.USD)
}
