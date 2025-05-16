package integration

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/linode/linodego"
	k8scondition "github.com/linode/linodego/k8s/pkg/condition"
)

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
		createOpts.Label = "go-lke-test-wait"
		createOpts.NodePools = []linodego.LKENodePoolCreateOptions{
			{Count: 3, Type: "g6-standard-1"},
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

func TestLKECluster_GetFound_smoke(t *testing.T) {
	client, lkeCluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = "go-lke-test-found"
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

func TestLKECluster_Enterprise_smoke(t *testing.T) {
	client, lkeCluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Tier = "enterprise"
		createOpts.Region = "us-lax"
		createOpts.K8sVersion = ""
	}}, "fixtures/TestLKECluster_Enterprise_smoke")
	defer teardown()
	i, err := client.GetLKECluster(context.Background(), lkeCluster.ID)
	if err != nil {
		t.Errorf("Error getting lkeCluster, expected struct, got %v and error %v", i, err)
	}
	if i.ID != lkeCluster.ID {
		t.Errorf("Expected a specific lkeCluster, but got a different one %v", i)
	}
	if i.Tier != "enterprise" {
		t.Errorf("Expected a lkeCluster to have enterprise tier")
	}
}

func TestLKECluster_Update(t *testing.T) {
	client, cluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = "go-lke-test-update"
		createOpts.K8sVersion = "1.31"
	}}, "fixtures/TestLKECluster_Update")
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	updatedTags := []string{"test=true"}
	updatedLabel := cluster.Label + "-updated"
	updatedK8sVersion := "1.32"

	updatedCluster, err := client.UpdateLKECluster(context.Background(), cluster.ID, linodego.LKEClusterUpdateOptions{
		Tags:       &updatedTags,
		Label:      updatedLabel,
		K8sVersion: updatedK8sVersion,
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

	// Update the LKE cluster to HA
	// This needs to be done in a separate API request from the K8s version upgrade
	isHA := true
	updatedControlPlane := &linodego.LKEClusterControlPlaneOptions{HighAvailability: &isHA}

	updatedCluster, err = client.UpdateLKECluster(context.Background(), cluster.ID, linodego.LKEClusterUpdateOptions{
		ControlPlane: updatedControlPlane,
	})
	if err != nil {
		t.Fatalf("failed to update LKE Cluster (%d): %s", cluster.ID, err)
	}

	if !reflect.DeepEqual(*updatedControlPlane.HighAvailability, updatedCluster.ControlPlane.HighAvailability) {
		t.Errorf("expected control plane to be updated to %#v; got %#v", updatedControlPlane, updatedCluster.ControlPlane)
	}
}

func TestLKECluster_Nodes_Recycle(t *testing.T) {
	client, cluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = "go-lke-test-recycle"
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
		createOpts.Label = "go-lke-test-apiend"
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
		createOpts.Label = "go-lke-test-kube-get"
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

func TestLKECluster_Kubeconfig_Delete(t *testing.T) {
	client, lkeCluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = "go-lke-test-kube-delete"
	}}, "fixtures/TestLKECluster_Kubeconfig_Delete")
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

	delete_err := client.DeleteLKEClusterKubeconfig(context.Background(), lkeCluster.ID)
	if err != nil {
		t.Errorf("Error deleting lkeCluster Kubeconfig, got error %v", delete_err)
	}
}

