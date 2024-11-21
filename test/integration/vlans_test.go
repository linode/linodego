package integration

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/linode/linodego"
)

func TestVLANs_List_smoke(t *testing.T) {
	vlanName := "go-vlan-test-list"
	instancePrefix := "go-ins-test-list"

	client, fixturesTeardown := createTestClient(t, "fixtures/TestVLANs_List")
	defer fixturesTeardown()

	var instances []*linodego.Instance
	for i := 0; i < 2; i++ {
		instance, instanceTeardown, err := createVLANInstance(t, client, fmt.Sprintf("%s-%d", instancePrefix, i), vlanName)
		if err != nil {
			t.Fatal(err)
		}
		defer instanceTeardown()

		instances = append(instances, instance)
	}

	for _, instance := range instances {
		if _, err := client.WaitForInstanceStatus(context.Background(), instance.ID, linodego.InstanceRunning, 240); err != nil {
			t.Error(err)
		}
	}

	vlans, err := client.ListVLANs(context.TODO(), &linodego.ListOptions{
		Filter: fmt.Sprintf("{\"label\": \"%s\"}", vlanName),
	})
	if err != nil {
		t.Error(err)
	}

	if len(vlans) < 1 {
		t.Error("expected list of vlans but got none")
	}

	for _, instance := range instances {
		if !vlanHasLinodeID(&vlans[0], instance.ID) {
			t.Errorf("instance %d not found in vlan", instance.ID)
		}
	}
}

func TestVLANs_GetIPAMAddress(t *testing.T) {
	vlanName := "go-vlan-test-ipam"
	instancePrefix := "go-ins-test-ipam"

	client, fixturesTeardown := createTestClient(t, "fixtures/TestVLANs_GetIPAMAddress")
	defer fixturesTeardown()

	instance, instanceTeardown, err := createVLANInstance(t, client, instancePrefix, vlanName)
	if err != nil {
		t.Fatal(err)
	}
	defer instanceTeardown()

	_, err = client.WaitForInstanceStatus(context.Background(), instance.ID, linodego.InstanceRunning, 240)
	if err != nil {
		t.Error(err)
	}

	ipam, err := client.GetVLANIPAMAddress(context.Background(), instance.ID, vlanName)
	if err != nil {
		t.Error(err)
	}

	r, _ := regexp.Compile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}/[0-9]{2}`)
	if !r.MatchString(ipam) {
		t.Errorf("Result does not match regular expression: %s", ipam)
	}
}

func createVLANInstance(t *testing.T, client *linodego.Client, instanceName, vlanName string) (*linodego.Instance, func(), error) {
	t.Helper()

	trueBool := true

	instance, err := createInstance(t, client, true, func(client *linodego.Client, opts *linodego.InstanceCreateOptions) {
		opts.Interfaces = []linodego.InstanceConfigInterfaceCreateOptions{
			{
				Label:       vlanName,
				Purpose:     linodego.InterfacePurposeVLAN,
				IPAMAddress: "10.0.0.1/24",
			},
		}

		opts.Booted = &trueBool
		opts.Label = instanceName
		opts.Region = getRegionsWithCaps(t, client, []string{"Vlans"})[0]
	})
	if err != nil {
		return nil, nil, err
	}

	teardown := func() {
		if terr := client.DeleteInstance(context.Background(), instance.ID); terr != nil {
			t.Errorf("Error deleting test Instance: %s", terr)
		}
	}

	return instance, teardown, err
}

func vlanHasLinodeID(vlan *linodego.VLAN, linodeID int) bool {
	for _, id := range vlan.Linodes {
		if id == linodeID {
			return true
		}
	}

	return false
}
