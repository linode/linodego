package linodego

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

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
	tests := []struct {
		name        string
		inputURL    string
		wantBaseURL string
		wantErr     string
	}{
		{
			name:        "Standard v4 URL",
			inputURL:    "https://api.test1.linode.com/",
			wantBaseURL: "https://api.test1.linode.com/v4",
		},
		{
			name:        "Beta v4 URL",
			inputURL:    "https://api.test2.linode.com/v4beta",
			wantBaseURL: "https://api.test2.linode.com/v4beta",
		},
		{
			name:     "Missing scheme",
			inputURL: "api.test3.linode.com/v4",
			wantErr:  "need both scheme and host in API URL, got \"api.test3.linode.com/v4\"",
		},
		{
			name:     "Missing host",
			inputURL: "https://",
			wantErr:  "need both scheme and host in API URL, got \"https://\"",
		},
		{
			name:     "Invalid URL",
			inputURL: "ht!tp://bad_url",
			wantErr:  "failed to parse URL: parse \"ht!tp://bad_url\": first path segment in URL cannot contain colon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(nil)

			_, err := client.UseURL(tt.inputURL)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
				// skip further checks if error was expected
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if client.resty.BaseURL != tt.wantBaseURL {
				t.Fatalf("mismatched base url: got %s, want %s", client.resty.BaseURL, tt.wantBaseURL)
			}
		})
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

func TestDoRequest_BeforeRequestSuccess(t *testing.T) {
	var capturedRequest *http.Request

	handler := func(w http.ResponseWriter, r *http.Request) {
		capturedRequest = r // Capture the request to inspect it later
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &httpClient{
		httpClient: server.Client(),
	}

	// Define a mutator that successfully modifies the request
	mutator := func(req *http.Request) error {
		req.Header.Set("X-Custom-Header", "CustomValue")
		return nil
	}

	client.httpOnBeforeRequest(mutator)

	err := client.doRequest(context.Background(), http.MethodGet, server.URL, RequestParams{})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Check if the header was successfully added to the captured request
	if reqHeader := capturedRequest.Header.Get("X-Custom-Header"); reqHeader != "CustomValue" {
		t.Fatalf("expected X-Custom-Header to be set to CustomValue, got: %v", reqHeader)
	}
}

func TestDoRequest_BeforeRequestError(t *testing.T) {
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

	client.httpOnBeforeRequest(mutator)

	err := client.doRequest(context.Background(), http.MethodGet, server.URL, RequestParams{})
	expectedErr := "failed to mutate before request"
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("expected error %q, got: %v", expectedErr, err)
	}
}

func TestDoRequest_AfterResponseSuccess(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// Create a custom RoundTripper to capture the response
	tr := &testRoundTripper{
		Transport: server.Client().Transport,
	}
	client := &httpClient{
		httpClient: &http.Client{Transport: tr},
	}

	mutator := func(resp *http.Response) error {
		resp.Header.Set("X-Modified-Header", "ModifiedValue")
		return nil
	}

	client.httpOnAfterResponse(mutator)

	err := client.doRequest(context.Background(), http.MethodGet, server.URL, RequestParams{})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Check if the header was successfully added to the response
	if respHeader := tr.Response.Header.Get("X-Modified-Header"); respHeader != "ModifiedValue" {
		t.Fatalf("expected X-Modified-Header to be set to ModifiedValue, got: %v", respHeader)
	}
}

func TestDoRequest_AfterResponseError(t *testing.T) {
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

	mutator := func(resp *http.Response) error {
		return errors.New("mutator error")
	}

	client.httpOnAfterResponse(mutator)

	err := client.doRequest(context.Background(), http.MethodGet, server.URL, RequestParams{})
	expectedErr := "failed to mutate after response"
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("expected error %q, got: %v", expectedErr, err)
	}
}

