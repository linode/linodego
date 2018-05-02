package golinode

import (
	"fmt"
)

type InstanceIPAddressResponse struct {
	IPv4 *InstanceIPv4Response
	IPv6 *InstanceIPv6Response
}

type InstanceIPv4Response struct {
	Public  []*InstanceIP
	Private []*InstanceIP
	Shared  []*InstanceIP
}

type InstanceIP struct {
	Address    string
	Gateway    string
	SubnetMask string
	Prefix     int
	Type       string
	Public     bool
	RDNS       string
	LinodeID   int `json:"linode_id"`
	Region     string
}

type InstanceIPv6Response struct {
	LinkLocal *InstanceIP `json:"link_local"`
	SLAAC     *InstanceIP
	Global    []*IPv6Range
}

type IPv6Range struct {
	Range  string
	Region string
}

// GetInstanceIPAddress gets the template with the provided ID
func (c *Client) GetInstanceIPAddress(linodeID int, ipaddress string) (*InstanceIPAddressResponse, error) {
	e, err := c.InstanceIPs.EndpointWithID(linodeID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, ipaddress)
	r, err := c.R().SetResult(&InstanceIPAddressResponse{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*InstanceIPAddressResponse), nil
}
