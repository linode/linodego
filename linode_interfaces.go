package linodego

type LinodeInterface struct {
	FirewallID   *int             `json:"firewall_id"`
	DefaultRoute *DefaultRule     `json:"default_route"`
	VPC          *VPCInterface    `json:"vpc"`
	Public       *PublicInterface `json:"public"`
	VLAN         *VLANInterface   `json:"vlan"`
}

type DefaultRule struct {
	IPv4 *bool `json:"ipv4"`
	IPv6 *bool `json:"ipv6"`
}

type VPCInterface struct {
	SubnetID int                     `json:"subnet_id"`
	IPv4     *LinodeInterfaceVPCIPv4 `json:"ipv4"`
}

type LinodeInterfaceVPCIPv4 struct {
	Addresses []LinodeInterfaceVPCIPv4Address `json:"addresses"`
	Ranges    []LinodeInterfaceIPv6Range      `json:"ranges"`
}

type LinodeInterfaceVPCIPv4Address struct {
	Address      string  `json:"address"`
	Primary      *bool   `json:"primary"`
	NAT11Address *string `json:"nat_1_1_address"`
}

type PublicInterface struct {
	IPv4 *LinodeInterfacePublicIPv4 `json:"ipv4"`
	IPv6 *LinodeInterfacePublicIPv6 `json:"ipv6"`
}

type LinodeInterfacePublicIPv4 struct {
	Addresses []LinodeInterfacePublicIPv4Address `json:"addresses"`
}

type LinodeInterfacePublicIPv6 struct {
	Ranges []LinodeInterfaceIPv6Range `json:"ranges"`
}

type LinodeInterfacePublicIPv4Address struct {
	Address string `json:"address"`
	Primary *bool  `json:"primary"`
}

type LinodeInterfaceIPv6Range struct {
	Range string `json:"range"`
}

type VLANInterface struct {
	VLANLabel   string  `json:"vlan_label"`
	IPAMAddress *string `json:"ipam_address"`
}
