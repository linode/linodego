package golinode

// LinodeBackup represents a linode backup
type LinodeBackup struct {
	Enabled      bool
	Availability string
	Schedule     struct {
		Day    string
		Window string
	}
	LastBackup *LinodeSnapshot
	Disks      []*LinodeDisk
}
