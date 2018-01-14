package golinode

import (
	"time"
)

// LinodeSnapshotsPagedResponse represents a linode API response for listing
type LinodeSnapshotsPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeSnapshot
}

// LinodeSnapshot represents a linode backup snapshot
type LinodeSnapshot struct {
	CreatedStr string `json:"created"`
	UpdatedStr string `json:"updated"`

	ID       int
	Label    string
	Status   string
	Type     string
	Created  *time.Time `json:"-"`
	Updated  *time.Time `json:"-"`
	Finished string
	Configs  []string
}
