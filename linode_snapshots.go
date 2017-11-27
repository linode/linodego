package golinode

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
