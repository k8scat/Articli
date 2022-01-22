package juejin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListAllDrafts(t *testing.T) {
	setupClient(t)
	ids, err := client.ListAllDrafts()
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
