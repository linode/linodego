package integration

import (
	"context"
	"github.com/linode/linodego"
	"strings"
	"sync/atomic"
	"testing"
)

func TestCache_RegionList(t *testing.T) {
	validateResult := func(r []linodego.Region, err error) {
		if err != nil {
			t.Fatal(err)
		}

		if len(r) == 0 {
			t.Fatalf("expected a list of regions - %v", r)
		}
	}

	client, teardown := createTestClient(t, "fixtures/TestCache_RegionList")
	defer teardown()

	// Collect request number
	totalRequests := int64(0)

	client.OnBeforeRequest(func(request *linodego.Request) error {
		if !strings.Contains(request.URL, "regions") {
			return nil
		}

		atomic.AddInt64(&totalRequests, 1)
		return nil
	})

	// First request (no cache)
	validateResult(client.ListRegions(context.Background(), nil))

	// Second request (cached)
	validateResult(client.ListRegions(context.Background(), nil))

	// Clear cache
	client.ClearCache()

	// Third request (non-cached)
	validateResult(client.ListRegions(context.Background(), nil))

	// Fourth request (cache disabled)
	client.UseCache(false)
	validateResult(client.ListRegions(context.Background(), nil))

	// Fifth request (cache disabled)
	validateResult(client.ListRegions(context.Background(), nil))

	// Validate request count
	if totalRequests != 4 {
		t.Fatalf("expected 4 requests, got %d", totalRequests)
	}
}
