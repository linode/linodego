package golinode

import "testing"

func TestListStackscripts(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	stackscripts, err := client.ListStackscripts()
	if err != nil {
		t.Errorf("Error listing stackscripts, expected struct - error %v", err)
	}
	if len(stackscripts) == 0 {
		t.Errorf("Expected a list of public stackscripts - %v", stackscripts)
	}
}
