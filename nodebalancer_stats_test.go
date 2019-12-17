package linodego_test

import (
	"context"
	"testing"
	"time"
)

func TestGetNodeBalancerStats(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	client, nodebalancer, teardown, err := setupNodeBalancer(t, "fixtures/TestGetNodeBalancerStats")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	ticker := time.NewTicker(10 * time.Second)
	timer := time.NewTimer(570 * time.Second)
	defer ticker.Stop()

poll:
	for {
		select {
		case <-ticker.C:
			_, err = client.GetNodeBalancerStats(context.Background(), nodebalancer.ID)
			if err == nil { //stats are now returning
				break poll
			}
		case <-timer.C:
			t.Fatal("Error getting stats, polling timed out")
		}
	}
	currentTime := nodebalancer.Created
	currentYear := currentTime.Year()
	currentMonth := int(currentTime.Month())
	_, err = client.GetInstanceStatsByDate(
		context.Background(), nodebalancer.ID, currentYear, currentMonth)
	if err != nil {
		t.Errorf("Error getting stats by date, expected struct, got error %v", err)
	}
}
