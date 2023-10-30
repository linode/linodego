package integration

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/linode/linodego"
	. "github.com/linode/linodego"
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

func createVPCWithSubnet(t *testing.T, client *linodego.Client, vpcModifier ...vpcModifier) (
	*linodego.VPC,
	*linodego.VPCSubnet,
	func(),
	error,
) {
	t.Helper()
	vpc, vpcTeardown, err := createVPC(t, client, vpcModifier...)
	if err != nil {
		if vpcTeardown != nil {
			vpcTeardown()
		}
		t.Fatal(err)
	}
	createOpts := linodego.VPCSubnetCreateOptions{
		Label: "linodego-vpc-test-" + getUniqueText(),
		IPv4:  TestSubnet,
	}
	vpcSubnet, err := client.CreateVPCSubnet(context.Background(), createOpts, vpc.ID)
	if err != nil {
		vpcTeardown()
		t.Fatal(formatVPCSubnetError(err, "creating", &vpc.ID, nil))
	}

	teardown := func() {
		err = client.DeleteVPCSubnet(context.Background(), vpc.ID, vpcSubnet.ID)
		if err != nil {
			t.Error(formatVPCSubnetError(err, "deleting", &vpc.ID, &vpcSubnet.ID))
		}
		vpcTeardown()
	}
	return vpc, vpcSubnet, teardown, err
}

func setupVPCWithSubnet(
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
	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	vpc, vpcSubnet, vpcSubnetTeardown, err := createVPCWithSubnet(t, client)
	if err != nil {
		if vpcSubnetTeardown != nil {
			vpcSubnetTeardown()
		}
		fixtureTeardown()
		t.Fatal(err)
	}
	teardown := func() {
		vpcSubnetTeardown()
		fixtureTeardown()
	}
	return client, vpc, vpcSubnet, teardown, err
}

func TestVPC_Subnet_Create(t *testing.T) {
	_, _, vpcSubnet, teardown, err := setupVPCWithSubnet(t, "fixtures/TestVPC_Subnet_Create")
	defer teardown()
	if err != nil {
		t.Error(formatVPCSubnetError(err, "setting up", nil, nil))
	}
	vpcSubnetCheck(vpcSubnet, t)
	opts := vpcSubnet.GetCreateOptions()
	vpcSubnetCreateOptionsCheck(&opts, vpcSubnet, t)
}

func TestVPC_Subnet_Update(t *testing.T) {
	client, vpc, vpcSubnet, teardown, err := setupVPCWithSubnet(t, "fixtures/TestVPC_Subnet_Update")
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
	client, vpc, vpcSubnet, teardown, err := setupVPCWithSubnet(t, "fixtures/TestVPC_Subnet_List")
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

func TestVPC_Subnet_Create_Invalid_data(t *testing.T) {
	client, vpc, teardown, err := setupVPC(t, "fixtures/TestVPC_Subnet_Create_Invalid_data")
	defer teardown()
	if err != nil {
		t.Error(formatVPCSubnetError(err, "setting up", nil, nil))
	}

	createOpts := linodego.VPCSubnetCreateOptions{
		Label: "linodego-vpc-test_invalid_label" + getUniqueText(),
		IPv4:  TestSubnet,
	}
	_, err = client.CreateVPCSubnet(context.Background(), createOpts, vpc.ID)
	e, _ := err.(*Error)

	if e.Code != 400 {
		t.Errorf("should have received a 400 Code with invalid label, got %v", e.Code)
	}
	expectedErrorMessage := "Label must include only ASCII letters, numbers, and dashes"
	if !strings.Contains(e.Message, expectedErrorMessage) {
		t.Errorf("Wrong error message displayed should have contained, %s", expectedErrorMessage)
	}
}

func TestVPC_Subnet_Update_Invalid_data(t *testing.T) {
	client, vpc, vpcSubnet, teardown, err := setupVPCWithSubnet(t, "fixtures/TestVPC_Subnet_Update_Invalid_Label")
	defer teardown()
	if err != nil {
		t.Error(formatVPCSubnetError(err, "setting up", nil, nil))
	}
	vpcSubnetCheck(vpcSubnet, t)

	opts := vpcSubnet.GetUpdateOptions()
	vpcSubnetUpdateOptionsCheck(&opts, vpcSubnet, t)

	opts.Label = "invalid_label"
	_, err = client.UpdateVPCSubnet(
		context.Background(),
		vpc.ID,
		vpcSubnet.ID,
		opts,
	)

	e, _ := err.(*Error)

	if e.Code != 400 {
		t.Errorf("should have received a 400 Code with invalid label, got %v", e.Code)
	}
	expectedErrorMessage := "Label must include only ASCII letters, numbers, and dashes"
	if !strings.Contains(e.Message, expectedErrorMessage) {
		t.Errorf("Wrong error message displayed should have contained, %s", expectedErrorMessage)
	}
}

func TestVPC_Subnet_WithInstance(t *testing.T) {
	client, vpc, vpcSubnet, inst, config, teardown := setupInstanceWith3Interfaces(t, "fixtures/TestVPC_Subnet_WithInstance")
	defer teardown()

	// Refresh the subnet to show the assigned instance/interface
	refreshedSubnet, err := client.GetVPCSubnet(context.Background(), vpc.ID, vpcSubnet.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(refreshedSubnet.Linodes) != 1 {
		t.Fatalf("expected 1 assigned linode, got %d", len(refreshedSubnet.Linodes))
	}

	targetLinode := refreshedSubnet.Linodes[0]
	if targetLinode.ID != inst.ID {
		t.Fatalf("expected assigned instance to have id %d, got %d", inst.ID, targetLinode.ID)
	}

	if len(targetLinode.Interfaces) != 1 {
		t.Fatalf("expected 1 assigned interface, got %d", len(targetLinode.Interfaces))
	}

	targetInterface := targetLinode.Interfaces[0]

	if targetInterface.ID != config.Interfaces[2].ID {
		t.Fatalf("interface ID mismatch, expected %d for %d", config.Interfaces[2].ID, targetInterface.ID)
	}

	// Ensure the NAT 1:1 information is reflected in the IP configuration of this instance
	networking, err := client.GetInstanceIPAddresses(context.Background(), inst.ID)
	if err != nil {
		t.Fatal(err)
	}

	nat1To1 := networking.IPv4.Public[0].VPCNAT1To1

	if nat1To1 == nil {
		t.Fatalf("expected VPCNAT1To1 to contain data, got nil")
	}

	if nat1To1.SubnetID != refreshedSubnet.ID {
		t.Fatal("IP/subnet id mismatch")
	}

	if nat1To1.VPCID != vpc.ID {
		t.Fatal("IP/VPC id mismatch")
	}

	if nat1To1.Address != config.Interfaces[2].IPv4.VPC {
		t.Fatalf("nat_1_1 subnet IP mismatch")
	}

}
