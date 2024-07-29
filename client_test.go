package linodego

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
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

	testResp := instanceResponse{
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

func TestDoRequest_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &httpClient{
		httpClient: server.Client(),
	}

	params := RequestParams{
		Response: &map[string]string{},
	}

	err := client.doRequest(context.Background(), http.MethodGet, server.URL, params)
	if err != nil {
		t.Fatal(cmp.Diff(nil, err))
	}

	expected := "success"
	actual := (*params.Response.(*map[string]string))["message"]
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("response mismatch (-expected +actual):\n%s", diff)
	}
}

func TestDoRequest_FailedEncodeBody(t *testing.T) {
	client := &httpClient{
		httpClient: http.DefaultClient,
	}

	params := RequestParams{
		Body: map[string]interface{}{
			"invalid": func() {},
		},
	}

	err := client.doRequest(context.Background(), http.MethodPost, "http://example.com", params)
	expectedErr := "failed to encode body"
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("expected error %q, got: %v", expectedErr, err)
	}
}

func TestDoRequest_FailedCreateRequest(t *testing.T) {
	client := &httpClient{
		httpClient: http.DefaultClient,
	}

	// Create a request with an invalid URL to simulate a request creation failure
	err := client.doRequest(context.Background(), http.MethodGet, "http://invalid url", RequestParams{})
	expectedErr := "failed to create request"
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("expected error %q, got: %v", expectedErr, err)
	}
}

func TestDoRequest_Non2xxStatusCode(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "error", http.StatusInternalServerError)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &httpClient{
		httpClient: server.Client(),
	}

	err := client.doRequest(context.Background(), http.MethodGet, server.URL, RequestParams{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	httpError, ok := err.(Error)
	if !ok {
		t.Fatalf("expected error to be of type Error, got %T", err)
	}
	if httpError.Code != http.StatusInternalServerError {
		t.Fatalf("expected status code %d, got %d", http.StatusInternalServerError, httpError.Code)
	}
	if !strings.Contains(httpError.Message, "error") {
		t.Fatalf("expected error message to contain %q, got %v", "error", httpError.Message)
	}
}

func TestDoRequest_FailedDecodeResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`invalid json`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &httpClient{
		httpClient: server.Client(),
	}

	params := RequestParams{
		Response: &map[string]string{},
	}

	err := client.doRequest(context.Background(), http.MethodGet, server.URL, params)
	expectedErr := "failed to decode response"
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("expected error %q, got: %v", expectedErr, err)
	}
}

func TestDoRequest_MutatorError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &httpClient{
		httpClient: server.Client(),
	}

	mutator := func(req *http.Request) error {
		return errors.New("mutator error")
	}

	err := client.doRequest(context.Background(), http.MethodGet, server.URL, RequestParams{}, mutator)
	expectedErr := "failed to mutate request"
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("expected error %q, got: %v", expectedErr, err)
	}
}
