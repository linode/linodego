package integration

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
)

func TestTwoFactor_CreateSecret(t *testing.T) {
	client := createMockClient(t)

	expectedResponse := linodego.TwoFactorSecret{
		Secret: "verysecureverysecureverysecure",
	}

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "/profile/tfa-enable"),
		httpmock.NewJsonResponderOrPanic(200, expectedResponse))

	secret, err := client.CreateTwoFactorSecret(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*secret, expectedResponse) {
		t.Fatalf("returned value did not match expected response: %s", cmp.Diff(*secret, expectedResponse))
	}
}

func TestTwoFactor_Disable(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "/profile/tfa-disable"),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.DisableTwoFactor(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestTwoFactor_Confirm(t *testing.T) {
	client := createMockClient(t)

	request := linodego.ConfirmTwoFactorOptions{TFACode: "reallycoolandlegittfacode"}
	response := linodego.ConfirmTwoFactorResponse{Scratch: "really cool scratch code"}

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "/profile/tfa-enable-confirm"),
		mockRequestBodyValidate(t, request, response))

	runResult, err := client.ConfirmTwoFactor(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*runResult, response) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(*runResult, response))
	}
}
