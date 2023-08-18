package linodego

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/linode/linodego/internal/parseabletime"
)

// The details and enrollment information of a Beta program that a customer is enrolled in.
type CustomerBetaProgram struct {
	Label       string     `json:"label"`
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Started     *time.Time `json:"-"`
	Ended       *time.Time `json:"-"`

	// Date the customer was enrolled in the beta program.
	Enrolled *time.Time `json:"-"`
}

// CustomerBetaProgramCreateOpts fields are those accepted by CreateCustomerBetaProgram
type CustomerBetaProgramCreateOpts struct {
	ID string `json:"id"`
}

// CustomerBetasPagedResponse represents a paginated Customer Beta Programs API response
type CustomerBetasPagedResponse struct {
	*PageOptions
	Data []CustomerBetaProgram `json:"data"`
}

// endpoint gets the endpoint URL for CustomerBetaProgram
func (CustomerBetasPagedResponse) endpoint(_ ...any) string {
	return "/account/betas"
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (cBeta *CustomerBetaProgram) UnmarshalJSON(b []byte) error {
	type Mask CustomerBetaProgram

	p := struct {
		*Mask
		Started  *parseabletime.ParseableTime `json:"started"`
		Ended    *parseabletime.ParseableTime `json:"ended"`
		Enrolled *parseabletime.ParseableTime `json:"enrolled"`
	}{
		Mask: (*Mask)(cBeta),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	cBeta.Started = (*time.Time)(p.Started)
	cBeta.Ended = (*time.Time)(p.Ended)
	cBeta.Enrolled = (*time.Time)(p.Enrolled)

	return nil
}

func (resp *CustomerBetasPagedResponse) castResult(r *resty.Request, e string) (int, int, error) {
	res, err := coupleAPIErrors(r.SetResult(CustomerBetasPagedResponse{}).Get(e))
	if err != nil {
		return 0, 0, err
	}
	castedRes := res.Result().(*CustomerBetasPagedResponse)
	resp.Data = append(resp.Data, castedRes.Data...)
	return castedRes.Pages, castedRes.Results, nil
}

// ListCustomerBetaPrograms lists all beta programs a customer is enrolled in.
func (c *Client) ListCustomerBetaPrograms(ctx context.Context, opts *ListOptions) ([]CustomerBetaProgram, error) {
	response := CustomerBetasPagedResponse{}
	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetCustomerBetaProgram gets the details of a beta program a customer is enrolled in.
func (c *Client) GetCustomerBetaProgram(ctx context.Context, betaID string) (*CustomerBetaProgram, error) {
	req := c.R(ctx).SetResult(&CustomerBetaProgram{})
	betaID = url.PathEscape(betaID)
	b := fmt.Sprintf("/account/betas/%s", betaID)
	r, err := coupleAPIErrors(req.Get(b))
	if err != nil {
		return nil, err
	}

	return r.Result().(*CustomerBetaProgram), nil
}

// CreateCustomerBetaProgram enrolls a customer in a beta program.
func (c *Client) CreateCustomerBetaProgram(ctx context.Context, opts CustomerBetaProgramCreateOpts) (*CustomerBetaProgram, error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	e := "account/betas"
	req := c.R(ctx).SetResult(&CustomerBetaProgram{}).SetBody(string(body))
	r, err := coupleAPIErrors(req.Post(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*CustomerBetaProgram), nil
}
