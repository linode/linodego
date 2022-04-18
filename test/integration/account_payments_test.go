package integration

import (
	"context"
	"testing"

	. "github.com/linode/linodego"
)

func TestPayment_GetMissing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestPayment_GetMissing")
	defer teardown()

	i, err := client.GetPayment(context.Background(), -1)
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

func TestPayment_GetFound(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestPayment_GetFound")
	defer teardown()

	p, err := client.ListPayments(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing payments, expected struct, got error %v", err)
	}
	if len(p) == 0 {
		t.Errorf("Expected a list of payments, but got none %v", p)
	}

	client, teardown = createTestClient(t, "fixtures/TestGetPayment_found")
	defer teardown()

	i, err := client.GetPayment(context.Background(), p[0].ID)
	if err != nil {
		t.Errorf("Error getting payment, expected struct, got %v and error %v", i, err)
	}
	if i.ID != p[0].ID {
		t.Errorf("Expected a specific payment, but got a different one %v", i)
	}

	assertDateSet(t, i.Date)
}

func TestPayments_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestPayments_List")
	defer teardown()

	i, err := client.ListPayments(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing payments, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of payments, but got none %v", i)
	}
}
