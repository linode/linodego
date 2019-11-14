package linodego_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	. "github.com/linode/linodego"
	"golang.org/x/oauth2"
)

var testingMode = recorder.ModeDisabled
var debugAPI = false
var validTestAPIKey = "NOTANAPIKEY"
var testingPollDuration = time.Duration(15000)

func init() {
	if apiToken, ok := os.LookupEnv("LINODE_TOKEN"); ok {
		validTestAPIKey = apiToken
	}

	if apiDebug, ok := os.LookupEnv("LINODE_DEBUG"); ok {
		if parsed, err := strconv.ParseBool(apiDebug); err == nil {
			debugAPI = parsed
			log.Println("[INFO] LINODE_DEBUG being set to", debugAPI)
		} else {
			log.Println("[WARN] LINODE_DEBUG should be an integer, 0 or 1")
		}
	}

	if envFixtureMode, ok := os.LookupEnv("LINODE_FIXTURE_MODE"); ok {
		if envFixtureMode == "record" {
			log.Printf("[INFO] LINODE_FIXTURE_MODE %s will be used for tests", envFixtureMode)
			testingMode = recorder.ModeRecording
		} else if envFixtureMode == "play" {
			log.Printf("[INFO] LINODE_FIXTURE_MODE %s will be used for tests", envFixtureMode)
			testingMode = recorder.ModeReplaying
			testingPollDuration = 1
		}
	}
}

// testRecorder returns a go-vcr recorder and an associated function that the caller must defer
func testRecorder(t *testing.T, fixturesYaml string, testingMode recorder.Mode) (r *recorder.Recorder, recordStopper func()) {
	if t != nil {
		t.Helper()
	}

	r, err := recorder.NewAsMode(fixturesYaml, testingMode, nil)
	if err != nil {
		log.Fatalln(err)
	}

	recordStopper = func() {
		r.Stop()
	}
	return
}

// createTestClient is a testing helper to creates a linodego.Client initialized using
// environment variables and configured to record or playback testing fixtures.
// The returned function should be deferred by the caller to ensure the fixture
// recording is properly closed.
func createTestClient(t *testing.T, fixturesYaml string) (*Client, func()) {
	var (
		c      Client
		apiKey *string
	)
	if t != nil {
		t.Helper()
	}

	apiKey = &validTestAPIKey

	var recordStopper func()
	var r http.RoundTripper

	if testing.Short() {
		apiKey = nil
	}

	if len(fixturesYaml) > 0 {
		r, recordStopper = testRecorder(t, fixturesYaml, testingMode)
	} else {
		r = nil
		recordStopper = func() {}
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *apiKey})
	oc := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
			Base:   r,
		},
	}

	c = NewClient(oc)
	c.SetDebug(debugAPI)
	c.SetPollDelay(testingPollDuration)

	return &c, recordStopper
}

func TestClientAliases(t *testing.T) {
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

func TestClient_APIResponseBadGateway(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestClient_APIResponseBadGateway")
	defer teardown()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected Client to handle 502 from API Server")
		}
	}()

	_, err := client.ListImages(context.Background(), nil)

	if err == nil {
		t.Errorf("Error should be thrown on 502 Response from API")
	}

	responseError, ok := err.(*Error)

	if !ok {
		t.Errorf("Error type did not match the expected result")
	}

	if !strings.Contains(responseError.Message, "Unexpected Content-Type") {
		t.Errorf("Error message does not contain: \"Unexpected Content-Type\"")
	}
}
