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
    "type": "user",
    "service_type": "dbaas",
    "description": "A test alert for dbaas service",
    "scope": "entity",
    "regions": [],
    "entity_ids": [
        "12345"
    ],
    "has_more_resources": false,
    "alert_channels": [
        {
            "id": 10000,
            "label": "Read-Write Channel",
            "type": "email",
            "url": "/monitor/alert-channels/10000"
        }
    ],
    "rule_criteria": {
        "rules": [
            {
                "aggregate_function": "avg",
                "dimension_filters": [
                    {
                        "dimension_label": "node_type",
                        "label": "Node Type",
                        "operator": "eq",
                        "value": "primary"
                    }
                ],
                "label": "High CPU Usage",
                "metric": "cpu_usage",
                "operator": "gt",
                "threshold": 90,
                "unit": "percent"
            }
        ]
    },
    "trigger_conditions": {
        "criteria_condition": "ALL",
        "evaluation_period_seconds": 300,
        "polling_interval_seconds": 60,
        "trigger_occurrences": 3
    },
    "class": "",
    "status": "enabled",
    "entities": {
        "url": "/monitor/services/dbaas/alert-definitions/123/entities",
        "count": 0,
        "has_more_resources": false
    },
    "created": "2024-01-01T00:00:00",
    "updated": "2024-01-01T00:00:00",
    "updated_by": "tester"
	}`

	monitorAlertDefinitionListResponse = `{
    "data": [
        {
            "id": 123,
            "label": "test-alert-definition",
            "severity": 1,
            "type": "user",
            "service_type": "dbaas",
            "description": "A test alert for dbaas service",
            "status": "enabled",
            "scope": "entity",
            "regions": [],
            "entity_ids": [
                "12345"
            ],
            "has_more_resources": true,
            "alert_channels": [
                {
                    "id": 10000,
                    "label": "Read-Write Channel",
                    "type": "email",
                    "url": "/monitor/alert-channels/10000"
                }
            ],
            "rule_criteria": {
                "rules": [
                    {
                        "aggregate_function": "avg",
                        "dimension_filters": [
                            {
                                "dimension_label": "node_type",
                                "label": "Node Type",
                                "operator": "eq",
                                "value": "primary"
                            }
                        ],
                        "label": "High CPU Usage",
                        "metric": "cpu_usage",
                        "operator": "gt",
                        "threshold": 90,
                        "unit": "percent"
                    }
                ]
            },
            "trigger_conditions": {
                "criteria_condition": "ALL",
                "evaluation_period_seconds": 300,
                "polling_interval_seconds": 60,
                "trigger_occurrences": 3
            },
            "class": "",
            "entities": {
                "url": "/monitor/services/dbaas/alert-definitions/123/entities",
                "count": 2,
                "has_more_resources": true
            },
            "created": "2024-01-01T00:00:00",
            "updated": "2024-01-01T00:00:00",
            "updated_by": "tester"
        }
    ],
    "page": 1,
    "pages": 1,
    "results": 1
	}`

	monitorAlertDefinitionUpdateResponse = `{
    "id": 123,
    "label": "test-alert-definition-renamed",
    "severity": 2,
    "type": "user",
    "service_type": "dbaas",
    "status": "disabled",
    "scope": "entity",
    "description": "A test alert for dbaas service",
    "regions": [],
    "entity_ids": [
        "12345",
        "45678"
    ],
    "has_more_resources": true,
    "alert_channels": [
        {
            "id": 10000,
            "label": "Read-Write Channel",
            "type": "email",
            "url": "/monitor/alert-channels/10000"
        }
    ],
    "rule_criteria": {
        "rules": [
            {
                "aggregate_function": "avg",
                "dimension_filters": [
                    {
                        "dimension_label": "node_type",
                        "label": "Node Type",
                        "operator": "eq",
                        "value": "primary"
                    }
                ],
                "label": "High CPU Usage",
                "metric": "cpu_usage",
                "operator": "gt",
                "threshold": 90,
                "unit": "percent"
            }
        ]
    },
    "trigger_conditions": {
        "criteria_condition": "ALL",
        "evaluation_period_seconds": 300,
        "polling_interval_seconds": 60,
        "trigger_occurrences": 3
    },
    "class": "",
    "entities": {
        "url": "/monitor/services/dbaas/alert-definitions/123/entities",
        "count": 2,
        "has_more_resources": true
    },
    "created": "2024-01-01T00:00:00",
    "updated": "2024-01-01T00:00:00",
    "updated_by": "tester"
	}`

	monitorAlertDefinitionUpdateLabelOnlyResponseSingleLine = `{"id": 123, "label": "test-alert-definition-renamed-one-line", "severity": 1, "type": "some_type", "service_type": "dbaas", "status": "enabled", "entity_ids": ["12345"], "channel_ids": [1], "is_enabled": true}`

	monitorAlertDefinitionEntitiesListResponse = `{
		"data": [
			{
				"id": "1",
				"label": "mydatabase-1",
				"url": "/v4/databases/mysql/instances/1",
				"type": "dbaas"
			},
			{
				"id": "2",
				"label": "mydatabase-2",
				"url": "/v4/databases/mysql/instances/2",
				"type": "dbaas"
			},
			{
				"id": "3",
				"label": "mydatabase-3",
				"url": "/v4/databases/mysql/instances/3",
				"type": "dbaas"
			}
		],
		"page": 1,
		"pages": 1,
		"results": 3
	}`
)

func TestCreateMonitorAlertDefinition(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("monitor/services/dbaas/alert-definitions", json.RawMessage(monitorAlertDefinitionGetResponse))

	createOpts := linodego.AlertDefinitionCreateOptions{
		Label:      "test-alert-definition",
		Severity:   int(linodego.SeverityLow),
		Scope:      linodego.AlertDefinitionScopeEntity,
		Regions:    []string{},
		ChannelIDs: []int{1},
		EntityIDs:  []string{"12345"},
	}

	alert, err := base.Client.CreateMonitorAlertDefinition(context.Background(), testMonitorAlertDefinitionServiceType, createOpts)
	assert.NoError(t, err)
	assert.NotNil(t, alert)
	assert.Equal(t, "test-alert-definition", alert.Label)
	assert.Equal(t, testMonitorAlertDefinitionID, alert.ID)
	assert.Equal(t, linodego.AlertDefinitionScopeEntity, alert.Scope)
	assert.Empty(t, alert.Regions)
	assert.Equal(t, "/monitor/services/dbaas/alert-definitions/123/entities", alert.Entities.URL)
	assert.Equal(t, 0, alert.Entities.Count)
	assert.False(t, alert.Entities.HasMoreResources)
	assert.Equal(t, []string{"12345"}, alert.EntityIDs)
	assert.False(t, alert.HasMoreResources)
	assert.NotNil(t, alert.AlertChannels)
	assert.NotNil(t, alert.RuleCriteria)
	assert.NotNil(t, alert.RuleCriteria.Rules)
	assert.NotNil(t, alert.TriggerConditions)
}

func TestCreateMonitorAlertDefinitionWithIdempotency(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("monitor/services/dbaas/alert-definitions", json.RawMessage(monitorAlertDefinitionGetResponse))

	createOpts := linodego.AlertDefinitionCreateOptions{
		Label:      "test-alert-definition",
		Severity:   int(linodego.SeverityLow),
		Scope:      linodego.AlertDefinitionScopeEntity,
		Regions:    []string{},
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
	assert.Equal(t, linodego.AlertDefinitionScopeEntity, alert.Scope)
	assert.Empty(t, alert.Regions)
	assert.Equal(t, "/monitor/services/dbaas/alert-definitions/123/entities", alert.Entities.URL)
	assert.Equal(t, 0, alert.Entities.Count)
	assert.False(t, alert.Entities.HasMoreResources)
}

func TestGetMonitorAlertDefinition(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/services/dbaas/alert-definitions/123", json.RawMessage(monitorAlertDefinitionGetResponse))

	alert, err := base.Client.GetMonitorAlertDefinition(context.Background(), testMonitorAlertDefinitionServiceType, testMonitorAlertDefinitionID)
	assert.NoError(t, err)
	assert.NotNil(t, alert)
	assert.Equal(t, "test-alert-definition", alert.Label)
	assert.Equal(t, testMonitorAlertDefinitionID, alert.ID)
	assert.Equal(t, linodego.AlertDefinitionScopeEntity, alert.Scope)
	assert.Empty(t, alert.Regions)
	assert.Equal(t, "/monitor/services/dbaas/alert-definitions/123/entities", alert.Entities.URL)
	assert.Equal(t, 0, alert.Entities.Count)
	assert.False(t, alert.Entities.HasMoreResources)
	assert.Equal(t, []string{"12345"}, alert.EntityIDs)
	assert.False(t, alert.HasMoreResources)
	assert.NotNil(t, alert.AlertChannels)
	assert.NotNil(t, alert.RuleCriteria)
	assert.NotNil(t, alert.RuleCriteria.Rules)
	assert.NotNil(t, alert.TriggerConditions)
}

func TestListMonitorAlertDefinitions(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("monitor/services/dbaas/alert-definitions", json.RawMessage(monitorAlertDefinitionListResponse))

	alerts, err := base.Client.ListMonitorAlertDefinitions(context.Background(), testMonitorAlertDefinitionServiceType, nil)
	assert.NoError(t, err)
	assert.Len(t, alerts, 1)
	assert.Equal(t, "test-alert-definition", alerts[0].Label)
	assert.Equal(t, testMonitorAlertDefinitionID, alerts[0].ID)
	assert.Equal(t, linodego.AlertDefinitionScopeEntity, alerts[0].Scope)
	assert.Empty(t, alerts[0].Regions)
	assert.Equal(t, "/monitor/services/dbaas/alert-definitions/123/entities", alerts[0].Entities.URL)
	assert.Equal(t, 2, alerts[0].Entities.Count)
	assert.True(t, alerts[0].Entities.HasMoreResources)
	assert.Equal(t, []string{"12345"}, alerts[0].EntityIDs)
	assert.True(t, alerts[0].HasMoreResources)
	assert.NotNil(t, alerts[0].AlertChannels)
	assert.NotNil(t, alerts[0].RuleCriteria)
	assert.NotNil(t, alerts[0].RuleCriteria.Rules)
	assert.NotNil(t, alerts[0].TriggerConditions)
}

func TestUpdateMonitorAlertDefinition(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
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
	assert.Equal(t, linodego.AlertDefinitionScopeEntity, alert.Scope)
	assert.Empty(t, alert.Regions)
	assert.Equal(t, "/monitor/services/dbaas/alert-definitions/123/entities", alert.Entities.URL)
	assert.Equal(t, 2, alert.Entities.Count)
	assert.True(t, alert.Entities.HasMoreResources)
	assert.True(t, alert.HasMoreResources)
}

func TestUpdateMonitorAlertDefinition_LabelOnly(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
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
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("monitor/services/dbaas/alert-definitions/123", nil)

	err := base.Client.DeleteMonitorAlertDefinition(context.Background(), testMonitorAlertDefinitionServiceType, testMonitorAlertDefinitionID)
	assert.NoError(t, err)
}

func TestListMonitorAlertDefinitionEntities(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet(
		"monitor/services/dbaas/alert-definitions/123/entities",
		json.RawMessage(monitorAlertDefinitionEntitiesListResponse),
	)

	entities, err := base.Client.ListMonitorAlertDefinitionEntities(
		context.Background(),
		testMonitorAlertDefinitionServiceType,
		testMonitorAlertDefinitionID, nil,
	)
	assert.NoError(t, err)
	assert.Len(t, entities, 3)

	assert.Equal(t, "1", entities[0].ID)
	assert.Equal(t, "mydatabase-1", entities[0].Label)
	assert.Equal(t, "/v4/databases/mysql/instances/1", entities[0].URL)
	assert.Equal(t, "dbaas", entities[0].Type)

	assert.Equal(t, "2", entities[1].ID)
	assert.Equal(t, "mydatabase-2", entities[1].Label)
	assert.Equal(t, "/v4/databases/mysql/instances/2", entities[1].URL)
	assert.Equal(t, "dbaas", entities[1].Type)

	assert.Equal(t, "3", entities[2].ID)
	assert.Equal(t, "mydatabase-3", entities[2].Label)
	assert.Equal(t, "/v4/databases/mysql/instances/3", entities[2].URL)
	assert.Equal(t, "dbaas", entities[2].Type)
}
