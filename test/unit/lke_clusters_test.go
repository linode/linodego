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
	assert.Equal(t, 123, clusters[0].SubnetID)
	assert.Equal(t, 456, clusters[0].VpcID)
	assert.Equal(t, linodego.LKEClusterStackIPv4, clusters[0].StackType)
	assert.Equal(t, false, clusters[0].ControlPlane.AuditLogsEnabled)
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
	assert.Equal(t, 123, cluster.SubnetID)
	assert.Equal(t, 456, cluster.VpcID)
	assert.Equal(t, linodego.LKEClusterStackIPv4, cluster.StackType)
	assert.Equal(t, false, cluster.ControlPlane.AuditLogsEnabled)
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
		SubnetID:   linodego.Pointer(123),
		VpcID:      linodego.Pointer(456),
		StackType:  linodego.Pointer(linodego.LKEClusterStackIPv4),
		ControlPlane: &linodego.LKEClusterControlPlaneOptions{
			AuditLogsEnabled: linodego.Pointer(false),
		},
	}

	base.MockPost("lke/clusters", fixtureData)

	cluster, err := base.Client.CreateLKECluster(context.Background(), createOptions)
	assert.NoError(t, err)
	assert.Equal(t, "new-cluster", cluster.Label)
	assert.Equal(t, "us-west", cluster.Region)
	assert.Equal(t, 123, cluster.SubnetID)
	assert.Equal(t, 456, cluster.VpcID)
	assert.Equal(t, linodego.LKEClusterStackIPv4, cluster.StackType)
	assert.Equal(t, false, cluster.ControlPlane.AuditLogsEnabled)
}

func TestLKECluster_Create_Enterprise_RuleSetIDs(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("lke_cluster_enterprise_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	createOptions := linodego.LKEClusterCreateOptions{
		Label:      "enterprise-cluster",
		Region:     "us-east",
		K8sVersion: "1.31",
		Tier:       "enterprise",
		SubnetID:   linodego.Pointer(2010),
		VpcID:      linodego.Pointer(1010),
		StackType:  linodego.Pointer(linodego.LKEClusterDualStack),
		ControlPlane: &linodego.LKEClusterControlPlaneOptions{
			HighAvailability: linodego.Pointer(true),
			AuditLogsEnabled: linodego.Pointer(true),
		},
	}

	base.MockPost("lke/clusters", fixtureData)

	cluster, err := base.Client.CreateLKECluster(context.Background(), createOptions)
	assert.NoError(t, err)
	assert.Equal(t, 3010, cluster.ID)
	assert.Equal(t, "enterprise", cluster.Tier)
	assert.Equal(t, 2010, cluster.SubnetID)
	assert.Equal(t, 1010, cluster.VpcID)
	assert.Equal(t, linodego.LKEClusterDualStack, cluster.StackType)

	// Validate ruleset_ids deserialization
	assert.NotNil(t, cluster.RuleSetIDs, "RuleSetIDs should not be nil for enterprise clusters")
	assert.Equal(t, 4010, cluster.RuleSetIDs.Inbound)
	assert.Equal(t, 4011, cluster.RuleSetIDs.Outbound)
}

func TestLKECluster_Get_NoRuleSetIDs(t *testing.T) {
	// Standard clusters do not return ruleset_ids; the field should be nil
	fixtureData, err := fixtures.GetFixture("lke_cluster_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("lke/clusters/123", fixtureData)

	cluster, err := base.Client.GetLKECluster(context.Background(), 123)
	assert.NoError(t, err)
	assert.Nil(t, cluster.RuleSetIDs, "RuleSetIDs should be nil for standard clusters")
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
		ControlPlane: &linodego.LKEClusterControlPlaneOptions{
			AuditLogsEnabled: linodego.Pointer(true),
		},
	}

	base.MockPut("lke/clusters/123", fixtureData)

	cluster, err := base.Client.UpdateLKECluster(context.Background(), 123, updateOptions)
	assert.NoError(t, err)
	assert.Equal(t, "updated-cluster", cluster.Label)
	assert.Equal(t, true, cluster.ControlPlane.AuditLogsEnabled)
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
