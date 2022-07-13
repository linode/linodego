package integration

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
)

func TestPhoneNumber_SendVerificationCode(t *testing.T) {
	client := createMockClient(t)

	requestData := linodego.SendPhoneNumberVerificationCodeOptions{
		ISOCode:     "US",
		PhoneNumber: "137-137-1337",
	}

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "/profile/phone-number"),
		mockRequestBodyValidate(t, requestData, nil))

	if err := client.SendPhoneNumberVerificationCode(context.Background(), requestData); err != nil {
		t.Fatal(err)
	}
}

func TestPhoneNumber_Delete(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, "/profile/phone-number"),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.DeletePhoneNumber(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestPhoneNumber_Verify(t *testing.T) {
	client := createMockClient(t)

	requestData := linodego.VerifyPhoneNumberOptions{
		OTPCode: "123456",
	}

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "/profile/phone-number/verify"),
		mockRequestBodyValidate(t, requestData, nil))

	if err := client.VerifyPhoneNumber(context.Background(), requestData); err != nil {
		t.Fatal(err)
	}
}
