package linodego

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego/internal/testutil"
)

func TestClient_SetAPIVersion(t *testing.T) {
	defaultURL := "https://api.linode.com/v4"

	baseURL := "api.very.cool.com"
	apiVersion := "v4beta"
	expectedHost := fmt.Sprintf("https://%s/%s", baseURL, apiVersion)

	updatedBaseURL := "api.more.cool.com"
	updatedAPIVersion := "v4beta_changed"
	updatedExpectedHost := fmt.Sprintf("https://%s/%s", updatedBaseURL, updatedAPIVersion)

	protocolBaseURL := "http://api.more.cool.com"
	protocolAPIVersion := "v4_http"
	protocolExpectedHost := fmt.Sprintf("%s/%s", protocolBaseURL, protocolAPIVersion)

	client := NewClient(nil)

	if client.resty.BaseURL != defaultURL {
		t.Fatal(cmp.Diff(client.resty.BaseURL, defaultURL))
	}

	client.SetBaseURL(baseURL)
	client.SetAPIVersion(apiVersion)

	if client.resty.BaseURL != expectedHost {
		t.Fatal(cmp.Diff(client.resty.BaseURL, expectedHost))
	}

	// Ensure setting twice does not cause conflicts
	client.SetBaseURL(updatedBaseURL)
	client.SetAPIVersion(updatedAPIVersion)

	if client.resty.BaseURL != updatedExpectedHost {
		t.Fatal(cmp.Diff(client.resty.BaseURL, updatedExpectedHost))
	}

	// Revert
	client.SetBaseURL(baseURL)
	client.SetAPIVersion(apiVersion)

	if client.resty.BaseURL != expectedHost {
		t.Fatal(cmp.Diff(client.resty.BaseURL, expectedHost))
	}

	// Custom protocol
	client.SetBaseURL(protocolBaseURL)
	client.SetAPIVersion(protocolAPIVersion)

	if client.resty.BaseURL != protocolExpectedHost {
		t.Fatal(cmp.Diff(client.resty.BaseURL, expectedHost))
	}
}

func TestClient_NewFromEnv(t *testing.T) {
	file := createTestConfig(t, configNewFromEnv)

	// This is cool
	t.Setenv(APIEnvVar, "")
	t.Setenv(APIConfigEnvVar, file.Name())
	t.Setenv(APIConfigProfileEnvVar, "cool")

	client, err := NewClientFromEnv(nil)
	if err != nil {
		t.Fatal(err)
	}

	if client.selectedProfile != "cool" {
		t.Fatalf("mismatched profile: %s != %s", client.selectedProfile, "cool")
	}

	if client.loadedProfile != "" {
		t.Fatal("expected empty loaded profile")
	}

	if err := client.UseProfile("cool"); err != nil {
		t.Fatal(err)
	}

	if client.loadedProfile != "cool" {
		t.Fatal("expected cool as loaded profile")
	}
}

func TestClient_NewFromEnvToken(t *testing.T) {
	t.Setenv(APIEnvVar, "blah")

	client, err := NewClientFromEnv(nil)
	if err != nil {
		t.Fatal(err)
	}

	if client.resty.Header.Get("Authorization") != "Bearer blah" {
		t.Fatal("token not found in auth header: blah")
	}
}

func TestClient_UseURL(t *testing.T) {
	client := NewClient(nil)

	if _, err := client.UseURL("https://api.test1.linode.com/"); err != nil {
		t.Fatal(err)
	}

	if client.baseURL != "api.test1.linode.com" {
		t.Fatalf("mismatched base url: %s", client.baseURL)
	}

	if client.apiVersion != "v4" {
		t.Fatalf("mismatched api version: %s", client.apiVersion)
	}

	if _, err := client.UseURL("https://api.test2.linode.com/v4beta"); err != nil {
		t.Fatal(err)
	}

	if client.baseURL != "api.test2.linode.com" {
		t.Fatalf("mismatched base url: %s", client.baseURL)
	}

	if client.apiVersion != "v4beta" {
		t.Fatalf("mismatched api version: %s", client.apiVersion)
	}
}

const configNewFromEnv = `
[default]
api_url = api.cool.linode.com
api_version = v4beta

[cool]
token = blah
`

func TestDebugLogSanitization(t *testing.T) {
	type instanceResponse struct {
		ID     int    `json:"id"`
		Region string `json:"region"`
		Label  string `json:"label"`
	}

	var testResp = instanceResponse{
		ID:     100,
		Region: "test-central",
		Label:  "this-is-a-test-linode",
	}
	var lgr bytes.Buffer

	plainTextToken := "NOTANAPIKEY"

	mockClient := testutil.CreateMockClient(t, NewClient)
	logger := testutil.CreateLogger()
	mockClient.SetLogger(logger)
	logger.L.SetOutput(&lgr)

	mockClient.SetDebug(true)
	if !mockClient.resty.Debug {
		t.Fatal("debug should be enabled")
	}
	mockClient.SetHeader("Authorization", fmt.Sprintf("Bearer %s", plainTextToken))

	if mockClient.resty.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", plainTextToken) {
		t.Fatal("token not found in auth header")
	}

	httpmock.RegisterRegexpResponder("GET", testutil.MockRequestURL("/linode/instances"),
		httpmock.NewJsonResponderOrPanic(200, &testResp))

	result, err := doGETRequest[instanceResponse](
		context.Background(),
		mockClient,
		"/linode/instances",
	)
	if err != nil {
		t.Fatal(err)
	}

	logInfo := lgr.String()
	if !strings.Contains(logInfo, "Bearer *******************************") {
		t.Fatal("sanitized bearer token was expected")
	}

	if !reflect.DeepEqual(*result, testResp) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(result, testResponse))
	}
}
