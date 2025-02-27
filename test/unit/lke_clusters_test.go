package unit

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestLKECluster_Regenerate(t *testing.T) {
	client := createMockClient(t)

	requestData := linodego.LKEClusterRegenerateOptions{
		KubeConfig:   true,
		ServiceToken: false,
	}

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "clusters/1234/regenerate"),
		mockRequestBodyValidate(t, requestData, nil))

	if _, err := client.RegenerateLKECluster(context.Background(), 1234, requestData); err != nil {
		t.Fatal(err)
	}
}

func TestLKECluster_DeleteServiceToken(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, "clusters/1234/servicetoken"), httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteLKEClusterServiceToken(context.Background(), 1234); err != nil {
		t.Fatal(err)
	}
}

func TestLKECluster_KubeconfigDelete(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, "clusters/1234/kubeconfig"), httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteLKEClusterKubeconfig(context.Background(), 1234); err != nil {
		t.Fatal(err)
	}
}

func TestLKECluster_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_cluster_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("lke/clusters", fixtureData)

	clusters, err := base.Client.ListLKEClusters(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, clusters, 2)

	assert.Equal(t, 123, clusters[0].ID)
	assert.Equal(t, "test-cluster", clusters[0].Label)
	assert.Equal(t, "us-east", clusters[0].Region)
}

func TestLKECluster_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_cluster_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("lke/clusters/123", fixtureData)

	cluster, err := base.Client.GetLKECluster(context.Background(), 123)
	assert.NoError(t, err)
	assert.Equal(t, 123, cluster.ID)
	assert.Equal(t, "test-cluster", cluster.Label)
}

func TestLKECluster_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_cluster_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	createOptions := linodego.LKEClusterCreateOptions{
		Label:      "new-cluster",
		Region:     "us-west",
		K8sVersion: "1.22",
		Tags:       []string{"tag1"},
	}

	base.MockPost("lke/clusters", fixtureData)

	cluster, err := base.Client.CreateLKECluster(context.Background(), createOptions)
	assert.NoError(t, err)
	assert.Equal(t, "new-cluster", cluster.Label)
	assert.Equal(t, "us-west", cluster.Region)
}

func TestLKECluster_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_cluster_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	updateOptions := linodego.LKEClusterUpdateOptions{
		Label: "updated-cluster",
		Tags:  &[]string{"new-tag"},
	}

	base.MockPut("lke/clusters/123", fixtureData)

	cluster, err := base.Client.UpdateLKECluster(context.Background(), 123, updateOptions)
	assert.NoError(t, err)
	assert.Equal(t, "updated-cluster", cluster.Label)
}

func TestLKECluster_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("lke/clusters/123", nil)

	err := base.Client.DeleteLKECluster(context.Background(), 123)
	assert.NoError(t, err)
}

func TestLKECluster_GetKubeconfig(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_cluster_kubeconfig")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("lke/clusters/123/kubeconfig", fixtureData)

	kubeconfig, err := base.Client.GetLKEClusterKubeconfig(context.Background(), 123)
	assert.NoError(t, err)
	assert.NotEmpty(t, kubeconfig.KubeConfig)
}

func TestLKECluster_DeleteKubeconfig(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("lke/clusters/123/kubeconfig", nil)

	err := base.Client.DeleteLKEClusterKubeconfig(context.Background(), 123)
	assert.NoError(t, err)
}

func TestLKECluster_GetDashboard(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_cluster_dashboard")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("lke/clusters/123/dashboard", fixtureData)

	dashboard, err := base.Client.GetLKEClusterDashboard(context.Background(), 123)
	assert.NoError(t, err)
	assert.Equal(t, "https://dashboard.example.com", dashboard.URL)
}

func TestLKECluster_GetAPLConsoleURL(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_cluster_apl")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("lke/clusters/123", fixtureData)

	url, err := base.Client.GetLKEClusterAPLConsoleURL(context.Background(), 123)
	assert.NoError(t, err)
	assert.Equal(t, "https://console.lke123.akamai-apl.net", url)
}

func TestLKECluster_GetAPLHealthCheckURL(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_cluster_apl")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("lke/clusters/123", fixtureData)

	url, err := base.Client.GetLKEClusterAPLHealthCheckURL(context.Background(), 123)
	assert.NoError(t, err)
	assert.Equal(t, "https://auth.lke123.akamai-apl.net/ready", url)
}
