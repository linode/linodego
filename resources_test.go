package golinode

import (
	"fmt"
	"testing"
)

func TestResourceEndpoint(t *testing.T) {
	apiKey := "MYFAKEAPIKEY"
	client, err := NewClient(&apiKey, nil)
	if err != nil {
		t.Error("Could not create new client in test")
	}
	r := client.Resource("distributions")
	e, err := r.Endpoint()
	if err != nil {
		t.Error("Got error when querying for distributions endpoint")
	}
	if e != distributionsEndpoint {
		t.Errorf("Distributions endpoint did not match '%s'", distributionsEndpoint)
	}
}
func TestResourceTemplatedEndpointWithID(t *testing.T) {
	apiKey := "MYFAKEAPIKEY"
	client, err := NewClient(&apiKey, nil)
	backupID := 1234255
	e, err := client.Backups.EndpointWithID(backupID)
	if err != nil {
		t.Error("Got error when getting endpoint with id for backups")
	}
	if e != fmt.Sprintf("linode/instances/%d/backups", backupID) {
		t.Errorf("Backups endpoint did not contain backup ID '%d'", backupID)
	}
}
