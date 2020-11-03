package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

var (
	testVLANCreateOpts = linodego.VLANCreateOptions{
		Description: "linodego-testing",
		Region:      "ca-central",
	}
)

type vlanModifier func(*linodego.VLANCreateOptions)

func TestListVLANs(t *testing.T) {
	client, _, teardown, err := setupVLAN(t, []vlanModifier{}, "fixtures/TestListVLAN")
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	vlans, err := client.ListVLANs(context.TODO(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(vlans) == 0 {
		t.Error("expected list of vlans but got none")
	}
}

func TestGetVLAN(t *testing.T) {
	description := "testing123"
	cidrBlock := "0.0.0.0/0"

	client, created, teardown, err := setupVLAN(t, []vlanModifier{
		func(opts *linodego.VLANCreateOptions) {
			opts.Description = description
			opts.CIDRBlock = cidrBlock
		},
	}, "fixtures/TestGetVLAN")
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	vlan, err := client.GetVLAN(context.TODO(), created.ID)
	if err != nil {
		t.Error(err)
	}

	if len(vlan.Linodes) != 0 {
		t.Error("expected no linodes to be attached to vlan")
	}
	if vlan.CIDRBlock != cidrBlock {
		t.Errorf("expected cidr block to be '%s'; got '%s'", cidrBlock, vlan.CIDRBlock)
	}
	if vlan.Description != description {
		t.Errorf("expected description to be '%s'; got '%s'", description, vlan.Description)
	}
}

func TestAttachVLAN(t *testing.T) {
	client, vlan, teardown, err := setupVLAN(t, []vlanModifier{}, "fixtures/TestAttachVLAN")
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	instance, err := createInstance(t, client, func(opts *linodego.InstanceCreateOptions) {
		opts.Region = "ca-central"
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := client.DeleteInstance(context.TODO(), instance.ID); err != nil {
			t.Error(err)
		}
	}()

	vlan, err = client.AttachVLAN(context.TODO(), vlan.ID, linodego.VLANAttachOptions{
		Linodes: []int{instance.ID},
	})
	if err != nil {
		t.Fatal(err)
	}

	if !(len(vlan.Linodes) == 1 && vlan.Linodes[0] == instance.ID) {
		t.Errorf("expected linode %d to have been attached to the vlan", instance.ID)
	}
}

func TestDetachVLAN(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestDetachVLAN")
	defer fixtureTeardown()

	instance, err := createInstance(t, client, func(opts *linodego.InstanceCreateOptions) {
		opts.Region = "ca-central"
	})
	if err != nil {
		t.Error(err)
	}
	teardownInstance := func() {
		if err := client.DeleteInstance(context.TODO(), instance.ID); err != nil {
			t.Error(err)
		}
	}

	vlan, teardownVlan, err := createVLAN(t, client, func(opts *linodego.VLANCreateOptions) {
		opts.Linodes = []int{instance.ID}
	})
	if err != nil {
		defer teardownInstance()
		t.Error(err)
	}
	defer func() {
		teardownInstance()
		teardownVlan()
	}()

	vlan, err = client.GetVLAN(context.TODO(), vlan.ID)
	if err != nil {
		t.Error(err)
	}

	if !(len(vlan.Linodes) == 1 && vlan.Linodes[0] == instance.ID) {
		t.Errorf("expected linode %d to be attached to vlan", instance.ID)
	}

	if vlan, err = client.DetachVLAN(context.TODO(), vlan.ID, linodego.VLANDetachOptions{
		Linodes: []int{instance.ID},
	}); err != nil {
		t.Error(err)
	}

	if len(vlan.Linodes) != 0 {
		t.Errorf("expect linode %d to be detached from vlan", instance.ID)
	}
}

func createVLAN(t *testing.T, client *linodego.Client, vlanModifiers ...vlanModifier) (*linodego.VLAN, func(), error) {
	t.Helper()

	createOpts := testVLANCreateOpts
	for _, modifier := range vlanModifiers {
		modifier(&createOpts)
	}

	vlan, err := client.CreateVLAN(context.TODO(), createOpts)
	if err != nil {
		t.Errorf("failed to create vlan: %s", err)
	}

	teardown := func() {
		if err := client.DeleteVLAN(context.TODO(), vlan.ID); err != nil {
			t.Errorf("failed to delete vlan: %s", err)
		}
	}
	return vlan, teardown, nil
}

func setupVLAN(t *testing.T, vlanModifiers []vlanModifier, fixturesYaml string) (*linodego.Client, *linodego.VLAN, func(), error) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	vlan, vlanTeardown, err := createVLAN(t, client, vlanModifiers...)

	teardown := func() {
		vlanTeardown()
		fixtureTeardown()
	}
	return client, vlan, teardown, err
}
