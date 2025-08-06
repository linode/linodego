package integration

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/linode/linodego"
	. "github.com/linode/linodego"
)

type vpcModifier func(*linodego.Client, *linodego.VPCCreateOptions)

func formatVPCError(err error, action string, vpcID *int) error {
	if err == nil {
		return nil
	}
	if vpcID == nil {
		return fmt.Errorf(
			"an error occurs when %v the VPC(s): %v",
			action,
			err,
		)
	}
	return fmt.Errorf(
		"an error occurs when %v the VPC %v: %v",
		action,
		*vpcID,
		err,
	)
}

func createVPC(t *testing.T, client *linodego.Client, vpcModifier ...vpcModifier) (*linodego.VPC, func(), error) {
	t.Helper()
	createOpts := linodego.VPCCreateOptions{
		Label:  "go-test-vpc-" + getUniqueText(),
		Region: getRegionsWithCaps(t, client, []string{"VPCs"})[0],
	}

	for _, mod := range vpcModifier {
		mod(client, &createOpts)
	}

	vpc, err := client.CreateVPC(context.Background(), createOpts)
	if err != nil {
		t.Fatal(formatVPCError(err, "creating", nil))
	}

	teardown := func() {
		if err := client.DeleteVPC(context.Background(), vpc.ID); err != nil {
			t.Error(formatVPCError(err, "deleting", &vpc.ID))
		}
	}
	return vpc, teardown, err
}

func createVPC_invalid_label(t *testing.T, client *linodego.Client) error {
	t.Helper()
	createOpts := linodego.VPCCreateOptions{
		Label:  "gotest_vpc_invalid_label" + getUniqueText(),
		Region: getRegionsWithCaps(t, client, []string{"VPCs"})[0],
	}
	_, err := client.CreateVPC(context.Background(), createOpts)

	return err
}

func setupVPC(t *testing.T, fixturesYaml string) (
	*linodego.Client,
	*linodego.VPC,
	func(),
	error,
) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	vpc, vpcTeardown, err := createVPC(t, client)

	teardown := func() {
		vpcTeardown()
		fixtureTeardown()
	}
	return client, vpc, teardown, err
}

func vpcCheck(vpc *linodego.VPC, t *testing.T) {
	if vpc.ID == 0 {
		t.Errorf("expected a VPC ID, but got 0")
	}
	assertDateSet(t, vpc.Created)
	assertDateSet(t, vpc.Updated)
}

func vpcCreateOptionsCheck(
	opts *linodego.VPCCreateOptions,
	vpc *linodego.VPC,
	t *testing.T,
) {
	good := (opts.Description == vpc.Description &&
		opts.Label == vpc.Label &&
		opts.Region == vpc.Region &&
		len(opts.Subnets) == len(vpc.Subnets))

	for i := 0; i < minInt(len(opts.Subnets), len(vpc.Subnets)); i++ {
		good = good && (opts.Subnets[i].IPv4 == vpc.Subnets[i].IPv4 &&
			opts.Subnets[i].Label == vpc.Subnets[i].Label)
	}

	if !good {
		t.Error(
			"the VPC instance and the VPC creation options instance are mismatched",
		)
	}
}

func vpcUpdateOptionsCheck(
	opts *linodego.VPCUpdateOptions,
	vpc *linodego.VPC,
	t *testing.T,
) {
	if !(opts.Description == vpc.Description && opts.Label == vpc.Label) {
		t.Error("the VPC instance and VPC Update Options instance are mismatched")
	}
}

func TestVPC_CreateGet_smoke(t *testing.T) {
	client, vpc, teardown, err := setupVPC(t, "fixtures/TestVPC_CreateGet")
	defer teardown()
	if err != nil {
		t.Error(formatVPCError(err, "setting up", nil))
	}
	vpcCheck(vpc, t)
	opts := vpc.GetCreateOptions()
	vpcCreateOptionsCheck(&opts, vpc, t)
	client.GetVPC(context.TODO(), vpc.ID)
}

func TestVPC_Update(t *testing.T) {
	client, vpc, teardown, err := setupVPC(t, "fixtures/TestVPC_Update")
	defer teardown()
	if err != nil {
		t.Error(formatVPCError(err, "setting up", nil))
	}
	vpcCheck(vpc, t)

	opts := vpc.GetUpdateOptions()
	vpcUpdateOptionsCheck(&opts, vpc, t)

	updatedDescription := "updated description"
	updatedLabel := "updated-label"

	opts.Description = updatedDescription
	opts.Label = updatedLabel
	updatedVPC, err := client.UpdateVPC(context.Background(), vpc.ID, opts)
	if err != nil {
		t.Error(formatVPCError(err, "updating", &vpc.ID))
	}
	vpcUpdateOptionsCheck(&opts, updatedVPC, t)
}

func TestVPC_List(t *testing.T) {
	client, vpc, teardown, err := setupVPC(t, "fixtures/TestVPC_List")
	defer teardown()
	if err != nil {
		t.Error(formatVPCError(err, "setting up", nil))
	}
	vpcCheck(vpc, t)

	vpcs, err := client.ListVPCs(context.Background(), nil)
	if err != nil {
		t.Error(formatVPCError(err, "listing", nil))
	}

	found := false
	for _, v := range vpcs {
		if v.ID == vpc.ID {
			found = true
		}
	}

	if !found {
		t.Errorf("vpc %v not found in list", vpc.ID)
	}
}

