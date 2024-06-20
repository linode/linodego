package integration

import (
    "context"
    "testing"
)

func TestAccountTransfer_Get(t *testing.T) {
    client, teardown := createTestClient(t, "fixtures/TestAccountTransfer_Get")
    defer teardown()

    transfer, err := client.GetAccountTransfer(context.Background())
    if err != nil {
        t.Fatalf("Error getting Account Transfer, expected struct, got error %v", err)
    }

    if transfer.Billable == 0 && transfer.Quota == 0 && transfer.Used == 0 {
        t.Fatalf("Expected non-zero values for Billable, Quota, and Used.")
    }

    if len(transfer.RegionTransfers) == 0 {
        t.Fatalf("Expected to see region transfers.")
    }

    for _, regionTransfer := range transfer.RegionTransfers {
        if regionTransfer.ID == "" {
            t.Errorf("Expected region ID to be non-empty.")
        }
        if regionTransfer.Billable == 0 && regionTransfer.Quota == 0 && regionTransfer.Used == 0 {
            t.Errorf("Expected non-zero values for Billable, Quota, and Used in region %s.", regionTransfer.ID)
        }
    }
}
