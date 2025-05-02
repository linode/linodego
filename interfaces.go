package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

type LinodeInterface struct {
	ID           int                    `json:"id"`
	Version      int                    `json:"version"`
	MACAddress   string                 `json:"mac_address"`
	Created      *time.Time             `json:"-"`
	Updated      *time.Time             `json:"-"`
	DefaultRoute InterfacesDefaultRoute `json:"default_route"`
	Public       *InterfacesPublic      `json:"public,omitempty"`
	VPC          *InterfacesVPC         `json:"vpc,omitempty"`
	VLAN         *InterfacesVLAN        `json:"vlan,omitempty"`
}

type InterfacesDefaultRoute struct {
	IPv4 *bool `json:"ipv4,omitempty"`
	IPv6 *bool `json:"ipv6,omitempty"`
}

type InterfacesPublic struct {
	IPv4 *InterfacesIPv4 `json:"ipv4,omitempty"`
	IPv6 *InterfacesIPv6 `json:"ipv6,omitempty"`
}

type InterfacesVPC struct {
	VPCID    int             `json:"vpc_id"`
	SubnetID int             `json:"subnet_id"`
	IPv4     *InterfacesIPv4 `json:"ipv4,omitempty"`
}

type InterfacesVLAN struct {
	Label       string  `json:"vlan_label"`
	IPAMAddress *string `json:"ipam_address,omitempty"`
}

type InterfacesIPv4 struct {
	Addresses []*InterfacesIPv4Address `json:"addresses,omitempty"`
	Shared    []*InterfacesIPv4Shared  `json:"shared,omitempty"`
}

type InterfacesIPv6 struct {
	Ranges []*InterfacesIPv6Range `json:"ranges,omitempty"`
	Shared []*InterfacesIPv6Range `json:"shared,omitempty"`
	Slaac  []*InterfacesIPv6Slaac `json:"slaac,omitempty"`
}

type InterfacesIPv4Address struct {
	Address *string `json:"address,omitempty"`
	Primary *bool   `json:"primary,omitempty"`
}

type InterfacesIPv4Shared struct {
	Address  *string `json:"address,omitempty"`
	LinodeID *string `json:"linode_id,omitempty"`
}

type InterfacesIPv6Range struct {
	Range       *string `json:"range,omitempty"`
	RouteTarget *string `json:"route_target,omitempty"`
}

type InterfacesIPv6Slaac struct {
	Prefix  *int    `json:"prefix,omitempty"`
	Address *string `json:"address,omitempty"`
}

type LinodeInterfaceCreateOptions struct {
	FirewallID   *int                    `json:"firewall_id,omitempty"`
	DefaultRoute *InterfacesDefaultRoute `json:"default_route,omitempty"`
	Public       *InterfacesCreatePublic `json:"public,omitempty"`
	VPC          *InterfacesCreateVPC    `json:"vpc,omitempty"`
	VLAN         *InterfacesCreateVLAN   `json:"vlan,omitempty"`
}

type InterfacesCreatePublic struct {
	IPv4 []*InterfacesIPv4       `json:"ipv4,omitempty"`
	IPv6 []*InterfacesIPv6Create `json:"ipv6,omitempty"`
}

type InterfacesIPv6Create struct {
	Ranges []*InterfacesIPv6CreateRange `json:"ranges,omitempty"`
}

type InterfacesIPv6CreateRange struct {
	Range string `json:"range"`
}

type InterfacesCreateVLAN struct {
	IpamAddress string `json:"ipam_address"`
	VLANLabel   string `json:"vlan_label"`
}

type InterfacesCreateVPC struct {
	IPv4     []*InterfacesIPv4 `json:"ipv4,omitempty"`
	SubnetID int               `json:"subnet_id"`
}

type LinodeInterfaceUpdateOptions struct {
	DefaultRoute *InterfacesDefaultRoute `json:"default_route,omitempty"`
	Public       *InterfacesCreatePublic `json:"public,omitempty"`
	VPC          *InterfacesCreateVPC    `json:"vpc,omitempty"`
	VLAN         *InterfacesCreateVLAN   `json:"vlan,omitempty"`
}

