package unit

import (
	"context"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
