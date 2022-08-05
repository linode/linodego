package linodego

/**
 * Pagination and Filtering types and helpers
 */

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// PageOptions are the pagination parameters for List endpoints
type PageOptions struct {
	Page    int `url:"page,omitempty" json:"page"`
	Pages   int `url:"pages,omitempty" json:"pages"`
	Results int `url:"results,omitempty" json:"results"`
}

// ListOptions are the pagination and filtering (TODO) parameters for endpoints
type ListOptions struct {
	*PageOptions
	PageSize int
	Filter   string
}

// NewListOptions simplified construction of ListOptions using only
// the two writable properties, Page and Filter
func NewListOptions(page int, filter string) *ListOptions {
	return &ListOptions{PageOptions: &PageOptions{Page: page}, Filter: filter}
}

func applyListOptionsToRequest(opts *ListOptions, req *resty.Request) {
	if opts != nil {
		if opts.PageOptions != nil && opts.Page > 0 {
			req.SetQueryParam("page", strconv.Itoa(opts.Page))
		}

		if opts.PageSize > 0 {
			req.SetQueryParam("page_size", strconv.Itoa(opts.PageSize))
		}

		if len(opts.Filter) > 0 {
			req.SetHeader("X-Filter", opts.Filter)
		}
	}
}

type PagedResponse interface {
	endpoint(*Client) string
	castResult(*resty.Request, string) (int, int, error)
}

// listHelper abstracts fetching and pagination for GET endpoints that
// do not require any Ids (top level endpoints).
// When opts (or opts.Page) is nil, all pages will be fetched and
// returned in a single (endpoint-specific)PagedResponse
// opts.results and opts.pages will be updated from the API response
// nolint
func (c *Client) listHelper(ctx context.Context, pager PagedResponse, opts *ListOptions) error {
	req := c.R(ctx)
	applyListOptionsToRequest(opts, req)

	pages, results, err := pager.castResult(req, pager.endpoint(c))
	if err != nil {
		return err
	}
	if opts != nil {
		if opts.PageOptions == nil {
			opts.PageOptions = &PageOptions{}
		}
		if opts.Page == 0 {
			for page := 2; page <= pages; page++ {
				opts.Page = page
				if err := c.newListHelper(ctx, pager, opts); err != nil {
					return err
				}
			}
		}
		opts.Results = results
		opts.Pages = pages
	}
	for page := 2; page <= pages; page++ {
		newOpts := ListOptions{PageOptions: &PageOptions{Page: page}}
		if err := c.newListHelper(ctx, pager, &newOpts); err != nil {
			return err
		}
	}

	return nil
}

/**
case *ImagesPagedResponse:
case *StackscriptsPagedResponse:
case *InstancesPagedResponse:
case *RegionsPagedResponse:
case *VolumesPagedResponse:
case *DatabasesPagedResponse:
case *DatabaseEnginesPagedResponse:
case *DatabaseTypesPagedResponse:
case *MySQLDatabasesPagedResponse:
case *MongoDatabasesPagedResponse:
case *PostgresDatabasesPagedResponse:
case *DomainsPagedResponse:
case *EventsPagedResponse:
case *FirewallsPagedResponse:
case *LKEClustersPagedResponse:
case *LKEVersionsPagedResponse:
case *LongviewSubscriptionsPagedResponse:
case *LongviewClientsPagedResponse:
case *IPAddressesPagedResponse:
case *IPv6PoolsPagedResponse:
case *IPv6RangesPagedResponse:
case *SSHKeysPagedResponse:
case *TicketsPagedResponse:
case *InvoicesPagedResponse:
case *NotificationsPagedResponse:
case *OAuthClientsPagedResponse:
case *PaymentsPagedResponse:
case *NodeBalancersPagedResponse:
case *TagsPagedResponse:
case *TokensPagedResponse:
case *UsersPagedResponse:
case *ObjectStorageBucketsPagedResponse:
case *ObjectStorageClustersPagedResponse:
case *ObjectStorageKeysPagedResponse:
case *VLANsPagedResponse:

case ProfileAppsPagedResponse:
case ProfileWhitelistPagedResponse:
case ManagedContactsPagedResponse:
case ManagedCredentialsPagedResponse:
case ManagedIssuesPagedResponse:
case ManagedLinodeSettingsPagedResponse:
case ManagedServicesPagedResponse:
**/

