package unit

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
)

func TestLKECluster_Regenerate(t *testing.T) {
	client := createMockClient(t)

	requestData := linodego.LKEClusterRegenerateOptions{
		KubeConfig:   true,
		ServiceToken: false,
	}

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "clusters/1234/regenerate"),
		mockRequestBodyValidate(t, requestData, nil))

	if _, err := client.RegenerateLKECluster(context.Background(), 1234, requestData); err != nil {
		t.Fatal(err)
	}
}

func TestLKECluster_DeleteServiceToken(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, "clusters/1234/servicetoken"), httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteLKEClusterServiceToken(context.Background(), 1234); err != nil {
		t.Fatal(err)
	}
}
