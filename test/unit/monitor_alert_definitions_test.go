package unit

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

const (
	testMonitorAlertDefinitionServiceType = "dbaas"
	testMonitorAlertDefinitionID          = 123

	monitorAlertDefinitionGetResponse = `{
		"id": 123,
		"label": "test-alert-definition",
		"severity": 1,
		"type": "some_type",
		"service_type": "dbaas",
		"status": "enabled",
		"entity_ids": ["12345"],
		"channel_ids": [1],
		"is_enabled": true
	}`

	monitorAlertDefinitionListResponse = `{
		"data": [{
			"id": 123,
			"label": "test-alert-definition",
			"severity": 1,
			"type": "some_type",
			"service_type": "dbaas",
			"status": "enabled",
			"entity_ids": ["12345"],
			"channel_ids": [1],
			"is_enabled": true
		}],
		"page": 1,
		"pages": 1,
		"results": 1
	}`

	monitorAlertDefinitionUpdateResponse = `{
		"id": 123,
		"label": "test-alert-definition-renamed",
		"severity": 2,
		"type": "some_type",
		"service_type": "dbaas",
		"status": "disabled",
		"entity_ids": ["12345"],
		"channel_ids": [1, 2],
		"is_enabled": false
	}`

	monitorAlertDefinitionUpdateLabelOnlyResponseSingleLine = `{"id": 123, "label": "test-alert-definition-renamed-one-line", "severity": 1, "type": "some_type", "service_type": "dbaas", "status": "enabled", "entity_ids": ["12345"], "channel_ids": [1], "is_enabled": true}`
)

func TestCreateMonitorAlertDefinition(t *testing.T) {
	var base ClientBaseCase
	defer base.TearDown(t)

	base.MockPost("monitor/services/dbaas/alert-definitions", json.RawMessage(monitorAlertDefinitionGetResponse))

	createOpts := linodego.AlertDefinitionCreateOptions{
		Label:      "test-alert-definition",
		Severity:   int(linodego.SeverityLow),
		ChannelIDs: []int{1},
		EntityIDs:  []string{"12345"},
	}

	alert, err := base.Client.CreateMonitorAlertDefinition(context.Background(), testMonitorAlertDefinitionServiceType, createOpts)
	assert.NoError(t, err)
	assert.NotNil(t, alert)
	assert.Equal(t, "test-alert-definition", alert.Label)
	assert.Equal(t, testMonitorAlertDefinitionID, alert.ID)
}

func TestCreateMonitorAlertDefinitionWithIdempotency(t *testing.T) {
	var base ClientBaseCase
	defer base.TearDown(t)

	base.MockPost("monitor/services/dbaas/alert-definitions", json.RawMessage(monitorAlertDefinitionGetResponse))

	createOpts := linodego.AlertDefinitionCreateOptions{
		Label:      "test-alert-definition",
		Severity:   int(linodego.SeverityLow),
		ChannelIDs: []int{1},
		EntityIDs:  []string{"12345"},
	}

	alert, err := base.Client.CreateMonitorAlertDefinitionWithIdempotency(
		context.Background(),
		testMonitorAlertDefinitionServiceType,
		createOpts,
		"idempotency-key",
	)
	assert.NoError(t, err)
	assert.NotNil(t, alert)
	assert.Equal(t, "test-alert-definition", alert.Label)
	assert.Equal(t, testMonitorAlertDefinitionID, alert.ID)
}

func TestGetMonitorAlertDefinition(t *testing.T) {
	var base ClientBaseCase
	defer base.TearDown(t)

	base.MockGet("monitor/services/dbaas/alert-definitions/123", json.RawMessage(monitorAlertDefinitionGetResponse))

	alert, err := base.Client.GetMonitorAlertDefinition(context.Background(), testMonitorAlertDefinitionServiceType, testMonitorAlertDefinitionID)
	assert.NoError(t, err)
	assert.NotNil(t, alert)
	assert.Equal(t, "test-alert-definition", alert.Label)
	assert.Equal(t, testMonitorAlertDefinitionID, alert.ID)
}

func TestListMonitorAlertDefinitions(t *testing.T) {
	var base ClientBaseCase
	defer base.TearDown(t)

	base.MockGet("monitor/services/dbaas/alert-definitions", json.RawMessage(monitorAlertDefinitionListResponse))

	alerts, err := base.Client.ListMonitorAlertDefinitions(context.Background(), testMonitorAlertDefinitionServiceType, nil)
	assert.NoError(t, err)
	assert.Len(t, alerts, 1)
	assert.Equal(t, "test-alert-definition", alerts[0].Label)
	assert.Equal(t, testMonitorAlertDefinitionID, alerts[0].ID)
}

func TestUpdateMonitorAlertDefinition(t *testing.T) {
	var base ClientBaseCase
	defer base.TearDown(t)

	base.MockPut("monitor/services/dbaas/alert-definitions/123", json.RawMessage(monitorAlertDefinitionUpdateResponse))

	updateOpts := linodego.AlertDefinitionUpdateOptions{
		Label:      "test-alert-definition-renamed",
		Severity:   int(linodego.SeverityLow),
		ChannelIDs: []int{1, 2},
	}

	alert, err := base.Client.UpdateMonitorAlertDefinition(
		context.Background(),
		testMonitorAlertDefinitionServiceType,
		testMonitorAlertDefinitionID,
		updateOpts,
	)
	assert.NoError(t, err)
	assert.NotNil(t, alert)
	assert.Equal(t, "test-alert-definition-renamed", alert.Label)
	assert.Equal(t, 2, alert.Severity)
}

func TestUpdateMonitorAlertDefinition_LabelOnly(t *testing.T) {
	var base ClientBaseCase
	defer base.TearDown(t)

	// Mock a PUT that returns the single-line fixture
	base.MockPut("monitor/services/dbaas/alert-definitions/123", json.RawMessage(monitorAlertDefinitionUpdateLabelOnlyResponseSingleLine))

	updateOpts := linodego.AlertDefinitionUpdateOptions{
		Label: "test-alert-definition-renamed-one-line",
	}

	alert, err := base.Client.UpdateMonitorAlertDefinition(
		context.Background(),
		testMonitorAlertDefinitionServiceType,
		testMonitorAlertDefinitionID,
		updateOpts,
	)
	assert.NoError(t, err)
	assert.NotNil(t, alert)
	assert.Equal(t, "test-alert-definition-renamed-one-line", alert.Label)
	assert.Equal(t, testMonitorAlertDefinitionID, alert.ID)
}

func TestDeleteMonitorAlertDefinition(t *testing.T) {
	var base ClientBaseCase
	defer base.TearDown(t)

	base.MockDelete("monitor/services/dbaas/alert-definitions/123", nil)

	err := base.Client.DeleteMonitorAlertDefinition(context.Background(), testMonitorAlertDefinitionServiceType, testMonitorAlertDefinitionID)
	assert.NoError(t, err)
}
