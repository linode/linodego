package golinode

// LinodeStackScript represents a linode stack script
type LinodeStackScript struct {
	ID                int
	Username          string
	Label             string
	Description       string
	Distributions     []*LinodeDistribution
	DeploymentsTotal  int
	DeploymentsActive int
	IsPublic          bool
	Created           string
	Updated           string
	RevNote           string
	UserDefinedFields *map[string]string
}

// LinodeStackScriptPager implements LinodeResponsePager
type LinodeStackScriptPager struct {
	Page, Pages int
	Data     []*LinodeStackScript
}

func (*p LinodeStackScriptPager) Results() interface{} {
	return p.Data
}

func (c *Client) ListStackScripts() (*LinodeResponsePager, err) {

}
