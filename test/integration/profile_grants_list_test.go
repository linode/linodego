package integration

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"reflect"
	"testing"
)

func TestGrantsList(t *testing.T) {
	//username := usernamePrefix + "grantslist"
	client := createMockClient(t)
	accessLevel := linodego.AccessLevelReadOnly
	desiredResponse := linodego.GrantsListResponse{
		Database: []linodego.GrantedEntity{
			{
				ID:          1,
				Label:       "example-entity-label",
				Permissions: "read_write",
			},
		},
		Domain: []linodego.GrantedEntity{
			{
				ID:          1,
				Label:       "example-entity-label",
				Permissions: "read_only",
			},
		},
		Firewall: []linodego.GrantedEntity{
			{
				ID:          1,
				Label:       "example-entity-label",
				Permissions: "read_only",
			},
		},
		Image: []linodego.GrantedEntity{
			{
				ID:          1,
				Label:       "example-entity-label",
				Permissions: "read_only",
			},
		},
		Linode: []linodego.GrantedEntity{
			{
				ID:          1,
				Label:       "example-entity-label",
				Permissions: "read_write",
			},
		},
		Longview: []linodego.GrantedEntity{
			{
				ID:          1,
				Label:       "example-entity-label",
				Permissions: "read_write",
			},
		},
		NodeBalancer: []linodego.GrantedEntity{
			{
				ID:          1,
				Label:       "example-entity-label",
				Permissions: "read_only",
			},
		},
		StackScript: []linodego.GrantedEntity{
			{
				ID:          1,
				Label:       "example-entity-label",
				Permissions: "read_only",
			},
		},
		Volume: []linodego.GrantedEntity{
			{
				ID:          1,
				Label:       "example-entity-label",
				Permissions: "read_only",
			},
		},

		Global: linodego.GlobalUserGrants{
			AccountAccess:        &accessLevel,
			AddDomains:           false,
			AddDatabases:         true,
			AddFirewalls:         false,
			AddImages:            true,
			AddLinodes:           true,
			AddLongview:          true,
			AddNodeBalancers:     true,
			AddStackScripts:      true,
			AddVolumes:           false,
			CancelAccount:        false,
			LongviewSubscription: true,
		},
	}

	httpmock.RegisterRegexpResponder(
		"GET",
		mockRequestURL(t, "/profile/grants"),
		httpmock.NewJsonResponderOrPanic(200, &desiredResponse),
	)
	grants, err := client.GrantsList(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*grants, desiredResponse) {
		t.Fatalf(
			"actual response does not equal desired response: %s",
			cmp.Diff(grants, desiredResponse),
		)
	}
}
