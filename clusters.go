package linodego

import (
	"context"
	"fmt"
)

// Cluster represents a linode cluster object
type Cluster struct {
	ID               string `json:"id"`
	Domain           string `json:"domain"`
	Status           string `json:"status"`
	Region           string `json:"region"`
	StaticSiteDomain string `json:"static_site_domain"`
}

// ClustersPagedResponse represents a linode API response for listing
type ClustersPagedResponse struct {
	*PageOptions
	Data []Cluster `json:"data"`
}

// endpoint gets the endpoint URL for Cluster
func (ClustersPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.Clusters.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// appendData appends Clusters when processing paginated Cluster responses
func (resp *ClustersPagedResponse) appendData(r *ClustersPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListClusters lists Clusters
func (c *Client) ListClusters(ctx context.Context, opts *ListOptions) ([]Cluster, error) {
	response := ClustersPagedResponse{}
	err := c.listHelper(ctx, &response, opts)
	for i := range response.Data {
		response.Data[i].fixDates()
	}
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *Cluster) fixDates() *Cluster {
	return v
}

// GetCluster gets the template with the provided ID
func (c *Client) GetCluster(ctx context.Context, id string) (*Cluster, error) {
	e, err := c.Clusters.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, id)
	r, err := coupleAPIErrors(c.R(ctx).SetResult(&Cluster{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*Cluster).fixDates(), nil
}
