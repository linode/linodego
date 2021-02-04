package integration

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/linode/linodego"
)

func TestUpdateUserGrants(t *testing.T) {
	username := usernamePrefix + "updateusergrants"

	client, _, teardown := setupUser(t, []userModifier{
		func(createOpts *linodego.UserCreateOptions) {
			createOpts.Username = username
			createOpts.Email = usernamePrefix + "updateusergrants@example.com"
			createOpts.Restricted = true
		},
	}, "fixtures/TestUpdateUserGrants")
	defer teardown()

	globalGrants := linodego.GlobalUserGrants{
		AccountAccess:    linodego.AccessLevelReadOnly,
		AddDomains:       false,
		AddImages:        true,
		AddLinodes:       false,
		AddLongview:      true,
		AddNodeBalancers: false,
		AddStackScripts:  true,
		AddVolumes:       true,
		CancelAccount:    false,
	}

	expectedUserGrants := linodego.UserGrants{
		Global:       globalGrants,
		Domain:       []linodego.GrantedEntity{},
		Image:        []linodego.GrantedEntity{},
		Linode:       []linodego.GrantedEntity{},
		Longview:     []linodego.GrantedEntity{},
		NodeBalancer: []linodego.GrantedEntity{},
		StackScript:  []linodego.GrantedEntity{},
		Volume:       []linodego.GrantedEntity{},
	}
	grants, err := client.UpdateUserGrants(context.TODO(), username, linodego.UserGrantsUpdateOptions{
		Global: globalGrants,
	})
	if err != nil {
		t.Fatalf("failed to get user grants: %s", err)
	}

	if !cmp.Equal(grants, &expectedUserGrants) {
		t.Errorf("expected rules to match test rules, but got diff: %s", cmp.Diff(grants, &expectedUserGrants))
	}
}
