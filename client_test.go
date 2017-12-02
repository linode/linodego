package golinode

import (
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
)

const (
	debugAPI = false
)

func createTestClient(debug bool) (*Client, error) {
	fakeAPIKey := "NOTANAPIKEY"

	r, err := recorder.NewAsMode("test/fixtures", recorder.ModeReplaying, nil)
	if err != nil {
		return nil, err
	}
	defer r.Stop() // Make sure recorder is stopped once done with it

	c, err := NewClient(&fakeAPIKey, r)
	if err != nil {
		return nil, err
	}
	c.SetDebug(debug)
	return c, nil
}

func TestNewClient(t *testing.T) {
	testAPIKey := "NOTAREALAPIKEY"

	client, err := NewClient(&testAPIKey, nil)
	if err != nil {
		t.Error("Expected Client got error", err)
	}
	if client == nil {
		t.Error("Expected Client got nil")
	}
}

func TestNewClientErrors(t *testing.T) {
	testAPIKey := ""

	client, err := NewClient(&testAPIKey, nil)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if client != nil {
		t.Error("Expected error, got", client)
	}
}
