package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeBalancerStats_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_stats_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nodebalancers/123/stats", fixtureData)

	stats, err := base.Client.GetNodeBalancerStats(context.Background(), 123)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, stats)

	assert.Equal(t, "NodeBalancer Stats", stats.Title)

	assert.Len(t, stats.Data.Connections, 2)
	assert.Len(t, stats.Data.Connections[0], 2)

	assert.Len(t, stats.Data.Traffic.In, 2)
	assert.Len(t, stats.Data.Traffic.Out, 2)

	assert.Equal(t, 1000.0, stats.Data.Traffic.In[0][0])
	assert.Equal(t, 500.0, stats.Data.Traffic.Out[0][0])
}
