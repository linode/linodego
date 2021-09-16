package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/linode/linodego"
)

var testDomainCreateOpts = linodego.DomainCreateOptions{
	Type:     linodego.DomainTypeMaster,
	SOAEmail: "example@example.com",
}

func TestCreateDomain(t *testing.T) {
	_, domain, teardown, err := setupDomain(t, "fixtures/TestCreateDomain")
	defer teardown()

	if err != nil {
		t.Errorf("Error creating domain: %v", err)
	}

	// when comparing fixtures to random value Domain will differ
	if domain.SOAEmail != testDomainCreateOpts.SOAEmail {
		t.Errorf("Domain returned does not match domain create request")
	}
}

func TestUpdateDomain(t *testing.T) {
	client, domain, teardown, err := setupDomain(t, "fixtures/TestUpdateDomain")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	updateOpts := linodego.DomainUpdateOptions{
		Domain: "linodego-renamed-domain.com",
	}
	domain, err = client.UpdateDomain(context.Background(), domain.ID, updateOpts)
	if err != nil {
		t.Errorf("Error renaming domain, %s", err)
	} else if domain.Domain != updateOpts.Domain {
		t.Errorf("Error renaming domain: Domain does not match")
	}
}

func TestListDomains(t *testing.T) {
	client, _, teardown, err := setupDomain(t, "fixtures/TestListDomains")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	domains, err := client.ListDomains(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing domains, expected struct, got error %v", err)
	}
	if len(domains) == 0 {
		t.Errorf("Expected a list of domains, but got %v", domains)
	}
}

func TestGetDomain(t *testing.T) {
	client, domain, teardown, err := setupDomain(t, "fixtures/TestGetDomain")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	_, err = client.GetDomain(context.Background(), domain.ID)
	if err != nil {
		t.Errorf("Error getting domain %d, expected *Domain, got error %v", domain.ID, err)
	}
}

func TestGetDomainZoneFile(t *testing.T) {
	client, domain, teardown, err := setupDomain(t, "fixtures/TestGetDomainZoneFile")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	_, err = client.GetDomainZoneFile(context.Background(), domain.ID)
	if err != nil {
		t.Errorf("failed to get domain zone file %d, expected *DomainZoneFile, got error %v", domain.ID, err)
	}
}

func setupDomain(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.Domain, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	createOpts := testDomainCreateOpts
	createOpts.Domain = fmt.Sprintf("%s-linodego-testing.com", randLabel())

	domain, err := client.CreateDomain(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error creating domain, expected struct, got error %v", err)
	}

	teardown := func() {
		if err := client.DeleteDomain(context.Background(), domain.ID); err != nil {
			t.Errorf("Expected to delete a domain, but got %v", err)
		}
		fixtureTeardown()
	}
	return client, domain, teardown, err
}
