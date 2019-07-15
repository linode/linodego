package linodego

import (
	"context"
	"encoding/json"
	"fmt"
)

// ObjKeys represents a linode object storage key object
type ObjKey struct {
	ID        int    `json:"id"`
	Label     string `json:"label"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

// ObjKeyCreateOptions fields are those accepted by CreateObjKey
type ObjKeyCreateOptions struct {
	Label string `json:"label"`
}

// ObjKeysPagedResponse represents a linode API response for listing
type ObjKeysPagedResponse struct {
	*PageOptions
	Data []ObjKey `json:"data"`
}

// endpoint gets the endpoint URL for Object Storage keys
func (ObjKeysPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.ObjKeys.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// appendData appends ObjKeys when processing paginated Objkey responses
func (resp *ObjKeysPagedResponse) appendData(r *ObjKeysPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListObjkeys lists Objkeys
func (c *Client) ListObjKeys(ctx context.Context, opts *ListOptions) ([]ObjKey, error) {
	response := ObjKeysPagedResponse{}
	err := c.listHelper(ctx, &response, opts)
	for i := range response.Data {
		response.Data[i].fixDates()
	}
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// CreateObjKey creates a ObjKey
func (c *Client) CreateObjKey(ctx context.Context, createOpts ObjKeyCreateOptions) (*ObjKey, error) {
	var body string
	e, err := c.ObjKeys.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&ObjKey{})

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
	return r.Result().(*ObjKey).fixDates(), nil
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *ObjKey) fixDates() *ObjKey {
	return v
}

// GetObjKey gets the template with the provided ID
func (c *Client) GetObjKey(ctx context.Context, id string) (*ObjKey, error) {
	e, err := c.ObjKeys.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, id)
	r, err := coupleAPIErrors(c.R(ctx).SetResult(&ObjKey{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*ObjKey).fixDates(), nil
}
