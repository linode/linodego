package unit

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/linode/linodego/internal/testutil"
	"github.com/stretchr/testify/require"
)

var testChildAccount = linodego.ChildAccount{
	Address1:          "123 Main Street",
	Address2:          "Suite A",
	Balance:           200,
	BalanceUninvoiced: 145,
	Capabilities: []string{
		"Linodes",
		"NodeBalancers",
		"Block Storage",
		"Object Storage",
	},
	City:    "Philadelphia",
	Company: "Linode LLC",
	Country: "US",
	CreditCard: &linodego.CreditCard{
		Expiry:   "11/2022",
		LastFour: "1111",
	},
	Email:     "john.smith@linode.com",
	EUUID:     "E1AF5EEC-526F-487D-B317EBEB34C87D71",
	FirstName: "John",
	LastName:  "Smith",
	Phone:     "215-555-1212",
	State:     "PA",
	TaxID:     "ATU99999999",
	Zip:       "19102-1234",
}

func TestAccountChild_list(t *testing.T) {
	client := createMockClient(t)

	desiredResponse := map[string]any{
		"page":    1,
		"pages":   1,
		"results": 1,
		"data": []linodego.ChildAccount{
			testChildAccount,
		},
	}

	httpmock.RegisterRegexpResponder(
		"GET",
		testutil.MockRequestURL("/account/child-accounts"),
		httpmock.NewJsonResponderOrPanic(200, &desiredResponse),
	)

	accounts, err := client.ListChildAccounts(context.Background(), nil)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(accounts, desiredResponse["data"]))
}

func TestAccountChild_get(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder(
		"GET",
		testutil.MockRequestURL(fmt.Sprintf("/account/child-accounts/%s", testChildAccount.EUUID)),
		httpmock.NewJsonResponderOrPanic(200, &testChildAccount),
	)

	account, err := client.GetChildAccount(context.Background(), testChildAccount.EUUID)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(*account, testChildAccount))
}

func TestAccountChild_createToken(t *testing.T) {
	client := createMockClient(t)

	desiredResponse := linodego.ChildAccountToken{
		ID:     123,
		Scopes: "*",
		Label:  "child_token",
	}

	httpmock.RegisterRegexpResponder(
		"POST",
		testutil.MockRequestURL(
			fmt.Sprintf("/account/child-accounts/%s/token", testChildAccount.EUUID),
		),
		httpmock.NewJsonResponderOrPanic(200, &desiredResponse),
	)

	token, err := client.CreateChildAccountToken(context.Background(), testChildAccount.EUUID)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(*token, desiredResponse))
}
