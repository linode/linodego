package unit

import (
	"context"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestAccountServiceTransfer_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_service_transfers_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/service-transfers", fixtureData)

	transfers, err := base.Client.ListAccountServiceTransfer(context.Background(), nil)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(transfers))
	ast := transfers[0]
	assert.Equal(t, time.Time(time.Date(2021, time.February, 11, 16, 37, 3, 0, time.UTC)), *ast.Created)
	assert.Equal(t, time.Time(time.Date(2021, time.February, 12, 16, 37, 3, 0, time.UTC)), *ast.Expiry)
	assert.Equal(t, time.Time(time.Date(2021, time.February, 11, 16, 37, 3, 0, time.UTC)), *ast.Updated)
	assert.Equal(t, 111, ast.Entities.Linodes[0])
	assert.Equal(t, 222, ast.Entities.Linodes[1])
	assert.Equal(t, true, ast.IsSender)
	assert.Equal(t, linodego.AccountServiceTransferStatus("pending"), ast.Status)
	assert.Equal(t, "123E4567-E89B-12D3-A456-426614174000", ast.Token)
}

func TestAccountServiceTransfer_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_service_transfers_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/service-transfers/123E4567-E89B-12D3-A456-426614174000", fixtureData)

	ast, err := base.Client.GetAccountServiceTransfer(context.Background(), "123E4567-E89B-12D3-A456-426614174000")
	assert.NoError(t, err)

	assert.Equal(t, time.Time(time.Date(2021, time.February, 11, 16, 37, 3, 0, time.UTC)), *ast.Created)
	assert.Equal(t, time.Time(time.Date(2021, time.February, 12, 16, 37, 3, 0, time.UTC)), *ast.Expiry)
	assert.Equal(t, time.Time(time.Date(2021, time.February, 11, 16, 37, 3, 0, time.UTC)), *ast.Updated)
	assert.Equal(t, 111, ast.Entities.Linodes[0])
	assert.Equal(t, 222, ast.Entities.Linodes[1])
	assert.Equal(t, true, ast.IsSender)
	assert.Equal(t, linodego.AccountServiceTransferStatus("pending"), ast.Status)
	assert.Equal(t, "123E4567-E89B-12D3-A456-426614174000", ast.Token)
}

func TestAccountServiceTransfer_Request(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_service_transfers_request")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.AccountServiceTransferRequestOptions{
		Entities: linodego.AccountServiceTransferEntity{
			Linodes: []int{111, 222},
		},
	}

	base.MockPost("account/service-transfers", fixtureData)

	ast, err := base.Client.RequestAccountServiceTransfer(context.Background(), requestData)
	assert.NoError(t, err)

	assert.Equal(t, time.Time(time.Date(2021, time.February, 11, 16, 37, 3, 0, time.UTC)), *ast.Created)
	assert.Equal(t, time.Time(time.Date(2021, time.February, 12, 16, 37, 3, 0, time.UTC)), *ast.Expiry)
	assert.Equal(t, time.Time(time.Date(2021, time.February, 11, 16, 37, 3, 0, time.UTC)), *ast.Updated)
	assert.Equal(t, 111, ast.Entities.Linodes[0])
	assert.Equal(t, 222, ast.Entities.Linodes[1])
	assert.Equal(t, true, ast.IsSender)
	assert.Equal(t, linodego.AccountServiceTransferStatus("pending"), ast.Status)
	assert.Equal(t, "123E4567-E89B-12D3-A456-426614174000", ast.Token)
}

func TestAccountServiceTransfer_Accept(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST",
		mockRequestURL(t, "account/service-transfers/123E4567-E89B-12D3-A456-426614174000/accept"),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.AcceptAccountServiceTransfer(context.Background(), "123E4567-E89B-12D3-A456-426614174000"); err != nil {
		t.Fatal(err)
	}
}

func TestAccountServiceTransfer_Cancel(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("DELETE",
		mockRequestURL(t, "account/service-transfers/123E4567-E89B-12D3-A456-426614174000"),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.CancelAccountServiceTransfer(context.Background(), "123E4567-E89B-12D3-A456-426614174000"); err != nil {
		t.Fatal(err)
	}
}
