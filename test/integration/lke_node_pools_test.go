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

func TestLKENodePool_GetMissing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLKENodePool_GetMissing")
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

func TestLKENodePool_GetFound(t *testing.T) {
	client, lkeCluster, pool, teardown, err := setupLKENodePool(t, "fixtures/TestLKENodePool_GetFound", nil)
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

func TestLKENodePools_List(t *testing.T) {
	client, lkeCluster, _, teardown, err := setupLKENodePool(t, "fixtures/TestLKENodePools_List", nil)
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

func TestLKENodePoolNode_Delete(t *testing.T) {
	client, lkeCluster, nodePool, teardown, err := setupLKENodePool(t, "fixtures/TestLKENodePoolNode_Delete", nil)
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

func TestLKENodePool_Update(t *testing.T) {
	client, lkeCluster, nodePool, teardown, err := setupLKENodePool(t, "fixtures/TestLKENodePool_Update", nil)
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
	if len(updated.Labels) != 0 {
		t.Errorf("expected labels to be empty; got %v", updated.Labels)
	}
	if len(updated.Taints) != 0 {
		t.Errorf("expected taints to be empty; got %v", updated.Taints)
	}

	updatedTags = []string{"bar", "foo", "test"}
	updatedLabels := linodego.LKENodePoolLabels{"foo": "bar"}
	updatedTaints := []linodego.LKENodePoolTaint{{
		Key:    "foo",
		Value:  "bar",
		Effect: linodego.LKENodePoolTaintEffectNoSchedule,
	}}
	updated, err = client.UpdateLKENodePool(context.TODO(), lkeCluster.ID, nodePool.ID, linodego.LKENodePoolUpdateOptions{
		Count:  3,              // upsize
		Tags:   &updatedTags,   // repopulate tags
		Labels: &updatedLabels, // set a label
		Taints: &updatedTaints, // set a taint
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
	if diff := cmp.Diff(updatedLabels, updated.Labels); diff != "" {
		t.Errorf("unexpected labels:\n%s", diff)
	}
	if diff := cmp.Diff(updatedTaints, updated.Taints); diff != "" {
		t.Errorf("unexpected taints:\n%s", diff)
	}

	if updated.Count != 3 {
		t.Errorf("expected count to be 3; got %d", updated.Count)
	}
}

func TestLKENodePool_CreateWithLabelsAndTaints(t *testing.T) {
	createOpts := &linodego.LKENodePoolCreateOptions{
		Labels: linodego.LKENodePoolLabels{"foo": "bar"},
		Taints: []linodego.LKENodePoolTaint{{
			Key:    "foo",
			Value:  "bar",
			Effect: linodego.LKENodePoolTaintEffectNoSchedule,
		}},
	}
	_, _, nodePool, teardown, err := setupLKENodePool(t, "fixtures/TestLKENodePool_CreateWithLabelsAndTaints", createOpts)
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	if diff := cmp.Diff(createOpts.Labels, nodePool.Labels); diff != "" {
		t.Errorf("unexpected labels:\n%s", diff)
	}
	if diff := cmp.Diff(createOpts.Taints, nodePool.Taints); diff != "" {
		t.Errorf("unexpected taints:\n%s", diff)
	}
}

func setupLKENodePool(t *testing.T, fixturesYaml string, nodePoolCreateOpts *linodego.LKENodePoolCreateOptions) (*linodego.Client, *linodego.LKECluster, *linodego.LKENodePool, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, lkeCluster, fixtureTeardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = "go-lke-test-def"
		if nodePoolCreateOpts != nil {
			createOpts.NodePools = []linodego.LKENodePoolCreateOptions{*nodePoolCreateOpts}
		}
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
