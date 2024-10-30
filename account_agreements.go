package linodego

import "context"

// AccountAgreements represents the agreements and their acceptance status for an Account
type AccountAgreements struct {
	EUModel                bool `json:"eu_model"`
	MasterServiceAgreement bool `json:"master_service_agreement"`
	PrivacyPolicy          bool `json:"privacy_policy"`
}

// AccountAgreementsUpdateOptions fields are those accepted by UpdateAccountAgreements
type AccountAgreementsUpdateOptions struct {
	EUModel                bool `json:"eu_model,omitempty"`
	MasterServiceAgreement bool `json:"master_service_agreement,omitempty"`
	PrivacyPolicy          bool `json:"privacy_policy,omitempty"`
}

// GetUpdateOptions converts an AccountAgreements to AccountAgreementsUpdateOptions for use in UpdateAccountAgreements
func (i AccountAgreements) GetUpdateOptions() (o AccountAgreementsUpdateOptions) {
	o.EUModel = i.EUModel
	o.MasterServiceAgreement = i.MasterServiceAgreement
	o.PrivacyPolicy = i.PrivacyPolicy

	return
}

// GetAccountAgreements gets all agreements and their acceptance status for the Account.
func (c *Client) GetAccountAgreements(ctx context.Context) (*AccountAgreements, error) {
	e := "account/agreements"

	response, err := doGETRequest[AccountAgreements](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// AcknowledgeAccountAgreements acknowledges account agreements for the Account
func (c *Client) AcknowledgeAccountAgreements(ctx context.Context, opts AccountAgreementsUpdateOptions) error {
	e := "account/agreements"
	_, err := doPOSTRequest[AccountAgreements](ctx, c, e, opts)
	if err != nil {
		return err
	}

	return nil
}
