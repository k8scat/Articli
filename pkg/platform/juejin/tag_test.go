package juejin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTags(t *testing.T) {
	setupClient(t)
	count, tags, err := client.ListTags("docker", 0)
	assert.Nil(t, err)
	fmt.Printf("count: %d\n", count)
	fmt.Printf("get: %d\n", len(tags))
	for _, tag := range tags {
		fmt.Printf("id: %s, name: %s\n", tag.ID, tag.Name)
	}
}

func TestListAllTags(t *testing.T) {
	setupClient(t)
	tags, err := client.ListAllTags("")
	assert.Nil(t, err)
	fmt.Printf("get: %d\n", len(tags))
	for _, tag := range tags {
		fmt.Printf("id: %s, name: %s\n", tag.ID, tag.Name)
	}
}

func TestListAllCategories(t *testing.T) {
	setupClient(t)
	categories, err := client.ListAllCategories()
	assert.Nil(t, err)
	for _, c := range categories {
		fmt.Printf("id: %s, name: %s\n", c.ID, c.Name)
	}
}
