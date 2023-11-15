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
				IO: [][]float64{},
				Swap: [][]float64{},
			},
			NetV4: linodego.StatsNet{
				In: [][]float64{},
				Out: [][]float64{},
				PrivateIn: [][]float64{},
				PrivateOut: [][]float64{},
			},
			NetV6: linodego.StatsNet{
				In: [][]float64{},
				Out: [][]float64{},
				PrivateIn: [][]float64{},
				PrivateOut: [][]float64{},
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
