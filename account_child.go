package linodego

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/url"
)

// ChildAccount represents an account under the current account.
// NOTE: This is an alias to prevent any future breaking changes.
type ChildAccount = Account

// ChildAccountToken represents a short-lived token created using
// the CreateChildAccountToken(...) function.
// NOTE: This is an alias to prevent any future breaking changes.
type ChildAccountToken = Token

// ChildAccountsPagedResponse represents a Linode API response
// for listing child accounts under the current account.
type ChildAccountsPagedResponse struct {
	*PageOptions
	Data []ChildAccount `json:"data"`
}

// endpoint gets the endpoint URL for Instance
func (ChildAccountsPagedResponse) endpoint(_ ...any) string {
	return "account/child-accounts"
}

func (resp *ChildAccountsPagedResponse) castResult(r *resty.Request, e string) (int, int, error) {
	res, err := coupleAPIErrors(r.SetResult(ChildAccountsPagedResponse{}).Get(e))
	if err != nil {
		return 0, 0, err
	}

	castedRes := res.Result().(*ChildAccountsPagedResponse)
	resp.Data = append(resp.Data, castedRes.Data...)
	return castedRes.Pages, castedRes.Results, nil
}

// ListChildAccounts lists child accounts under the current account.
func (c *Client) ListChildAccounts(ctx context.Context, opts *ListOptions) ([]ChildAccount, error) {
	response := ChildAccountsPagedResponse{}

	err := c.listHelper(ctx, &response, opts)

	return response.Data, err
}

// GetChildAccount gets a single child accounts under the current account.
func (c *Client) GetChildAccount(ctx context.Context, euuid string) (*ChildAccount, error) {
	e := fmt.Sprintf("account/child-accounts/%s", url.PathEscape(euuid))
	req := c.R(ctx).SetResult(ChildAccount{})
	r, err := coupleAPIErrors(req.Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*ChildAccount), nil
}

// CreateChildAccountToken creates a short-lived token that can be used to
// access the Linode API under a child account.
// The attributes of this token are not currently configurable.
func (c *Client) CreateChildAccountToken(ctx context.Context, euuid string) (*ChildAccountToken, error) {
	e := fmt.Sprintf("account/child-accounts/%s/token", url.PathEscape(euuid))

	req := c.R(ctx).SetResult(&ChildAccountToken{})
	r, err := coupleAPIErrors(req.Post(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*ChildAccountToken), nil
}
