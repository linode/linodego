package unit

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	monitorAlertChannelListResponse = `{
		"data": [{
			"id": 123,
			"label": "alert notification channel",
			"channel_type": "email",
			"type": "user",
			"details": {
				"email": {
					"usernames": [
						"admin-user1",
						"admin-user2"
					],
					"recipient_type": "user"
				}
			},
			"alerts": {
				"url": "/monitor/alert-channels/123/alerts",
				"type": "alerts-definitions",
				"alert_count": 0
			},
			"created": "2024-01-01T00:00:00",
			"updated": "2024-01-01T00:00:00",
			"created_by": "tester",
			"updated_by": "tester"
		}],
		"page": 1,
		"pages": 1,
		"results": 1
	}`

	monitorAlertChannelCreateResponse = `{
		"id": 10000,
		"label": "My Email Alert Channel",
		"channel_type": "email",
		"type": "user",
		"created": "2026-06-23T09:43:00",
		"created_by": "johndoe",
		"updated": "2026-06-23T09:43:00",
		"updated_by": "johndoe",
		"details": {
			"email": {
				"recipient_type": "user",
				"usernames": ["johndoe", "janedoe"]
			}
		},
		"alerts": {
			"alert_count": 0,
			"type": "alert-definitions",
			"url": "/monitor/alert-channels/10000/alerts"
		}
	}`
)

func TestListAlertChannels(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/alert-channels", json.RawMessage(monitorAlertChannelListResponse))

	channels, err := base.Client.ListAlertChannels(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, channels, 1)

	channel := channels[0]
	assert.Equal(t, 123, channel.ID)
	assert.Equal(t, "alert notification channel", channel.Label)
	assert.Equal(t, linodego.EmailAlertNotification, channel.ChannelType)
	assert.Equal(t, linodego.UserAlertChannel, channel.Type)
	require.NotNil(t, channel.Details.Email)
	assert.Equal(t, []string{"admin-user1", "admin-user2"}, channel.Details.Email.Usernames)
	assert.Equal(t, "user", channel.Details.Email.RecipientType)
	assert.Equal(t, 0, channel.Alerts.AlertCount)
	assert.Equal(t, "/monitor/alert-channels/123/alerts", channel.Alerts.URL)
}

func TestCreateAlertChannel(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("monitor/alert-channels", json.RawMessage(monitorAlertChannelCreateResponse))

	opts := linodego.AlertChannelCreateOptions{
		ChannelType: linodego.EmailAlertNotification,
		Label:       linodego.Pointer("My Email Alert Channel"),
		Details: linodego.AlertChannelDetailsOptions{
			Email: &linodego.EmailChannelCreateOptions{
				Usernames: []string{"johndoe", "janedoe"},
			},
		},
	}

	channel, err := base.Client.CreateAlertChannel(context.Background(), opts)
	require.NoError(t, err)
	require.NotNil(t, channel)

	assert.Equal(t, 10000, channel.ID)
	assert.Equal(t, "My Email Alert Channel", channel.Label)
	assert.Equal(t, linodego.EmailAlertNotification, channel.ChannelType)
	assert.Equal(t, linodego.UserAlertChannel, channel.Type)
	require.NotNil(t, channel.Details.Email)
	assert.Equal(t, []string{"johndoe", "janedoe"}, channel.Details.Email.Usernames)
	assert.Equal(t, "user", channel.Details.Email.RecipientType)
	assert.Equal(t, 0, channel.Alerts.AlertCount)
	assert.Equal(t, "/monitor/alert-channels/10000/alerts", channel.Alerts.URL)
}
