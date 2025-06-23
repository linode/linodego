package unit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestDomainRecord_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("domainrecord_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	domainID := 1234

	base.MockGet(formatMockAPIPath("domains/%d/records", domainID), fixtureData)

	domainsRecords, err := base.Client.ListDomainRecords(context.Background(), domainID, &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.Len(t, domainsRecords, 1)

	domainRecord := domainsRecords[0]

	assert.Equal(t, 123456, domainRecord.ID)
	assert.Equal(t, "test", domainRecord.Name)
	assert.Equal(t, 80, domainRecord.Port)
	assert.Equal(t, 50, domainRecord.Priority)
	assert.Nil(t, domainRecord.Protocol)
	assert.Nil(t, domainRecord.Service)
	assert.Nil(t, domainRecord.Tag)
	assert.Equal(t, "192.0.2.0", domainRecord.Target)
	assert.Equal(t, 604800, domainRecord.TTLSec)
	assert.Equal(t, linodego.DomainRecordType("A"), domainRecord.Type)
	assert.Equal(t, "2018-01-01T00:01:01Z", domainRecord.Created.Format(time.RFC3339))
	assert.Equal(t, "2018-01-01T00:01:01Z", domainRecord.Updated.Format(time.RFC3339))
	assert.Equal(t, 50, domainRecord.Weight)
}

func TestDomainRecord_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("domainrecord_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	domainID := 1234

	domainRecordID := 123456

	base.MockGet(formatMockAPIPath("domains/%d/records/%d", domainID, domainRecordID), fixtureData)

	domainRecord, err := base.Client.GetDomainRecord(context.Background(), domainID, domainRecordID)
	assert.NoError(t, err)

	assert.Equal(t, 123456, domainRecord.ID)
	assert.Equal(t, "test", domainRecord.Name)
	assert.Equal(t, 80, domainRecord.Port)
	assert.Equal(t, 50, domainRecord.Priority)
	assert.Nil(t, domainRecord.Protocol)
	assert.Nil(t, domainRecord.Service)
	assert.Nil(t, domainRecord.Tag)
	assert.Equal(t, "192.0.2.0", domainRecord.Target)
	assert.Equal(t, 604800, domainRecord.TTLSec)
	assert.Equal(t, linodego.DomainRecordType("A"), domainRecord.Type)
	assert.Equal(t, 50, domainRecord.Weight)
	assert.Equal(t, "2018-01-01T00:01:01Z", domainRecord.Created.Format(time.RFC3339))
	assert.Equal(t, "2018-01-01T00:01:01Z", domainRecord.Updated.Format(time.RFC3339))
}

func TestDomainRecord_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("domainrecord_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	domainID := 1234

	priority := 50
	weight := 50
	port := 80
	service := (*string)(nil)
	protocol := (*string)(nil)
	tag := (*string)(nil)
	ttlSec := 604800

	requestData := linodego.DomainRecordCreateOptions{
		Type:     linodego.RecordTypeA,
		Name:     "test",
		Target:   "192.0.2.0",
		Priority: &priority,
		Weight:   &weight,
		Port:     &port,
		Service:  service,
		Protocol: protocol,
		TTLSec:   &ttlSec,
		Tag:      tag,
	}

	base.MockPost(formatMockAPIPath("domains/%d/records", domainID), fixtureData)

	domainRecord, err := base.Client.CreateDomainRecord(context.Background(), domainID, requestData)
	assert.NoError(t, err)

	assert.Equal(t, 123456, domainRecord.ID)
	assert.Equal(t, "test", domainRecord.Name)
	assert.Equal(t, 80, domainRecord.Port)
	assert.Equal(t, 50, domainRecord.Priority)
	assert.Nil(t, domainRecord.Protocol)
	assert.Nil(t, domainRecord.Service)
	assert.Nil(t, domainRecord.Tag)
	assert.Equal(t, "192.0.2.0", domainRecord.Target)
	assert.Equal(t, 604800, domainRecord.TTLSec)
	assert.Equal(t, linodego.DomainRecordType("A"), domainRecord.Type)
	assert.Equal(t, 50, domainRecord.Weight)
	assert.Equal(t, "2018-01-01T00:01:01Z", domainRecord.Created.Format(time.RFC3339))
	assert.Equal(t, "2018-01-01T00:01:01Z", domainRecord.Updated.Format(time.RFC3339))
}

func TestDomainRecord_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("domainrecord_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	domainID := 1234

	priority := 50
	weight := 50
	port := 80
	service := (*string)(nil)
	protocol := (*string)(nil)
	tag := (*string)(nil)
	recordType := linodego.RecordTypeA
	name := "test"
	target := "192.0.2.0"
	ttlSec := 604800

	requestData := linodego.DomainRecordUpdateOptions{
		Type:     &recordType,
		Name:     &name,
		Target:   &target,
		Priority: &priority,
		Weight:   &weight,
		Port:     &port,
		Service:  service,
		Protocol: protocol,
		TTLSec:   &ttlSec,
		Tag:      tag,
	}

	domainRecordID := 123456

	base.MockPut(formatMockAPIPath("domains/%d/records/%d", domainID, domainRecordID), fixtureData)

	domainRecord, err := base.Client.UpdateDomainRecord(context.Background(), domainID, domainRecordID, requestData)
	assert.NoError(t, err)

	assert.Equal(t, 123456, domainRecord.ID)
	assert.Equal(t, "test", domainRecord.Name)
	assert.Equal(t, 80, domainRecord.Port)
	assert.Equal(t, 50, domainRecord.Priority)
	assert.Nil(t, domainRecord.Protocol)
	assert.Nil(t, domainRecord.Service)
	assert.Nil(t, domainRecord.Tag)
	assert.Equal(t, "192.0.2.0", domainRecord.Target)
	assert.Equal(t, 604800, domainRecord.TTLSec)
	assert.Equal(t, linodego.DomainRecordType("A"), domainRecord.Type)
	assert.Equal(t, 50, domainRecord.Weight)
	assert.Equal(t, "2018-01-01T00:01:01Z", domainRecord.Created.Format(time.RFC3339))
	assert.Equal(t, "2018-01-01T00:01:01Z", domainRecord.Updated.Format(time.RFC3339))
}

func TestDomainRecord_Delete(t *testing.T) {
	client := createMockClient(t)

	domainID := 1234

	domainRecordID := 123456

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, fmt.Sprintf("domains/%d/records/%d", domainID, domainRecordID)),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteDomainRecord(context.Background(), domainID, domainRecordID); err != nil {
		t.Fatal(err)
	}
}
