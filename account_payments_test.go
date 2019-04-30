package linodego_test

import (
	"context"

	. "github.com/linode/linodego"

	"testing"
)

func TestGetPayment_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetPayment_missing")
	defer teardown()

	i, err := client.GetPayment(context.Background(), "does-not-exist")
	if err == nil {
		t.Errorf("should have received an error requesting a missing payment, got %v", i)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing payment, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing payment, got %v", e.Code)
	}
}

func TestGetPayment_found(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetPayment_found")
	defer teardown()

	i, err := client.GetPayment(context.Background(), "linode/ubuntu16.04lts")
	if err != nil {
		t.Errorf("Error getting payment, expected struct, got %v and error %v", i, err)
	}
	if i.ID != "linode/ubuntu16.04lts" {
		t.Errorf("Expected a specific payment, but got a different one %v", i)
	}
}
func TestListPayments(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestListPayments")
	defer teardown()

	i, err := client.ListPayments(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing payments, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of payments, but got none %v", i)
	}
}
