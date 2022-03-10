package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

type DatabasesPagedResponse struct {
	*PageOptions
	Data []Database `json:"data"`
}

func (DatabasesPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.Databases.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// A Database is a managed DBaaS instance
type Database struct {
	ID          int           `json:"id"`
	Status      string        `json:"status"`
	Label       string        `json:"label"`
	Hosts       DatabaseHosts `json:"hosts"`
	Region      string        `json:"region"`
	Type        string        `json:"type"`
	Engine      string        `json:"engine"`
	Version     string        `json:"version"`
	ClusterSize int           `json:"cluster_size"`
	Encrypted   bool          `json:"encrypted"`
	AllowList   []string      `json:"allow_list"`
	InstanceURI string        `json:"instance_uri"`
	Created     *time.Time    `json:"-"`
	Updated     *time.Time    `json:"-"`
}

// DatabaseHost for Primary/Secondary of Database
type DatabaseHosts struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary,omitempty"`
}

// UnmasrhalJSON for Database responses
func (d *Database) UnmarshalJSON(b []byte) error {
	type Mask Database

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		Mask: (*Mask)(d),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	d.Created = (*time.Time)(p.Created)
	d.Updated = (*time.Time)(p.Updated)
	return nil
}

func (c *Client) ListDatabases(ctx context.Context, opts *ListOptions) ([]Database, error) {
	response := DatabasesPagedResponse{}

	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
