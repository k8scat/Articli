package juejin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListAllDrafts(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_JUEJIN_COOKIE"))
	if err != nil {
		t.Fail()
		return
	}

	_, err = client.ListAllDrafts()
	assert.Nil(t, err)
}
