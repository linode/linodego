package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// AlertDefinition represents an ACLP Alert Definition object
type AlertDefinition struct {
	ID                int                    `json:"id"`
	Label             string                 `json:"label"`
	Severity          int                    `json:"severity"`
	Type              string                 `json:"type"`
	ServiceType       string                 `json:"service_type"`
	Status            string                 `json:"status"`
	HasMoreResources  bool                   `json:"has_more_resources"`
	Rule              *Rule                  `json:"rule"`
	RuleCriteria      *RuleCriteria          `json:"rule_criteria"`
	TriggerConditions *TriggerConditions     `json:"trigger_conditions"`
	AlertChannels     []AlertChannelEnvelope `json:"alert_channels"`
	Created           *time.Time             `json:"-"`
	Updated           *time.Time             `json:"-"`
	UpdatedBy         string                 `json:"updated_by"`
	CreatedBy         string                 `json:"created_by"`
	EntityIDs         []string               `json:"entity_ids"`
	Description       string                 `json:"description"`
	Class             string                 `json:"class"`
}

// Backwards-compatible alias
type MonitorAlertDefinition = AlertDefinition

// TriggerConditions represents the trigger conditions for an alert.
type TriggerConditions struct {
	CriteriaCondition       string `json:"criteria_condition,omitempty"`
	EvaluationPeriodSeconds int    `json:"evaluation_period_seconds,omitempty"`
	PollingIntervalSeconds  int    `json:"polling_interval_seconds,omitempty"`
	TriggerOccurrences      int    `json:"trigger_occurrences,omitempty"`
}

// RuleCriteria represents the rule criteria for an alert.
type RuleCriteria struct {
	Rules []Rule `json:"rules,omitempty"`
}

// Rule represents a single rule for an alert.
type Rule struct {
	AggregateFunction string            `json:"aggregate_function,omitempty"`
	DimensionFilters  []DimensionFilter `json:"dimension_filters,omitempty"`
	Label             string            `json:"label,omitempty"`
	Metric            string            `json:"metric,omitempty"`
	Operator          string            `json:"operator,omitempty"`
	Threshold         *float64          `json:"threshold,omitempty"`
	Unit              *string           `json:"unit,omitempty"`
}

// DimensionFilter represents a single dimension filter used inside a Rule.
type DimensionFilter struct {
	DimensionLabel string      `json:"dimension_label"`
	Label          string      `json:"label"`
	Operator       string      `json:"operator"`
	Value          interface{} `json:"value"`
}

// AlertChannelEnvelope represents a single alert channel entry returned inside alert definition
type AlertChannelEnvelope struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

// AlertType represents the type of alert: "user" or "system"
type AlertType string

const (
	AlertTypeUser   AlertType = "user"
	AlertTypeSystem AlertType = "system"
)

// Severity represents the severity level of an alert.
// 0 = Severe, 1 = Medium, 2 = Low, 3 = Info
type Severity int

const (
	SeveritySevere Severity = 0
	SeverityMedium Severity = 1
	SeverityLow    Severity = 2
	SeverityInfo   Severity = 3
)

// AlertDefinitionCreateOptions are the options used to create a new alert definition.
type AlertDefinitionCreateOptions struct {
	ServiceType       string             `json:"service_type"`                 // mandatory
	Label             string             `json:"label"`                        // mandatory
	Severity          int                `json:"severity"`                     // mandatory
	ChannelIDs        []int              `json:"channel_ids"`                  // mandatory
	RuleCriteria      *RuleCriteria      `json:"rule_criteria,omitempty"`      // optional
	TriggerConditions *TriggerConditions `json:"trigger_conditions,omitempty"` // optional
	EntityIDs         []string           `json:"entity_ids,omitempty"`         // optional
	Description       string             `json:"description,omitempty"`        // optional
}

