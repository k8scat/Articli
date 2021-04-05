package juejin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTags(t *testing.T) {
	setup(t)
	count, tags, err := c.ListTags("docker", 0)
	assert.Nil(t, err)
	fmt.Printf("count: %d\n", count)
	fmt.Printf("get: %d\n", len(tags))
	for _, tag := range tags {
		fmt.Printf("id: %s, name: %s\n", tag.ID, tag.Name)
	}
}

func TestListAllTags(t *testing.T) {
	setup(t)
	tags, err := c.ListAllTags("")
	assert.Nil(t, err)
	fmt.Printf("get: %d\n", len(tags))
	for _, tag := range tags {
		fmt.Printf("id: %s, name: %s\n", tag.ID, tag.Name)
	}
}

func TestListAllCategories(t *testing.T) {
	setup(t)
	categories, err := c.ListAllCategories()
	assert.Nil(t, err)
	for _, c := range categories {
		fmt.Printf("id: %s, name: %s\n", c.ID, c.Name)
	}
}
