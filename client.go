package golinode

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-resty/resty"
)

const (
	// APIHost Linode API hostname
	APIHost = "api.linode.com"
	// APIVersion Linode API version
	APIVersion = "v4"
	// APIProto connect to API with http(s)
	APIProto = "https"
	// Version of golinode
	Version = "1.0.0"
	// APIEnvVar environment var to check for API key
	APIEnvVar = "LINODE_API_KEY"
)

// Client is a wrapper around the Resty client
type Client struct {
	apiKey    string
	resty     *resty.Client
	resources map[string]*Resource

	Images        *Resource
	Disks         *Resource
	Configs       *Resource
	Backups       *Resource
	Instances     *Resource
	Regions       *Resource
	StackScripts  *Resource
	Volumes       *Resource
	Kernels       *Resource
	Types         *Resource
	Domains       *Resource
	Longview      *Resource
	NodeBalancers *Resource
	Support       *Resource
	Account       *Resource
	Profile       *Resource
	Managed       *Resource
}

// R wraps resty's R method
func (c *Client) R() *resty.Request {
	return c.resty.R()
}

// SetDebug sets the debug on resty's client
func (c *Client) SetDebug(debug bool) *Client {
	c.resty.SetDebug(debug)
	return c
}

// Resource looks up a resource by name
func (c Client) Resource(resourceName string) *Resource {
	selectedResource, ok := c.resources[resourceName]
	if !ok {
		log.Fatalf("Could not find resource named '%s', exiting.", resourceName)
	}
	return selectedResource
}

// PageOptions are the pagination parameters for List endpoints
type PageOptions struct {
	Page    int `url:"page,omitempty"`
	Pages   int `url:"per_page,omitempty"`
	Results int `url:"results,omitempty"`
}

// ListOptions are the pagination and filtering (TODO) parameters for endpoints
type ListOptions struct {
	*PageOptions
	Filter string
}

// NewClient factory to create new Client struct
func NewClient(codeAPIKey *string, transport http.RoundTripper) (*Client, error) {
	linodeAPIKey := ""

	if codeAPIKey != nil {
		linodeAPIKey = *codeAPIKey
	} else if envAPIKey, ok := os.LookupEnv(APIEnvVar); ok {
		linodeAPIKey = envAPIKey
	}

	if len(linodeAPIKey) == 0 || linodeAPIKey == "" {
		log.Print("Could not find LINODE_API_KEY, authenticated endpoints will fail.")
	}

	restyClient := resty.New().
		SetHostURL(fmt.Sprintf("%s://%s/%s", APIProto, APIHost, APIVersion)).
		SetAuthToken(linodeAPIKey).
		SetTransport(transport).
		SetHeader("User-Agent", fmt.Sprintf("go-linode %s https://github.com/chiefy/go-linode", Version))

	resources := map[string]*Resource{
		stackscriptsName:  NewResource(stackscriptsName, stackscriptsEndpoint, false),
		imagesName:        NewResource(imagesName, imagesEndpoint, false),
		instancesName:     NewResource(instancesName, instancesEndpoint, false),
		regionsName:       NewResource(regionsName, regionsEndpoint, false),
		disksName:         NewResource(disksName, disksEndpoint, true),
		configsName:       NewResource(configsName, configsEndpoint, true),
		backupsName:       NewResource(backupsName, backupsEndpoint, true),
		volumesName:       NewResource(volumesName, volumesEndpoint, false),
		kernelsName:       NewResource(kernelsName, kernelsEndpoint, false),
		typesName:         NewResource(typesName, typesEndpoint, false),
		domainsName:       NewResource(domainsName, domainsEndpoint, false),
		longviewName:      NewResource(longviewName, longviewEndpoint, false),
		nodebalancersName: NewResource(nodebalancersName, nodebalancersEndpoint, false),
		supportName:       NewResource(supportName, supportEndpoint, false),
		accountName:       NewResource(accountName, accountEndpoint, false),
		profileName:       NewResource(profileName, profileEndpoint, false),
		managedName:       NewResource(managedName, managedEndpoint, false),
	}

	return &Client{
		apiKey:    linodeAPIKey,
		resty:     restyClient,
		resources: resources,

		Images:        resources[imagesName],
		StackScripts:  resources[stackscriptsName],
		Instances:     resources[instancesName],
		Regions:       resources[regionsName],
		Disks:         resources[disksName],
		Configs:       resources[configsName],
		Backups:       resources[backupsName],
		Volumes:       resources[volumesName],
		Kernels:       resources[kernelsName],
		Types:         resources[typesName],
		Domains:       resources[domainsName],
		Longview:      resources[longviewName],
		NodeBalancers: resources[nodebalancersName],
		Support:       resources[supportName],
		Account:       resources[accountName],
		Profile:       resources[profileName],
		Managed:       resources[managedName],
	}, nil
}

type PagedResponse struct {
	ListResponse
	*PageOptions
}

type ListResponse interface {
	Endpoint(*Client) string
	AppendData(*resty.Response)
	SetResult(*resty.Request)
	ListHelper(*resty.Request, *ListOptions) error
}

// ListHelper abstracts fetching and pagination for GETmany endpoints
func (c *Client) ListHelper(i interface{}, opts *ListOptions) error {
	req := c.R()
	if opts != nil {
		req.SetQueryParam("page", strconv.Itoa(opts.Page))
	}

	var (
		err     error
		e       string
		pages   int
		results int
		r       *resty.Response
	)

	switch v := i.(type) {
	case LinodeKernelsPagedResponse:
		e = v.Endpoint(c)
		req.SetResult(v) // Can I just set PagedResponse instead of specific type?
		r, err = req.Get(e)
		if err != nil {
			return err
		}

		pages = r.Result().(*LinodeKernelsPagedResponse).Pages
		results = r.Result().(*LinodeKernelsPagedResponse).Results
		v.AppendData(r.Result().(*LinodeKernelsPagedResponse))
	case LinodeTypesPagedResponse:
		e = v.Endpoint(c)
		req.SetResult(v) // Can I just set PagedResponse instead of specific type?
		r, err = req.Get(e)
		if err != nil {
			return err
		}

		pages = r.Result().(*LinodeTypesPagedResponse).Pages
		results = r.Result().(*LinodeTypesPagedResponse).Results
		v.AppendData(r.Result().(*LinodeTypesPagedResponse))

	default:
		panic("what")
	}

	if opts == nil {
		for page := 2; page <= pages; page = page + 1 {
			c.ListHelper(i, &ListOptions{PageOptions: &PageOptions{Page: page}})
		}
	} else {
		opts.Results = results
		opts.Pages = pages
	}

	return nil
}
