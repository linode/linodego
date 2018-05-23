package linodego

import (
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
)

var testingMode = recorder.ModeReplaying
var debugAPI = false
var validTestAPIKey = "NOTANAPIKEY"

func createTestClient(debug bool) (*Client, error) {
	var (
		fixturesYaml = "test/fixtures"
		c            Client
		apiKey       *string
	)

	apiKey = &validTestAPIKey

	if testing.Short() {
		apiKey = nil
		fixturesYaml = "test/fixtures_short"
	}

	r, err := recorder.NewAsMode(fixturesYaml, testingMode, nil)
	if err != nil {
		return nil, err
	}

	c = NewClient(apiKey, r)

	defer r.Stop() // Make sure recorder is stopped once done with it

	if err != nil {
		return nil, err
	}
	c.SetDebug(debug)
	return &c, nil
}

func TestClientAliases(t *testing.T) {
	client := NewClient(&validTestAPIKey, nil)

	if client.Images == nil {
		t.Error("Expected alias for Images to return a *Resource")
	}
	if client.Instances == nil {
		t.Error("Expected alias for Instances to return a *Resource")
	}
	if client.InstanceSnapshots == nil {
		t.Error("Expected alias for Backups to return a *Resource")
	}
	if client.StackScripts == nil {
		t.Error("Expected alias for StackScripts to return a *Resource")
	}
	if client.Regions == nil {
		t.Error("Expected alias for Regions to return a *Resource")
	}
	if client.Volumes == nil {
		t.Error("Expected alias for Volumes to return a *Resource")
	}
}
