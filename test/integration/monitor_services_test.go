package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonitorServices_Get_smoke(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestMonitorServices_Get")
	defer teardown()

	// List the all the regsitered ACLP services
	monitorServicesClientList, listErr := client.ListMonitorServices(context.Background(), nil)
	if listErr != nil {
		t.Errorf("Error listing monitor services:%s", listErr)
	}

	// found := false
	for _, services := range monitorServicesClientList {
		validateServiceTypes(t, services)
	}

	// Get the details of the registered ACLP services based on serviceType
	monitorServicesClient, getErr := client.ListMonitorServiceByType(context.Background(), "dbaas", nil)
	if getErr != nil {
		t.Errorf("Error getting monitor services : %s", getErr)
	}

	found := false
	for _, element := range monitorServicesClient {
		if element.ServiceType == "dbaas" {
			found = true
		}
	}

	if !found {
		t.Errorf("Monitor service not found in list.")
	}
}

func TestMonitorServices_GetNotAllowedServiceType(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestMonitorNotAllowedServiceType_Get")
	defer teardown()

	_, getErr := client.ListMonitorServiceByType(context.Background(), "saas", nil)
	require.Error(t, getErr)
	assert.Contains(t, getErr.Error(), "[404] Not found")
}

func validateServiceTypes(
	t *testing.T,
	serviceType linodego.MonitorService,
) {
	require.NotEmpty(t, serviceType.ServiceType)
	require.NotEmpty(t, serviceType.Label)
}
