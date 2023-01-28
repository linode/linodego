package linodego

import (
	"context"
	"log"
)

type GrantsListResponse = UserGrants

func (c *Client) GrantsList(ctx context.Context) (*GrantsListResponse, error) {
	e := "profile/grants"
	r, err := coupleAPIErrors(c.R(ctx).SetResult(GrantsListResponse{}).Get(e))
	if err != nil {
		return nil, err
	}

	if r.StatusCode() == 204 {
		log.Printf(
			"[WARN] The user has a full account access, " +
				"the instance of the struct would be empty",
		)
	}

	return r.Result().(*GrantsListResponse), err
}
