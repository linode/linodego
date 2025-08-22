package linodego

import (
	"context"
)

// InstanceConfigInterface contains information about a configuration's network interface
type InstanceConfigInterface struct {
	ID          int                          `json:"id"`
	IPAMAddress string                       `json:"ipam_address"`
	Label       string                       `json:"label"`
	Purpose     ConfigInterfacePurpose       `json:"purpose"`
	Primary     bool                         `json:"primary"`
	Active      bool                         `json:"active"`
	VPCID       *int                         `json:"vpc_id"`
	SubnetID    *int                         `json:"subnet_id"`
	IPv4        *VPCIPv4                     `json:"ipv4"`
	IPv6        *InstanceConfigInterfaceIPv6 `json:"ipv6"`
	IPRanges    []string                     `json:"ip_ranges"`
}

// InstanceConfigInterfaceIPv6 represents the IPv6 configuration of a Linode interface.
type InstanceConfigInterfaceIPv6 struct {
	SLAAC    []InstanceConfigInterfaceIPv6SLAAC `json:"slaac"`
	Ranges   []InstanceConfigInterfaceIPv6Range `json:"ranges"`
	IsPublic bool                               `json:"is_public"`
}

// InstanceConfigInterfaceIPv6SLAAC represents a single IPv6 SLAAC under
// a Linode interface.
type InstanceConfigInterfaceIPv6SLAAC struct {
	Range   string `json:"range"`
	Address string `json:"address"`
}

// InstanceConfigInterfaceIPv6Range represents a single IPv6 range under a Linode interface.
type InstanceConfigInterfaceIPv6Range struct {
	Range string `json:"range"`
}

type VPCIPv4 struct {
	VPC     string  `json:"vpc,omitempty"`
	NAT1To1 *string `json:"nat_1_1,omitempty"`
}

type InstanceConfigInterfaceCreateOptions struct {
	IPAMAddress string                                    `json:"ipam_address,omitempty"`
	Label       string                                    `json:"label,omitempty"`
	Purpose     ConfigInterfacePurpose                    `json:"purpose,omitempty"`
	Primary     bool                                      `json:"primary,omitempty"`
	SubnetID    *int                                      `json:"subnet_id,omitempty"`
	IPv4        *VPCIPv4                                  `json:"ipv4,omitempty"`
	IPv6        *InstanceConfigInterfaceCreateOptionsIPv6 `json:"ipv6,omitempty"`
	IPRanges    []string                                  `json:"ip_ranges,omitempty"`
}

// InstanceConfigInterfaceCreateOptionsIPv6 represents the IPv6 configuration of a Linode interface
// specified during creation.
type InstanceConfigInterfaceCreateOptionsIPv6 struct {
	SLAAC    []InstanceConfigInterfaceCreateOptionsIPv6SLAAC `json:"slaac,omitempty"`
	Ranges   []InstanceConfigInterfaceCreateOptionsIPv6Range `json:"ranges,omitempty"`
	IsPublic *bool                                           `json:"is_public,omitempty"`
}

// InstanceConfigInterfaceCreateOptionsIPv6SLAAC represents a single IPv6 SLAAC of a Linode interface
// specified during creation.
type InstanceConfigInterfaceCreateOptionsIPv6SLAAC struct {
	Range string `json:"range"`
}

// InstanceConfigInterfaceCreateOptionsIPv6Range represents a single IPv6 ranges of a Linode interface
// specified during creation.
type InstanceConfigInterfaceCreateOptionsIPv6Range struct {
	Range *string `json:"range,omitempty"`
}

type InstanceConfigInterfaceUpdateOptions struct {
	Primary  bool                                      `json:"primary,omitempty"`
	IPv4     *VPCIPv4                                  `json:"ipv4,omitempty"`
	IPv6     *InstanceConfigInterfaceUpdateOptionsIPv6 `json:"ipv6,omitempty"`
	IPRanges *[]string                                 `json:"ip_ranges,omitempty"`
}

// InstanceConfigInterfaceUpdateOptionsIPv6 represents the IPv6 configuration of a Linode interface
// specified during updates.
type InstanceConfigInterfaceUpdateOptionsIPv6 struct {
	SLAAC    *[]InstanceConfigInterfaceUpdateOptionsIPv6SLAAC `json:"slaac,omitempty"`
	Ranges   *[]InstanceConfigInterfaceUpdateOptionsIPv6Range `json:"ranges,omitempty"`
	IsPublic *bool                                            `json:"is_public,omitempty"`
}

// InstanceConfigInterfaceUpdateOptionsIPv6SLAAC represents a single IPv6 SLAAC of a Linode interface
// specified during updates.
type InstanceConfigInterfaceUpdateOptionsIPv6SLAAC struct {
	Range *string `json:"range,omitempty"`
}

// InstanceConfigInterfaceUpdateOptionsIPv6Range represents a single IPv6 ranges of a Linode interface
// specified during updates.
type InstanceConfigInterfaceUpdateOptionsIPv6Range struct {
	Range *string `json:"range,omitempty"`
}

type InstanceConfigInterfacesReorderOptions struct {
	IDs []int `json:"ids"`
}

func getInstanceConfigInterfacesCreateOptionsList(
	interfaces []InstanceConfigInterface,
) []InstanceConfigInterfaceCreateOptions {
	interfaceOptsList := make([]InstanceConfigInterfaceCreateOptions, len(interfaces))
	for index, configInterface := range interfaces {
		interfaceOptsList[index] = configInterface.GetCreateOptions()
	}

	return interfaceOptsList
}

