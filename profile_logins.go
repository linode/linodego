package linodego

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// Profile represents a Profile object
type ProfileLogin struct {
	Datetime   *time.Time `json:"_"`
	ID         int        `json:"id"`
	IP         string     `json:"ip"`
	Restricted bool       `json:"restricted"`
	Status     string     `json:"status"`
	Username   string     `json:"username"`
}

type ProfileLoginsPagedResponse struct {
	*PageOptions
	Data []ProfileLogin `json:"data"`
}

func (ProfileLoginsPagedResponse) endpoint(_ ...any) string {
	return "profile/logins"
}

func (resp *ProfileLoginsPagedResponse) castResult(r *resty.Request, e string) (int, int, error) {
	res, err := coupleAPIErrors(r.SetResult(ProfileLoginsPagedResponse{}).Get(e))
	if err != nil {
		return 0, 0, err
	}
	castedRes := res.Result().(*ProfileLoginsPagedResponse)
	resp.Data = append(resp.Data, castedRes.Data...)
	return castedRes.Pages, castedRes.Results, nil
}

// GetProfileLogin returns the Profile Login of the authenticated user
func (c *Client) GetProfileLogin(ctx context.Context, id int) (*ProfileLogin, error) {
	e := fmt.Sprintf("profile/logins/%d", id)

	req := c.R(ctx).SetResult(&ProfileLogin{})
	r, err := coupleAPIErrors(req.Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*ProfileLogin), nil
}

// ListProfileLogins lists Profile Logins of the authenticated user
func (c *Client) ListProfileLogins(ctx context.Context, opts *ListOptions) ([]ProfileLogin, error) {
	response := ProfileLoginsPagedResponse{}
	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
