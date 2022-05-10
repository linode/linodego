package integration

import (
	"context"
	"net/url"
	"reflect"
	"testing"

	"github.com/linode/linodego"
	k8scondition "github.com/linode/linodego/k8s/pkg/condition"
)

var testLKEClusterCreateOpts = linodego.LKEClusterCreateOptions{
	Label:      label,
	Region:     "us-southeast",
	K8sVersion: "1.23",
	Tags:       []string{"testing"},
	NodePools:  []linodego.LKENodePoolCreateOptions{{Count: 1, Type: "g6-standard-2", Tags: []string{"test"}}},
}

func TestLKECluster_GetMissing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLKECluster_GetMissing")
	defer teardown()

	i, err := client.GetLKECluster(context.Background(), 0)
	if err == nil {
		t.Errorf("should have received an error requesting a missing lkeCluster, got %v", i)
	}
	e, ok := err.(*linodego.Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing lkeCluster, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing lkeCluster, got %v", e.Code)
	}
}

func TestLKECluster_WaitForReady(t *testing.T) {
	client, cluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
		createOpts.NodePools = []linodego.LKENodePoolCreateOptions{
			{Count: 3, Type: "g6-standard-2"},
		}
	}}, "fixtures/TestLKECluster_WaitForReady")
	defer teardown()

	wrapper, teardownClusterClient := transportRecorderWrapper(t, "fixtures/TestLKECluster_WaitForReady_Cluster")
	defer teardownClusterClient()

	if err = k8scondition.WaitForLKEClusterReady(context.Background(), *client, cluster.ID, linodego.LKEClusterPollOptions{
		Retry:            true,
		TimeoutSeconds:   10 * 60,
		TransportWrapper: wrapper,
	}); err != nil {
		t.Errorf("Error waiting for the LKE cluster pools to be ready: %s", err)
	}
}

func TestLKECluster_GetFound(t *testing.T) {
	client, lkeCluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
	}}, "fixtures/TestLKECluster_GetFound")
	defer teardown()
	i, err := client.GetLKECluster(context.Background(), lkeCluster.ID)
	if err != nil {
		t.Errorf("Error getting lkeCluster, expected struct, got %v and error %v", i, err)
	}
	if i.ID != lkeCluster.ID {
		t.Errorf("Expected a specific lkeCluster, but got a different one %v", i)
	}
}

func TestLKECluster_Update(t *testing.T) {
	client, cluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
	}}, "fixtures/TestLKECluster_Update")
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	updatedTags := []string{"test=true"}
	updatedLabel := "new" + cluster.Label
	updatedK8sVersion := "1.23"
	updatedControlPlane := &linodego.LKEClusterControlPlane{HighAvailability: true}
	updatedCluster, err := client.UpdateLKECluster(context.TODO(), cluster.ID, linodego.LKEClusterUpdateOptions{
		Tags:         &updatedTags,
		Label:        updatedLabel,
		K8sVersion:   updatedK8sVersion,
		ControlPlane: updatedControlPlane,
	})
	if err != nil {
		t.Fatalf("failed to update LKE Cluster (%d): %s", cluster.ID, err)
	}

	if updatedCluster.Label != updatedLabel {
		t.Errorf("expected label to be updated to %q; got %q", updatedLabel, updatedCluster.Label)
	}

	if updatedCluster.K8sVersion != updatedK8sVersion {
		t.Errorf("expected k8s version to be updated to %q; got %q", updatedK8sVersion, updatedCluster.K8sVersion)
	}

	if !reflect.DeepEqual(updatedTags, updatedCluster.Tags) {
		t.Errorf("expected tags to be updated to %#v; got %#v", updatedTags, updatedCluster.Tags)
	}

	if !reflect.DeepEqual(*updatedControlPlane, updatedCluster.ControlPlane) {
		t.Errorf("expected control plane to be updated to %#v; got %#v", updatedControlPlane, updatedCluster.ControlPlane)
	}
}

func TestLKECluster_Nodes_Recycle(t *testing.T) {
	client, cluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
	}}, "fixtures/TestLKECluster_Nodes_Recycle")
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	err = client.RecycleLKEClusterNodes(context.TODO(), cluster.ID)
	if err != nil {
		t.Errorf("failed to recycle LKE cluster: %s", err)
	}
}

func TestLKECluster_APIEndpoints_List(t *testing.T) {
	client, lkeCluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
	}}, "fixtures/TestLKECluster_APIEndpoints_List")
	defer teardown()

	if err != nil {
		t.Error(err)
	}

	i, err := client.ListLKEClusterAPIEndpoints(context.Background(), lkeCluster.ID, nil)
	if err != nil {
		t.Errorf("Error listing lkeClusterAPIEndpoints, expected struct, got error %v", err)
	}
	if len(i) <= 0 {
		t.Errorf("Expected some lkeClusterAPIEndpoints, but got none %v", i)
	}
}

func TestLKECluster_Kubeconfig_Get(t *testing.T) {
	client, lkeCluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
	}}, "fixtures/TestLKECluster_Kubeconfig_Get")
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

func TestLKECluster_Dashboard_Get(t *testing.T) {
	client, lkeCluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
	}}, "fixtures/TestLKECluster_Dashboard_Get")
	defer teardown()

	_, err = client.WaitForLKEClusterStatus(context.Background(), lkeCluster.ID, linodego.LKEClusterReady, 180)
	if err != nil {
		t.Errorf("Error waiting for LKECluster readiness: %s", err)
	}
	i, err := client.GetLKEClusterDashboard(context.Background(), lkeCluster.ID)
	if err != nil {
		t.Errorf("Error getting LKE cluster dashboard URL, expected struct, got %v and error %v", i, err)
	}

	if len(i.URL) == 0 {
		t.Errorf("Expected an LKE cluster dashboard URL, but got empty string %v", i)
	}

	if _, err := url.ParseRequestURI(i.URL); err != nil {
		t.Errorf("invalid url: %s", err)
	}
}

func TestLKEClusters_List(t *testing.T) {
	client, _, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
	}}, "fixtures/TestLKEClusters_List")
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

func TestLKEVersion_GetMissing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLKEVersion_GetMissing")
	defer teardown()

	i, err := client.GetLKEVersion(context.Background(), "does-not-exist")
	if err == nil {
		t.Errorf("should have received an error requesting a missing version, got %v", i)
	}
	e, ok := err.(*linodego.Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing version, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing version, got %v", e.Code)
	}
}

func TestLKEVersion_GetFound(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLKEVersion_GetFound")
	defer teardown()

	i, err := client.GetLKEVersion(context.Background(), "1.23")
	if err != nil {
		t.Errorf("Error getting version, expected struct, got %v and error %v", i, err)
	}
	if i.ID != "1.23" {
		t.Errorf("Expected a specific version, but got a different one %v", i)
	}
}

func TestLKEVersions_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLKEVersions_List")
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