func (i InstanceConfigInterface) GetCreateOptions() InstanceConfigInterfaceCreateOptions {
	opts := InstanceConfigInterfaceCreateOptions{
		Label:    i.Label,
		Purpose:  i.Purpose,
		Primary:  i.Primary,
		SubnetID: i.SubnetID,
	}

	if len(i.IPRanges) > 0 {
		opts.IPRanges = i.IPRanges
	}

	if i.IPv4 != nil {
		opts.IPv4 = &VPCIPv4{
			VPC:     i.IPv4.VPC,
			NAT1To1: i.IPv4.NAT1To1,
		}
	}

	if i.IPv6 != nil {
		ipv6 := *i.IPv6

		opts.IPv6 = &InstanceConfigInterfaceCreateOptionsIPv6{
			SLAAC: mapSlice(
				ipv6.SLAAC,
				func(i InstanceConfigInterfaceIPv6SLAAC) InstanceConfigInterfaceCreateOptionsIPv6SLAAC {
					return InstanceConfigInterfaceCreateOptionsIPv6SLAAC{
						Range: i.Range,
					}
				},
			),
			Ranges: mapSlice(
				ipv6.Ranges,
				func(i InstanceConfigInterfaceIPv6Range) InstanceConfigInterfaceCreateOptionsIPv6Range {
					return InstanceConfigInterfaceCreateOptionsIPv6Range{
						Range: copyValue(&i.Range),
					}
				},
			),
			IsPublic: copyValue(&ipv6.IsPublic),
		}
	}

	opts.IPAMAddress = i.IPAMAddress

	return opts
}

func (i InstanceConfigInterface) GetUpdateOptions() InstanceConfigInterfaceUpdateOptions {
	opts := InstanceConfigInterfaceUpdateOptions{
		Primary: i.Primary,
	}

	if i.Purpose == InterfacePurposeVPC {
		if i.IPv4 != nil {
			opts.IPv4 = &VPCIPv4{
				VPC:     i.IPv4.VPC,
				NAT1To1: i.IPv4.NAT1To1,
			}
		}

		if i.IPv6 != nil {
			ipv6 := *i.IPv6

			newSLAAC := mapSlice(
				ipv6.SLAAC,
				func(i InstanceConfigInterfaceIPv6SLAAC) InstanceConfigInterfaceUpdateOptionsIPv6SLAAC {
					return InstanceConfigInterfaceUpdateOptionsIPv6SLAAC{
						Range: copyValue(&i.Range),
					}
				},
			)

			newRanges := mapSlice(
				ipv6.Ranges,
				func(i InstanceConfigInterfaceIPv6Range) InstanceConfigInterfaceUpdateOptionsIPv6Range {
					return InstanceConfigInterfaceUpdateOptionsIPv6Range{
						Range: copyValue(&i.Range),
					}
				},
			)

			opts.IPv6 = &InstanceConfigInterfaceUpdateOptionsIPv6{
				SLAAC:    &newSLAAC,
				Ranges:   &newRanges,
				IsPublic: copyValue(&ipv6.IsPublic),
			}
		}
	}

	if i.IPRanges != nil {
		// Copy the slice to prevent accidental
		// mutations
		copiedIPRanges := make([]string, len(i.IPRanges))
		copy(copiedIPRanges, i.IPRanges)

		opts.IPRanges = &copiedIPRanges
	}

	return opts
}

func (c *Client) AppendInstanceConfigInterface(
	ctx context.Context,
	linodeID int,
	configID int,
	opts InstanceConfigInterfaceCreateOptions,
) (*InstanceConfigInterface, error) {
	e := formatAPIPath("/linode/instances/%d/configs/%d/interfaces", linodeID, configID)
	return doPOSTRequest[InstanceConfigInterface](ctx, c, e, opts)
}

func (c *Client) GetInstanceConfigInterface(
	ctx context.Context,
	linodeID int,
	configID int,
	interfaceID int,
) (*InstanceConfigInterface, error) {
	e := formatAPIPath(
		"linode/instances/%d/configs/%d/interfaces/%d",
		linodeID,
		configID,
		interfaceID,
	)

	return doGETRequest[InstanceConfigInterface](ctx, c, e)
}

func (c *Client) ListInstanceConfigInterfaces(
	ctx context.Context,
	linodeID int,
	configID int,
) ([]InstanceConfigInterface, error) {
	e := formatAPIPath(
		"linode/instances/%d/configs/%d/interfaces",
		linodeID,
		configID,
	)

	response, err := doGETRequest[[]InstanceConfigInterface](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return *response, nil
}

func (c *Client) UpdateInstanceConfigInterface(
	ctx context.Context,
	linodeID int,
	configID int,
	interfaceID int,
	opts InstanceConfigInterfaceUpdateOptions,
) (*InstanceConfigInterface, error) {
	e := formatAPIPath(
		"linode/instances/%d/configs/%d/interfaces/%d",
		linodeID,
		configID,
		interfaceID,
	)

	return doPUTRequest[InstanceConfigInterface](ctx, c, e, opts)
}

func (c *Client) DeleteInstanceConfigInterface(
	ctx context.Context,
	linodeID int,
	configID int,
	interfaceID int,
) error {
	e := formatAPIPath(
		"linode/instances/%d/configs/%d/interfaces/%d",
		linodeID,
		configID,
		interfaceID,
	)

	return doDELETERequest(ctx, c, e)
}

func (c *Client) ReorderInstanceConfigInterfaces(
	ctx context.Context,
	linodeID int,
	configID int,
	opts InstanceConfigInterfacesReorderOptions,
) error {
	e := formatAPIPath(
		"linode/instances/%d/configs/%d/interfaces/order",
		linodeID,
		configID,
	)

	return doPOSTRequestNoResponseBody(ctx, c, e, opts)
}
