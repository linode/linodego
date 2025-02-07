package unit

import (
	"context"
	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListDatabasePostgreSQL_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("postgresql_databases_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/postgresql/instances", fixtureData)
	databases, err := base.Client.ListPostgresDatabases(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, databases, "Expected non-empty postgresql database list")
}

func TestDatabasePostgreSQL_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("postgresql_database_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/postgresql/instances/123", fixtureData)

	db, err := base.Client.GetPostgresDatabase(context.Background(), 123)
	assert.NoError(t, err)

	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, "postgresql", db.Engine)
	assert.Equal(t, 123, db.ID)
	assert.Equal(t, "example-db", db.Label)
	assert.Equal(t, linodego.DatabasePlatform("rdbms-default"), db.Platform)
	assert.Equal(t, 3306, db.Port)
	assert.Equal(t, "us-east", db.Region)
	assert.Equal(t, linodego.DatabaseStatus("active"), db.Status)
	assert.Equal(t, 15, db.TotalDiskSizeGB)
	assert.Equal(t, "g6-dedicated-2", db.Type)
	assert.Equal(t, 3, db.Updates.Duration)
	assert.Equal(t, linodego.DatabaseDayOfWeek(1), db.Updates.DayOfWeek)
	assert.Equal(t, linodego.DatabaseMaintenanceFrequency("weekly"), db.Updates.Frequency)
	assert.Equal(t, 0, db.Updates.HourOfDay)
	assert.Equal(t, 2, db.UsedDiskSizeGB)
	assert.Equal(t, "13.2", db.Version)
}

func TestDatabasePostgreSQL_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("postgresql_database_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.PostgresUpdateOptions{
		Label: "example-db-updated",
	}

	base.MockPut("databases/postgresql/instances/123", fixtureData)

	db, err := base.Client.UpdatePostgresDatabase(context.Background(), 123, requestData)
	assert.NoError(t, err)

	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, "postgresql", db.Engine)
	assert.Equal(t, 123, db.ID)
	assert.Equal(t, "example-db-updated", db.Label)
	assert.Equal(t, linodego.DatabasePlatform("rdbms-default"), db.Platform)
	assert.Equal(t, 3306, db.Port)
	assert.Equal(t, "us-east", db.Region)
	assert.Equal(t, linodego.DatabaseStatus("active"), db.Status)
	assert.Equal(t, 15, db.TotalDiskSizeGB)
	assert.Equal(t, "g6-dedicated-2", db.Type)
	assert.Equal(t, 3, db.Updates.Duration)
	assert.Equal(t, linodego.DatabaseDayOfWeek(1), db.Updates.DayOfWeek)
	assert.Equal(t, linodego.DatabaseMaintenanceFrequency("weekly"), db.Updates.Frequency)
	assert.Equal(t, 0, db.Updates.HourOfDay)
	assert.Equal(t, 2, db.UsedDiskSizeGB)
	assert.Equal(t, "13.2", db.Version)
}

func TestDatabasePostgreSQL_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("postgresql_database_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.PostgresCreateOptions{
		Label:  "example-db-created",
		Region: "us-east",
		Type:   "g6-dedicated-2",
		Engine: "postgresql",
	}

	base.MockPost("databases/postgresql/instances", fixtureData)

	db, err := base.Client.CreatePostgresDatabase(context.Background(), requestData)
	assert.NoError(t, err)

	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, "postgresql", db.Engine)
	assert.Equal(t, 123, db.ID)
	assert.Equal(t, "example-db-created", db.Label)
	assert.Equal(t, linodego.DatabasePlatform("rdbms-default"), db.Platform)
	assert.Equal(t, 3306, db.Port)
	assert.Equal(t, "us-east", db.Region)
	assert.Equal(t, linodego.DatabaseStatus("active"), db.Status)
	assert.Equal(t, 15, db.TotalDiskSizeGB)
	assert.Equal(t, "g6-dedicated-2", db.Type)
	assert.Equal(t, 3, db.Updates.Duration)
	assert.Equal(t, linodego.DatabaseDayOfWeek(1), db.Updates.DayOfWeek)
	assert.Equal(t, linodego.DatabaseMaintenanceFrequency("weekly"), db.Updates.Frequency)
	assert.Equal(t, 0, db.Updates.HourOfDay)
	assert.Equal(t, 2, db.UsedDiskSizeGB)
	assert.Equal(t, "13.2", db.Version)
}

func TestDatabasePostgreSQL_Delete(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, "databases/postgresql/instances/123"), httpmock.NewStringResponder(200, "{}"))

	if err := client.DeletePostgresDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabasePostgreSQL_SSL_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("postgresql_database_ssl_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/postgresql/instances/123/ssl", fixtureData)

	ssl, err := base.Client.GetPostgresDatabaseSSL(context.Background(), 123)
	assert.NoError(t, err)

	expectedCACertificate := []byte("-----BEGIN CERTIFICATE-----\nThis is a test certificate\n-----END CERTIFICATE-----\n")

	assert.Equal(t, expectedCACertificate, ssl.CACertificate)
}

func TestDatabasePostgreSQL_Credentials_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("postgresql_database_credentials_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/postgresql/instances/123/credentials", fixtureData)

	creds, err := base.Client.GetPostgresDatabaseCredentials(context.Background(), 123)
	assert.NoError(t, err)

	assert.Equal(t, "linroot", creds.Username)
	assert.Equal(t, "s3cur3P@ssw0rd", creds.Password)
}

func TestDatabasePostgreSQL_Credentials_Reset(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/postgresql/instances/123/credentials/reset"), httpmock.NewStringResponder(200, "{}"))

	if err := client.ResetPostgresDatabaseCredentials(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabasePostgreSQL_Patch(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/postgresql/instances/123/patch"), httpmock.NewStringResponder(200, "{}"))

	if err := client.PatchPostgresDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}
