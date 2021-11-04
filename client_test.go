package linodego

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestClient_SetAPIVersion(t *testing.T) {
	defaultURL := "https://api.linode.com/v4"
	baseURL := "api.very.cool.com"
	apiVersion := "v4beta"
	expectedHost := fmt.Sprintf("https://%s/%s", baseURL, apiVersion)

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
	client.SetBaseURL(baseURL)
	client.SetAPIVersion(apiVersion)

	if client.resty.HostURL != expectedHost {
		t.Fatal(cmp.Diff(client.resty.HostURL, expectedHost))
	}
}
