package integration

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/linode/linodego"
)

var testDomainRecordCreateOpts = linodego.DomainRecordCreateOptions{
	Target: "127.0.0.1",
	Type:   linodego.RecordTypeA,
	Name:   "a",
}

func TestDomainRecord_Create(t *testing.T) {
	_, _, record, teardown, err := setupDomainRecord(t, "fixtures/TestDomainRecord_Create")
	defer teardown()

	if err != nil {
		t.Errorf("Error creating domain record, got error %v", err)
	}

	expected := testDomainRecordCreateOpts

	// cant compare Target, fixture IPs are sanitized
	if record.Name != expected.Name || record.Type != expected.Type {
		t.Errorf("DomainRecord did not match CreateOptions")
	}
}

func TestDomainRecord_Update(t *testing.T) {
	client, domain, record, teardown, err := setupDomainRecord(t, "fixtures/TestDomainRecord_Update")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	updateOpts := linodego.DomainRecordUpdateOptions{
		Name: "renamed",
	}
	recordUpdated, err := client.UpdateDomainRecord(context.Background(), domain.ID, record.ID, updateOpts)
	if err != nil {
		t.Errorf("Error updating domain record, %s", err)
	}
	if recordUpdated.Name != "renamed" || record.Type != recordUpdated.Type || recordUpdated.Target != record.Target {
		t.Errorf("DomainRecord did not match UpdateOptions")
	}
}

func TestDomainRecords_List(t *testing.T) {
	client, domain, record, teardown, err := setupDomainRecord(t, "fixtures/TestDomainRecords_List")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	filter, err := json.Marshal(map[string]interface{}{
		"name": record.Name,
	})
	if err != nil {
		t.Error(err)
	}

	listOpts := linodego.NewListOptions(0, string(filter))
	records, err := client.ListDomainRecords(context.Background(), domain.ID, listOpts)
	if err != nil {
		t.Errorf("Error listing domains records, expected array, got error %v", err)
	}
	if len(records) != 1 {
		t.Errorf("Expected ListDomainRecords to match one result")
	}
}

func TestDomainRecords_ListMultiplePages(t *testing.T) {
	client, domain, record, teardown, err := setupDomainRecord(t, "fixtures/TestDomainRecords_ListMultiplePages")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	filter, err := json.Marshal(map[string]interface{}{
		"name": record.Name,
	})
	if err != nil {
		t.Error(err)
	}
	listOpts := linodego.NewListOptions(0, string(filter))
	records, err := client.ListDomainRecords(context.Background(), domain.ID, listOpts)
	if err != nil {
		t.Errorf("Error listing domains records, expected array, got error %v", err)
	}
	if len(records) != 1 {
		t.Errorf("Expected ListDomainRecords to match one result")
	}
}

func TestDomainRecord_Get(t *testing.T) {
	client, domain, record, teardown, err := setupDomainRecord(t, "fixtures/TestDomainRecord_Get")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	recordGot, err := client.GetDomainRecord(context.Background(), domain.ID, record.ID)
	if recordGot.Name != record.Name {
		t.Errorf("GetDomainRecord did not get the expected record")
	}
	if err != nil {
		t.Errorf("Error getting domain %d, got error %v", domain.ID, err)
	}
}

func setupDomainRecord(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.Domain, *linodego.DomainRecord, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, domain, fixtureTeardown, err := setupDomain(t, fixturesYaml)
	if err != nil {
		t.Errorf("Error creating domain, got error %v", err)
	}

	createOpts := testDomainRecordCreateOpts
	record, err := client.CreateDomainRecord(context.Background(), domain.ID, createOpts)
	if err != nil {
		t.Errorf("Error creating domain record, got error %v", err)
	}

	teardown := func() {
		// delete the DomainRecord to exercise the code
		if err := client.DeleteDomainRecord(context.Background(), domain.ID, record.ID); err != nil {
			t.Errorf("Expected to delete a domain record, but got %v", err)
		}
		fixtureTeardown()
	}
	return client, domain, record, teardown, err
}