// listHelperWithID abstracts fetching and pagination for GET endpoints that
// require an Id (second level endpoints).
// When opts (or opts.Page) is nil, all pages will be fetched and
// returned in a single (endpoint-specific)PagedResponse
// opts.results and opts.pages will be updated from the API response
// nolint
func (c *Client) listHelperWithID(ctx context.Context, i interface{}, idRaw interface{}, opts *ListOptions) error {
	var (
		err     error
		pages   int
		results int
		r       *resty.Response
	)

	req := c.R(ctx)
	applyListOptionsToRequest(opts, req)

	id, _ := idRaw.(int)

	switch v := i.(type) {
	case *DomainRecordsPagedResponse:
		if r, err = coupleAPIErrors(req.SetResult(DomainRecordsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			response, ok := r.Result().(*DomainRecordsPagedResponse)
			if !ok {
				return fmt.Errorf("response is not a *DomainRecordsPagedResponse")
			}
			pages = response.Pages
			results = response.Results
			v.appendData(response)
		}
	case *FirewallDevicesPagedResponse:
		if r, err = coupleAPIErrors(req.SetResult(FirewallDevicesPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*FirewallDevicesPagedResponse).Pages
			results = r.Result().(*FirewallDevicesPagedResponse).Results
			v.appendData(r.Result().(*FirewallDevicesPagedResponse))
		}
	case *InstanceConfigsPagedResponse:
		if r, err = coupleAPIErrors(req.SetResult(InstanceConfigsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*InstanceConfigsPagedResponse).Pages
			results = r.Result().(*InstanceConfigsPagedResponse).Results
			v.appendData(r.Result().(*InstanceConfigsPagedResponse))
		}
	case *InstanceDisksPagedResponse:
		if r, err = coupleAPIErrors(req.SetResult(InstanceDisksPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*InstanceDisksPagedResponse).Pages
			results = r.Result().(*InstanceDisksPagedResponse).Results
			v.appendData(r.Result().(*InstanceDisksPagedResponse))
		}
	case *InstanceVolumesPagedResponse:
		if r, err = coupleAPIErrors(req.SetResult(InstanceVolumesPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*InstanceVolumesPagedResponse).Pages
			results = r.Result().(*InstanceVolumesPagedResponse).Results
			v.appendData(r.Result().(*InstanceVolumesPagedResponse))
		}
	case *InvoiceItemsPagedResponse:
		if r, err = coupleAPIErrors(req.SetResult(InvoiceItemsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*InvoiceItemsPagedResponse).Pages
			results = r.Result().(*InvoiceItemsPagedResponse).Results
			v.appendData(r.Result().(*InvoiceItemsPagedResponse))
		}
	case *LKEClusterAPIEndpointsPagedResponse:
		if r, err = coupleAPIErrors(req.SetResult(LKEClusterAPIEndpointsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*LKEClusterAPIEndpointsPagedResponse).Pages
			results = r.Result().(*LKEClusterAPIEndpointsPagedResponse).Results
			v.appendData(r.Result().(*LKEClusterAPIEndpointsPagedResponse))
		}
	case *LKENodePoolsPagedResponse:
		if r, err = coupleAPIErrors(req.SetResult(LKENodePoolsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*LKENodePoolsPagedResponse).Pages
			results = r.Result().(*LKENodePoolsPagedResponse).Results
			v.appendData(r.Result().(*LKENodePoolsPagedResponse))
		}
	case *MySQLDatabaseBackupsPagedResponse:
		if r, err = coupleAPIErrors(req.SetResult(MySQLDatabaseBackupsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*MySQLDatabaseBackupsPagedResponse).Pages
			results = r.Result().(*MySQLDatabaseBackupsPagedResponse).Results
			v.appendData(r.Result().(*MySQLDatabaseBackupsPagedResponse))
		}
	case *MongoDatabaseBackupsPagedResponse:
		if r, err = coupleAPIErrors(req.SetResult(MongoDatabaseBackupsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*MongoDatabaseBackupsPagedResponse).Pages
			results = r.Result().(*MongoDatabaseBackupsPagedResponse).Results
			v.appendData(r.Result().(*MongoDatabaseBackupsPagedResponse))
		}
	case *PostgresDatabaseBackupsPagedResponse:
		if r, err = coupleAPIErrors(req.SetResult(PostgresDatabaseBackupsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*PostgresDatabaseBackupsPagedResponse).Pages
			results = r.Result().(*PostgresDatabaseBackupsPagedResponse).Results
			v.appendData(r.Result().(*PostgresDatabaseBackupsPagedResponse))
		}
	case *NodeBalancerConfigsPagedResponse:
		if r, err = coupleAPIErrors(req.SetResult(NodeBalancerConfigsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*NodeBalancerConfigsPagedResponse).Pages
			results = r.Result().(*NodeBalancerConfigsPagedResponse).Results
			v.appendData(r.Result().(*NodeBalancerConfigsPagedResponse))
		}
	case *TaggedObjectsPagedResponse:
		idStr := idRaw.(string)

		if r, err = coupleAPIErrors(req.SetResult(TaggedObjectsPagedResponse{}).Get(v.endpointWithID(c, idStr))); err == nil {
			pages = r.Result().(*TaggedObjectsPagedResponse).Pages
			results = r.Result().(*TaggedObjectsPagedResponse).Results
			v.appendData(r.Result().(*TaggedObjectsPagedResponse))
		}
	/**
	case TicketAttachmentsPagedResponse:
		if r, err = req.SetResult(v).Get(v.endpoint(c)); r.Error() != nil {
			return NewError(r)
		} else if err == nil {
			pages = r.Result().(*TicketAttachmentsPagedResponse).Pages
			results = r.Result().(*TicketAttachmentsPagedResponse).Results
			v.appendData(r.Result().(*TicketAttachmentsPagedResponse))
		}
	case TicketRepliesPagedResponse:
		if r, err = req.SetResult(v).Get(v.endpoint(c)); r.Error() != nil {
			return NewError(r)
		} else if err == nil {
			pages = r.Result().(*TicketRepliesPagedResponse).Pages
			results = r.Result().(*TicketRepliesPagedResponse).Results
			v.appendData(r.Result().(*TicketRepliesPagedResponse))
		}
	**/
	default:
		log.Fatalf("Unknown listHelperWithID interface{} %T used", i)
	}

	if err != nil {
		return err
	}

	if opts == nil {
		for page := 2; page <= pages; page++ {
			if err := c.listHelperWithID(ctx, i, id, &ListOptions{PageOptions: &PageOptions{Page: page}}); err != nil {
				return err
			}
		}
	} else {
		if opts.PageOptions == nil {
			opts.PageOptions = &PageOptions{}
		}
		if opts.Page == 0 {
			for page := 2; page <= pages; page++ {
				opts.Page = page
				if err := c.listHelperWithID(ctx, i, id, opts); err != nil {
					return err
				}
			}
		}
		opts.Results = results
		opts.Pages = pages
	}

	return nil
}

// listHelperWithTwoIDs abstracts fetching and pagination for GET endpoints that
// require twos IDs (third level endpoints).
// When opts (or opts.Page) is nil, all pages will be fetched and
// returned in a single (endpoint-specific)PagedResponse
// opts.results and opts.pages will be updated from the API response
// nolint
func (c *Client) listHelperWithTwoIDs(ctx context.Context, i interface{}, firstID, secondID int, opts *ListOptions) error {
	var (
		err     error
		pages   int
		results int
		r       *resty.Response
	)

	req := c.R(ctx)
	applyListOptionsToRequest(opts, req)

	switch v := i.(type) {
	case *NodeBalancerNodesPagedResponse:
		if r, err = coupleAPIErrors(req.SetResult(NodeBalancerNodesPagedResponse{}).Get(v.endpointWithTwoIDs(c, firstID, secondID))); err == nil {
			pages = r.Result().(*NodeBalancerNodesPagedResponse).Pages
			results = r.Result().(*NodeBalancerNodesPagedResponse).Results
			v.appendData(r.Result().(*NodeBalancerNodesPagedResponse))
		}
	default:
		log.Fatalf("Unknown listHelperWithTwoIDs interface{} %T used", i)
	}

	if err != nil {
		return err
	}

	if opts == nil {
		for page := 2; page <= pages; page++ {
			if err := c.listHelper(ctx, i, &ListOptions{PageOptions: &PageOptions{Page: page}}); err != nil {
				return err
			}
		}
	} else {
		if opts.PageOptions == nil {
			opts.PageOptions = &PageOptions{}
		}
		if opts.Page == 0 {
			for page := 2; page <= pages; page++ {
				opts.Page = page
				if err := c.listHelperWithTwoIDs(ctx, i, firstID, secondID, opts); err != nil {
					return err
				}
			}
		}
		opts.Results = results
		opts.Pages = pages
	}

	return nil
}
