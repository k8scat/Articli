package juejin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListAllDrafts(t *testing.T) {
	setup(t)
	ids, err := c.ListAllDrafts()
	assert.Nil(t, err)
	for _, id := range ids {
		fmt.Printf("id: %s\n", id)
	}
}

func TestCreateDraft(t *testing.T) {
	// req := &CreateDraftRequest{
	// 	Title: "",
	// }
}
