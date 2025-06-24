package integration

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/linode/linodego"
)

func TestUserGrants_Update(t *testing.T) {
	username := usernamePrefix + "updateusergrants"

	client, _, teardown := setupUser(t, []userModifier{
		func(createOpts *linodego.UserCreateOptions) {
			createOpts.Username = username
			createOpts.Email = usernamePrefix + "updateusergrants@example.com"
			createOpts.Restricted = true
		},
	}, "fixtures/TestUserGrants_Update")
	defer teardown()

	globalGrants := linodego.GlobalUserGrants{
		AccountAccess:      linodego.Pointer(linodego.AccessLevelReadOnly),
		AddDomains:         false,
		AddDatabases:       true,
		AddFirewalls:       true,
		AddImages:          true,
		AddLinodes:         false,
		AddLongview:        true,
		AddNodeBalancers:   false,
		AddPlacementGroups: false,
		AddStackScripts:    true,
		AddVolumes:         true,
		AddVPCs:            true,
		CancelAccount:      false,
	}

	grants, err := client.UpdateUserGrants(context.TODO(), username, linodego.UserGrantsUpdateOptions{
		Global: globalGrants,
	})
	if err != nil {
		t.Fatalf("failed to get user grants: %s", err)
	}

	if !cmp.Equal(grants.Global, globalGrants) {
		t.Errorf("expected rules to match test rules, but got diff: %s", cmp.Diff(grants.Global, globalGrants))
	}
}

func TestUserGrants_UpdateNoAccess(t *testing.T) {
	username := usernamePrefix + "updateusergrantsna"

	client, _, teardown := setupUser(t, []userModifier{
		func(createOpts *linodego.UserCreateOptions) {
			createOpts.Username = username
			createOpts.Email = usernamePrefix + "updateusergrants@example.com"
			createOpts.Restricted = true
		},
	}, "fixtures/TestUserGrants_UpdateNoAccess")
	defer teardown()

	globalGrants := linodego.GlobalUserGrants{
		AccountAccess: nil,
	}

	grants, err := client.UpdateUserGrants(context.TODO(), username, linodego.UserGrantsUpdateOptions{
		Global: globalGrants,
	})
	if err != nil {
		t.Fatalf("failed to get user grants: %s", err)
	}

	if !cmp.Equal(grants.Global, globalGrants) {
		t.Errorf("expected rules to match test rules, but got diff: %s", cmp.Diff(grants.Global, globalGrants))
	}

	// Ensure all grants are no access
	grantFields := [][]linodego.GrantedEntity{
		grants.Domain,
		grants.Firewall,
		grants.Image,
		grants.Linode,
		grants.Longview,
		grants.NodeBalancer,
		grants.PlacementGroup,
		grants.StackScript,
		grants.Volume,
		grants.VPC,
	}

	for _, grantField := range grantFields {
		for _, grant := range grantField {
			if grant.Permissions != "" {
				t.Errorf("expected permissions to be nil, but got %s", grant.Permissions)
			}
		}
	}
}
