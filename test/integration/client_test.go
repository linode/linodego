package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
)

func TestClient_NGINXRetry(t *testing.T) {
	client := createMockClient(t)

	// Recreate the NGINX LB error
	nginxErrorFunc := func(request *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(400, "")

		resp.Header.Add("Server", "nginx")
		resp.Header.Set("Content-Type", "text/html")

		return resp, nil
	}

	step := 0

	httpmock.RegisterRegexpResponder("PUT",
		mockRequestURL(t, "/profile"), func(request *http.Request) (*http.Response, error) {
			if step == 0 {
				step++
				return nginxErrorFunc(request)
			}

			step++
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
