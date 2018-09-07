package linodego_test

import (
	"context"
	"testing"
)

func TestListEvents_resizing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	client, teardown := createTestClient(t, "fixtures/TestListEvents_resizing")
	defer teardown()

	events, err := client.ListEvents(context.Background(), nil)
	if err != nil {
		t.Errorf("Error getting Events, expected struct, got error %v", err)
	}

	// TODO(displague) this test is dependent on specific fixture data, mock it here, or just test
	// fixDates directly
	if events[2].TimeRemaining != nil {
		t.Errorf("Error listing Events, expected resize event time_remaining to be nil, got %v", events[2].TimeRemaining)
	}

	if events[1].TimeRemaining == nil || *events[1].TimeRemaining != 0 {
		t.Errorf("Error listing Events, expected resize event time_remaining to be 0 seconds, got %v", events[1].TimeRemaining)
	}

	if events[0].TimeRemaining == nil || *events[0].TimeRemaining != 60+23 {
		t.Errorf("Error listing Events, expected resize event time_remaining to be 83 seconds, got %v", events[0].TimeRemaining)
	}
}
