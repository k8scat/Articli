package oschina

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestListCategories(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_OSCHINA_COOKIE"))
	if err != nil {
		t.Fail()
		return
	}

	_, err = client.ListCategories()
	assert.Nil(t, err)
}
