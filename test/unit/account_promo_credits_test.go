package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestAccountPromoCredits_Add(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_promo_credits_add_promo_code")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("account/promo-codes", fixtureData)

	requestData := linodego.PromoCodeCreateOptions{
		PromoCode: "supercoolpromocode",
	}

	promoCode, err := base.Client.AddPromoCode(context.Background(), requestData)
	if err != nil {
		t.Fatalf("Error adding promo code: %v", err)
	}

	assert.Equal(t, "10.00", promoCode.CreditMonthlyCap)
	assert.Equal(t, "50.00", promoCode.CreditRemaining)
	assert.Equal(t, "Receive up to $10 off your services every month for 6 months! Unused credits will expire once this promotion period ends.", promoCode.Description)
	assert.Equal(t, "https://linode.com/10_a_month_promotion.svg", promoCode.ImageURL)
	assert.Equal(t, "all", promoCode.ServiceType)
	assert.Equal(t, "$10 off your Linode a month!", promoCode.Summary)
	assert.Equal(t, "10.00", promoCode.ThisMonthCreditRemaining)
}
