package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/linode/linodego"
)

func TestListLongviewClients(t *testing.T) {
	// Load the mock fixture for Longview clients
	fixtureData, err := fixtures.GetFixture("longview_clients_list")
	assert.NoError(t, err, "Expected no error when getting fixture")

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request for the longview clients endpoint
	base.MockGet("longview/clients", fixtureData)

	// Call the ListLongviewClients method
	clients, err := base.Client.ListLongviewClients(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err, "Expected no error when listing longview clients")
	assert.NotEmpty(t, clients, "Expected non-empty longview clients list")

	// Validate the first longview client details
	assert.Equal(t, 123, clients[0].ID, "Expected client ID to match")
	assert.Equal(t, "apache_client", clients[0].Label, "Expected client label to match")
	assert.Equal(t, "API_KEY_123", clients[0].APIKey, "Expected API key to match")
	assert.Equal(t, "install_code_123", clients[0].InstallCode, "Expected install code to match")

	// Validate the Apps field
	assert.NotNil(t, clients[0].Apps.Apache, "Expected apache app to be non-nil")
	assert.NotNil(t, clients[0].Apps.MySQL, "Expected mysql app to be non-nil")
	assert.NotNil(t, clients[0].Apps.NginX, "Expected nginx app to be non-nil")

	// Validate the created and updated time for the first client
	expectedCreatedTime, err := time.Parse(time.RFC3339, "2025-01-23T00:00:00Z")
	assert.NoError(t, err, "Expected no error when parsing created time")
	assert.Equal(t, expectedCreatedTime, *clients[0].Created, "Expected created time to match")

	expectedUpdatedTime, err := time.Parse(time.RFC3339, "2025-01-23T00:00:00Z")
	assert.NoError(t, err, "Expected no error when parsing updated time")
	assert.Equal(t, expectedUpdatedTime, *clients[0].Updated, "Expected updated time to match")
}

func TestGetLongviewClient(t *testing.T) {
	// Load the mock fixture for a single longview client
	fixtureData, err := fixtures.GetFixture("longview_client_single")
	assert.NoError(t, err, "Expected no error when getting fixture")

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request for a single longview client
	base.MockGet("longview/clients/123", fixtureData)

	// Call the GetLongviewClient method
	client, err := base.Client.GetLongviewClient(context.Background(), 123)
	assert.NoError(t, err, "Expected no error when getting longview client")
	assert.NotNil(t, client, "Expected non-nil longview client")

	// Validate the client details
	assert.Equal(t, 123, client.ID, "Expected client ID to match")
	assert.Equal(t, "apache_client", client.Label, "Expected client label to match")
	assert.Equal(t, "API_KEY_123", client.APIKey, "Expected API key to match")
	assert.Equal(t, "install_code_123", client.InstallCode, "Expected install code to match")

	// Validate the Apps field
	assert.NotNil(t, client.Apps.Apache, "Expected apache app to be non-nil")
	assert.NotNil(t, client.Apps.MySQL, "Expected mysql app to be non-nil")
	assert.NotNil(t, client.Apps.NginX, "Expected nginx app to be non-nil")

	// Validate the created and updated time for the client
	expectedCreatedTime, err := time.Parse(time.RFC3339, "2025-01-23T00:00:00Z")
	assert.NoError(t, err, "Expected no error when parsing created time")
	assert.Equal(t, expectedCreatedTime, *client.Created, "Expected created time to match")

	expectedUpdatedTime, err := time.Parse(time.RFC3339, "2025-01-23T00:00:00Z")
	assert.NoError(t, err, "Expected no error when parsing updated time")
	assert.Equal(t, expectedUpdatedTime, *client.Updated, "Expected updated time to match")
}

func TestGetLongviewPlan(t *testing.T) {
	// Load the mock fixture for Longview plan
	fixtureData, err := fixtures.GetFixture("longview_plan")
	assert.NoError(t, err, "Expected no error when getting fixture")

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request for the longview plan endpoint
	base.MockGet("longview/plan", fixtureData)

	// Call the GetLongviewPlan method
	plan, err := base.Client.GetLongviewPlan(context.Background())
	assert.NoError(t, err, "Expected no error when getting longview plan")
	assert.NotNil(t, plan, "Expected non-nil longview plan")

	// Validate the plan details
	assert.Equal(t, "longview-plan-id", plan.ID, "Expected plan ID to match")
	assert.Equal(t, "Longview Plan", plan.Label, "Expected plan label to match")
	assert.Equal(t, 5, plan.ClientsIncluded, "Expected number of clients included to match")
	assert.Equal(t, 50.00, plan.Price.Hourly, "Expected hourly price to match")
	assert.Equal(t, 500.00, plan.Price.Monthly, "Expected monthly price to match")
}
