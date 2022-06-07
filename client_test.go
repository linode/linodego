package linodego

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
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

	if client.resty.HostURL != defaultURL {
		t.Fatal(cmp.Diff(client.resty.HostURL, defaultURL))
	}

	client.SetBaseURL(baseURL)
	client.SetAPIVersion(apiVersion)

	if client.resty.HostURL != expectedHost {
		t.Fatal(cmp.Diff(client.resty.HostURL, expectedHost))
	}

	// Ensure setting twice does not cause conflicts
	client.SetBaseURL(updatedBaseURL)
	client.SetAPIVersion(updatedAPIVersion)

	if client.resty.HostURL != updatedExpectedHost {
		t.Fatal(cmp.Diff(client.resty.HostURL, updatedExpectedHost))
	}

	// Revert
	client.SetBaseURL(baseURL)
	client.SetAPIVersion(apiVersion)

	if client.resty.HostURL != expectedHost {
		t.Fatal(cmp.Diff(client.resty.HostURL, expectedHost))
	}

	// Custom protocol
	client.SetBaseURL(protocolBaseURL)
	client.SetAPIVersion(protocolAPIVersion)

	if client.resty.HostURL != protocolExpectedHost {
		t.Fatal(cmp.Diff(client.resty.HostURL, expectedHost))
	}
}

func TestClient_NewFromEnv(t *testing.T) {
	file := createTestConfig(t, configNewFromEnv)

	// This is cool
	t.Setenv(APIConfigEnvVar, file.Name())
	t.Setenv(APIConfigProfileEnvVar, "cool")

	client, err := NewClientFromEnv(nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedURL := "https://api.cool.linode.com/v4beta"

	if client.resty.HostURL != expectedURL {
		t.Fatalf("mismatched host url: %s != %s", client.resty.HostURL, expectedURL)
	}

	if client.resty.Header.Get("Authorization") != "Bearer blah" {
		t.Fatalf("token not found in auth header: %s", "blah")
	}
}

const configNewFromEnv = `
[default]
api_url = api.cool.linode.com
api_version = v4beta

[cool]
token = blah
`
