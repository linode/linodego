package linodego_test

import (
	"context"
	"testing"
)

func TestListClusters(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestListClusters")
	defer teardown()

	clusters, err := client.ListClusters(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing clusters, expected struct - error %v", err)
	}
	if len(clusters) == 0 {
		t.Errorf("Expected a list of clusters - %v", clusters)
	}
}
