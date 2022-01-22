package oschina

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListAllTechnicalFields(t *testing.T) {
	setup(t)
	fields, err := c.ListAllTechnicalFields()
	assert.Nil(t, err)
	for i, f := range fields {
		fmt.Printf("%d. name: %s, val: %s\n", i+1, f.Name, f.ID)
	}
}

func TestListAllCategories(t *testing.T) {
	setup(t)
	categories, err := c.ListAllCategories()
	assert.Nil(t, err)
	for i, cate := range categories {
		fmt.Printf("%d. name: %s, val: %s\n", i+1, cate.Name, cate.ID)
	}
}
