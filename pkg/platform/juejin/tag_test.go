package juejin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTags(t *testing.T) {
	setupClient(t)
	tags, count, err := client.ListTags("docker", StartCursor)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(tags))
	assert.NotEqual(t, 0, count)
}

func TestListAllCategories(t *testing.T) {
	setupClient(t)
	categories, err := client.ListCategories()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(categories))
}
