package unit

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestListDatabaseMySQL_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("mysql_databases_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/mysql/instances", fixtureData)
	databases, err := base.Client.ListMySQLDatabases(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, databases, "Expected non-empty mysql database list")
}

func TestDatabaseMySQL_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("mysql_database_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/mysql/instances/123", fixtureData)

	db, err := base.Client.GetMySQLDatabase(context.Background(), 123)
	assert.NoError(t, err)

	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, "mysql", db.Engine)
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
	assert.Equal(t, "8.0.26", db.Version)
}

func TestDatabaseMySQL_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("mysql_database_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.MySQLUpdateOptions{
		Label: "example-db-updated",
	}

	base.MockPut("databases/mysql/instances/123", fixtureData)

	db, err := base.Client.UpdateMySQLDatabase(context.Background(), 123, requestData)
	assert.NoError(t, err)

	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, "mysql", db.Engine)
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
	assert.Equal(t, "8.0.26", db.Version)
}

func TestDatabaseMySQL_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("mysql_database_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.MySQLCreateOptions{
		Label:  "example-db-created",
		Region: "us-east",
		Type:   "g6-dedicated-2",
		Engine: "mysql",
	}

	base.MockPost("databases/mysql/instances", fixtureData)

	db, err := base.Client.CreateMySQLDatabase(context.Background(), requestData)
	assert.NoError(t, err)

	assert.Equal(t, 3, db.ClusterSize)
	assert.Equal(t, "mysql", db.Engine)
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
	assert.Equal(t, "8.0.26", db.Version)
}

func TestDatabaseMySQL_Delete(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, "databases/mysql/instances/123"), httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteMySQLDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseMySQL_SSL_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("mysql_database_ssl_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/mysql/instances/123/ssl", fixtureData)

	ssl, err := base.Client.GetMySQLDatabaseSSL(context.Background(), 123)
	assert.NoError(t, err)

	expectedCACertificate := []byte("-----BEGIN CERTIFICATE-----\nThis is a test certificate\n-----END CERTIFICATE-----\n")

	assert.Equal(t, expectedCACertificate, ssl.CACertificate)
}

func TestDatabaseMySQL_Credentials_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("mysql_database_credentials_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("databases/mysql/instances/123/credentials", fixtureData)

	creds, err := base.Client.GetMySQLDatabaseCredentials(context.Background(), 123)
	assert.NoError(t, err)

	assert.Equal(t, "linroot", creds.Username)
	assert.Equal(t, "s3cur3P@ssw0rd", creds.Password)
}

func TestDatabaseMySQL_Credentials_Reset(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/mysql/instances/123/credentials/reset"), httpmock.NewStringResponder(200, "{}"))

	if err := client.ResetMySQLDatabaseCredentials(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseMySQL_Patch(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/mysql/instances/123/patch"), httpmock.NewStringResponder(200, "{}"))

	if err := client.PatchMySQLDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseMySQL_Suspend(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/mysql/instances/123/suspend"), httpmock.NewStringResponder(200, "{}"))

	if err := client.SuspendMySQLDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseMySQL_Resume(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "databases/mysql/instances/123/resume"), httpmock.NewStringResponder(200, "{}"))

	if err := client.ResumeMySQLDatabase(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}
