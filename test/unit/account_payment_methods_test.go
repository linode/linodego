package unit

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestAccountPaymentMethods_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_payment_methods_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/payment-methods/123", fixtureData)

	pm, err := base.Client.GetPaymentMethod(context.Background(), 123)
	assert.NoError(t, err)

	assert.Equal(t, 123, pm.ID)
	assert.Equal(t, true, pm.IsDefault)
	assert.Equal(t, "credit_card", pm.Type)
	assert.Equal(t, linodego.PaymentMethodDataCreditCard{
		CardType: "Discover",
		Expiry:   "06/2022",
		LastFour: "1234",
	}, pm.Data)
}

func TestAccountPaymentMethods_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_payment_methods_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/payment-methods", fixtureData)

	methods, err := base.Client.ListPaymentMethods(context.Background(), nil)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(methods))
	pm := methods[0]
	assert.Equal(t, 123, pm.ID)
	assert.Equal(t, true, pm.IsDefault)
	assert.Equal(t, "credit_card", pm.Type)
	assert.Equal(t, linodego.PaymentMethodDataCreditCard{
		CardType: "Discover",
		Expiry:   "06/2022",
		LastFour: "1234",
	}, pm.Data)
}

func TestAccountPaymentMethods_Add(t *testing.T) {
	client := createMockClient(t)

	card := linodego.PaymentMethodCreateOptionsData{
		CardNumber:  "1234123412341234",
		CVV:         "123",
		ExpiryMonth: 3,
		ExpiryYear:  2028,
	}

	requestData := linodego.PaymentMethodCreateOptions{
		Data:      &card,
		IsDefault: true,
		Type:      "credit_card",
	}

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "account/payment-methods"),
		mockRequestBodyValidate(t, requestData, nil))

	if err := client.AddPaymentMethod(context.Background(), requestData); err != nil {
		t.Fatal(err)
	}
}

func TestAccountPaymentMethods_Delete(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, "account/payment-methods/123"), httpmock.NewStringResponder(200, "{}"))

	if err := client.DeletePaymentMethod(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestAccountPaymentMethods_SetDefault(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "account/payment-methods/123"),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.SetDefaultPaymentMethod(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}
