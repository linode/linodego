package unit

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestFirewallRuleSets_List(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	response := map[string]any{
		"data": []map[string]any{
			{
				"id":          101,
				"label":       "allow-ssh",
				"description": "Allow SSH inbound",
				"type":        "inbound",
				"rules": []map[string]any{
					{
						"action": "ACCEPT",
						"addresses": map[string]any{
							"ipv4": []string{"pl::vpcs:primary"},
							"ipv6": []string{"pl::vpcs:primary"},
						},
						"label":       "ssh",
						"ports":       "22",
						"protocol":    "TCP",
						"description": "Allow inbound SSH",
					},
				},
				"version":            3,
				"is_service_defined": false,
				"created":            "2024-01-01T00:00:01",
				"updated":            "2024-01-02T00:00:02",
				"deleted":            nil,
			},
		},
		"page":    1,
		"pages":   1,
		"results": 1,
	}

	base.MockGet(formatMockAPIPath("networking/firewalls/rulesets"), response)

	ruleSets, err := base.Client.ListFirewallRuleSets(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, ruleSets, 1)

	rs := ruleSets[0]
	assert.Equal(t, 101, rs.ID)
	assert.Equal(t, "allow-ssh", rs.Label)
	assert.Equal(t, linodego.FirewallRuleSetTypeInbound, rs.Type)
	assert.False(t, rs.IsServiceDefined)
	assert.Len(t, rs.Rules, 1)

	rule := rs.Rules[0]
	assert.Equal(t, "ACCEPT", rule.Action)
	if assert.NotNil(t, rule.Addresses.IPv4) {
		assert.Equal(t, []string{"pl::vpcs:primary"}, *rule.Addresses.IPv4)
	}
	if assert.NotNil(t, rule.Addresses.IPv6) {
		assert.Equal(t, []string{"pl::vpcs:primary"}, *rule.Addresses.IPv6)
	}

	if assert.NotNil(t, rs.Created) {
		assert.Equal(t, time.Date(2024, time.January, 1, 0, 0, 1, 0, time.UTC), rs.Created.UTC())
	}
	if assert.NotNil(t, rs.Updated) {
		assert.Equal(t, time.Date(2024, time.January, 2, 0, 0, 2, 0, time.UTC), rs.Updated.UTC())
	}
	assert.Nil(t, rs.Deleted)
}

func TestFirewallRuleSets_Get(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	ruleSetID := 202
	response := map[string]any{
		"id":          ruleSetID,
		"label":       "allow-egress",
		"description": "Allow outbound traffic",
		"type":        "outbound",
		"rules": []map[string]any{
			{
				"ruleset": 77,
			},
		},
		"version":            4,
		"is_service_defined": true,
		"created":            "2024-02-01T01:01:01",
		"updated":            "2024-02-02T02:02:02",
		"deleted":            nil,
	}

	base.MockGet(formatMockAPIPath("networking/firewalls/rulesets/%d", ruleSetID), response)

	rs, err := base.Client.GetFirewallRuleSet(context.Background(), ruleSetID)
	assert.NoError(t, err)
	assert.Equal(t, ruleSetID, rs.ID)
	assert.Equal(t, linodego.FirewallRuleSetTypeOutbound, rs.Type)
	assert.True(t, rs.IsServiceDefined)
	if assert.Len(t, rs.Rules, 1) {
		assert.Equal(t, 77, rs.Rules[0].RuleSet)
	}
	if assert.NotNil(t, rs.Created) {
		assert.Equal(t, time.Date(2024, time.February, 1, 1, 1, 1, 0, time.UTC), rs.Created.UTC())
	}
	if assert.NotNil(t, rs.Updated) {
		assert.Equal(t, time.Date(2024, time.February, 2, 2, 2, 2, 0, time.UTC), rs.Updated.UTC())
	}
}

