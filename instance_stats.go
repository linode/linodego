package linodego

import (
	"context"
)

type Net struct {
	In         [][]float64 `json:"in"`
	Out        [][]float64 `json:"out"`
	PrivateIn  [][]float64 `json:"private_in"`
	PrivateOut [][]float64 `json:"private_out"`
}

type IO struct {
	IO   [][]float64 `json:"io"`
	Swap [][]float64 `json:"swap"`
}

type Data struct {
	CPU   [][]float64 `json:"cpu"`
	IO    *IO         `json:"io"`
	Netv4 *Net        `json:"netv4"`
	Netv6 *Net        `json:"netv6"`
}

type InstanceStats struct {
	Title string `json:"title"`
	Data  *Data  `json:"data"`
}

// endpointWithID gets the endpoint URL for InstanceStats of a given Instance
func (InstanceStats) endpointWithID(c *Client, id int) string {
	endpoint, err := c.InstanceStats.endpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// GetInstanceStats gets the template with the provided ID
func (c *Client) GetInstanceStats(ctx context.Context, linodeID int) (*InstanceStats, error) {
	e, err := c.InstanceStats.endpointWithID(linodeID)
	if err != nil {
		return nil, err
	}
	r, err := coupleAPIErrors(c.R(ctx).SetResult(&InstanceStats{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*InstanceStats), nil
}
