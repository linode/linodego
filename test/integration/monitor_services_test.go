package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestMonitorServices_Get_smoke(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestMonitorServices_Get")
	defer teardown()

	MonitorServicesClientList, listErr := client.ListMonitorServices(context.Background(), nil)
	if listErr != nil {
		t.Errorf("Error listing monitor services:%s", listErr)
	}

	// found := false
	for _, services := range MonitorServicesClientList {
		validateServiceTypes(t, services)
	}
}

func validateServiceTypes(
	t *testing.T,
	serviceType linodego.MonitorServices,
) {
	require.NotEmpty(t, serviceType.ServiceType)
	require.NotEmpty(t, serviceType.Label)
}
