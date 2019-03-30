package linodego

import (
	"context"
)

type Net struct {
    In []float64 `json:"in"`
    Out []float64 `json:"out"`
    PrivateIn []float64 `json:"private_in"`
    PrivateOut []float64 `json:"private_out"`
}

type IO struct {
    IO []float64 `json:"io"`
    Swap []float64 `json:"swap"`
}

type Data struct {
    CPU []float64 `json:"cpu"`
    IO *IO `json:"io"`
    Netv4 *Net `json:"netv4"`
    Netv6 *Net `json:"netv6"`
}

type InstanceStatsResponse struct {
    Title string `json:"title"`
    Data *Data `json:"data"`
}

// endpoint gets the endpoint URL for Stats
func (InstanceStatsResponse) endpoint(c *Client) string {
	endpoint, err := c.InstanceStats.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}


// Get Stats
func (c *Client) GetStats(ctx context.Context) (InstanceStatsResponse, error) {
	response := InstanceStatsResponse{}
	return response, nil
}

