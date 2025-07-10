package integration

import (
	"testing"

	"github.com/linode/linodego"
	"golang.org/x/net/context"
)

func TestMonitorAPI_Fetch_Entity_Metrics(t *testing.T) {
	mClient, entityIDs, teardown, err := setup(t, "fixtures/TestMonitorAPI_Get_Entity_Metrics")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	opts := linodego.EntityMetricsFetchOptions{
		EntityIDs: entityIDs,
		Metrics: []linodego.EntityMetric{
			{
				Name:              "avg_read_iops",
				AggregateFunction: linodego.AggregateFunctionAvg,
			},
		},
		RelativeTimeDuration: linodego.MetricRelativeTimeDuration{
			Unit:  linodego.MetricTimeUnitHr,
			Value: 1,
		},
	}

	mClient.SetDebug(true)

	metrics, err := mClient.FetchEntityMetrics(context.Background(), "dbaas", &opts)
	if err != nil {
		t.Errorf("Error fetching the metrics for the entity: %s", err)
	}

	if len(metrics.Data.Result) < 1 {
		t.Errorf("No metric returned.")
	}
}

func setup(t *testing.T, fixturesYaml string) (*linodego.MonitorClient, []any, func(), error) {
	t.Helper()

	// create a DB entity to generate token
	client, _, teardown, err := setupPostgresDatabase(t, nil, fixturesYaml)
	if err != nil {
		t.Error(err)
	}

	dbs, err := client.ListDatabases(context.Background(), nil)
	if len(dbs) == 0 {
		t.Errorf("Expected a list of Databases, but got none: %v", err)
	}
	if err != nil {
		t.Errorf("Error listing Databases, expected struct, got error %v", err)
	}

	var entityIDs []any
	for _, db := range dbs {
		entityIDs = append(entityIDs, db.ID)
	}

	// Create a JWE token for the given entity IDs
	createOpts := linodego.MonitorTokenCreateOptions{
		EntityIDs: entityIDs,
	}

	token, createErr := client.CreateMonitorServiceTokenForServiceType(context.Background(), "dbaas", createOpts)
	if createErr != nil {
		t.Errorf("Error creating token : %s", createErr)
	}

	if token == nil {
		t.Errorf("Error generating token. Did not get token back.")
	}

	mClient, fixtureTeardown := createTestMonitorClient(t, fixturesYaml, token)

	td := func() {
		teardown()
		fixtureTeardown()
	}
	return mClient, entityIDs, td, err
}
