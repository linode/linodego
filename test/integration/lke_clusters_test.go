package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/linode/linodego/pkg/errors"
)

var (
	testLKEClusterCreateOpts = linodego.LKEClusterCreateOptions{
		Label:      label,
		Region:     "us-central",
		K8sVersion: "1.17",
		Tags:       []string{"testing"},
		NodePools:  []linodego.LKEClusterPoolCreateOptions{{Count: 1, Type: "g6-standard-2"}},
	}
)

func TestGetLKECluster_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetLKECluster_missing")
	defer teardown()

	i, err := client.GetLKECluster(context.Background(), 0)
	if err == nil {
		t.Errorf("should have received an error requesting a missing lkeCluster, got %v", i)
	}
	e, ok := err.(*errors.Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing lkeCluster, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing lkeCluster, got %v", e.Code)
	}
}

func TestWaitForLKEClusterReady(t *testing.T) {
	client, cluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
		createOpts.NodePools = []linodego.LKEClusterPoolCreateOptions{
			{Count: 3, Type: "g6-standard-2"},
		}
	}}, "fixtures/TestWaitForLKEClusterReady")
	defer teardown()

	wrapper, teardownClusterClient := transportRecorderWrapper(t, "fixtures/TestWaitForLKEClusterReadyClusterClient")
	defer teardownClusterClient()

	if err = client.WaitForLKEClusterReady(context.Background(), cluster.ID, linodego.LKEClusterPollOptions{
		TimeoutSeconds:   5 * 60,
		TransportWrapper: wrapper,
	}); err != nil {
		t.Errorf("Error waiting for the LKE cluster pools to be ready: %s", err)
	}
}

func TestGetLKECluster_found(t *testing.T) {
	client, lkeCluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
	}}, "fixtures/TestGetLKECluster_found")
	defer teardown()
	i, err := client.GetLKECluster(context.Background(), lkeCluster.ID)
	if err != nil {
		t.Errorf("Error getting lkeCluster, expected struct, got %v and error %v", i, err)
	}
	if i.ID != lkeCluster.ID {
		t.Errorf("Expected a specific lkeCluster, but got a different one %v", i)
	}
}

func TestListLKEClusterAPIEndpoints(t *testing.T) {
	client, lkeCluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
	}}, "fixtures/TestListLKEClusterAPIEndpoints")
	defer teardown()

	if err != nil {
		t.Error(err)
	}

	i, err := client.ListLKEClusterAPIEndpoints(context.Background(), lkeCluster.ID, nil)
	if err != nil {
		t.Errorf("Error listing lkeClusterAPIEndpoints, expected struct, got error %v", err)
	}
	if len(i) != 1 {
		t.Errorf("Expected a single lkeClusterAPIEndpoints, but got none %v", i)
	}
}

func TestGetLKEClusterKubeconfig(t *testing.T) {
	client, lkeCluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
	}}, "fixtures/TestGetLKEClusterKubeconfig")
	defer teardown()

	_, err = client.WaitForLKEClusterStatus(context.Background(), lkeCluster.ID, linodego.LKEClusterReady, 180)
	if err != nil {
		t.Errorf("Error waiting for LKECluster readiness: %s", err)
	}
	i, err := client.GetLKEClusterKubeconfig(context.Background(), lkeCluster.ID)
	if err != nil {
		t.Errorf("Error getting lkeCluster Kubeconfig, expected struct, got %v and error %v", i, err)
	}
	if len(i.KubeConfig) == 0 {
		t.Errorf("Expected an lkeCluster Kubeconfig, but got empty string %v", i)
	}
}

func TestListLKEClusters(t *testing.T) {
	client, _, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
	}}, "fixtures/TestListLKEClusters")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	// @TODO filter on the known label, API docs say this is supported, but it
	// errors
	i, err := client.ListLKEClusters(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing lkeClusters, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of lkeClusters, but got none %v", i)
	}
}

func TestGetLKEVersion_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetLKEVersion_missing")
	defer teardown()

	i, err := client.GetLKEVersion(context.Background(), "does-not-exist")
	if err == nil {
		t.Errorf("should have received an error requesting a missing version, got %v", i)
	}
	e, ok := err.(*errors.Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing version, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing version, got %v", e.Code)
	}
}

func TestGetLKEVersion_found(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetLKEVersion_found")
	defer teardown()

	i, err := client.GetLKEVersion(context.Background(), "1.16")
	if err != nil {
		t.Errorf("Error getting version, expected struct, got %v and error %v", i, err)
	}
	if i.ID != "1.16" {
		t.Errorf("Expected a specific version, but got a different one %v", i)
	}
}
func TestListLKEVersions(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestListLKEVersions")
	defer teardown()

	i, err := client.ListLKEVersions(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing versions, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of versions, but got none %v", i)
	}
}

type clusterModifier func(*linodego.LKEClusterCreateOptions)

func setupLKECluster(t *testing.T, clusterModifiers []clusterModifier, fixturesYaml string) (*linodego.Client, *linodego.LKECluster, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	createOpts := testLKEClusterCreateOpts
	for _, modifier := range clusterModifiers {
		modifier(&createOpts)
	}
	lkeCluster, err := client.CreateLKECluster(context.Background(), createOpts)
	if err != nil {
		t.Errorf("failed to create LKE cluster: %s", err)
	}

	teardown := func() {
		if err := client.DeleteLKECluster(context.Background(), lkeCluster.ID); err != nil {
			t.Errorf("failed to delete LKE cluster: %s", err)
		}
		fixtureTeardown()
	}
	return client, lkeCluster, teardown, err
}
