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
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego/v2/internal/testutil"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, hc *http.Client) Client {
	t.Helper()

	client, err := NewClient(hc)
	require.NoError(t, err)
	return client
}

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

	client := newTestClient(t, nil)

	if client.hostURL != defaultURL {
		t.Fatal(cmp.Diff(client.hostURL, defaultURL))
	}

	client.SetBaseURL(baseURL)
	client.SetAPIVersion(apiVersion)

	if client.hostURL != expectedHost {
		t.Fatal(cmp.Diff(client.hostURL, expectedHost))
	}

	// Ensure setting twice does not cause conflicts
	client.SetBaseURL(updatedBaseURL)
	client.SetAPIVersion(updatedAPIVersion)

	if client.hostURL != updatedExpectedHost {
		t.Fatal(cmp.Diff(client.hostURL, updatedExpectedHost))
	}

	// Revert
	client.SetBaseURL(baseURL)
	client.SetAPIVersion(apiVersion)

	if client.hostURL != expectedHost {
		t.Fatal(cmp.Diff(client.hostURL, expectedHost))
	}

	// Custom protocol
	client.SetBaseURL(protocolBaseURL)
	client.SetAPIVersion(protocolAPIVersion)

	if client.hostURL != protocolExpectedHost {
		t.Fatal(cmp.Diff(client.hostURL, expectedHost))
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

	if client.header.Get("Authorization") != "Bearer blah" {
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
			client := newTestClient(t, nil)

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

			if client.hostURL != tt.wantBaseURL {
				t.Fatalf("mismatched base url: got %s, want %s", client.hostURL, tt.wantBaseURL)
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

	mockClient := testutil.CreateMockClientWithError(t, NewClient)
	logger := testutil.CreateLogger()
	mockClient.SetLogger(logger)
	logger.L.SetOutput(&lgr)

	mockClient.SetDebug(true)
	if !mockClient.debug {
		t.Fatal("debug should be enabled")
	}
	mockClient.SetHeader("Authorization", fmt.Sprintf("Bearer %s", plainTextToken))

	if mockClient.header.Get("Authorization") != fmt.Sprintf("Bearer %s", plainTextToken) {
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
		if r.URL.Path == "/v4/foo/bar" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"message":"success"}`))
		} else {
			http.NotFound(w, r)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := newTestClient(t, server.Client())
	client.SetBaseURL(server.URL)

	params := requestParams{
		Response: &map[string]string{},
	}

	err := client.doRequest(context.Background(), http.MethodGet, "/foo/bar", params, nil) // Pass only the endpoint
	if err != nil {
		t.Fatal(cmp.Diff(nil, err))
	}

	expected := "success"
	actual := (*params.Response.(*map[string]string))["message"]
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("response mismatch (-expected +actual):\n%s", diff)
	}
}

func TestDoRequest_FailedCreateRequest(t *testing.T) {
	client := newTestClient(t, nil)

	// Create a request with an invalid method to simulate a request creation failure
	err := client.doRequest(context.Background(), "bad method", "/foo/bar", requestParams{}, nil)
	expectedErr := "failed to create request"
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("expected error %q, got: %v", expectedErr, err)
	}
}

func TestDoRequest_Non2xxStatusCode(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "error", http.StatusInternalServerError) // Simulate a 500 Internal Server Error
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := newTestClient(t, server.Client())
	client.SetBaseURL(server.URL)

	err := client.doRequest(context.Background(), http.MethodGet, "/foo/bar", requestParams{}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	httpError, ok := err.(*Error)
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
		_, _ = w.Write([]byte(`invalid json`)) // Simulate invalid JSON
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := newTestClient(t, server.Client())
	client.SetBaseURL(server.URL)

	params := requestParams{
		Response: &map[string]string{},
	}

	err := client.doRequest(context.Background(), http.MethodGet, "/foo/bar", params, nil)
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

	client := newTestClient(t, server.Client())
	client.SetBaseURL(server.URL)

	mutator := func(req *http.Request) error {
		req.Header.Set("X-Custom-Header", "CustomValue")
		return nil
	}

	client.OnBeforeRequest(mutator)

	err := client.doRequest(context.Background(), http.MethodGet, "/foo/bar", requestParams{}, nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

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

	client := newTestClient(t, server.Client())
	client.SetBaseURL(server.URL)

	mutator := func(req *http.Request) error {
		return errors.New("mutator error")
	}

	client.OnBeforeRequest(mutator)

	err := client.doRequest(context.Background(), http.MethodGet, "/foo/bar", requestParams{}, nil)
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
	client := newTestClient(t, &http.Client{Transport: tr})
	client.SetBaseURL(server.URL)

	mutator := func(resp *http.Response) error {
		resp.Header.Set("X-Modified-Header", "ModifiedValue")
		return nil
	}

	client.OnAfterResponse(mutator)

	err := client.doRequest(context.Background(), http.MethodGet, "/foo/bar", requestParams{}, nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

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

	client := newTestClient(t, server.Client())
	client.SetBaseURL(server.URL)

	mutator := func(resp *http.Response) error {
		return errors.New("mutator error")
	}

	client.OnAfterResponse(mutator)

	err := client.doRequest(context.Background(), http.MethodGet, "/foo/bar", requestParams{}, nil)
	expectedErr := "failed to mutate after response"
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("expected error %q, got: %v", expectedErr, err)
	}
}

func TestDoRequestLogging_Success(t *testing.T) {
	var logBuffer bytes.Buffer
	logger := createLogger()
	logger.l.SetOutput(&logBuffer) // Redirect log output to buffer

	client := newTestClient(t, nil)
	client.SetDebug(true)
	client.SetLogger(logger)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	client.SetBaseURL(server.URL)
	defer server.Close()

	params := requestParams{
		Response: &map[string]string{},
	}
	endpoint := "/foo/bar"

	err := client.doRequest(context.Background(), http.MethodGet, endpoint, params, nil)
	if err != nil {
		t.Fatal(cmp.Diff(nil, err))
	}

	logInfo := logBuffer.String()

	expectedRequestParts := []string{
		"GET /v4/foo/bar HTTP/1.1",
		"Accept: application/json",
		"Authorization: Bearer *******************************",
		"Content-Type: application/json",
		"User-Agent: linodego/dev https://github.com/linode/linodego",
	}

	expectedResponseParts := []string{
		"STATUS: 200 OK",
		"PROTO: HTTP/1.1",
		"Content-Length: 21",
		"Content-Type: application/json",
		`"message": "success"`,
	}

	for _, part := range expectedRequestParts {
		if !strings.Contains(logInfo, part) {
			t.Fatalf("expected request part %q not found in logs", part)
		}
	}
	for _, part := range expectedResponseParts {
		if !strings.Contains(logInfo, part) {
			t.Fatalf("expected response part %q not found in logs", part)
		}
	}
}

func TestDoRequestLogging_Error(t *testing.T) {
	var logBuffer bytes.Buffer
	logger := createLogger()
	logger.l.SetOutput(&logBuffer) // Redirect log output to buffer

	client := newTestClient(t, nil)
	client.SetDebug(true)
	client.SetLogger(logger)

	// Create a request with an invalid method to simulate a request creation failure
	err := client.doRequest(context.Background(), "bad method", "/foo/bar", requestParams{}, nil)
	expectedErr := "failed to create request"
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("expected error %q, got: %v", expectedErr, err)
	}

	logInfo := logBuffer.String()
	expectedLog := "ERROR failed to create request"

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

	var logBuf bytes.Buffer
	prevWriter := log.Writer()
	log.SetOutput(&logBuf)
	defer log.SetOutput(prevWriter)
	_, err = NewClient(&http.Client{Transport: tr})
	require.NoError(t, err)
	require.Contains(t, logBuf.String(), "[WARN] Custom root certificate is not supported with a custom transport")
}

func TestClient_CustomRootCAWithMissingFile(t *testing.T) {
	t.Setenv(APIHostCert, "/does/not/exist.pem")

	_, err := NewClient(nil)
	require.ErrorContains(t, err, "failed to read root certificate")
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

			client := newTestClient(t, test.httpClient)
			transport, ok := client.httpClient.Transport.(*http.Transport)
			if !ok {
				t.Fatal("expected *http.Transport")
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

	if client.hostURL != defaultURL {
		t.Fatal(cmp.Diff(client.hostURL, defaultURL))
	}

	client.SetBaseURL(baseURL)
	client.SetAPIVersion(apiVersion)

	if client.hostURL != expectedHost {
		t.Fatal(cmp.Diff(client.hostURL, expectedHost))
	}

	// Ensure setting twice does not cause conflicts
	client.SetBaseURL(updatedBaseURL)
	client.SetAPIVersion(updatedAPIVersion)

	if client.hostURL != updatedExpectedHost {
		t.Fatal(cmp.Diff(client.hostURL, updatedExpectedHost))
	}

	// Revert
	client.SetBaseURL(baseURL)
	client.SetAPIVersion(apiVersion)

	if client.hostURL != expectedHost {
		t.Fatal(cmp.Diff(client.hostURL, expectedHost))
	}

	// Custom protocol
	client.SetBaseURL(protocolBaseURL)
	client.SetAPIVersion(protocolAPIVersion)

	if client.hostURL != protocolExpectedHost {
		t.Fatal(cmp.Diff(client.hostURL, expectedHost))
	}
}

func TestMonitorClient_SetRootCertificateWithCustomRoundTripper(t *testing.T) {
	caFile, err := os.CreateTemp(t.TempDir(), "linodego_test_ca_*")
	if err != nil {
		t.Fatalf("Failed to create temp ca file: %s", err)
	}
	defer os.Remove(caFile.Name())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tr := &testRoundTripper{Transport: server.Client().Transport}
	client := NewMonitorClient(&http.Client{Transport: tr})

	err = client.SetRootCertificate(caFile.Name())
	require.ErrorContains(t, err, "current transport is not an *http.Transport instance")
}

func TestMonitorClient_SetRootCertificateWithMissingFile(t *testing.T) {
	client := NewMonitorClient(nil)

	err := client.SetRootCertificate("/does/not/exist.pem")
	require.ErrorContains(t, err, "failed to read root certificate")
}

func TestMonitorClient_SetRootCertificateWithoutCustomRoundTripper(t *testing.T) {
	caFile, err := os.CreateTemp(t.TempDir(), "linodego_test_ca_*")
	if err != nil {
		t.Fatalf("Failed to create temp ca file: %s", err)
	}
	defer os.Remove(caFile.Name())

	tests := []struct {
		name       string
		httpClient *http.Client
	}{
		{"default http client", nil},
		{"timeout http client", &http.Client{Timeout: time.Second}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewMonitorClient(test.httpClient)

			err := client.SetRootCertificate(caFile.Name())
			require.NoError(t, err)

			transport, ok := client.httpClient.Transport.(*http.Transport)
			if !ok {
				t.Fatal("expected *http.Transport")
			}
			if transport.TLSClientConfig == nil || transport.TLSClientConfig.RootCAs == nil {
				t.Error("expected root CAs to be set")
			}
		})
	}
}

func TestRedactHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header
		wantVal map[string]string
	}{
		{
			name: "redacts authorization header",
			headers: http.Header{
				"Authorization": []string{"Bearer supersecrettoken"},
				"Content-Type":  []string{"application/json"},
			},
			wantVal: map[string]string{
				"Authorization": redactHeadersMap["Authorization"],
				"Content-Type":  "application/json",
			},
		},
		{
			name: "leaves non-sensitive headers unchanged",
			headers: http.Header{
				"Content-Type": []string{"application/json"},
				"Accept":       []string{"application/json"},
			},
			wantVal: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
		},
		{
			name:    "handles empty headers",
			headers: http.Header{},
			wantVal: map[string]string{},
		},
		{
			name: "does not mutate original headers",
			headers: http.Header{
				"Authorization": []string{"Bearer supersecrettoken"},
			},
			wantVal: map[string]string{
				"Authorization": redactHeadersMap["Authorization"],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalAuth := tt.headers.Get("Authorization")

			result := redactHeaders(tt.headers)

			// Verify expected values in result
			for key, expectedVal := range tt.wantVal {
				if got := result.Get(key); got != expectedVal {
					t.Errorf("redactHeaders() header %q = %q, want %q", key, got, expectedVal)
				}
			}

			// Verify original was not mutated
			if tt.headers.Get("Authorization") != originalAuth {
				t.Error("redactHeaders() mutated the original headers")
			}
		})
	}
}

func TestEnableLogSanitization(t *testing.T) {
	mockClient := testutil.CreateMockClientWithError(t, NewClient)
	mockClient.SetDebug(true)

	plainTextToken := "supersecrettoken"
	mockClient.SetToken(plainTextToken)

	var logBuf bytes.Buffer
	logger := testutil.CreateLogger()
	logger.L.SetOutput(&logBuf)
	mockClient.SetLogger(logger)

	httpmock.RegisterResponder("GET", "=~.*",
		httpmock.NewStringResponder(200, `{}`).HeaderSet(http.Header{
			"Authorization": []string{"Bearer " + plainTextToken},
		}))

	err := mockClient.doRequest(context.Background(), http.MethodGet, "/test", requestParams{}, nil)
	require.NoError(t, err)

	logOutput := logBuf.String()

	// Verify token is not present in either request or response logs
	if strings.Contains(logOutput, plainTextToken) {
		t.Errorf("log output contains raw token %q, expected it to be redacted", plainTextToken)
	}

	// Verify Authorization header still appears (as redacted value) in request log
	if !strings.Contains(logOutput, "Authorization") {
		t.Error("expected Authorization header to appear in request log output")
	}
}

func TestDoRequest_RetryCountZero_StillExecutes(t *testing.T) {
	var called atomic.Bool

	handler := func(w http.ResponseWriter, r *http.Request) {
		called.Store(true)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1}`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := newTestClient(t, nil)
	client.SetBaseURL(server.URL)
	client.SetRetryCount(0)

	type result struct {
		ID int `json:"id"`
	}

	var got result

	err := client.doRequest(context.Background(), http.MethodGet, "/test", requestParams{
		Response: &got,
	}, nil)
	require.NoError(t, err, "doRequest should not return an error")
	require.True(t, called.Load(), "server handler should have been called even with retryCount=0")
	require.Equal(t, 1, got.ID, "response should have been decoded")
}
