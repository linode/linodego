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
		"label": "Email Alert Channel",
		"channel_type": "email",
		"type": "user",
		"created": "2026-06-23T09:43:00",
		"created_by": "tester",
		"updated": "2026-06-23T09:43:00",
		"updated_by": "tester",
		"details": {
			"email": {
				"recipient_type": "user",
				"usernames": ["tester1", "tester2"]
			}
		},
		"alerts": {
			"alert_count": 0,
			"type": "alert-definitions",
			"url": "/monitor/alert-channels/10000/alerts"
		}
	}`

	monitorAlertChannelCreateWebhookResponse = `{
		"id": 10001,
		"label": "Webhook Alert Channel",
		"channel_type": "webhook",
		"type": "user",
		"created": "2026-06-23T09:43:00",
		"created_by": "tester",
		"updated": "2026-06-23T09:43:00",
		"updated_by": "tester",
		"details": {
			"webhook": {
				"endpoint_url": "https://example.com/test",
				"authentication": {
					"type": "basic"
				},
				"data_compression": "none",
				"client_certificate_details": {
					"tls_hostname": "example.com",
					"client_ca_certificate": "ca-cert",
					"client_certificate": "client-cert"
				},
				"custom_headers": [
					{
						"name": "x-trace-id",
						"value": "1234"
					}
				]
			}
		},
		"alerts": {
			"alert_count": 0,
			"type": "alert-definitions",
			"url": "/monitor/alert-channels/10001/alerts"
		}
	}`

	monitorAlertChannelUpdateResponse = `{
		"id": 10000,
		"label": "Email Alert Channel Updated",
		"channel_type": "email",
		"type": "user",
		"created": "2026-06-23T09:43:00",
		"created_by": "tester",
		"updated": "2026-06-24T09:43:00",
		"updated_by": "tester",
		"details": {
			"email": {
				"recipient_type": "user",
				"usernames": ["tester"]
			}
		},
		"alerts": {
			"alert_count": 0,
			"type": "alert-definitions",
			"url": "/monitor/alert-channels/10000/alerts"
		}
	}`

	monitorAlertChannelListAlertsForChannelResponse = `{
		"data": [
			{
				"id": 10000,
				"label": "Dbaas alert definition",
				"service_type": "dbaas",
				"type": "system",
				"url": "/monitor/services/dbaas/alerts-definitions/10000"
			}
		],
		"page": 1,
		"pages": 1,
		"results": 1
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
		Label:       linodego.Pointer("Email Alert Channel"),
		Details: linodego.AlertChannelDetailsOptions{
			Email: &linodego.EmailChannelCreateOptions{
				Usernames: []string{"tester1", "tester2"},
			},
		},
	}

	channel, err := base.Client.CreateAlertChannel(context.Background(), opts)
	require.NoError(t, err)
	require.NotNil(t, channel)

	assert.Equal(t, 10000, channel.ID)
	assert.Equal(t, "Email Alert Channel", channel.Label)
	assert.Equal(t, linodego.EmailAlertNotification, channel.ChannelType)
	assert.Equal(t, linodego.UserAlertChannel, channel.Type)
	require.NotNil(t, channel.Details.Email)
	assert.Equal(t, []string{"tester1", "tester2"}, channel.Details.Email.Usernames)
	assert.Equal(t, "user", channel.Details.Email.RecipientType)
	assert.Equal(t, 0, channel.Alerts.AlertCount)
	assert.Equal(t, "/monitor/alert-channels/10000/alerts", channel.Alerts.URL)
}

func TestCreateAlertChannelWebhook(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("monitor/alert-channels", json.RawMessage(monitorAlertChannelCreateWebhookResponse))

	authType := linodego.WebhookAuthenticationTypeBasic
	dataCompression := linodego.WebhookDataCompressionNone
	tlsHostname := "example.com"

	opts := linodego.AlertChannelCreateOptions{
		ChannelType: linodego.WebhookAlertNotification,
		Label:       linodego.Pointer("Webhook Alert Channel"),
		Details: linodego.AlertChannelDetailsOptions{
			Webhook: &linodego.WebhookChannelCreateOptions{
				EndpointURL: "https://example.com/test",
				Authentication: &linodego.WebhookChannelAuthenticationCreateOptions{
					Type: &authType,
					Details: &linodego.WebhookChannelAuthenticationBasicDetails{
						BasicAuthenticationUser:     "webhook-user",
						BasicAuthenticationPassword: "webhook-pass",
					},
				},
				DataCompression: &dataCompression,
				ClientCertificateDetails: &linodego.WebhookChannelClientCertificateCreateOptions{
					TLSHostname:         &tlsHostname,
					ClientCACertificate: "ca-cert",
					ClientCertificate:   "client-cert",
					ClientPrivateKey:    "client-private-key",
				},
				CustomHeaders: []linodego.WebhookChannelCustomHeader{
					{
						Name:  "x-trace-id",
						Value: "1234",
					},
				},
			},
		},
	}

	channel, err := base.Client.CreateAlertChannel(context.Background(), opts)
	require.NoError(t, err)
	require.NotNil(t, channel)

	assert.Equal(t, 10001, channel.ID)
	assert.Equal(t, "Webhook Alert Channel", channel.Label)
	assert.Equal(t, linodego.WebhookAlertNotification, channel.ChannelType)
	assert.Equal(t, linodego.UserAlertChannel, channel.Type)
	require.NotNil(t, channel.Details.Webhook)
	assert.Equal(t, "https://example.com/test", channel.Details.Webhook.EndpointURL)
	assert.Equal(t, linodego.WebhookAuthenticationTypeBasic, channel.Details.Webhook.Authentication.Type)
	assert.Equal(t, linodego.WebhookDataCompressionNone, channel.Details.Webhook.DataCompression)
	require.NotNil(t, channel.Details.Webhook.ClientCertificateDetails)
	assert.Equal(t, "example.com", channel.Details.Webhook.ClientCertificateDetails.TLSHostname)
	assert.Equal(t, "ca-cert", channel.Details.Webhook.ClientCertificateDetails.ClientCACertificate)
	assert.Equal(t, "client-cert", channel.Details.Webhook.ClientCertificateDetails.ClientCertificate)
	assert.Equal(t, []linodego.WebhookChannelCustomHeader{{
		Name:  "x-trace-id",
		Value: "1234",
	}}, channel.Details.Webhook.CustomHeaders)
	assert.Equal(t, 0, channel.Alerts.AlertCount)
	assert.Equal(t, "/monitor/alert-channels/10001/alerts", channel.Alerts.URL)
}

