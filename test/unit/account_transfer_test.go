package unit

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
)

func TestAccount_getTransfer(t *testing.T) {
	client := createMockClient(t)

	desiredResponse := linodego.AccountTransfer{
		Billable: 123,
		Quota:    456,
		Used:     789,
		RegionTransfers: []linodego.AccountTransferRegion{
			{
				ID:       "us-southeast",
				Billable: 987,
				Quota:    654,
				Used:     3211,
			},
		},
	}

	httpmock.RegisterRegexpResponder("GET", mockRequestURL(t, "/account/transfer"),
		httpmock.NewJsonResponderOrPanic(200, &desiredResponse))

	questions, err := client.GetAccountTransfer(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*questions, desiredResponse) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(questions, desiredResponse))
	}
}
