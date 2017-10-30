package golinode

import "testing"

const (
	testAPIKey = "asdfasdfasdf"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(testAPIKey)
	if err != nil {
		t.Error("Expected Client got error", err)
	}
	if client == nil {
		t.Error("Expected Client got nil")
	}
}

func TestNewClientErrors(t *testing.T) {
	client, err := NewClient("")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if client != nil {
		t.Error("Expected error, got", client)
	}
}