func TestDeleteAlertChannel(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("monitor/alert-channels/10000", nil)

	err := base.Client.DeleteAlertChannel(context.Background(), 10000)
	assert.NoError(t, err)
}

func TestUpdateAlertChannel(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("monitor/alert-channels/10000", json.RawMessage(monitorAlertChannelUpdateResponse))

	label := "Email Alert Channel Updated"
	opts := linodego.AlertChannelUpdateOptions{
		Label: &label,
	}

	channel, err := base.Client.UpdateAlertChannel(context.Background(), 10000, opts)
	require.NoError(t, err)
	require.NotNil(t, channel)

	assert.Equal(t, 10000, channel.ID)
	assert.Equal(t, label, channel.Label)
	assert.Equal(t, linodego.EmailAlertNotification, channel.ChannelType)
	assert.Equal(t, linodego.UserAlertChannel, channel.Type)
	require.NotNil(t, channel.Details.Email)
	assert.Equal(t, []string{"tester"}, channel.Details.Email.Usernames)
	assert.Equal(t, "user", channel.Details.Email.RecipientType)
	assert.Equal(t, 0, channel.Alerts.AlertCount)
	assert.Equal(t, "/monitor/alert-channels/10000/alerts", channel.Alerts.URL)
}
func TestGetAlertChannel(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/alert-channels/10000", json.RawMessage(monitorAlertChannelCreateResponse))

	channel, err := base.Client.GetAlertChannel(context.Background(), 10000)
	require.NoError(t, err)
	require.NotNil(t, channel)

	assert.Equal(t, 10000, channel.ID)
	assert.Equal(t, "Email Alert Channel", channel.Label)
	assert.Equal(t, linodego.EmailAlertNotification, channel.ChannelType)
	assert.Equal(t, linodego.UserAlertChannel, channel.Type)
	require.NotNil(t, channel.Details.Email)
	assert.Equal(t, []string{"tester1", "tester2"}, channel.Details.Email.Usernames)
	assert.Equal(t, "user", channel.Details.Email.RecipientType)
	assert.Equal(t, 0, channel.Alerts.AlertCount)
	assert.Equal(t, "/monitor/alert-channels/10000/alerts", channel.Alerts.URL)
}

func TestListAlertsForAlertChannel(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/alert-channels/123/alerts", json.RawMessage(monitorAlertChannelListAlertsForChannelResponse))

	alerts, err := base.Client.ListAlertsForChannel(context.Background(), 123, nil)
	require.NoError(t, err)
	require.Len(t, alerts, 1)

	assert.Equal(t, 10000, alerts[0].ID)
	assert.Equal(t, "dbaas", alerts[0].ServiceType)
	assert.Equal(t, "system", alerts[0].Type)
}

func TestVerifyAlertChannel(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	monitorVerifyResponse := `{
		"webhook": {
			"http_status": 200,
			"response": "OK"
		}
	}`

	base.MockPost("monitor/alert-channels/verify", json.RawMessage(monitorVerifyResponse))

	authType := linodego.WebhookAuthenticationTypeBasic
	dataCompression := linodego.WebhookDataCompressionGZIP
	tlsHostname := "example.com"

	opts := linodego.WebhookChannelCreateOptions{
		EndpointURL: "https://example.com/test",
		Authentication: &linodego.WebhookChannelAuthenticationCreateOptions{
			Type: &authType,
			Details: &linodego.WebhookChannelAuthenticationBasicDetails{
				BasicAuthenticationUser:     "webhook-user",
				BasicAuthenticationPassword: "webhook-pass",
			},
		},
		DataCompression: &dataCompression,
		ClientCertificateDetails: &linodego.WebhookChannelClientCertificateCreateOptions{
			TLSHostname:         &tlsHostname,
			ClientCACertificate: "ca-cert",
			ClientCertificate:   "client-cert",
			ClientPrivateKey:    "client-private-key",
		},
		CustomHeaders: []linodego.WebhookChannelCustomHeader{
			{
				Name:  "x-trace-id",
				Value: "1234",
			},
		},
	}

	resp, err := base.Client.VerifyAlertChannel(context.Background(), opts)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, 200, resp.Webhook.HTTPStatus)
	assert.Equal(t, "OK", resp.Webhook.Response)
}
