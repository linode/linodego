package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestEntities_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("entities_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet(formatMockAPIPath("entities"), fixtureData)

	entities, err := base.Client.ListEntities(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.Equal(t, 7, entities[0].ID)
}
