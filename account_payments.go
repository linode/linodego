package linodego

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/linode/linodego/internal/parseabletime"
)

// Payment represents a Payment object
type Payment struct {
	// The unique ID of the Payment
	ID int `json:"id"`

	// The amount, in US dollars, of the Payment.
	USD json.Number `json:"usd,Number"`

	// When the Payment was made.
	Date *time.Time `json:"-"`
}

// PaymentCreateOptions fields are those accepted by CreatePayment
type PaymentCreateOptions struct {
	// CVV (Card Verification Value) of the credit card to be used for the Payment
	CVV string `json:"cvv,omitempty"`

	// The amount, in US dollars, of the Payment
	USD json.Number `json:"usd,Number"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *Payment) UnmarshalJSON(b []byte) error {
	type Mask Payment

	p := struct {
		*Mask
		Date *parseabletime.ParseableTime `json:"date"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Date = (*time.Time)(p.Date)

	return nil
}

// GetCreateOptions converts a Payment to PaymentCreateOptions for use in CreatePayment
func (i Payment) GetCreateOptions() (o PaymentCreateOptions) {
	o.USD = i.USD
	return
}

// PaymentsPagedResponse represents a paginated Payment API response
type PaymentsPagedResponse struct {
	*PageOptions
	Data []Payment `json:"data"`
}

// endpoint gets the endpoint URL for Payment
func (PaymentsPagedResponse) endpoint(c *Client, _ ...any) string {
	endpoint, err := c.Payments.Endpoint()
	if err != nil {
		panic(err)
	}

	return endpoint
}

func (resp *PaymentsPagedResponse) castResult(r *resty.Request, e string) (int, int, error) {
	res, err := coupleAPIErrors(r.SetResult(PaymentsPagedResponse{}).Get(e))
	if err != nil {
		return 0, 0, err
	}
	castedRes := res.Result().(*PaymentsPagedResponse)
	resp.Data = append(resp.Data, castedRes.Data...)
	return castedRes.Pages, castedRes.Results, nil
}

// ListPayments lists Payments
func (c *Client) ListPayments(ctx context.Context, opts *ListOptions) ([]Payment, error) {
	response := PaymentsPagedResponse{}
	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetPayment gets the payment with the provided ID
func (c *Client) GetPayment(ctx context.Context, id int) (*Payment, error) {
	e, err := c.Payments.Endpoint()
	if err != nil {
		return nil, err
	}

	e = fmt.Sprintf("%s/%d", e, id)
	r, err := coupleAPIErrors(c.R(ctx).SetResult(&Payment{}).Get(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*Payment), nil
}

// CreatePayment creates a Payment
func (c *Client) CreatePayment(ctx context.Context, createOpts PaymentCreateOptions) (*Payment, error) {
	var body string

	e, err := c.Payments.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&Payment{})

	if bodyData, err := json.Marshal(createOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, NewError(err)
	}

	r, err := coupleAPIErrors(req.
		SetBody(body).
		Post(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*Payment), nil
}
