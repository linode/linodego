package integration

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/linode/linodego"
)

var testIPv6RangeCreateOptions = IPv6RangeCreateOptions{
	PrefixLength: 64,
}

// TestGetIPv6Range should return an IPv6 Range by id.
func TestListIPv6Range_instance(t *testing.T) {
	client, ipRange, _, teardown, err := setupIPv6RangeInstance(t, []ipv6RangeModifier{}, "fixtures/TestListIPv6Range_instance")
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	result, err := client.ListIPv6Ranges(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing IPv6 Ranges, expected struct, got error %v", err)
	}

	for _, r := range result {
		if r.RouteTarget == ipRange.RouteTarget && fmt.Sprintf("%s/%d", r.Range, r.Prefix) == ipRange.Range {
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

type ipv6RangeModifier func(options *IPv6RangeCreateOptions)

func createIPv6Range(t *testing.T, client *Client, ipv6RangeModifiers ...ipv6RangeModifier) (*IPv6Range, func(), error) {
	t.Helper()

	createOpts := testIPv6RangeCreateOptions
	for _, modifier := range ipv6RangeModifiers {
		modifier(&createOpts)
	}

	ipRange, err := client.CreateIPv6Range(context.Background(), createOpts)
	if err != nil {
		t.Errorf("failed to create ipv6 range: %s", err)
	}

	teardown := func() {
		rangeSegments := strings.Split(ipRange.Range, "/")

		if err := client.DeleteIPv6Range(context.Background(),
			strings.Join(rangeSegments[:len(rangeSegments)-1], "/")); err != nil {
			t.Errorf("failed to delete ipv6 range: %s", err)
		}
	}
	return ipRange, teardown, nil
}

func setupIPv6RangeInstance(t *testing.T, ipv6RangeModifiers []ipv6RangeModifier, fixturesYaml string) (*Client, *IPv6Range, *Instance, func(), error) {
	t.Helper()

	client, instance, instanceTeardown, err := setupInstance(t, fixturesYaml)
	if err != nil {
		t.Fatal(err)
	}

	ipv6RangeModifiers = append(ipv6RangeModifiers, func(options *IPv6RangeCreateOptions) {
		options.LinodeID = instance.ID
	})

	ipRange, rangeTeardown, err := createIPv6Range(t, client, ipv6RangeModifiers...)

	teardown := func() {
		rangeTeardown()
		instanceTeardown()
	}
	return client, ipRange, instance, teardown, err
}
