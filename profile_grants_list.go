package linodego

import (
	"context"
)

type GrantsListResponse = UserGrants

func (c *Client) GrantsList(ctx context.Context) (*GrantsListResponse, error) {
	e := "profile/grants"
	r, err := coupleAPIErrors(c.R(ctx).SetResult(GrantsListResponse{}).Get(e))
	if err != nil {
		return nil, err
	}
	// TODO: handle 204 empty content response
	return r.Result().(*GrantsListResponse), err
}
