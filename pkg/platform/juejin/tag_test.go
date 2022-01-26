package juejin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTags(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_JUEJIN_COOKIE"))
	if err != nil {
		t.Fail()
		return
	}

	tags, count, err := client.ListTags("docker", StartCursor)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(tags))
	assert.NotEqual(t, 0, count)
}

func TestListAllCategories(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_JUEJIN_COOKIE"))
	if err != nil {
		t.Fail()
		return
	}

	categories, err := client.ListCategories()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(categories))
}