type LinodeInterfacesUpgrade struct {
	ConfigID   int               `json:"config_id,omitempty"`
	DryRun     *bool             `json:"dry_run,omitempty"`
	Interfaces []LinodeInterface `json:"interfaces"`
}

type LinodeInterfacesUpgradeOptions struct {
	ConfigID int   `json:"config_id,omitempty"`
	DryRun   *bool `json:"dry_run,omitempty"`
}

type InterfaceSettings struct {
	NetworkHelper bool                  `json:"network_helper"`
	DefaultRoute  *SettingsDefaultRoute `json:"default_route,omitempty"`
}

type InterfaceSettingsUpdateOptions struct {
	NetworkHelper bool                  `json:"network_helper"`
	DefaultRoute  *SettingsDefaultRoute `json:"default_route,omitempty"`
}

type SettingsDefaultRoute struct {
	IPv4InterfaceID          int   `json:"ipv4_interface_id,omitempty"`
	IPv4EligibleInterfaceIDs []int `json:"ipv4_eligible_interface_ids,omitempty"`
	IPv6InterfaceID          int   `json:"ipv6_interface_id,omitempty"`
	IPv6EligibleInterfaceIDs []int `json:"ipv6_eligible_interface_ids,omitempty"`
}

func (i *LinodeInterface) UnmarshalJSON(b []byte) error {
	type Mask LinodeInterface

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Created = (*time.Time)(p.Created)
	i.Updated = (*time.Time)(p.Updated)

	return nil
}

func (c *Client) ListInterfaces(ctx context.Context, linodeID int, opts *ListOptions) ([]LinodeInterface, error) {
	e := formatAPIPath("linode/instances/%d/interfaces", linodeID)
	return getPaginatedResults[LinodeInterface](ctx, c, e, opts)
}

func (c *Client) GetInterface(ctx context.Context, linodeID int, interfaceID int) (*LinodeInterface, error) {
	e := formatAPIPath("linode/instances/%d/interfaces/%d", linodeID, interfaceID)
	return doGETRequest[LinodeInterface](ctx, c, e)
}

func (c *Client) CreateInterface(ctx context.Context, linodeID int, opts LinodeInterfaceCreateOptions) (*LinodeInterface, error) {
	e := formatAPIPath("linode/instances/%d/interfaces", linodeID)
	return doPOSTRequest[LinodeInterface](ctx, c, e, opts)
}

func (c *Client) UpdateInterface(ctx context.Context, linodeID int, interfaceID int, opts LinodeInterfaceUpdateOptions) (*LinodeInterface, error) {
	e := formatAPIPath("linode/instances/%d/interfaces/%d", linodeID, interfaceID)
	return doPUTRequest[LinodeInterface](ctx, c, e, opts)
}

func (c *Client) DeleteInterface(ctx context.Context, linodeID int, interfaceID int) error {
	e := formatAPIPath("linode/instances/%d/interfaces/%d", linodeID, interfaceID)
	return doDELETERequest(ctx, c, e)
}

func (c *Client) UpgradeInterfaces(ctx context.Context, linodeID int, opts LinodeInterfacesUpgradeOptions) (*LinodeInterfacesUpgrade, error) {
	e := formatAPIPath("linode/instances/%d/upgrade-interfaces", linodeID)
	return doPOSTRequest[LinodeInterfacesUpgrade](ctx, c, e, opts)
}

func (c *Client) ListInterfaceFirewalls(ctx context.Context, linodeID int, interfaceID int, opts *ListOptions) ([]Firewall, error) {
	e := formatAPIPath("linode/instances/%d/interfaces/%d/firewalls", linodeID, interfaceID)
	return getPaginatedResults[Firewall](ctx, c, e, opts)
}

func (c *Client) GetInterfaceSettings(ctx context.Context, linodeID int) (*InterfaceSettings, error) {
	e := formatAPIPath("linode/instances/%d/interfaces/settings", linodeID)
	return doGETRequest[InterfaceSettings](ctx, c, e)
}

func (c *Client) UpdateInterfaceSettings(ctx context.Context, linodeID int, opts InterfaceSettingsUpdateOptions) (*InterfaceSettings, error) {
	e := formatAPIPath("linode/instances/%d/interfaces/settings", linodeID)
	return doPUTRequest[InterfaceSettings](ctx, c, e, opts)
}
