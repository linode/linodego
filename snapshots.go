package golinode

import (
	"time"
)

// LinodeSnapshot represents a linode backup snapshot
type LinodeSnapshot struct {
	CreatedStr  string `json:"created"`
	UpdatedStr  string `json:"updated"`
	FinishedStr string `json:"finished"`

	ID           int
	Label        string
	Status       string
	Type         string
	Region       string
	Created      *time.Time `json:"-"`
	Updated      *time.Time `json:"-"`
	Finished     *time.Time `json:"-"`
	Configs      []string
	Disks        []*LinodeInstanceDisk
	Availability string
}

func (l *LinodeSnapshot) fixDates() *LinodeSnapshot {
	l.Created, _ = parseDates(l.CreatedStr)
	l.Updated, _ = parseDates(l.UpdatedStr)
	l.Finished, _ = parseDates(l.FinishedStr)
	return l
}
