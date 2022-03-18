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

type MySQLUpdateOptions struct {
	Label     string   `json:"label"`
	AllowList []string `json:"allow_list"`
}

type MySQLDatabasesPagedResponse struct {
	*PageOptions
	Data []Database `json:"data"`
}

func (MySQLDatabasesPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.MySQL.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (resp *MySQLDatabasesPagedResponse) appendData(r *MySQLDatabasesPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

func (c *Client) ListMySQLDatabases(ctx context.Context, opts *ListOptions) ([]Database, error) {
	response := MySQLDatabasesPagedResponse{}

	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *Client) GetMySQLDatabase(ctx context.Context, id int) (*Database, error) {
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

func (c *Client) CreateMySQLDatabase(ctx context.Context, createOpts MySQLCreateOptions) (*Database, error) {
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

func (c *Client) DeleteMySQLDatabase(ctx context.Context, id int) error {
	e, err := c.MySQL.Endpoint()
	if err != nil {
		return err
	}

	req := c.R(ctx)

	e = fmt.Sprintf("%s/%d", e, id)
	_, err = coupleAPIErrors(req.Delete(e))
	return err
}

func (c *Client) UpdateMySQLDatabase(ctx context.Context, id int, opts MySQLUpdateOptions) (*Database, error) {
	e, err := c.MySQL.Endpoint()
	if err != nil {
		return nil, err
	}
	req := c.R(ctx).SetResult(&Database{})

	bodyData, err := json.Marshal(opts)
	if err != nil {
		return nil, NewError(err)
	}

	body := string(bodyData)

	e = fmt.Sprintf("%s/%d", e, id)
	r, err := coupleAPIErrors(req.SetBody(body).Put(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*Database), nil
}

func (c *Client) GetMySQLDatabaseSSL(ctx context.Context, id int) (*DatabaseSSL, error) {
	e, err := c.MySQL.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx)

	e = fmt.Sprintf("%s/%d/ssl", e, id)
	r, err := coupleAPIErrors(req.SetResult(&DatabaseSSL{}).Get(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*DatabaseSSL), nil
}

func (c *Client) GetMySQLDatabaseCredentials(ctx context.Context, id int) (*DatabaseCredential, error) {
	e, err := c.MySQL.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx)

	e = fmt.Sprintf("%s/%d/credentials", e, id)
	r, err := coupleAPIErrors(req.SetResult(&DatabaseCredential{}).Get(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*DatabaseCredential), nil
}

func (c *Client) ResetMySQLDatabaseCredentials(ctx context.Context, id int) error {
	e, err := c.MySQL.Endpoint()
	if err != nil {
		return err
	}

	req := c.R(ctx)

	e = fmt.Sprintf("%s/%d/credentials/reset", e, id)
	_, err = coupleAPIErrors(req.Post(e))
	if err != nil {
		return err
	}

	return nil
}
