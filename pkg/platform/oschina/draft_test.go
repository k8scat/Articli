package oschina

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListAllTechnicalFields(t *testing.T) {
	setupClient(t)
	fields, err := client.ListAllTechnicalFields()
	assert.Nil(t, err)
	for i, f := range fields {
		fmt.Printf("%d. name: %s, val: %s\n", i+1, f.Name, f.ID)
	}
}

func TestListAllCategories(t *testing.T) {
	setupClient(t)
	categories, err := client.ListAllCategories()
	assert.Nil(t, err)
	for i, cate := range categories {
		fmt.Printf("%d. name: %s, val: %s\n", i+1, cate.Name, cate.ID)
	}
}
