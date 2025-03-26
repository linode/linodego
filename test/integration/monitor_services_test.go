package integration

import (
	"context"
	"log"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestMonitorServices_Get_smoke(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestMonitorServices_Get")
	defer teardown()

	monitorServicesClientList, listErr := client.ListMonitorServices(context.Background(), nil)
	if listErr != nil {
		t.Errorf("Error listing monitor services:%s", listErr)
	}

	// found := false
	for _, services := range monitorServicesClientList {
		validateServiceTypes(t, services)
	}

	monitorServicesClient, getErr := client.GetMonitorServiceByType(context.Background(), "dbaas", nil)
	if getErr != nil {
		t.Errorf("Error getting monitor services : %s", getErr)
	}

	found := false
	for _, element := range monitorServicesClient {
		log.Printf("[INFO] element.ServiceType : %#+v", element.ServiceType)
		log.Printf("[INFO] monitorServicesClient.ServiceType : %#+v", monitorServicesClient)
		if element.ServiceType == "dbaas" {
			log.Printf("[WARN] event.Created is nil when API returned: %#+v", element.ServiceType)
			found = true
		}
	}

	if !found {
		t.Errorf("Monitor service not found in list.")
	}

}

func validateServiceTypes(
	t *testing.T,
	serviceType linodego.MonitorServices,
) {
	require.NotEmpty(t, serviceType.ServiceType)
	require.NotEmpty(t, serviceType.Label)
}
