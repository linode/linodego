package unit

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestAccountAgreements_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_agreements_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/agreements", fixtureData)

	agreements, err := base.Client.GetAccountAgreements(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, true, agreements.EUModel)
	assert.Equal(t, true, agreements.PrivacyPolicy)
	assert.Equal(t, true, agreements.MasterServiceAgreement)
}

func TestAccountAgreements_Acknowledge(t *testing.T) {
	client := createMockClient(t)

	euModel := true
	masterServiceAgreement := true
	privacyPolicy := true

	requestData := linodego.AccountAgreementsUpdateOptions{
		EUModel:                &euModel,
		MasterServiceAgreement: &masterServiceAgreement,
		PrivacyPolicy:          &privacyPolicy,
	}

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "account/agreements"),
		mockRequestBodyValidate(t, requestData, nil))

	if err := client.AcknowledgeAccountAgreements(context.Background(), requestData); err != nil {
		t.Fatal(err)
	}
}