func TestFirewallRuleSets_Create(t *testing.T) {
	client := createMockClient(t)

	req := linodego.RuleSetCreateOptions{
		Label:       "allow-vpc",
		Description: "Allow VPC ingress",
		Type:        linodego.FirewallRuleSetTypeInbound,
		Rules: []linodego.FirewallRule{
			{
				Action: "ACCEPT",
				Addresses: linodego.NetworkAddresses{
					IPv4: &[]string{"pl::vpcs:primary"},
				},
				Label:    "ingress",
				Protocol: linodego.NetworkProtocol("TCP"),
				Ports:    "443",
			},
		},
	}

	response := map[string]any{
		"id":          303,
		"label":       req.Label,
		"description": req.Description,
		"type":        string(req.Type),
		"rules": []map[string]any{
			{
				"action": "ACCEPT",
				"addresses": map[string]any{
					"ipv4": []string{"pl::vpcs:primary"},
				},
				"label":    "ingress",
				"ports":    "443",
				"protocol": "TCP",
			},
		},
		"version":            1,
		"is_service_defined": false,
	}

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "networking/firewalls/rulesets"),
		mockRequestBodyValidate(t, req, response))

	rs, err := client.CreateFirewallRuleSet(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, 303, rs.ID)
	assert.Equal(t, req.Label, rs.Label)
	assert.Len(t, rs.Rules, 1)
}

func TestFirewallRuleSets_Update(t *testing.T) {
	client := createMockClient(t)

	ruleSetID := 404
	label := "updated-egress"
	description := "Updated description"
	rules := []linodego.FirewallRule{
		{
			Action: "ACCEPT",
			Addresses: linodego.NetworkAddresses{
				IPv6: &[]string{"pl::system:obj-storage"},
			},
			Label:    "egress",
			Protocol: linodego.NetworkProtocol("UDP"),
			Ports:    "53",
		},
	}

	req := linodego.RuleSetUpdateOptions{
		Label:       &label,
		Description: &description,
		Rules:       &rules,
	}

	response := map[string]any{
		"id":          ruleSetID,
		"label":       label,
		"description": description,
		"type":        "outbound",
		"rules": []map[string]any{
			{
				"action": "ACCEPT",
				"addresses": map[string]any{
					"ipv6": []string{"pl::system:obj-storage"},
				},
				"label":    "egress",
				"ports":    "53",
				"protocol": "UDP",
			},
		},
		"version":            6,
		"is_service_defined": false,
	}

	httpmock.RegisterRegexpResponder("PUT", mockRequestURL(t, fmt.Sprintf("networking/firewalls/rulesets/%d", ruleSetID)),
		mockRequestBodyValidate(t, req, response))

	rs, err := client.UpdateFirewallRuleSet(context.Background(), ruleSetID, req)
	assert.NoError(t, err)
	assert.Equal(t, label, rs.Label)
	assert.Len(t, rs.Rules, 1)
}

func TestFirewallRuleSets_Delete(t *testing.T) {
	client := createMockClient(t)

	ruleSetID := 505
	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, fmt.Sprintf("networking/firewalls/rulesets/%d", ruleSetID)),
		httpmock.NewStringResponder(200, "{}"))

	err := client.DeleteFirewallRuleSet(context.Background(), ruleSetID)
	assert.NoError(t, err)
}

func TestRuleSet_UnmarshalJSON(t *testing.T) {
	var rs linodego.RuleSet

	raw := []byte(`{
		"id": 42,
		"label": "combined",
		"type": "inbound",
		"description": "Combined rule set",
		"rules": [
			{"ruleset": 12}
		],
		"version": 2,
		"is_service_defined": false,
		"created": "2023-03-03T03:03:03",
		"updated": "2023-04-04T04:04:04",
		"deleted": "2023-05-05T05:05:05"
	}`)

	assert.NoError(t, json.Unmarshal(raw, &rs))
	assert.Equal(t, 42, rs.ID)
	if assert.NotNil(t, rs.Created) {
		assert.Equal(t, time.Date(2023, time.March, 3, 3, 3, 3, 0, time.UTC), rs.Created.UTC())
	}
	if assert.NotNil(t, rs.Updated) {
		assert.Equal(t, time.Date(2023, time.April, 4, 4, 4, 4, 0, time.UTC), rs.Updated.UTC())
	}
	if assert.NotNil(t, rs.Deleted) {
		assert.Equal(t, time.Date(2023, time.May, 5, 5, 5, 5, 0, time.UTC), rs.Deleted.UTC())
	}
	if assert.Len(t, rs.Rules, 1) {
		assert.Equal(t, 12, rs.Rules[0].RuleSet)
	}
}

func TestRuleSet_UnmarshalJSONNumericServiceDefined(t *testing.T) {
	var rs linodego.RuleSet

	raw := []byte(`{
		"id": 99,
		"label": "numeric",
		"type": "inbound",
		"rules": [],
		"version": 1,
		"is_service_defined": 1
	}`)

	assert.NoError(t, json.Unmarshal(raw, &rs))
	assert.True(t, rs.IsServiceDefined)
}
