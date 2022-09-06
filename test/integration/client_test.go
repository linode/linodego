package integration

import (
	"context"
	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"net/http"
	"testing"
)

func TestClient_Aliases(t *testing.T) {
	client, _ := createTestClient(t, "")

	if client.Images == nil {
		t.Error("Expected alias for Images to return a *Resource")
	}
	if client.Instances == nil {
		t.Error("Expected alias for Instances to return a *Resource")
	}
	if client.InstanceSnapshots == nil {
		t.Error("Expected alias for Backups to return a *Resource")
	}
	if client.StackScripts == nil {
		t.Error("Expected alias for StackScripts to return a *Resource")
	}
	if client.Regions == nil {
		t.Error("Expected alias for Regions to return a *Resource")
	}
	if client.Volumes == nil {
		t.Error("Expected alias for Volumes to return a *Resource")
	}
}

func TestClient_NGINXRetry(t *testing.T) {
	client := createMockClient(t)

	// Recreate the NGINX LB error
	nginxErrorFunc := func(request *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(400, nil)
		if err != nil {
			return nil, err
		}

		resp.Header.Add("Server", "nginx")

		return resp, nil
	}

	step := 0

	httpmock.RegisterRegexpResponder("PUT",
		mockRequestURL(t, "/profile"), func(request *http.Request) (*http.Response, error) {
			if step == 0 {
				step = 1
				return nginxErrorFunc(request)
			}

			step = 2
			return httpmock.NewJsonResponse(200, nil)
		})

	if _, err := client.UpdateProfile(context.Background(),
		linodego.ProfileUpdateOptions{}); err != nil {
		t.Fatal(err)
	}

	if step != 2 {
		t.Fatalf("retry checks did not finish")
	}
}
