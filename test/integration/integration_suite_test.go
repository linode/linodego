package integration

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/linode/linodego"
	"golang.org/x/oauth2"
	"k8s.io/client-go/transport"
)

var testingMode = recorder.ModeDisabled
var debugAPI = false
var validTestAPIKey = "NOTANAPIKEY"

var testingPollDuration = time.Duration(15000)
var testingMaxRetryTime = time.Duration(30) * time.Second

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
			testingMaxRetryTime = time.Duration(1) * time.Nanosecond
		}
	}
}

// testRecorder returns a go-vcr recorder and an associated function that the caller must defer
func testRecorder(t *testing.T, fixturesYaml string, testingMode recorder.Mode, realTransport http.RoundTripper) (r *recorder.Recorder, recordStopper func()) {
	if t != nil {
		t.Helper()
	}

	r, err := recorder.NewAsMode(fixturesYaml, testingMode, realTransport)
	if err != nil {
		log.Fatalln(err)
	}

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		return nil
	})

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Response.Headers, "Date")
		delete(i.Response.Headers, "Retry-After")
		delete(i.Response.Headers, "X-Customer-Uuid")
		delete(i.Response.Headers, "X-Ratelimit-Reset")
		delete(i.Response.Headers, "X-Ratelimit-Remaining")
		delete(i.Response.Headers, "X-Spec-Version")
		return nil
	})

	r.AddFilter(func(i *cassette.Interaction) error {
		re := regexp.MustCompile(`"access_key": "[[:alnum:]]*"`)
		i.Response.Body = re.ReplaceAllString(i.Response.Body, `"access_key": "[SANITIZED]"`)
		re = regexp.MustCompile(`"secret_key": "[[:alnum:]]*"`)
		i.Response.Body = re.ReplaceAllString(i.Response.Body, `"secret_key": "[SANITIZED]"`)
		re = regexp.MustCompile("20[0-9]{2}-[01][0-9]-[0-3][0-9]T[0-2][0-9]:[0-9]{2}:[0-9]{2}")
		i.Response.Body = re.ReplaceAllString(i.Response.Body, "2018-01-02T03:04:05")
		re = regexp.MustCompile("192\\.168\\.((1?[0-9][0-9]?|2[0-4][0-9]|25[0-5])\\.)(1?[0-9][0-9]?|2[0-4][0-9]|25[0-5])")
		i.Response.Body = re.ReplaceAllString(i.Response.Body, "192.168.030.040")
		re = regexp.MustCompile("^192\\.168/!s/((1?[0-9][0-9]?|2[0-4][0-9]|25[0-5])\\.){3}(1?[0-9][0-9]?|2[0-4][0-9]|25[0-5])")
		i.Response.Body = re.ReplaceAllString(i.Response.Body, "0.0.0.0")
		re = regexp.MustCompile("nb-[0-9]{1,3}-[0-9]{1,3}-[0-9]{1,3}-[0-9]{1,3}")
		i.Response.Body = re.ReplaceAllString(i.Response.Body, "nb-0-0-0-0")
		return nil
	})

	r.AddSaveFilter(func(i *cassette.Interaction) error {
		re := regexp.MustCompile("AWSAccessKeyId=[[:alnum:]]{20}")
		i.Response.Body = re.ReplaceAllString(i.Response.Body, "AWSAccessKeyID=SANITIZED")
		i.Request.URL = re.ReplaceAllString(i.Request.URL, "AWSAccessKeyID=SANITIZED")
		return nil
	})

	recordStopper = func() {
		r.Stop()
	}
	return
}

// createTestClient is a testing helper to creates a linodego.Client initialized using
// environment variables and configured to record or playback testing fixtures.
// The returned function should be deferred by the caller to ensure the fixture
// recording is properly closed.
func createTestClient(t *testing.T, fixturesYaml string) (*linodego.Client, func()) {
	var (
		c      linodego.Client
		apiKey *string
	)
	if t != nil {
		t.Helper()
	}

	apiKey = &validTestAPIKey

	var recordStopper func()
	var r http.RoundTripper

	if len(fixturesYaml) > 0 {
		r, recordStopper = testRecorder(t, fixturesYaml, testingMode, nil)
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

	c = linodego.NewClient(oc)
	c.SetDebug(debugAPI).
		SetPollDelay(testingPollDuration).
		SetRetryMaxWaitTime(testingMaxRetryTime)

	return &c, recordStopper
}

// transportRecordWrapper returns a tranport.WrapperFunc which provides the test
// recorder as an http.RoundTripper.
func transportRecorderWrapper(t *testing.T, fixtureYaml string) (transport.WrapperFunc, func()) {
	t.Helper()

	rec, teardown := testRecorder(t, fixtureYaml, testingMode, nil)
	return func(r http.RoundTripper) http.RoundTripper {
		rec.SetTransport(r)
		return rec
	}, teardown
}
