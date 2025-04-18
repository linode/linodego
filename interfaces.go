package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

type LinodeInterface struct {
	ID           int                    `json:"id"`
	MACAddress   string                 `json:"mac_address"`
	Created      *time.Time             `json:"-"`
	Updated      *time.Time             `json:"-"`
	DefaultRoute InterfacesDefaultRoute `json:"default_route"`
	Version      int                    `json:"version"`
	VPC          InterfacesVPC          `json:"vpc"`
	Public       InterfacesPublic       `json:"public"`
	VLAN         InterfacesVLAN         `json:"vlan"`
}

type InterfacesDefaultRoute struct {
	IPv4 bool `json:"ipv4,omitempty"`
	IPv6 bool `json:"ipv6,omitempty"`
}

type InterfacesVPC struct {
	SubnetID int             `json:"subnet_id"`
	IPv4     *InterfacesIPv4 `json:"ipv4,omitempty"`
}

type InterfacesPublic struct {
	IPv4 *InterfacesIPv4 `json:"ipv4,omitempty"`
	IPv6 *InterfacesIPv6 `json:"ipv6,omitempty"`
}

type InterfacesVLAN struct {
	Label       string  `json:"vlan_label"`
	IPAMAddress *string `json:"ipam_address,omitempty"`
}

type InterfacesIPv4 struct {
	Addresses []*InterfacesAddress `json:"addresses,omitempty"`
}

type InterfacesIPv6 struct {
	Addresses []*InterfacesAddress `json:"addresses,omitempty"`
	Ranges    []*InterfacesRange   `json:"ranges,omitempty"`
}

type InterfacesAddress struct {
	Address *string `json:"address,omitempty"`
	Prefix  *string `json:"prefix,omitempty"`
	Primary bool    `json:"primary,omitempty"`
	NAT1to1 *string `json:"nat_1_1_address,omitempty"`
}

type InterfacesRange struct {
	Range       *string `json:"range,omitempty"`
	RouteTarget *string `json:"route_target,omitempty"`
}

type LinodeInterfaceCreateOptions struct {
	FirewallID   int                     `json:"int,omitempty"`
	DefaultRoute *InterfacesDefaultRoute `json:"default_route,omitempty"`
	VPC          *InterfacesVPC          `json:"vpc,omitempty"`
	Public       *InterfacesPublic       `json:"public,omitempty"`
	VLAN         *InterfacesVLAN         `json:"vlan,omitempty"`
}

type LinodeInterfaceUpdateOptions struct {
	DefaultRoute *InterfacesDefaultRoute `json:"default_route,omitempty"`
	VPC          *InterfacesVPC          `json:"vpc,omitempty"`
	Public       *InterfacesPublic       `json:"public,omitempty"`
	VLAN         *InterfacesVLAN         `json:"vlan,omitempty"`
}

type LinodeInterfacesUpgrade struct {
	ConfigID   int               `json:"config_id,omitempty"`
	DryRun     *bool             `json:"dry_run,omitempty"`
	Interfaces []LinodeInterface `json:"interfaces"`
}

type LinodeInterfacesUpgradeOptions struct {
	ConfigID int  `json:"config_id,omitempty"`
	DryRun   bool `json:"dry_run,omitempty"`
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
