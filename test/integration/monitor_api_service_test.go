package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
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
				Name:              "read_iops",
				AggregateFunction: linodego.AggregateFunctionAvg,
			},
		},
		RelativeTimeDuration: &linodego.MetricRelativeTimeDuration{
			Unit:  linodego.MetricTimeUnitHr,
			Value: 1,
		},
	}

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

	client, clientFixtureTeardown := createTestClient(t, "fixtures/TestMonitorAPI_Get_Entity_Metrics_ListDB")

	dbs, err := client.ListDatabases(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing Databases, expected struct, got error %v", err)
	}

	var teardown func()
	if len(dbs) < 1 {
		// create a DB entity to generate token
		client, _, teardown, err = setupPostgresDatabase(t, nil, "fixtures/TestMonitorAPI_Get_Entity_Metrics_setupPostgres")
		if err != nil {
			t.Error(err)
		}
	}

	// refresh the DB list
	dbs, err = client.ListDatabases(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing Databases, expected struct, got error %v", err)
	}
	if len(dbs) == 0 {
		t.Errorf("Expected a list of Databases, but got none: %v", err)
	}

	var entityIDs []any
	for _, db := range dbs {
		// DB entities must from the same region
		if db.Region == dbs[0].Region {
			entityIDs = append(entityIDs, db.ID)
		}
	}

	// Create a JWE token for the given entity IDs
	createOpts := linodego.MonitorTokenCreateOptions{
		EntityIDs: entityIDs,
	}

	token, createErr := client.CreateMonitorServiceTokenForServiceType(context.Background(), "dbaas", createOpts)
	if createErr != nil {
		t.Errorf("Error creating monitor-api token : %s", createErr)
	}

	if token == nil {
		t.Errorf("Error generating token. Did not get token back.")
	}

	mClient, fixtureTeardown := createTestMonitorClient(t, fixturesYaml, token)

	var td func()
	if teardown != nil {
		td = func() {
			clientFixtureTeardown()
			teardown()
			fixtureTeardown()
		}
	} else {
		td = func() {
			clientFixtureTeardown()
			fixtureTeardown()
		}
	}

	return mClient, entityIDs, td, err
}
