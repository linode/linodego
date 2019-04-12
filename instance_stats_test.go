package linodego_test

import (
	"fmt"
    "context"
	"testing"
    "time"
)


func TestGetInstanceStats(t *testing.T) {
    fmt.Println("Started TestGetInstanceStats")
    if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	client, instance, teardown, err := setupInstance(t, "fixtures/TestGetInstanceStats")
	defer teardown()
    if err != nil {
        t.Error(err)
    }

    err = client.BootInstance(context.Background(), instance.ID, 0)
    if err != nil {
        t.Error(err)
    }

    ticker := time.NewTicker(10 * time.Second)
    timer := time.NewTimer(570 * time.Second)
    defer ticker.Stop()

    // Test GetInstanceStats
    poll:
    for {
        select {
        case <-ticker.C:
            _, err = client.GetInstanceStats(context.Background(), instance.ID)
            if err == nil { // stats are now returning
                break poll
            }
        case <-timer.C:
            t.Fatal("Error getting stats, polling timed out")
        }
    }

    // test GetInstanceStatsByDate
    // No need to poll, since we know that if we get to this point,
    // the instance is returning stats
    currentTime := time.Now()
    currentYear := currentTime.Year()
    currentMonth := int(currentTime.Month())
    _, err = client.GetInstanceStatsByDate(
        context.Background(), instance.ID, currentYear, currentMonth)
    if err != nil {
        t.Errorf("Error getting stats by date, expected struct, got error %v", err)
    }
}
