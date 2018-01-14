package golinode

import (
	"fmt"
	"time"
)

// LinodeStackscriptsPagedResponse represents a linode API response for listing
type LinodeStackscriptsPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeStackscript
}

// LinodeStackscript represents a linode stack script
type LinodeStackscript struct {
	CreatedStr string `json:"created"`
	UpdatedStr string `json:"updated"`

	ID                int
	Username          string
	Label             string
	Images            []string
	Description       string
	Distributions     []*LinodeDistribution
	DeploymentsTotal  int
	DeploymentsActive int
	IsPublic          bool
	Created           *time.Time `json:"-"`
	Updated           *time.Time `json:"-"`
	RevNote           string
	Script            string
	UserDefinedFields *map[string]string
	UserGravatarID    string
}

func (l *LinodeStackscript) fixDates() *LinodeStackscript {
	l.Created, _ = parseDates(l.CreatedStr)
	l.Updated, _ = parseDates(l.UpdatedStr)
	return l
}

// ListStackscripts gets all public stackscripts
func (c *Client) ListStackscripts() ([]*LinodeStackscript, error) {
	e, err := c.StackScripts.Endpoint()
	if err != nil {
		return nil, err
	}

	r, err := c.R().
		SetResult(&LinodeStackscriptsPagedResponse{}).
		Get(e)
	if err != nil {
		return nil, err
	}

	ss := r.Result().(*LinodeStackscriptsPagedResponse).Data
	for _, s := range ss {
		s.fixDates()
	}
	return ss, nil
}

// GetStackscript returns a stackscript with specified id
func (c *Client) GetStackscript(id int) (*LinodeStackscript, error) {
	e, err := c.StackScripts.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, id)

	r, err := c.R().
		SetResult(&LinodeStackscript{}).
		Get(e)

	if err != nil {
		return nil, err
	}
	d := r.Result().(*LinodeStackscript)
	return d, nil
}
