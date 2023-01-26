package linodego

import (
	"context"
	"encoding/json"
)

type GrantsListResponse = UserGrants


func (c *Client) GrantsList(ctx context.Context) (*GrantsListResponse, error) {
	e := "profile/grants"
	r, err := coupleAPIErrors(c.R(ctx).Get(e))
	if err != nil {
		return nil, err
	}

	if r.StatusCode() == 204 {
		return nil, nil
	}
	var result GrantsListResponse
	err = json.Unmarshal(r.Body(), &result)

	return &result, err
}