func TestDoRequestLogging_Success(t *testing.T) {
	var logBuffer bytes.Buffer
	logger := createLogger()
	logger.l.SetOutput(&logBuffer) // Redirect log output to buffer

	client := &httpClient{
		httpClient: http.DefaultClient,
		debug:      true,
		logger:     logger,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	params := RequestParams{
		Response: &map[string]string{},
	}

	err := client.doRequest(context.Background(), http.MethodGet, server.URL, params)
	if err != nil {
		t.Fatal(cmp.Diff(nil, err))
	}

	logInfo := logBuffer.String()
	logInfoWithoutTimestamps := removeTimestamps(logInfo)

	// Expected logs with templates filled in
	expectedRequestLog := "DEBUG RESTY Sending request:\nMethod: GET\nURL: " + server.URL + "\nHeaders: map[Accept:[application/json] Content-Type:[application/json]]\nBody: "
	expectedResponseLog := "DEBUG RESTY Received response:\nStatus: 200 OK\nHeaders: map[Content-Length:[21] Content-Type:[text/plain; charset=utf-8]]\nBody: {\"message\":\"success\"}"

	if !strings.Contains(logInfo, expectedRequestLog) {
		t.Fatalf("expected log %q not found in logs", expectedRequestLog)
	}
	if !strings.Contains(logInfoWithoutTimestamps, expectedResponseLog) {
		t.Fatalf("expected log %q not found in logs", expectedResponseLog)
	}
}

func TestDoRequestLogging_Error(t *testing.T) {
	var logBuffer bytes.Buffer
	logger := createLogger()
	logger.l.SetOutput(&logBuffer) // Redirect log output to buffer

	client := &httpClient{
		httpClient: http.DefaultClient,
		debug:      true,
		logger:     logger,
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

	logInfo := logBuffer.String()
	expectedLog := "ERROR RESTY failed to encode body"

	if !strings.Contains(logInfo, expectedLog) {
		t.Fatalf("expected log %q not found in logs", expectedLog)
	}
}

func removeTimestamps(log string) string {
	lines := strings.Split(log, "\n")
	var filteredLines []string
	for _, line := range lines {
		// Find the index of the "Date:" substring
		if index := strings.Index(line, "Date:"); index != -1 {
			// Cut off everything after "Date:"
			trimmedLine := strings.TrimSpace(line[:index])
			filteredLines = append(filteredLines, trimmedLine+"]")
		} else {
			filteredLines = append(filteredLines, line)
		}
	}
	return strings.Join(filteredLines, "\n")
}

type testRoundTripper struct {
	Transport http.RoundTripper
	Response  *http.Response
}

func (t *testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Transport.RoundTrip(req)
	if err == nil {
		t.Response = resp
	}
	return resp, err
}

func TestClient_CustomRootCAWithCustomRoundTripper(t *testing.T) {
	caFile, err := os.CreateTemp(t.TempDir(), "linodego_test_ca_*")
	if err != nil {
		t.Fatalf("Failed to create temp ca file: %s", err)
	}
	defer os.Remove(caFile.Name())

	t.Setenv(APIHostCert, caFile.Name())

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// Create a custom RoundTripper
	tr := &testRoundTripper{
		Transport: server.Client().Transport,
	}

	buf := new(strings.Builder)
	log.SetOutput(buf)

	NewClient(&http.Client{Transport: tr})

	expectedLog := "Custom transport is not allowed with a custom root CA"

	if !strings.Contains(buf.String(), expectedLog) {
		t.Fatalf("expected log %q not found in logs", expectedLog)
	}

	log.SetOutput(os.Stderr)
}

func TestClient_CustomRootCAWithoutCustomRoundTripper(t *testing.T) {
	caFile, err := os.CreateTemp(t.TempDir(), "linodego_test_ca_*")
	if err != nil {
		t.Fatalf("Failed to create temp ca file: %s", err)
	}
	defer os.Remove(caFile.Name())

	tests := []struct {
		name       string
		setCA      bool
		httpClient *http.Client
	}{
		{"do not set CA", false, nil},
		{"set CA", true, nil},
		{"do not set CA, use timeout", false, &http.Client{Timeout: time.Second}},
		{"set CA, use timeout", true, &http.Client{Timeout: time.Second}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setCA {
				t.Setenv(APIHostCert, caFile.Name())
			}

			client := NewClient(test.httpClient)
			transport, err := client.resty.Transport()
			if err != nil {
				t.Fatal(err)
			}
			if test.setCA && (transport.TLSClientConfig == nil || transport.TLSClientConfig.RootCAs == nil) {
				t.Error("expected root CAs to be set")
			}
			if !test.setCA && transport.TLSClientConfig != nil {
				t.Errorf("didn't set a custom CA, but client TLS config is not nil: %#v", transport.TLSClientConfig)
			}
		})
	}
}

func TestMonitorClient_SetAPIBasics(t *testing.T) {
	defaultURL := "https://monitor-api.linode.com/v2beta"

	baseURL := "api.very.cool.com"
	apiVersion := "v4beta"
	expectedHost := fmt.Sprintf("https://%s/%s", baseURL, apiVersion)

	updatedBaseURL := "api.more.cool.com"
	updatedAPIVersion := "v4beta_changed"
	updatedExpectedHost := fmt.Sprintf("https://%s/%s", updatedBaseURL, updatedAPIVersion)

	protocolBaseURL := "http://api.more.cool.com"
	protocolAPIVersion := "v4_http"
	protocolExpectedHost := fmt.Sprintf("%s/%s", protocolBaseURL, protocolAPIVersion)

	client := NewMonitorClient(nil)

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
