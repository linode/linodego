package integration

import (
	"context"
	"fmt"

	. "github.com/linode/linodego"
	"github.com/linode/linodego/pkg/errors"

	"testing"
)

func TestGetIPAddress_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetIPAddress_missing")
	defer teardown()

	doesNotExist := "010.020.030.040"
	i, err := client.GetIPAddress(context.Background(), doesNotExist)
	if err == nil {
		t.Errorf("should have received an error requesting a missing ipaddress, got %v", i)
	}
	e, ok := err.(*errors.Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing ipaddress, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing ipaddress, got %v", e.Code)
	}
}

func TestGetIPAddress_found(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestGetIPAddress_found")
	defer teardown()
	if err != nil {
		t.Errorf("Error creating IPAddress test Instance, got error %v", err)
	}

	address := instance.IPv4[0].String()
	i, err := client.GetIPAddress(context.Background(), address)
	if err != nil {
		t.Errorf("Error getting ipaddress, expected struct, got %v and error %v", i, err)
	}
	if i.Address != address {
		t.Errorf("Expected a specific ipaddress, but got a different one %v", i)
	}
}
func TestListIPAddresses(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestListIPAddresses")
	defer teardown()
	if err != nil {
		t.Errorf("Error creating IPAddress test Instance, got error %v", err)
	}

	filter := fmt.Sprintf("{\"linode_id\":%d}", instance.ID)
	i, err := client.ListIPAddresses(context.Background(), NewListOptions(0, filter))
	if err != nil {
		t.Errorf("Error listing ipaddresses, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of ipaddresses, but got none %v", i)
	}
}

func TestGetInstanceIPAddresses(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestGetInstanceIPAddresses")
	defer teardown()
	if err != nil {
		t.Errorf("Error creating IPAddress test Instance, got error %v", err)
	}

	i, err := client.GetInstanceIPAddresses(context.Background(), instance.ID)
	if err != nil {
		t.Errorf("Error listing ipaddresses, expected struct, got error %v", err)
	}
	if i.IPv4.Public[0].Address != instance.IPv4[0].String() {
		t.Errorf("Expected matching ipaddresses with GetInstanceIPAddress Instance IPAddress but got %v", i)
	}
}

func TestUpdateIPAddress(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestUpdateIPAddress")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	address := instance.IPv4[0].String()
	i, err := client.GetInstanceIPAddresses(context.Background(), instance.ID)
	if err != nil {
		t.Errorf("Error getting ipaddress: %s", err)
	}
	rdns := i.IPv4.Public[0].RDNS

	updateOpts := IPAddressUpdateOptions{
		RDNS: &rdns,
	}

	_, err = client.UpdateIPAddress(context.Background(), address, updateOpts)
	if err != nil {
		t.Error(err)
	}
}
