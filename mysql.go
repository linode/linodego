package linodego

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// MySQLCreateOptions fields are used when creating a new MySQL Database
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

// MySQLUpdateOptions fields are used when altering the existing MySQL Database
type MySQLUpdateOptions struct {
	Label     string   `json:"label"`
	AllowList []string `json:"allow_list"`
}

// MySQLDatabaseBackup is information for interacting with a backup for the existing MySQL Database
type MySQLDatabaseBackup struct {
	ID      int
	Label   string
	Type    string
	Created time.Time
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

// List all MySQL Databases associated with the account
func (c *Client) ListMySQLDatabases(ctx context.Context, opts *ListOptions) ([]Database, error) {
	response := MySQLDatabasesPagedResponse{}

	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

type MySQLDatabaseBackupsPagedResponse struct {
	*PageOptions
	Data []MySQLDatabaseBackup `json:"data"`
}

func (MySQLDatabaseBackupsPagedResponse) endpointWithID(c *Client, id int) string {
	endpoint, err := c.MySQL.Endpoint()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/%d/backups", endpoint, id)
}

func (resp *MySQLDatabaseBackupsPagedResponse) appendData(r *MySQLDatabaseBackupsPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// List all MySQL Database Backups associated with the given MySQL Database
func (c *Client) ListMySQLDatabaseBackups(ctx context.Context, databaseID int, opts *ListOptions) ([]MySQLDatabaseBackup, error) {
	response := MySQLDatabaseBackupsPagedResponse{}

	err := c.listHelperWithID(ctx, &response, databaseID, opts)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// Get a single MySQL Database matching the id
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

// Create a new MySQL Database using the createOpts as configuration, returns the new MySQL Database
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

// Delete an existing MySQL Database with the given id
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

// Update the given MySQL Database with the provided opts, returns the Database with the new settings
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

// Get the SSL Certificate for the given MySQL Database
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

// Get the Root Credentials for the given MySQL Database
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

// Reset the Root Credentials for the given MySQL Database (may take a few seconds to work)
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

// Get a specific MySQL Database Backup with the given ids
func (c *Client) GetMySQLDatabaseBackup(ctx context.Context, databaseID int, backupID int) (*MySQLDatabaseBackup, error) {
	e, err := c.MySQL.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx)

	e = fmt.Sprintf("%s/%d/backups/%d", e, databaseID, backupID)
	r, err := coupleAPIErrors(req.SetResult(&MySQLDatabaseBackup{}).Get(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*MySQLDatabaseBackup), nil
}

// Restore the given MySQL Database with the given Backup
func (c *Client) RestoreMySQLDatabaseBackup(ctx context.Context, databaseID int, backupID int) error {
	e, err := c.MySQL.Endpoint()
	if err != nil {
		return err
	}

	req := c.R(ctx)

	e = fmt.Sprintf("%s/%d/backups/%d/restore", e, databaseID, backupID)
	_, err = coupleAPIErrors(req.Post(e))
	if err != nil {
		return err
	}
	return nil
}
