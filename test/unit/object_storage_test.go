package unit

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestObjectStorage_Cancel(t *testing.T) {
	client := createMockClient(t)

	desiredResponse := make(map[string]interface{})

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "/object-storage/cancel"),
		httpmock.NewJsonResponderOrPanic(200, &desiredResponse))

	err := client.CancelObjectStorage(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
