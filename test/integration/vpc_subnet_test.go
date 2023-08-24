package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/linode/linodego"
)

const (
	TestSubnet = "192.168.0.0/25"
)

func formatVPCSubnetError(err error, action string, vpcID, vpcSubnetID *int) error {
	if err == nil {
		return nil
	}
	vpcMsg := ""
	if vpcID != nil {
		vpcMsg = fmt.Sprintf(" in VPC %v", *vpcID)
	}
	if vpcSubnetID == nil {
		return fmt.Errorf(
			"an error occurs when %v the subnet(s)%v: %v",
			action,
			vpcMsg,
			err,
		)
	}
	return fmt.Errorf(
		"an error occurs when %v the subnet %v%v: %v",
		action,
		*vpcSubnetID,
		vpcMsg,
		err,
	)
}

func vpcSubnetCheck(vpcSubnet *linodego.VPCSubnet, t *testing.T) {
	if vpcSubnet.ID == 0 {
		t.Error("expected a VPC subnet ID, but got 0")
	}
	assertDateSet(t, vpcSubnet.Created)
	assertDateSet(t, vpcSubnet.Updated)
}

func vpcSubnetCreateOptionsCheck(
	opts *linodego.VPCSubnetCreateOptions,
	vpcSubnet *linodego.VPCSubnet,
	t *testing.T,
) {
	if !(opts.IPv4 == vpcSubnet.IPv4 && opts.Label == vpcSubnet.Label) {
		t.Error(
			"the VPC subnet instance and the VPC subnet " +
				"create options instance are mismatched",
		)
	}
}

func vpcSubnetUpdateOptionsCheck(
	opts *linodego.VPCSubnetUpdateOptions,
	vpcSubnet *linodego.VPCSubnet,
	t *testing.T,
) {
	if !(opts.Label == vpcSubnet.Label) {
		t.Error(
			"the VPC subnet instance and the VPC subnet " +
				"update options instance are mismatched",
		)
	}
}

func setupVPCSubnet(
	t *testing.T,
	fixturesYaml string,
) (
	*linodego.Client,
	*linodego.VPC,
	*linodego.VPCSubnet,
	func(),
	error,
) {
	t.Helper()
	client, vpc, fixtureTeardown, err := setupVPC(t, fixturesYaml)
	if err != nil {
		if fixtureTeardown != nil {
			fixtureTeardown()
		}
		t.Fatal(err)
	}
	createOpts := linodego.VPCSubnetCreateOptions{
		Label: "linodego-vpc-test-" + getUniqueText(),
		IPv4:  TestSubnet,
	}
	vpcSubnet, err := client.CreateVPCSubnet(context.Background(), createOpts, vpc.ID)
	if err != nil {
		fixtureTeardown()
		t.Fatal(formatVPCSubnetError(err, "creating", &vpc.ID, nil))
	}

	teardown := func() {
		err = client.DeleteVPCSubnet(context.Background(), vpc.ID, vpcSubnet.ID)
		if err != nil {
			t.Error(formatVPCSubnetError(err, "deleting", &vpc.ID, &vpcSubnet.ID))
		}
		fixtureTeardown()
	}
	return client, vpc, vpcSubnet, teardown, err
}

func TestVPC_Subnet_Create(t *testing.T) {
	_, _, vpcSubnet, teardown, err := setupVPCSubnet(t, "fixtures/TestVPC_Subnet_Create")
	defer teardown()
	if err != nil {
		t.Error(formatVPCSubnetError(err, "setting up", nil, nil))
	}
	vpcSubnetCheck(vpcSubnet, t)
	opts := vpcSubnet.GetCreateOptions()
	vpcSubnetCreateOptionsCheck(&opts, vpcSubnet, t)
}

func TestVPC_Subnet_Update(t *testing.T) {
	client, vpc, vpcSubnet, teardown, err := setupVPCSubnet(t, "fixtures/TestVPC_Subnet_Update")
	defer teardown()
	if err != nil {
		t.Error(formatVPCSubnetError(err, "setting up", nil, nil))
	}
	vpcSubnetCheck(vpcSubnet, t)

	opts := vpcSubnet.GetUpdateOptions()
	vpcSubnetUpdateOptionsCheck(&opts, vpcSubnet, t)

	updatedVPCSubnet, err := client.UpdateVPCSubnet(
		context.Background(),
		vpc.ID,
		vpcSubnet.ID,
		opts,
	)
	if err != nil {
		t.Error(formatVPCSubnetError(err, "updating", &vpc.ID, &vpcSubnet.ID))
	}

	vpcSubnetUpdateOptionsCheck(&opts, updatedVPCSubnet, t)
}

func TestVPC_Subnet_List(t *testing.T) {
	client, vpc, vpcSubnet, teardown, err := setupVPCSubnet(t, "fixtures/TestVPC_Subnet_List")
	defer teardown()
	if err != nil {
		t.Error(formatVPCSubnetError(err, "setting up", nil, nil))
	}
	vpcSubnetCheck(vpcSubnet, t)
	opts := vpcSubnet.GetCreateOptions()
	vpcSubnetCreateOptionsCheck(&opts, vpcSubnet, t)

	vpcSubnets, err := client.ListVPCSubnet(context.Background(), vpc.ID, nil)

	found := false
	for _, v := range vpcSubnets {
		if v.ID == vpcSubnet.ID {
			found = true
		}
	}

	if !found {
		t.Errorf("the VPC %v subnet %v not found in list", vpc.ID, vpcSubnet.ID)
	}
}
