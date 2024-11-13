package unit

import (
	"context"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProfileApps_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_apps_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("profile/apps/123", fixtureData)

	app, err := base.Client.GetProfileApp(context.Background(), 123)
	assert.NoError(t, err)

	assert.Equal(t, 123, app.ID)
	assert.Equal(t, "example-app", app.Label)
	assert.Equal(t, "linodes:read_only", app.Scopes)
	assert.Equal(t, "example.org", app.Website)
}

func TestProfileApps_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_apps_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("profile/apps", fixtureData)

	apps, err := base.Client.ListProfileApps(context.Background(), nil)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(apps))
	app := apps[0]

	assert.Equal(t, 123, app.ID)
	assert.Equal(t, "example-app", app.Label)
	assert.Equal(t, "linodes:read_only", app.Scopes)
	assert.Equal(t, "example.org", app.Website)
}

func TestProfileApps_Delete(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, "profile/apps/123"), httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteProfileApp(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}
