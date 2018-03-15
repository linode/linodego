package golinode

import (
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
)

const (
	debugAPI = false
)

var validTestAPIKey = "NOTANAPIKEY"

func createTestClient(debug bool) (*Client, error) {
	r, err := recorder.NewAsMode("test/fixtures", recorder.ModeReplaying, nil)
	if err != nil {
		return nil, err
	}
	defer r.Stop() // Make sure recorder is stopped once done with it

	c, err := NewClient(&validTestAPIKey, r)
	if err != nil {
		return nil, err
	}
	c.SetDebug(debug)
	return c, nil
}

func TestNewClient(t *testing.T) {
	client, err := NewClient(&validTestAPIKey, nil)
	if err != nil {
		t.Error("Expected Client got error", err)
	}
	if client == nil {
		t.Error("Expected Client got nil")
	}
}

func TestNewClientErrors(t *testing.T) {
	blankAPIKey := ""

	client, err := NewClient(&blankAPIKey, nil)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if client != nil {
		t.Error("Expected error, got", client)
	}
}

func TestClientAliases(t *testing.T) {
	client, err := NewClient(&validTestAPIKey, nil)
	if err != nil {
		t.Error("Expected client got error", err)
	}
	if client.Images == nil {
		t.Error("Expected alias for Distributions to return a *Resource")
	}
	if client.Instances == nil {
		t.Error("Expected alias for Instances to return a *Resource")
	}
	if client.Backups == nil {
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
