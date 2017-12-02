package golinode

// LinodeInstancesPagedResponse represents a linode API response for listing
type LinodeStackscriptsPagedResponse struct {
	Page, Pages, Results int
	data                 []*LinodeStackscript
}

// Data returns data collection from paged response
func (r LinodeStackscriptsPagedResponse) Data() ([]*LinodeStackscript, error) {
	return r.data, nil
}

// LinodeStackscript represents a linode stack script
type LinodeStackscript struct {
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