// AlertDefinitionUpdateOptions are the options used to update an alert definition.
type AlertDefinitionUpdateOptions struct {
	ServiceType       string             `json:"service_type"`                 // mandatory, must not be empty
	AlertID           int                `json:"alert_id"`                     // mandatory, must not be zero
	Label             *string            `json:"label,omitempty"`              // optional
	Severity          *int               `json:"severity,omitempty"`           // optional, should be int to match AlertDefinition
	Description       *string            `json:"description,omitempty"`        // optional
	RuleCriteria      *RuleCriteria      `json:"rule_criteria,omitempty"`      // optional
	TriggerConditions *TriggerConditions `json:"trigger_conditions,omitempty"` // optional
	EntityIDs         []string           `json:"entity_ids,omitempty"`         // optional
	ChannelIDs        []int              `json:"channel_ids,omitempty"`        // optional
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *AlertDefinition) UnmarshalJSON(b []byte) error {
	type Mask AlertDefinition

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Created = (*time.Time)(p.Created)
	i.Updated = (*time.Time)(p.Updated)

	return nil
}

// ListMonitorAlertDefinitions gets a paginated list of ACLP Monitor Alert Definitions.
func (c *Client) ListMonitorAlertDefinitions(ctx context.Context, serviceType string, opts *ListOptions) ([]MonitorAlertDefinition, error) {
	var endpoint string
	if serviceType != "" {
		endpoint = formatAPIV4BetaPath("monitor/services/%s/alert-definitions", serviceType)
	} else {
		endpoint = formatAPIV4BetaPath("monitor/alert-definitions")
	}
	return getPaginatedResults[AlertDefinition](ctx, c, endpoint, opts)
}

// GetMonitorAlertDefinition gets an ACLP Monitor Alert Definition.
func (c *Client) GetMonitorAlertDefinition(ctx context.Context, serviceType string, alertID int) (*MonitorAlertDefinition, error) {
	e := formatAPIV4BetaPath("monitor/services/%s/alert-definitions/%d", serviceType, alertID)
	return doGETRequest[AlertDefinition](ctx, c, e)
}

// CreateMonitorAlertDefinition creates an ACLP Monitor Alert Definition.
func (c *Client) CreateMonitorAlertDefinition(ctx context.Context, serviceType string, opts AlertDefinitionCreateOptions) (*MonitorAlertDefinition, error) {
	e := formatAPIV4BetaPath("monitor/services/%s/alert-definitions", serviceType)
	return doPOSTRequest[AlertDefinition](ctx, c, e, opts)
}

// CreateMonitorAlertDefinitionWithIdempotency creates an ACLP Monitor Alert Definition
// and optionally sends an Idempotency-Key header to make the request idempotent.
func (c *Client) CreateMonitorAlertDefinitionWithIdempotency(ctx context.Context, serviceType string, opts AlertDefinitionCreateOptions, idempotencyKey string) (*MonitorAlertDefinition, error) {
	e := formatAPIV4BetaPath("monitor/services/%s/alert-definitions", serviceType)

	var result AlertDefinition
	req := c.R(ctx).SetResult(&result)

	if idempotencyKey != "" {
		req.SetHeader("Idempotency-Key", idempotencyKey)
	}

	body, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	req.SetBody(string(body))

	r, err := coupleAPIErrors(req.Post(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*AlertDefinition), nil
}

// UpdateMonitorAlertDefinition updates an ACLP Monitor Alert Definition.
func (c *Client) UpdateMonitorAlertDefinition(ctx context.Context, serviceType string, alertID int, opts AlertDefinitionUpdateOptions) (*AlertDefinition, error) {
	e := formatAPIV4BetaPath("monitor/services/%s/alert-definitions/%d", serviceType, alertID)
	return doPUTRequest[AlertDefinition](ctx, c, e, opts)
}

// DeleteMonitorAlertDefinition deletes an ACLP Monitor Alert Definition.
func (c *Client) DeleteMonitorAlertDefinition(ctx context.Context, serviceType string, alertID int) error {
	e := formatAPIV4BetaPath("monitor/services/%s/alert-definitions/%d", serviceType, alertID)
	return doDELETERequest(ctx, c, e)
}
