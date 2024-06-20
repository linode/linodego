package integration

import (
    "context"
    "testing"

    "github.com/linode/linodego"
)

func TestAccountSettings_Get(t *testing.T) {
    client, teardown := createTestClient(t, "fixtures/TestAccountSettings")
    defer teardown()

    settings, err := client.GetAccountSettings(context.Background())
    if err != nil {
        t.Fatalf("Error getting Account Settings, expected struct, got error %v", err)
    }

    if settings.BackupsEnabled != true {
        t.Fatalf("Expected BackupsEnabled to be true, got %v", settings.BackupsEnabled)
    }

    if settings.Managed != true {
        t.Fatalf("Expected Managed to be true, got %v", settings.Managed)
    }

    if settings.NetworkHelper != true {
        t.Fatalf("Expected NetworkHelper to be true, got %v", settings.NetworkHelper)
    }

    if settings.LongviewSubscription == nil || *settings.LongviewSubscription != "longview-3" {
        t.Fatalf("Expected LongviewSubscription to be 'longview-3', got %v", settings.LongviewSubscription)
    }

    if settings.ObjectStorage == nil || *settings.ObjectStorage != "active" {
        t.Fatalf("Expected ObjectStorage to be 'active', got %v", settings.ObjectStorage)
    }
}

func TestAccountSettings_Update(t *testing.T) {
    client, teardown := createTestClient(t, "fixtures/TestAccountSettings")
    defer teardown()

    opts := linodego.AccountSettingsUpdateOptions{
        BackupsEnabled:      Bool(false),
        LongviewSubscription: String("longview-10"),
        NetworkHelper:       Bool(false),
    }

    settings, err := client.UpdateAccountSettings(context.Background(), opts)
    if err != nil {
        t.Fatalf("Error updating Account Settings, expected struct, got error %v", err)
    }

    if settings.BackupsEnabled != false {
        t.Fatalf("Expected BackupsEnabled to be false, got %v", settings.BackupsEnabled)
    }

    if settings.NetworkHelper != false {
        t.Fatalf("Expected NetworkHelper to be false, got %v", settings.NetworkHelper)
    }

    if settings.LongviewSubscription == nil || *settings.LongviewSubscription != "longview-10" {
        t.Fatalf("Expected LongviewSubscription to be 'longview-10', got %v", settings.LongviewSubscription)
    }
}

func Bool(v bool) *bool { return &v }
func String(v string) *string { return &v }
