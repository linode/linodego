package linodego

import (
	"context"
	"encoding/json"
	"fmt"
)

type MySQLCreateOptions struct {
	Label           string   `json:"label"`
	Region          string   `json:"region"`
	Type            string   `json:"type"`
	Engine          string   `json:"engine"`
	Encrypted       bool     `json:"encrypted"`
	ClusterSize     int      `json:"cluster_size"`
	ReplicationType string   `json:"replication_type"`
	SSLConnection   bool     `json:"ssl_connection"`
	AllowList       []string `json:"allow_list"`
}

func (c *Client) CreateMySQL(ctx context.Context, createOpts MySQLCreateOptions) (*Database, error) {
	var body string
	e, err := c.MySQL.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&Database{})

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
	return r.Result().(*Database), nil
}

func (c *Client) GetMySQL(ctx context.Context, id int) (*Database, error) {
	e, err := c.MySQL.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx)

	e = fmt.Sprintf("%s/%d", e, id)
	r, err := coupleAPIErrors(req.SetResult(&Database{}).Get(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*Database), nil
}

func (c *Client) DeleteMySQL(ctx context.Context, id int) error {
	e, err := c.MySQL.Endpoint()
	if err != nil {
		return err
	}

	req := c.R(ctx)

	e = fmt.Sprintf("%s/%d", e, id)
	_, err = coupleAPIErrors(req.Delete(e))
	return err
}
