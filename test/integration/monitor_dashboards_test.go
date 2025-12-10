package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonitorDashboards_Get_smoke(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestMonitorDashboards_Get")
	defer teardown()

	// Get list of all ACLP Dashboards
	monitorDashboardsClientList, listErr := client.ListMonitorDashboards(context.Background(), nil)
	if listErr != nil {
		t.Errorf("Error listing monitor dashboards:%s", listErr)
	}

	// validating the details of the Dashboards
	for _, dashboards := range monitorDashboardsClientList {
		validateDashboards(t, dashboards)
	}

	// Get an ACLP Dashboard by dashboardID
	monitorDashhboardClient, getErr := client.GetMonitorDashboard(context.Background(), 1)
	if getErr != nil {
		t.Errorf("Error getting dashboard by ID :%s", getErr)
	}

	found := false
	for _, element := range monitorDashboardsClientList {
		if element.ServiceType == monitorDashhboardClient.ServiceType {
			found = true
		}
	}

	if !found {
		t.Errorf("Monitor dashboard not found in list.")
	}

	// List ACLP Dashboards by serviceType
	monitorDashhboardClientST, listErr := client.ListMonitorDashboardsByServiceType(context.Background(), "dbaas", nil)
	if listErr != nil {
		t.Errorf("Error listing monitor dashboards:%s", listErr)
	}

	found_st := false
	for _, element := range monitorDashhboardClientST {
		if element.ServiceType == monitorDashhboardClient.ServiceType {
			found_st = true
		}
	}

	if !found_st {
		t.Errorf("Monitor dashboard not found in list.")
	}
}

func TestMonitorDashboards_GetNotExistingDashboardID(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestMonitorInvalidDashboard_Get")
	defer teardown()

	_, getErr := client.GetMonitorDashboard(context.Background(), 999999)
	require.Error(t, getErr)
	assert.Contains(t, getErr.Error(), "[404] Not found")
}

func validateDashboards(
	t *testing.T,
	dashboards linodego.MonitorDashboard,
) {
	require.NotEmpty(t, dashboards.ID)
	require.NotEmpty(t, dashboards.ServiceType)
	require.NotEmpty(t, dashboards.Label)
	require.NotEmpty(t, dashboards.Created)
	require.NotEmpty(t, dashboards.Updated)
	require.NotEmpty(t, dashboards.Widgets)

	// Validate group_by and filters for each widget
	for _, widget := range dashboards.Widgets {
		require.NotNil(t, widget.GroupBy, "Expected group_by to be present in widget")
		require.GreaterOrEqual(t, len(widget.GroupBy), 0, "group_by should be a slice (possibly empty)")
		// filters is optional
		if widget.Filters != nil {
			require.GreaterOrEqual(t, len(widget.Filters), 0, "filters should be a slice (possibly empty)")
		}
	}
}
