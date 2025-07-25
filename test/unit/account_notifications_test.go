package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestAccountNotifications_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_notifications_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/notifications", fixtureData)

	notifications, err := base.Client.ListNotifications(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.NotEmpty(t, notifications, "Expected notifications to be returned.")

	// Assertions for the first notification in the list
	notification := notifications[0]
	assert.Equal(t, "You have an important ticket open!", notification.Label, "Expected notification label to be 'You have an important ticket open!'")
	assert.Equal(t, "You have an important ticket open!", notification.Message, "Expected notification message to be 'You have an important ticket open!'")
	assert.Equal(t, linodego.NotificationSeverity("major"), notification.Severity, "Expected notification severity to be 'major'")
	assert.Equal(t, linodego.NotificationType("ticket_important"), notification.Type, "Expected notification type to be 'ticket_important'")

	// Validate entity within notification
	assert.Equal(t, 3456, notification.Entity.ID, "Expected ticket ID to be 3456.")
	assert.Equal(t, "Linode not booting.", notification.Entity.Label, "Expected entity label to be 'Linode not booting.'")
	assert.Equal(t, "ticket", notification.Entity.Type, "Expected entity type to be 'ticket'.")
	assert.Equal(t, "/support/tickets/3456", notification.Entity.URL, "Expected entity URL to be '/support/tickets/3456'")
}
