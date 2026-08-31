package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/v2/internal/parseabletime"
)

type NATGateway struct {
	ID                       int                 `json:"id"`
	Region                   string              `json:"region"`
	Addresses                []NATGatewayAddress `json:"addresses"`
	AddressAutoscaleMax      int                 `json:"address_autoscale_max"`
	DefaultPortsPerInterface int                 `json:"default_ports_per_interface"`
	Label                    string              `json:"label"`
	PortsetAssignments       int                 `json:"portset_assignments"`
	PortsetCapacity          int                 `json:"portset_capacity"`
	VPCSubnet                NATGatewayVPCSubnet `json:"vpc_subnet"`
	Created                  *time.Time          `json:"-"`
	Updated                  *time.Time          `json:"-"`
}

func (n *NATGateway) UnmarshalJSON(b []byte) error {
	type Mask NATGateway

	p := struct {
		*Mask

		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		Mask: (*Mask)(n),
	}
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	n.Created = (*time.Time)(p.Created)
	n.Updated = (*time.Time)(p.Updated)

	return nil
}

type NATGatewayAddress struct {
	Address string `json:"address"`
}

type NATGatewayAddressCreateOptions struct {
	Address string `json:"address"`
}

type NATGatewayVPCSubnet struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Label    string `json:"label"`
	URL      string `json:"url"`
	VPCID    int    `json:"vpc_id"`
	VPCLabel string `json:"vpc_label"`
}

type NATGatewayCreateOptions struct {
	Region                   string                           `json:"region"`
	Addresses                []NATGatewayAddressCreateOptions `json:"addresses"`
	DefaultPortsPerInterface *int                             `json:"default_ports_per_interface,omitzero"`
	Label                    string                           `json:"label"`
	UseAutoscaling           *bool                            `json:"use_autoscaling,omitzero"`
	VPCSubnetID              *int                             `json:"vpc_subnet_id,omitzero"`
}

type NATGatewayUpdateOptions struct {
	Label *string `json:"label,omitempty"`
}

type NATGatewayAddAddressOptions struct {
	Address string `json:"address"`
}

type NATGatewayAddressObject struct {
	Address            string `json:"address"`
	InUse              bool   `json:"in_use"`
	InterfaceCount     int    `json:"interface_count"`
	InterfaceURL       string `json:"interface_url"`
	PortsetAssignments int    `json:"portset_assignments"`
	PortsetCapacity    int    `json:"portset_capacity"`
}

type NATGatewayInterface struct {
	ID        int                       `json:"id"`
	Linode    NATGatewayInterfaceLinode `json:"linode"`
	Addresses []string                  `json:"addresses"`

	// NOTE: This field may not be available for all customers
	Portsets []NATGatewayPortset `json:"portsets"`
}

type NATGatewayInterfaceLinode struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

type NATGatewayPortset struct {
	Address string                  `json:"address"`
	Ports   []NATGatewayPortsetPort `json:"ports"`
}

type NATGatewayPortsetPort struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type NATGatewayType struct {
	ID    string        `json:"id"`
	Label string        `json:"label"`
	Price baseTypePrice `json:"price"`
}

// TODO: replace with baseType once the region_prices and transfer fields are in use and returned by the API

type NATGatewaySettings struct {
	AllowedPortsPerInterface                 []int `json:"allowed_ports_per_interface"`
	MaximumAutoscalingAddressesPerNATGateway int   `json:"maximum_autoscaling_addresses_per_natgateway"`
	MaximumReservedAddressesPerNATGateway    int   `json:"maximum_reserved_addresses_per_natgateway"`
}

// ListNATGateways lists NAT Gateways
func (c *Client) ListNATGateways(ctx context.Context, opts *ListOptions) ([]NATGateway, error) {
	return getPaginatedResults[NATGateway](ctx, c, "networking/natgateways", opts)
}

