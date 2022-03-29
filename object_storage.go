package linodego

import (
	"context"
	"fmt"
)

// ObjectStorageTransfer is an object matching the response of object-storage/transfer
type ObjectStorageTransfer struct {
	AmmountUsed int `json:"used"`
}

// CancelObjectStorage cancels and removes all object storage from the Account
func (c *Client) CancelObjectStorage(ctx context.Context) error {
	e, err := c.ObjectStorage.Endpoint()
	if err != nil {
		return err
	}
	e = fmt.Sprintf("%s/cancel", e)
	_, err = coupleAPIErrors(c.R(ctx).Post(e))
	if err != nil {
		return err
	}
	return nil
}

// GetObjectStorageTransfer returns the amount of outbound data transferred used by the Account
func (c *Client) GetObjectStorageTransfer(ctx context.Context) (*ObjectStorageTransfer, error) {
	e, err := c.ObjectStorage.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/transfer", e)
	r, err := coupleAPIErrors(c.R(ctx).SetResult(&ObjectStorageTransfer{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*ObjectStorageTransfer), nil
}
