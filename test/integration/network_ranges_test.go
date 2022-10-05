package integration

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/linode/linodego"
)

var testIPv6RangeCreateOptions = IPv6RangeCreateOptions{
	PrefixLength: 64,
}

func TestIPv6Range_Instance_List(t *testing.T) {
	client, ipRange, inst, err := setupIPv6RangeInstance(t, []ipv6RangeModifier{}, "fixtures/TestIPv6Range_Instance_List")
	if err != nil {
		t.Fatal(err)
	}

	filter := Filter{}
	filter.AddField(Eq, "region", inst.Region)
	filterStr, err := filter.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.ListIPv6Ranges(context.Background(), &ListOptions{Filter: string(filterStr)})
	if err != nil {
		t.Errorf("Error listing IPv6 Ranges, expected struct, got error %v", err)
	}

	for _, r := range result {
		if r.RouteTarget == ipRange.RouteTarget && fmt.Sprintf("%s/%d", r.Range, r.Prefix) == ipRange.Range {
			// Ensure GET returns the correct details
			rangeView, err := client.GetIPv6Range(context.Background(), r.Range)
			if err != nil {
				t.Errorf("failed to get ipv6 range: %s", err)
			}

			rangeCommonFields := []IPv6Range{
				{Range: r.Range, Prefix: r.Prefix, Region: r.Region},
				{Range: rangeView.Range, Prefix: rangeView.Prefix, Region: rangeView.Region},
			}

			if !reflect.DeepEqual(rangeCommonFields[0], rangeCommonFields[1]) {
				t.Errorf("ipv6 range view does not match result from list: %s", cmp.Diff(rangeCommonFields[0], rangeCommonFields[1]))
			}

			return
		}
	}

	t.Errorf("failed to find ipv6 range with matching range")
}

func TestIPv6Range_Share(t *testing.T) {
	client, ipRange, origInst, err := setupIPv6RangeInstance(t, []ipv6RangeModifier{}, "fixtures/TestIPv6Range_Share")
	if err != nil {
		t.Fatal(err)
	}

	inst, err := createInstance(t, client, func(inst *InstanceCreateOptions) {
		inst.Label = "go-ins-test-share6"
		inst.Region = origInst.Region
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := client.DeleteInstance(context.Background(), inst.ID); err != nil {
			if t != nil {
				t.Errorf("Error deleting test Instance: %s", err)
			}
		}
	})

	// Share the ip with the new instance
	err = client.ShareIPAddresses(context.Background(), IPAddressesShareOptions{
		LinodeID: origInst.ID,
		IPs: []string{
			strings.TrimSuffix(ipRange.Range, "/64"),
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	err = client.ShareIPAddresses(context.Background(), IPAddressesShareOptions{
		LinodeID: inst.ID,
		IPs: []string{
			strings.TrimSuffix(ipRange.Range, "/64"),
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	ips, err := client.GetInstanceIPAddresses(context.Background(), inst.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(ips.IPv6.Global) < 1 {
		t.Fatal("invalid number of global ipv6 addresses")
	}

	if ips.IPv6.Global[0].Range != strings.TrimSuffix(ipRange.Range, "/64") {
		t.Fatal("ipv6 address does not match")
	}
}

type ipv6RangeModifier func(options *IPv6RangeCreateOptions)

func createIPv6Range(t *testing.T, client *Client, ipv6RangeModifiers ...ipv6RangeModifier) (*IPv6Range, error) {
	t.Helper()

	createOpts := testIPv6RangeCreateOptions
	for _, modifier := range ipv6RangeModifiers {
		modifier(&createOpts)
	}

	ipRange, err := client.CreateIPv6Range(context.Background(), createOpts)
	if err != nil {
		t.Errorf("failed to create ipv6 range: %s", err)
	}

	t.Cleanup(func() {
		rangeSegments := strings.Split(ipRange.Range, "/")

		if err := client.DeleteIPv6Range(context.Background(),
			strings.Join(rangeSegments[:len(rangeSegments)-1], "/")); err != nil {
			t.Errorf("failed to delete ipv6 range: %s", err)
		}
	})

	return ipRange, nil
}

func setupIPv6RangeInstance(t *testing.T, ipv6RangeModifiers []ipv6RangeModifier, fixturesYaml string) (*Client, *IPv6Range, *Instance, error) {
	t.Helper()

	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	t.Cleanup(fixtureTeardown)

	instance, err := createInstance(t, client, func(inst *InstanceCreateOptions) {
		inst.Label = "go-ins-test-range6"
		inst.Region = "eu-west"
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := client.DeleteInstance(context.Background(), instance.ID); err != nil {
			if t != nil {
				t.Errorf("Error deleting test Instance: %s", err)
			}
		}
	})

	ipv6RangeModifiers = append(ipv6RangeModifiers, func(options *IPv6RangeCreateOptions) {
		options.LinodeID = instance.ID
	})

	ipRange, err := createIPv6Range(t, client, ipv6RangeModifiers...)

	return client, ipRange, instance, err
}
