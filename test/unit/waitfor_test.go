package unit

import (
	"context"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego/v2"
)

func TestWaitForInstanceStatusUsesContextDeadline(t *testing.T) {
	client := createMockClient(t)
	client.SetPollDelay(time.Millisecond)

	httpmock.RegisterRegexpResponder("GET", mockRequestURL(t, "linode/instances/123"),
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, linodego.Instance{ID: 123, Status: linodego.InstanceProvisioning})
		})

	_, err := client.WaitForInstanceStatus(waitTestContext(t, 5*time.Millisecond), 123, linodego.InstanceRunning)
	if err == nil {
		t.Fatal("expected deadline error")
	}
}

func TestWaitForVolumeStatusSuccess(t *testing.T) {
	client := createMockClient(t)
	client.SetPollDelay(time.Millisecond)

	step := 0
	httpmock.RegisterRegexpResponder("GET", mockRequestURL(t, "volumes/123"),
		func(_ *http.Request) (*http.Response, error) {
			step++

			status := linodego.VolumeCreating
			if step == 2 {
				status = linodego.VolumeActive
			}

			return httpmock.NewJsonResponse(http.StatusOK, linodego.Volume{ID: 123, Status: status})
		})

	volume, err := client.WaitForVolumeStatus(waitTestContext(t, time.Second), 123, linodego.VolumeActive)
	if err != nil {
		t.Fatal(err)
	}

	if volume.Status != linodego.VolumeActive {
		t.Fatalf("expected volume to be active, got %s", volume.Status)
	}

	if step != 2 {
		t.Fatalf("expected 2 polls, got %d", step)
	}
}

func TestWaitForVolumeLinodeIDNil(t *testing.T) {
	client := createMockClient(t)
	client.SetPollDelay(time.Millisecond)

	httpmock.RegisterRegexpResponder("GET", mockRequestURL(t, "volumes/123"),
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, linodego.Volume{ID: 123, LinodeID: nil})
		})

	volume, err := client.WaitForVolumeLinodeID(waitTestContext(t, time.Second), 123, nil)
	if err != nil {
		t.Fatal(err)
	}

	if volume.LinodeID != nil {
		t.Fatalf("expected nil LinodeID, got %v", volume.LinodeID)
	}
}

func TestEventPollerWaitForFinished(t *testing.T) {
	client := createMockClient(t)
	client.SetPollDelay(time.Millisecond)

	listEventsCount := 0
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`/[a-zA-Z0-9]+/account/events(\?.*)?$`),
		func(_ *http.Request) (*http.Response, error) {
			listEventsCount++

			events := []linodego.Event{{
				ID:     111,
				Status: linodego.EventFinished,
				Action: linodego.ActionLinodeBoot,
				Entity: &linodego.EventEntity{ID: 123, Type: linodego.EntityLinode},
			}}

			if listEventsCount > 1 {
				events = append(events, linodego.Event{
					ID:     456,
					Status: linodego.EventStarted,
					Action: linodego.ActionLinodeBoot,
					Entity: &linodego.EventEntity{ID: 123, Type: linodego.EntityLinode},
				})
			}

			return httpmock.NewJsonResponse(http.StatusOK, map[string]any{
				"data":    events,
				"page":    1,
				"pages":   1,
				"results": len(events),
			})
		})

	poller, err := client.NewEventPoller(
		context.Background(),
		123,
		linodego.EntityLinode,
		linodego.ActionLinodeBoot,
	)
	if err != nil {
		t.Fatal(err)
	}

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`/[a-zA-Z0-9]+/account/events/456$`),
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, linodego.Event{ID: 456, Status: linodego.EventFinished})
		})

	event, err := poller.WaitForFinished(waitTestContext(t, time.Second))
	if err != nil {
		t.Fatal(err)
	}

	if event.Status != linodego.EventFinished {
		t.Fatalf("expected event to be finished, got %s", event.Status)
	}
}

func waitTestContext(t *testing.T, timeout time.Duration) context.Context {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(cancel)

	return ctx
}
