package linodego

import (
	"context"
	"fmt"
)

type ReserveIPOptions struct {
	Region string `json:"region"`
}

func (c *Client) GetReservedIPs(ctx context.Context, opts *ListOptions) ([]InstanceIP, error) {
	e := formatAPIPath("networking/reserved/ips")
	response, err := getPaginatedResults[InstanceIP](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	fmt.Println("Reserved IPs Response:")
	for i, ip := range response {
		fmt.Printf("  Address no: %d\n", i+1)
		fmt.Printf("  Address: %s\n", ip.Address)
		fmt.Printf("  Reserved: %s\n", ip.Gateway)
		fmt.Printf("  Subnet Mask: %s\n", ip.SubnetMask)
		fmt.Printf("  Prefix: %d\n", ip.Prefix)
		fmt.Printf("  Type: %s\n", ip.Type)
		fmt.Printf("  Public: %s\n", ip.Public)
		fmt.Printf("  RDNS: %s\n", ip.RDNS)
		fmt.Printf("  linode_id: %d\n", ip.LinodeID)
		fmt.Printf("  Region: %s\n", ip.Region)
		fmt.Printf("  VPC_NAT1To_1: %s\n", ip.VPCNAT1To1)
		fmt.Printf("  Reserved: %d\n", ip.Reserved)

	}

	return response, nil
}

func (c *Client) GetReservedIPAddress(ctx context.Context, id string) (*InstanceIP, error) {
	e := formatAPIPath("networking/reserved/ips/%s", id)
	response, err := doGETRequest[InstanceIP](ctx, c, e)
	if err != nil {
		return nil, err
	}

	fmt.Println("Reserved IP Response:")
	fmt.Printf("  Address: %s\n", response.Address)
	fmt.Printf("  Gateway: %s\n", response.Gateway)
	fmt.Printf("  Subnet Mask: %s\n", response.SubnetMask)
	fmt.Printf("  Prefix: %d\n", response.Prefix)
	fmt.Printf("  Type: %s\n", response.Type)
	fmt.Printf("  Public: %t\n", response.Public) // Assuming Public is a boolean
	fmt.Printf("  RDNS: %s\n", response.RDNS)
	fmt.Printf("  Linode ID: %d\n", response.LinodeID)
	fmt.Printf("  Region: %s\n", response.Region)
	fmt.Printf("  VPC NAT 1-to-1: %t\n", response.VPCNAT1To1) // Assuming VPCNAT1To1 is a boolean
	fmt.Printf("  Reserved: %t\n", response.Reserved)         // Assuming Reserved is a boolean

	return response, nil
}

func (c *Client) ReserveIPAddress(ctx context.Context, opts ReserveIPOptions) (*InstanceIP, error) {
	e := "networking/reserved/ips"
	response, err := doPOSTRequest[InstanceIP](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) DeleteReservedIPAddress(ctx context.Context, ipAddress string) error {
	e := formatAPIPath("networking/reserved/ips/%s", ipAddress)
	err := doDELETERequest(ctx, c, e)
	return err
}
