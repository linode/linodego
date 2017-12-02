package golinode

// LinodeInstancesPagedResponse represents a linode API response for listing
type LinodeSnapshotsPagedResponse struct {
	Page, Pages, Results int
	data                 []*LinodeSnapshot
}

// Data returns data collection from paged response
func (r LinodeSnapshotsPagedResponse) Data() ([]*LinodeSnapshot, error) {
	return r.data, nil
}

// LinodeSnapshot represents a linode backup snapshot
type LinodeSnapshot struct {
	ID       int
	Label    string
	Status   string
	Type     string
	Created  string
	Updated  string
	Finished string
	Configs  []string
}
