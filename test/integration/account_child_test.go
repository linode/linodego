package integration

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"reflect"
	"testing"
)

// NOTE: The tests in this file are not implemented as E2E tests because
// child accounts are not currently self-serve.

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

	desiredResponse := linodego.ChildAccountsPagedResponse{
		PageOptions: &linodego.PageOptions{
			Page:    1,
			Pages:   1,
			Results: 1,
		},
		Data: []linodego.ChildAccount{testChildAccount},
	}

	httpmock.RegisterRegexpResponder("GET", mockRequestURL(t, "/account/child-accounts"),
		httpmock.NewJsonResponderOrPanic(200, &desiredResponse))

	accounts, err := client.ListChildAccounts(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(accounts, desiredResponse.Data) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(accounts, desiredResponse))
	}
}

func TestAccountChild_get(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder(
		"GET",
		mockRequestURL(t, fmt.Sprintf("/account/child-accounts/%s", testChildAccount.EUUID)),
		httpmock.NewJsonResponderOrPanic(200, &testChildAccount),
	)

	account, err := client.GetChildAccount(context.Background(), testChildAccount.EUUID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*account, testChildAccount) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(account, testChildAccount))
	}
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
		mockRequestURL(
			t,
			fmt.Sprintf("/account/child-accounts/%s/token", testChildAccount.EUUID),
		),
		httpmock.NewJsonResponderOrPanic(200, &desiredResponse),
	)

	token, err := client.CreateChildAccountToken(context.Background(), testChildAccount.EUUID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*token, desiredResponse) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(token, desiredResponse))
	}
}
