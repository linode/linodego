package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// MonitorAlertDefinition represents an ACLP Alert Definition object
type MonitorAlertDefinition struct {
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

// AlertTriggerConditions represents the trigger conditions for an alert.
type TriggerConditions struct {
	CriteriaCondition       string `json:"criteria_condition,omitempty"`
	EvaluationPeriodSeconds int    `json:"evaluation_period_seconds,omitempty"`
	PollingIntervalSeconds  int    `json:"polling_interval_seconds,omitempty"`
	TriggerOccurrences      int    `json:"trigger_occurrences,omitempty"`
}

// AlertRuleCriteria represents the rule criteria for an alert.
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

const (
	MonitorAlertDefinitionSeverityCritical = 1
	MonitorAlertDefinitionSeverityMajor    = 2
	MonitorAlertDefinitionSeverityMinor    = 3
)

const (
	MonitorAlertDefinitionStatusEnabled  = "enabled"
	MonitorAlertDefinitionStatusDisabled = "disabled"
)
// AlertChannelEnvelope represents a single alert channel entry returned inside alert definition
type AlertChannelEnvelope struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

// MonitorAlertDefinitionCreateOptions are the options used to create a new alert definition.
type MonitorAlertDefinitionCreateOptions struct {
	Label             string             `json:"label"`
	Severity          int                `json:"severity"`
	Class             string             `json:"class"`
	Type              string             `json:"type"`
	Description       string             `json:"description,omitempty"`
	ChannelIDs        []int              `json:"channel_ids"`
	EntityIDs         []string           `json:"entity_ids"`
	IsEnabled         bool               `json:"is_enabled"`
	TriggerConditions *TriggerConditions `json:"trigger_conditions,omitempty"`
	Rule              *Rule              `json:"rule,omitempty"`
	RuleCriteria      *RuleCriteria      `json:"rule_criteria,omitempty"`
}

// MonitorAlertDefinitionUpdateOptions are the options used to update an alert definition.
type MonitorAlertDefinitionUpdateOptions struct {
	Label             string             `json:"label,omitempty"`
	Severity          int                `json:"severity,omitempty"`
	Class             string             `json:"class,omitempty"`
	Description       string             `json:"description,omitempty"`
	ChannelIDs        []int              `json:"channel_ids,omitempty"`
	EntityIDs         []string           `json:"entity_ids,omitempty"`
	IsEnabled         *bool              `json:"is_enabled,omitempty"`
	TriggerConditions *TriggerConditions `json:"trigger_conditions,omitempty"`
	Rule              *Rule              `json:"rule,omitempty"`
	RuleCriteria      *RuleCriteria      `json:"rule_criteria,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *MonitorAlertDefinition) UnmarshalJSON(b []byte) error {
	type Mask MonitorAlertDefinition

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
	return getPaginatedResults[MonitorAlertDefinition](ctx, c, endpoint, opts)
}

// GetMonitorAlertDefinition gets an ACLP Monitor Alert Definition.
func (c *Client) GetMonitorAlertDefinition(ctx context.Context, serviceType string, alertID int) (*MonitorAlertDefinition, error) {
	e := formatAPIV4BetaPath("monitor/services/%s/alert-definitions/%d", serviceType, alertID)
	return doGETRequest[MonitorAlertDefinition](ctx, c, e)
}

// CreateMonitorAlertDefinition creates an ACLP Monitor Alert Definition.
func (c *Client) CreateMonitorAlertDefinition(ctx context.Context, serviceType string, opts MonitorAlertDefinitionCreateOptions) (*MonitorAlertDefinition, error) {
	e := formatAPIV4BetaPath("monitor/services/%s/alert-definitions", serviceType)
	return doPOSTRequest[MonitorAlertDefinition](ctx, c, e, opts)
}

// CreateMonitorAlertDefinitionWithIdempotency creates an ACLP Monitor Alert Definition
// and optionally sends an Idempotency-Key header to make the request idempotent.
func (c *Client) CreateMonitorAlertDefinitionWithIdempotency(ctx context.Context, serviceType string, opts MonitorAlertDefinitionCreateOptions, idempotencyKey string) (*MonitorAlertDefinition, error) {
	e := formatAPIV4BetaPath("monitor/services/%s/alert-definitions", serviceType)

	var result MonitorAlertDefinition
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

	return r.Result().(*MonitorAlertDefinition), nil
}

// UpdateMonitorAlertDefinition updates an ACLP Monitor Alert Definition.
func (c *Client) UpdateMonitorAlertDefinition(ctx context.Context, serviceType string, alertID int, opts MonitorAlertDefinitionUpdateOptions) (*MonitorAlertDefinition, error) {
	e := formatAPIV4BetaPath("monitor/services/%s/alert-definitions/%d", serviceType, alertID)
	return doPUTRequest[MonitorAlertDefinition, MonitorAlertDefinitionUpdateOptions](ctx, c, e, opts)
}

// DeleteMonitorAlertDefinition deletes an ACLP Monitor Alert Definition.
func (c *Client) DeleteMonitorAlertDefinition(ctx context.Context, serviceType string, alertID int) error {
	e := formatAPIV4BetaPath("monitor/services/%s/alert-definitions/%d", serviceType, alertID)
	return doDELETERequest(ctx, c, e)
}
