package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestType_GetMissing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestType_GetMissing")
	defer teardown()

	i, err := client.GetType(context.Background(), "does-not-exist")
	if err == nil {
		t.Errorf("should have received an error requesting a missing image, got %v", i)
	}
	e, ok := err.(*linodego.Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing image, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing image, got %v", e.Code)
	}
}

func TestType_GetFound(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestType_GetFound")
	defer teardown()

	i, err := client.GetType(context.Background(), "g6-standard-1")
	if err != nil {
		t.Errorf("Error getting image, expected struct, got %v and error %v", i, err)
	}
	if i.ID != "g6-standard-1" {
		t.Errorf("Expected a specific image, but got a different one %v", i)
	}
}

func TestTypes_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestTypes_List")
	defer teardown()

	i, err := client.ListTypes(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing images, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of images, but got none %v", i)
	}
}

func TestTypes_RegionSpecific(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestTypes_RegionSpecific")
	defer teardown()

	validateOverride := func(override linodego.LinodeRegionPrice) {
		if override.ID == "" {
			t.Fatal("Expected region; got nil")
		}

		if override.Monthly <= 0 {
			t.Fatalf("Expected monthly cost; got %f", override.Monthly)
		}

		if override.Hourly <= 0 {
			t.Fatalf("Expected hourly cost; got %f", override.Hourly)
		}
	}

	types, err := client.ListTypes(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing images, expected struct, got error %v", err)
	}

	var targetType *linodego.LinodeType
	for _, t := range types {
		if t.RegionPrices != nil && len(t.RegionPrices) > 0 {
			targetType = &t
		}
	}

	if targetType == nil {
		t.Fatal("expected type with region override, got none")
	}

	// Validate overrides
	for _, override := range targetType.RegionPrices {
		validateOverride(override)
	}

	for _, override := range targetType.Addons.Backups.RegionPrices {
		validateOverride(override)
	}
}
