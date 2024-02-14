package integration

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
)

func TestInstanceStats_Get(t *testing.T) {
	client := createMockClient(t)

	desiredResponse := linodego.InstanceStats{
		Title: "test_title",
		Data: linodego.InstanceStatsData{
			CPU: [][]float64{},
			IO: linodego.StatsIO{
				IO:   [][]float64{{1.2, 2.3}, {3.4, 4.5}},
				Swap: [][]float64{{14, 2.3}, {34, 4.5}},
			},
			NetV4: linodego.StatsNet{
				In:         [][]float64{{1.2, 2.3}, {3.4, 4.5}},
				Out:        [][]float64{{1, 2}, {3, 4}},
				PrivateIn:  [][]float64{{2, 3}, {4, 5}},
				PrivateOut: [][]float64{{12.1, 2.33}, {4.4, 4.5}},
			},
			NetV6: linodego.StatsNet{
				In:         [][]float64{{1.2, .3}, {3.4, .5}},
				Out:        [][]float64{{0, 2.3}, {3, 4.55}},
				PrivateIn:  [][]float64{{1.24, 3}, {3, 5}},
				PrivateOut: [][]float64{{1, 6}, {7, 8}},
			},
		},
	}

	httpmock.RegisterRegexpResponder("GET", mockRequestURL(t, "/instances/36183732/stats"),
		httpmock.NewJsonResponderOrPanic(200, &desiredResponse))

	questions, err := client.GetInstanceStats(context.Background(), 36183732)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*questions, desiredResponse) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(questions, desiredResponse))
	}
}
