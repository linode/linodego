package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount_getTransfer(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_transfer_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/transfer", fixtureData)

	transferInfo, err := base.Client.GetAccountTransfer(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, 0, transferInfo.Billable)
	assert.Equal(t, 9141, transferInfo.Quota)
	assert.Equal(t, 2, transferInfo.Used)

	assert.Len(t, transferInfo.RegionTransfers, 1)
	assert.Equal(t, "us-east", transferInfo.RegionTransfers[0].ID)
	assert.Equal(t, 0, transferInfo.RegionTransfers[0].Billable)
	assert.Equal(t, 5010, transferInfo.RegionTransfers[0].Quota)
	assert.Equal(t, 1, transferInfo.RegionTransfers[0].Used)
}
