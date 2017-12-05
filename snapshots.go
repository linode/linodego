package golinode

// LinodeSnapshotsPagedResponse represents a linode API response for listing
type LinodeSnapshotsPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeSnapshot
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
