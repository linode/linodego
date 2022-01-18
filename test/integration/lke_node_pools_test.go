package integration

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/linode/linodego"
)

var testLKENodePoolCreateOpts = linodego.LKENodePoolCreateOptions{
	Type:  "g6-standard-2",
	Count: 2,
	Disks: []linodego.LKENodePoolDisk{
		{
			Size: 1000,
			Type: "ext4",
		},
	},
	Tags: []string{"testing"},
}

func TestGetLKENodePool_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetLKENodePool_missing")
	defer teardown()

	i, err := client.GetLKENodePool(context.Background(), 0, 0)
	if err == nil {
		t.Errorf("should have received an error requesting a missing LKENodePool, got %v", i)
	}
	e, ok := err.(*linodego.Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing LKENodePool, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing LKENodePool, got %v", e.Code)
	}
}

func TestGetLKENodePool_found(t *testing.T) {
	client, lkeCluster, pool, teardown, err := setupLKENodePool(t, "fixtures/TestGetLKENodePool_found")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	i, err := client.GetLKENodePool(context.Background(), lkeCluster.ID, pool.ID)
	if err != nil {
		t.Errorf("Error getting LKENodePool, expected struct, got %v and error %v", i, err)
	}

	if i.ID != pool.ID {
		t.Errorf("Expected a specific LKENodePool, but got a different one %v", i)
	}
	if i.Count != 2 {
		t.Errorf("expected count to be 2; got %d", i.Count)
	}
	if i.Type != "g6-standard-2" {
		t.Errorf("expected type to be g6-standard-2; got %s", i.Type)
	}
	if diff := cmp.Diff(linodego.LKENodePoolAutoscaler{
		Min:     2,
		Max:     2,
		Enabled: false,
	}, pool.Autoscaler); diff != "" {
		t.Errorf("unexpected autoscaler:\n%s", diff)
	}
	if diff := cmp.Diff([]string{"testing"}, i.Tags); diff != "" {
		t.Errorf("unexpected tags:\n%s", diff)
	}
}

func TestListLKENodePools(t *testing.T) {
	client, lkeCluster, _, teardown, err := setupLKENodePool(t, "fixtures/TestListLKENodePools")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	i, err := client.ListLKENodePools(context.Background(), lkeCluster.ID, nil)
	if err != nil {
		t.Errorf("Error listing LKENodePool, expected struct, got error %v", err)
	}
	if len(i) != 2 {
		t.Errorf("Expected two LKENodePool, but got %#v", i)
	}
}

func TestDeleteLKENodePoolNode(t *testing.T) {
	client, lkeCluster, nodePool, teardown, err := setupLKENodePool(t, "fixtures/TestDeleteLKENodePoolNode")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	linodes := nodePool.Linodes
	err = client.DeleteLKENodePoolNode(context.TODO(), lkeCluster.ID, linodes[0].ID)
	if err != nil {
		t.Errorf("failed to delete node %q: %s", linodes[0].ID, err)
	}

	nodePool, err = client.GetLKENodePool(context.TODO(), lkeCluster.ID, nodePool.ID)
	if err != nil {
		t.Errorf("failed to get updated node pool: %s", err)
	}

	if !(len(nodePool.Linodes) == 1 && nodePool.Linodes[0].ID == linodes[1].ID) {
		t.Errorf("expected node pool to have 1 linode (%s); got %v", linodes[1].ID, nodePool.Linodes)
	}
}

func TestUpdateLKENodePool(t *testing.T) {
	client, lkeCluster, nodePool, teardown, err := setupLKENodePool(t, "fixtures/TestUpdateLKENodePool")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	updatedAutoscaler := linodego.LKENodePoolAutoscaler{
		Enabled: true,
		Min:     2,
		Max:     5,
	}
	updatedTags := []string{}
	updated, err := client.UpdateLKENodePool(context.TODO(), lkeCluster.ID, nodePool.ID, linodego.LKENodePoolUpdateOptions{
		Count:      2,            // downsize
		Tags:       &updatedTags, // remove all tags
		Autoscaler: &updatedAutoscaler,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(updated.Tags) != 0 {
		t.Errorf("expected tags to be empty; got %v", updated.Tags)
	}
	if updated.Count != 2 {
		t.Errorf("expected count to be 2; got %d", updated.Count)
	}

	updatedTags = []string{"bar", "foo", "test"}
	updated, err = client.UpdateLKENodePool(context.TODO(), lkeCluster.ID, nodePool.ID, linodego.LKENodePoolUpdateOptions{
		Count: 3,            // upsize
		Tags:  &updatedTags, // repopulate tags
	})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(updatedTags, updated.Tags); diff != "" {
		t.Errorf("unexpected tags:\n%s", diff)
	}
	if diff := cmp.Diff(updatedAutoscaler, updated.Autoscaler); diff != "" {
		t.Errorf("unexpected autoscaler:\n%s", diff)
	}
	if updated.Count != 3 {
		t.Errorf("expected count to be 3; got %d", updated.Count)
	}
}

func setupLKENodePool(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.LKECluster, *linodego.LKENodePool, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, lkeCluster, fixtureTeardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
	}}, fixturesYaml)
	if err != nil {
		t.Errorf("Error creating lkeCluster, got error %v", err)
	}

	pool, err := client.CreateLKENodePool(context.Background(), lkeCluster.ID, testLKENodePoolCreateOpts)
	if err != nil {
		t.Errorf("Error creating LKE Node Pool, got error %v", err)
	}

	teardown := func() {
		// delete the LKENodePool to exercise the code
		if err := client.DeleteLKENodePool(context.Background(), lkeCluster.ID, pool.ID); err != nil {
			t.Errorf("Expected to delete a LKE Node Pool, but got %v", err)
		}
		fixtureTeardown()
	}
	return client, lkeCluster, pool, teardown, err
}