// GetNATGateway gets the NAT Gateway with the specified id
func (c *Client) GetNATGateway(ctx context.Context, id int) (*NATGateway, error) {
	e := formatAPIPath("networking/natgateways/%d", id)
	return doGETRequest[NATGateway](ctx, c, e)
}

// CreateNATGateway creates a new NAT Gateway
func (c *Client) CreateNATGateway(ctx context.Context, opts NATGatewayCreateOptions) (*NATGateway, error) {
	return doPOSTRequest[NATGateway](ctx, c, "networking/natgateways", opts)
}

// UpdateNATGateway updates the NAT Gateway with the specified id
func (c *Client) UpdateNATGateway(ctx context.Context, opts NATGatewayUpdateOptions, id int) (*NATGateway, error) {
	e := formatAPIPath("networking/natgateways/%d", id)
	return doPUTRequest[NATGateway](ctx, c, e, opts)
}

// DeleteNATGateway deletes the NAT Gateway with the specified id
func (c *Client) DeleteNATGateway(ctx context.Context, id int) error {
	e := formatAPIPath("networking/natgateways/%d", id)
	return doDELETERequest(ctx, c, e)
}

// NATGatewayAddAddress adds an addresd to the NAT Gateway with the specified id
func (c *Client) NATGatewayAddAddress(ctx context.Context, opts NATGatewayAddAddressOptions, id int) (*NATGatewayAddressObject, error) {
	e := formatAPIPath("networking/natgateways/%d/addresses", id)
	return doPOSTRequest[NATGatewayAddressObject](ctx, c, e, opts)
}

// NATGatewayListAddresses lists a NAT Gateway's addresses
func (c *Client) NATGatewayListAddresses(ctx context.Context, opts *ListOptions, id int) ([]NATGatewayAddressObject, error) {
	e := formatAPIPath("networking/natgateways/%d/addresses", id)
	return getPaginatedResults[NATGatewayAddressObject](ctx, c, e, opts)
}

// NATGatewayGetAddress gets the specified address from the NAT Gateway with the provided id
func (c *Client) NATGatewayGetAddress(ctx context.Context, id int, address string) (*NATGatewayAddressObject, error) {
	e := formatAPIPath("networking/natgateways/%d/addresses/%s", id, address)
	return doGETRequest[NATGatewayAddressObject](ctx, c, e)
}

// NATGatewayDeleteAddress deletes the specified Address from the NAT Gateway with the specified id
func (c *Client) NATGatewayDeleteAddress(ctx context.Context, id int, address string) error {
	e := formatAPIPath("networking/natgateways/%d/addresses/%s", id, address)
	return doDELETERequest(ctx, c, e)
}

// NATGatewayListAddressInterfaces lists a NAT Gateway's address' interfaces
func (c *Client) NATGatewayListAddressInterfaces(ctx context.Context, opts *ListOptions, id int, address string) ([]NATGatewayInterface, error) {
	e := formatAPIPath("networking/natgateways/%d/addresses/%s/interfaces", id, address)
	return getPaginatedResults[NATGatewayInterface](ctx, c, e, opts)
}

// NATGatewayListInterfaces lists a NAT Gateway's interfaces
func (c *Client) NATGatewayListInterfaces(ctx context.Context, opts *ListOptions, id int) ([]NATGatewayInterface, error) {
	e := formatAPIPath("networking/natgateways/%d/interfaces", id)
	return getPaginatedResults[NATGatewayInterface](ctx, c, e, opts)
}

// NATGatewayGetTypes gets the pricing for NAT Gateways
func (c *Client) NATGatewayGetTypes(ctx context.Context, opts *ListOptions) ([]NATGatewayType, error) {
	return getPaginatedResults[NATGatewayType](ctx, c, "networking/natgateways/types", opts)
}

// NATGatewayGetSettings gets the user settings and limits for NAT Gateway
func (c *Client) NATGatewayGetSettings(ctx context.Context) (*NATGatewaySettings, error) {
	return doGETRequest[NATGatewaySettings](ctx, c, "networking/natgateways/settings")
}
