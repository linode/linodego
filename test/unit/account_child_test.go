package unit

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var fixtures *TestFixtures

func TestMain(m *testing.M) {
	fixtures = NewTestFixtures()

	code := m.Run()

	os.Exit(code)
}

func TestAccountChild_list(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_child_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/child-accounts", fixtureData)

	accounts, err := base.Client.ListChildAccounts(context.Background(), nil)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(accounts), "Expected one child account")

	if len(accounts) > 0 {
		childAccount := accounts[0]
		assert.Equal(t, "123 Main Street", childAccount.Address1)
		assert.Equal(t, "Suite A", childAccount.Address2)
		assert.Equal(t, float32(200), childAccount.Balance)
		assert.Equal(t, float32(145), childAccount.BalanceUninvoiced)
		assert.Equal(t, "San Diego", childAccount.City)
		assert.Equal(t, "john.smith@linode.com", childAccount.Email)
		assert.Equal(t, "858-555-1212", childAccount.Phone)
		assert.Equal(t, "CA", childAccount.State)
		assert.Equal(t, "92111-1234", childAccount.Zip)
	}
}

func TestAccountChild_get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_child_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/child-accounts/A1BC2DEF-34GH-567I-J890KLMN12O34P56", fixtureData)

	account, err := base.Client.GetChildAccount(context.Background(), "A1BC2DEF-34GH-567I-J890KLMN12O34P56")
	assert.NoError(t, err)

	assert.Equal(t, "John", account.FirstName)
	assert.Equal(t, "Smith", account.LastName)
	assert.Equal(t, "123 Main Street", account.Address1)
	assert.Equal(t, "Suite A", account.Address2)
	assert.Equal(t, float32(200), account.Balance)
	assert.Equal(t, float32(145), account.BalanceUninvoiced)
	assert.Equal(t, "San Diego", account.City)
	assert.Equal(t, "john.smith@linode.com", account.Email)
	assert.Equal(t, "858-555-1212", account.Phone)
	assert.Equal(t, "CA", account.State)
	assert.Equal(t, "92111-1234", account.Zip)
	assert.Equal(t, "US", account.Country)
	assert.Equal(t, "external", account.BillingSource)
	assert.Equal(t, []string{"Linodes", "NodeBalancers", "Block Storage", "Object Storage"}, account.Capabilities)

	if account.CreditCard != nil {
		assert.Equal(t, "11/2024", account.CreditCard.Expiry)
		assert.Equal(t, "1111", account.CreditCard.LastFour)
	}

	assert.Equal(t, "A1BC2DEF-34GH-567I-J890KLMN12O34P56", account.EUUID)
}

func TestAccountChild_createToken(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_child_create_token")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("account/child-accounts/A1BC2DEF-34GH-567I-J890KLMN12O34P56/token", fixtureData)

	token, err := base.Client.CreateChildAccountToken(context.Background(), "A1BC2DEF-34GH-567I-J890KLMN12O34P56")
	assert.NoError(t, err)

	// Assertions for the created token data
	assert.Equal(t, 918, token.ID)
	assert.Equal(t, "parent1_1234_2024-05-01T00:01:01", token.Label)
	assert.Equal(t, "2024-05-01T00:01:01Z", token.Created.Format(time.RFC3339))
	assert.Equal(t, "2024-05-01T00:16:01Z", token.Expiry.Format(time.RFC3339))
	assert.Equal(t, "*", token.Scopes)
	assert.Equal(t, "abcdefghijklmnop", token.Token)
}
