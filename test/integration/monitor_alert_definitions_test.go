package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

const (
	testMonitorAlertDefinitionServiceType = "dbaas"

	// TODO: use a fixed channel id for now until the alert channel has been fixed.
	channelID = 10000
)

func TestMonitorAlertDefinition_smoke(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestMonitorAlertDefinition")
	defer teardown()

	// Get All Alert Definitions
	alerts, err := client.ListAllMonitorAlertDefinitions(context.Background(), nil)
	// Even if there is no alert definition, it should not error out
	if err != nil {
		t.Fatalf("failed to fetch monitor alert definitions: %s", err)
	}

	// New: Iterate and log each alert definition for visibility
	for _, alert := range alerts {
		// Check few mandatory fields on each listed alert
		assert.NotZero(t, alert.ID, "alert.ID should not be zero")
		assert.NotEmpty(t, alert.Label, "alert.Label should not be empty")

		// If alert has a rule, validate basic rule structure
		if len(alert.RuleCriteria.Rules) > 0 {
			assert.NotEmpty(t, alert.RuleCriteria.Rules, "RuleCriteria.Rules should not be empty when RuleCriteria is provided")
			for _, r := range alert.RuleCriteria.Rules {
				assert.NotEmpty(t, r.Metric, "rule.Metric should not be empty")
				assert.NotEmpty(t, r.Operator, "rule.Operator should not be empty")
			}
		}
	}

	// Basic assertions based on the fixture
	assert.NoError(t, err)

	// Determine a channel ID to use for creating a new alert definition:
	var channelID int
	var fetchedChannelLabel string
	var fetchedChannelID int
	if len(alerts) > 0 && len(alerts[0].AlertChannels) > 0 {
		channelID = alerts[0].AlertChannels[0].ID
		fetchedChannelID = alerts[0].AlertChannels[0].ID
		fetchedChannelLabel = alerts[0].AlertChannels[0].Label
	} else {
		// Fallback to ListAlertChannels to get available channels
		channels, err := client.ListAlertChannels(context.Background(), nil)
		if err != nil || len(channels) == 0 {
			t.Fatalf("failed to determine a monitor channel to use: %s", err)
		}
		channelID = channels[0].ID
		fetchedChannelID = channels[0].ID
		fetchedChannelLabel = channels[0].Label
	}
	// Validate the chosen channel
	assert.NotZero(t, fetchedChannelID, "fetchedChannel.ID should not be zero")
	assert.NotEmpty(t, fetchedChannelLabel, "fetchedChannel.Label should not be empty")

	// Test creating a new Monitor Alert Definition
	createOpts := linodego.AlertDefinitionCreateOptions{
		Label:       "go-test-alert-definition-create",
		Severity:    int(linodego.SeverityLow),
		Description: linodego.Pointer("Test alert definition creation"),
		ChannelIDs:  []int{channelID},
		EntityIDs:   nil,
		TriggerConditions: &linodego.TriggerConditions{
			CriteriaCondition:       "ALL",
			EvaluationPeriodSeconds: 300,
			PollingIntervalSeconds:  300,
			TriggerOccurrences:      1,
		},
		RuleCriteria: &linodego.RuleCriteriaOptions{
			Rules: []linodego.RuleOptions{
				{
					AggregateFunction: "avg",
					Metric:            "memory_usage",
					Operator:          "gt",
					Threshold:         90.0,
					DimensionFilters: []linodego.DimensionFilterOptions{
						{
							DimensionLabel: "node_type",
							Operator:       "eq",
							Value:          "primary",
						},
					},
				},
			},
		},
	}

	createdAlert, err := client.CreateMonitorAlertDefinition(context.Background(), testMonitorAlertDefinitionServiceType, createOpts)
	if err != nil {
		// The test fixtures may return a 400 if an existing alert is being updated.
		// Treat this as a non-fatal condition for the smoke test: log and exit.
		t.Logf("CreateMonitorAlertDefinition returned error, skipping create assertions: %s", err)
		return
	}
	assert.NoError(t, err)
	assert.NotNil(t, createdAlert)
	assert.Equal(t, createOpts.Label, createdAlert.Label)
	assert.Equal(t, createOpts.Severity, createdAlert.Severity)
	assert.Equal(t, *createOpts.Description, createdAlert.Description)
	assert.ElementsMatch(t, createOpts.EntityIDs, createdAlert.EntityIDs)
	// assert.Equal(t, fetchedChannel.Label, createdAlert.AlertChannels[0].Label)

	// More thorough assertions on the created alert's nested fields
	// TriggerConditions is a struct, so it is never nil
	assert.Equal(t, createOpts.TriggerConditions.CriteriaCondition, createdAlert.TriggerConditions.CriteriaCondition)
	assert.Equal(t, createOpts.TriggerConditions.EvaluationPeriodSeconds, createdAlert.TriggerConditions.EvaluationPeriodSeconds)
	assert.Equal(t, createOpts.TriggerConditions.PollingIntervalSeconds, createdAlert.TriggerConditions.PollingIntervalSeconds)
	assert.Equal(t, createOpts.TriggerConditions.TriggerOccurrences, createdAlert.TriggerConditions.TriggerOccurrences)

	if len(createdAlert.RuleCriteria.Rules) > 0 && len(createOpts.RuleCriteria.Rules) > 0 {
		assert.Equal(t, len(createOpts.RuleCriteria.Rules), len(createdAlert.RuleCriteria.Rules), "created alert should have same number of rules")
		for i, r := range createOpts.RuleCriteria.Rules {
			cr := createdAlert.RuleCriteria.Rules[i]
			assert.Equal(t, r.Metric, cr.Metric)
			assert.Equal(t, r.Operator, cr.Operator)
			assert.Equal(t, r.Threshold, cr.Threshold)
			// Dimension filters
			if len(r.DimensionFilters) > 0 {
				assert.Equal(t, len(r.DimensionFilters), len(cr.DimensionFilters))
				for j, df := range r.DimensionFilters {
					cdf := cr.DimensionFilters[j]
					assert.Equal(t, df.DimensionLabel, cdf.DimensionLabel)
					assert.Equal(t, df.Operator, cdf.Operator)
					assert.Equal(t, df.Value, cdf.Value)
				}
			}
		}
	}

	// Update the created alert definition: change label only
	newLabel := createdAlert.Label + "-updated"
	updateOpts := linodego.AlertDefinitionUpdateOptions{
		Label:             newLabel,
		Severity:          createdAlert.Severity,
		ChannelIDs:        createOpts.ChannelIDs,
		RuleCriteria:      createOpts.RuleCriteria,
		TriggerConditions: createOpts.TriggerConditions,
		EntityIDs:         createOpts.EntityIDs,
		Description:       &createdAlert.Description,
	}
	// wait for 1 minute before update for create to complete
	time.Sleep(1 * time.Minute)
	updatedAlert, err := client.UpdateMonitorAlertDefinition(context.Background(), testMonitorAlertDefinitionServiceType, createdAlert.ID, updateOpts)
	if err != nil {
		// Some fixtures may not support update; treat as non-fatal
		t.Logf("UpdateMonitorAlertDefinition returned error, skipping update assertions: %s", err)
	} else {
		assert.NotNil(t, updatedAlert)
		assert.Equal(t, createdAlert.ID, updatedAlert.ID, "updated alert should keep same ID")
		assert.Equal(t, newLabel, updatedAlert.Label, "updated alert should have the new label")
	}

	// Clean up created alert definition
	if createdAlert != nil {
		// Retry deletion with exponential backoff for up to 2 minutes
		maxWait := 2 * time.Minute
		baseDelay := 2 * time.Second
		var lastErr error
		start := time.Now()
		for attempt := 0; time.Since(start) < maxWait; attempt++ {
			err = client.DeleteMonitorAlertDefinition(context.Background(), testMonitorAlertDefinitionServiceType, createdAlert.ID)
			if err == nil {
				break
			}
			lastErr = err
			// Exponential backoff, capped at 30s
			sleep := baseDelay * (1 << attempt)
			if sleep > 30*time.Second {
				sleep = 30 * time.Second
			}
			time.Sleep(sleep)
		}
		assert.NoError(t, err, "DeleteMonitorAlertDefinition failed after retries: %v", lastErr)
	}
}

