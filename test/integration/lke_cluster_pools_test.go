package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/linode/linodego/pkg/errors"
)

var (
	testLKEClusterPoolCreateOpts = linodego.LKEClusterPoolCreateOptions{
		Type:  "g6-standard-2",
		Count: 1,
	}
)

func TestGetLKEClusterPool_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetLKEClusterPool_missing")
	defer teardown()

	i, err := client.GetLKEClusterPool(context.Background(), 0, 0)
	if err == nil {
		t.Errorf("should have received an error requesting a missing lkeClusterPool, got %v", i)
	}
	e, ok := err.(*errors.Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing lkeClusterPool, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing lkeClusterPool, got %v", e.Code)
	}
}

func TestGetLKEClusterPool_found(t *testing.T) {
	client, lkeCluster, pool, teardown, err := setupLKEClusterPool(t, "fixtures/TestGetLKEClusterPool_found")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	i, err := client.GetLKEClusterPool(context.Background(), lkeCluster.ID, pool.ID)
	if err != nil {
		t.Errorf("Error getting lkeClusterPool, expected struct, got %v and error %v", i, err)
	}
	if i.ID != pool.ID {
		t.Errorf("Expected a specific lkeClusterPool, but got a different one %v", i)
	}
}

func TestListLKEClusterPools(t *testing.T) {
	client, lkeCluster, _, teardown, err := setupLKEClusterPool(t, "fixtures/TestListLKEClusterPools")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	i, err := client.ListLKEClusterPools(context.Background(), lkeCluster.ID, nil)
	if err != nil {
		t.Errorf("Error listing lkeClusterPools, expected struct, got error %v", err)
	}
	if len(i) != 1 {
		t.Errorf("Expected a single lkeClusterPools, but got none %v", i)
	}
}

func setupLKEClusterPool(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.LKECluster, *linodego.LKEClusterPool, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, lkeCluster, fixtureTeardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
	}}, fixturesYaml)
	if err != nil {
		t.Errorf("Error creating lkeCluster, got error %v", err)
	}

	pools, err := client.ListLKEClusterPools(context.Background(), lkeCluster.ID, nil)
	if err != nil {
		t.Errorf("Error creating LKECluster Pool, got error %v", err)
	}
	if len(pools) == 0 {
		t.Errorf("Missing initial LKECluster Pool")
	}
	pool := &pools[0]
	teardown := func() {
		// delete the LKEClusterPool to exercise the code
		if err := client.DeleteLKEClusterPool(context.Background(), lkeCluster.ID, pool.ID); err != nil {
			t.Errorf("Expected to delete a LKECluster Pool, but got %v", err)
		}
		fixtureTeardown()
	}
	return client, lkeCluster, pool, teardown, err
}
