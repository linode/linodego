//go:build ignore
// +build ignore

package integration

import (
	"context"
	"testing"

	. "github.com/linode/linodego"
)

func TestTemplate_GetMissing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestTemplate_GetMissing")
	defer teardown()

	i, err := client.GetTemplate(context.Background(), "does-not-exist")
	if err == nil {
		t.Errorf("should have received an error requesting a missing template, got %v", i)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing template, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing template, got %v", e.Code)
	}
}

func TestTemplate_GetFound(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestTemplate_GetFound")
	defer teardown()

	i, err := client.GetTemplate(context.Background(), "linode/ubuntu16.04lts")
	if err != nil {
		t.Errorf("Error getting template, expected struct, got %v and error %v", i, err)
	}
	if i.ID != "linode/ubuntu16.04lts" {
		t.Errorf("Expected a specific template, but got a different one %v", i)
	}
}

func TestTemplates_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestTemplates_List")
	defer teardown()

	i, err := client.ListTemplates(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing templates, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of templates, but got none %v", i)
	}
}
