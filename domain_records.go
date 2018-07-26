package linodego

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty"
)

// DomainRecord represents a DomainRecord object
type DomainRecord struct {
	ID       int
	Type     DomainRecordType
	Name     string
	Target   string
	Priority int
	Weight   int
	Port     int
	Service  *string
	Protocol *string
	TTLSec   int `json:"ttl_sec"`
	Tag      *string
}

type DomainRecordCreateOptions struct {
	Type     DomainRecordType `json:"type"`
	Name     string           `json:"name"`
	Target   string           `json:"target"`
	Priority int              `json:"priority"`
	Weight   int              `json:"weight"`
	Port     int              `json:"port"`
	Service  *string          `json:"service,omitempty"`
	Protocol *string          `json:"protocol,omitempty"`
	TTLSec   int              `json:"ttl_sec"`
	Tag      *string          `json:"tag,omitempty"`
}

type DomainRecordUpdateOptions DomainRecordCreateOptions

type DomainRecordType string

const (
	RecordTypeA     DomainRecordType = "A"
	RecordTypeAAAA  DomainRecordType = "AAAA"
	RecordTypeNS    DomainRecordType = "NS"
	RecordTypeMX    DomainRecordType = "MX"
	RecordTypeCNAME DomainRecordType = "CNAME"
	RecordTypeTXT   DomainRecordType = "TXT"
	RecordTypeSRV   DomainRecordType = "SRV"
	RecordTypePTR   DomainRecordType = "PTR"
	RecordTypeCAA   DomainRecordType = "CAA"
)

func (d DomainRecord) GetUpdateOptions() (du DomainRecordUpdateOptions) {
	du.Type = d.Type
	du.Name = d.Name
	du.Target = d.Target
	du.Priority = d.Priority
	du.Weight = du.Weight
	du.Port = du.Port
	du.Service = du.Service
	du.Protocol = du.Protocol
	du.TTLSec = du.TTLSec
	du.Tag = d.Tag
	return
}

// DomainRecordsPagedResponse represents a paginated DomainRecord API response
type DomainRecordsPagedResponse struct {
	*PageOptions
	Data []*DomainRecord
}

// endpoint gets the endpoint URL for InstanceConfig
func (DomainRecordsPagedResponse) endpointWithID(c *Client, id int) string {
	endpoint, err := c.DomainRecords.endpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// appendData appends DomainRecords when processing paginated DomainRecord responses
func (resp *DomainRecordsPagedResponse) appendData(r *DomainRecordsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// setResult sets the Resty response type of DomainRecord
func (DomainRecordsPagedResponse) setResult(r *resty.Request) {
	r.SetResult(DomainRecordsPagedResponse{})
}

// ListDomainRecords lists DomainRecords
func (c *Client) ListDomainRecords(ctx context.Context, opts *ListOptions) ([]*DomainRecord, error) {
	response := DomainRecordsPagedResponse{}
	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *DomainRecord) fixDates() *DomainRecord {
	// v.Created, _ = parseDates(v.CreatedStr)
	// v.Updated, _ = parseDates(v.UpdatedStr)
	return v
}

// GetDomainRecord gets the domainrecord with the provided ID
func (c *Client) GetDomainRecord(ctx context.Context, id string) (*DomainRecord, error) {
	e, err := c.DomainRecords.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, id)
	r, err := c.R(ctx).SetResult(&DomainRecord{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*DomainRecord), nil
}

// CreateDomainRecord creates a DomainRecord
func (c *Client) CreateDomainRecord(ctx context.Context, domainrecord *DomainRecordCreateOptions) (*DomainRecord, error) {
	var body string
	e, err := c.DomainRecords.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&DomainRecord{})

	if bodyData, err := json.Marshal(domainrecord); err == nil {
		body = string(bodyData)
	} else {
		return nil, NewError(err)
	}

	r, err := coupleAPIErrors(req.
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e))

	if err != nil {
		return nil, err
	}
	return r.Result().(*DomainRecord).fixDates(), nil
}

// UpdateDomainRecord updates the DomainRecord with the specified id
func (c *Client) UpdateDomainRecord(ctx context.Context, id int, domainrecord DomainRecordUpdateOptions) (*DomainRecord, error) {
	var body string
	e, err := c.DomainRecords.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, id)

	req := c.R(ctx).SetResult(&DomainRecord{})

	if bodyData, err := json.Marshal(domainrecord); err == nil {
		body = string(bodyData)
	} else {
		return nil, NewError(err)
	}

	r, err := coupleAPIErrors(req.
		SetBody(body).
		Put(e))

	if err != nil {
		return nil, err
	}
	return r.Result().(*DomainRecord).fixDates(), nil
}

// DeleteDomainRecord deletes the DomainRecord with the specified id
func (c *Client) DeleteDomainRecord(ctx context.Context, id int) error {
	e, err := c.DomainRecords.Endpoint()
	if err != nil {
		return err
	}
	e = fmt.Sprintf("%s/%d", e, id)

	if _, err := coupleAPIErrors(c.R(ctx).Delete(e)); err != nil {
		return err
	}

	return nil
}
