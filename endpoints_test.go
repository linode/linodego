package golinode

import (
	"fmt"
	"testing"
)

func TestResourceEndpoint(t *testing.T) {
	apiKey := "MYFAKEAPIKEY"
	client, err := NewClient(&apiKey, nil)
	r, err := client.getResource("distributions")
	if err != nil {
		t.Errorf("Expected resource struct for distributions, got err %v", err)
	}
	e, err := r.Endpoint(nil)
	if err != nil {
		t.Errorf("Expected endpoint for resource distributions, got err %v", err)
	}
	if e != distributionsEndpoint {
		t.Errorf("Distributions endpoint did not match '%s'", distributionsEndpoint)
	}
}
func TestResourceTemplatedEndpointWithID(t *testing.T) {
	apiKey := "MYFAKEAPIKEY"
	client, err := NewClient(&apiKey, nil)
	backupID := 1234255
	r, err := client.getResource("backups")
	if err != nil {
		t.Errorf("Expected resource struct for backups, got err %v", err)
	}
	e, err := r.Endpoint(backupID)
	if err != nil {
		t.Errorf("Expected endpoint string, got err %v", err)
	}
	if e != fmt.Sprintf("linode/instances/%d/backups", backupID) {
		t.Errorf("Backups endpoint did not contain backup ID '%d'", backupID)
	}
	if e, err := r.Endpoint(nil); err == nil {
		t.Errorf("Sent nil backup ID and didn't get error when generating endpoint, got %v", e)
	}
}