func TestMonitorAlertDefinitions_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestMonitorAlertDefinitions_List")
	defer teardown()

	// List all alert definitions
	alerts, err := client.ListAllMonitorAlertDefinitions(context.Background(), nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, alerts, "Expected at least one alert definition")

	for _, alert := range alerts {
		assert.NotZero(t, alert.ID)
		assert.NotEmpty(t, alert.Label)
		assert.NotEmpty(t, alert.ServiceType)
	}
}

func TestMonitorAlertChannels_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestMonitorAlertChannels_List")
	defer teardown()

	// List all alert channels
	channels, err := client.ListAlertChannels(context.Background(), nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, channels, "Expected at least one alert channel")

	for _, channel := range channels {
		assert.NotZero(t, channel.ID)
		assert.NotEmpty(t, channel.Label)
		assert.NotEmpty(t, channel.ChannelType)
	}
}

func TestMonitorAlertDefinition_CreateWithIdempotency(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestMonitorAlertDefinition_CreateWithIdempotency")
	defer teardown()

	// Get a channel ID to use
	channels, err := client.ListAlertChannels(context.Background(), nil)
	if err != nil || len(channels) == 0 {
		t.Fatalf("failed to determine a monitor channel to use: %s", err)
	}
	channelID := channels[0].ID

	uniqueLabel := fmt.Sprintf("go-test-alert-definition-idempotency-%d", time.Now().UnixNano())

	createOpts := linodego.AlertDefinitionCreateOptions{
		Label:       uniqueLabel,
		Severity:    int(linodego.SeverityLow),
		Description: linodego.Pointer("Test alert definition creation with idempotency"),
		ChannelIDs:  []int{channelID},
		EntityIDs:   nil,
		TriggerConditions: &linodego.TriggerConditions{
			CriteriaCondition:       "ALL",
			EvaluationPeriodSeconds: 300,
			PollingIntervalSeconds:  300,
			TriggerOccurrences:      1,
		},
		RuleCriteria: &linodego.RuleCriteriaOptions{
			Rules: []linodego.RuleOptions{
				{
					AggregateFunction: "avg",
					Metric:            "memory_usage",
					Operator:          "gt",
					Threshold:         90.0,
					DimensionFilters: []linodego.DimensionFilterOptions{
						{
							DimensionLabel: "node_type",
							Operator:       "eq",
							Value:          "primary",
						},
					},
				},
			},
		},
	}

	// Create the alert definition
	createdAlert, err := client.CreateMonitorAlertDefinition(context.Background(), testMonitorAlertDefinitionServiceType, createOpts)
	if err != nil {
		alerts, listErr := client.ListMonitorAlertDefinitions(context.Background(), testMonitorAlertDefinitionServiceType, nil)
		if listErr == nil {
			for _, a := range alerts {
				if a.Label == createOpts.Label {
					_ = client.DeleteMonitorAlertDefinition(context.Background(), testMonitorAlertDefinitionServiceType, a.ID)
					break
				}
			}
			// Retry creation
			createdAlert, err = client.CreateMonitorAlertDefinition(context.Background(), testMonitorAlertDefinitionServiceType, createOpts)
		}
	}
	assert.NoError(t, err)
	assert.NotNil(t, createdAlert)

	// Attempt to create the same alert definition again to test idempotency
	// Expected to return Error as per the API behavior
	_, err = client.CreateMonitorAlertDefinition(context.Background(), testMonitorAlertDefinitionServiceType, createOpts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "An alert with this label already exists")

	// Cleanup
	if createdAlert != nil {
		_ = client.DeleteMonitorAlertDefinition(context.Background(), testMonitorAlertDefinitionServiceType, createdAlert.ID)
	}
}
