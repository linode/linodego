package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestListLongviewSubscriptions(t *testing.T) {
	// Load the mock fixture for Longview subscriptions
	fixtureData, err := fixtures.GetFixture("longview_subscriptions_list")
	assert.NoError(t, err, "Expected no error when getting fixture")

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request for the Longview subscriptions endpoint
	base.MockGet("longview/subscriptions", fixtureData)

	// Call the ListLongviewSubscriptions method
	subscriptions, err := base.Client.ListLongviewSubscriptions(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err, "Expected no error when listing Longview subscriptions")
	assert.NotEmpty(t, subscriptions, "Expected non-empty Longview subscriptions list")

	// Validate the first subscription's details
	assert.Equal(t, "longview-1", subscriptions[0].ID, "Expected subscription ID to match")
	assert.Equal(t, "Longview Pro", subscriptions[0].Label, "Expected subscription label to match")
	assert.Equal(t, 3, subscriptions[0].ClientsIncluded, "Expected clients included to match")
	assert.NotNil(t, subscriptions[0].Price, "Expected price to be non-nil")
	assert.Equal(t, float32(10.00), subscriptions[0].Price.Monthly, "Expected monthly price to match")
	assert.Equal(t, float32(0.01), subscriptions[0].Price.Hourly, "Expected hourly price to match")
}

func TestGetLongviewSubscription(t *testing.T) {
	// Load the mock fixture for a single Longview subscription
	fixtureData, err := fixtures.GetFixture("longview_subscription_get")
	assert.NoError(t, err, "Expected no error when getting fixture")

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request for a single Longview subscription
	subscriptionID := "longview-1"
	base.MockGet("longview/subscriptions/"+subscriptionID, fixtureData)

	// Call the GetLongviewSubscription method
	subscription, err := base.Client.GetLongviewSubscription(context.Background(), subscriptionID)
	assert.NoError(t, err, "Expected no error when getting Longview subscription")
	assert.NotNil(t, subscription, "Expected non-nil Longview subscription")

	// Validate the subscription's details
	assert.Equal(t, "longview-1", subscription.ID, "Expected subscription ID to match")
	assert.Equal(t, "Longview Pro", subscription.Label, "Expected subscription label to match")
	assert.Equal(t, 3, subscription.ClientsIncluded, "Expected clients included to match")
	assert.NotNil(t, subscription.Price, "Expected price to be non-nil")
	assert.Equal(t, float32(10.00), subscription.Price.Monthly, "Expected monthly price to match")
	assert.Equal(t, float32(0.01), subscription.Price.Hourly, "Expected hourly price to match")
}
