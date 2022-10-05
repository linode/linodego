package integration

import (
	"context"
	"testing"

	. "github.com/linode/linodego"
)

func TestTag_Create(t *testing.T) {
	client, instance, teardown, err := setupTaggedInstance(t, "fixtures/TestTag_Create")
	defer teardown()
	if err != nil {
		t.Errorf("Error creating test Instance: %s", err)
	}

	if instance.Tags[0] != "go-tag-test" {
		t.Errorf("should have created a tagged instance, got %v", instance.Tags)
	}

	updateOpts := instance.GetUpdateOptions()
	if updateOpts.Tags == nil {
		updateOpts.Tags = new([]string)
	}
	newTags := append(*updateOpts.Tags, "go-tag-test-bar")
	updateOpts.Tags = &newTags
	updateOpts.Backups = nil
	instance, err = client.UpdateInstance(context.Background(), instance.ID, updateOpts)
	if err != nil {
		t.Errorf("should have updated instance tags, got %q", err)
	}

	tag, err := client.CreateTag(context.Background(), TagCreateOptions{Label: "go-tag-test-foo", Linodes: []int{instance.ID}})
	if err != nil {
		t.Errorf("should have created a tag, got %q", err)
	}
	tags, err := client.ListTags(context.Background(), nil)
	if err != nil {
		t.Errorf("should have listed tags, got %q", err)
	}
	found := false
	for _, t := range tags {
		if t.Label == tag.Label {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("should have found created tag, %q", tag.Label)
	}
	x, err := client.ListTaggedObjects(context.Background(), "go-tag-test", nil)
	if err != nil {
		t.Errorf("should have listed tagged objects, got %q", err)
	}
	if len(x) == 0 || x[0].Type != "linode" || x[0].Data.(Instance).ID != instance.ID {
		t.Errorf("should have found instance in tagged objects list, got %v", x)
	}

	so, err := x.SortedObjects()
	if err != nil {
		t.Errorf("should have sorted tagged objects list, got %q", err)
	}

	if len(so.Instances) == 0 || so.Instances[0].ID != instance.ID {
		t.Errorf("should have found instance in sorted tagged objects list, got %v", so)
	}

	for _, tag := range []string{"go-tag-test", "go-tag-test-foo", "go-tag-test-bar"} {
		if err := client.DeleteTag(context.Background(), tag); err != nil {
			t.Error(err)
		}
	}
}

func TestTag_ListTaggedObjects_Missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestTag_ListTaggedObjects_Missing")
	defer teardown()

	i, err := client.ListTaggedObjects(context.Background(), "does-not-exist", nil)
	if err == nil {
		t.Errorf("should have received an error requesting a missing tag, got %v", i)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing tag, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing tag, got %v", e.Code)
	}
}

func setupTaggedInstance(t *testing.T, fixturesYaml string) (*Client, *Instance, func(), error) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	createOpts := InstanceCreateOptions{
		Label:  "go-ins-test-tag",
		Region: "us-southeast",
		Type:   "g6-nanode-1",
		Tags:   []string{"go-tag-test"},
	}
	instance, err := client.CreateInstance(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error creating test Instance: %s", err)
	}

	teardown := func() {
		if err := client.DeleteInstance(context.Background(), instance.ID); err != nil {
			t.Errorf("Error deleting test Instance: %s", err)
		}
		fixtureTeardown()
	}
	return client, instance, teardown, err
}