func TestLKECluster_Dashboard_Get(t *testing.T) {
	client, lkeCluster, teardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = "go-lke-test-dash"
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
		createOpts.Label = "go-lke-test-list"
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

	i, err := client.GetLKEVersion(context.Background(), "1.29")
	if err != nil {
		t.Errorf("Error getting version, expected struct, got %v and error %v", i, err)
	}

	if i.ID != "1.29" {
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

func TestLKECluster_APLEnabled_smoke(t *testing.T) {
	client, lkeCluster, teardown, err := setupLKECluster(t, []clusterModifier{
		func(createOpts *linodego.LKEClusterCreateOptions) {
			createOpts.Label = "go-lke-test-apl-enabled"
		},
		func(createOpts *linodego.LKEClusterCreateOptions) {
			createOpts.APLEnabled = true
		},
		func(createOpts *linodego.LKEClusterCreateOptions) {
			// NOTE: g6-dedicated-4 is the minimum APL-compatible Linode type
			createOpts.NodePools = []linodego.LKENodePoolCreateOptions{{Count: 3, Type: "g6-dedicated-4", Tags: []string{"test"}}}
		},
	},
		"fixtures/TestLKECluster_APLEnabled")
	defer teardown()

	expectedConsoleURL := fmt.Sprintf("https://console.lke%d.akamai-apl.net", lkeCluster.ID)
	consoleURL, err := client.GetLKEClusterAPLConsoleURL(context.Background(), lkeCluster.ID)
	if err != nil {
		t.Errorf("Error getting LKE APL console URL, expected string, got %v and error %v", consoleURL, err)
	}
	if consoleURL != expectedConsoleURL {
		t.Errorf("Expected an APL console URL %v, but got a different one %v", expectedConsoleURL, consoleURL)
	}

	expectedHealthCheckURL := fmt.Sprintf("https://auth.lke%d.akamai-apl.net/ready", lkeCluster.ID)
	healthCheckURL, err := client.GetLKEClusterAPLHealthCheckURL(context.Background(), lkeCluster.ID)
	if err != nil {
		t.Errorf("Error getting LKE APL health check URL, expected string, got %v and error %v", healthCheckURL, err)
	}
	if healthCheckURL != expectedHealthCheckURL {
		t.Errorf("Expected an APL health check URL %v, but got a different one %v", expectedHealthCheckURL, healthCheckURL)
	}
}

func TestLKETierVersion_ListAndGet(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestLKETierVersion_ListAndGet")
	defer teardown()

	testCases := []string{"standard", "enterprise"}

	for _, tier := range testCases {
		t.Run(fmt.Sprintf("Tier=%s", tier), func(t *testing.T) {
			versions, err := client.ListLKETierVersions(context.Background(), tier, nil)
			if err != nil {
				t.Fatalf("Error listing versions: %v", err)
			}

			if len(versions) == 0 {
				t.Fatalf("Expected a list of versions for tier %s, but got none", tier)
			}

			for _, version := range versions {
				if string(version.Tier) != tier {
					t.Errorf("Expected version tier %q, but got %q", tier, version.Tier)
				}
			}

			v, err := client.GetLKETierVersion(context.Background(), tier, versions[0].ID)
			if err != nil {
				t.Fatalf("Error getting version %s for tier %s: %v", versions[0].ID, tier, err)
			}

			if v.ID != versions[0].ID {
				t.Errorf("Expected version ID %q, but got %q", versions[0].ID, v.ID)
			}
		})
	}
}

type clusterModifier func(*linodego.LKEClusterCreateOptions)

func setupLKECluster(t *testing.T, clusterModifiers []clusterModifier, fixturesYaml string) (*linodego.Client, *linodego.LKECluster, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	createOpts := linodego.LKEClusterCreateOptions{
		Label:      label,
		Tier:       "standard", // default, can be overridden
		Tags:       []string{"testing"},
		Region:     "", // region will be resolved below
		K8sVersion: "", // will be resolved if empty
		NodePools: []linodego.LKENodePoolCreateOptions{{
			Count: 1,
			Type:  "g6-standard-2",
			Tags:  []string{"test"},
		}},
	}

	for _, modifier := range clusterModifiers {
		modifier(&createOpts)
	}

	if createOpts.Region == "" {
		createOpts.Region = getRegionsWithCaps(t, client, []string{"Kubernetes", "LA Disk Encryption"})[0]
	}

	if createOpts.K8sVersion == "" {
		createOpts.K8sVersion = getK8sVersion(t, client, createOpts.Tier)
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

func getK8sVersion(t *testing.T, client *linodego.Client, tier string) string {
	t.Helper()

	versions, err := client.ListLKETierVersions(context.Background(), tier, nil)
	if err != nil {
		t.Fatalf("Error listing versions for tier %q: %v", tier, err)
	}

	if len(versions) == 0 {
		t.Fatalf("Expected a list of versions for tier %q, but got none", tier)
	}

	return versions[0].ID
}
