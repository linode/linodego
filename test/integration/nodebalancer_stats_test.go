package integration

import (
	"context"
	"testing"
	"time"

	"github.com/linode/linodego/pkg/errors"
)

func TestGetNodeBalancerStats(t *testing.T) {
	client, nodebalancer, teardown, err := setupNodeBalancer(t, "fixtures/TestGetNodeBalancerStats")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	ticker := time.NewTicker(10 * time.Second)
	timer := time.NewTimer(120 * time.Second)
	defer ticker.Stop()

poll:
	for {
		select {
		case <-ticker.C:
			_, err = client.GetNodeBalancerStats(context.Background(), nodebalancer.ID)
			if err != nil {
				// Possible that the call succeeded but that stats aren't available (HTTP: 4XX)
				if v, ok := err.(*errors.Error); ok {
					if v.Code == 400 && v.Message == "Stats are unavailable at this time." {
						break poll
					}
					// Otherwise, let's call it fatal
					t.Fatal(err)
				}
			}
			if err == nil { //stats are now returning
				break poll
			}
		case <-timer.C:
			t.Fatal("Error getting stats, polling timed out")
		}
	}
}
