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

    // Need to wait a while for the new linode to boot and start returning stats
    time.Sleep(9 * time.Minute)

    _, err = client.GetInstanceStats(context.Background(), instance.ID)
    if err != nil {
        t.Errorf("Error getting stats, expected struct, got error %v", err)
    }

    currentTime := time.Now()
    currentYear := currentTime.Year()
    currentMonth := int(currentTime.Month())
    _, err = client.GetInstanceStatsByDate(
        context.Background(), instance.ID, currentYear, currentMonth)
    if err != nil {
        t.Errorf("Error getting stats, expected struct, got error %v", err)
    }
}
