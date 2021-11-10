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

	updatedBaseURL := "api.more.cool.com"
	updatedAPIVersion := "v4beta_changed"
	updatedExpectedHost := fmt.Sprintf("https://%s/%s", updatedBaseURL, updatedAPIVersion)

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
}
