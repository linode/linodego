package unit

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceStats_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_stats_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/36183732/stats", fixtureData)

	stats, err := base.Client.GetInstanceStats(context.Background(), 36183732)
	assert.NoError(t, err)

	assert.Equal(t, "linode.com - my-linode (linode123456) - day (5 min avg)", stats.Title)

	fmt.Printf("Stats: %+v\n", stats) //TODO:: Debug for assertion remove later

	assert.Len(t, stats.Data.CPU, 1)

	assert.Len(t, stats.Data.IO.IO, 1)
	assert.Len(t, stats.Data.IO.Swap, 1)

	assert.Len(t, stats.Data.NetV4.In, 1)
	assert.Len(t, stats.Data.NetV4.Out, 1)
	assert.Len(t, stats.Data.NetV4.PrivateIn, 1)
	assert.Len(t, stats.Data.NetV4.PrivateOut, 1)

	assert.Len(t, stats.Data.NetV6.In, 1)
	assert.Len(t, stats.Data.NetV6.Out, 1)
	assert.Len(t, stats.Data.NetV6.PrivateIn, 1)
	assert.Len(t, stats.Data.NetV6.PrivateOut, 1)

	assert.Equal(t, 0.42, stats.Data.CPU[0][1])
	assert.Equal(t, 0.19, stats.Data.IO.IO[0][1])
	assert.Equal(t, 0.0, stats.Data.IO.Swap[0][1])
	assert.Equal(t, 2004.36, stats.Data.NetV4.In[0][1])
	assert.Equal(t, 3928.91, stats.Data.NetV4.Out[0][1])
	assert.Equal(t, 195.18, stats.Data.NetV6.PrivateIn[0][1])
	assert.Equal(t, 5.6, stats.Data.NetV6.PrivateOut[0][1])
}
