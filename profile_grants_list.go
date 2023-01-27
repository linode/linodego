package linodego

import (
	"context"
	"encoding/json"
	"log"
)

type GrantsListResponse = UserGrants

func (c *Client) GrantsList(ctx context.Context) (*GrantsListResponse, error) {
	e := "profile/grants"
	r, err := coupleAPIErrors(c.R(ctx).Get(e))
	if err != nil {
		return nil, err
	}

	var result GrantsListResponse
	err = nil
	if r.StatusCode() == 204 {
		log.Printf(
			"[WARN] The user has a full account access, " +
				"the instance of the struct would be empty",
		)
	} else {
		err = json.Unmarshal(r.Body(), &result)
	}
	return &result, err
}
