package unit

import (
	"context"
	"fmt"
	"github.com/jarcoal/httpmock"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestDomain_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("domain_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet(formatMockAPIPath("domains"), fixtureData)

	domains, err := base.Client.ListDomains(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.Len(t, domains, 1)

	domain := domains[0]
	assert.Equal(t, 1234, domain.ID)
	assert.Equal(t, "example.org", domain.Domain)
	assert.Equal(t, 300, domain.ExpireSec)
	assert.Equal(t, 300, domain.RefreshSec)
	assert.Equal(t, 300, domain.RetrySec)
	assert.Equal(t, "admin@example.org", domain.SOAEmail)
	assert.Equal(t, linodego.DomainStatus("active"), domain.Status)
	assert.Equal(t, []string{"example tag", "another example"}, domain.Tags)
	assert.Equal(t, 300, domain.TTLSec)
	assert.Equal(t, linodego.DomainType("master"), domain.Type)
	assert.Empty(t, domain.MasterIPs)
	assert.Empty(t, domain.Description)
}

func TestDomain_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("domain_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	domainID := 1234
	base.MockGet(formatMockAPIPath("domains/%d", domainID), fixtureData)

	domain, err := base.Client.GetDomain(context.Background(), domainID)

	assert.NoError(t, err)

	assert.Equal(t, 1234, domain.ID)
	assert.Equal(t, "example.org", domain.Domain)
	assert.Equal(t, 300, domain.ExpireSec)
	assert.Equal(t, 300, domain.RefreshSec)
	assert.Equal(t, 300, domain.RetrySec)
	assert.Equal(t, "admin@example.org", domain.SOAEmail)
	assert.Equal(t, linodego.DomainStatus("active"), domain.Status)
	assert.Equal(t, []string{"example tag", "another example"}, domain.Tags)
	assert.Equal(t, 300, domain.TTLSec)
	assert.Equal(t, linodego.DomainType("master"), domain.Type)
	assert.Empty(t, domain.MasterIPs)
	assert.Empty(t, domain.Description)
}

func TestDomain_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("domain_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	description := ""
	soaEmail := "admin@example.org"
	retrySec := 300
	expireSec := 300
	refreshSec := 300
	ttlSec := 300
	status := linodego.DomainStatus("active")

	requestData := linodego.DomainCreateOptions{
		Domain:      "example.org",
		Type:        linodego.DomainType("master"),
		Description: &description,
		SOAEmail:    &soaEmail,
		RetrySec:    &retrySec,
		MasterIPs:   []string{},
		AXfrIPs:     []string{},
		Tags:        []string{"example tag", "another example"},
		ExpireSec:   &expireSec,
		RefreshSec:  &refreshSec,
		TTLSec:      &ttlSec,
		Status:      &status,
	}

	base.MockPost(formatMockAPIPath("domains"), fixtureData)

	domain, err := base.Client.CreateDomain(context.Background(), requestData)

	assert.NoError(t, err)

	assert.Equal(t, 1234, domain.ID)
	assert.Equal(t, "example.org", domain.Domain)
	assert.Equal(t, 300, domain.ExpireSec)
	assert.Equal(t, 300, domain.RefreshSec)
	assert.Equal(t, 300, domain.RetrySec)
	assert.Equal(t, "admin@example.org", domain.SOAEmail)
	assert.Equal(t, linodego.DomainStatus("active"), domain.Status)
	assert.Equal(t, []string{"example tag", "another example"}, domain.Tags)
	assert.Equal(t, 300, domain.TTLSec)
	assert.Equal(t, linodego.DomainType("master"), domain.Type)
	assert.Empty(t, domain.MasterIPs)
	assert.Empty(t, domain.Description)
}

func TestDomain_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("domain_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	domainID := 1234

	domainName := "example.org"
	domainType := linodego.DomainType("master")
	description := ""
	soaEmail := "admin@example.org"
	retrySec := 300
	expireSec := 300
	refreshSec := 300
	ttlSec := 300
	status := linodego.DomainStatus("active")

	requestData := linodego.DomainUpdateOptions{
		Domain:      &domainName,
		Type:        &domainType,
		Description: &description,
		SOAEmail:    &soaEmail,
		RetrySec:    &retrySec,
		MasterIPs:   []string{},
		AXfrIPs:     []string{},
		Tags:        []string{"example tag", "another example"},
		ExpireSec:   &expireSec,
		RefreshSec:  &refreshSec,
		TTLSec:      &ttlSec,
		Status:      &status,
	}

	base.MockPut(formatMockAPIPath("domains/%d", domainID), fixtureData)

	domain, err := base.Client.UpdateDomain(context.Background(), domainID, requestData)

	assert.NoError(t, err)

	assert.Equal(t, 1234, domain.ID)
	assert.Equal(t, "example.org", domain.Domain)
	assert.Equal(t, 300, domain.ExpireSec)
	assert.Equal(t, 300, domain.RefreshSec)
	assert.Equal(t, 300, domain.RetrySec)
	assert.Equal(t, "admin@example.org", domain.SOAEmail)
	assert.Equal(t, linodego.DomainStatus("active"), domain.Status)
	assert.Equal(t, []string{"example tag", "another example"}, domain.Tags)
	assert.Equal(t, 300, domain.TTLSec)
	assert.Equal(t, linodego.DomainType("master"), domain.Type)
	assert.Empty(t, domain.MasterIPs)
	assert.Empty(t, domain.Description)
}

func TestDomain_Delete(t *testing.T) {
	client := createMockClient(t)

	domainID := 1234

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, fmt.Sprintf("domains/%d", domainID)),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteDomain(context.Background(), domainID); err != nil {
		t.Fatal(err)
	}
}

func TestDomain_GetDomainZoneFile(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("domain_get_domainzonefile")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	domainID := 1234
	base.MockGet(formatMockAPIPath("domains/%d/zone-file", domainID), fixtureData)

	domain, err := base.Client.GetDomainZoneFile(context.Background(), domainID)
	assert.NoError(t, err)

	expectedZoneFile := []string{
		"; example.com [123]",
		"$TTL 864000",
		"@  IN  SOA  ns1.linode.com. user.example.com. 2021000066 14400 14400 1209600 86400",
		"@    NS  ns1.linode.com.",
		"@    NS  ns2.linode.com.",
		"@    NS  ns3.linode.com.",
		"@    NS  ns4.linode.com.",
		"@    NS  ns5.linode.com.",
	}

	assert.Equal(t, expectedZoneFile, domain.ZoneFile)
}

func TestDomain_Clone(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("domain_clone")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.DomainCloneOptions{
		Domain: "linodego-domain-clone.com",
	}

	domainToCloneID := 123
	base.MockPost(formatMockAPIPath("domains/%d/clone", domainToCloneID), fixtureData)

	domain, err := base.Client.CloneDomain(context.Background(), domainToCloneID, requestData)

	assert.NoError(t, err)

	assert.Equal(t, "linodego-domain-clone.com", domain.Domain)
	assert.Equal(t, "admin@example.org", domain.SOAEmail)
}

func TestDomain_Import(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("domain_import")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.DomainImportOptions{
		Domain:           "linodego-domain-import.com",
		RemoteNameserver: "linodego-domain-import-nameserver.com",
	}

	base.MockPost("domains/import", fixtureData)

	domain, err := base.Client.ImportDomain(context.Background(), requestData)

	assert.NoError(t, err)

	assert.Equal(t, "example.org", domain.Domain)
	assert.Equal(t, "admin@example.org", domain.SOAEmail)
}