func TestVPC_Create_Invalid_data(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestVPC_Create_Invalid")
	defer teardown()
	err := createVPC_invalid_label(t, client)

	e, _ := err.(*Error)

	if e.Code != 400 {
		t.Errorf("should have received a 400 Code with invalid label, got %v", e.Code)
	}
	expectedErrorMessage := "Label must include only ASCII letters, numbers, and dashes"
	if !strings.Contains(e.Message, expectedErrorMessage) {
		t.Errorf("Wrong error message displayed should have contained, %s", expectedErrorMessage)
	}
}

func TestVPC_Update_Invalid_data(t *testing.T) {
	client, vpc, teardown, err := setupVPC(t, "fixtures/TestVPC_Update_Invalid")
	defer teardown()
	if err != nil {
		t.Error(formatVPCError(err, "setting up", nil))
	}
	vpcCheck(vpc, t)

	opts := vpc.GetUpdateOptions()
	vpcUpdateOptionsCheck(&opts, vpc, t)

	updatedDescription := "updated description"
	updatedLabel := "updated_invalid_label"

	opts.Description = updatedDescription
	opts.Label = updatedLabel

	_, err = client.UpdateVPC(context.Background(), vpc.ID, opts)

	e, _ := err.(*Error)

	if e.Code != 400 {
		t.Errorf("should have received a 400 Code with invalid label, got %v", e.Code)
	}
	expectedErrorMessage := "Label must include only ASCII letters, numbers, and dashes"
	if !strings.Contains(e.Message, expectedErrorMessage) {
		t.Errorf("Wrong error message displayed should have contained, %s", expectedErrorMessage)
	}
}

func TestVPC_ListAllIPAddresses(t *testing.T) {
	client, _, _, instance, config, teardown := setupInstanceWithVPCAndNATOneToOne(
		t, "fixtures/TestVPC_ListAllIPAddresses",
	)
	defer teardown()

	vpcIPs, err := client.ListAllVPCIPAddresses(
		context.Background(),
		linodego.NewListOptions(1, fmt.Sprintf("{\"linode_id\": %d}", instance.ID)),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(vpcIPs) == 0 {
		t.Fatal("expecting 1 VPC IP address, but got 0")
	}

	if *vpcIPs[0].Address != config.Interfaces[0].IPv4.VPC {
		t.Fatalf(
			"expecting VPC IP address on Linode %d to be %q, but got %q",
			instance.ID, *vpcIPs[0].Address, config.Interfaces[0].IPv4.VPC,
		)
	}
}

func TestVPC_ListIPAddresses(t *testing.T) {
	client, vpc, _, instance, config, teardown := setupInstanceWithVPCAndNATOneToOne(
		t, "fixtures/TestVPC_ListIPAddresses",
	)
	defer teardown()

	vpcIPs, err := client.ListVPCIPAddresses(
		context.Background(),
		vpc.ID,
		linodego.NewListOptions(1, fmt.Sprintf("{\"linode_id\": %d}", instance.ID)),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(vpcIPs) == 0 {
		t.Fatal("expecting 1 VPC IP address, but got 0")
	}

	if *vpcIPs[0].Address != config.Interfaces[0].IPv4.VPC {
		t.Fatalf(
			"expecting VPC IP address on Linode %d to be %q, but got %q",
			instance.ID, *vpcIPs[0].Address, config.Interfaces[0].IPv4.VPC,
		)
	}
}

func TestVPC_ListAllIPv6Addresses(t *testing.T) {
	client, vpc, _, instance, config, teardown := setupInstanceWithDualStackVPCAndNAT11(
		t, "fixtures/TestVPC_ListAllIPv6Addresses",
	)
	defer teardown()

	vpcIPs, err := client.ListAllVPCIPv6Addresses(
		context.Background(),
		linodego.NewListOptions(1, fmt.Sprintf("{\"linode_id\": %d}", instance.ID)),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(vpcIPs) == 0 {
		t.Fatal("expecting 1 VPC IP address, but got 0")
	}

	require.NotNil(t, vpc.IPv6[0].Range)
	require.Equal(t, *vpcIPs[0].IPv6Range, config.Interfaces[0].IPv6.SLAAC[0].Range)
	require.Equal(t, vpcIPs[0].IPv6Addresses[0].SLAACAddress, config.Interfaces[0].IPv6.SLAAC[0].Address)
	require.True(t, *vpcIPs[0].IPv6IsPublic)
}

func TestVPC_ListIPv6Addresses(t *testing.T) {
	client, vpc, _, _, config, teardown := setupInstanceWithDualStackVPCAndNAT11(
		t, "fixtures/TestVPC_ListIPv6Addresses",
	)
	defer teardown()

	vpcIPs, err := client.ListVPCIPv6Addresses(
		context.Background(),
		vpc.ID,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(vpcIPs) == 0 {
		t.Fatal("expecting 1 VPC IP address, but got 0")
	}

	require.NotNil(t, vpc.IPv6[0].Range)
	require.Equal(t, *vpcIPs[0].IPv6Range, config.Interfaces[0].IPv6.SLAAC[0].Range)
	require.Equal(t, vpcIPs[0].IPv6Addresses[0].SLAACAddress, config.Interfaces[0].IPv6.SLAAC[0].Address)
	require.True(t, *vpcIPs[0].IPv6IsPublic)
}
