package linodego

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego/internal/testutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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

	client := NewClient(server.Client())
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
	client := NewClient(nil)

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

	client := NewClient(server.Client())
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

	client := NewClient(server.Client())
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

	client := NewClient(server.Client())
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

	client := NewClient(server.Client())
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
	client := NewClient(&http.Client{Transport: tr})
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

	client := NewClient(server.Client())
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

	client := NewClient(nil)
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

	err := client.doRequest(context.Background(), http.MethodGet, server.URL, params, nil)
	if err != nil {
		t.Fatal(cmp.Diff(nil, err))
	}

	logInfo := logBuffer.String()

	expectedRequestParts := []string{
		"GET /v4/" + server.URL + " " + "HTTP/1.1",
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

	client := NewClient(nil)
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
