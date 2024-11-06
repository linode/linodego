package unit

import (
	"context"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccount_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.AccountUpdateOptions{
		City:  "Cambridge",
		State: "MA",
	}

	base.MockPut("account", fixtureData)

	accountInfo, err := base.Client.UpdateAccount(context.Background(), requestData)
	assert.NoError(t, err)

	assert.Equal(t, "John", accountInfo.FirstName)
	assert.Equal(t, "Smith", accountInfo.LastName)
	assert.Equal(t, "john.smith@linode.com", accountInfo.Email)
	assert.Equal(t, "Linode LLC", accountInfo.Company)
	assert.Equal(t, "123 Main Street", accountInfo.Address1)
	assert.Equal(t, "Suite A", accountInfo.Address2)
	assert.Equal(t, float32(200), accountInfo.Balance)
	assert.Equal(t, float32(145), accountInfo.BalanceUninvoiced)
	assert.Equal(t, "19102-1234", accountInfo.Zip)
	assert.Equal(t, "US", accountInfo.Country)
	assert.Equal(t, "ATU99999999", accountInfo.TaxID)
	assert.Equal(t, "215-555-1212", accountInfo.Phone)
	if accountInfo.CreditCard != nil {
		assert.Equal(t, "11/2022", accountInfo.CreditCard.Expiry)
		assert.Equal(t, "1111", accountInfo.CreditCard.LastFour)
	}
	assert.Equal(t, "E1AF5EEC-526F-487D-B317EBEB34C87D71", accountInfo.EUUID)
	assert.Equal(t, "akamai", accountInfo.BillingSource)
	assert.Equal(t, []string{"Linodes", "NodeBalancers", "Block Storage", "Object Storage", "Placement Groups", "Block Storage Encryption"}, accountInfo.Capabilities)

	assert.Equal(t, "Cambridge", accountInfo.City)
	assert.Equal(t, "MA", accountInfo.State)
}
